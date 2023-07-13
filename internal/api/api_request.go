package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type (
	Request interface {
		json.Unmarshaler
		Validate(message *Message) (ok bool)
	}
)

var (
	ErrInvalidJson = errors.New("json invalid")
)

func ReadRequest(body io.ReadCloser, req Request) error {
	defer func() {
		_ = body.Close()
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	err = json.Unmarshal(data, req)
	if err != nil {
		return fmt.Errorf("unmarshal fails: %s: %w", err.Error(), ErrInvalidJson)
	}

	return nil
}
