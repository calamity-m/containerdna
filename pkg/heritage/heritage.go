package heritage

import (
	"errors"

	"github.com/calamity-m/containerdna/pkg/containers"
	"github.com/sirupsen/logrus"
)

func ValidateHeritage(relaxed bool, child string, parents ...string) (bool, error) {

	childCh := make(chan containers.Image)
	parentsCh := make(chan containers.Image, len(parents))

	// These channels can close at the end of our function call
	defer close(childCh)
	defer close(parentsCh)

	// Get child
	go func() {
		img, _ := containers.GetImage(child)
		childCh <- *img
	}()

	// Get Parents
	for _, parent := range parents {
		go func() {
			img, _ := containers.GetImage(parent)
			parentsCh <- *img
		}()
	}

	errs := make([]error, 0)

	// Wait on our child and grab it.
	childImg := <-childCh
	if childImg.Err != nil {
		logrus.Debugf("Enountered error while fetching child layers: %v", childImg.Err)
		errs = append(errs, childImg.Err)
	}

	// Grab our parents, until we have a total of len(parents) parents.
	parentImgs := make([]containers.Image, 0, len(parents))
	for i := 0; i < len(parents); i++ {
		img := <-parentsCh
		if img.Err != nil {
			logrus.Debugf("Enountered error while fetching parent layers: %v", img.Err)
			errs = append(errs, img.Err)
		}
		parentImgs = append(parentImgs, img)
	}

	// Fail if we have any errors
	if len(errs) != 0 {
		return false, errors.Join(errs...)
	}

	return validateChildParentsImage(relaxed, childImg, parentImgs...), nil

}

func validateChildParentsImage(relaxed bool, child containers.Image, parents ...containers.Image) bool {
	if len(child.Layers) == 0 {
		return false
	}

	parentsMap := map[int]containers.Image{}
	for i, kv := range parents {
		if !relaxed && (len(kv.Layers) > len(child.Layers)) {
			// If our parent has more layers than child we can exist early
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

	// Look, I won't lie there is a better way to do this my but my brain was fried and I am nothin' but lazy.
	tracker := false
	for _, val := range state {
		if relaxed && val {
			return true
		} else if val {
			tracker = true
		} else {
			tracker = false
		}
	}

	return tracker
}
