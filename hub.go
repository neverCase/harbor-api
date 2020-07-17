package harbor_api

type Config struct {
	Url      string `json:"url"`
	Admin    string `json:"admin"`
	Password string `json:"password"`
}

type HubGetter interface {
	Hub() HarborInterface
}

type HubInterface interface {
	Get(url string) HarborInterface
}

type hub struct {
	harbors map[string]HarborInterface
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

func (h *hub) Get(url string) HarborInterface {
	return h.harbors[url]
}