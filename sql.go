package sqm

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"
)

func buildUpdate(fields []string, values []interface{}, paramList *[]interface{}) string {

	var parts []string
	for index, f := range fields {

		switch reflect.TypeOf(values[index]).Name() {
		case "string":
			*paramList = append(*paramList, values[index].(string))
			parts = append(parts, f+"=$"+strconv.Itoa(len(*paramList)))

		case "NullString":
			nullString := values[index].(sql.NullString)
			if nullString.String == "" {
				parts = append(parts, f+"=NULL")
			} else {
				*paramList = append(*paramList, nullString.String)
				parts = append(parts, f+"=$"+strconv.Itoa(len(*paramList)))
			}

		default:
			parts = append(parts, f+"=''")
		}

	}

	return strings.Join(parts, ", ")
}

func parseInsertValues(values []interface{}, paramList *[]interface{}) string {
	var valuesSQL string

	for index, value := range values {
		switch reflect.TypeOf(value).Name() {
		case "string":
			*paramList = append(*paramList, value.(string))
			valuesSQL += "$" + strconv.Itoa(len(*paramList))
		case "NullString":
			nullString := value.(sql.NullString)
			if nullString.String == "" {
				valuesSQL += "NULL"
			} else {
				*paramList = append(*paramList, nullString.String)
				valuesSQL += "$" + strconv.Itoa(len(*paramList))
			}

		default:
			valuesSQL += "''"
		}

		if index != len(values)-1 {
			valuesSQL += ", "
		}

	}

	return valuesSQL
}

func buildInStmt(params []string, paramList *[]interface{}, isNot bool) string {
	var sqlCMD string
	if isNot {
		sqlCMD += " NOT"
	}
	sqlCMD += ` IN( `

	for index, param := range params {
		*paramList = append(*paramList, param)
		sqlCMD += "$" + strconv.Itoa(len(*paramList))
		if index != len(params)-1 {
			sqlCMD += ", "
		}
	}
	if len(params) == 0 {
		sqlCMD += "''"
	}

	sqlCMD += " )"

	return sqlCMD
}

//TODO: return error when query params are invalid or insufficient
//TODO: test this against injections
func (q *Query) toSQL(qT int) (string, []interface{}) {
	query := ""

	var fields []string

	var paramList []interface{}

	for _, field := range q.fields {
		fields = append(fields, field.db)
	}
	// Start
	switch qT {
	case queryTypeSelect:
		query += "SELECT " + strings.Join(fields, ", ") + "\nFROM " + q.table
	case queryTypeUpdate:
		query += "UPDATE " + q.table + "\n SET " + buildUpdate(fields, q.values, &paramList)
	case queryTypeDelete:
		query += "DELETE\nFROM " + q.table
	case queryTypeCount:
		query += "SELECT COUNT(*)\nFROM " + q.table
	case queryTypeInsert:
		query += "INSERT INTO " + q.table + "(" + strings.Join(fields, ", ") + ") VALUES (" + parseInsertValues(q.values, &paramList) + ")"
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
				paramList = append(paramList, condition.params[0])
				query += condition.field + " = $" + strconv.Itoa(len(paramList))

			case NotEqual:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " != $" + strconv.Itoa(len(paramList))

			case Like:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " LIKE $" + strconv.Itoa(len(paramList))
			case NotLike:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " NOT LIKE $" + strconv.Itoa(len(paramList))
			case In:
				query += condition.field + buildInStmt(condition.params, &paramList, false)
			case NotIn:
				query += condition.field + buildInStmt(condition.params, &paramList, true)
			case IsNull:
				query += condition.field + " IS NULL"
			case IsNotNull:
				query += condition.field + " IS NOT NULL"
			case Between:
				paramList = append(paramList, condition.params[0], condition.params[1])
				query += condition.field + " BETWEEN $" + strconv.Itoa(len(paramList)-1) + " AND $" + strconv.Itoa(len(paramList))
			case Greater:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " > $" + strconv.Itoa(len(paramList))
			case GreaterEqual:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " >= $" + strconv.Itoa(len(paramList))
			case Lower:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " < $" + strconv.Itoa(len(paramList))
			case LowerEqual:
				paramList = append(paramList, condition.params[0])
				query += condition.field + " <= $" + strconv.Itoa(len(paramList))
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
			// paramList = append(paramList, oB.field)
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

	return query, paramList
}
