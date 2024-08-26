package helpers

type ServerError struct {
	Message string
	Status  int
}

func (s *ServerError) ErrorMessage() any {
	return map[string]string{
		"error:": s.Message,
	}
}

// func (s *ServerError) Error() string {
// 	return s.Message
// }

func NewServerError(message string, status int) *ServerError {
	return &ServerError{Message: message, Status: status}
}
