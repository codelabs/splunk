package hec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// HTTPEventCollector ...
// Interface here
type HTTPEventCollector interface {
	GetSource() string
	GetSourceType() string
	GetHost() string
	GetAuthorization() string
	Post(*http.Request) (string, error)
}

// HEC ...
type HEC struct {
	token      string
	source     string
	sourcetype string
	host       string
}

// Event - is of type map
type Event map[string]interface{}

// NewHEC ...
func NewHEC(tkn string, src string, srctype string, host string) (*HEC, error) {

	// Enable microseconds for logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Error out if missing token
	if tkn == "" {
		return nil, errors.New(missToken)
	}
	// Error out if missing source
	if src == "" {
		return nil, errors.New(missSrc)
	}
	// Error out if missing sourcetype
	if srctype == "" {
		return nil, errors.New(missSrctype)
	}
	// Get hostname if not passed
	if host == "" {

		// Error out if unable to get hostname
		var err error
		if host, err = os.Hostname(); err != nil {
			log.Printf("token=%s source=%s sourcetype=%s host=%s error=%v", tkn, src, srctype, host, err)
			return nil, err
		}
	}

	log.Printf("token=%s source=%s sourcetype=%s host=%s", tkn, src, srctype, host)

	var hec = &HEC{
		token:      tkn,
		source:     src,
		sourcetype: srctype,
		host:       host,
	}

	return hec, nil
}

// GetSource ...
func (h *HEC) GetSource() string {
	return h.source
}

// GetSourceType ...
func (h *HEC) GetSourceType() string {
	return h.sourcetype
}

// GetHost ...
func (h *HEC) GetHost() string {
	return h.host
}

// GetAuthorization ...
func (h *HEC) GetAuthorization() string {
	return "Splunk " + h.token
}

// CreateRequest ...
func CreateRequest(method string, url string, headers map[string]interface{}, body string) (*http.Request, error) {

	request, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		log.Println("Error on creating request " + err.Error())
		return nil, err
	}

	// Set Headers
	for k, v := range headers {
		request.Header.Set(k, v.(string))
	}

	return request, nil
}

// Post ...
func (h *HEC) Post(req *http.Request) (string, error) {

	var res *http.Response
	var err error

	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		log.Println("Error on client Do " + err.Error())
		return "", err
	}

	defer res.Body.Close()

	var result []byte
	if result, err = ioutil.ReadAll(res.Body); err != nil {
		return "", err
	}

	return string(result), err
}

// PostHecEvent ...
func PostHecEvent(h HTTPEventCollector, s *Server, e string) (string, error) {

	var (
		err      error         // Holds error information
		event    Event         // Holds the unmarshalled json
		hecEvent Event         // Final event to send to splunk
		hecbyte  []byte        // byte-ified final event
		req      *http.Request // Request object
	)

	hecEvent = Event{
		"source":     h.GetSource(),
		"sourcetype": h.GetSourceType(),
		"host":       h.GetHost(),
		"time":       time.Now().UnixNano() / 1000, // epoch in microseconds
	}

	if err = json.Unmarshal([]byte(e), &event); err != nil {
		log.Println(err)
		return "", err
	}

	hecEvent["event"] = event

	if hecbyte, err = json.Marshal(hecEvent); err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("url=%s event=%s", s.GetHecPostURL(), string(hecbyte))

	var hdrs = map[string]interface{}{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": h.GetAuthorization(),
	}

	if req, err = CreateRequest("POST", s.GetHecPostURL(), hdrs, string(hecbyte)); err != nil {
		return "", err
	}

	var result string
	result, err = h.Post(req)
	return result, err
}
