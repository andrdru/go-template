package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	// Message .
	// init with NewMessage()
	Message struct {
		Data          any
		ErrorCode     int
		ErrorMessages []string
		ErrorMaps     map[string][]string
	}

	// swagger:model
	MessageError struct {
		// http response code
		ErrorCode int `json:"code"`
		// errors messages as list
		ErrorMessages []string `json:"messages"`
		// errors messages as map of lists
		ErrorMaps map[string][]string `json:"maps"`
	}

	Options struct {
		code    int
		message string
		field   string
	}

	Option func(*Options)

	// swagger:model
	Empty struct{}
)

var _ json.Marshaler = &Message{}

func NewMessage() *Message {
	return &Message{
		ErrorCode: http.StatusOK,
	}
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	if m.ErrorCode == http.StatusOK {
		if m.Data == nil {
			m.Data = struct{}{}
		}
		return json.Marshal(m.Data)
	}

	if m.ErrorCode == 0 {
		m.ErrorCode = http.StatusInternalServerError
	}

	return json.Marshal(MessageError{
		ErrorCode:     m.ErrorCode,
		ErrorMessages: m.ErrorMessages,
		ErrorMaps:     m.ErrorMaps,
	})
}

func (m *Message) Return(w http.ResponseWriter) error {
	var data, err = m.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")

	// code should not be 0
	// defined by SetError() or MarshalJSON()
	w.WriteHeader(m.ErrorCode)

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

// SetError set error with list and map
func (m *Message) SetError(options ...Option) {
	var args = &Options{}

	var opt Option
	for _, opt = range options {
		opt(args)
	}

	if m.ErrorCode == 0 ||
		m.ErrorCode == http.StatusOK ||
		m.ErrorCode < http.StatusInternalServerError && args.code >= http.StatusInternalServerError {
		m.ErrorCode = args.code
	}

	if args.field == "" && args.message == "" {
		return
	}

	var parts = make([]string, 0, 2)
	if args.field != "" {
		parts = append(parts, args.field)
	}
	if args.message != "" {
		parts = append(parts, args.message)
	}

	m.ErrorMessages = append(m.ErrorMessages, strings.Join(parts, ": "))

	if len(parts) == 2 {
		if m.ErrorMaps == nil {
			m.ErrorMaps = make(map[string][]string)
		}

		m.ErrorMaps[args.field] = append(m.ErrorMaps[args.field], args.message)
	}
}

func MapError(field string, message string) Option {
	return func(args *Options) {
		args.field = field
		args.message = message
	}
}

func Error(message string) Option {
	return func(args *Options) {
		args.message = message
	}
}

func Code(code int) Option {
	return func(args *Options) {
		args.code = code
	}
}
