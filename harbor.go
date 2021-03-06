package harbor_api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"
)

type HarborGetter interface {
	Harbor() HarborInterface
}

type HarborInterface interface {
	Http(method string, url string) (res *http.Response, err error)
	Login() error
	Projects() (res []Project, err error)
	Repositories(projectId int) (res []RepoRecord, err error)
	Tags(imageName string) (res []TagDetail, err error)
	TagOne(imageName, tagName string) (res TagDetail, err error)
	Watch(opt Option) (watch.Interface, error)
}

func NewHarbor(url, admin, password string) HarborInterface {
	h := &harbor{
		url:      url,
		admin:    admin,
		password: password,
		timeout:  10,
	}
	h.images = NewImages(context.Background(), h.TagOne)
	return h
}

type harbor struct {
	url      string
	admin    string
	password string
	timeout  int

	images Images
}

type HarborUrlSuffix string

const (
	Login        HarborUrlSuffix = "login"
	SystemInfo   HarborUrlSuffix = "api/systeminfo"
	Projects     HarborUrlSuffix = "api/projects"                         // api/projects?page=1&page_size=15
	Repositories HarborUrlSuffix = "api/repositories?&project_id=%d"      // api/repositories?page=1&page_size=15&project_id=2
	Tags         HarborUrlSuffix = "api/repositories/%s/tags?detail=true" // api/repositories/helix-saga/redis-slave/tags?detail=true
	TagOne       HarborUrlSuffix = "api/repositories/%s/tags/%s"          // api/repositories/helix-saga/go-all/tags/latest
)

func (h *harbor) Http(method string, url string) (res *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, nil); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	if res, err = httpClient.Do(req); err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (h *harbor) Login() error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/%v", h.url, Login), nil)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	resp, err = httpClient.Do(req)
	if err != nil {
		klog.V(2).Info(err)
	}
	_ = resp
	return err
}

func (h *harbor) Projects() (res []Project, err error) {
	var resp *http.Response
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, Projects)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	klog.Info(res)
	return res, nil
}

func (h *harbor) Repositories(projectId int) (res []RepoRecord, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(Repositories), projectId)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	klog.Info(res)
	return res, nil
}

func (h *harbor) Tags(imageName string) (res []TagDetail, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(Tags), imageName)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	klog.Info(res)
	return res, nil
}

func (h *harbor) TagOne(imageName, tagName string) (res TagDetail, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(TagOne), imageName, tagName)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	//klog.Info(res)
	return res, nil
}

func (h *harbor) Watch(opt Option) (watch.Interface, error) {
	image, err := h.images.Image(opt)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return image.Watch(), nil
}
