// Package assert provides a small collection of test assertion helpers built
// on top of the standard [testing] package.
//
// The functions in this package are designed to be used directly inside Go
// test functions. Each helper accepts a *testing.T, calls t.Helper so that
// failure messages point to the call site, and reports errors via
// t.Errorf rather than t.Fatalf, allowing a test to record multiple
// failures in a single run.
//
// # Available Assertions
//
//	AssertEqual(t, got, want)    // deep equality via reflect.DeepEqual
//	AssertNil(t, value)         // value must be nil
//	AssertNotNil(t, value)      // value must not be nil
//	AssertTrue(t, condition)    // boolean must be true
//	AssertFalse(t, condition)   // boolean must be false
//
// # Usage
//
//	func TestAdd(t *testing.T) {
//	    result := Add(2, 3)
//	    assert.AssertEqual(t, result, 5)
//	}
//
// assert is used throughout the replify test suite as a lightweight
// alternative to external assertion libraries, keeping test dependencies
// minimal while still providing clear failure messages.
package assert
