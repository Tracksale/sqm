package sqm

import "fmt"

var ErrorInvalidType = fmt.Errorf("Invalid Type")
var ErrorMultipleResults = fmt.Errorf("Passed a pointer for a struct but returned multiple results, check query or use a slice")
var ErrorNoRows = fmt.Errorf("No Rows Affected")
