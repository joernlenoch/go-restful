package restful_test

import (
	"github.com/joernlenoch/go-restful"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareCount(t *testing.T) {
	t.Parallel()

	query, args, err := restful.PrepareCount(restful.Config{
		Fields: restful.Fields{
			restful.Field("name"),
			restful.Field("age"),
		},
		Table: "user",
	}, restful.Request{
		Fields: "name",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT COUNT(*) FROM user", query)
	assert.Empty(t, args, "should not have arguments")
}

func TestPrepareCount_FilterSearch(t *testing.T) {
	t.Parallel()

	query, args, err := restful.PrepareCount(restful.Config{
		Fields: restful.Fields{
			restful.Field("name").Required(),
			restful.Field("age").Searchable(),
			restful.Field("roles").QueryBy("SELECT * FROM roles WHERE roles.name = name").Searchable(),
		},
		Table:    "user",
		Distinct: true,
	}, restful.Request{
		Filter: "name~=a*sd",
		Search: "10",
	})

	assert.NoError(t, err, "must not throw errors")
	assert.Equal(t, "SELECT DISTINCT COUNT(*) FROM user WHERE name LIKE :name0 AND (age LIKE :__restful_search OR roles LIKE :__restful_search)", query)
	assert.Equal(t, 2, len(args), "should have 1 arguments")
	assert.Equal(t, "%a%sd%", args["name0"], "should have transformed args")
}
