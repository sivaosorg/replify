package test

import (
	"testing"

	"github.com/sivaosorg/replify"
	"github.com/sivaosorg/replify/pkg/fj"
)

// fjTestJSON is a shared JSON fixture used across fj integration tests.
const fjTestJSON = `{
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
	"tags": ["go","json","fast","go","json"],
	"users": [
		{"id":1,"name":"Alice","role":"admin","active":true},
		{"id":2,"name":"Bob","role":"user","active":false},
		{"id":3,"name":"Charlie","role":"admin","active":true}
	]
}`

// ///////////////////////////
// Section: BodyCtx
// ///////////////////////////

func TestBodyCtx_ReturnsValidContext(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	ctx := r.BodyCtx()
	if ctx.Kind() != fj.String {
		// The body is stored as the raw JSON string, so Kind() == String is expected.
		// Just ensure the context exists.
		if !ctx.Exists() {
			t.Error("BodyCtx() returned a non-existent context")
		}
	}
}

func TestBodyCtx_Nil(t *testing.T) {
	r := replify.WrapOk("ok", nil)
	// Should not panic; context may or may not exist.
	_ = r.BodyCtx()
}

// ///////////////////////////
// Section: QueryBody
// ///////////////////////////

func TestQueryBody_TopLevel(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	ctx := r.QueryBody("store.owner")
	if ctx.String() != "Alice" {
		t.Errorf("QueryBody(store.owner) = %q; want %q", ctx.String(), "Alice")
	}
}

func TestQueryBody_ArrayIndex(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	ctx := r.QueryBody("store.books.0.title")
	if ctx.String() != "The Go Programming Language" {
		t.Errorf("QueryBody(books.0.title) = %q", ctx.String())
	}
}

func TestQueryBody_Missing(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	ctx := r.QueryBody("nonexistent.path")
	if ctx.Exists() {
		t.Error("QueryBody(missing) should not exist")
	}
}

// ///////////////////////////
// Section: QueryBodyMul
// ///////////////////////////

func TestQueryBodyMul_MultiplePaths(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	results := r.QueryBodyMul("store.owner", "store.books.#")
	if len(results) != 2 {
		t.Fatalf("QueryBodyMul len = %d; want 2", len(results))
	}
	if results[0].String() != "Alice" {
		t.Errorf("QueryBodyMul[0] = %q; want Alice", results[0].String())
	}
}

// ///////////////////////////
// Section: ValidBody
// ///////////////////////////

func TestValidBody_ValidJSON(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if !r.ValidBody() {
		t.Error("ValidBody() = false; want true for valid JSON string")
	}
}

func TestValidBody_Struct(t *testing.T) {
	type demo struct {
		Name string `json:"name"`
	}
	r := replify.WrapOk("ok", demo{Name: "Alice"})
	if !r.ValidBody() {
		t.Error("ValidBody() = false for struct body; want true")
	}
}

// ///////////////////////////
// Section: SearchBody
// ///////////////////////////

func TestSearchBody_Keyword(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	hits := r.SearchBody("tech")
	if len(hits) != 2 {
		t.Errorf("SearchBody(tech) len = %d; want 2", len(hits))
	}
}

func TestSearchBody_NoMatch(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	hits := r.SearchBody("zzznomatch")
	if len(hits) != 0 {
		t.Errorf("SearchBody(no match) len = %d; want 0", len(hits))
	}
}

// ///////////////////////////
// Section: SearchBodyMatch
// ///////////////////////////

func TestSearchBodyMatch_Wildcard(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	hits := r.SearchBodyMatch("Al*")
	if len(hits) == 0 {
		t.Error("SearchBodyMatch(Al*) returned no results")
	}
}

// ///////////////////////////
// Section: SearchBodyByKey
// ///////////////////////////

func TestSearchBodyByKey_Author(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	authors := r.SearchBodyByKey("author")
	if len(authors) != 4 {
		t.Errorf("SearchBodyByKey(author) len = %d; want 4", len(authors))
	}
}

