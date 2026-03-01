package test

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/sivaosorg/replify/pkg/encoding"
)

// ///////////////////////////
// Helpers / shared fixtures
// ///////////////////////////

type sampleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ///////////////////////////
// Section: JSON()
// ///////////////////////////

func TestJSON_Nil(t *testing.T) {
	if got := encoding.JSON(nil); got != "" {
		t.Errorf("JSON(nil) = %q; want %q", got, "")
	}
}

func TestJSON_String(t *testing.T) {
	// Strings are returned as-is (no quoting) in jsonSafe.
	if got := encoding.JSON("hello"); got != "hello" {
		t.Errorf("JSON(string) = %q; want %q", got, "hello")
	}
}

func TestJSON_Bool(t *testing.T) {
	if got := encoding.JSON(true); got != "true" {
		t.Errorf("JSON(true) = %q; want %q", got, "true")
	}
	if got := encoding.JSON(false); got != "false" {
		t.Errorf("JSON(false) = %q; want %q", got, "false")
	}
}

func TestJSON_IntegerScalars(t *testing.T) {
	cases := []struct {
		input any
		want  string
	}{
		{int(42), "42"},
		{int8(8), "8"},
		{int16(16), "16"},
		{int32(32), "32"},
		{int64(-64), "-64"},
		{uint(10), "10"},
		{uint8(255), "255"},
		{uint16(1000), "1000"},
		{uint32(100000), "100000"},
		{uint64(999), "999"},
	}
	for _, tc := range cases {
		if got := encoding.JSON(tc.input); got != tc.want {
			t.Errorf("JSON(%v) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestJSON_Float64(t *testing.T) {
	if got := encoding.JSON(float64(3.14)); got != "3.14" {
		t.Errorf("JSON(3.14) = %q; want %q", got, "3.14")
	}
}

func TestJSON_Float32(t *testing.T) {
	got := encoding.JSON(float32(1.5))
	if got != "1.5" {
		t.Errorf("JSON(float32(1.5)) = %q; want %q", got, "1.5")
	}
}

func TestJSON_NaN(t *testing.T) {
	// floatsUseNullForNonFinite == true => "null"
	if got := encoding.JSON(math.NaN()); got != "null" {
		t.Errorf("JSON(NaN) = %q; want %q", got, "null")
	}
}

func TestJSON_Inf(t *testing.T) {
	if got := encoding.JSON(math.Inf(1)); got != "null" {
		t.Errorf("JSON(+Inf) = %q; want %q", got, "null")
	}
	if got := encoding.JSON(math.Inf(-1)); got != "null" {
		t.Errorf("JSON(-Inf) = %q; want %q", got, "null")
	}
}

func TestJSON_Complex128(t *testing.T) {
	got := encoding.JSON(complex(1.0, 2.0))
	if got != `{"real":1,"imag":2}` {
		t.Errorf("JSON(complex128) = %q; want %q", got, `{"real":1,"imag":2}`)
	}
}

func TestJSON_Complex64(t *testing.T) {
	got := encoding.JSON(complex64(complex(3.0, 4.0)))
	if got != `{"real":3,"imag":4}` {
		t.Errorf("JSON(complex64) = %q; want %q", got, `{"real":3,"imag":4}`)
	}
}

func TestJSON_Struct(t *testing.T) {
	s := sampleStruct{Name: "Alice", Age: 30}
	got := encoding.JSON(s)
	want := `{"name":"Alice","age":30}`
	if got != want {
		t.Errorf("JSON(struct) = %q; want %q", got, want)
	}
}

func TestJSON_Map(t *testing.T) {
	m := map[string]int{"a": 1}
	got := encoding.JSON(m)
	if got != `{"a":1}` {
		t.Errorf("JSON(map) = %q; want %q", got, `{"a":1}`)
	}
}

func TestJSON_Slice(t *testing.T) {
	got := encoding.JSON([]int{1, 2, 3})
	if got != `[1,2,3]` {
		t.Errorf("JSON(slice) = %q; want %q", got, `[1,2,3]`)
	}
}

func TestJSON_NilPointer(t *testing.T) {
	var p *sampleStruct
	if got := encoding.JSON(p); got != "null" {
		t.Errorf("JSON(nil pointer) = %q; want %q", got, "null")
	}
}

func TestJSON_NilMap(t *testing.T) {
	var m map[string]int
	if got := encoding.JSON(m); got != "null" {
		t.Errorf("JSON(nil map) = %q; want %q", got, "null")
	}
}

func TestJSON_NilSlice(t *testing.T) {
	var sl []int
	if got := encoding.JSON(sl); got != "null" {
		t.Errorf("JSON(nil slice) = %q; want %q", got, "null")
	}
}

func TestJSON_RawMessageValid(t *testing.T) {
	rm := json.RawMessage(`{"key":"value"}`)
	if got := encoding.JSON(rm); got != `{"key":"value"}` {
		t.Errorf("JSON(RawMessage valid) = %q; want %q", got, `{"key":"value"}`)
	}
}

func TestJSON_RawMessageInvalid(t *testing.T) {
	rm := json.RawMessage(`{invalid}`)
	// invalid raw message => ""
	if got := encoding.JSON(rm); got != "" {
		t.Errorf("JSON(RawMessage invalid) = %q; want %q", got, "")
	}
}

func TestJSON_RawMessageNil(t *testing.T) {
	var rm json.RawMessage
	if got := encoding.JSON(rm); got != "null" {
		t.Errorf("JSON(nil RawMessage) = %q; want %q", got, "null")
	}
}

// ///////////////////////////
// Section: JSONToken()
// ///////////////////////////

func TestJSONToken_Nil(t *testing.T) {
	got, err := encoding.JSONToken(nil)
	if err == nil {
		t.Error("JSONToken(nil) expected error, got nil")
	}
	if got != "" {
		t.Errorf("JSONToken(nil) = %q; want %q", got, "")
	}
}

func TestJSONToken_String(t *testing.T) {
	// Strings are returned as-is (no quoting) in jsonSafeToken.
	got, err := encoding.JSONToken("hello")
	if err != nil {
		t.Fatalf("JSONToken(string) unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("JSONToken(string) = %q; want %q", got, "hello")
	}
}

func TestJSONToken_Bool(t *testing.T) {
	got, err := encoding.JSONToken(true)
	if err != nil {
		t.Fatalf("JSONToken(true) unexpected error: %v", err)
	}
	if got != "true" {
		t.Errorf("JSONToken(true) = %q; want %q", got, "true")
	}
}

func TestJSONToken_IntegerScalars(t *testing.T) {
	cases := []struct {
		input any
		want  string
	}{
		{int(42), "42"},
		{int8(-8), "-8"},
		{int16(16), "16"},
		{int32(32), "32"},
		{int64(64), "64"},
		{uint(7), "7"},
		{uint8(255), "255"},
		{uint16(1000), "1000"},
		{uint32(100000), "100000"},
		{uint64(999), "999"},
	}
	for _, tc := range cases {
		got, err := encoding.JSONToken(tc.input)
		if err != nil {
			t.Errorf("JSONToken(%v) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("JSONToken(%v) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestJSONToken_Float64(t *testing.T) {
	got, err := encoding.JSONToken(float64(2.718))
	if err != nil {
		t.Fatalf("JSONToken(float64) unexpected error: %v", err)
	}
	if got != "2.718" {
		t.Errorf("JSONToken(float64) = %q; want %q", got, "2.718")
	}
}

func TestJSONToken_NaN(t *testing.T) {
	// floatsUseNullForNonFinite == true => "null", no error
	got, err := encoding.JSONToken(math.NaN())
	if err != nil {
		t.Fatalf("JSONToken(NaN) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONToken(NaN) = %q; want %q", got, "null")
	}
}

func TestJSONToken_Inf(t *testing.T) {
	got, err := encoding.JSONToken(math.Inf(1))
	if err != nil {
		t.Fatalf("JSONToken(+Inf) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONToken(+Inf) = %q; want %q", got, "null")
	}
}

func TestJSONToken_Complex128(t *testing.T) {
	got, err := encoding.JSONToken(complex(1.0, 2.0))
	if err != nil {
		t.Fatalf("JSONToken(complex128) unexpected error: %v", err)
	}
	if got != `{"real":1,"imag":2}` {
		t.Errorf("JSONToken(complex128) = %q; want %q", got, `{"real":1,"imag":2}`)
	}
}

func TestJSONToken_Struct(t *testing.T) {
	s := sampleStruct{Name: "Bob", Age: 25}
	got, err := encoding.JSONToken(s)
	if err != nil {
		t.Fatalf("JSONToken(struct) unexpected error: %v", err)
	}
	want := `{"name":"Bob","age":25}`
	if got != want {
		t.Errorf("JSONToken(struct) = %q; want %q", got, want)
	}
}

func TestJSONToken_Slice(t *testing.T) {
	got, err := encoding.JSONToken([]string{"a", "b"})
	if err != nil {
		t.Fatalf("JSONToken(slice) unexpected error: %v", err)
	}
	if got != `["a","b"]` {
		t.Errorf("JSONToken(slice) = %q; want %q", got, `["a","b"]`)
	}
}

func TestJSONToken_NilPointer(t *testing.T) {
	var p *sampleStruct
	got, err := encoding.JSONToken(p)
	if err != nil {
		t.Fatalf("JSONToken(nil pointer) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONToken(nil pointer) = %q; want %q", got, "null")
	}
}

func TestJSONToken_NilMap(t *testing.T) {
	var m map[string]int
	got, err := encoding.JSONToken(m)
	if err != nil {
		t.Fatalf("JSONToken(nil map) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONToken(nil map) = %q; want %q", got, "null")
	}
}

func TestJSONToken_RawMessageValid(t *testing.T) {
	rm := json.RawMessage(`[1,2,3]`)
	got, err := encoding.JSONToken(rm)
	if err != nil {
		t.Fatalf("JSONToken(RawMessage valid) unexpected error: %v", err)
	}
	if got != `[1,2,3]` {
		t.Errorf("JSONToken(RawMessage valid) = %q; want %q", got, `[1,2,3]`)
	}
}

func TestJSONToken_RawMessageInvalid(t *testing.T) {
	rm := json.RawMessage(`{bad json}`)
	_, err := encoding.JSONToken(rm)
	if err == nil {
		t.Error("JSONToken(invalid RawMessage) expected error, got nil")
	}
}

func TestJSONToken_RawMessageNil(t *testing.T) {
	var rm json.RawMessage
	got, err := encoding.JSONToken(rm)
	if err != nil {
		t.Fatalf("JSONToken(nil RawMessage) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONToken(nil RawMessage) = %q; want %q", got, "null")
	}
}

// ///////////////////////////
// Section: JSONPretty()
// ///////////////////////////

func TestJSONPretty_Nil(t *testing.T) {
	if got := encoding.JSONPretty(nil); got != "" {
		t.Errorf("JSONPretty(nil) = %q; want %q", got, "")
	}
}

func TestJSONPretty_String(t *testing.T) {
	// Strings are returned as-is (no quoting).
	if got := encoding.JSONPretty("world"); got != "world" {
		t.Errorf("JSONPretty(string) = %q; want %q", got, "world")
	}
}

func TestJSONPretty_Bool(t *testing.T) {
	if got := encoding.JSONPretty(false); got != "false" {
		t.Errorf("JSONPretty(false) = %q; want %q", got, "false")
	}
}

func TestJSONPretty_Struct_IsIndented(t *testing.T) {
	s := sampleStruct{Name: "Carol", Age: 40}
	got := encoding.JSONPretty(s)
	// Must contain newlines and indentation (4-space indent used by MarshalIndent).
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPretty(struct) = %q; expected indented JSON with newlines", got)
	}
	if !strings.Contains(got, "    ") {
		t.Errorf("JSONPretty(struct) = %q; expected 4-space indentation", got)
	}
	if !strings.Contains(got, `"name"`) || !strings.Contains(got, `"Carol"`) {
		t.Errorf("JSONPretty(struct) = %q; missing expected fields", got)
	}
}

func TestJSONPretty_Map_IsIndented(t *testing.T) {
	m := map[string]int{"x": 99}
	got := encoding.JSONPretty(m)
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPretty(map) = %q; expected indented JSON", got)
	}
}

func TestJSONPretty_Slice(t *testing.T) {
	got := encoding.JSONPretty([]int{1, 2})
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPretty(slice) = %q; expected indented JSON", got)
	}
}

func TestJSONPretty_NilPointer(t *testing.T) {
	var p *sampleStruct
	if got := encoding.JSONPretty(p); got != "null" {
		t.Errorf("JSONPretty(nil pointer) = %q; want %q", got, "null")
	}
}

func TestJSONPretty_RawMessageValid_IsIndented(t *testing.T) {
	rm := json.RawMessage(`{"a":1,"b":2}`)
	got := encoding.JSONPretty(rm)
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPretty(RawMessage) = %q; expected indented JSON", got)
	}
}

func TestJSONPretty_RawMessageInvalid(t *testing.T) {
	rm := json.RawMessage(`{bad}`)
	if got := encoding.JSONPretty(rm); got != "" {
		t.Errorf("JSONPretty(invalid RawMessage) = %q; want %q", got, "")
	}
}

func TestJSONPretty_NaN(t *testing.T) {
	if got := encoding.JSONPretty(math.NaN()); got != "null" {
		t.Errorf("JSONPretty(NaN) = %q; want %q", got, "null")
	}
}

func TestJSONPretty_IntegerScalar(t *testing.T) {
	if got := encoding.JSONPretty(int64(7)); got != "7" {
		t.Errorf("JSONPretty(int64) = %q; want %q", got, "7")
	}
}

// ///////////////////////////
// Section: JSONPrettyToken()
// ///////////////////////////

func TestJSONPrettyToken_Nil(t *testing.T) {
	_, err := encoding.JSONPrettyToken(nil)
	if err == nil {
		t.Error("JSONPrettyToken(nil) expected error, got nil")
	}
}

func TestJSONPrettyToken_String(t *testing.T) {
	got, err := encoding.JSONPrettyToken("test")
	if err != nil {
		t.Fatalf("JSONPrettyToken(string) unexpected error: %v", err)
	}
	if got != "test" {
		t.Errorf("JSONPrettyToken(string) = %q; want %q", got, "test")
	}
}

func TestJSONPrettyToken_Bool(t *testing.T) {
	got, err := encoding.JSONPrettyToken(true)
	if err != nil {
		t.Fatalf("JSONPrettyToken(bool) unexpected error: %v", err)
	}
	if got != "true" {
		t.Errorf("JSONPrettyToken(bool) = %q; want %q", got, "true")
	}
}

func TestJSONPrettyToken_IntegerScalar(t *testing.T) {
	got, err := encoding.JSONPrettyToken(int(100))
	if err != nil {
		t.Fatalf("JSONPrettyToken(int) unexpected error: %v", err)
	}
	if got != "100" {
		t.Errorf("JSONPrettyToken(int) = %q; want %q", got, "100")
	}
}

func TestJSONPrettyToken_Float64(t *testing.T) {
	got, err := encoding.JSONPrettyToken(float64(1.23))
	if err != nil {
		t.Fatalf("JSONPrettyToken(float64) unexpected error: %v", err)
	}
	if got != "1.23" {
		t.Errorf("JSONPrettyToken(float64) = %q; want %q", got, "1.23")
	}
}

func TestJSONPrettyToken_NaN(t *testing.T) {
	got, err := encoding.JSONPrettyToken(math.NaN())
	if err != nil {
		t.Fatalf("JSONPrettyToken(NaN) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONPrettyToken(NaN) = %q; want %q", got, "null")
	}
}

func TestJSONPrettyToken_Struct_IsIndented(t *testing.T) {
	s := sampleStruct{Name: "Dave", Age: 20}
	got, err := encoding.JSONPrettyToken(s)
	if err != nil {
		t.Fatalf("JSONPrettyToken(struct) unexpected error: %v", err)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPrettyToken(struct) = %q; expected indented JSON with newlines", got)
	}
	if !strings.Contains(got, "    ") {
		t.Errorf("JSONPrettyToken(struct) = %q; expected 4-space indentation", got)
	}
	if !strings.Contains(got, `"name"`) || !strings.Contains(got, `"Dave"`) {
		t.Errorf("JSONPrettyToken(struct) = %q; missing expected fields", got)
	}
}

func TestJSONPrettyToken_Slice_IsIndented(t *testing.T) {
	got, err := encoding.JSONPrettyToken([]int{4, 5, 6})
	if err != nil {
		t.Fatalf("JSONPrettyToken(slice) unexpected error: %v", err)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPrettyToken(slice) = %q; expected indented JSON", got)
	}
}

func TestJSONPrettyToken_NilPointer(t *testing.T) {
	var p *sampleStruct
	got, err := encoding.JSONPrettyToken(p)
	if err != nil {
		t.Fatalf("JSONPrettyToken(nil pointer) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONPrettyToken(nil pointer) = %q; want %q", got, "null")
	}
}

func TestJSONPrettyToken_NilMap(t *testing.T) {
	var m map[string]int
	got, err := encoding.JSONPrettyToken(m)
	if err != nil {
		t.Fatalf("JSONPrettyToken(nil map) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONPrettyToken(nil map) = %q; want %q", got, "null")
	}
}

func TestJSONPrettyToken_RawMessageValid_IsIndented(t *testing.T) {
	rm := json.RawMessage(`{"c":3}`)
	got, err := encoding.JSONPrettyToken(rm)
	if err != nil {
		t.Fatalf("JSONPrettyToken(RawMessage) unexpected error: %v", err)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("JSONPrettyToken(RawMessage) = %q; expected indented JSON", got)
	}
}

func TestJSONPrettyToken_RawMessageInvalid(t *testing.T) {
	rm := json.RawMessage(`{not valid}`)
	_, err := encoding.JSONPrettyToken(rm)
	if err == nil {
		t.Error("JSONPrettyToken(invalid RawMessage) expected error, got nil")
	}
}

func TestJSONPrettyToken_RawMessageNil(t *testing.T) {
	var rm json.RawMessage
	got, err := encoding.JSONPrettyToken(rm)
	if err != nil {
		t.Fatalf("JSONPrettyToken(nil RawMessage) unexpected error: %v", err)
	}
	if got != "null" {
		t.Errorf("JSONPrettyToken(nil RawMessage) = %q; want %q", got, "null")
	}
}

func TestJSONPrettyToken_Complex128(t *testing.T) {
	got, err := encoding.JSONPrettyToken(complex(2.0, 3.0))
	if err != nil {
		t.Fatalf("JSONPrettyToken(complex128) unexpected error: %v", err)
	}
	if got != `{"real":2,"imag":3}` {
		t.Errorf("JSONPrettyToken(complex128) = %q; want %q", got, `{"real":2,"imag":3}`)
	}
}
