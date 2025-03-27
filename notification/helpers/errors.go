package helpers

import "log"

type ServerError struct {
	Message string
	Status  int
}

func (s *ServerError) ErrorMessage() any {
	return map[string]string{
		"error:": s.Message,
	}
}

func LogPanic(p any) {
	// log.Pnaic converts error to string and then panic . so it breaks panicRecoverer middleware
	log.Printf("logging pnic: %v\n", p)
	panic(p)
}

// func (s *ServerError) Error() string {
// 	return s.Message
// }

func NewServerError(message string, status int) *ServerError {
	return &ServerError{Message: message, Status: status}
}
