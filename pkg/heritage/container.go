package heritage

import (
	"context"
	"fmt"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
)

func GetImageReference(name string) (types.ImageReference, error) {
	// docker-daemon:test:test
	// docker-daemon:alpine:latest
	// docker://alpine
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		logrus.Errorf("Failed to get image reference for %s, err: %v", name, err)
		return nil, err
	}

	return ref, err
}

func GetImageLayers(ref types.ImageReference) ([]types.BlobInfo, error) {
	ctx := context.Background()

	img, err := ref.NewImage(ctx, nil)
	defer func(img types.ImageCloser) {
		err := img.Close()
		if err != nil {
			panic(err)
		}
	}(img)

	if err != nil {
		logrus.Errorf("Failed to parse into Image with given image ref: %v", ref)
		return nil, err
	}

	info := img.LayerInfos()
	return info, nil
}

func GetImageManifest(ctx context.Context, closer types.ImageCloser) []byte {
	manifest, s, err := closer.Manifest(ctx)
	if err != nil {
		panic(err)
	}
	logrus.Debugf(fmt.Sprintf("Retrun from Manifest get: %s", s))
	return manifest
}
