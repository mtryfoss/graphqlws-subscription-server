package gss

import (
	"encoding/json"
	"errors"
)

type RequestData struct {
	Users   []string    `json:"users"`
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

func (d *RequestData) Validate() error {
	if d.Channel == "" {
		return errors.New("require channel")
	}
	if d.Payload == nil {
		return errors.New("require payload")
	}
	return nil
}

func NewRequestDataFromBytes(b []byte) (*RequestData, error) {
	data := &RequestData{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, errors.New("cannot parse invalid JSON request data")
	}
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return data, nil
}
