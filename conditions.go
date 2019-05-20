package sqm

import "fmt"

// Possible condition operations
const (
	// I greatly undersestimated this task magic prime number
	// Will fix this later
	internalClose = -2
	internalOpen  = -1
	internalOr    = 0

	Equal     = 1
	NotEqual  = 2
	Like      = 3
	NotLike   = 4
	In        = 5
	NotIn     = 6
	IsNull    = 7
	IsNotNull = 8
)

type groupStruct struct {
	conditions []interface{}
}

// Internal representation of condition
type conditionStruct struct {
	field         string
	conditionType int
	params        []string
}

// Or is a Small hack for api readability
var Or = conditionStruct{conditionType: internalOr}

func internalR(interfaces ...interface{}) []conditionStruct {

	var stack []conditionStruct

	for _, i := range interfaces {
		switch i.(type) {
		case conditionStruct:
			stack = append(stack, i.(conditionStruct))
		case groupStruct:
			stack = append(stack, conditionStruct{conditionType: internalOpen})

			tmpStack := internalR(i.(groupStruct).conditions...)

			for _, tmpC := range tmpStack {
				stack = append(stack, tmpC)
			}
			stack = append(stack, conditionStruct{conditionType: internalClose})
		}
	}

	return stack
}

// Where adds conditions to the stack machine
func (q *Query) Where(interfaces ...interface{}) *Query {
	conditions := internalR(interfaces...)

	for _, c := range conditions {
		q.conditionStack = append(q.conditionStack, c)
	}

	return q
}

// C is a shorthand for Condition and is a typesafe way
// of expressing a node in the condition tree
func C(field string, conditionType int, params ...interface{}) conditionStruct {
	//Accept anything that can be stringified
	c := conditionStruct{
		field:         field,
		conditionType: conditionType,
	}

	for _, param := range params {
		// Transform anything transformable to string
		c.params = append(c.params, fmt.Sprintf("%v", param))
	}

	return c
}

// Group stuff
func Group(items ...interface{}) groupStruct {
	return groupStruct{items}
}
