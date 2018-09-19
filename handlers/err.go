package handlers

import (
	"fmt"
)

// APIErr APIErr
type APIErr struct {
	Massage string
	Code    int
}

// NewError NewError
func NewError(code int, msg string) *APIErr {
	return &APIErr{
		Code:    code,
		Massage: msg,
	}
}

func (p *APIErr) Error() string {
	return fmt.Sprintf("%s:%d", p.Massage, p.Code)
}
