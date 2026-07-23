// Package codegen provides a zero-dependency library for generating random
// codes (numbers and letters) with configurable length. It is designed for
// use in fulfillment services and is thread-safe for concurrent use by
// multiple goroutines.
//
// # Basic usage
//
//	g, err := codegen.New(
//	    codegen.WithLength(12),
//	    codegen.WithCharset(codegen.CharsetAlphanumericUpper),
//	    codegen.WithPrefix("ORD-"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	code, err := g.Generate()
//	// code == "ORD-A3BF9KP2XQ17"
package codegen
