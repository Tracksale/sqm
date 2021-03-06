package sqm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	queryTypeSelect = 1
	queryTypeUpdate = 2
	queryTypeDelete = 3
	queryTypeCount  = 4
	queryTypeInsert = 5
)

// Query Internal query representation
type Query struct {
	conn *sql.DB

	table          string
	conditionStack []ConditionStruct
	orderBy        []orderBy
	limit          *int
	offset         *int
	debug          bool
	fields         []field

	values []interface{}
}

// Using an specified db connection and table
func Using(db *sql.DB, table string) *Query {
	return &Query{
		conn:  db,
		table: table,
		debug: false,
	}
}

//Debug enable SQL log
func (q *Query) Debug() *Query {
	q.debug = true
	return q
}

type field struct {
	sField reflect.StructField
	db     string
}

func (q Query) log(msg string) {
	if q.debug {
		ansiColor := "\033[1;33m%s\033[0m"
		fmt.Printf(ansiColor, strings.Replace(msg, "\n", " ", -1))
		fmt.Println("")
	}
}

func getFields(rT reflect.Type) []field {
	var fs []field
	// element
	if rT.Kind() != reflect.Struct {
		return fs
	}

	for j := 0; j < rT.NumField(); j++ {
		f := rT.Field(j)

		db, exists := f.Tag.Lookup("db")
		if exists {
			fs = append(fs, field{f, db})
		} else {
			db, exists := f.Tag.Lookup("json")
			if exists {
				db := strings.Split(db, ",")[0]
				fs = append(fs, field{f, db})
			}
		}
	}

	return fs
}

// TODO: Check if fields are valid and writable

// Select Starts a select query chain
func (q *Query) Select(i interface{}) error {

	rV := reflect.ValueOf(i)

	// Only accept pointers
	if rV.Kind() != reflect.Ptr {
		return ErrorInvalidType
	}

	// Follow pointer
	rV = rV.Elem()

	var rT reflect.Type
	var isCollection = false

	if rV.Kind() == reflect.Slice {
		// Follow slice into type
		rT = rV.Type().Elem()
		isCollection = true
	} else if rV.Kind() == reflect.Struct {
		rT = rV.Type()
	}

	if rT.Kind() != reflect.Struct {
		return ErrorInvalidType
	}

	q.fields = getFields(rT)

	var mappings []interface{}

	for _, field := range q.fields {
		sF := field.sField
		var tmpField reflect.Value
		sfKind := sF.Type.Kind()

		if sfKind == reflect.Map || sfKind == reflect.Slice || sfKind == reflect.Struct || sfKind == reflect.Interface {
			tmpField = reflect.New(reflect.TypeOf([]byte{}))
		} else {
			tmpField = reflect.New(sF.Type)
		}

		// What the actual fuck
		a := tmpField.Elem().Addr().Interface()

		mappings = append(mappings, a)
	}

	sql, paramList := q.toSQL(queryTypeSelect)
	q.log(sql)

	rows, err := q.conn.Query(sql, paramList...)
	if err != nil {
		return err
	}

	items := reflect.MakeSlice(reflect.SliceOf(rT), 0, 1)
	for rows.Next() {
		err = rows.Scan(mappings...)

		// Return if we found a single row error
		if err != nil {
			return err
		}

		item := reflect.New(rT)
		for j := 0; j < rT.NumField(); j++ {
			f := item.Elem().Field(j)
			scanRes := reflect.ValueOf(mappings[j]).Elem()
			fKind := f.Type().Kind()

			switch fKind {
			case reflect.Map, reflect.Slice, reflect.Struct, reflect.Interface:
				tmpParse := reflect.New(f.Type())
				json.Unmarshal(scanRes.Bytes(), tmpParse.Interface())
				f.Set(tmpParse.Elem())
			case reflect.Ptr:
				if scanRes.String() == "" || scanRes.String() == "<nil>" {
					f.Set(reflect.ValueOf(nil))
				} else {
					f.Set(scanRes)
				}
			default:
				f.Set(scanRes)
			}
		}

		items = reflect.Append(items, item.Elem())
	}

	if !isCollection {
		if items.Len() > 1 {
			return ErrorMultipleResults
		}
		if items.Len() == 0 {
			return ErrorNoRows
		}
		rV.Set(items.Index(0))

	} else {
		rV.Set(items)
	}

	return nil
}

//Count ...
func (q *Query) Count(count *int) error {
	query, paramList := q.toSQL(queryTypeCount)
	q.log(query)

	rows, err := q.conn.Query(query, paramList...)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(count)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Accept slices
// TODO: Prepared Statements - use db.Query instead of db.Exec - ensure against SQL injection
// 			create a map(or something) and construct db.Query

// Insert Starts an insert query chain
func (q *Query) Insert(i interface{}) (int64, error) {
	var rowsAffected int64

	rV := reflect.ValueOf(i)

	if rV.Kind() == reflect.Ptr {
		rV = reflect.Indirect(rV)
	}

	// Only accept structs
	if rV.Kind() != reflect.Struct {
		return rowsAffected, ErrorInvalidType
	}

	q.fields = getFields(rV.Type())

	for j := 0; j < rV.NumField(); j++ {
		q.values = append(q.values, prepareInput(rV.Field(j)))
	}

	sql, paramList := q.toSQL(queryTypeInsert)
	q.log(sql)

	result, err := q.conn.Exec(sql, paramList...)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

// Update Starts an update query chain
func (q *Query) Update(i interface{}) (int64, error) {
	var rowsAffected int64

	rV := reflect.ValueOf(i)

	if rV.Kind() == reflect.Ptr {
		rV = reflect.Indirect(reflect.ValueOf(i))
	}

	// Only accept structs
	if rV.Kind() != reflect.Struct {
		return rowsAffected, ErrorInvalidType
	}

	q.fields = getFields(rV.Type())

	for j := 0; j < rV.NumField(); j++ {
		q.values = append(q.values, prepareInput(rV.Field(j)))
	}

	sql, paramList := q.toSQL(queryTypeUpdate)
	q.log(sql)

	result, err := q.conn.Exec(sql, paramList...)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

func prepareInput(field reflect.Value) interface{} {
	fieldKind := field.Type().Kind()

	switch fieldKind {

	case reflect.Map, reflect.Slice, reflect.Struct, reflect.Interface:
		value, _ := json.Marshal(field.Interface())
		return sql.NullString{String: fmt.Sprintf("%v", string(value))}

	case reflect.Ptr:
		indirectValue := reflect.Indirect(field)

		if indirectValue.Kind() == reflect.Invalid {
			return sql.NullString{String: "", Valid: false}
		}
		return sql.NullString{String: fmt.Sprintf("%v", indirectValue)}

	default:
		return fmt.Sprintf("%v", field)
	}
}

// Delete starts a delete from query chain
func (q *Query) Delete() (int64, error) {
	var rowsAffected int64

	sql, paramList := q.toSQL(queryTypeDelete)
	q.log(sql)

	result, err := q.conn.Exec(sql, paramList...)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}
