package harbor_api

import (
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/controller/artifact"
	"github.com/goharbor/harbor/src/controller/tag"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type fakeConfig struct {
	url      string
	admin    string
	password string
	timeout  int
}

var fc = fakeConfig{
	url:      "http://harbor.domain.com",
	admin:    "admin",
	password: "pwd",
	timeout:  10,
}

func TestNewHarbor(t *testing.T) {
	type args struct {
		url      string
		admin    string
		password string
	}
	tests := []struct {
		name string
		args args
		want HarborInterface
	}{
		//{
		//	name: "NewHarbor_case1",
		//	args: args{
		//		url:      fc.url,
		//		admin:    fc.admin,
		//		password: fc.password,
		//	},
		//	want: NewHarbor(fc.url, fc.admin, fc.password),
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHarbor(tt.args.url, tt.args.admin, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHarbor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_harbor_Login(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		//{
		//	name: "Login_case1",
		//	fields: fields{
		//		url:      fc.url,
		//		admin:    fc.admin,
		//		password: fc.password,
		//		timeout:  fc.timeout,
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
			}
			if err := h.Login(); (err != nil) != tt.wantErr {
				t.Errorf("harbor.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_harbor_Projects(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []models.Project
		wantErr bool
	}{
		{
			name: "Projects_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
			}
			res, err := h.Projects()
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("harbor.Projects() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			if err != nil {
				t.Errorf("harbor.Projects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			zaplogger.Sugar().Infow("project", "content", res)
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Projects() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
	zaplogger.Sugar().Infof("Test_harbor_Projects End \n\n\n")
}

func Test_harbor_Http(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
	}
	type args struct {
		method string
		url    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *http.Response
		wantErr bool
	}{
		//{
		//	name: "http_case1",
		//	fields: fields{
		//		url:      fc.url,
		//		admin:    fc.admin,
		//		password: fc.password,
		//		timeout:  fc.timeout,
		//	},
		//	args: args{
		//		method: "GET",
		//		url:    fc.url,
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
			}
			gotRes, err := h.Http(tt.args.method, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Http() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = gotRes
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Http() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_harbor_Repositories(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
	}
	type args struct {
		projectId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []models.RepoRecord
		wantErr bool
	}{
		{
			name: "Repositories_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
			}
			projects, err := h.Projects()
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Repositories() -> h.Projects()  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(projects) == 0 {
				t.Errorf("harbor.Repositories() -> h.Projects()  len(projects) == 0")
				return
			}
			gotRes, err := h.Repositories(projects[0].Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Repositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = gotRes
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Repositories() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
	zaplogger.Sugar().Infof("Test_harbor_Repositories End \n\n\n")
}

func Test_harbor_Tags(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
	}
	type args struct {
		imageName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []*tag.Tag
		wantErr bool
	}{
		{
			name: "Tags_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
			}
			projects, err := h.Projects()
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Projects()  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(projects) == 0 {
				//t.Errorf("harbor.Tags() -> h.Projects()  len(projects) == 0")
				return
			}
			repositories, err := h.Repositories(projects[0].Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Repositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(repositories) == 0 {
				//t.Errorf("harbor.Tags() -> h.Repositories()  len(projects) == 0")
				return
			}
			name := strings.Replace(repositories[0].Name, fmt.Sprintf("%s/", projects[0].Name), "", -1)
			gotRes, err := h.Tags(projects[0].Name, name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range gotRes {
				zaplogger.Sugar().Infof("tag image -> %s:%s", repositories[0].Name, v.Name)
			}
			_ = gotRes
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Tags() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
	zaplogger.Sugar().Infof("Test_harbor_Tags End \n\n\n")
}

func Test_harbor_Artifacts(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
		images   Images
	}
	type args struct {
		projectName    string
		repositoryName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []artifact.Artifact
		wantErr bool
	}{
		{
			name: "Artifacts_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
				images:   tt.fields.images,
			}
			projects, err := h.Projects()
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Projects()  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(projects) == 0 {
				//t.Errorf("harbor.Tags() -> h.Projects()  len(projects) == 0")
				return
			}
			repositories, err := h.Repositories(projects[0].Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Repositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(repositories) == 0 {
				//t.Errorf("harbor.Tags() -> h.Repositories()  len(projects) == 0")
				return
			}
			name := strings.Replace(repositories[0].Name, fmt.Sprintf("%s/", projects[0].Name), "", -1)
			gotRes, err := h.Artifacts(projects[0].Name, name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Artifacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range gotRes {
				for _, v2 := range v.Tags {
					zaplogger.Sugar().Infow("tags", "tag", v2.Tag)
				}
			}
			_ = gotRes
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Artifacts() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
	zaplogger.Sugar().Infof("Test_harbor_Artifacts End \n\n\n")
}

func Test_harbor_Watch(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
		images   Images
	}
	type args struct {
		opt Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    watch.Interface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
				images:   tt.fields.images,
			}
			got, err := h.Watch(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("harbor.Watch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_harbor_References(t *testing.T) {
	type fields struct {
		url      string
		admin    string
		password string
		timeout  int
		images   Images
	}
	type args struct {
		projectName    string
		repositoryName string
		digestOrTag    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes artifact.Artifact
		wantErr bool
	}{
		{
			name: "References_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:      tt.fields.url,
				admin:    tt.fields.admin,
				password: tt.fields.password,
				timeout:  tt.fields.timeout,
				images:   tt.fields.images,
			}
			projects, err := h.Projects()
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Projects()  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(projects) == 0 {
				//t.Errorf("harbor.Tags() -> h.Projects()  len(projects) == 0")
				return
			}
			repositories, err := h.Repositories(projects[0].Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() -> h.Repositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(repositories) == 0 {
				//t.Errorf("harbor.Tags() -> h.Repositories()  len(projects) == 0")
				return
			}
			name := strings.Replace(repositories[0].Name, fmt.Sprintf("%s/", projects[0].Name), "", -1)
			tags, err := h.Tags(projects[0].Name, name)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotRes, err := h.References(projects[0].Name, name, tags[0].Name)
			if err != nil {
				t.Errorf("harbor.References() error = %v", err)
				return
			}
			for _, v := range gotRes.Tags {
				zaplogger.Sugar().Infow("tags", "repository", repositories[0].Name, "tag", v.Tag)
			}
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.References() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
	zaplogger.Sugar().Infof("Test_harbor_References End \n\n\n")
}
