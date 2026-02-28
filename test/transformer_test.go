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

// /////////////////////////////////////////////////////////
// Section: @project – field projection and renaming
// /////////////////////////////////////////////////////////

func TestExtProject_PickFields(t *testing.T) {
	got := fj.Get(`{"name":"Alice","age":30,"city":"NY"}`, `@project:{"pick":["name","age"]}`).Raw()
	want := `{"name":"Alice","age":30}`
	if got != want {
		t.Errorf("project pick: got %q; want %q", got, want)
	}
}

func TestExtProject_RenameFields(t *testing.T) {
	got := fj.Get(`{"name":"Alice","age":30}`, `@project:{"rename":{"name":"fullName"}}`).Raw()
	want := `{"fullName":"Alice","age":30}`
	if got != want {
		t.Errorf("project rename: got %q; want %q", got, want)
	}
}

func TestExtProject_PickAndRename(t *testing.T) {
	got := fj.Get(`{"name":"Alice","age":30,"city":"NY"}`,
		`@project:{"pick":["name","age"],"rename":{"name":"fullName","age":"years"}}`).Raw()
	want := `{"fullName":"Alice","years":30}`
	if got != want {
		t.Errorf("project pick+rename: got %q; want %q", got, want)
	}
}

func TestExtProject_NonObject_PassThrough(t *testing.T) {
	input := `[1,2,3]`
	got := fj.Get(input, `@project:{"pick":["a"]}`).Raw()
	if got != input {
		t.Errorf("project non-object: got %q; want %q", got, input)
	}
}

func TestExtProject_NoArg_ReturnsAll(t *testing.T) {
	input := `{"a":1,"b":2}`
	got := fj.Get(input, `@project`).Raw()
	if got != input {
		t.Errorf("project no-arg: got %q; want %q", got, input)
	}
}

// /////////////////////////////////////////////////////////
// Section: @filter – conditional array filtering
// /////////////////////////////////////////////////////////

func TestExtFilter_EqualityString(t *testing.T) {
	got := fj.Get(`[{"status":"active"},{"status":"inactive"},{"status":"active"}]`,
		`@filter:{"key":"status","value":"active"}`).Raw()
	want := `[{"status":"active"},{"status":"active"}]`
	if got != want {
		t.Errorf("filter eq string: got %q; want %q", got, want)
	}
}

func TestExtFilter_EqualityNumber(t *testing.T) {
	got := fj.Get(`[{"age":30},{"age":25},{"age":30}]`,
		`@filter:{"key":"age","value":30}`).Raw()
	want := `[{"age":30},{"age":30}]`
	if got != want {
		t.Errorf("filter eq number: got %q; want %q", got, want)
	}
}

func TestExtFilter_OpGt(t *testing.T) {
	got := fj.Get(`[{"age":30},{"age":25},{"age":40}]`,
		`@filter:{"key":"age","op":"gt","value":28}`).Raw()
	want := `[{"age":30},{"age":40}]`
	if got != want {
		t.Errorf("filter gt: got %q; want %q", got, want)
	}
}

func TestExtFilter_OpLte(t *testing.T) {
	got := fj.Get(`[{"score":10},{"score":5},{"score":20}]`,
		`@filter:{"key":"score","op":"lte","value":10}`).Raw()
	want := `[{"score":10},{"score":5}]`
	if got != want {
		t.Errorf("filter lte: got %q; want %q", got, want)
	}
}

func TestExtFilter_OpContains(t *testing.T) {
	got := fj.Get(`[{"name":"Alice"},{"name":"Bob"},{"name":"Alicia"}]`,
		`@filter:{"key":"name","op":"contains","value":"Ali"}`).Raw()
	want := `[{"name":"Alice"},{"name":"Alicia"}]`
	if got != want {
		t.Errorf("filter contains: got %q; want %q", got, want)
	}
}

func TestExtFilter_OpNe(t *testing.T) {
	got := fj.Get(`[{"x":1},{"x":2},{"x":1}]`,
		`@filter:{"key":"x","op":"ne","value":1}`).Raw()
	want := `[{"x":2}]`
	if got != want {
		t.Errorf("filter ne: got %q; want %q", got, want)
	}
}

func TestExtFilter_NonArray_PassThrough(t *testing.T) {
	input := `{"a":1}`
	got := fj.Get(input, `@filter:{"key":"a","value":1}`).Raw()
	if got != input {
		t.Errorf("filter non-array: got %q; want %q", got, input)
	}
}

