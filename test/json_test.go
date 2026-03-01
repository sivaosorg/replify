package test

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/sivaosorg/replify"
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

// ///////////////////////////
// Section: wrapper.JSON() and wrapper.JSONPretty()
// ///////////////////////////

// wrapperJSONFixture is the JSON body used across wrapper JSON tests.
const wrapperJSONFixture = `{"store":{"owner":"Alice"},"ratings":[5,3,4]}`

// TestWrapperJSON_ValidJSONStringBody ensures that JSON() inlines a valid JSON
// string body as a JSON object (not an escaped string literal) in the output.
func TestWrapperJSON_ValidJSONStringBody(t *testing.T) {
w := replify.New().WithBody(wrapperJSONFixture)
out := w.JSON()

var parsed map[string]any
if err := json.Unmarshal([]byte(out), &parsed); err != nil {
t.Fatalf("JSON() output is not valid JSON: %v\noutput: %s", err, out)
}

// The "data" field must be a JSON object, not a string.
data, ok := parsed["data"].(map[string]any)
if !ok {
t.Fatalf("JSON() data field type = %T; want map[string]any", parsed["data"])
}
if _, hasStore := data["store"]; !hasStore {
t.Errorf("JSON() data field missing 'store' key; got keys in data: %v", keysOf(data))
}
}

// TestWrapperJSON_NoEscapedCharsInDataField verifies that JSON() does not
// produce escaped newlines or tabs in the "data" field when the body is valid JSON.
func TestWrapperJSON_NoEscapedCharsInDataField(t *testing.T) {
body := "{\n\t\"key\": \"value\"\n}"
w := replify.New().WithBody(body)
out := w.JSON()
if strings.Contains(out, `\n`) || strings.Contains(out, `\t`) {
t.Errorf("JSON() output contains escaped whitespace; want compact inline JSON\noutput: %s", out)
}
}

// TestWrapperJSON_BytesBody verifies that a []byte body containing valid JSON
// is also inlined correctly (not double-encoded).
func TestWrapperJSON_BytesBody(t *testing.T) {
w := replify.New().WithBody([]byte(wrapperJSONFixture))
out := w.JSON()

var parsed map[string]any
if err := json.Unmarshal([]byte(out), &parsed); err != nil {
t.Fatalf("JSON() (bytes body) is not valid JSON: %v", err)
}
if _, ok := parsed["data"].(map[string]any); !ok {
t.Errorf("JSON() (bytes body) data field type = %T; want map[string]any", parsed["data"])
}
}

// TestWrapperJSON_NonJSONStringBody verifies that a non-JSON string body is
// kept as a JSON string value (not inlined as raw JSON).
func TestWrapperJSON_NonJSONStringBody(t *testing.T) {
w := replify.New().WithBody("hello world")
out := w.JSON()

var parsed map[string]any
if err := json.Unmarshal([]byte(out), &parsed); err != nil {
t.Fatalf("JSON() (non-JSON body) is not valid JSON: %v", err)
}
if s, ok := parsed["data"].(string); !ok || s != "hello world" {
t.Errorf("JSON() (non-JSON body) data = %v (%T); want string %q", parsed["data"], parsed["data"], "hello world")
}
}

// TestWrapperJSON_EmptyBody verifies that JSON() works when no body is set.
func TestWrapperJSON_EmptyBody(t *testing.T) {
w := replify.New()
out := w.JSON()

var parsed map[string]any
if err := json.Unmarshal([]byte(out), &parsed); err != nil {
t.Fatalf("JSON() (empty body) is not valid JSON: %v", err)
}
if _, hasData := parsed["data"]; hasData {
t.Errorf("JSON() (empty body) should not have 'data' field; got: %v", parsed["data"])
}
}

// TestWrapperJSONPretty_ValidJSONStringBody verifies that JSONPretty() produces
// properly indented output with the data field inlined as a JSON object.
func TestWrapperJSONPretty_ValidJSONStringBody(t *testing.T) {
w := replify.New().WithBody(wrapperJSONFixture)
out := w.JSONPretty()

if !strings.Contains(out, "    ") {
t.Errorf("JSONPretty() output is not indented\noutput: %s", out)
}

var parsed map[string]any
if err := json.Unmarshal([]byte(out), &parsed); err != nil {
t.Fatalf("JSONPretty() output is not valid JSON: %v\noutput: %s", err, out)
}
if _, ok := parsed["data"].(map[string]any); !ok {
t.Errorf("JSONPretty() data field type = %T; want map[string]any", parsed["data"])
}
}

// TestWrapperJSONBytes_ValidJSONBody verifies that JSONBytes() returns a non-nil
// byte slice equal to JSON() when the body is valid JSON.
func TestWrapperJSONBytes_ValidJSONBody(t *testing.T) {
w := replify.New().WithBody(wrapperJSONFixture)
b := w.JSONBytes()
if b == nil {
t.Fatal("JSONBytes() = nil; want non-nil byte slice")
}
if string(b) != w.JSON() {
t.Errorf("JSONBytes() != []byte(JSON())")
}
}

