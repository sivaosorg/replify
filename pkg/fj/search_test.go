package fj

import (
	"testing"
)

// ///////////////////////////
// Shared test JSON documents
// ///////////////////////////

const searchTestJSON = `{
	"store": {
		"books": [
			{"id":1,"title":"The Go Programming Language","author":"Donovan","price":34.99,"genre":"tech"},
			{"id":2,"title":"Clean Code","author":"Martin","price":29.99,"genre":"tech"},
			{"id":3,"title":"Harry Potter","author":"Rowling","price":14.99,"genre":"fiction"},
			{"id":4,"title":"Dune","author":"Herbert","price":12.99,"genre":"fiction"}
		],
		"owner": "Alice"
	},
	"ratings": [5, 3, 4, 5, 1, 4],
	"tags": ["go","json","fast","go","json"]
}`

// ///////////////////////////
// Section: Search
// ///////////////////////////

func TestSearch_FindsMatchingLeaves(t *testing.T) {
	results := Search(searchTestJSON, "tech")
	if len(results) != 2 {
		t.Fatalf("Search() len = %d; want 2", len(results))
	}
	for _, r := range results {
		if r.String() != "tech" {
			t.Errorf("Search() got %q; want \"tech\"", r.String())
		}
	}
}

func TestSearch_EmptyKeywordMatchesAllLeaves(t *testing.T) {
	json := `{"a":1,"b":"x"}`
	results := Search(json, "")
	if len(results) != 2 {
		t.Errorf("Search(empty) len = %d; want 2", len(results))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	results := Search(searchTestJSON, "zzznomatch")
	if len(results) != 0 {
		t.Errorf("Search(no match) len = %d; want 0", len(results))
	}
}

func TestSearch_EmptyJSON(t *testing.T) {
	results := Search("", "anything")
	if len(results) != 0 {
		t.Errorf("Search(empty json) len = %d; want 0", len(results))
	}
}

func TestSearch_PartialSubstring(t *testing.T) {
	json := `["Alice","Bob","Albany","Charlie"]`
	results := Search(json, "Al")
	if len(results) != 2 {
		t.Fatalf("Search(Al) len = %d; want 2", len(results))
	}
}

// ///////////////////////////
// Section: SearchByKey
// ///////////////////////////

func TestSearchByKey_SingleKey(t *testing.T) {
	results := SearchByKey(searchTestJSON, "author")
	// 4 books each have an "author" field
	if len(results) != 4 {
		t.Fatalf("SearchByKey(author) len = %d; want 4", len(results))
	}
}

func TestSearchByKey_MultipleKeys(t *testing.T) {
	results := SearchByKey(searchTestJSON, "title", "owner")
	// 4 title fields + 1 owner field
	if len(results) != 5 {
		t.Fatalf("SearchByKey(title,owner) len = %d; want 5", len(results))
	}
}

func TestSearchByKey_NoKeysProvided(t *testing.T) {
	results := SearchByKey(searchTestJSON)
	if len(results) != 0 {
		t.Errorf("SearchByKey(no keys) len = %d; want 0", len(results))
	}
}

func TestSearchByKey_KeyNotFound(t *testing.T) {
	results := SearchByKey(searchTestJSON, "nonexistent")
	if len(results) != 0 {
		t.Errorf("SearchByKey(nonexistent) len = %d; want 0", len(results))
	}
}

func TestSearchByKey_NestedKey(t *testing.T) {
	json := `{"a":{"x":1},"b":{"x":2},"c":{"y":3}}`
	results := SearchByKey(json, "x")
	if len(results) != 2 {
		t.Fatalf("SearchByKey nested len = %d; want 2", len(results))
	}
}

// ///////////////////////////
// Section: Contains
// ///////////////////////////

func TestContains_Exists(t *testing.T) {
	json := `{"msg":"hello world"}`
	if !Contains(json, "msg", "world") {
		t.Error("Contains(world) = false; want true")
	}
}

func TestContains_NotExists(t *testing.T) {
	json := `{"msg":"hello world"}`
	if Contains(json, "msg", "xyz") {
		t.Error("Contains(xyz) = true; want false")
	}
}

func TestContains_MissingPath(t *testing.T) {
	json := `{"msg":"hello"}`
	if Contains(json, "missing", "hello") {
		t.Error("Contains(missing path) = true; want false")
	}
}

func TestContains_NumericValue(t *testing.T) {
	json := `{"code":200}`
	if !Contains(json, "code", "200") {
		t.Error("Contains(200) = false; want true")
	}
}

// ///////////////////////////
// Section: FindPath
// ///////////////////////////

func TestFindPath_ObjectField(t *testing.T) {
	json := `{"user":{"name":"Alice","age":30}}`
	got := FindPath(json, "Alice")
	if got != "user.name" {
		t.Errorf("FindPath(Alice) = %q; want %q", got, "user.name")
	}
}

func TestFindPath_ArrayElement(t *testing.T) {
	json := `{"items":["a","b","c"]}`
	got := FindPath(json, "b")
	if got != "items.1" {
		t.Errorf("FindPath(b) = %q; want %q", got, "items.1")
	}
}

func TestFindPath_NotFound(t *testing.T) {
	got := FindPath(searchTestJSON, "zzznomatch")
	if got != "" {
		t.Errorf("FindPath(not found) = %q; want \"\"", got)
	}
}

func TestFindPath_ReturnsFirst(t *testing.T) {
	// "go" appears at tags.0 and tags.3
	json := `{"tags":["go","json","fast","go"]}`
	got := FindPath(json, "go")
	if got != "tags.0" {
		t.Errorf("FindPath(go) = %q; want %q", got, "tags.0")
	}
}

// ///////////////////////////
// Section: FindPaths
// ///////////////////////////

func TestFindPaths_MultiplePaths(t *testing.T) {
	json := `{"a":"x","b":{"c":"x","d":"y"}}`
	paths := FindPaths(json, "x")
	if len(paths) != 2 {
		t.Fatalf("FindPaths(x) len = %d; want 2", len(paths))
	}
	if paths[0] != "a" || paths[1] != "b.c" {
		t.Errorf("FindPaths(x) = %v; want [a b.c]", paths)
	}
}

func TestFindPaths_NotFound(t *testing.T) {
	paths := FindPaths(searchTestJSON, "zzznomatch")
	if len(paths) != 0 {
		t.Errorf("FindPaths(not found) len = %d; want 0", len(paths))
	}
}

func TestFindPaths_ArrayIndices(t *testing.T) {
	json := `["go","json","go"]`
	paths := FindPaths(json, "go")
	if len(paths) != 2 {
		t.Fatalf("FindPaths array indices len = %d; want 2", len(paths))
	}
	if paths[0] != "0" || paths[1] != "2" {
		t.Errorf("FindPaths array = %v; want [0 2]", paths)
	}
}

// ///////////////////////////
// Section: Count
// ///////////////////////////

func TestCount_Array(t *testing.T) {
	if n := Count(searchTestJSON, "store.books"); n != 4 {
		t.Errorf("Count(books) = %d; want 4", n)
	}
}

func TestCount_Scalar(t *testing.T) {
	json := `{"name":"Alice"}`
	if n := Count(json, "name"); n != 1 {
		t.Errorf("Count(scalar) = %d; want 1", n)
	}
}

func TestCount_MissingPath(t *testing.T) {
	if n := Count(searchTestJSON, "missing.path"); n != 0 {
		t.Errorf("Count(missing) = %d; want 0", n)
	}
}

func TestCount_EmptyArray(t *testing.T) {
	json := `{"items":[]}`
	if n := Count(json, "items"); n != 0 {
		t.Errorf("Count(empty array) = %d; want 0", n)
	}
}

// ///////////////////////////
// Section: Sum
// ///////////////////////////

func TestSum_Array(t *testing.T) {
	if s := Sum(searchTestJSON, "ratings"); s != 22.0 {
		t.Errorf("Sum(ratings) = %f; want 22.0", s)
	}
}

func TestSum_SingleNumber(t *testing.T) {
	json := `{"price":9.99}`
	if s := Sum(json, "price"); s != 9.99 {
		t.Errorf("Sum(price) = %f; want 9.99", s)
	}
}

func TestSum_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	if s := Sum(json, "tags"); s != 0 {
		t.Errorf("Sum(no numbers) = %f; want 0", s)
	}
}

