package handlers

import "fmt"

var (
	// JSONSuccess is a map with a message key and a success message as value
	JSONSuccess = map[string]string{"message": "OK"}
	// JSONNotFound is a map with a message key and a not found message as value
	JSONNotFound = map[string]string{"message": "Not found"}
)

// JSONMessageError returns a map with a message key and the given error as value
func JSONMessageError(err error) map[string]string {
	return map[string]string{"message": fmt.Sprintf("%s", err)}
}

// JSONMessageText returns a map with a message key and the given text as value
func JSONMessageText(text string) map[string]string {
	return map[string]string{"message": text}
}