// keysOf returns the keys of a map as a slice (helper for test error messages).
func keysOf(m map[string]any) []string {
keys := make([]string, 0, len(m))
for k := range m {
keys = append(keys, k)
}
return keys
}

// ///////////////////////////
// Section: encoding.NormalizeJSON()
// ///////////////////////////

// TestNormalizeJSON_AlreadyValid verifies that already-valid JSON is returned
// unchanged (fast path, no modification).
func TestNormalizeJSON_AlreadyValid(t *testing.T) {
input := `{"key":"value","n":42}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(valid) unexpected error: %v", err)
}
if got != input {
t.Errorf("NormalizeJSON(valid) = %q; want %q (unchanged)", got, input)
}
}

// TestNormalizeJSON_Empty verifies that an empty string returns an error.
func TestNormalizeJSON_Empty(t *testing.T) {
_, err := encoding.NormalizeJSON("")
if err == nil {
t.Error("NormalizeJSON(\"\") expected error; got nil")
}
}

// TestNormalizeJSON_WhitespaceOnly verifies that a whitespace-only string returns an error.
func TestNormalizeJSON_WhitespaceOnly(t *testing.T) {
_, err := encoding.NormalizeJSON("   ")
if err == nil {
t.Error("NormalizeJSON(whitespace) expected error; got nil")
}
}

// TestNormalizeJSON_EscapedStructuralQuotes verifies that literal `\"` sequences
// used as structural key/value delimiters are unescaped to produce valid JSON.
func TestNormalizeJSON_EscapedStructuralQuotes(t *testing.T) {
// Simulate a raw string with \"key\" artifacts.
input := `{\"key\": \"value\"}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(escaped structural quotes) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON result is not valid JSON: %q", got)
}
want := `{"key": "value"}`
if got != want {
t.Errorf("NormalizeJSON = %q; want %q", got, want)
}
}

// TestNormalizeJSON_MixedEscapedKeys verifies a realistic case where only some keys
// are escaped with `\"` and the rest of the object is normal JSON.
func TestNormalizeJSON_MixedEscapedKeys(t *testing.T) {
input := `{\"store\": {"owner": "Alice"}}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(mixed escaped keys) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON result is not valid JSON: %q", got)
}
}

// TestNormalizeJSON_UnfixableInput verifies that an input that cannot be normalized
// to valid JSON returns an error.
func TestNormalizeJSON_UnfixableInput(t *testing.T) {
_, err := encoding.NormalizeJSON("{this is not json at all}")
if err == nil {
t.Error("NormalizeJSON(unfixable) expected error; got nil")
}
}

// TestNormalizeJSON_ObjectArray verifies normalization of a more complex structure
// with escaped outer keys.
func TestNormalizeJSON_ObjectArray(t *testing.T) {
input := `{\"items\": [1, 2, 3], \"count\": 3}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(object+array) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON result is not valid JSON: %q", got)
}
}

// ///////////////////////////
// Section: wrapper.WithNormalizedBody()
// ///////////////////////////

// TestWithNormalizedBody_ValidJSON verifies that a valid JSON string is accepted
// without modification and downstream functions work correctly.
func TestWithNormalizedBody_ValidJSON(t *testing.T) {
input := `{"store":{"owner":"Alice"}}`
w, err := replify.New().WithNormalizedBody(input)
if err != nil {
t.Fatalf("WithNormalizedBody(valid) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody(valid): IsJSONBody() = false; want true")
}
if w.Body() != input {
t.Errorf("WithNormalizedBody(valid): Body() = %v; want %q", w.Body(), input)
}
}

// TestWithNormalizedBody_EscapedQuotes verifies that escaped structural quotes are
// fixed so that IsJSONBody(), JSON(), and JSONPretty() all behave correctly.
func TestWithNormalizedBody_EscapedQuotes(t *testing.T) {
input := `{\"store\": {\"owner\": \"Alice\"}}`
w, err := replify.New().WithNormalizedBody(input)
if err != nil {
t.Fatalf("WithNormalizedBody(escaped) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody(escaped): IsJSONBody() = false; want true")
}
// JSON() output must be valid JSON with data inlined as an object.
jsonOut := w.JSON()
var parsed map[string]any
if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
t.Fatalf("WithNormalizedBody(escaped): JSON() output is not valid JSON: %v\noutput: %s", err, jsonOut)
}
if _, ok := parsed["data"].(map[string]any); !ok {
t.Errorf("WithNormalizedBody(escaped): JSON() data field type = %T; want map[string]any", parsed["data"])
}
}

// TestWithNormalizedBody_InvalidInput verifies that an input that cannot be normalized
// returns an error and leaves the wrapper body unchanged.
func TestWithNormalizedBody_InvalidInput(t *testing.T) {
w := replify.New()
w.WithBody("original")
_, err := w.WithNormalizedBody("{this is not json}")
if err == nil {
t.Error("WithNormalizedBody(invalid) expected error; got nil")
}
}

// TestWithNormalizedBody_Empty verifies that an empty string returns an error.
func TestWithNormalizedBody_Empty(t *testing.T) {
_, err := replify.New().WithNormalizedBody("")
if err == nil {
t.Error("WithNormalizedBody(\"\") expected error; got nil")
}
}

