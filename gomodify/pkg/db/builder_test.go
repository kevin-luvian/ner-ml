package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/kevin-luvian/gomodify/pkg/assert"
)

func TestDB_WhereQuery(t *testing.T) {
	type wants struct {
		where string
		args  []interface{}
	}
	testCases := []struct {
		name string
		arg  GetDBParam
		want wants
	}{{
		name: "success empty",
		arg: GetDBParam{
			Filters: []Filter{},
		},
		want: wants{
			where: "WHERE (1=1)",
		},
	}, {
		name: "success equal",
		arg: GetDBParam{
			Filters: []Filter{{
				Field: "id",
				Value: 1,
			}},
		},
		want: wants{
			where: "WHERE (id = ?)",
			args:  []interface{}{1},
		},
	}, {
		name: "success array",
		arg: GetDBParam{
			Filters: []Filter{{
				Field: "id",
				Value: []int{1, 2, 3},
			}},
		},
		want: wants{
			where: "WHERE (id IN (?,?,?))",
			args:  []interface{}{1, 2, 3},
		},
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.WhereQuery()
			assert.Equal(t, tt.want.where, got.where)
			assert.Equal(t, tt.want.args, got.args)
		})
	}
}

func TestDB_SortQueries(t *testing.T) {
	type wants struct {
		sort string
	}
	testCases := []struct {
		name string
		arg  GetDBParam
		want wants
	}{{
		name: "success empty",
		arg: GetDBParam{
			Sorts: []Sort{},
		},
		want: wants{
			sort: "",
		},
	}, {
		name: "success",
		arg: GetDBParam{
			Sorts: []Sort{{
				Field: "id",
				Asc:   true,
			}},
		},
		want: wants{
			sort: "ORDER BY \"id\" ASC",
		},
	}, {
		name: "success ordered",
		arg: GetDBParam{
			Sorts: []Sort{{
				Field: "id",
				Asc:   true,
				Order: 1,
			}, {
				Field: "name",
				Asc:   false,
			}},
		},
		want: wants{
			sort: "ORDER BY \"name\" DESC, \"id\" ASC",
		},
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.SortQuery()
			assert.Equal(t, tt.want.sort, got.sort)
		})
	}
}

func TestDB_sanitizeSearchString(t *testing.T) {
	testCases := []struct {
		name   string
		arg    string
		expect string
	}{{
		name:   "success",
		arg:    "search-abc",
		expect: "search-abc",
	}, {
		name:   "success replaced",
		arg:    "test_table:abc",
		expect: "test\\_table::abc",
	}, {
		name:   "success replaced",
		arg:    "test_table_1:abc:def",
		expect: "test\\_table\\_1::abc::def",
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			sanitized := sanitizeSearchString(tt.arg)
			assert.Equal(t, tt.expect, sanitized)
		})
	}
}

func TestDB_getSearchFormula(t *testing.T) {
	testCases := []struct {
		name string
		arg  GetDBParam
		want sq.Or
	}{{
		name: "success",
		arg: GetDBParam{
			Search: Search{
				Query:  "search",
				Fields: []string{"table_1", "table_2"},
			},
		},
		want: sq.Or{
			sq.Like{"LOWER(table_1)": "%" + "search" + "%"},
			sq.Like{"LOWER(table_2)": "%" + "search" + "%"},
		},
	}, {
		name: "success no fields",
		arg: GetDBParam{
			Search: Search{
				Query:  "search",
				Fields: []string{},
			},
		},
		want: sq.Or{},
	}, {
		name: "success no query",
		arg: GetDBParam{
			Search: Search{
				Query:  "",
				Fields: []string{"table_1"},
			},
		},
		want: sq.Or{},
	}, {
		name: "success replaced",
		arg: GetDBParam{
			Search: Search{
				Query:  "search:table_one",
				Fields: []string{"table_1"},
			},
		},
		want: sq.Or{
			sq.Like{"LOWER(table_1)": "%" + "search::table\\_one" + "%"},
		},
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.getSearchFormula()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDB_sqSelectAndTrim(t *testing.T) {
	sl := sqSelect()
	q, args, err := sl.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, len(args), 0)
	assert.Equal(t, "SELECT * FROM x", q)

	// panic trimming
	assert.NoPanic(t, func() {
		q = trimSelect(q)
		assert.Equal(t, "", q)
	})

	q = trimSelect("SELECT * FROM x w")
	assert.Equal(t, "w", q)

	sl = sqSelect().Where(sq.Eq{"id": 1})
	q, args, err = sl.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))

	// no panic trimming
	assert.NoPanic(t, func() {
		q = trimSelect(q)
	})

	assert.Equal(t, "WHERE id = ?", q)
}

func TestDB_getWhereFormula(t *testing.T) {
	type wants struct {
		query string
		args  []interface{}
	}
	testCases := []struct {
		name string
		arg  Filter
		want wants
	}{{
		name: "equal",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: EqualOp,
		},
		want: wants{
			query: "id = ?",
			args:  []interface{}{1},
		},
	}, {
		name: "equals",
		arg: Filter{
			Field: "id",
			Value: []int{1, 2, 3},
		},
		want: wants{
			query: "id IN (?,?,?)",
			args:  []interface{}{1, 2, 3},
		},
	}, {
		name: "not_null",
		arg: Filter{
			Field:    "id",
			Value:    []int{1, 2, 3},
			Operator: NotNullOp,
		},
		want: wants{
			query: "id IS NOT NULL",
		},
	}, {
		name: "is_null",
		arg: Filter{
			Field:    "id",
			Operator: IsNullOp,
		},
		want: wants{
			query: "id IS NULL",
		},
	}, {
		name: "not_equal",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: NotEqualOp,
		},
		want: wants{
			query: "id <> ?",
			args:  []interface{}{1},
		},
	}, {
		name: "greater_than",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: GreaterThanOp,
		},
		want: wants{
			query: "id > ?",
			args:  []interface{}{1},
		},
	}, {
		name: "lower_than",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: LowerThanOp,
		},
		want: wants{
			query: "id < ?",
			args:  []interface{}{1},
		},
	}, {
		name: "greater_than_equal",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: GreaterThanEqualOp,
		},
		want: wants{
			query: "id >= ?",
			args:  []interface{}{1},
		},
	}, {
		name: "lower_than_equal",
		arg: Filter{
			Field:    "id",
			Value:    1,
			Operator: LowerThanEqualOp,
		},
		want: wants{
			query: "id <= ?",
			args:  []interface{}{1},
		},
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			q, args, err := tt.arg.getWhereFormula().ToSql()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.query, q)
			assert.Equal(t, tt.want.args, args)
		})
	}
}
