package harbor_api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

const (
	maxQueuedEvents  = 1000
	maxRemovedChan   = 1000
	loopTickTimeInMs = 500
)

type RequestHandler func(imageName, tag string) (res TagDetail, err error)

type Images interface {
	Image(opt Option) (Image, error)
}

type images struct {
	mu sync.Mutex

	images      map[string]Image
	removedChan chan string

	handler RequestHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewImages(ctx context.Context, handler RequestHandler) Images {
	subCtx, cancel := context.WithCancel(ctx)
	is := &images{
		images:      make(map[string]Image, 0),
		removedChan: make(chan string, maxRemovedChan),
		handler:     handler,
		ctx:         subCtx,
		cancel:      cancel,
	}
	go is.Loop()
	return is
}

func (is *images) Loop() {
	for {
		select {
		case name := <-is.removedChan:
			is.mu.Lock()
			t, ok := is.images[name]
			if !ok {
				is.mu.Unlock()
				continue
			}
			delete(is.images, name)
			is.mu.Unlock()
			t.Shutdown()
		case <-is.ctx.Done():
			for _, v := range is.images {
				v.Shutdown()
			}
			return
		}
	}
}

func (is *images) Image(opt Option) (Image, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	if t, ok := is.images[opt.ImageName()]; ok {
		return t, nil
	} else {
		i, err := NewImage(is.ctx, opt, is.removedChan, is.handler)
		if err != nil {
			return nil, err
		}
		is.images[opt.ImageName()] = i
		return i, nil
	}
}

type Image interface {
	Watch() watch.Interface
	Shutdown()
}

type image struct {
	once sync.Once

	opt          Option
	handler      RequestHandler
	broadcasters *watch.Broadcaster

	ctx    context.Context
	cancel context.CancelFunc
}

func NewImage(ctx context.Context, opt Option, removedChan chan<- string, handler RequestHandler) (Image, error) {
	// todo: is it necessary to check whether the harbor image was existed?
	subCtx, cancel := context.WithCancel(ctx)
	i := &image{
		opt:          opt,
		handler:      handler,
		broadcasters: watch.NewBroadcaster(maxQueuedEvents, watch.DropIfChannelFull),
		ctx:          subCtx,
		cancel:       cancel,
	}
	go i.Loop(removedChan)
	return i, nil
}

func (i *image) Loop(removedChan chan<- string) {
	defer i.Shutdown()
	tick := time.NewTicker(time.Millisecond * loopTickTimeInMs)
	defer tick.Stop()
	for {
		select {
		case <-i.ctx.Done():
			return
		case <-tick.C:
			res, err := i.handler(fmt.Sprintf("%s/%s", i.opt.Project, i.opt.Repository), i.opt.Tag)
			if err != nil {
				// todo check whether the error was like `{"code":404,"message":"resource: xxxxx not found"}`
				klog.V(2).Info(err)
				select {
				case removedChan <- i.opt.ImageName():
					klog.Infof("Loop send removedChan:%s success", i.opt.ImageName())
				case <-time.After(time.Second * 1):
					klog.Infof("Loop send removedChan:%s timout", i.opt.ImageName())
				}
				continue
			}
			if i.opt.sha256 == "" || res.Digest != i.opt.sha256 {
				i.opt.sha256 = res.Digest
				i.broadcasters.Action(watch.Modified, i.opt)
			}
		}
	}
}

func (i *image) Watch() watch.Interface {
	return i.broadcasters.Watch()
}

func (i *image) Shutdown() {
	i.once.Do(func() {
		i.cancel()
		i.broadcasters.Shutdown()
	})
}

// image: harbor.domain.com/helix-saga/go-all:latest
// imageID: docker-pullable://harbor.domain.com/helix-saga/go-all@sha256:27d6aa8f9d040c5e85c61a093ad2dc769e57440e8240c3294f47093e97d96c9a
func ConvertDockerImageIdToHarbor(s string) {

}
