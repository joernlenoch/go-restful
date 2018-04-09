package restful

import "fmt"

// Representation of an set of fields.

type OrderType int

const (
	OrderNone OrderType = iota
	DESC      OrderType = iota
	ASC       OrderType = iota
)

type (
	Fields []field

	field struct {
		Name         string
		Query        string
		IsRequired   bool
		IsSearchable bool
		Order OrderType
	}
)

func (f field) QueryBy(q string) field {
	f.Query = q
	return f
}

func (f field) Required() field {
	f.IsRequired = true
	return f
}

func (f field) Searchable() field {
	f.IsSearchable = true
	return f
}

// Mark this field as default order
func (f field) OrderBy(o OrderType) field {
	f.Order = o
	return f
}

func Field(name string) field {
	return field{Name: name}
}

func (f field) String() string {
	if len(f.Query) > 0 {
		return fmt.Sprintf("%s AS '%s'", f.Query, f.Name)
	}

	return fmt.Sprintf("%s", f.Name)
}
