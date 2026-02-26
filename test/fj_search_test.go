package test

import (
	"testing"

	"github.com/sivaosorg/replify/pkg/fj"
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
	results := fj.Search(searchTestJSON, "tech")
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
	results := fj.Search(json, "")
	if len(results) != 2 {
		t.Errorf("Search(empty) len = %d; want 2", len(results))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	results := fj.Search(searchTestJSON, "zzznomatch")
	if len(results) != 0 {
		t.Errorf("Search(no match) len = %d; want 0", len(results))
	}
}

func TestSearch_EmptyJSON(t *testing.T) {
	results := fj.Search("", "anything")
	if len(results) != 0 {
		t.Errorf("Search(empty json) len = %d; want 0", len(results))
	}
}

func TestSearch_PartialSubstring(t *testing.T) {
	json := `["Alice","Bob","Albany","Charlie"]`
	results := fj.Search(json, "Al")
	if len(results) != 2 {
		t.Fatalf("Search(Al) len = %d; want 2", len(results))
	}
}

// ///////////////////////////
// Section: SearchByKey
// ///////////////////////////

func TestSearchByKey_SingleKey(t *testing.T) {
	results := fj.SearchByKey(searchTestJSON, "author")
	// 4 books each have an "author" field
	if len(results) != 4 {
		t.Fatalf("SearchByKey(author) len = %d; want 4", len(results))
	}
}

func TestSearchByKey_MultipleKeys(t *testing.T) {
	results := fj.SearchByKey(searchTestJSON, "title", "owner")
	// 4 title fields + 1 owner field
	if len(results) != 5 {
		t.Fatalf("SearchByKey(title,owner) len = %d; want 5", len(results))
	}
}

func TestSearchByKey_NoKeysProvided(t *testing.T) {
	results := fj.SearchByKey(searchTestJSON)
	if len(results) != 0 {
		t.Errorf("SearchByKey(no keys) len = %d; want 0", len(results))
	}
}

func TestSearchByKey_KeyNotFound(t *testing.T) {
	results := fj.SearchByKey(searchTestJSON, "nonexistent")
	if len(results) != 0 {
		t.Errorf("SearchByKey(nonexistent) len = %d; want 0", len(results))
	}
}

func TestSearchByKey_NestedKey(t *testing.T) {
	json := `{"a":{"x":1},"b":{"x":2},"c":{"y":3}}`
	results := fj.SearchByKey(json, "x")
	if len(results) != 2 {
		t.Fatalf("SearchByKey nested len = %d; want 2", len(results))
	}
}

// ///////////////////////////
// Section: Contains
// ///////////////////////////

func TestContains_Exists(t *testing.T) {
	json := `{"msg":"hello world"}`
	if !fj.Contains(json, "msg", "world") {
		t.Error("Contains(world) = false; want true")
	}
}

func TestContains_NotExists(t *testing.T) {
	json := `{"msg":"hello world"}`
	if fj.Contains(json, "msg", "xyz") {
		t.Error("Contains(xyz) = true; want false")
	}
}

func TestContains_MissingPath(t *testing.T) {
	json := `{"msg":"hello"}`
	if fj.Contains(json, "missing", "hello") {
		t.Error("Contains(missing path) = true; want false")
	}
}

func TestContains_NumericValue(t *testing.T) {
	json := `{"code":200}`
	if !fj.Contains(json, "code", "200") {
		t.Error("Contains(200) = false; want true")
	}
}

// ///////////////////////////
// Section: FindPath
// ///////////////////////////

func TestFindPath_ObjectField(t *testing.T) {
	json := `{"user":{"name":"Alice","age":30}}`
	got := fj.FindPath(json, "Alice")
	if got != "user.name" {
		t.Errorf("FindPath(Alice) = %q; want %q", got, "user.name")
	}
}

func TestFindPath_ArrayElement(t *testing.T) {
	json := `{"items":["a","b","c"]}`
	got := fj.FindPath(json, "b")
	if got != "items.1" {
		t.Errorf("FindPath(b) = %q; want %q", got, "items.1")
	}
}

func TestFindPath_NotFound(t *testing.T) {
	got := fj.FindPath(searchTestJSON, "zzznomatch")
	if got != "" {
		t.Errorf("FindPath(not found) = %q; want \"\"", got)
	}
}

func TestFindPath_ReturnsFirst(t *testing.T) {
	// "go" appears at tags.0 and tags.3
	json := `{"tags":["go","json","fast","go"]}`
	got := fj.FindPath(json, "go")
	if got != "tags.0" {
		t.Errorf("FindPath(go) = %q; want %q", got, "tags.0")
	}
}

// ///////////////////////////
// Section: FindPaths
// ///////////////////////////

func TestFindPaths_MultiplePaths(t *testing.T) {
	json := `{"a":"x","b":{"c":"x","d":"y"}}`
	paths := fj.FindPaths(json, "x")
	if len(paths) != 2 {
		t.Fatalf("FindPaths(x) len = %d; want 2", len(paths))
	}
	if paths[0] != "a" || paths[1] != "b.c" {
		t.Errorf("FindPaths(x) = %v; want [a b.c]", paths)
	}
}

func TestFindPaths_NotFound(t *testing.T) {
	paths := fj.FindPaths(searchTestJSON, "zzznomatch")
	if len(paths) != 0 {
		t.Errorf("FindPaths(not found) len = %d; want 0", len(paths))
	}
}

func TestFindPaths_ArrayIndices(t *testing.T) {
	json := `["go","json","go"]`
	paths := fj.FindPaths(json, "go")
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
	if n := fj.Count(searchTestJSON, "store.books"); n != 4 {
		t.Errorf("Count(books) = %d; want 4", n)
	}
}

func TestCount_Scalar(t *testing.T) {
	json := `{"name":"Alice"}`
	if n := fj.Count(json, "name"); n != 1 {
		t.Errorf("Count(scalar) = %d; want 1", n)
	}
}

func TestCount_MissingPath(t *testing.T) {
	if n := fj.Count(searchTestJSON, "missing.path"); n != 0 {
		t.Errorf("Count(missing) = %d; want 0", n)
	}
}

func TestCount_EmptyArray(t *testing.T) {
	json := `{"items":[]}`
	if n := fj.Count(json, "items"); n != 0 {
		t.Errorf("Count(empty array) = %d; want 0", n)
	}
}

// ///////////////////////////
// Section: Sum
// ///////////////////////////

func TestSum_Array(t *testing.T) {
	if s := fj.Sum(searchTestJSON, "ratings"); s != 22.0 {
		t.Errorf("Sum(ratings) = %f; want 22.0", s)
	}
}

func TestSum_SingleNumber(t *testing.T) {
	json := `{"price":9.99}`
	if s := fj.Sum(json, "price"); s != 9.99 {
		t.Errorf("Sum(price) = %f; want 9.99", s)
	}
}

func TestSum_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	if s := fj.Sum(json, "tags"); s != 0 {
		t.Errorf("Sum(no numbers) = %f; want 0", s)
	}
}

