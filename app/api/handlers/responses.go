package handlers

// Message is a generic message struct for returning JSON messages
type Message struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// MessageSuccess is a default success message
var MessageSuccess = Message{Message: "OK"}

// MessageNotFound is a default not found message
var MessageNotFound = Message{Message: "Not found"}
