package enum

import "fmt"

const (
	InviteStatusSent     = "sent"     // Invite is sent
	InviteStatusAccepted = "accepted" // Invite is accepted
	InviteStatusRejected = "rejected" // Invite is rejected
)

var InviteStatuses = []string{
	InviteStatusSent,
	InviteStatusAccepted,
	InviteStatusRejected,
}

var ErrorInviteStatusNotSupported = fmt.Errorf("invite status not supported, must be one of: %v", GetAllInviteStatuses())

func ParseInviteStatus(status string) (string, error) {
	for _, s := range InviteStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorInviteStatusNotSupported)
}

func GetAllInviteStatuses() []string {
	return InviteStatuses
}
