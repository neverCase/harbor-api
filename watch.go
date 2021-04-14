package harbor_api

import (
	"context"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/goharbor/harbor/src/controller/artifact"
	"k8s.io/apimachinery/pkg/watch"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	maxQueuedEvents  = 1000
	maxRemovedChan   = 1000
	loopTickTimeInMs = 1500
)

type RequestHandler func(projectName string, repositoryName string, digestOrTag string) (res artifact.Artifact, err error)

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
	i := &images{
		images:      make(map[string]Image, 0),
		removedChan: make(chan string, maxRemovedChan),
		handler:     handler,
		ctx:         subCtx,
		cancel:      cancel,
	}
	go i.Loop()
	return i
}

func (images *images) Loop() {
	for {
		select {
		case name := <-images.removedChan:
			images.mu.Lock()
			t, ok := images.images[name]
			if !ok {
				images.mu.Unlock()
				continue
			}
			delete(images.images, name)
			images.mu.Unlock()
			t.Shutdown()
		case <-images.ctx.Done():
			for _, v := range images.images {
				v.Shutdown()
			}
			return
		}
	}
}

func (images *images) Image(opt Option) (Image, error) {
	images.mu.Lock()
	defer images.mu.Unlock()
	if t, ok := images.images[opt.ImageName()]; ok {
		return t, nil
	} else {
		i, err := NewImage(images.ctx, opt, images.removedChan, images.handler)
		if err != nil {
			return nil, err
		}
		images.images[opt.ImageName()] = i
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
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			res, err := i.handler(i.opt.Project, i.opt.Repository, i.opt.Tag)
			if err != nil {
				// todo check whether the error was like `{"code":404,"message":"resource: xxxxx not found"}`
				zaplogger.Sugar().Error(err)
				select {
				case removedChan <- i.opt.ImageName():
					zaplogger.Sugar().Infof("Loop send removedChan:%s success", i.opt.ImageName())
				case <-time.After(time.Second * 1):
					zaplogger.Sugar().Infof("Loop send removedChan:%s timout", i.opt.ImageName())
				}
				continue
			}
			if i.opt.Sha256 == "" || res.Digest != i.opt.Sha256 {
				i.opt.Sha256 = res.Digest
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
func GetHashFromDockerImageId(s string) string {
	t := strings.Split(s, "@")
	if len(t) != 2 {
		return ""
	}
	return t[1]
}