// ///////////////////////////
// Section: SearchBodyByKeyPattern
// ///////////////////////////

func TestSearchBodyByKeyPattern_Wildcard(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	hits := r.SearchBodyByKeyPattern("auth*")
	if len(hits) != 4 {
		t.Errorf("SearchBodyByKeyPattern(auth*) len = %d; want 4", len(hits))
	}
}

// ///////////////////////////
// Section: BodyContains / BodyContainsMatch
// ///////////////////////////

func TestBodyContains_True(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if !r.BodyContains("store.owner", "Ali") {
		t.Error("BodyContains(store.owner, Ali) = false; want true")
	}
}

func TestBodyContains_False(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if r.BodyContains("store.owner", "xyz") {
		t.Error("BodyContains(store.owner, xyz) = true; want false")
	}
}

func TestBodyContainsMatch_True(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if !r.BodyContainsMatch("store.owner", "Al*") {
		t.Error("BodyContainsMatch(store.owner, Al*) = false; want true")
	}
}

// ///////////////////////////
// Section: FindBodyPath / FindBodyPaths
// ///////////////////////////

func TestFindBodyPath_Found(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	path := r.FindBodyPath("Rowling")
	if path == "" {
		t.Error("FindBodyPath(Rowling) = empty; want a path")
	}
}

func TestFindBodyPaths_Found(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	paths := r.FindBodyPaths("tech")
	if len(paths) != 2 {
		t.Errorf("FindBodyPaths(tech) len = %d; want 2", len(paths))
	}
}

// ///////////////////////////
// Section: FindBodyPathMatch / FindBodyPathsMatch
// ///////////////////////////

func TestFindBodyPathMatch_Found(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	path := r.FindBodyPathMatch("Row*")
	if path == "" {
		t.Error("FindBodyPathMatch(Row*) = empty; want a path")
	}
}

func TestFindBodyPathsMatch_Multiple(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	paths := r.FindBodyPathsMatch("D*") // Donovan + Dune
	if len(paths) != 2 {
		t.Errorf("FindBodyPathsMatch(D*) len = %d; want 2", len(paths))
	}
}

// ///////////////////////////
// Section: CountBody
// ///////////////////////////

func TestCountBody_Array(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if n := r.CountBody("store.books"); n != 4 {
		t.Errorf("CountBody(store.books) = %d; want 4", n)
	}
}

func TestCountBody_Scalar(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if n := r.CountBody("store.owner"); n != 1 {
		t.Errorf("CountBody(store.owner) = %d; want 1", n)
	}
}

func TestCountBody_Missing(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if n := r.CountBody("missing"); n != 0 {
		t.Errorf("CountBody(missing) = %d; want 0", n)
	}
}

// ///////////////////////////
// Section: SumBody / MinBody / MaxBody / AvgBody
// ///////////////////////////

func TestSumBody(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	if s := r.SumBody("ratings"); s != 22 {
		t.Errorf("SumBody(ratings) = %f; want 22", s)
	}
}

func TestMinBody(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	v, ok := r.MinBody("ratings")
	if !ok || v != 1 {
		t.Errorf("MinBody(ratings) = (%f, %v); want (1, true)", v, ok)
	}
}

func TestMaxBody(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	v, ok := r.MaxBody("ratings")
	if !ok || v != 5 {
		t.Errorf("MaxBody(ratings) = (%f, %v); want (5, true)", v, ok)
	}
}

func TestAvgBody(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	avg, ok := r.AvgBody("ratings")
	if !ok {
		t.Fatal("AvgBody(ratings) ok = false; want true")
	}
	// [5,3,4,5,1,4] sum=22, n=6 → 3.6667
	if avg < 3.6 || avg > 3.7 {
		t.Errorf("AvgBody(ratings) = %f; want ~3.667", avg)
	}
}

// ///////////////////////////
// Section: CollectBodyFloat64
// ///////////////////////////

