package repository_test

import (
	"testing"

	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
	sqlscope "github.com/kamalyes/go-sqlbuilder/scope"
	"github.com/stretchr/testify/assert"
)

func TestQueryWithScopeFiltersKeepsBusinessScopeFields(t *testing.T) {
	data := sqlscope.NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*sqlscope.ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery().
		AddFilterIfNotEmpty("region_code", "old").
		AddFilterIfNotEmpty("name", "keep").
		WithScopeFilters(sqlscope.NewSQLScopeAdapter(data), func(q *sqlrepo.Query) {
			q.AddFilterIfNotEmpty("region_code", "MM")
			q.AddFilterIfNotEmpty("platform_id", "P1")
		})

	assert.Len(t, query.Filters, 3)
	assert.Equal(t, "name", query.Filters[0].Field)
	assert.Equal(t, "region_code", query.Filters[1].Field)
	assert.Equal(t, "MM", query.Filters[1].Value)
	assert.Equal(t, "platform_id", query.Filters[2].Field)
	assert.Equal(t, "P1", query.Filters[2].Value)
	assert.NotNil(t, query.FilterGroup)

	where, args := query.BuildWhereClause()
	assert.Contains(t, where, "tenant_id = ?")
	assert.Contains(t, where, "region_code = ?")
	assert.Contains(t, where, "platform_id = ?")
	assert.Contains(t, args, "T001")
	assert.Contains(t, args, "MM")
	assert.Contains(t, args, "P1")
}

func TestQueryWithFiltersSkipsNilBuilders(t *testing.T) {
	query := sqlrepo.NewQuery().WithFilters(nil, func(q *sqlrepo.Query) {
		q.AddFilter(sqlrepo.NewFilter("status", sqlconstants.OP_EQ, "active"))
	})

	assert.Len(t, query.Filters, 1)
	assert.Equal(t, "status", query.Filters[0].Field)
}
