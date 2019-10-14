package restful_test

import (
	"github.com/joernlenoch/go-restful"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TODO

func TestBuilder(t *testing.T) {
  restful.Select().Distinct()
    .Fields(
      restful.Field("company_id").Required(),
      restful.Field("name"),
      restful.Field("roles").QueryBy(`CONCAT("[", GROUP_CONCAT(JSON_QUOTE(role)),"]")`),
    )
    .Where("").And("").Or("")
    .From("user_company").Join("company").Using("company_id")
    .GroupBy("company_id")
}
*/

func TestPrepare_OptionalFields(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name"),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT name FROM user", query)
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepare_Fields(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name"),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT name, age FROM user", query)
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepare_NoOptionalFields(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Required(),
			restful.Field("age").Required(),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
	})

	assert.NoError(t, err, "must not throw an error")
	assert.NotEmpty(t, query, "should not be empty")
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepare_AltQuery(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").QueryBy("user.name").Required(),
			restful.Field("age"),
		},
		Table: "`user`",
	}, restful.Request{})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT user.name AS 'name', age FROM `user`", query)
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepare_FilterSearch(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Required(),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{
		Filter: "name~=a*sd",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT name, age FROM user WHERE name LIKE :name0", query)
	assert.Equal(t, 1, len(args), "should have 1 arguments")
	assert.Equal(t, "%a%sd%", args["name0"], "should have transformed args")
}

func TestPrepare_Fields_Error(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name"),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{
		Fields: "does_not_exist",
	})

	assert.Error(t, err, "must not throw errors")
	assert.Empty(t, query, "should not return a query")
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepare_Injection(t *testing.T) {
	t.Parallel()

	query, args, err := restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name"),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
		Filter: "name~=\" AND 1;",
	})

	assert.Error(t, err, "must throw errors")

	//
	// Search with multiple fields
	//

	query, args, err = restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Searchable(),
			restful.Field("identifier").Searchable(),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name,identifier",
		Search: "hallo*test",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT name, identifier FROM user WHERE (name LIKE :__restful_search OR identifier LIKE :__restful_search)", query)
	assert.Equal(t, "%hallo%test%", args["__restful_search"])

	//
	// Search with mutliple valid fields, but only one selected
	//

	query, args, err = restful.Prepare(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Searchable(),
			restful.Field("identifier").Searchable(),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
		Search: "hallo*test",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT name FROM user WHERE (name LIKE :__restful_search OR identifier LIKE :__restful_search)", query)
	assert.Equal(t, "%hallo%test%", args["__restful_search"])
}

func TestCount_Basic(t *testing.T) {
	t.Parallel()

	query, _, err := restful.Count(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Searchable(),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT COUNT(*) FROM (SELECT name FROM user) t", query)
}

func TestCount_WithFilter(t *testing.T) {
	t.Parallel()

	query, _, err := restful.Count(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Searchable(),
			restful.Field("age").QueryBy("JSON_QUERY(age)"),
		},
		Table: "user",
	}, restful.Request{
		Filter: "age=4",
		Order:  "-name",
		Limit:  10,
		Offset: 5,
		Search: "name=%test%",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT COUNT(*) FROM (SELECT name, age FROM user WHERE age = :age0 AND name LIKE :__restful_search) t", query)
}
