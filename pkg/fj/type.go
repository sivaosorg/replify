package fj

import (
	"unsafe"
)

// Type represents the different possible types for a JSON value.
// It is used to indicate the specific type of a JSON value, such as a string, number, boolean, etc.
type Type int

// Context represents a JSON value returned from the Get() function.
// It stores information about a specific JSON element, including its type,
// raw string data, string representation, numeric value, index in the original JSON,
// and the indexes of elements that match a path containing a '#'.
type Context struct {
	// kind is the JSON type (such as String, Number, Object, etc.).
	kind Type

	// raw contains the raw JSON string that has not been processed or parsed.
	raw string

	// str contains the string value of the JSON element, if it is a string type.
	str string

	// num contains the numeric value of the JSON element, if it is a number type.
	num float64

	// idx holds the position of the raw JSON value in the original JSON string.
	// A value of 0 means the index is unknown.
	idx int

	// idxs holds the indices of all elements that match a path containing the '#' query character.
	idxs []int

	// err stores any error encountered during processing or parsing of the JSON element.
	// This field is used to capture issues such as invalid JSON syntax, type mismatches,
	// or any other error that may occur while retrieving or interpreting the JSON value.
	// If no error occurred, this field will be nil.
	err error
}

// qCtx is a simplified version of the Context struct,
// primarily used to store intermediate results for JSON path processing.
type qCtx struct {
	// arr stores a slice of Context elements representing an array result.
	arr []Context

	// elems stores a slice of interface{} elements, used for intermediary results.
	elems []interface{}

	// ops maps string keys to Context values, used for handling operations on specific paths.
	ops map[string]Context

	// opResults maps string keys to interface{} values, used for operation results on specific paths.
	opResults map[string]interface{}

	// vn stores a byte value for a specific operation, likely used for flagging or identifying states.
	vn byte
}

// wildcard represents a wildcard element in a JSON path query.
// It is used to represent patterns that match various JSON values in path expressions.
type wildcard struct {
	// part represents a specific segment of the wildcard pattern.
	part string

	// path represents the path expression that includes the wildcard.
	path string

	// pipe represents an operation or transformation to be applied to the wildcard.
	pipe string

	// piped indicates whether the wildcard has been piped for further processing.
	piped bool

	// wild indicates whether this wildcard is a true wildcard (e.g., "*").
	wild bool

	// more indicates whether there are more operations or elements to process after this wildcard.
	more bool
}

// meta represents a more complex metadata structure for handling JSON path queries
// that may involve operations, conditions, and logging.
type meta struct {
	// part represents a segment of the deeper query structure.
	part string

	// path represents the full path for the deeper query structure.
	path string

	// pipe represents any operations or transformations applied to the deeper query structure.
	pipe string

	// piped indicates if the deeper query structure has been piped for further processing.
	piped bool

	// more indicates whether there are additional elements or operations in the deeper query.
	more bool

	// arch is a flag indicating some form of structure or architecture flag.
	arch bool

	// logOk indicates whether logging is allowed for this deeper structure.
	logOk bool

	// logKey stores the key used for logging in the deeper structure.
	logKey string

	// query holds the query-related data for the deeper structure.
	query struct {
		// on is a flag that indicates whether the query is active.
		on bool

		// all indicates if the query applies to all elements, not just a single one.
		all bool

		// path is the path to query.
		path string

		// opt specifies an option related to the query.
		opt string

		// val represents the value to search for or match in the query.
		val string
	}
}

// parser holds the state and configuration for parsing JSON data.
type parser struct {
	// json is the raw JSON string that needs to be parsed.
	json string

	// val stores the parsed context value that represents a JSON element.
	val Context

	// pipe holds the pipe operation string to be applied during parsing.
	pipe string

	// piped indicates whether the value has been piped for additional processing.
	piped bool

	// calc indicates whether calculations or transformations should be performed during parsing.
	calc bool

	// lines indicates whether the JSON data should be processed line by line.
	lines bool
}

// stringHeader is a custom struct that mimics the reflect.stringHeader type
// and is used for low-level manipulation of string data in Go.
// It provides access to the internal representation of a string.
type stringHeader struct {
	data unsafe.Pointer // data is a pointer to the underlying data of the string.
	n    int            // n is the length of the string.
}

// sliceHeader is a custom struct that mimics the reflect.sliceHeader type
// and is used for low-level manipulation of slices in Go.
// It provides access to the internal representation of a slice.
type sliceHeader struct {
	data unsafe.Pointer // data is a pointer to the underlying array of the slice.
	n    int            // n is the number of elements in the slice.
	cap  int            // cap is the maximum number of elements the slice can hold before resizing.
}

// sel represents a selection in a JSON path query that specifies a name and path.
type sel struct {
	name string // name represents the name of the selector or key in the JSON path.
	path string //  path represents the full path expression for the selector.
}
