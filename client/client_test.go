package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSimpleHttpClient(t *testing.T) {
	c := NewSimpleHttpClient(nil)
	if c == nil {
		t.Errorf("Expected an object, received nil")
	}
}

func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	defer ts.Close()
	c := NewSimpleHttpClient(nil)
	r, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
	if err != nil {
		t.Errorf("New request failed: %#v", err)
	}

	resp, body, err := c.Do(r)
	if err != nil {
		t.Errorf("Do failed: %#v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, received %v", resp.StatusCode)
	}

	if string(body) != "hello" {
		t.Errorf("Expected body 'hello', received %s", string(body))
	}
}
