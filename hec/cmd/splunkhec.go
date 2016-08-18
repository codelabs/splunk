package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/codelabs/splunk/hec"
)

// RCFilecfg - List of tokens
type RCFilecfg struct {
	Tokens            []Record
	DefaultTokenIndex int `json:"default_token_index"`
}

// Record - token structure
type Record struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	Source     string `json:"src"`
	Sourcetype string `json:"srctype"`
	Server     string `json:"server"`
	Port       int    `json:"port"`
}

func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func getRCFile() string {

	// Get user's home directory and construct
	// rcfile path
	var (
		u      *user.User // user information
		err    error      // error
		rcfile string     // full path to rcfile
	)

	if u, err = user.Current(); err != nil {
		log.Fatal(err)
	}

	rcfile = u.HomeDir + "/.splunkhecrc"
	log.Println("Config File: " + rcfile)

	return rcfile
}

func loadHecRCFile() RCFilecfg {

	var rcfile = getRCFile() // full path to .splunkhecrc file
	var err error            // error handler
	var fd []byte            // File Descriptor for rcfile

	if fd, err = ioutil.ReadFile(rcfile); err != nil {
		log.Fatal(err)
	}

	var cfg RCFilecfg
	json.Unmarshal(fd, &cfg)

	log.Println(cfg)
	return cfg
}

func main() {

	// Command Line Flags
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\t main [ -name=<index name> ]\n")
		flag.PrintDefaults()
	}
	name := flag.String("name", "", "Index name of the token")
	event := flag.String("event", "", "json event to post to splunk")
	flag.Parse()

	// Commad Line Flag Validations
	if *event == "" {
		log.Fatal("Missing json event")
	}

	if !isValidJSON(*event) {
		log.Fatal("Not a valid json event")
	}

	// Load the rcfile
	var cfg = loadHecRCFile()
	var tokenrec Record

	// If default token is not selected, then look at command line
	// param for name of the token to use
	if cfg.DefaultTokenIndex == 0 {

		if *name == "" {
			log.Fatal(errors.New("Missing index name in command line argument"))
		}

		for _, rec := range cfg.Tokens {
			if rec.Name == *name {
				tokenrec = rec
				break
			}
		}
	} else {
		tokenrec = cfg.Tokens[cfg.DefaultTokenIndex-1]
	}

	fmt.Println(tokenrec)

	var h *hec.HEC
	var result string
	var err error

	if h, err = hec.NewHEC(tokenrec.Token, tokenrec.Source, tokenrec.Sourcetype, ""); err != nil {
		log.Fatal(err)
	}

	var splunk = hec.NewServer(tokenrec.Server, tokenrec.Port)

	if result, err = hec.PostHecEvent(h, splunk, *event); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}
