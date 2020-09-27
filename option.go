package harbor_api

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Option struct {
	APIVersion string
	Kind       string

	Project     string
	Repository  string
	Tag         string
	sha256      string
	ExpiredTime int64
}

func (o Option) ImageName() string {
	return fmt.Sprintf("%s/%s:%s", o.Project, o.Repository, o.Tag)
}

func (o Option) GetObjectKind() schema.ObjectKind {
	return o
}

func (o Option) DeepCopyObject() runtime.Object {
	a := o
	return a
}

func (o Option) SetGroupVersionKind(kind schema.GroupVersionKind) {
	o.APIVersion = kind.Version
	o.Kind = kind.Kind
}

func (o Option) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(o.APIVersion, o.Kind)
}
