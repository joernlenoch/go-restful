package restful

import (
	"fmt"
	"strings"
)

func PrepareCount(cfg Config, req Request, searchTarget string) (query string, args map[string]interface{}, err error) {

	args = map[string]interface{}{}

	// Add the fixed (or default) fields
	if len(cfg.Fields) == 0 {
		err = ErrNoFields
		return
	}

	var filter string

	if filter, err = prepareFilter(req.Filter, &args, cfg.Fields); err != nil {
		return
	}

	var search string
	if search, err = prepareSearch(cfg.Fields, &args, req.Search); err != nil {
		return
	}

	// Merge the filter params and the custom ones
	if cfg.AdditionalParams != nil {
		for k, v := range cfg.AdditionalParams {
			args[k] = v
		}
	}

	// Build the query
	query = "SELECT COUNT("

	if cfg.Distinct && searchTarget != "*" {
		query += "DISTINCT "
	}

	query = fmt.Sprintf("%s%s) FROM %s", query, searchTarget, cfg.Table)

	//
	// WHERE
	//

	requirements := []string{}
	if len(cfg.Where) > 0 {
		requirements = append(requirements, cfg.Where)
	}

	if len(filter) > 0 {
		requirements = append(requirements, filter)
	}

	if len(search) > 0 {
		requirements = append(requirements, search)
	}

	if len(requirements) > 0 {
		query += " WHERE " + strings.Join(requirements, " AND ")
	}

	//
	// GROUP
	//

	if len(cfg.GroupBy) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", cfg.GroupBy)
	}

	return query, args, nil
}
