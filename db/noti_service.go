package db

import (
	"log"
)

// NotificationService interface wraps around external noti functionality.
type NotificationService interface {
	NewNoti(t Todo) error
}

// StubbedNotificationService is a concrete mock of the external NotificationService.
type StubbedNotificationService struct{}

// NewNotificationService initialises the NotificationService.
func NewNotificationService() NotificationService {
	return &StubbedNotificationService{}
}

// NewNoti creates a new notification and sends it to the notification servivce for notifying.
func (sns *StubbedNotificationService) NewNoti(t Todo) error {
	log.Printf("STUBBED NOTI SERVICE: todo %s noti: %v", t.ID, t)
	return nil
}
