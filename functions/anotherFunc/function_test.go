package anotherFunc

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const method = "GET"
const url = "/"
const status = http.StatusOK
const body = string("Hello World!\n")

func respBodyConversion(source []byte) string {
	return string(source)
}

// Should be left as is?
func TestBrezBaze(t *testing.T) {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(BrezBaze)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != status {
		t.Errorf("Status: %v | %v", resp.StatusCode, status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if respBodyConversion(respBody) != string(body) {
		t.Errorf("Body:\n--> %v\n--> %v", body, "Hello, World!\n")
	}
}
