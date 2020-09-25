package harbor_api

type Option struct {
	Project    string
	Repository string
	Tag        string

	ExpiredTime int64
}
