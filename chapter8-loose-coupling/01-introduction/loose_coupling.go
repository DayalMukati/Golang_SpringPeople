package main

import "fmt"

// Loosely Coupled Example
type FraudChecker interface {
	Check(amount int) bool
}

type Notifier interface {
	Send(message string)
}

func ProcessPayment(amount int, fc FraudChecker, n Notifier) {
	if fc.Check(amount) {
		n.Send("Payment processed successfully")
	}
}

// Implementations
type SimpleFraudChecker struct{}

func (s SimpleFraudChecker) Check(amount int) bool {
	return amount < 100000 // approve if less than 100k
}

type SimpleNotifier struct{}

func (s SimpleNotifier) Send(message string) {
	fmt.Println("Notification:", message)
}

func main() {
	fc := SimpleFraudChecker{}
	n := SimpleNotifier{}
	ProcessPayment(5000, fc, n)
}