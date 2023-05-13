package db

type Operator string

const (
	NotNullOp          Operator = "IS NOT NULL"
	IsNullOp           Operator = "IS NULL"
	EqualOp            Operator = "="
	NotEqualOp         Operator = "<>"
	GreaterThanOp      Operator = ">"
	LowerThanOp        Operator = "<"
	GreaterThanEqualOp Operator = ">="
	LowerThanEqualOp   Operator = "<="
)

type GetDBParam struct {
	Search       Search
	Filters      []Filter
	Sorts        []Sort
	Limit        int
	Offset       int
	DisableGet   bool
	DisableCount bool

	where string
	sort  string
	args  []interface{}
}

type Search struct {
	Query  string
	Fields []string
}

type Sort struct {
	Order int
	Field string
	Asc   bool
}

type Filter struct {
	Field    string
	Value    interface{}
	Operator Operator
}
