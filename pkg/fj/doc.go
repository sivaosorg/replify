// Package fj (Fast JSON) provides a fast and simple way to retrieve, query,
// and transform values from a JSON document without unmarshalling the entire
// structure into Go types.
//
// fj uses a dot-notation path syntax supporting wildcards, array indexing,
// conditional queries, multi-selectors, pipe operators, and a rich set of
// built-in transformers. Custom transformers can be registered at runtime.
//
// # Path Syntax
//
//	user.name              → object field access
//	roles.0.name           → array index access
//	roles.#.name           → iterate all array elements
//	roles.#(roleName=="Admin").roleId  → conditional query
//	name.@uppercase        → built-in transformer
//	name.@word:upper       → transformer with argument
//	{id,name}              → multi-selector (new object)
//	[id,name]              → multi-selector (new array)
//
// # Basic Usage
//
//	ctx := fj.Get(json, "user.roles.#.roleName")
//	fmt.Println(ctx.String()) // ["Admin","Editor"]
//
// # Custom Transformers
//
//	fj.AddTransformer("word", func(json, arg string) string {
//	    if arg == "upper" { return strings.ToUpper(json) }
//	    return json
//	})
//
// fj is safe for concurrent use by multiple goroutines.
package fj