func TestCollectBodyFloat64(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	nums := r.CollectBodyFloat64("ratings")
	if len(nums) != 6 {
		t.Errorf("CollectBodyFloat64(ratings) len = %d; want 6", len(nums))
	}
}

// ///////////////////////////
// Section: FilterBody
// ///////////////////////////

func TestFilterBody_Genre(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	fiction := r.FilterBody("store.books", func(ctx fj.Context) bool {
		return ctx.Get("genre").String() == "fiction"
	})
	if len(fiction) != 2 {
		t.Errorf("FilterBody(fiction) len = %d; want 2", len(fiction))
	}
}

// ///////////////////////////
// Section: FirstBody
// ///////////////////////////

func TestFirstBody_Found(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	cheap := r.FirstBody("store.books", func(ctx fj.Context) bool {
		return ctx.Get("price").Float64() < 20
	})
	if !cheap.Exists() {
		t.Error("FirstBody(price<20) not found")
	}
	if cheap.Get("title").String() != "Harry Potter" {
		t.Errorf("FirstBody(price<20) title = %q; want Harry Potter", cheap.Get("title").String())
	}
}

func TestFirstBody_NotFound(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	ctx := r.FirstBody("store.books", func(c fj.Context) bool {
		return c.Get("price").Float64() > 1000
	})
	if ctx.Exists() {
		t.Error("FirstBody(price>1000) should not exist")
	}
}

// ///////////////////////////
// Section: DistinctBody
// ///////////////////////////

func TestDistinctBody_Tags(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	unique := r.DistinctBody("tags")
	if len(unique) != 3 {
		t.Errorf("DistinctBody(tags) len = %d; want 3", len(unique))
	}
}

// ///////////////////////////
// Section: PluckBody
// ///////////////////////////

func TestPluckBody_IdName(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	rows := r.PluckBody("users", "id", "name")
	if len(rows) != 3 {
		t.Fatalf("PluckBody(users, id, name) len = %d; want 3", len(rows))
	}
	if !rows[0].Get("id").Exists() || !rows[0].Get("name").Exists() {
		t.Error("PluckBody first row missing id or name")
	}
	if rows[0].Get("role").Exists() {
		t.Error("PluckBody should not include 'role' field")
	}
}

// ///////////////////////////
// Section: GroupByBody
// ///////////////////////////

func TestGroupByBody_Genre(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	groups := r.GroupByBody("store.books", "genre")
	if len(groups["tech"]) != 2 {
		t.Errorf("GroupByBody tech len = %d; want 2", len(groups["tech"]))
	}
	if len(groups["fiction"]) != 2 {
		t.Errorf("GroupByBody fiction len = %d; want 2", len(groups["fiction"]))
	}
}

func TestGroupByBody_Role(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	groups := r.GroupByBody("users", "role")
	if len(groups["admin"]) != 2 {
		t.Errorf("GroupByBody admin len = %d; want 2", len(groups["admin"]))
	}
}

// ///////////////////////////
// Section: SortBodyBy
// ///////////////////////////

func TestSortBodyBy_PriceAscending(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	sorted := r.SortBodyBy("store.books", "price", true)
	if len(sorted) != 4 {
		t.Fatalf("SortBodyBy(price asc) len = %d; want 4", len(sorted))
	}
	if sorted[0].Get("price").Float64() != 12.99 {
		t.Errorf("SortBodyBy first = %f; want 12.99", sorted[0].Get("price").Float64())
	}
}

func TestSortBodyBy_PriceDescending(t *testing.T) {
	r := replify.WrapOk("ok", fjTestJSON)
	sorted := r.SortBodyBy("store.books", "price", false)
	if len(sorted) != 4 {
		t.Fatalf("SortBodyBy(price desc) len = %d; want 4", len(sorted))
	}
	if sorted[0].Get("price").Float64() != 34.99 {
		t.Errorf("SortBodyBy(price desc) first = %f; want 34.99", sorted[0].Get("price").Float64())
	}
}
