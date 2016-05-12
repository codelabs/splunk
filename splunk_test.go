package splunk

import "testing"

func TestGetAccount(t *testing.T) {
	var acct = Account{"admin", "changeme"}
	var exp = "username=admin&password=changeme"
	var rec = acct.GetAccount()
	if rec != exp {
		t.Error("Received " + rec + " Expected " + exp)
	}
}

func TestGetURL(t *testing.T) {
	var head = Head{"localhost", 5500}
	var exp = "https://localhost:5500"
	var rec = head.GetURL()
	if rec != exp {
		t.Error("Received " + rec + " Expected " + exp)
	}
}
