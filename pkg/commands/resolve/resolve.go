package resolve

import (
	"fmt"
	"io"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/frigus02/kyml/pkg/k8syaml"
	"github.com/frigus02/kyml/pkg/resolve"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type resolveOptions struct{}

// NewCmdResolve creates a new resolve command.
func NewCmdResolve(in io.Reader, out io.Writer) *cobra.Command {
	var o resolveOptions

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve image tags to their distribution digest",
		Long: `Resolve image tags in Kubernetes YAML documents to their distribution digest. Data is read from stdin and printed to stdout.

This can be helpful if you tag the same image multiple times, e.g. because you build for every commit and use the commit sha as the Docker tag. Resolving the tag to the content digest before sending the manifests to Kubernetes makes sure your services only restart if the image actually changed.

In case an image is multi platform, it is resolved to the linux amd64 variant.`,
		Example: `  # Resolve image tags before deploying to cluster
  kyml cat feature/* | kyml resolve | kubectl apply -f -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(in, out, resolveImage)
		},
	}

	return cmd
}

// Validate validates resolve command.
func (o *resolveOptions) Validate(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("this command takes no positional arguments")
	}

	return nil
}

// Run runs resolve command.
func (o *resolveOptions) Run(in io.Reader, out io.Writer, resolveImage imageResolver) error {
	documents, err := cat.StreamDecodeOnly(in)
	if err != nil {
		return err
	}

	resolvedImageMap := make(map[string]string)
	for _, doc := range documents {
		if pathToPodSpec := getPathToPodSpec(doc.GroupVersionKind()); pathToPodSpec != nil {
			obj := doc.UnstructuredContent()

			pathToInitContainers := append(pathToPodSpec, "initContainers")
			if err := resolveImagesInContainers(obj, resolveImage, resolvedImageMap, pathToInitContainers...); err != nil {
				return err
			}

			pathToContainers := append(pathToPodSpec, "containers")
			if err := resolveImagesInContainers(obj, resolveImage, resolvedImageMap, pathToContainers...); err != nil {
				return err
			}
		}
	}

	return k8syaml.Encode(out, documents)
}

type imageResolver func(imageRef string) (resolveImage string, err error)

func resolveImage(imageRef string) (string, error) {
	return resolve.Resolve(imageRef)
}

func resolveImagesInContainers(
	obj map[string]interface{},
	resolveImage imageResolver,
	resolvedImageMap map[string]string,
	fields ...string,
) error {
	containers, found, err := unstructured.NestedSlice(obj, fields...)
	if !found || err != nil {
		return nil
	}

	for _, container := range containers {
		container, ok := container.(map[string]interface{})
		if !ok {
			return nil
		}

		image, ok := container["image"].(string)
		if !ok {
			return nil
		}

		resolvedImage, ok := resolvedImageMap[image]
		if !ok {
			resolvedImage, err = resolveImage(image)
			if err != nil {
				return err
			}

			if resolvedImage == "" {
				return fmt.Errorf("image %s not found", image)
			}

			resolvedImageMap[image] = resolvedImage
		}

		container["image"] = resolvedImage
	}

	return unstructured.SetNestedSlice(obj, containers, fields...)
}