// /////////////////////////////////////////////////////////
// Section: @default – null / missing field injection
// /////////////////////////////////////////////////////////

func TestExtDefault_NullField(t *testing.T) {
	got := fj.Get(`{"name":"Alice","role":null}`,
		`@default:{"role":"user","active":true}`).Raw()
	want := `{"name":"Alice","role":"user","active":true}`
	if got != want {
		t.Errorf("default null: got %q; want %q", got, want)
	}
}

func TestExtDefault_MissingField(t *testing.T) {
	got := fj.Get(`{"name":"Alice"}`,
		`@default:{"role":"user"}`).Raw()
	want := `{"name":"Alice","role":"user"}`
	if got != want {
		t.Errorf("default missing: got %q; want %q", got, want)
	}
}

func TestExtDefault_ExistingNotOverwritten(t *testing.T) {
	got := fj.Get(`{"name":"Alice","role":"admin"}`,
		`@default:{"role":"user"}`).Raw()
	want := `{"name":"Alice","role":"admin"}`
	if got != want {
		t.Errorf("default existing not overwritten: got %q; want %q", got, want)
	}
}

func TestExtDefault_NonObject_PassThrough(t *testing.T) {
	input := `[1,2,3]`
	got := fj.Get(input, `@default:{"a":1}`).Raw()
	if got != input {
		t.Errorf("default non-object: got %q; want %q", got, input)
	}
}

// /////////////////////////////////////////////////////////
// Section: @coerce – type coercion
// /////////////////////////////////////////////////////////

func TestExtCoerce_NumberToString(t *testing.T) {
	got := fj.Get(`42`, `@coerce:{"to":"string"}`).Raw()
	want := `"42"`
	if got != want {
		t.Errorf("coerce number→string: got %q; want %q", got, want)
	}
}

func TestExtCoerce_StringToNumber(t *testing.T) {
	got := fj.Get(`"99"`, `@coerce:{"to":"number"}`).Raw()
	want := `99`
	if got != want {
		t.Errorf("coerce string→number: got %q; want %q", got, want)
	}
}

func TestExtCoerce_NumberToBool_True(t *testing.T) {
	got := fj.Get(`1`, `@coerce:{"to":"bool"}`).Raw()
	want := `true`
	if got != want {
		t.Errorf("coerce number→bool true: got %q; want %q", got, want)
	}
}

func TestExtCoerce_NumberToBool_False(t *testing.T) {
	got := fj.Get(`0`, `@coerce:{"to":"bool"}`).Raw()
	want := `false`
	if got != want {
		t.Errorf("coerce number→bool false: got %q; want %q", got, want)
	}
}

func TestExtCoerce_ObjectPassThrough(t *testing.T) {
	input := `{"a":1}`
	got := fj.Get(input, `@coerce:{"to":"string"}`).Raw()
	if got != input {
		t.Errorf("coerce object pass-through: got %q; want %q", got, input)
	}
}

// /////////////////////////////////////////////////////////
// Section: @count – count elements
// /////////////////////////////////////////////////////////

func TestExtCount_Array(t *testing.T) {
	got := fj.Get(`[1,2,3]`, `@count`).Raw()
	want := `3`
	if got != want {
		t.Errorf("count array: got %q; want %q", got, want)
	}
}

func TestExtCount_Object(t *testing.T) {
	got := fj.Get(`{"a":1,"b":2}`, `@count`).Raw()
	want := `2`
	if got != want {
		t.Errorf("count object: got %q; want %q", got, want)
	}
}

func TestExtCount_EmptyArray(t *testing.T) {
	got := fj.Get(`[]`, `@count`).Raw()
	want := `0`
	if got != want {
		t.Errorf("count empty array: got %q; want %q", got, want)
	}
}

func TestExtCount_Scalar(t *testing.T) {
	got := fj.Get(`"hello"`, `@count`).Raw()
	want := `0`
	if got != want {
		t.Errorf("count scalar: got %q; want %q", got, want)
	}
}

// /////////////////////////////////////////////////////////
// Section: @first / @last – array head/tail
// /////////////////////////////////////////////////////////

func TestExtFirst_Normal(t *testing.T) {
	got := fj.Get(`[10,20,30]`, `@first`).Raw()
	want := `10`
	if got != want {
		t.Errorf("first: got %q; want %q", got, want)
	}
}

