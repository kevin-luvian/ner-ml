package db

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/kevin-luvian/gomodify/pkg/logging"
	"github.com/lib/pq"
)

// WhereQuery enable where query
func (g GetDBParam) GetQuery(query string) (string, []interface{}) {
	qWhere := fmt.Sprintf("%s %s", g.where, g.sort)
	query = fmt.Sprintf(query, qWhere)

	query, args, err := sqlx.In(query, g.args...)
	if err != nil {
		// caused by improper query / missing arg
		panic(err)
	}
	return query, args
}

// WhereQuery enable where query
func (g GetDBParam) WhereQuery() GetDBParam {
	filters := sq.And{}
	search := g.getSearchFormula()

	if len(search) > 0 {
		filters = append(filters, search)
	}

	for _, f := range g.Filters {
		filters = append(filters, f.getWhereFormula())
	}

	q, args, err := sqSelect().Where(filters).ToSql()
	if err != nil {
		logging.Errorln(err)
		return g
	}

	g.where = trimSelect(q)
	g.args = args

	return g
}

// SortQuery enable sorting. returned args always empty
func (g GetDBParam) SortQuery() GetDBParam {
	if len(g.Sorts) == 0 {
		return g
	}

	// sort ordering asc
	sort.Slice(g.Sorts, func(i, j int) bool {
		return g.Sorts[i].Order-g.Sorts[j].Order < 0
	})

	ordering := make([]string, len(g.Sorts))
	for i := range g.Sorts {
		dir := "ASC"
		if !g.Sorts[i].Asc {
			dir = "DESC"
		}
		ordering[i] = pq.QuoteIdentifier(g.Sorts[i].Field) + " " + dir
	}

	q, _, err := sqSelect().OrderBy(ordering...).ToSql()
	if err != nil {
		logging.Errorln(err)
		panic(err)
	}

	g.sort = trimSelect(q)

	return g
}

func (g *GetDBParam) getSearchFormula() sq.Or {
	search := sq.Or{}

	if len(g.Search.Fields) == 0 {
		return search
	}

	if g.Search.Query != "" {
		s := "%" + sanitizeSearchString(strings.ToLower(g.Search.Query)) + "%"

		for _, f := range g.Search.Fields {
			search = append(search, sq.Like{fmt.Sprintf("LOWER(%s)", f): s})
		}
	}

	return search
}

func sanitizeSearchString(s string) string {
	re := regexp.MustCompile(":")
	newString := re.ReplaceAllString(s, "::")

	re = regexp.MustCompile("_")
	return re.ReplaceAllString(newString, "\\_")
}

func sqSelect() sq.SelectBuilder {
	return sq.Select("*").From("x")
}

// trim "SELECT * FROM x ..."
func trimSelect(q string) string {
	if len(q) < 16 {
		return ""
	}
	return q[16:]
}

func (f *Filter) getWhereFormula() sq.Sqlizer {
	switch v := f.Value.(type) {
	default:
		switch f.Operator {
		case NotNullOp:
			return sq.NotEq{f.Field: nil}
		case IsNullOp:
			return sq.Eq{f.Field: nil}
		case EqualOp:
			return sq.Eq{f.Field: v}
		case NotEqualOp:
			return sq.NotEq{f.Field: v}
		case GreaterThanOp:
			return sq.Gt{f.Field: v}
		case LowerThanOp:
			return sq.Lt{f.Field: v}
		case GreaterThanEqualOp:
			return sq.GtOrEq{f.Field: v}
		case LowerThanEqualOp:
			return sq.LtOrEq{f.Field: v}
		default:
			return sq.Eq{f.Field: v}
		}
	}
}
