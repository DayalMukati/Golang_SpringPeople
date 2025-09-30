package main

import "fmt"

// Abstraction
type Notifier interface {
	Send(message string) error
}

// Twilio implementation
type TwilioNotifier struct{}

func (t TwilioNotifier) Send(message string) error {
	fmt.Println("Twilio SMS sent:", message)
	return nil
}

// AWS SNS implementation
type SNSNotifier struct{}

func (s SNSNotifier) Send(message string) error {
	fmt.Println("AWS SNS sent:", message)
	return nil
}

// Payment service depends only on Notifier interface
func processPayment(n Notifier) {
	fmt.Println("Payment processed successfully")
	n.Send("Payment successful notification")
}

func main() {
	// Choose implementation
	twilio := TwilioNotifier{}
	aws := SNSNotifier{}

	processPayment(twilio) // uses Twilio
	processPayment(aws)    // uses AWS SNS
}