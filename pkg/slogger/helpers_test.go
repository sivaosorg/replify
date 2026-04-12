package slogger

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// =============================================================================
// Basic Assertions
// =============================================================================

// assertEqual compares two values for deep equality.
func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\n  expected: %#v\n  actual:   %#v", expected, actual)
	}
}

// assertNotEqual checks that two values are not equal.
func assertNotEqual(t *testing.T, notExpected, actual interface{}) {
	t.Helper()
	if reflect.DeepEqual(notExpected, actual) {
		t.Errorf("expected values to differ, but both are: %#v", actual)
	}
}

// assertTrue checks that value is true.
func assertTrue(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Error("expected true, got false")
	}
}

// assertFalse checks that value is false.
func assertFalse(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Error("expected false, got true")
	}
}

// =============================================================================
// Nil Assertions
// =============================================================================

// assertNil checks that value is nil.
func assertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		return
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if !v.IsNil() {
			t.Errorf("expected nil, got %#v", value)
		}
	default:
		t.Errorf("expected nil, got %#v", value)
	}
}

// assertNotNil checks that value is not nil.
func assertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Error("expected non-nil, got nil")
		return
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			t.Error("expected non-nil, got nil")
		}
	}
}

// =============================================================================
// Error Assertions
// =============================================================================

// assertError checks that err is not nil.
func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// assertNoError checks that err is nil.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

// assertErrorIs checks that err matches target using errors.Is.
func assertErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("expected error %v, got %v", target, err)
	}
}

// assertErrorContains checks that err message contains substring.
func assertErrorContains(t *testing.T, err error, substring string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", substring)
		return
	}
	if !strings.Contains(err.Error(), substring) {
		t.Errorf("expected error containing %q, got %q", substring, err.Error())
	}
}

// =============================================================================
// String Assertions
// =============================================================================

// assertContains checks that s contains substring.
func assertContains(t *testing.T, s, substring string) {
	t.Helper()
	if !strings.Contains(s, substring) {
		t.Errorf("expected %q to contain %q", s, substring)
	}
}

// assertNotContains checks that s does not contain substring.
func assertNotContains(t *testing.T, s, substring string) {
	t.Helper()
	if strings.Contains(s, substring) {
		t.Errorf("expected %q to not contain %q", s, substring)
	}
}

// assertEmpty checks that value is empty (string, slice, map).
func assertEmpty(t *testing.T, value interface{}) {
	t.Helper()
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		if v.Len() != 0 {
			t.Errorf("expected empty, got length %d", v.Len())
		}
	default:
		t.Errorf("assertEmpty not supported for type %T", value)
	}
}

// assertNotEmpty checks that value is not empty.
func assertNotEmpty(t *testing.T, value interface{}) {
	t.Helper()
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		if v.Len() == 0 {
			t.Error("expected non-empty, got empty")
		}
	default:
		t.Errorf("assertNotEmpty not supported for type %T", value)
	}
}

// =============================================================================
// Length Assertions
// =============================================================================

// assertLen checks that value has expected length.
func assertLen(t *testing.T, value interface{}, expectedLen int) {
	t.Helper()
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		if v.Len() != expectedLen {
			t.Errorf("expected length %d, got %d", expectedLen, v.Len())
		}
	default:
		t.Errorf("assertLen not supported for type %T", value)
	}
}

// =============================================================================
// Panic Assertions
// =============================================================================

// assertPanics checks that f panics.
func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but did not panic")
		}
	}()
	f()
}

// assertNotPanics checks that f does not panic.
func assertNotPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected no panic, but panicked with: %v", r)
		}
	}()
	f()
}

// assertPanicsWith checks that f panics with specific value.
func assertPanicsWith(t *testing.T, f func(), expected interface{}) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, but did not panic")
			return
		}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected panic with %v, got %v", expected, r)
		}
	}()
	f()
}

// =============================================================================
// Type Assertions
// =============================================================================

// assertType checks that value is of expected type.
func assertType(t *testing.T, expectedType, value interface{}) {
	t.Helper()
	expected := reflect.TypeOf(expectedType)
	actual := reflect.TypeOf(value)
	if expected != actual {
		t.Errorf("expected type %v, got %v", expected, actual)
	}
}

// assertImplements checks that value implements interface.
func assertImplements(t *testing.T, interfacePtr, value interface{}) {
	t.Helper()
	interfaceType := reflect.TypeOf(interfacePtr).Elem()
	valueType := reflect.TypeOf(value)
	if !valueType.Implements(interfaceType) {
		t.Errorf("expected %v to implement %v", valueType, interfaceType)
	}
}

// =============================================================================
// Comparison Assertions
// =============================================================================

// assertGreater checks that a > b.
func assertGreater(t *testing.T, a, b int) {
	t.Helper()
	if a <= b {
		t.Errorf("expected %d > %d", a, b)
	}
}

// assertGreaterOrEqual checks that a >= b.
func assertGreaterOrEqual(t *testing.T, a, b int) {
	t.Helper()
	if a < b {
		t.Errorf("expected %d >= %d", a, b)
	}
}

// assertLess checks that a < b.
func assertLess(t *testing.T, a, b int) {
	t.Helper()
	if a >= b {
		t.Errorf("expected %d < %d", a, b)
	}
}

// assertLessOrEqual checks that a <= b.
func assertLessOrEqual(t *testing.T, a, b int) {
	t.Helper()
	if a > b {
		t.Errorf("expected %d <= %d", a, b)
	}
}

// =============================================================================
// Fatal Variants (stop test immediately)
// =============================================================================

// requireEqual is like assertEqual but fails immediately.
func requireEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n  expected: %#v\n  actual:   %#v", expected, actual)
	}
}

// requireNoError is like assertNoError but fails immediately.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// requireNotNil is like assertNotNil but fails immediately.
func requireNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Fatal("expected non-nil, got nil")
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			t.Fatal("expected non-nil, got nil")
		}
	}
}
