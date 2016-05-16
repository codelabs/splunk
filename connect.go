package splunk

import (
	"log"
	"strconv"
)

// Fetcher -  fetches data from a remote end point
type Fetcher interface {
	Fetch(url string, body string) (id string, err error)
}

// User - User account to make splunk api call
type User struct {
	username string
	password string
}

// Fetch - Implements Fetcher interface, fetches data from requested
// splunk head
func (u *User) Fetch(url string, body string) (id string, err error) {

	log.Println("url=" + url + " content=" + body)
	var str = "Foo"
	return str, err
}

// Head - Constructs URL for the specified head and port
type Head struct {
	name string
	port int
}

// GetURL - gives url for the splunk head
func (h *Head) GetURL() string {
	return "https://" + h.name + ":" + strconv.Itoa(h.port)
}

// Connect ...
func Connect(f Fetcher, head Head) (string, error) {

	var id string
	var err error

	id, err = f.Fetch(head.GetURL(), id)
	return id, err
}
