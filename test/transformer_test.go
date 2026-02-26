package test

import (
	"strings"
	"testing"

	"github.com/sivaosorg/replify/pkg/fj"
)

// /////////////////////////////////////////
// Struct-based Transformer implementations
// /////////////////////////////////////////

// repeatTransformer is a struct-based Transformer that repeats the input JSON n times,
// demonstrating stateful logic that is impossible with a plain function literal.
type repeatTransformer struct {
	n int
}

func (t *repeatTransformer) Apply(json, arg string) string {
	var b strings.Builder
	for i := 0; i < t.n; i++ {
		b.WriteString(json)
	}
	return b.String()
}

// uppercaseStructTransformer is a stateless struct implementation of Transformer.
type uppercaseStructTransformer struct{}

func (t *uppercaseStructTransformer) Apply(json, arg string) string {
	return strings.ToUpper(json)
}

// /////////////////////
// Section: Transformer
// /////////////////////

// TestTransformerInterface_StructImpl verifies that a struct implementing the
// Transformer interface can be registered and invoked through fj.Get.
func TestTransformerInterface_StructImpl(t *testing.T) {
	fj.AddTransformer("testUpperStruct", &uppercaseStructTransformer{})

	got := fj.Get(`{"name":"alice"}`, "@testUpperStruct").Raw()
	want := `{"NAME":"ALICE"}`
	if got != want {
		t.Errorf("struct Transformer: got %q; want %q", got, want)
	}
}

// TestTransformerInterface_StatefulStructImpl verifies that a stateful struct
// Transformer (one that carries configuration) works correctly.
func TestTransformerInterface_StatefulStructImpl(t *testing.T) {
	fj.AddTransformer("testRepeat2", &repeatTransformer{n: 2})

	// Input is a JSON number; repeating "42" twice yields "4242", also a valid JSON number.
	got := fj.Get(`42`, "@testRepeat2").Raw()
	want := `4242`
	if got != want {
		t.Errorf("stateful Transformer: got %q; want %q", got, want)
	}
}

// TestTransformerInterface_FuncAdapter verifies that TransformerFunc satisfies
// the Transformer interface and behaves identically to a plain function.
func TestTransformerInterface_FuncAdapter(t *testing.T) {
	fj.AddTransformer("testFuncAdapter", fj.TransformerFunc(func(json, arg string) string {
		return strings.ToLower(json)
	}))

	got := fj.Get(`{"NAME":"ALICE"}`, "@testFuncAdapter").Raw()
	want := `{"name":"alice"}`
	if got != want {
		t.Errorf("TransformerFunc adapter: got %q; want %q", got, want)
	}
}

// TestTransformerFunc_ApplyMethod verifies that TransformerFunc.Apply calls
// the underlying function, satisfying the Transformer interface contract.
func TestTransformerFunc_ApplyMethod(t *testing.T) {
	fn := fj.TransformerFunc(func(json, arg string) string {
		return "wrapped:" + json
	})
	got := fn.Apply("input", "")
	want := "wrapped:input"
	if got != want {
		t.Errorf("Apply() = %q; want %q", got, want)
	}
}

// TestTransformerInterface_Overwrite verifies that re-registering a name with a
// different Transformer replaces the previous one.
func TestTransformerInterface_Overwrite(t *testing.T) {
	fj.AddTransformer("testOverwrite", fj.TransformerFunc(func(json, arg string) string {
		return `"first"`
	}))
	fj.AddTransformer("testOverwrite", fj.TransformerFunc(func(json, arg string) string {
		return `"second"`
	}))

	got := fj.Get(`{}`, "@testOverwrite").Raw()
	want := `"second"`
	if got != want {
		t.Errorf("overwrite: got %q; want %q", got, want)
	}
}

// TestIsTransformerRegistered verifies the registry check works for both
// struct-based and function-based registrations.
func TestIsTransformerRegistered(t *testing.T) {
	fj.AddTransformer("testRegistered", &uppercaseStructTransformer{})
	if !fj.IsTransformerRegistered("testRegistered") {
		t.Error("IsTransformerRegistered: expected true for registered name")
	}
	if fj.IsTransformerRegistered("testNotRegistered_xyz") {
		t.Error("IsTransformerRegistered: expected false for unregistered name")
	}
}

// TestBuiltinTransformers_StillWork verifies that existing built-in transformers
// remain functional after the registry was refactored to use the Transformer interface.
func TestBuiltinTransformers_StillWork(t *testing.T) {
	tests := []struct {
		name string
		json string
		path string
		want string
	}{
		{"uppercase", `{"key":"value"}`, "@uppercase", `{"KEY":"VALUE"}`},
		{"lowercase", `{"KEY":"VALUE"}`, "@lowercase", `{"key":"value"}`},
		{"minify", `{"a": 1}`, "@minify", `{"a":1}`},
		{"valid_true", `{"a":1}`, "@valid", "true"},
		{"valid_false", `{bad}`, "@valid", "false"},
		{"reverse_array", `[1,2,3]`, "@reverse", `[3,2,1]`},
		{"keys", `{"x":1,"y":2}`, "@keys", `["x","y"]`},
		{"values", `{"x":1,"y":2}`, "@values", `[1,2]`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fj.Get(tc.json, tc.path).Raw()
			if got != tc.want {
				t.Errorf("Get(%q, %q) = %q; want %q", tc.json, tc.path, got, tc.want)
			}
		})
	}
}
