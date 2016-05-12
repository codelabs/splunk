package splunk

import (
	"strconv"
)

// Response ...
type Response string

// WebService ...
type WebService interface {
	Post(string, string) (Response, error)
}

// Head ...
type Head struct {
	server string
	port   int
}

// Account ...
type Account struct {
	user string
	pswd string
}

// GetAccount ...
func (a *Account) GetAccount() string {
	return "username=" + a.user + "&password=" + a.pswd
}

// GetURL ...
func (h *Head) GetURL() string {
	return "https://" + h.server + ":" + strconv.Itoa(h.port)
}
