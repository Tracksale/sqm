package sqm

import (
	"strconv"
	"strings"
)

func buildUpdate(fields []string, values []string) string {

	var parts []string
	for index, f := range fields {
		parts = append(parts, f+"='"+values[index]+"'")
	}

	return strings.Join(parts, ", ")
}

func (q *Query) toSQL(qT int) string {
	query := ""

	var fields []string

	for _, field := range q.fields {
		fields = append(fields, field.db)
	}
	// Start
	switch qT {
	case queryTypeSelect:
		query += "SELECT " + strings.Join(fields, ", ") + "\nFROM " + q.table
	case queryTypeUpdate:
		query += "UPDATE " + q.table + "\n SET " + buildUpdate(fields, q.values)
	case queryTypeDelete:
		query += "DELETE\nFROM " + q.table
	case queryTypeCount:
		query += "SELECT COUNT(*)\nFROM " + q.table
	case queryTypeInsert:
		query += "INSERT INTO " + q.table + "(" + strings.Join(fields, ", ") + ") VALUES ('" + strings.Join(q.values, "', '") + "')"
	}

	// Where conditions
	if len(q.conditionStack) > 0 {
		query += "\nWHERE"

		query += "\n\t"

		tmpStack := append([]ConditionStruct(nil), q.conditionStack...)

		for len(tmpStack) > 0 {
			//Pop
			condition := tmpStack[0]
			tmpStack = tmpStack[1:]

			switch condition.conditionType {
			case internalOpen:
				query += "("
			case internalClose:
				query += ")"
			case internalOr:
				query += " OR "
			case Equal:
				query += condition.field + " = '" + condition.params[0] + "'"
			case NotEqual:
				query += condition.field + " != '" + condition.params[0] + "'"
			case Like:
				query += condition.field + " LIKE '" + condition.params[0] + "'"
			case NotLike:
				query += condition.field + " NOT LIKE '" + condition.params[0] + "'"
			case In:
				query += condition.field + ` IN ('` + strings.Join(condition.params, "', '") + `')`
			case NotIn:
				query += condition.field + ` NOT IN ('` + strings.Join(condition.params, "', '") + `')`
			case IsNull:
				query += condition.field + " IS NULL"
			case IsNotNull:
				query += condition.field + " IS NOT NULL"
			}

			// I am so sorry, i will find a better way
			if len(tmpStack) > 0 {
				if tmpStack[0].conditionType != internalOr &&
					tmpStack[0].conditionType != internalClose &&
					condition.conditionType != internalOpen &&
					condition.conditionType != internalOr {
					query += " AND "
				}
			}
		}
	}

	// OrderBy
	if len(q.orderBy) > 0 {
		query += "\nORDER BY "

		oBs := []string{}

		for _, oB := range q.orderBy {
			tmpOb := "\n\t" + oB.field + " "
			if oB.direction == Asc {
				tmpOb += "ASC"
			} else {
				tmpOb += "DESC"
			}
			oBs = append(oBs, tmpOb)
		}

		query += strings.Join(oBs, ",")
	}

	// Limit
	if q.limit != nil {
		query += "\nLIMIT " + strconv.Itoa(*q.limit)
	}

	// Offset
	if q.offset != nil {
		query += "\nOFFSET " + strconv.Itoa(*q.offset)
	}

	return query
}
