package notification

import (
	"context"
	"log"
)

type NotificationService interface {
	SendPushNotification(ctx context.Context, agentID string, message string) error
	SendSMS(ctx context.Context, phoneNumber string, message string) error
}

type MockService struct{}

func (m *MockService) SendPushNotification(ctx context.Context, agentID string, message string) error {
	log.Printf("Sending push notification to agent %s: %s", agentID, message)
	return nil
}

func (m *MockService) SendSMS(ctx context.Context, phoneNumber string, message string) error {
	log.Printf("Sending SMS to %s: %s", phoneNumber, message)
	return nil
}
