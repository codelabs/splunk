package splunk

import "testing"

type testUser struct {
	username string
	password string
}

func (tu *testUser) Fetch(url string, body string) (id string, err error) {
	var str = "Bar"
	return str, err
}

func TestGetURL(t *testing.T) {

	var head = Head{
		name: "localhost",
		port: 5500,
	}
	var exp = "https://localhost:5500"
	var rec = head.GetURL()

	if rec != exp {
		t.Error("Received=" + rec + " Expected=" + exp)
	}
}

func TestConnect(t *testing.T) {

	var tu = &testUser{
		username: "admin",
		password: "changeme",
	}

	var head = Head{
		name: "localhost",
		port: 5500,
	}

	var id string
	var err error

	id, err = Connect(tu, head)
	if err != nil {
		t.Error(err)
	}

	if id != "Bar" {
		t.Error("Received=" + id + " Expected=Bar")
	}
}
