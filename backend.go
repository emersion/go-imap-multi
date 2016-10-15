package multi

import (
	"errors"

	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	childs map[string][]backend.Backend
}

func New() *Backend {
	return &Backend{}
}

func (be *Backend) Use(ref string, child backend.Backend) {
	be.childs[ref] = append(be.childs[ref], child)
}

func (be *Backend) Login(username, password string) (backend.User, error) {
	if len(be.childs[""]) == 0 {
		return nil, errors.New("No root backend available")
	}

	root := be.childs[""][0]

	uu, err := root.Login(username, password)
	if err != nil {
		return nil, err
	}

	u := &user{
		be: be,
		childs: make(map[string][]backend.User),
		username: username,
		password: password,
	}

	for ref, childs := range be.childs {
		u.childs[ref] = make([]backend.User, len(childs))
	}
	u.childs[""][0] = uu

	return u, nil
}
