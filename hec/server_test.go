package hec

import "testing"

func TestNewServer(t *testing.T) {

	var s = NewServer("localhost", 8080)
	if s == nil {
		t.Error("Uable to create server")
	}
}

func TestGetHecPostURL(t *testing.T) {

	var s = NewServer("localhost", 8080)
	var expected = "http://localhost:8080/services/collector/event"
	var got = s.GetHecPostURL()

	if got != expected {
		t.Errorf("Expected [%s] Got [%s]", expected, got)
	}
}
