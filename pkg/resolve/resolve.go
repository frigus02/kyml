package resolve

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Resolve takes a Docker image reference (not an image ID) and attempts to
// resolve this to the image distribution digest. It does this using after
// another "docker inspect" and "docker manifest inspect". If any of these
// succeed, it returns the image reference minus tag plus resolved digest.
func Resolve(imageRef string) (string, error) {
	resolved, err := resolveWithDockerInspect(imageRef, execCmd)
	if resolved != "" || err != nil {
		return resolved, err
	}

	return resolveWithDockerManifestInspect(imageRef, execCmd)
}

type commandExecutor func(name string, arg ...string) (out []byte, err error)

func execCmd(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	return cmd.Output()
}

func resolveWithDockerInspect(imageRef string, execCmd commandExecutor) (string, error) {
	out, err := execCmd("docker", "inspect", "--format", "{{json .RepoDigests}}", imageRef)
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			stderr := string(err.Stderr)
			if strings.Contains(stderr, "No such object:") {
				return "", nil
			}

			return "", fmt.Errorf("docker inspect: %v (stderr: %s)", err, stderr)
		}

		return "", fmt.Errorf("docker inspect: %v", err)
	}

	var repoDigests []string
	if err = json.Unmarshal(out, &repoDigests); err != nil {
		return "", fmt.Errorf("json decode repo digests: %v", err)
	}

	digestPrefix := removeTagAndDigest(imageRef) + "@"
	for _, digest := range repoDigests {
		if strings.HasPrefix(digest, digestPrefix) {
			return digest, nil
		}
	}

	return "", nil
}

func resolveWithDockerManifestInspect(imageRef string, execCmd commandExecutor) (string, error) {
	out, err := execCmd("docker", "manifest", "inspect", "--verbose", imageRef)
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			stderr := string(err.Stderr)
			if strings.Contains(stderr, "no such manifest:") {
				return "", nil
			}

			return "", fmt.Errorf("docker manifest inspect: %v (stderr: %s)", err, stderr)
		}

		return "", fmt.Errorf("docker manifest inspect: %v", err)
	}

	outStr := string(out)
	if outStr[0] == '[' {
		// Manifest lists contain different manifests for different platforms.
		// To support this, kyml would need to pick the manifest for the
		// target Kubernetes platform.
		// https://docs.docker.com/registry/spec/manifest-v2-2/#manifest-list
		return "", fmt.Errorf("registry returned manifest list for %s, which is not supported yet", imageRef)
	}

	var result struct {
		Descriptor struct {
			Digest string `json:"digest"`
		} `json:"Descriptor"`
	}
	if err = json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("json decode result: %v", err)
	}

	repo := removeTagAndDigest(imageRef)
	return repo + "@" + result.Descriptor.Digest, nil
}

func removeTagAndDigest(imageRef string) string {
	indexAt := strings.Index(imageRef, "@")
	if indexAt > -1 {
		imageRef = imageRef[0:indexAt]
	}

	lastIndexColon := strings.LastIndex(imageRef, ":")
	if lastIndexColon > strings.LastIndex(imageRef, "/") {
		imageRef = imageRef[0:lastIndexColon]
	}

	return imageRef
}
