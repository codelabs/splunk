package splunk

import "log"

type logger struct {
	items map[string]string
}

func (l *logger) generate(msg string) {
	log.Print(msg)
}
