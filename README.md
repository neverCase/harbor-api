# harbor-api

This is a simaple and easy way for you to use the api of the harbor in go.

### notice
* The structures inside the harbor_types.go were all copied from github.com/goharbor/harbor/src/common/models.
* Because there are some errors when we import [goharbor/harbor](https://github.com/goharbor/harbor/tree/master/src) in go.mod. 
* Otherwise, we would import the types from [github.com/goharbor/harbor/src/common/models](https://github.com/goharbor/harbor/tree/master/src/common/models) instead of copying.


## Features
*	Projects() (res []Project, err error)
*   Repositories(projectId int) (res []RepoRecord, err error)
*   Tags(imageName string) (res []TagDetail, err error)

## Usage
```
h := NewHarbor(url, admin, password)
res, err := h.Projects()
...
```
