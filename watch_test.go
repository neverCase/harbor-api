package harbor_api

import (
	"fmt"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	h := NewHarbor(fc.url, fc.admin, fc.password)
	opt := Option{
		Project:    "helix-saga",
		Repository: "go-all",
		Tag:        "latest",
	}

	// test TagOne
	tag, err := h.TagOne(fmt.Sprintf("%s/%s", opt.Project, opt.Repository), opt.Tag)
	if err != nil {
		t.Errorf("harbor.TagOne err:%v", err)
	}
	fmt.Println("tag-info:", tag)

	result, err := h.Watch(opt)
	if err != nil {
		t.Errorf("harbor.Watch err:%v", err)
	}
	for {
		select {
		case obj, isClose := <-result.ResultChan():
			if !isClose {
				fmt.Println("ResultChan close")
			}
			fmt.Println("obj:", obj)
		case <-time.After(time.Second * 3):
			fmt.Println("graceful stop after 3 seconds")
			return
		}
	}
}
