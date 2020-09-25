package harbor_api

import (
	"k8s.io/apimachinery/pkg/watch"
)

// NewWatch returns a k8s.io watch.Interface
func NewWatch(h *harbor, opt Option) watch.Interface {
	return &harborWatch{
		h:      h,
		opt:    opt,
		result: make(chan watch.Event, 1000),
	}
}

// harborWatch implements the k8s.io/apimachinery/pkg/watch.Interface
type harborWatch struct {
	h      *harbor
	opt    Option
	result chan watch.Event
}

func (w *harborWatch) Stop() {

}

func (w *harborWatch) ResultChan() <-chan watch.Event {
	return w.result
}
