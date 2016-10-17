package multi

import (
	"errors"
	"strings"

	"github.com/emersion/go-imap/backend"
)

type user struct {
	be *Backend
	childs map[string][]backend.User
	username, password string
}

func (u *user) ensureLoggedIn(reference string, i int) error {
	if u.childs[reference][i] != nil {
		return nil
	}

	be := u.be.childs[reference][i]

	if child, err := be.Login(u.username, u.password); err != nil {
		return err
	} else {
		u.childs[reference][i] = child
		return nil
	}
}

func (u *user) Username() string {
	return u.username
}

func (u *user) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	var mailboxes []backend.Mailbox
	for ref, childs := range u.childs {
		for i, child := range childs {
			if err := u.ensureLoggedIn(ref, i); err != nil {
				return nil, err
			}

			childMailboxes, err := child.ListMailboxes(subscribed)
			if err != nil {
				return nil, err
			}

			for _, m := range childMailboxes {
				if ref != "" {
					m = &mailbox{m, ref}
				}
				mailboxes = append(mailboxes, m)
			}
		}
	}

	return mailboxes, nil
}

func (u *user) GetMailbox(name string) (backend.Mailbox, error) {
	for ref, childs := range u.childs {
		if !strings.HasPrefix(name, ref) {
			continue
		}
		name := strings.TrimPrefix(name, ref)

		for i, child := range childs {
			if err := u.ensureLoggedIn(ref, i); err != nil {
				return nil, err
			}

			if mailbox, _ := child.GetMailbox(name); mailbox != nil {
				return mailbox, nil
			}
		}
	}

	return nil, errors.New("No such mailbox")
}

func (u *user) CreateMailbox(name string) error {
	child := u.childs[""][0]
	for ref, childs := range u.childs {
		if len(childs) == 0 {
			continue
		}

		if ref != "" && strings.HasPrefix(name, ref) {
			if err := u.ensureLoggedIn(ref, 0); err != nil {
				return nil
			}

			child = childs[0]
			name = strings.TrimPrefix(name, ref)
			break
		}
	}

	return child.CreateMailbox(name)
}

func (u *user) DeleteMailbox(name string) error {
	return nil // TODO
}

func (u *user) RenameMailbox(existingName, newName string) error {
		return nil // TODO
}

func (u *user) Logout() error {
	for _, childs := range u.childs {
		for _, child := range childs {
			if child == nil {
				continue
			}

			if err := child.Logout(); err != nil {
				return err
			}
		}
	}

	return nil
}