func TestSum_MixedArray(t *testing.T) {
	json := `{"data":[1,"ignore",2,null,3]}`
	if s := Sum(json, "data"); s != 6.0 {
		t.Errorf("Sum(mixed) = %f; want 6.0", s)
	}
}

// ///////////////////////////
// Section: Min
// ///////////////////////////

func TestMin_Array(t *testing.T) {
	v, ok := Min(searchTestJSON, "ratings")
	if !ok || v != 1.0 {
		t.Errorf("Min(ratings) = (%f, %v); want (1.0, true)", v, ok)
	}
}

func TestMin_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	_, ok := Min(json, "tags")
	if ok {
		t.Error("Min(no numbers) ok = true; want false")
	}
}

func TestMin_MissingPath(t *testing.T) {
	_, ok := Min(searchTestJSON, "missing")
	if ok {
		t.Error("Min(missing) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Max
// ///////////////////////////

func TestMax_Array(t *testing.T) {
	v, ok := Max(searchTestJSON, "ratings")
	if !ok || v != 5.0 {
		t.Errorf("Max(ratings) = (%f, %v); want (5.0, true)", v, ok)
	}
}

func TestMax_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	_, ok := Max(json, "tags")
	if ok {
		t.Error("Max(no numbers) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Avg
// ///////////////////////////

func TestAvg_Array(t *testing.T) {
	// ratings: 5+3+4+5+1+4 = 22 / 6 ≈ 3.666...
	v, ok := Avg(searchTestJSON, "ratings")
	if !ok {
		t.Fatal("Avg(ratings) ok = false; want true")
	}
	expected := 22.0 / 6.0
	if v != expected {
		t.Errorf("Avg(ratings) = %f; want %f", v, expected)
	}
}

func TestAvg_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	_, ok := Avg(json, "tags")
	if ok {
		t.Error("Avg(no numbers) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Filter
// ///////////////////////////

func TestFilter_KeepsMatching(t *testing.T) {
	results := Filter(searchTestJSON, "ratings", func(ctx Context) bool {
		return ctx.Float64() >= 4
	})
	// ratings >= 4: 5, 4, 5, 4 → 4 values
	if len(results) != 4 {
		t.Errorf("Filter(>=4) len = %d; want 4", len(results))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	results := Filter(searchTestJSON, "ratings", func(ctx Context) bool {
		return ctx.Float64() > 10
	})
	if len(results) != 0 {
		t.Errorf("Filter(>10) len = %d; want 0", len(results))
	}
}

func TestFilter_MissingPath(t *testing.T) {
	results := Filter(searchTestJSON, "missing", func(_ Context) bool { return true })
	if len(results) != 0 {
		t.Errorf("Filter(missing) len = %d; want 0", len(results))
	}
}

func TestFilter_OnObjects(t *testing.T) {
	results := Filter(searchTestJSON, "store.books", func(ctx Context) bool {
		return ctx.Get("genre").String() == "fiction"
	})
	// 2 fiction books
	if len(results) != 2 {
		t.Errorf("Filter(fiction) len = %d; want 2", len(results))
	}
}

// ///////////////////////////
// Section: First
// ///////////////////////////

func TestFirst_ReturnsFirstMatch(t *testing.T) {
	ctx := First(searchTestJSON, "ratings", func(c Context) bool {
		return c.Float64() == 5
	})
	if !ctx.Exists() || ctx.Float64() != 5 {
		t.Errorf("First(==5) = %v; want 5", ctx.Float64())
	}
}

func TestFirst_NoMatch(t *testing.T) {
	ctx := First(searchTestJSON, "ratings", func(c Context) bool {
		return c.Float64() > 100
	})
	if ctx.Exists() {
		t.Errorf("First(>100) exists; want zero-value")
	}
}

func TestFirst_MissingPath(t *testing.T) {
	ctx := First(searchTestJSON, "missing", func(_ Context) bool { return true })
	if ctx.Exists() {
		t.Error("First(missing) exists; want zero-value")
	}
}

// ///////////////////////////
// Section: Distinct
// ///////////////////////////

func TestDistinct_DeduplicatesValues(t *testing.T) {
	results := Distinct(searchTestJSON, "tags")
	// original: ["go","json","fast","go","json"] → unique: go, json, fast
	if len(results) != 3 {
		t.Fatalf("Distinct(tags) len = %d; want 3", len(results))
	}
}

func TestDistinct_PreservesOrder(t *testing.T) {
	results := Distinct(searchTestJSON, "tags")
	expected := []string{"go", "json", "fast"}
	for i, r := range results {
		if r.String() != expected[i] {
			t.Errorf("Distinct[%d] = %q; want %q", i, r.String(), expected[i])
		}
	}
}

func TestDistinct_MissingPath(t *testing.T) {
	results := Distinct(searchTestJSON, "missing")
	if len(results) != 0 {
		t.Errorf("Distinct(missing) len = %d; want 0", len(results))
	}
}

func TestDistinct_AllUnique(t *testing.T) {
	json := `{"nums":[1,2,3,4,5]}`
	results := Distinct(json, "nums")
	if len(results) != 5 {
		t.Errorf("Distinct(all unique) len = %d; want 5", len(results))
	}
}

// ///////////////////////////
// Section: Pluck
// ///////////////////////////

func TestPluck_ExtractsFields(t *testing.T) {
	results := Pluck(searchTestJSON, "store.books", "id", "title")
	if len(results) != 4 {
		t.Fatalf("Pluck len = %d; want 4", len(results))
	}
	// Each result must be an object with exactly id and title
	for i, r := range results {
		if !r.IsObject() {
			t.Errorf("Pluck[%d] not an object", i)
		}
		if !r.Get("id").Exists() || !r.Get("title").Exists() {
			t.Errorf("Pluck[%d] missing id or title", i)
		}
		if r.Get("author").Exists() {
			t.Errorf("Pluck[%d] contains author but should not", i)
		}
	}
}

func TestPluck_NoFieldsProvided(t *testing.T) {
	results := Pluck(searchTestJSON, "store.books")
	if len(results) != 0 {
		t.Errorf("Pluck(no fields) len = %d; want 0", len(results))
	}
}

func TestPluck_MissingPath(t *testing.T) {
	results := Pluck(searchTestJSON, "missing", "id")
	if len(results) != 0 {
		t.Errorf("Pluck(missing path) len = %d; want 0", len(results))
	}
}

func TestPluck_MissingFieldIsOmitted(t *testing.T) {
	json := `{"items":[{"a":1},{"a":2,"b":3}]}`
	results := Pluck(json, "items", "a", "b")
	if len(results) != 2 {
		t.Fatalf("Pluck partial fields len = %d; want 2", len(results))
	}
	// First item has only "a"
	if results[0].Get("b").Exists() {
		t.Error("Pluck first item should not have 'b'")
	}
	// Second item has both "a" and "b"
	if !results[1].Get("b").Exists() {
		t.Error("Pluck second item should have 'b'")
	}
}

// ///////////////////////////
// Section: joinPath / itoa helpers
// ///////////////////////////

func TestJoinPath(t *testing.T) {
	tests := []struct {
		prefix, segment, want string
	}{
		{"", "a", "a"},
		{"a", "b", "a.b"},
		{"a.b", "c", "a.b.c"},
		{"", "0", "0"},
	}
	for _, tt := range tests {
		got := joinPath(tt.prefix, tt.segment)
		if got != tt.want {
			t.Errorf("joinPath(%q,%q) = %q; want %q", tt.prefix, tt.segment, got, tt.want)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "10"},
		{99, "99"},
		{100, "100"},
		{12345, "12345"},
	}
	for _, tt := range tests {
		got := itoa(tt.in)
		if got != tt.want {
			t.Errorf("itoa(%d) = %q; want %q", tt.in, got, tt.want)
		}
	}
}
