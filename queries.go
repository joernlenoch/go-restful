package restful

import (
	"fmt"
	"strings"
	"errors"
)

var (
	ErrNoFields = errors.New( "No fields selected")
	ErrFieldNotAllowed = errors.New("query contains fields that are not allowed")
	ErrFilterStructure = errors.New("the filter string does not match the allowed structure")
	ErrFilterNotAllowed = errors.New("the filter is not allowed")
	ErrOrderInvalidStructure = errors.New("The order string does not match the allowed structure")
	ErrOrderNotAllowed = errors.New("The order is not allowed")
)

const (
	LimitMax = 50
)

type (

	// The query builder configuration structure.
	Config struct {
		Fields           []field
		Distinct         bool
		Table            string
		Where            string
		GroupBy          string
		AdditionalParams Params
	}

	// Additional params that will be injected into the overall query building proces.
	Params map[string]interface{}

	// Structure that can be used in conjunction with request handling to simplify the
	// config collection process.
	Request struct {
		Fields string `json:"fields" form:"fields" query:"fields"`
		Filter string `json:"filter" form:"filter" query:"filter"`
		Sort   string `json:"sort" form:"sort" query:"sort"`
		Limit  uint   `json:"limit" form:"limit" query:"limit"`
		Offset uint   `json:"offset" form:"offset" query:"offset"`
	}
)

func Prepare(cfg Config, req Request) (query string, args map[string]interface{}, err error) {

	var fields Fields

	// Add the fixed (or default) fields
	if len(cfg.Fields) == 0 {
		err = ErrNoFields
		return
	}

	fields, err = selectFields(req.Fields, cfg.Fields)
	if err != nil {
		return
	}

	// Prepare the order
	order, err := prepareOrder(req.Sort, cfg.Fields)
	if err != nil {
		return "", nil, err
	}

	filter, args, err := prepareFilter(req.Filter, cfg.Fields)
	if err != nil {
		return "", nil, err
	}

	// Merge the filter params and the custom ones
	if cfg.AdditionalParams != nil {
		for k, v := range cfg.AdditionalParams {
			args[k] = v
		}
	}

	// Build the query
	query = "SELECT"

	if cfg.Distinct {
		query += " DISTINCT"
	}

	// Join the field configuration together and add to the query string.
	var fieldStr = fields[0].String()
	for _, f := range fields[1:] {
		fieldStr = fieldStr + ", " + f.String()
	}

	query = fmt.Sprintf("%s %s FROM %s", query, fieldStr, cfg.Table)

	if len(cfg.Where) > 0 || len(filter) > 0 {
		query += " WHERE " + cfg.Where

		if len(cfg.Where) > 0 && len(filter) > 0{
			query = query + " AND "
		}

		query = query + filter
	}


	if len(cfg.GroupBy) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", cfg.GroupBy)
	}

	if len(order) != 0 {
		query += " ORDER BY " + order
	}

	if req.Limit <= 0 || req.Limit > LimitMax {
		req.Limit = LimitMax
	}

	if req.Offset > 0 {
		query += fmt.Sprintf(" LIMIT %d,%d", req.Offset, req.Limit)
	} else {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	return query, args, nil
}


// Takes in a param filter string and creates a sql appropriate representation. Also
// ensures that only parameters are used that
func selectFields(raw string, fields Fields) (Fields, error) {

	selection := make(Fields, 0, len(fields))

	if raw == "" {
		return fields, nil
	}

	// Make sure to add all required fields
	for _, v := range fields {
		if v.IsRequired {
			selection = append(selection, v)
		}
	}

	parts := strings.Split(raw, ",")

	partsLoop:
	for _, part := range parts {

		// Skip anything that contains false data. We do not throw errors
		// as it makes it easier to have some custom field types, that must be extended manually.
		if ok := fieldRegex.MatchString(part); !ok {
			continue
		}

		for _, f := range fields {
			if part == f.Name {

				// Make sure it is not twice in there
				for _, s := range selection {
					if s == f {
						continue partsLoop
					}
				}

				// Add to the final list
				selection = append(selection, f)
				continue partsLoop
			}
		}

		// Throw an error if a field cannot be found
		// TODO ignored - maybe add warnings?
		// return nil, ErrFieldNotAllowed
	}

	if len(selection) == 0 {
		return nil, ErrNoFields
	}

	return selection, nil
}

// Append the table name when no dot is found to remove all ambiguity
func prependTableName(fields *[]string, table string) {
	for i, field := range *fields {
		if !strings.Contains(field, ".") {
			(*fields)[i] = fmt.Sprintf("%s.%s", table, field)
		}
	}
}


// Takes in a param filter string and creates a sql appropriate representation. Also
// ensures that only parameters are used that
func prepareFilter(filter string, valid Fields) (string, map[string]interface{}, error) {

	args := map[string]interface{}{}

	if filter == "" {
		return "", args, nil
	}

	parts := strings.Split(filter, ",")
	sql := make([]string, 0, len(parts))

	for i, part := range parts {

		// rgx.MatchString(part)
		matches := filterRegex.FindStringSubmatch(part)

		if len(matches) != 4 {
			return "", nil, ErrFilterStructure
		}

		// make sure that the given parameter is part of the valid list
		param, cmp, value := matches[1], matches[2], matches[3]

		isValid := false
		for _, v := range valid {
			if param == v.Name {
				isValid = true
				break
			}
		}

		if !isValid {
			return "", nil, ErrFilterNotAllowed
		}

		// Prepare the SQL string
		key := fmt.Sprintf("%s%d", param, i)

		if cmp != "~=" {
			sql = append(sql, fmt.Sprintf("%s %s :%s", param, cmp, key))
			args[key] = value
		} else {
			// Prepare the search parameters by adding an additional parameter

			sql = append(sql, fmt.Sprintf("%s LIKE :%s", param, key))
			search := strings.Replace(value, "*", "%", -1 )
			args[key] = "%"+ search +"%"
		}
	}

	return strings.Join(sql, " AND "), args, nil
}


func prepareOrder(raw string, valid Fields) (string, error) {

	if raw == "" {
		return "", nil
	}

	parts := strings.Split(raw, ",")
	order := make([]string, 0, len(parts))

	for _, part := range parts {

		matches := orderRegex.FindStringSubmatch(part)

		// Important! Remember that the first result is always the full match
		if len(matches) != 3 {
			return "", ErrOrderInvalidStructure
		}

		// Make sure that the given parameter is part of the valid list and that the field exists.
		mark, param := matches[1], matches[2]

		isValid := false
		for _, v := range valid {
			if param == v.Name {
				isValid = true
				break
			}
		}

		if !isValid {
			return "", ErrOrderNotAllowed
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


//
//
func Reduce(i interface{}, cfg Request) (map[string]interface{}, error) {

	// TBD

	return nil, nil
}