package splunk

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
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

	var req *http.Request
	var res *http.Response
	var logger = NewLogger()
	logger.generate("url=" + url + " content=" + body)

	req, err = http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		logger.generate("error")
	}

	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		logger.generate("error")
	}

	defer res.Body.Close()

	var result []byte
	if result, err = ioutil.ReadAll(res.Body); err != nil {
		logger.generate("error")
	}

	return string(result), err
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

// Connect - Connects to splunk server on provided host:port and user account
// details and returns instance of SessionMgr on success.
func Connect(f Fetcher, host string, port int, user string, pass string) (*SessionMgr, error) {

	var id string
	var err error
	var logger = NewLogger()

	// Prepare URL
	var authurl = "https://" + host + ":" + strconv.Itoa(port) + "/services/auth/login"
	logger.generate(authurl)

	var data = url.Values{}
	data.Add("username", user)
	data.Add("password", pass)
	data.Add("output_mode", "json")

	logger.generate(data.Encode())

	var session = &SessionMgr{
		host: host,
		port: port,
	}

	if id, err = f.Fetch(authurl, data.Encode()); err != nil {
		logger.generate("error")
		return session, err
	}

	session = &SessionMgr{
		host: host,
		port: port,
		sid:  id,
	}
	return session, err
}
