package main

// Tightly Coupled Example
func ProcessPayment(amount int) {
	fraudCheck(amount)    // Direct dependency
	sendNotification()    // Direct dependency
}

func fraudCheck(amount int) {
	// fraud logic here
}

func sendNotification() {
	// notification logic here
}

func main() {
	ProcessPayment(1000)
}