func TestExtFirst_Empty(t *testing.T) {
	got := fj.Get(`[]`, `@first`).Raw()
	want := `null`
	if got != want {
		t.Errorf("first empty: got %q; want %q", got, want)
	}
}

func TestExtFirst_NonArray(t *testing.T) {
	got := fj.Get(`{"a":1}`, `@first`).Raw()
	want := `null`
	if got != want {
		t.Errorf("first non-array: got %q; want %q", got, want)
	}
}

func TestExtLast_Normal(t *testing.T) {
	got := fj.Get(`[10,20,30]`, `@last`).Raw()
	want := `30`
	if got != want {
		t.Errorf("last: got %q; want %q", got, want)
	}
}

func TestExtLast_Empty(t *testing.T) {
	got := fj.Get(`[]`, `@last`).Raw()
	want := `null`
	if got != want {
		t.Errorf("last empty: got %q; want %q", got, want)
	}
}

// /////////////////////////////////////////////////////////
// Section: @sum / @min / @max – numeric aggregation
// /////////////////////////////////////////////////////////

func TestExtSum_Integers(t *testing.T) {
	got := fj.Get(`[1,2,3,4]`, `@sum`).Raw()
	want := `10`
	if got != want {
		t.Errorf("sum integers: got %q; want %q", got, want)
	}
}

func TestExtSum_MixedTypes(t *testing.T) {
	got := fj.Get(`[1.5,2.5,"x",null]`, `@sum`).Raw()
	want := `4`
	if got != want {
		t.Errorf("sum mixed: got %q; want %q", got, want)
	}
}

func TestExtSum_Empty(t *testing.T) {
	got := fj.Get(`[]`, `@sum`).Raw()
	want := `0`
	if got != want {
		t.Errorf("sum empty: got %q; want %q", got, want)
	}
}

func TestExtMin_Integers(t *testing.T) {
	got := fj.Get(`[3,1,4,1,5]`, `@min`).Raw()
	want := `1`
	if got != want {
		t.Errorf("min: got %q; want %q", got, want)
	}
}

func TestExtMin_Empty(t *testing.T) {
	got := fj.Get(`[]`, `@min`).Raw()
	want := `null`
	if got != want {
		t.Errorf("min empty: got %q; want %q", got, want)
	}
}

func TestExtMax_Integers(t *testing.T) {
	got := fj.Get(`[3,1,4,1,5]`, `@max`).Raw()
	want := `5`
	if got != want {
		t.Errorf("max: got %q; want %q", got, want)
	}
}

func TestExtMax_Floats(t *testing.T) {
	got := fj.Get(`[1.1,3.3,2.2]`, `@max`).Raw()
	want := `3.3`
	if got != want {
		t.Errorf("max floats: got %q; want %q", got, want)
	}
}

// /////////////////////////////////////////////////////////
// Section: @pluck – extract field from each array element
// /////////////////////////////////////////////////////////

func TestExtPluck_SimpleField(t *testing.T) {
	got := fj.Get(`[{"name":"Alice","age":30},{"name":"Bob","age":25}]`,
		`@pluck:name`).Raw()
	want := `["Alice","Bob"]`
	if got != want {
		t.Errorf("pluck name: got %q; want %q", got, want)
	}
}

func TestExtPluck_NestedField(t *testing.T) {
	got := fj.Get(`[{"addr":{"city":"NY"}},{"addr":{"city":"LA"}}]`,
		`@pluck:addr.city`).Raw()
	want := `["NY","LA"]`
	if got != want {
		t.Errorf("pluck nested: got %q; want %q", got, want)
	}
}

func TestExtPluck_MissingFieldOmitted(t *testing.T) {
	got := fj.Get(`[{"name":"Alice"},{"age":30}]`,
		`@pluck:name`).Raw()
	want := `["Alice"]`
	if got != want {
		t.Errorf("pluck missing field: got %q; want %q", got, want)
	}
}

func TestExtPluck_NonArray(t *testing.T) {
	got := fj.Get(`{"name":"Alice"}`, `@pluck:name`).Raw()
	want := `[]`
	if got != want {
		t.Errorf("pluck non-array: got %q; want %q", got, want)
	}
}

func TestExtPluck_EmptyArg(t *testing.T) {
	got := fj.Get(`[{"name":"Alice"}]`, `@pluck`).Raw()
	want := `[]`
	if got != want {
		t.Errorf("pluck empty arg: got %q; want %q", got, want)
	}
}
