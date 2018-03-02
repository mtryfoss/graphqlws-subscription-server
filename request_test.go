package gss

import "testing"

func TestValidation(t *testing.T) {
	r := &RequestData{}
	var err error
	err = r.Validate()
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "require channel" {
		t.Error("unexpected error message")
	}

	r.Channel = "foo"
	err = r.Validate()
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "require payload" {
		t.Error("unexpected error message")
	}

	r.Payload = "foo"
	err = r.Validate()
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "require field" {
		t.Error("unexpected error message")
	}
}

func TestReadFromBytes(t *testing.T) {
	d, err := NewRequestDataFromBytes([]byte("hoge{"))
	if d != nil {
		t.Error("data should not exists")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "cannot parse invalid JSON request data" {
		t.Error("unexpected error message")
	}
	d, err = NewRequestDataFromBytes([]byte(`{"users":["hoge","fuga"]}`))
	if d != nil {
		t.Error("data should not exists")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "require channel" {
		t.Error("unexpected error message")
	}

	d, err = NewRequestDataFromBytes([]byte(`{"users":["hoge","fuga"],"channel":"foo","field":"hey","payload":"bar"}`))
	if d == nil {
		t.Error("data should exists")
	}
	if err != nil {
		t.Error("error should not exists")
	}
	if d.Channel != "foo" {
		t.Error("channel data is wrong")
	}
	if len(d.Users) != 2 {
		t.Error("users length is wrong")
	}
	if d.Users[0] != "hoge" || d.Users[1] != "fuga" {
		t.Error("users data is wrong")
	}
	if d.Payload.(string) != "bar" {
		t.Error("payload data is wrong")
	}
}
