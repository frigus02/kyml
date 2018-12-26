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

	var digest string

	if string(out)[0] == '[' {
		// The registry returned a multi platform manifest list. In this case
		// we always return linux amd64.
		// See: https://blog.docker.com/2017/09/docker-official-images-now-multi-platform/
		// See: https://docs.docker.com/registry/spec/manifest-v2-2/#manifest-list

		var result []struct {
			Descriptor struct {
				Digest   string `json:"digest"`
				Platform struct {
					Architecture string `json:"architecture"`
					OS           string `json:"os"`
				} `json:"platform"`
			} `json:"Descriptor"`
		}
		if err = json.Unmarshal(out, &result); err != nil {
			return "", fmt.Errorf("json decode manifest list: %v", err)
		}

		for _, platformImage := range result {
			if platformImage.Descriptor.Platform.Architecture == "amd64" &&
				platformImage.Descriptor.Platform.OS == "linux" {
				digest = platformImage.Descriptor.Digest
				break
			}
		}

		if digest == "" {
			return "", nil
		}
	} else {
		var result struct {
			Descriptor struct {
				Digest string `json:"digest"`
			} `json:"Descriptor"`
		}
		if err = json.Unmarshal(out, &result); err != nil {
			return "", fmt.Errorf("json decode manifest: %v", err)
		}

		digest = result.Descriptor.Digest
	}

	return removeTagAndDigest(imageRef) + "@" + digest, nil
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
