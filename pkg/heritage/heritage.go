package heritage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
)

func Playground() {

	var s string

	s = "hh"

	fmt.Println(s)

	ch := make(chan []string, 3)

	defer close(ch)
	fmt.Print(ch)

	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	for elem := range queue {
		fmt.Println(elem)
		logrus.Info(elem)
	}

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func(id int) {
			fmt.Printf("start %d\n, also: %d", id, i)
			fmt.Printf("done %d\n", id)

			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println("done wg")

	fmt.Println(&wg)

}

type Image struct {
	Layers []types.BlobInfo
	Name   string
	Err    error
}

func ValidateWithChannelsNoWg(strict bool, child string, parentNames ...string) (bool, error) {

	errCh := make(chan error)
	resultCh := make(chan Image)

	for _, name := range parentNames {
		go func(name string, ch chan Image, errCh chan error) {
			fmt.Printf("I am starting: %s\n", name)
			parentRef, err := GetImageReference(name)
			if err != nil {
				errCh <- err
				resultCh <- Image{}
				return
			}
			lys, err := GetImageLayers(parentRef)
			if err != nil {
				errCh <- err
				return
			}
			parent := &Image{
				Layers: lys,
				Name:   name,
			}
			fmt.Printf("I got my result, tryna add it: %s\n", name)
			resultCh <- *parent
			errCh <- nil
			fmt.Printf("I done, finished my result: %s\n", name)
		}(name, resultCh, errCh)

	}

	errs := make([]error, 0, len(parentNames))
	parents := make([]Image, 0, len(parentNames))
	for i := 0; i < len(parentNames); i++ {
		parent := <-resultCh
		err := <-errCh

		parents = append(parents, parent)
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		logrus.Errorf("Encountered %d errors", len(errs))
		return false, errors.Join(errs...)
	}

	close(errCh)
	close(resultCh)

	fmt.Println(parents)

	return false, nil
}

func ValidateWithChannels(strict bool, child string, parent ...string) (bool, error) {
	// Define our waitgroup for fetching children and parents
	var wg sync.WaitGroup

	parents := make([]Image, len(parent))
	errCh := make(chan error)
	resultCh := make(chan Image)

	for _, name := range parent {
		wg.Add(1)

		go func(name string, ch chan Image, errCh chan error) {
			defer wg.Done()

			parentRef, err := GetImageReference(name)
			if err != nil {
				errCh <- err
				return
			}

			lys, err := GetImageLayers(parentRef)
			if err != nil {
				errCh <- err
				return
			}

			parent := &Image{
				Layers: lys,
				Name:   name,
			}

			resultCh <- *parent

		}(name, resultCh, errCh)

	}

	wg.Wait()

	close(errCh)
	close(resultCh)

	errStack := make([]error, 0)
	for err := range errCh {
		if err != nil {
			errStack = append(errStack, err)
		}
	}

	if len(errStack) != 0 {
		return false, errors.Join(errStack...)
	}

	for parent := range resultCh {
		parents = append(parents, parent)
	}

	return false, nil
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

func GetImageAsyncChannels(strict bool, child string, parents ...string) (bool, error) {

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

func GetImageAsyncWG(wg *sync.WaitGroup, resultCh chan Image, name string) {
	wg.Add(1)
	defer wg.Done()

	image, err := GetImage(name)

	if err != nil {
		resultCh <- Image{Err: err}
		return
	}

	resultCh <- *image
}

func UseGoFunc(strict bool, child string, parents ...string) (bool, error) {

	childCh := make(chan Image)
	parentsCh := make(chan Image, len(parents))

	var wg sync.WaitGroup

	// Get child
	go GetImageAsyncWG(&wg, childCh, child)

	// Get parents
	for _, parent := range parents {
		go GetImageAsyncWG(&wg, parentsCh, parent)
	}

	wg.Wait()
	close(childCh)
	close(parentsCh)

	// idk about errors but

	childImg := <-childCh
	parentImgs := make([]Image, 0, len(parents))
	errs := make([]error, 0)
	for img := range parentsCh {
		if img.Err != nil {
			// i unhappy
			fmt.Printf("Failed %v\n", img.Err)
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

func ValidateWithContextExitEarly(strict bool, child string, parent ...string) (bool, error) {

	// Define our waitgroup for fetching children and parents
	var wg sync.WaitGroup

	// Add weighting for our parents and child
	wg.Add(len(parent) + 1)

	parents := make([]Image, len(parent))
	// Mutex for parents working
	var parentMu sync.Mutex

	ctx, _ := context.WithCancel(context.Background())

	for _, name := range parent {

		go func() {
			defer wg.Done()

			parent, err := GetImage(name)
			if err != nil {
				// idk do something
			}

			parentMu.Lock()
			parents = append(parents, *parent)
			parentMu.Unlock()
		}()

	}

	wg.Wait()

	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return false, nil
}

func ValidateChildParents(child []types.BlobInfo, parent ...[]types.BlobInfo) bool {
	parents := map[int][]types.BlobInfo{}
	for i, kv := range parent {
		if len(kv) > len(child) {
			// We cannot have a parent with a larger number of layers than childred
			logrus.Debug("One of the supplied parents has more layers than children")
			return false
		}
		parents[i] = kv
	}

	state := map[int]bool{}
	for i := range parent {
		state[i] = false
	}

	for i := range child {
		for j := range parents {
			if i >= len(parents[j]) {
				// Just skip if we are out of bounds for this individual parent, state is already captured
				continue
			}

			if child[i].Digest == parents[j][i].Digest {
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
