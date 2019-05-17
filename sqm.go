package sqm

import (
	"database/sql"
	"fmt"
	"reflect"
)

const (
	queryTypeSelect = 1
	queryTypeUpdate = 2
	queryTypeDelete = 3
	queryTypeCount  = 4
	queryTypeInsert = 5
)

// Internal query representation
type query struct {
	conn *sql.DB

	table          string
	conditionStack []conditionStruct
	orderBy        []orderBy
	limit          *int
	offset         *int

	fields []field

	values []string
}

// Using an specified db connection and table
func Using(db *sql.DB, table string) *query {
	return &query{
		conn:  db,
		table: table,
	}
}

type field struct {
	sField reflect.StructField
	db     string
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
		}
	}

	return fs
}

var ErrorInvalidType = fmt.Errorf("Invalid Type")
var ErrorMultipleResults = fmt.Errorf("Passed a pointer for a struct but returned multiple results, check query or use a slice")

// TODO: Check if fields are valid and writable

// Select Starts a select query chain
func (q *query) Select(i interface{}) error {

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
		tmpField := reflect.New(sF.Type)

		// What the actual fuck
		a := tmpField.Elem().Addr().Interface()

		mappings = append(mappings, a)
	}

	sql := q.toSQL(queryTypeSelect)

	rows, err := q.conn.Query(sql)
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
			f.Set(reflect.ValueOf(mappings[j]).Elem())
		}

		items = reflect.Append(items, item.Elem())
	}

	if !isCollection {
		if items.Len() > 1 {
			return ErrorMultipleResults
		}

		rV.Set(items.Index(0))
	} else {
		rV.Set(items)
	}

	return nil
}

func (q *query) Count(count *int) error {
	query := q.toSQL(queryTypeCount)

	rows, err := q.conn.Query(query)
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
// TODO: Prepared Statements

// Update Starts an update query chain
func (q *query) Insert(i interface{}) (int64, error) {
	var rowsAffected int64

	rV := reflect.ValueOf(i)

	// Only accept structs
	if rV.Kind() != reflect.Struct {
		return rowsAffected, ErrorInvalidType
	}

	q.fields = getFields(rV.Type())

	for j := 0; j < rV.NumField(); j++ {
		q.values = append(q.values, fmt.Sprintf("%v", rV.Field(j)))
	}

	sql := q.toSQL(queryTypeInsert)

	result, err := q.conn.Exec(sql)
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
func (q *query) Update(i interface{}) (int64, error) {
	var rowsAffected int64

	rV := reflect.ValueOf(i)

	// Only accept structs
	if rV.Kind() != reflect.Struct {
		return rowsAffected, ErrorInvalidType
	}

	q.fields = getFields(rV.Type())

	for j := 0; j < rV.NumField(); j++ {
		q.values = append(q.values, fmt.Sprintf("%v", rV.Field(j)))
	}

	sql := q.toSQL(queryTypeUpdate)

	fmt.Println(sql)

	result, err := q.conn.Exec(sql)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

// DeleteFrom starts a delete from query chain
func (q *query) Delete() (int64, error) {
	var rowsAffected int64

	sql := q.toSQL(queryTypeDelete)

	result, err := q.conn.Exec(sql)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}