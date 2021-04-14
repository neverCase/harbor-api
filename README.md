# harbor-api

This is a simple and easy way for you to interact with Harbor.

### notice
* The structures inside the harbor.go were all imported from [goharbor/harbor](https://github.com/goharbor/harbor/tree/master/src)

## Features
*	Projects() (res []models.Project, err error)
*	Repositories(projectName string) (res []models.RepoRecord, err error)
*	Artifacts(projectName string, repositoryName string) (res []artifact.Artifact, err error)
*	Tags(projectName string, repositoryName string) (res []*tag.Tag, err error)
*	References(projectName string, repositoryName string, digestOrTag string) (res artifact.Artifact, err error)
*	Watch(opt Option) (watch.Interface, error), watch implements the k8s.io/apimachinery/pkg/watch.Interface, and it watches and compares the image's sha256 by the specific tag

## Usage
```
h := NewHarbor(url, admin, password)
res, err := h.Projects()
...
```

## Todo
* add api pagination in the future