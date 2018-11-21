package noDBExp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const url = "/"
const status = http.StatusOK
const body = string("Hello World!\n")
func respBodyConversion(source []byte)(string) {
	return string(source)
}


// Should be left as is?
func TestBrezBaze(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(BrezBaze)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status: %v | %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if respBodyConversion(respBody) != string(body) {
		t.Errorf("Body:\n--> %v\n--> %v", body, "Hello, World!\n")
	}
}
