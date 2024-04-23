package heritage

import (
	"errors"

	"github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
)

type Image struct {
	Layers []types.BlobInfo
	Name   string
	Err    error
}

func GetImage(image string) (*Image, error) {

	imageRef, err := GetImageReference(image)
	if err != nil {
		// idk
		return &Image{Err: err}, err
	}

	layers, err := GetImageLayers(imageRef)
	if err != nil {
		return &Image{Err: err}, err
	}

	return &Image{Layers: layers, Name: image}, nil
}

func ValidateHeritage(strict bool, child string, parents ...string) (bool, error) {

	childCh := make(chan Image)
	parentsCh := make(chan Image, len(parents))

	// These channels can close at the end of our function call
	defer close(childCh)
	defer close(parentsCh)

	// Get child
	go func() {
		img, _ := GetImage(child)
		childCh <- *img
	}()

	// Get Parents
	for _, parent := range parents {
		go func() {
			img, _ := GetImage(parent)
			parentsCh <- *img
		}()
	}

	childImg := <-childCh
	parentImgs := make([]Image, 0, len(parents))
	errs := make([]error, 0)

	for i := 0; i < len(parents); i++ {
		img := <-parentsCh
		if img.Err != nil {
			errs = append(errs, img.Err)
		}
		parentImgs = append(parentImgs, img)
	}

	if childImg.Err != nil {
		errs = append(errs, childImg.Err)
	}

	if len(errs) != 0 {
		return false, errors.Join(errs...)
	}

	return ValidateChildParentsImage(childImg, parentImgs...), nil

}

func ValidateChildParentsImage(child Image, parents ...Image) bool {
	parentsMap := map[int]Image{}
	for i, kv := range parents {
		if len(kv.Layers) > len(child.Layers) {
			// We cannot have a parent with a larger number of layers than childred
			logrus.Debugf("%s has more layers than child %s", kv.Name, child.Name)
			return false
		}
		parentsMap[i] = kv
	}

	state := map[int]bool{}
	for i := range parents {
		state[i] = false
	}

	for i := range child.Layers {
		for j := range parentsMap {
			if i >= len(parents[j].Layers) {
				// Just skip if we are out of bounds for this individual parent, state is already captured
				continue
			}

			if child.Layers[i].Digest == parentsMap[j].Layers[i].Digest {
				state[j] = true
			} else {
				state[j] = false
			}
		}
	}

	for _, v := range state {
		if !v {
			return false
		}
	}

	return true
}
