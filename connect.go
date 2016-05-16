package splunk

import (
	"log"
	"strconv"
)

// PkgLogger - Logging (Singleton) for entire package
type PkgLogger struct {
	*logger
}

// NewLogger - Creates singleton instance for logger
func NewLogger() *PkgLogger {
	var p PkgLogger
	if p.logger == nil {
		p.logger = &logger{
			items: make(map[string]string),
		}
	}

	return &p
}

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

// SessionMgr ...
type SessionMgr struct {
	host string // server name
	port int    // port number
	sid  string // session id
}

// GetURL - gives url for the splunk head
func (s *SessionMgr) GetURL() string {
	return "https://" + s.host + ":" + strconv.Itoa(s.port)
}

// GetSessionID - gives existing session id
func (s *SessionMgr) GetSessionID() string {
	var logger = NewLogger()
	logger.generate(s.sid)
	return s.sid
}

// Connect ...
func Connect(f Fetcher, host string, port int) (*SessionMgr, error) {

	var id string
	var err error
	var logger = NewLogger()

	// Prepare URL
	var url = "https://" + host + ":" + strconv.Itoa(port) + "services/auth/login"
	logger.generate(url)

	id, err = f.Fetch(url, id)

	var session = &SessionMgr{
		host: host,
		port: port,
		sid:  id,
	}
	return session, err
}