// TestNormalizeJSON_BOM verifies that a leading UTF-8 BOM is stripped so that
// the remaining content is recognized as valid JSON.
func TestNormalizeJSON_BOM(t *testing.T) {
bom := "\xEF\xBB\xBF"
input := bom + `{"key":"value"}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(BOM) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON(BOM) result is not valid JSON: %q", got)
}
if got != `{"key":"value"}` {
t.Errorf("NormalizeJSON(BOM) = %q; want %q", got, `{"key":"value"}`)
}
}

// TestNormalizeJSON_NullBytes verifies that embedded null bytes are removed so
// that the remaining content is recognized as valid JSON.
func TestNormalizeJSON_NullBytes(t *testing.T) {
input := "{\"key\"\x00:\"value\"}"
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(null bytes) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON(null bytes) result is not valid JSON: %q", got)
}
}

// TestNormalizeJSON_TrailingCommaObject verifies that a trailing comma inside an
// object is removed so that the result is valid JSON.
func TestNormalizeJSON_TrailingCommaObject(t *testing.T) {
input := `{"a":1,"b":2,}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(trailing comma object) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON(trailing comma object) result is not valid JSON: %q", got)
}
}

// TestNormalizeJSON_TrailingCommaArray verifies that a trailing comma inside an
// array is removed so that the result is valid JSON.
func TestNormalizeJSON_TrailingCommaArray(t *testing.T) {
input := `[1,2,3,]`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(trailing comma array) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON(trailing comma array) result is not valid JSON: %q", got)
}
}

// TestNormalizeJSON_EscapedQuotesPlusTrailingComma verifies that the cumulative
// multi-pass approach handles both escaped quotes and trailing commas together.
func TestNormalizeJSON_EscapedQuotesPlusTrailingComma(t *testing.T) {
input := `{\"a\":1,\"b\":2,}`
got, err := encoding.NormalizeJSON(input)
if err != nil {
t.Fatalf("NormalizeJSON(escaped+trailing comma) unexpected error: %v", err)
}
if !encoding.IsValidJSON(got) {
t.Errorf("NormalizeJSON(escaped+trailing comma) result is not valid JSON: %q", got)
}
}

// ///////////////////////////
// Section: wrapper.WithNormalizedBody() — additional input types
// ///////////////////////////

// TestWithNormalizedBody_Nil verifies that a nil input returns an error.
func TestWithNormalizedBody_Nil(t *testing.T) {
_, err := replify.New().WithNormalizedBody(nil)
if err == nil {
t.Error("WithNormalizedBody(nil) expected error; got nil")
}
}

// TestWithNormalizedBody_ByteSlice verifies that a []byte containing valid JSON
// is accepted and stored correctly.
func TestWithNormalizedBody_ByteSlice(t *testing.T) {
input := []byte(`{"key":"value"}`)
w, err := replify.New().WithNormalizedBody(input)
if err != nil {
t.Fatalf("WithNormalizedBody([]byte valid) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody([]byte valid): IsJSONBody() = false; want true")
}
}

// TestWithNormalizedBody_ByteSliceEscaped verifies that a []byte with escaped
// structural quotes is normalized and stored correctly.
func TestWithNormalizedBody_ByteSliceEscaped(t *testing.T) {
input := []byte(`{\"key\":\"value\"}`)
w, err := replify.New().WithNormalizedBody(input)
if err != nil {
t.Fatalf("WithNormalizedBody([]byte escaped) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody([]byte escaped): IsJSONBody() = false; want true")
}
}

// TestWithNormalizedBody_Struct verifies that an arbitrary Go struct is marshaled
// to JSON and stored as the body with IsJSONBody() returning true.
func TestWithNormalizedBody_Struct(t *testing.T) {
type payload struct {
Name string `json:"name"`
Age  int    `json:"age"`
}
w, err := replify.New().WithNormalizedBody(payload{Name: "Alice", Age: 30})
if err != nil {
t.Fatalf("WithNormalizedBody(struct) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody(struct): IsJSONBody() = false; want true")
}
jsonOut := w.JSON()
var parsed map[string]any
if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
t.Fatalf("WithNormalizedBody(struct): JSON() output is not valid JSON: %v\noutput: %s", err, jsonOut)
}
if _, ok := parsed["data"].(map[string]any); !ok {
t.Errorf("WithNormalizedBody(struct): JSON() data field type = %T; want map[string]any", parsed["data"])
}
}

// TestWithNormalizedBody_Map verifies that a map[string]any is marshaled to JSON
// and stored as the body with IsJSONBody() returning true.
func TestWithNormalizedBody_Map(t *testing.T) {
input := map[string]any{"x": 1, "y": true}
w, err := replify.New().WithNormalizedBody(input)
if err != nil {
t.Fatalf("WithNormalizedBody(map) unexpected error: %v", err)
}
if !w.IsJSONBody() {
t.Error("WithNormalizedBody(map): IsJSONBody() = false; want true")
}
}
