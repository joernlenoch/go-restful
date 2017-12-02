package restful

import (
	"fmt"
	"regexp"
	"strings"
)

// Takes in a param filter string and creates a sql appropriate representation. Also
// ensures that only parameters are used that
func PrepareFilter(filter string, valid []string) (string, map[string]interface{}, error) {

	args := map[string]interface{}{}

	if filter == "" {
		return "", args, nil
	}

	parts := strings.Split(filter, ",")
	sql := make([]string, 0, len(parts))

	for i, part := range parts {

		// Clean the element, remove anything that does not match a simple
		// variable.
		rgx, err := regexp.Compile("^([a-zA-Z0-9_]+)(!=|=|<|>|<=|>=|<>)([a-zA-ZäüöÄÜÖß0-9_:.-]+)$")
		if err != nil {
			return "", nil, ServerError(err, "Unable to compile regular expression")
		}

		// rgx.MatchString(part)
		matches := rgx.FindStringSubmatch(part)

		if len(matches) != 4 {
			return "", nil, BadRequest(Error{
				Message: fmt.Sprintf("The filter string does not match the allowed structure: %s", part),
				Reason:  "filter",
				DevInfo: fmt.Sprintf("The resulted matches where not 3 (+1): %#v", matches),
			})
		}

		// make sure that the given parameter is part of the valid list
		param, cmp, value := matches[1], matches[2], matches[3]

		isValid := false
		for _, v := range valid {
			if param == v {
				isValid = true
				break
			}
		}

		if !isValid {
			return "", nil, BadRequest(Error{
				Message: fmt.Sprintf("The filter '%s' is not allowed.", param),
				Reason:  param,
			})
		}

		// Prepare the SQL string
		key := fmt.Sprintf("%s%d", param, i)
		sql = append(sql, fmt.Sprintf("%s %s :%s", param, cmp, key))
		args[key] = value
	}

	return strings.Join(sql, " AND "), args, nil
}
