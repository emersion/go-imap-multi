package multi

import (
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type mailbox struct {
	backend.Mailbox
	reference string
}

func (m *mailbox) Name() string {
	return strings.TrimPrefix(m.Mailbox.Name(), m.reference)
}

func (m *mailbox) Info() (*imap.MailboxInfo, error) {
	info, err := m.Mailbox.Info()
	if err != nil {
		return nil, err
	}

	info.Name = m.Name()
	return info, nil
}

func (m *mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status, err := m.Mailbox.Status(items)
	if err != nil {
		return nil, err
	}

	status.Name = m.Name()
	return status, nil
}
