package scan

import (
	"context"
	"math"
	"sync"
)

type sizeWG struct {
	pool chan struct{}
	wg   sync.WaitGroup
}

func NewSizeWG(size int) *sizeWG {
	t := math.MaxInt16
	if size > 0 && size < t {
		t = size
	}
	return &sizeWG{
		pool: make(chan struct{}, t),
		wg:   sync.WaitGroup{},
	}

}
func (swg *sizeWG) addContext(ctx context.Context) {
	select {
	case <-ctx.Done():

	case swg.pool <- struct{}{}:
		break
	}
	swg.wg.Add(1)
}
func (swg *sizeWG) Add() {
	swg.addContext(context.Background())
}
func (swg *sizeWG) Done() {
	<-swg.pool
	swg.wg.Done()
}
func (swg *sizeWG) Wait() {
	swg.wg.Wait()
}
