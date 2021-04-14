package harbor_api

import (
	"context"
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/controller/artifact"
	"github.com/goharbor/harbor/src/controller/tag"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HarborGetter interface {
	Harbor() HarborInterface
}

type HarborInterface interface {
	Http(method string, url string) (res *http.Response, err error)
	Login() error
	Projects() (res []models.Project, err error)
	Repositories(projectName string) (res []models.RepoRecord, err error)
	Artifacts(projectName string, repositoryName string) (res []artifact.Artifact, err error)
	Tags(projectName string, repositoryName string) (res []*tag.Tag, err error)
	References(projectName string, repositoryName string, digestOrTag string) (res artifact.Artifact, err error)
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
	Login        HarborUrlSuffix = "c/login"
	Projects     HarborUrlSuffix = "api/v2.0/projects?page=1&page_size=55&with_detail=true"
	Repositories HarborUrlSuffix = "api/v2.0/projects/%s/repositories?page=1&page_size=50"
	Artifacts    HarborUrlSuffix = "api/v2.0/projects/%s/repositories/%s/artifacts?with_tag=true&with_scan_overview=false&with_label=false&with_immutable_status=false&page_size=50&page=1"
	References   HarborUrlSuffix = "api/v2.0/projects/%s/repositories/%s/artifacts/%s?with_tag=true&with_scan_overview=false&with_label=false&with_immutable_status=false"
	TagOne       HarborUrlSuffix = "api/repositories/%s/tags/%s" // api/repositories/helix-saga/go-all/tags/latest
)

func (h *harbor) Http(method string, url string) (res *http.Response, err error) {
	zaplogger.Sugar().Infow("harbor-api http", "method", method, "url", url)
	var req *http.Request
	if req, err = http.NewRequest(method, url, nil); err != nil {
		zaplogger.Sugar().Error(err)
		return res, err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	if res, err = httpClient.Do(req); err != nil {
		zaplogger.Sugar().Error(err)
	}
	return res, err
}

func (h *harbor) Login() error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	u := fmt.Sprintf("%s/%v", h.url, Login)
	zaplogger.Sugar().Info("url:", u)
	data := url.Values{}
	data.Set("principal", h.admin)
	data.Set("password", h.password)
	body := ioutil.NopCloser(strings.NewReader(data.Encode())) // endode v:[body struce]
	req, err = http.NewRequest("POST", u, body)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value") // setting post head
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	//resp, err = httpClient.PostForm(u, data)
	resp, err = httpClient.Do(req)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	zaplogger.Sugar().Info(resp)
	_ = resp
	return err
}

func (h *harbor) Projects() (res []models.Project, err error) {
	var resp *http.Response
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, Projects)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
	}
	return res, nil
}

func (h *harbor) Repositories(projectName string) (res []models.RepoRecord, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(Repositories), projectName)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
	}
	return res, nil
}

func (h *harbor) Artifacts(projectName string, repositoryName string) (res []artifact.Artifact, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(Artifacts), projectName, repositoryName)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
	}
	return res, nil
}

func (h *harbor) Tags(projectName string, repositoryName string) (res []*tag.Tag, err error) {
	data, err := h.Artifacts(projectName, repositoryName)
	if err != nil {
		return nil, err
	}
	res = make([]*tag.Tag, 0)
	for _, v := range data {
		res = append(res, v.Tags...)
	}
	return res, nil
}

func (h *harbor) References(projectName string, repositoryName string, digestOrTag string) (res artifact.Artifact, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(References), projectName, repositoryName, digestOrTag)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
	}
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
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			zaplogger.Sugar().Error(err)
			return res, err
		}
	}
	//zaplogger.Sugar().Info(res)
	return res, nil
}

func (h *harbor) Watch(opt Option) (watch.Interface, error) {
	image, err := h.images.Image(opt)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	return image.Watch(), nil
}
