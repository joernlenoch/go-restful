package restful

import (
	"fmt"
	"strings"
)

const LimitMax = 50

type (
	Config struct {
		*FilterConfig

		DefaultFields []string
		AllowedFields []string
		ValidFilters  []string
		Distinct      bool
		Table         string
	}

	Params map[string]interface{}

	Form struct {
		Fields string `form:"fields" query:"fields"`
		Filter string `form:"filter" query:"filter"`
		Sort   string `form:"sort" query:"sort"`
		Limit  uint   `form:"limit" query:"limit"`
		Offset uint   `form:"offset" query:"offset"`
	}

	FilterConfig struct {
		Fields           string
		Filter           string
		Sort             string
		Limit            uint
		Offset           uint
		GroupBy          string
		Where            string
		AdditionalParams Params
	}
)

func Prepare(cfg Config) (string, map[string]interface{}, error) {

	fields, err := PrepareFields(cfg.Fields, cfg.AllowedFields)
	if err != nil {
		return "", nil, err
	}

	if len(fields) == 0 {
		fields = strings.Join(cfg.DefaultFields, ",")
	}

	order, err := PrepareOrder(cfg.Sort, cfg.ValidFilters)
	if err != nil {
		return "", nil, err
	}

	filter, args, err := PrepareFilter(cfg.Filter, cfg.ValidFilters)
	if err != nil {
		return "", nil, err
	}

	// Merge the filter params and the custom ones
	if cfg.AdditionalParams != nil {
		for k, v := range cfg.AdditionalParams {
			args[k] = v
		}
	}

	query := "SELECT"

	if cfg.Distinct {
		query += " DISTINCT"
	}

	query = fmt.Sprintf("%s %s FROM %s", query, fields, cfg.Table)

	if len(cfg.Where) != 0 {
		query += " WHERE " + cfg.Where
	}

	if len(filter) != 0 {
		query += " HAVING " + filter
	}

	if len(cfg.GroupBy) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", cfg.GroupBy)
	}

	if len(order) != 0 {
		query += " ORDER BY " + order
	}

	if cfg.Limit <= 0 || cfg.Limit > LimitMax {
		cfg.Limit = LimitMax
	}

	if cfg.Offset > 0 {
		query += fmt.Sprintf(" LIMIT %d,%d", cfg.Offset, cfg.Limit)
	} else {
		query += fmt.Sprintf(" LIMIT %d", cfg.Limit)
	}

	return query, args, nil
}
