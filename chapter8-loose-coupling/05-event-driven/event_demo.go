package main

import "fmt"

type Event struct {
	Name string
	Data string
}

// Publisher: Payment
func payment(events chan Event) {
	fmt.Println("Payment: publishing TransactionCreated")
	events <- Event{Name: "TransactionCreated", Data: "TXN123"}
}

// Subscriber: Fraud
func fraud(events chan Event) {
	for e := range events {
		if e.Name == "TransactionCreated" {
			fmt.Println("Fraud: checking transaction", e.Data)
		}
	}
}

func main() {
	events := make(chan Event, 1)

	go fraud(events) // subscriber
	payment(events)  // publisher

	select {} // keep running
}