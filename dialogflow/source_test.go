package dialogflow

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewFileSource(t *testing.T) {
	source := NewFileSource("../examples/intents.yaml")
	defer func() {
		if err := source.Close(); err != nil {
			t.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(source)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "intents") {
		t.Fail()
	}

	data, err = ioutil.ReadAll(source)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "intents") {
		t.Fail()
	}
}

func TestNewURLSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("../examples/entities.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write(data); err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	source := NewURLSource(ts.URL)
	defer func() {
		if err := source.Close(); err != nil {
			t.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(source)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "entities") {
		t.Fail()
	}
}
