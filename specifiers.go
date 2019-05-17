package sqm

// Internal representation of orderBy
type orderBy struct {
	field     string
	direction int
}

// OrderBy options
var (
	Asc  = 1
	Desc = 2
)

// OrderBy ...orders by
func (q *query) OrderBy(field string, direction int) *query {
	q.orderBy = append(q.orderBy, orderBy{field, direction})

	return q
}

// Limit receives both limit and offset, if you dont want any offsets
// just use 0 as the second argument
func (q *query) Limit(limit int, offset int) *query {
	q.limit = &limit
	q.offset = &offset

	return q
}
