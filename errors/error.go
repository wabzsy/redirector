package errors

import (
	"fmt"
)

type Error struct {
	Tag    string
	Reason string
	Err    string
}

func (e *Error) Error() string {
	var result string

	if e.Tag != "" {
		result += fmt.Sprintf("[%s] ", e.Tag)
	}

	if e.Reason != "" {
		result += fmt.Sprintf("%s: ", e.Reason)
	}

	if e.Err != "" {
		result += fmt.Sprintf("%s", e.Err)
	}
	return result
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(e.Error()), nil
}

func (e *Error) UnmarshalJSON(bs []byte) error {
	e.Err = string(bs)
	return nil
}

func NewError(op, reason, text string) *Error {
	return &Error{
		Tag:    op,
		Reason: reason,
		Err:    text,
	}
}

func New(text string) *Error {
	return &Error{
		Err: text,
	}
}
