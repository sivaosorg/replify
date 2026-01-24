package conv

// ///////////////////////////
// Section: Default converter instance
// ///////////////////////////

// defaultConverter is the package-level converter instance used by
// all top-level conversion functions.  It is safe for concurrent use.
var defaultConverter = NewConverter()
