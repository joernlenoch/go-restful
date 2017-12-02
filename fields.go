package restful

import (
	"regexp"
	"strings"

	"fmt"
)

// Takes in a param filter string and creates a sql appropriate representation. Also
// ensures that only parameters are used that
func PrepareFields(raw string, valid []string) (string, error) {

	fields := []string{}

	if raw == "" {
		return "", nil
	}

	parts := strings.Split(raw, ",")

	for _, part := range parts {

		// Skip anything that contains false data. We do not throw errors
		// as it makes it easier to have some custom field types, that must be extended manually.
		ok, err := regexp.MatchString("^[a-zA-Z0-9_]*$", part)
		if err != nil {
			return "", ServerError(err, "Unable to compile regular expression")
		}

		if !ok {
			continue
		}

		isValid := false
		for _, v := range valid {
			if part == v {
				isValid = true
				break
			}
		}

		if !isValid {
			continue
		}

		// Prepare the SQL string
		fields = append(fields, part)
	}

	return strings.Join(fields, ", "), nil
}

// Append the table name when no dot is found to remove all ambiguity
func prependTableName(fields *[]string, table string) {
	for i, field := range *fields {
		if !strings.Contains(field, ".") {
			(*fields)[i] = fmt.Sprintf("%s.%s", table, field)
		}
	}
}
