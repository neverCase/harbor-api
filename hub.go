package harbor_api

import (
	"fmt"
	"regexp"
)

const (
	ErrorHarborUrlWasNotExisted = "error: the harbor url:%s was not existed"

	HttpPrefix  = "http://"
	HttpsPrefix = "https://"
)

type HubGetter interface {
	HarborHub() HubInterface
}

type HubInterface interface {
	List() []string
	Get(url string) (HarborInterface, error)
}

type hub struct {
	harbors map[string]HarborInterface
}

type Config struct {
	Url      string `json:"url"`
	Admin    string `json:"admin"`
	Password string `json:"password"`
}

func NewHub(c []Config) HubInterface {
	h := &hub{
		harbors: make(map[string]HarborInterface, 0),
	}
	for _, v := range c {
		h.harbors[v.Url] = NewHarbor(v.Url, v.Admin, v.Password)
	}
	return h
}

func (h *hub) List() []string {
	res := make([]string, 0)
	for k := range h.harbors {
		res = append(res, k)
	}
	return res
}

func (h *hub) Get(url string) (HarborInterface, error) {
	s := ConvertUrlToHttp(url)
	if t, ok := h.harbors[s]; ok {
		return t, nil
	}
	s = ConvertUrlToHttps(url)
	if t, ok := h.harbors[s]; ok {
		return t, nil
	}
	return nil, fmt.Errorf(ErrorHarborUrlWasNotExisted, url)
}

func ConvertUrlToHttp(in string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s|%s`, HttpPrefix, HttpsPrefix))
	out := re.ReplaceAll([]byte(in), []byte(HttpPrefix))
	return string(out)
}

func ConvertUrlToHttps(in string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s|%s`, HttpPrefix, HttpsPrefix))
	out := re.ReplaceAll([]byte(in), []byte(HttpsPrefix))
	return string(out)
}
