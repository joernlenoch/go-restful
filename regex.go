package restful

import (
	"log"
	"regexp"
)

func init() {
	var err error

	// Clean the element, remove anything that does not match a simple variable.
	orderRegex, err = regexp.Compile("^(-|\\+|)([a-zA-ZäüöÄÜÖß0-9_]+)$")
	if err != nil {
		log.Fatal("Unable to compile regular expression: ", err)
	}

	fieldRegex, err = regexp.Compile("^[a-zA-Z0-9_]*$")
	if err != nil {
		log.Fatal("Unable to compile regular expression: ", err)
	}

	filterRegex, err = regexp.Compile("^([a-zA-Z0-9_]+)(!=|~=|=|<|>|<=|>=|<>)([a-zA-ZäüöÄÜÖß0-9_:.-\\\\*]+)$")
	if err != nil {
		log.Fatal("Unable to compile regular expression: ", err)
	}
}

var (
	orderRegex  *regexp.Regexp
	filterRegex *regexp.Regexp
	fieldRegex  *regexp.Regexp
)
