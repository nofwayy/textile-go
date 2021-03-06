package core

import (
	"errors"
	"fmt"
	"time"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/textileio/textile-go/repo"
)

type NotificationInfo struct {
	Id        string    `json:"id"`
	Date      time.Time `json:"date"`
	ActorId   string    `json:"actor_id"`
	Username  string    `json:"username,omitempty"`
	Subject   string    `json:"subject"`
	SubjectId string    `json:"subject_id"`
	BlockId   string    `json:"block_id,omitempty"`
	Target    string    `json:"target,omitempty"`
	Type      string    `json:"type"`
	Body      string    `json:"body"`
	Read      bool      `json:"read"`
}

// Notifications lists notifications
func (t *Textile) Notifications(offset string, limit int) []NotificationInfo {
	infos := make([]NotificationInfo, 0)
	for _, note := range t.datastore.Notifications().List(offset, limit) {
		infos = append(infos, t.NotificationInfo(note))
	}
	return infos
}

// NotificationInfo returns the notification info object
func (t *Textile) NotificationInfo(note repo.Notification) NotificationInfo {
	return NotificationInfo{
		Id:        note.Id,
		Date:      note.Date,
		ActorId:   note.ActorId,
		Username:  t.ContactUsername(note.ActorId),
		Subject:   note.Subject,
		SubjectId: note.SubjectId,
		BlockId:   note.BlockId,
		Target:    note.Target,
		Type:      note.Type.Description(),
		Body:      note.Body,
		Read:      note.Read,
	}
}

// CountUnreadNotifications counts unread notifications
func (t *Textile) CountUnreadNotifications() int {
	return t.datastore.Notifications().CountUnread()
}

// ReadNotification marks a notification as read
func (t *Textile) ReadNotification(id string) error {
	return t.datastore.Notifications().Read(id)
}

// ReadAllNotifications marks all notification as read
func (t *Textile) ReadAllNotifications() error {
	return t.datastore.Notifications().ReadAll()
}

// AcceptThreadInviteViaNotification uses an invite notification to accept an invite to a thread
func (t *Textile) AcceptThreadInviteViaNotification(id string) (mh.Multihash, error) {
	notification := t.datastore.Notifications().Get(id)
	if notification == nil {
		return nil, errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != repo.InviteReceivedNotification {
		return nil, errors.New(fmt.Sprintf("notification not type invite"))
	}

	hash, err := t.AcceptThreadInvite(notification.BlockId)
	if err != nil {
		return nil, err
	}

	return hash, t.datastore.Notifications().Delete(id)
}

// IgnoreThreadInviteViaNotification uses an invite notification to ignore an invite to a thread
func (t *Textile) IgnoreThreadInviteViaNotification(id string) error {
	notification := t.datastore.Notifications().Get(id)
	if notification == nil {
		return errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != repo.InviteReceivedNotification {
		return errors.New(fmt.Sprintf("notification not type invite"))
	}

	if err := t.IgnoreThreadInvite(notification.BlockId); err != nil {
		return err
	}

	return t.datastore.Notifications().Delete(id)
}
