package splunk

import (
	"errors"
	"reflect"
	"testing"
)

type testUser struct {
	username string
	password string
}

func (tu *testUser) Fetch(url string, body string) (id string, err error) {

	if url == "https://localhost:5500/services/auth/login" {
		var str = "Bar"
		return str, err
	}

	if url == "https://localhost:8080/services/auth/login" {
		err = errors.New("404 Not Found")
		return "", err
	}

	return "", err
}

// Test User
var tu = &testUser{
	username: "admin",
	password: "changeme",
}

func TestConnect(t *testing.T) {

	session, err := Connect(tu, "localhost", 8080, tu.username, tu.password)
	if err == nil {
		t.Error("Connect fail")
	}

	session, err = Connect(tu, "localhost", 5500, tu.username, tu.password)
	if err != nil {
		t.Error(err)
	}

	var stype = reflect.TypeOf(session)

	if stype.Kind() != reflect.Ptr {
		t.Error("Invalid type returned")
	}
}

func TestGetURL(t *testing.T) {

	session, err := Connect(tu, "localhost", 5500, tu.username, tu.password)
	if err != nil {
		t.Error(err)
	}

	var exp = "https://localhost:5500"
	var rec = session.GetURL()

	if rec != exp {
		t.Error("Received=" + rec + " Expected=" + exp)
	}
}

func TestGetSessionID(t *testing.T) {

	session, err := Connect(tu, "localhost", 5500, tu.username, tu.password)
	if err != nil {
		t.Error(err)
	}

	var exp = "Bar"
	var rec = session.GetSessionID()
	if rec != exp {
		t.Error("Received=" + rec + " Expected=" + exp)
	}
}
