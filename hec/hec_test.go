package hec

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type TestHEC struct {
	token      string
	source     string
	sourcetype string
	host       string
}

func (t *TestHEC) GetSource() string {
	return t.source
}

func (t *TestHEC) GetSourceType() string {
	return t.sourcetype
}

func (t *TestHEC) GetHost() string {
	return t.host
}

func (t *TestHEC) GetAuthorization() string {
	return "Splunk " + t.token
}

func (t *TestHEC) Post(req *http.Request) (string, error) {

	if req != nil {
		return "ok", nil
	}

	fmt.Println("Referer = " + req.Referer())

	return "fail", errors.New("fail")
}

func TestNewHEC(t *testing.T) {

	var err error
	var token = ""
	var source = ""
	var sourcetype = ""
	var host = ""

	// UT-1: token missing
	_, err = NewHEC(token, source, sourcetype, host)
	if err != nil && err.Error() != missToken {
		t.Errorf("Expected [%s] Received [%s]", missToken, err)
	}

	// UT-2: source missing
	token = "foobar"
	_, err = NewHEC(token, source, sourcetype, host)
	if err != nil && err.Error() != missSrc {
		t.Errorf("Expected [%s] Received [%s]", missSrc, err)
	}

	// UT-3: sourcetype missing
	source = "/var/log/messages"
	_, err = NewHEC(token, source, sourcetype, host)
	if err != nil && err.Error() != missSrctype {
		t.Errorf("Expected [%s] Received [%s]", missSrctype, err)
	}

	var hec *HEC

	// UT-4: host missing
	sourcetype = "syslog"
	hec, err = NewHEC(token, source, sourcetype, host)
	if hec == nil && err != nil {
		t.Error(err)
	}

	// UT-5: host supplied
	host = "hec.client.com"
	hec, err = NewHEC(token, source, sourcetype, host)
	if hec == nil && err != nil {
		t.Error(err)
	}
}

func TestGetAuthorization(t *testing.T) {

	var token = "abcdefgh"
	hec, _ := NewHEC(token, "/var/log/messages", "syslog", "hec.corp.com")
	var expected = "Splunk " + token
	var got = hec.GetAuthorization()

	if got != expected {
		t.Errorf("Expected [%s] Got [%s]", expected, got)
	}
}

func TestCreateRequest(t *testing.T) {

	var token = "abcdefgh"
	hec, err := NewHEC(token, "/var/log/messages", "syslog", "hec.corp.com")

	// Test server where data is posted
	var s = NewServer("localhost", 8080)
	// Headers
	var hdrs = map[string]interface{}{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": hec.GetAuthorization(),
	}

	_, err = CreateRequest("WRONG METHOD", s.GetHecPostURL(), hdrs, "{}")
	if err == nil {
		t.Error("Request created for wrong method")
	}

	req, err := CreateRequest("POST", s.GetHecPostURL(), hdrs, "{}")
	if err != nil {
		t.Error(err)
	}

	var rtype = reflect.TypeOf(req)
	if rtype.Kind() != reflect.Ptr {
		t.Error("Invalid pointer returned")
	}
}

func TestGetSource(t *testing.T) {
	hec, _ := NewHEC("abcd", "/var/log/messages", "syslog", "")
	var expected = "/var/log/messages"
	var got = hec.GetSource()

	if got != expected {
		t.Errorf("Expected [%s] Got [%s]", expected, got)
	}
}

func TestGetSourceType(t *testing.T) {
	hec, _ := NewHEC("abcd", "/var/log/messages", "syslog", "")
	var expected = "syslog"
	var got = hec.GetSourceType()

	if got != expected {
		t.Errorf("Expected [%s] Got [%s]", expected, got)
	}
}

func TestGetHost(t *testing.T) {
	hec, _ := NewHEC("abcd", "/var/log/messages", "syslog", "hec.corp.com")
	var expected = "hec.corp.com"
	var got = hec.GetHost()

	if got != expected {
		t.Errorf("Expected [%s] Got [%s]", expected, got)
	}
}

func TestPostHecEvent(t *testing.T) {

	// Test HEC
	var hec = &TestHEC{
		token:      "zzzzyyyy",
		source:     "/var/log/messages",
		sourcetype: "syslog",
		host:       "hec.corp.com",
	}

	// Test server where data is posted
	var s = NewServer("localhost", 8080)

	// UT-1 Bad event string
	var (
		result string
		err    error
		event  = "{\"foo\":}"
	)

	if _, err = PostHecEvent(hec, s, event); err == nil {
		t.Error("Failed to validate improper json string " + event)
	}

	// UT-2 Good simple event string
	event = "{\"foo\":\"bar\"}"
	if result, err = PostHecEvent(hec, s, event); err != nil {
		t.Error(err)
	}

	if result != "ok" {
		t.Error("post fail")
	}

}
