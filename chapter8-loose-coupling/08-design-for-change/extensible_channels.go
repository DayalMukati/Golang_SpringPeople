package main

import "fmt"

// Abstraction
type Channel interface {
	Send(message string) error
}

// SMS Implementation
type SMSChannel struct{}

func (s SMSChannel) Send(message string) error {
	fmt.Println("SMS sent:", message)
	return nil
}

// Email Implementation
type EmailChannel struct{}

func (e EmailChannel) Send(message string) error {
	fmt.Println("Email sent:", message)
	return nil
}

// Notification Service
func notify(c Channel, msg string) {
	c.Send(msg)
}

func main() {
	// Use SMS
	notify(SMSChannel{}, "Payment of ₹1000 successful")

	// Switch to Email (no code change in Payment Service)
	notify(EmailChannel{}, "Payment of ₹1000 successful")
}