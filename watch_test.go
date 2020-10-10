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

func TestGetHashFromDockerImageId(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetHashFromDockerImageId_1",
			args: args{
				s: "docker-pullable://harbor.domain.com/helix-saga/go-all@sha256:27d6aa8f9d040c5e85c61a093ad2dc769e57440e8240c3294f47093e97d96c9a",
			},
			want: "sha256:27d6aa8f9d040c5e85c61a093ad2dc769e57440e8240c3294f47093e97d96c9a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHashFromDockerImageId(tt.args.s); got != tt.want {
				t.Errorf("GetHashFromDockerImageId() = %v, want %v", got, tt.want)
			}
		})
	}
}
