package restful

import "fmt"

// Representation of an set of fields.


type (

  Fields []field

  field struct {
    Name string
    Query string
    IsRequired bool
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

func Field(name string) field {
  return field{ Name: name }
}


func (f field) String() string {
  if len(f.Query) > 0 {
    return fmt.Sprintf("%s AS '%s'", f.Query, f.Name)
  }

  return fmt.Sprintf("%s", f.Name)
}


