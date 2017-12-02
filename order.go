package restful

import (
	"fmt"
	"regexp"
	"strings"
)

func PrepareOrder(raw string, valid []string) (string, error) {

	if raw == "" {
		return "", nil
	}

	parts := strings.Split(raw, ",")
	order := make([]string, 0, len(parts))

	for _, part := range parts {

		// Clean the element, remove anything that does not match a simple
		// variable.
		rgx, err := regexp.Compile("^(-|\\+|)([a-zA-ZäüöÄÜÖß0-9_]+)$")
		if err != nil {
			return "", ServerError(err, "Unable to compile regular expression")
		}

		// rgx.MatchString(part)
		matches := rgx.FindStringSubmatch(part)

		// Important! Remember that the first result is always the full match
		if len(matches) != 3 {
			return "", BadRequest(Error{
				Message: fmt.Sprintf("The mark string does not match the allowed structure: %s", part),
				Reason:  "mark",
				DevInfo: fmt.Sprintf("The resulted matches where not 2 (+1): %#v", matches),
			})
		}

		// make sure that the given parameter is part of the valid list
		mark, param := matches[1], matches[2]

		isValid := false
		for _, v := range valid {
			if param == v {
				isValid = true
				break
			}
		}

		if !isValid {
			return "", BadRequest(Error{
				Message: fmt.Sprintf("The mark '%s' is not allowed.", param),
				Reason:  param,
			})
		}

		// Prepare the SQL string
		key := "ASC"
		if mark == "-" {
			key = "DESC"
		}

		order = append(order, fmt.Sprintf("%s %s", param, key))
	}

	return strings.Join(order, ", "), nil
}
