package gss

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func notifyRequestTest(method, url, content, contentType string) (int, []byte) {
	req, _ := http.NewRequest(method, url, strings.NewReader(content))
	ctx := context.Background()
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, _ := client.Do(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	return resp.StatusCode, buf.Bytes()
}

type notifyTestCase struct {
	Label           string
	Method          string
	Body            string
	ContentType     string
	StatusCode      int
	ResponseSuccess bool
	ResponseError   string
}

func TestNotificationHandler(t *testing.T) {
	ch := make(chan *RequestData, 1)
	mux := http.NewServeMux()
	mux.Handle("/notify", NewNotifyHandler(ch))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, testCase := range []notifyTestCase{
		notifyTestCase{
			Label:           "invalid method",
			Method:          "GET",
			Body:            "foobar",
			ContentType:     "application/json",
			StatusCode:      http.StatusMethodNotAllowed,
			ResponseSuccess: false,
			ResponseError:   "not allowed method",
		},
		notifyTestCase{
			Label:           "invalid content-type",
			Method:          "POST",
			Body:            "foobar",
			ContentType:     "application/x-www-form-urlencoded",
			StatusCode:      http.StatusBadRequest,
			ResponseSuccess: false,
			ResponseError:   "Content-Type requires application/json",
		},
		notifyTestCase{
			Label:           "wrong body",
			Method:          "POST",
			Body:            "foobar",
			ContentType:     "application/json",
			StatusCode:      http.StatusBadRequest,
			ResponseSuccess: false,
			ResponseError:   "cannot parse invalid JSON request data",
		},
	} {
		t.Run(testCase.Label, func(t *testing.T) {
			statusCode, bufBytes := notifyRequestTest(testCase.Method, ts.URL+"/notify", testCase.Body, testCase.ContentType)
			if statusCode != testCase.StatusCode {
				t.Error(testCase.Label + ": /notify invalid statusCode: " + strconv.Itoa(statusCode))
			}
			resData := &NotificationResponse{}
			json.Unmarshal(bufBytes, resData)
			if resData.Success != testCase.ResponseSuccess {
				t.Error(testCase.Label + ": response.Success is wrong")
			}
			if len(resData.Errors) != 1 {
				t.Error(testCase.Label + ": response.Errors count should be 1")
			}
			if !resData.Success && resData.Errors[0] != testCase.ResponseError {
				t.Error(testCase.Label + ": unexpected error message")
			}
		})
	}

	for _, testCase := range []notifyTestCase{
		notifyTestCase{
			Label:           "success",
			Method:          "POST",
			Body:            `{"channel":"foo","users":["hoge","fuga"],"payload":{"aaa":"iii","uu":{"ee":"oo"}}}`,
			ContentType:     "application/json",
			StatusCode:      http.StatusOK,
			ResponseSuccess: true,
		},
	} {
		t.Run(testCase.Label, func(t *testing.T) {
			statusCode, bufBytes := notifyRequestTest(testCase.Method, ts.URL+"/notify", testCase.Body, testCase.ContentType)
			if statusCode != testCase.StatusCode {
				t.Error(testCase.Label + ": /notify invalid statusCode: " + strconv.Itoa(statusCode))
			}
			resData := &NotificationResponse{}
			json.Unmarshal(bufBytes, resData)
			if resData.Success != testCase.ResponseSuccess {
				t.Error(testCase.Label + ": response.Success is wrong. " + resData.Errors[0])
			}
			if len(resData.Errors) != 0 {
				t.Error(testCase.Label + ": response.Errors count should be 0")
			}
			data := <-ch
			if data.Channel != "foo" {
				t.Error(testCase.Label + ": data.Channel should be 'foo'")
			}
			if len(data.Users) != 2 {
				t.Error(testCase.Label + ": users cound should be 2")
			}
			if data.Payload == nil {
				t.Error(testCase.Label + ": payload should exists")
			}
		})
	}

}