func TestSum_MixedArray(t *testing.T) {
	json := `{"data":[1,"ignore",2,null,3]}`
	if s := fj.Sum(json, "data"); s != 6.0 {
		t.Errorf("Sum(mixed) = %f; want 6.0", s)
	}
}

// ///////////////////////////
// Section: Min
// ///////////////////////////

func TestMin_Array(t *testing.T) {
	v, ok := fj.Min(searchTestJSON, "ratings")
	if !ok || v != 1.0 {
		t.Errorf("Min(ratings) = (%f, %v); want (1.0, true)", v, ok)
	}
}

func TestMin_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	_, ok := fj.Min(json, "tags")
	if ok {
		t.Error("Min(no numbers) ok = true; want false")
	}
}

func TestMin_MissingPath(t *testing.T) {
	_, ok := fj.Min(searchTestJSON, "missing")
	if ok {
		t.Error("Min(missing) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Max
// ///////////////////////////

func TestMax_Array(t *testing.T) {
	v, ok := fj.Max(searchTestJSON, "ratings")
	if !ok || v != 5.0 {
		t.Errorf("Max(ratings) = (%f, %v); want (5.0, true)", v, ok)
	}
}

func TestMax_NoNumbers(t *testing.T) {
	json := `{"tags":["a","b"]}`
	_, ok := fj.Max(json, "tags")
	if ok {
		t.Error("Max(no numbers) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Avg
// ///////////////////////////

func TestAvg_Array(t *testing.T) {
	// ratings: 5+3+4+5+1+4 = 22 / 6 ≈ 3.666...
	v, ok := fj.Avg(searchTestJSON, "ratings")
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
	_, ok := fj.Avg(json, "tags")
	if ok {
		t.Error("Avg(no numbers) ok = true; want false")
	}
}

// ///////////////////////////
// Section: Filter
// ///////////////////////////

func TestFilter_KeepsMatching(t *testing.T) {
	results := fj.Filter(searchTestJSON, "ratings", func(ctx fj.Context) bool {
		return ctx.Float64() >= 4
	})
	// ratings >= 4: 5, 4, 5, 4 → 4 values
	if len(results) != 4 {
		t.Errorf("Filter(>=4) len = %d; want 4", len(results))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	results := fj.Filter(searchTestJSON, "ratings", func(ctx fj.Context) bool {
		return ctx.Float64() > 10
	})
	if len(results) != 0 {
		t.Errorf("Filter(>10) len = %d; want 0", len(results))
	}
}

func TestFilter_MissingPath(t *testing.T) {
	results := fj.Filter(searchTestJSON, "missing", func(_ fj.Context) bool { return true })
	if len(results) != 0 {
		t.Errorf("Filter(missing) len = %d; want 0", len(results))
	}
}

func TestFilter_OnObjects(t *testing.T) {
	results := fj.Filter(searchTestJSON, "store.books", func(ctx fj.Context) bool {
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
	ctx := fj.First(searchTestJSON, "ratings", func(c fj.Context) bool {
		return c.Float64() == 5
	})
	if !ctx.Exists() || ctx.Float64() != 5 {
		t.Errorf("First(==5) = %v; want 5", ctx.Float64())
	}
}

func TestFirst_NoMatch(t *testing.T) {
	ctx := fj.First(searchTestJSON, "ratings", func(c fj.Context) bool {
		return c.Float64() > 100
	})
	if ctx.Exists() {
		t.Errorf("First(>100) exists; want zero-value")
	}
}

func TestFirst_MissingPath(t *testing.T) {
	ctx := fj.First(searchTestJSON, "missing", func(_ fj.Context) bool { return true })
	if ctx.Exists() {
		t.Error("First(missing) exists; want zero-value")
	}
}

// ///////////////////////////
// Section: Distinct
// ///////////////////////////

func TestDistinct_DeduplicatesValues(t *testing.T) {
	results := fj.Distinct(searchTestJSON, "tags")
	// original: ["go","json","fast","go","json"] → unique: go, json, fast
	if len(results) != 3 {
		t.Fatalf("Distinct(tags) len = %d; want 3", len(results))
	}
}

func TestDistinct_PreservesOrder(t *testing.T) {
	results := fj.Distinct(searchTestJSON, "tags")
	expected := []string{"go", "json", "fast"}
	for i, r := range results {
		if r.String() != expected[i] {
			t.Errorf("Distinct[%d] = %q; want %q", i, r.String(), expected[i])
		}
	}
}

func TestDistinct_MissingPath(t *testing.T) {
	results := fj.Distinct(searchTestJSON, "missing")
	if len(results) != 0 {
		t.Errorf("Distinct(missing) len = %d; want 0", len(results))
	}
}

func TestDistinct_AllUnique(t *testing.T) {
	json := `{"nums":[1,2,3,4,5]}`
	results := fj.Distinct(json, "nums")
	if len(results) != 5 {
		t.Errorf("Distinct(all unique) len = %d; want 5", len(results))
	}
}

// ///////////////////////////
// Section: Pluck
// ///////////////////////////

func TestPluck_ExtractsFields(t *testing.T) {
	results := fj.Pluck(searchTestJSON, "store.books", "id", "title")
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
	results := fj.Pluck(searchTestJSON, "store.books")
	if len(results) != 0 {
		t.Errorf("Pluck(no fields) len = %d; want 0", len(results))
	}
}

func TestPluck_MissingPath(t *testing.T) {
	results := fj.Pluck(searchTestJSON, "missing", "id")
	if len(results) != 0 {
		t.Errorf("Pluck(missing path) len = %d; want 0", len(results))
	}
}

func TestPluck_MissingFieldIsOmitted(t *testing.T) {
	json := `{"items":[{"a":1},{"a":2,"b":3}]}`
	results := fj.Pluck(json, "items", "a", "b")
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
// Section: SearchMatch (match.Match integration)
// ///////////////////////////

func TestSearchMatch_WildcardStar(t *testing.T) {
	json := `["Alice","Albany","Bob","Alan","Charlie"]`
	results := fj.SearchMatch(json, "Al*")
	if len(results) != 3 {
		t.Fatalf("SearchMatch(Al*) len = %d; want 3", len(results))
	}
}

func TestSearchMatch_WildcardQuestion(t *testing.T) {
	json := `{"a":"cat","b":"bat","c":"car","d":"dog"}`
	results := fj.SearchMatch(json, "?at")
	// matches "cat" and "bat"
	if len(results) != 2 {
		t.Fatalf("SearchMatch(?at) len = %d; want 2", len(results))
	}
}

func TestSearchMatch_Exact(t *testing.T) {
	results := fj.SearchMatch(searchTestJSON, "tech")
	if len(results) != 2 {
		t.Fatalf("SearchMatch(exact) len = %d; want 2", len(results))
	}
}

func TestSearchMatch_StarMatchesAll(t *testing.T) {
	json := `{"a":1,"b":"x","c":true}`
	results := fj.SearchMatch(json, "*")
	if len(results) != 3 {
		t.Errorf("SearchMatch(*) len = %d; want 3", len(results))
	}
}

func TestSearchMatch_NoMatch(t *testing.T) {
	results := fj.SearchMatch(searchTestJSON, "zzz*")
	if len(results) != 0 {
		t.Errorf("SearchMatch(zzz*) len = %d; want 0", len(results))
	}
}

func TestSearchMatch_EmptyJSON(t *testing.T) {
	results := fj.SearchMatch("", "Al*")
	if len(results) != 0 {
		t.Errorf("SearchMatch(empty json) len = %d; want 0", len(results))
	}
}

// ///////////////////////////
// Section: SearchByKeyPattern (match.Match integration)
// ///////////////////////////

func TestSearchByKeyPattern_PrefixWildcard(t *testing.T) {
	json := `{"author":"Donovan","authority":"admin","title":"Go"}`
	results := fj.SearchByKeyPattern(json, "auth*")
	// matches "author" and "authority"
	if len(results) != 2 {
		t.Fatalf("SearchByKeyPattern(auth*) len = %d; want 2", len(results))
	}
}

func TestSearchByKeyPattern_SingleChar(t *testing.T) {
	json := `{"ab":1,"ac":2,"bc":3}`
	results := fj.SearchByKeyPattern(json, "a?")
	// matches "ab" and "ac"
	if len(results) != 2 {
		t.Fatalf("SearchByKeyPattern(a?) len = %d; want 2", len(results))
	}
}

func TestSearchByKeyPattern_StarMatchesAll(t *testing.T) {
	json := `{"x":1,"y":2,"z":3}`
	results := fj.SearchByKeyPattern(json, "*")
	if len(results) != 3 {
		t.Errorf("SearchByKeyPattern(*) len = %d; want 3", len(results))
	}
}

func TestSearchByKeyPattern_Nested(t *testing.T) {
	// All 4 books have "author"
	results := fj.SearchByKeyPattern(searchTestJSON, "author")
	if len(results) != 4 {
		t.Fatalf("SearchByKeyPattern(author) len = %d; want 4", len(results))
	}
}

func TestSearchByKeyPattern_NoMatch(t *testing.T) {
	results := fj.SearchByKeyPattern(searchTestJSON, "zzz*")
	if len(results) != 0 {
		t.Errorf("SearchByKeyPattern(zzz*) len = %d; want 0", len(results))
	}
}

// ///////////////////////////
// Section: ContainsMatch (match.Match integration)
// ///////////////////////////

func TestContainsMatch_Match(t *testing.T) {
	json := `{"email":"alice@example.com"}`
	if !fj.ContainsMatch(json, "email", "*@example.com") {
		t.Error("ContainsMatch(*@example.com) = false; want true")
	}
}

func TestContainsMatch_NoMatch(t *testing.T) {
	json := `{"email":"alice@example.com"}`
	if fj.ContainsMatch(json, "email", "*@other.com") {
		t.Error("ContainsMatch(*@other.com) = true; want false")
	}
}

func TestContainsMatch_MissingPath(t *testing.T) {
	if fj.ContainsMatch(searchTestJSON, "missing", "*") {
		t.Error("ContainsMatch(missing path) = true; want false")
	}
}

func TestContainsMatch_QuestionMark(t *testing.T) {
	json := `{"code":"A1"}`
	if !fj.ContainsMatch(json, "code", "?1") {
		t.Error("ContainsMatch(?1) = false; want true")
	}
}

// ///////////////////////////
// Section: FindPathMatch (match.Match integration)
// ///////////////////////////

func TestFindPathMatch_ObjectField(t *testing.T) {
	json := `{"user":{"name":"Alice","email":"alice@example.com"}}`
	got := fj.FindPathMatch(json, "Al*")
	if got != "user.name" {
		t.Errorf("FindPathMatch(Al*) = %q; want %q", got, "user.name")
	}
}

func TestFindPathMatch_ArrayElement(t *testing.T) {
	json := `{"items":["alpha","beta","almond"]}`
	got := fj.FindPathMatch(json, "al*")
	if got != "items.0" {
		t.Errorf("FindPathMatch(al*) = %q; want %q", got, "items.0")
	}
}

func TestFindPathMatch_NotFound(t *testing.T) {
	got := fj.FindPathMatch(searchTestJSON, "zzz*")
	if got != "" {
		t.Errorf("FindPathMatch(not found) = %q; want \"\"", got)
	}
}

func TestFindPathMatch_ReturnsFirst(t *testing.T) {
	// "Alice" appears before other Al* values in the JSON
	json := `{"a":"Alice","b":"Albany"}`
	got := fj.FindPathMatch(json, "Al*")
	if got != "a" {
		t.Errorf("FindPathMatch(first) = %q; want %q", got, "a")
	}
}

// ///////////////////////////
// Section: FindPathsMatch (match.Match integration)
// ///////////////////////////

func TestFindPathsMatch_Multiple(t *testing.T) {
	json := `{"a":"Alice","b":{"c":"Albany","d":"Bob"}}`
	paths := fj.FindPathsMatch(json, "Al*")
	if len(paths) != 2 {
		t.Fatalf("FindPathsMatch len = %d; want 2", len(paths))
	}
	if paths[0] != "a" || paths[1] != "b.c" {
		t.Errorf("FindPathsMatch = %v; want [a b.c]", paths)
	}
}

func TestFindPathsMatch_NotFound(t *testing.T) {
	paths := fj.FindPathsMatch(searchTestJSON, "zzz*")
	if len(paths) != 0 {
		t.Errorf("FindPathsMatch(not found) len = %d; want 0", len(paths))
	}
}

func TestFindPathsMatch_ArrayIndices(t *testing.T) {
	json := `["apple","apricot","banana","avocado"]`
	paths := fj.FindPathsMatch(json, "a*")
	if len(paths) != 3 {
		t.Fatalf("FindPathsMatch array len = %d; want 3", len(paths))
	}
}

// ///////////////////////////
// Section: CoerceTo (conv.Infer integration)
// ///////////////////////////

func TestCoerceTo_Int(t *testing.T) {
	json := `{"age":30}`
	ctx := fj.Get(json, "age")
	var age int
	if err := fj.CoerceTo(ctx, &age); err != nil {
		t.Fatalf("CoerceTo int: %v", err)
	}
	if age != 30 {
		t.Errorf("CoerceTo int = %d; want 30", age)
	}
}

func TestCoerceTo_Bool(t *testing.T) {
	json := `{"active":true}`
	ctx := fj.Get(json, "active")
	var active bool
	if err := fj.CoerceTo(ctx, &active); err != nil {
		t.Fatalf("CoerceTo bool: %v", err)
	}
	if !active {
		t.Error("CoerceTo bool = false; want true")
	}
}

func TestCoerceTo_String(t *testing.T) {
	json := `{"name":"Alice"}`
	ctx := fj.Get(json, "name")
	var name string
	if err := fj.CoerceTo(ctx, &name); err != nil {
		t.Fatalf("CoerceTo string: %v", err)
	}
	if name != "Alice" {
		t.Errorf("CoerceTo string = %q; want \"Alice\"", name)
	}
}

func TestCoerceTo_Float64(t *testing.T) {
	json := `{"price":34.99}`
	ctx := fj.Get(json, "price")
	var price float64
	if err := fj.CoerceTo(ctx, &price); err != nil {
		t.Fatalf("CoerceTo float64: %v", err)
	}
	if price != 34.99 {
		t.Errorf("CoerceTo float64 = %f; want 34.99", price)
	}
}

func TestCoerceTo_MissingContext(t *testing.T) {
	ctx := fj.Get(`{"a":1}`, "missing")
	var v int
	// Should not panic; error is acceptable
	_ = fj.CoerceTo(ctx, &v)
}

// ///////////////////////////
// Section: CollectFloat64 (conv.Float64 integration)
// ///////////////////////////

func TestCollectFloat64_NumericArray(t *testing.T) {
	json := `{"vals":[1,2,3,4,5]}`
	got := fj.CollectFloat64(json, "vals")
	if len(got) != 5 {
		t.Fatalf("CollectFloat64 len = %d; want 5", len(got))
	}
}

func TestCollectFloat64_StringEncodedNumbers(t *testing.T) {
	// conv.Float64 handles string-encoded numbers
	json := `{"vals":["10","20.5",30]}`
	got := fj.CollectFloat64(json, "vals")
	if len(got) != 3 {
		t.Fatalf("CollectFloat64(string numbers) len = %d; want 3", len(got))
	}
	if got[0] != 10 || got[1] != 20.5 || got[2] != 30 {
		t.Errorf("CollectFloat64(string numbers) = %v", got)
	}
}

func TestCollectFloat64_SkipsNonNumeric(t *testing.T) {
	json := `{"vals":[1,"skip",null,2,true]}`
	got := fj.CollectFloat64(json, "vals")
	// 1 and 2 are numeric; "skip" fails, null fails, true can be coerced to 1 by conv
	// At minimum 1 and 2 must be present
	found1, found2 := false, false
	for _, v := range got {
		if v == 1 {
			found1 = true
		}
		if v == 2 {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("CollectFloat64(mixed) = %v; want at least 1 and 2", got)
	}
}

func TestCollectFloat64_SingleScalar(t *testing.T) {
	json := `{"score":42}`
	got := fj.CollectFloat64(json, "score")
	if len(got) != 1 || got[0] != 42 {
		t.Errorf("CollectFloat64(scalar) = %v; want [42]", got)
	}
}

func TestCollectFloat64_MissingPath(t *testing.T) {
	got := fj.CollectFloat64(searchTestJSON, "missing")
	if len(got) != 0 {
		t.Errorf("CollectFloat64(missing) len = %d; want 0", len(got))
	}
}

// ///////////////////////////
// Section: GroupBy (conv.String integration)
// ///////////////////////////

func TestGroupBy_BasicGrouping(t *testing.T) {
	groups := fj.GroupBy(searchTestJSON, "store.books", "genre")
	if len(groups["tech"]) != 2 {
		t.Errorf("GroupBy tech len = %d; want 2", len(groups["tech"]))
	}
	if len(groups["fiction"]) != 2 {
		t.Errorf("GroupBy fiction len = %d; want 2", len(groups["fiction"]))
	}
}

func TestGroupBy_MissingKeyField(t *testing.T) {
	// Elements without the key go to ""
	json := `{"items":[{"a":1},{"a":2,"genre":"tech"},{"a":3}]}`
	groups := fj.GroupBy(json, "items", "genre")
	if len(groups["tech"]) != 1 {
		t.Errorf("GroupBy tech len = %d; want 1", len(groups["tech"]))
	}
	if len(groups[""]) != 2 {
		t.Errorf("GroupBy empty-key len = %d; want 2", len(groups[""]))
	}
}

func TestGroupBy_MissingPath(t *testing.T) {
	groups := fj.GroupBy(searchTestJSON, "missing", "genre")
	if len(groups) != 0 {
		t.Errorf("GroupBy(missing path) len = %d; want 0", len(groups))
	}
}

func TestGroupBy_NotArray(t *testing.T) {
	json := `{"val":"x"}`
	groups := fj.GroupBy(json, "val", "genre")
	if len(groups) != 0 {
		t.Errorf("GroupBy(not array) len = %d; want 0", len(groups))
	}
}

func TestGroupBy_NumericKey(t *testing.T) {
	json := `{"items":[{"score":1},{"score":2},{"score":1}]}`
	groups := fj.GroupBy(json, "items", "score")
	// conv.String(1.0) should produce "1"
	if len(groups["1"]) != 2 {
		t.Errorf("GroupBy(score=1) len = %d; want 2 (got keys: %v)", len(groups["1"]), groups)
	}
}

// ///////////////////////////
// Section: SortBy (conv integration)
// ///////////////////////////

func TestSortBy_NumericAscending(t *testing.T) {
	json := `{"items":[{"n":3},{"n":1},{"n":2}]}`
	sorted := fj.SortBy(json, "items", "n", true)
	if len(sorted) != 3 {
		t.Fatalf("SortBy len = %d; want 3", len(sorted))
	}
	if sorted[0].Get("n").Int64() != 1 || sorted[1].Get("n").Int64() != 2 || sorted[2].Get("n").Int64() != 3 {
		t.Errorf("SortBy ascending = %v", sorted)
	}
}

func TestSortBy_NumericDescending(t *testing.T) {
	json := `{"items":[{"n":3},{"n":1},{"n":2}]}`
	sorted := fj.SortBy(json, "items", "n", false)
	if len(sorted) != 3 {
		t.Fatalf("SortBy desc len = %d; want 3", len(sorted))
	}
	if sorted[0].Get("n").Int64() != 3 || sorted[1].Get("n").Int64() != 2 || sorted[2].Get("n").Int64() != 1 {
		t.Errorf("SortBy descending = %v", sorted)
	}
}

func TestSortBy_StringField(t *testing.T) {
	sorted := fj.SortBy(searchTestJSON, "store.books", "title", true)
	if len(sorted) != 4 {
		t.Fatalf("SortBy(title) len = %d; want 4", len(sorted))
	}
	// Alphabetically first should be "Clean Code"
	if sorted[0].Get("title").String() != "Clean Code" {
		t.Errorf("SortBy(title) first = %q; want %q", sorted[0].Get("title").String(), "Clean Code")
	}
}

func TestSortBy_PriceDescending(t *testing.T) {
	sorted := fj.SortBy(searchTestJSON, "store.books", "price", false)
	if len(sorted) != 4 {
		t.Fatalf("SortBy(price desc) len = %d; want 4", len(sorted))
	}
	// Most expensive first: 34.99
	if sorted[0].Get("price").Float64() != 34.99 {
		t.Errorf("SortBy(price desc) first = %f; want 34.99", sorted[0].Get("price").Float64())
	}
}

func TestSortBy_ScalarArray(t *testing.T) {
	json := `{"nums":[5,2,8,1,9,3]}`
	sorted := fj.SortBy(json, "nums", "", true)
	if len(sorted) != 6 {
		t.Fatalf("SortBy(scalar) len = %d; want 6", len(sorted))
	}
	if sorted[0].Float64() != 1 {
		t.Errorf("SortBy(scalar asc) first = %f; want 1", sorted[0].Float64())
	}
}

func TestSortBy_MissingPath(t *testing.T) {
	sorted := fj.SortBy(searchTestJSON, "missing", "n", true)
	if len(sorted) != 0 {
		t.Errorf("SortBy(missing) len = %d; want 0", len(sorted))
	}
}
