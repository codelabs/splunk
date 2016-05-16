package splunk

import "log"

type logger struct {
	items map[string]string
}

func (l *logger) generate(msg string) {
	log.Print(msg)
}

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
