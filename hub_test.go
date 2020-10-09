package harbor_api

import (
	"reflect"
	"testing"
)

var fakeConfig2 = []Config{
	{
		Url:      "https://www.domain1.com",
		Admin:    "admin",
		Password: "pwd",
	},
	{
		Url:      "http://111.222.333.11:8863",
		Admin:    "admin",
		Password: "pwd",
	},
}

func TestNewHub(t *testing.T) {
	type args struct {
		c []Config
	}
	tests := []struct {
		name string
		args args
		want HubInterface
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHub(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hub_List(t *testing.T) {
	type fields struct {
		harbors map[string]HarborInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Test_hub_List_1",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHub(fakeConfig2)
			if got := h.List(); len(got) != len(fakeConfig2) {
				t.Errorf("hub.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hub_Get(t *testing.T) {
	type fields struct {
		harbors map[string]HarborInterface
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    HarborInterface
		wantErr bool
	}{
		{
			name: "Test_hub_Get_1",
			args: args{
				url: "https://www.domain1.com",
			},
			wantErr: false,
		},
		{
			name: "Test_hub_Get_2",
			args: args{
				url: "http://www.domain1.com",
			},
			wantErr: false,
		},
		{
			name: "Test_hub_Get_3",
			args: args{
				url: "http://www.domain13333.com",
			},
			wantErr: true,
		},
		{
			name: "Test_hub_Get_4",
			args: args{
				url: "https://111.222.333.11:8863",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHub(fakeConfig2)
			got, err := h.Get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("hub.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("hub.Get() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestConvertUrlToHttp(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestConvertUrlToHttp_1",
			args: args{
				in: "http://harbor.domain.com",
			},
			want: "http://harbor.domain.com",
		},
		{
			name: "TestConvertUrlToHttp_2",
			args: args{
				in: "https://harbor.domain222.com",
			},
			want: "http://harbor.domain222.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertUrlToHttp(tt.args.in); got != tt.want {
				t.Errorf("ConvertUrlToHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertUrlToHttps(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestConvertUrlToHttps_1",
			args: args{
				in: "http://harbor.domain.com",
			},
			want: "https://harbor.domain.com",
		},
		{
			name: "TestConvertUrlToHttps_2",
			args: args{
				in: "https://harbor.domain222.com",
			},
			want: "https://harbor.domain222.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertUrlToHttps(tt.args.in); got != tt.want {
				t.Errorf("ConvertUrlToHttps() = %v, want %v", got, tt.want)
			}
		})
	}
}
