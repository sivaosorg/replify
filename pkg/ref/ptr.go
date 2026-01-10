package ref

// Ptr returns a pointer to the given value.
// This is useful for creating pointers to literals or values inline,
// especially when working with optional fields in structs or API requests.
//
// Example - API Request with optional fields:
//
//	type UpdateUserRequest struct {
//	    Name  *string `json:"name,omitempty"`
//	    Email *string `json:"email,omitempty"`
//	    Age   *int    `json:"age,omitempty"`
//	}
//
//	// Without Ptr - need temporary variables
//	name := "John Doe"
//	email := "john@example.com"
//	req := UpdateUserRequest{
//	    Name:  &name,
//	    Email: &email,
//	}
//
//	// With Ptr - clean and concise
//	req := UpdateUserRequest{
//	    Name:  ref.Ptr("John Doe"),
//	    Email: ref.Ptr("john@example.com"),
//	    Age:   ref.Ptr(30),
//	}
//
// Example - Database insert with optional fields:
//
//	type User struct {
//	    Name  string
//	    Email *string // Can be NULL
//	    Bio   *string // Can be NULL
//	}
//
//	user := User{
//	    Name:  "Jane Smith",
//	    Email: ref.Ptr("jane@example.com"),
//	    Bio:   ref.Ptr("Software Engineer"),
//	}
//
// Example - Partial updates:
//
//	// Only update name and email, leave other fields unchanged
//	update := UpdateUserRequest{
//	    Name:  ref.Ptr("New Name"),
//	    Email: ref.Ptr("new@example.com"),
//	    // Age is nil, won't be updated
//	}
func Ptr[T any](v T) *T {
	return &v
}

// Deref safely dereferences a pointer and returns its value.
// If the pointer is nil, it returns the zero value of type T.
// This is useful for safely accessing pointer values without explicit nil checks.
//
// Example - Safe access to optional config fields:
//
//	type Config struct {
//	    Host    *string
//	    Port    *int
//	    Timeout *int
//	}
//
//	config := Config{
//	    Host: ref.Ptr("localhost"),
//	    // Port and Timeout are nil
//	}
//
//	// Without Deref - manual nil checks
//	var port int
//	if config.Port != nil {
//	    port = *config.Port
//	} else {
//	    port = 0 // zero value
//	}
//
//	// With Deref - clean and safe
//	host := ref.Deref(config.Host)    // "localhost"
//	port := ref.Deref(config.Port)    // 0 (zero value for int)
//	timeout := ref.Deref(config.Timeout) // 0
//
// Example - Processing API response with optional fields:
//
//	type APIResponse struct {
//	    Status  string
//	    Message *string
//	    Data    *ResponseData
//	}
//
//	var resp APIResponse
//	json.Unmarshal(data, &resp)
//
//	message := ref.Deref(resp.Message) // "" if nil
//	if message != "" {
//	    fmt.Println("Message:", message)
//	}
//
// Example - Database query results:
//
//	type User struct {
//	    ID    int
//	    Name  string
//	    Email *string // Can be NULL in database
//	    Age   *int    // Can be NULL in database
//	}
//
//	var user User
//	db.QueryRow("SELECT id, name, email, age FROM users WHERE id = ?", 1).Scan(
//	    &user.ID, &user.Name, &user.Email, &user.Age,
//	)
//
//	email := ref.Deref(user.Email) // "" if NULL
//	age := ref.Deref(user.Age)     // 0 if NULL
//
// Note: If you need a custom default value instead of the zero value,
// use DerefOr instead.
func Deref[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

// DerefOr safely dereferences a pointer and returns its value.
// If the pointer is nil, it returns the provided default value instead of the zero value.
// This is useful when you need specific defaults for nil pointers rather than zero values.
//
// Example - Configuration with custom defaults:
//
//	type ServerConfig struct {
//	    Host     *string
//	    Port     *int
//	    Timeout  *time.Duration
//	    MaxConns *int
//	}
//
//	config := ServerConfig{
//	    Host: ref.Ptr("localhost"),
//	    // Port, Timeout, MaxConns are nil
//	}
//
//	// Apply defaults for nil values
//	host := ref.DerefOr(config.Host, "0.0.0.0")
//	port := ref.DerefOr(config.Port, 8080)
//	timeout := ref.DerefOr(config.Timeout, 30*time.Second)
//	maxConns := ref.DerefOr(config.MaxConns, 100)
//
//	fmt.Printf("Server: %s:%d (timeout: %v, max conns: %d)\n",
//	    host, port, timeout, maxConns)
//	// Output: Server: localhost:8080 (timeout: 30s, max conns: 100)
//
// Example - HTTP client options with defaults:
//
//	type HTTPOptions struct {
//	    Timeout    *time.Duration
//	    MaxRetries *int
//	    UserAgent  *string
//	}
//
//	func NewHTTPClient(opts HTTPOptions) *http.Client {
//	    return &http.Client{
//	        Timeout: ref.DerefOr(opts.Timeout, 30*time.Second),
//	        Transport: &http.Transport{
//	            MaxIdleConns: ref.DerefOr(opts.MaxRetries, 3),
//	        },
//	    }
//	}
//
//	// Use with partial options
//	client := NewHTTPClient(HTTPOptions{
//	    Timeout: ref.Ptr(10 * time.Second),
//	    // MaxRetries will use default: 3
//	})
//
// Example - Database query with default pagination:
//
//	type QueryOptions struct {
//	    Limit  *int
//	    Offset *int
//	    Sort   *string
//	}
//
//	func QueryUsers(opts QueryOptions) []User {
//	    limit := ref.DerefOr(opts.Limit, 10)      // Default: 10
//	    offset := ref.DerefOr(opts.Offset, 0)     // Default: 0
//	    sort := ref.DerefOr(opts.Sort, "id ASC")  // Default: "id ASC"
//
//	    query := fmt.Sprintf("SELECT * FROM users ORDER BY %s LIMIT %d OFFSET %d",
//	        sort, limit, offset)
//	    // Execute query...
//	}
//
//	// Query with custom limit only
//	users := QueryUsers(QueryOptions{
//	    Limit: ref.Ptr(20),
//	    // Offset and Sort use defaults
//	})
//
// Example - Feature flags with defaults:
//
//	type FeatureFlags struct {
//	    EnableCache   *bool
//	    EnableMetrics *bool
//	    EnableTracing *bool
//	}
//
//	flags := FeatureFlags{
//	    EnableCache: ref.Ptr(true),
//	    // EnableMetrics and EnableTracing are nil
//	}
//
//	cache := ref.DerefOr(flags.EnableCache, false)     // true
//	metrics := ref.DerefOr(flags.EnableMetrics, true)  // true (default)
//	tracing := ref.DerefOr(flags.EnableTracing, false) // false (default)
//
// Example - API rate limiting with defaults:
//
//	type RateLimitConfig struct {
//	    RequestsPerSecond *int
//	    BurstSize         *int
//	    Timeout           *time.Duration
//	}
//
//	func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
//	    return &RateLimiter{
//	        rps:     ref.DerefOr(cfg.RequestsPerSecond, 100),
//	        burst:   ref.DerefOr(cfg.BurstSize, 10),
//	        timeout: ref.DerefOr(cfg.Timeout, 5*time.Second),
//	    }
//	}
//
// Note: The difference between Deref and DerefOr:
//   - Deref(ptr) returns zero value if nil: Deref[int](nil) = 0
//   - DerefOr(ptr, 100) returns custom default if nil: DerefOr[int](nil, 100) = 100
func DerefOr[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// IsNil checks if a pointer is nil.
// This is useful for explicit nil checks in conditional logic or validation.
//
// Example - Validating required fields:
//
//	type CreateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func ValidateCreateUser(req CreateUserRequest) error {
//	    if ref.IsNil(req.Name) {
//	        return errors.New("name is required")
//	    }
//	    if ref.IsNil(req.Email) {
//	        return errors.New("email is required")
//	    }
//	    // Age is optional, no need to check
//	    return nil
//	}
//
// Example - Conditional processing:
//
//	type Config struct {
//	    Database *DatabaseConfig
//	    Cache    *CacheConfig
//	    Queue    *QueueConfig
//	}
//
//	func InitializeServices(cfg Config) {
//	    if !ref.IsNil(cfg.Database) {
//	        initDatabase(cfg.Database)
//	    }
//	    if !ref.IsNil(cfg.Cache) {
//	        initCache(cfg.Cache)
//	    }
//	    if !ref.IsNil(cfg.Queue) {
//	        initQueue(cfg.Queue)
//	    }
//	}
//
// Example - Checking optional API response fields:
//
//	type APIResponse struct {
//	    Data  *ResponseData
//	    Error *APIError
//	}
//
//	resp := callAPI()
//	if ref.IsNil(resp.Error) {
//	    // Success case
//	    processData(resp.Data)
//	} else {
//	    // Error case
//	    handleError(resp.Error)
//	}
func IsNil[T any](ptr *T) bool {
	return ptr == nil
}

// IsNotNil checks if a pointer is not nil.
// This is the inverse of IsNil and useful for positive condition checks.
//
// Example - Processing optional fields when present:
//
//	type User struct {
//	    ID      int
//	    Name    string
//	    Email   *string
//	    Phone   *string
//	    Address *Address
//	}
//
//	func DisplayUserInfo(user User) {
//	    fmt.Printf("User: %s\n", user.Name)
//
//	    if ref.IsNotNil(user.Email) {
//	        fmt.Printf("Email: %s\n", *user.Email)
//	    }
//	    if ref.IsNotNil(user.Phone) {
//	        fmt.Printf("Phone: %s\n", *user.Phone)
//	    }
//	    if ref.IsNotNil(user.Address) {
//	        fmt.Printf("Address: %v\n", *user.Address)
//	    }
//	}
//
// Example - Counting non-nil fields:
//
//	func CountProvidedFields(req UpdateUserRequest) int {
//	    count := 0
//	    if ref.IsNotNil(req.Name) {
//	        count++
//	    }
//	    if ref.IsNotNil(req.Email) {
//	        count++
//	    }
//	    if ref.IsNotNil(req.Age) {
//	        count++
//	    }
//	    return count
//	}
func IsNotNil[T any](ptr *T) bool {
	return ptr != nil
}

// Equal checks if two pointers point to equal values.
// Returns true if both are nil, or both are non-nil and their values are equal.
// Returns false if one is nil and the other is not.
//
// Example - Comparing optional configuration:
//
//	oldConfig := Config{
//	    Timeout: ref.Ptr(30 * time.Second),
//	    MaxRetries: ref.Ptr(3),
//	}
//
//	newConfig := Config{
//	    Timeout: ref.Ptr(30 * time.Second),
//	    MaxRetries: ref.Ptr(5),
//	}
//
//	if ref.Equal(oldConfig.Timeout, newConfig.Timeout) {
//	    fmt.Println("Timeout unchanged")
//	}
//	if !ref.Equal(oldConfig.MaxRetries, newConfig.MaxRetries) {
//	    fmt.Println("MaxRetries changed")
//	}
//
// Example - Detecting changes in optional fields:
//
//	type UpdateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func HasChanges(current User, update UpdateUserRequest) bool {
//	    if !ref.Equal(ref.Ptr(current.Name), update.Name) {
//	        return true
//	    }
//	    if !ref.Equal(current.Email, update.Email) {
//	        return true
//	    }
//	    if !ref.Equal(current.Age, update.Age) {
//	        return true
//	    }
//	    return false
//	}
//
// Example - Comparing API responses:
//
//	oldResp := &APIResponse{Status: "success"}
//	newResp := &APIResponse{Status: "success"}
//
//	if ref.Equal(ref.Ptr(oldResp.Status), ref.Ptr(newResp.Status)) {
//	    fmt.Println("Status is the same")
//	}
//
// Note: This function requires type T to be comparable.
// For non-comparable types (slices, maps, functions), use custom comparison logic.
func Equal[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Copy creates a new pointer with a copy of the value.
// If the input pointer is nil, returns nil.
// This is useful for creating independent copies of pointer values.
//
// Example - Deep copying configuration:
//
//	originalConfig := Config{
//	    Host: ref.Ptr("localhost"),
//	    Port: ref.Ptr(8080),
//	}
//
//	// Create independent copy
//	configCopy := Config{
//	    Host: ref.Copy(originalConfig.Host),
//	    Port: ref.Copy(originalConfig.Port),
//	}
//
//	// Modifying copy doesn't affect original
//	*configCopy.Host = "example.com"
//	fmt.Println(*originalConfig.Host) // Still "localhost"
//
// Example - Cloning user data:
//
//	type User struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func CloneUser(user User) User {
//	    return User{
//	        Name:  ref.Copy(user.Name),
//	        Email: ref.Copy(user.Email),
//	        Age:   ref.Copy(user.Age),
//	    }
//	}
//
//	original := User{
//	    Name:  ref.Ptr("John"),
//	    Email: ref.Ptr("john@example.com"),
//	}
//
//	clone := CloneUser(original)
//	*clone.Name = "Jane" // Doesn't affect original
//
// Example - Creating independent request copies:
//
//	originalReq := UpdateUserRequest{
//	    Name: ref.Ptr("Original Name"),
//	    Age:  ref.Ptr(25),
//	}
//
//	// Create copy for retry with modifications
//	retryReq := UpdateUserRequest{
//	    Name: ref.Copy(originalReq.Name),
//	    Age:  ref.Ptr(26), // Different value
//	}
func Copy[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	v := *ptr
	return &v
}

// CoalescePtr returns the first non-nil pointer from the arguments.
// If all pointers are nil, returns nil.
// This is useful for fallback chains and default value selection.
//
// Example - Configuration fallback chain:
//
//	type Config struct {
//	    Host *string
//	    Port *int
//	}
//
//	userConfig := Config{Port: ref.Ptr(3000)}
//	envConfig := Config{Host: ref.Ptr("env-host")}
//	defaultConfig := Config{
//	    Host: ref.Ptr("localhost"),
//	    Port: ref.Ptr(8080),
//	}
//
//	// Use first non-nil value: user > env > default
//	finalHost := ref.CoalescePtr(
//	    userConfig.Host,
//	    envConfig.Host,
//	    defaultConfig.Host,
//	) // Returns envConfig.Host
//
//	finalPort := ref.CoalescePtr(
//	    userConfig.Port,
//	    envConfig.Port,
//	    defaultConfig.Port,
//	) // Returns userConfig.Port
//
// Example - API response with multiple possible sources:
//
//	type Response struct {
//	    PrimaryData   *Data
//	    SecondaryData *Data
//	    CachedData    *Data
//	}
//
//	resp := getResponse()
//	data := ref.CoalescePtr(
//	    resp.PrimaryData,
//	    resp.SecondaryData,
//	    resp.CachedData,
//	) // Use first available data source
//
// Example - User preference with fallbacks:
//
//	type UserSettings struct {
//	    Theme    *string
//	    Language *string
//	}
//
//	userPref := UserSettings{Language: ref.Ptr("en")}
//	orgPref := UserSettings{Theme: ref.Ptr("dark")}
//	systemPref := UserSettings{
//	    Theme:    ref.Ptr("light"),
//	    Language: ref.Ptr("en-US"),
//	}
//
//	theme := ref.CoalescePtr(
//	    userPref.Theme,
//	    orgPref.Theme,
//	    systemPref.Theme,
//	) // "dark"
//
//	language := ref.CoalescePtr(
//	    userPref.Language,
//	    orgPref.Language,
//	    systemPref.Language,
//	) // "en"
func CoalescePtr[T any](ptrs ...*T) *T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return ptr
		}
	}
	return nil
}

// Coalesce returns the first non-nil pointer's value from the arguments.
// If all pointers are nil, returns the zero value of type T.
// This is similar to CoalescePtr but returns the dereferenced value.
//
// Example - Configuration with multiple fallbacks:
//
//	userTimeout := ref.Ptr(5 * time.Second)
//	var envTimeout *time.Duration // nil
//	defaultTimeout := ref.Ptr(30 * time.Second)
//
//	timeout := ref.Coalesce(
//	    userTimeout,
//	    envTimeout,
//	    defaultTimeout,
//	) // Returns 5 * time.Second
//
// Example - String with fallbacks:
//
//	var userTitle *string // nil
//	var orgTitle *string  // nil
//	defaultTitle := ref.Ptr("Untitled")
//
//	title := ref.Coalesce(
//	    userTitle,
//	    orgTitle,
//	    defaultTitle,
//	) // Returns "Untitled"
//
// Example - Database field with fallbacks:
//
//	type User struct {
//	    PreferredName *string
//	    DisplayName   *string
//	    Username      string
//	}
//
//	user := User{
//	    DisplayName: ref.Ptr("John Doe"),
//	    Username:    "john123",
//	}
//
//	name := ref.Coalesce(
//	    user.PreferredName,
//	    user.DisplayName,
//	    ref.Ptr(user.Username),
//	) // Returns "John Doe"
//
// Example - Numeric with zero value fallback:
//
//	var customLimit *int     // nil
//	var userLimit *int       // nil
//	var defaultLimit *int    // nil
//
//	limit := ref.Coalesce(
//	    customLimit,
//	    userLimit,
//	    defaultLimit,
//	) // Returns 0 (zero value for int)
func Coalesce[T any](ptrs ...*T) T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return *ptr
		}
	}
	var zero T
	return zero
}

// ToPtr converts a value to a pointer, but only if the value is not the zero value.
// If the value is the zero value, returns nil.
// This is useful for omitting zero values in API requests or database updates.
//
// Example - API request omitting zero values:
//
//	type UpdateUserRequest struct {
//	    Name  *string `json:"name,omitempty"`
//	    Age   *int    `json:"age,omitempty"`
//	    Score *int    `json:"score,omitempty"`
//	}
//
//	name := "John"
//	age := 0      // Zero value, should be omitted
//	score := 100
//
//	req := UpdateUserRequest{
//	    Name:  ref.ToPtr(name),   // &"John"
//	    Age:   ref.ToPtr(age),    // nil (zero value)
//	    Score: ref.ToPtr(score),  // &100
//	}
//
//	// JSON output: {"name":"John","score":100}
//	// Age is omitted because it's nil
//
// Example - Database partial update:
//
//	func UpdateUser(id int, name string, age int, email string) {
//	    update := map[string]interface{}{}
//
//	    if name := ref.ToPtr(name); name != nil {
//	        update["name"] = *name
//	    }
//	    if age := ref.ToPtr(age); age != nil {
//	        update["age"] = *age
//	    }
//	    if email := ref.ToPtr(email); email != nil {
//	        update["email"] = *email
//	    }
//
//	    // Only update non-zero fields
//	    db.Model(&User{}).Where("id = ?", id).Updates(update)
//	}
//
// Example - Form data validation:
//
//	type SearchFilter struct {
//	    Query     *string
//	    MinPrice  *float64
//	    MaxPrice  *float64
//	    InStock   *bool
//	}
//
//	query := "laptop"
//	minPrice := 0.0    // Zero, will be nil
//	maxPrice := 1000.0
//	inStock := false   // Zero value for bool, will be nil
//
//	filter := SearchFilter{
//	    Query:    ref.ToPtr(query),    // &"laptop"
//	    MinPrice: ref.ToPtr(minPrice), // nil
//	    MaxPrice: ref.ToPtr(maxPrice), // &1000.0
//	    InStock:  ref.ToPtr(inStock),  // nil
//	}
//
// Note: This function requires type T to be comparable to detect zero values.
func ToPtr[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// Map applies a function to the value pointed to by the pointer.
// If the pointer is nil, returns nil.
// This is useful for transforming pointer values without explicit nil checks.
//
// Example - String transformations:
//
//	name := ref.Ptr("john doe")
//
//	// Convert to uppercase
//	upperName := ref.Map(name, strings.ToUpper)
//	fmt.Println(*upperName) // "JOHN DOE"
//
//	// Nil input returns nil
//	var nilName *string
//	result := ref.Map(nilName, strings.ToUpper)
//	fmt.Println(result) // nil
//
// Example - Numeric transformations:
//
//	price := ref.Ptr(100.0)
//
//	// Apply tax
//	priceWithTax := ref.Map(price, func(p float64) float64 {
//	    return p * 1.1 // 10% tax
//	})
//	fmt.Println(*priceWithTax) // 110.0
//
// Example - Type conversions:
//
//	intValue := ref.Ptr(42)
//
//	// Convert int to string
//	stringValue := ref.Map(intValue, func(i int) string {
//	    return fmt.Sprintf("value: %d", i)
//	})
//	fmt.Println(*stringValue) // "value: 42"
//
// Example - Processing optional user input:
//
//	type UserInput struct {
//	    Email *string
//	}
//
//	input := UserInput{
//	    Email: ref.Ptr("  USER@EXAMPLE.COM  "),
//	}
//
//	// Normalize email
//	normalizedEmail := ref.Map(input.Email, func(email string) string {
//	    return strings.ToLower(strings.TrimSpace(email))
//	})
//	fmt.Println(*normalizedEmail) // "user@example.com"
//
// Example - Chain transformations:
//
//	age := ref.Ptr(25)
//	ageGroup := ref.Map(age, func(a int) string {
//	    if a < 18 {
//	        return "minor"
//	    } else if a < 65 {
//	        return "adult"
//	    }
//	    return "senior"
//	})
//	fmt.Println(*ageGroup) // "adult"
func Map[T any, R any](ptr *T, fn func(T) R) *R {
	if ptr == nil {
		return nil
	}
	result := fn(*ptr)
	return &result
}

// Filter returns the pointer if the predicate function returns true, otherwise returns nil.
// If the input pointer is nil, returns nil.
// This is useful for conditional filtering of pointer values.
//
// Example - Validating optional fields:
//
//	age := ref.Ptr(15)
//
//	// Only keep age if >= 18
//	validAge := ref.Filter(age, func(a int) bool {
//	    return a >= 18
//	})
//	fmt.Println(validAge) // nil (15 < 18)
//
//	adult := ref.Ptr(25)
//	validAdult := ref.Filter(adult, func(a int) bool {
//	    return a >= 18
//	})
//	fmt.Println(*validAdult) // 25
//
// Example - Email validation:
//
//	email := ref.Ptr("invalid-email")
//
//	validEmail := ref.Filter(email, func(e string) bool {
//	    return strings.Contains(e, "@")
//	})
//	fmt.Println(validEmail) // nil
//
//	goodEmail := ref.Ptr("user@example.com")
//	validGoodEmail := ref.Filter(goodEmail, func(e string) bool {
//	    return strings.Contains(e, "@")
//	})
//	fmt.Println(*validGoodEmail) // "user@example.com"
//
// Example - Price range filtering:
//
//	type Product struct {
//	    Name  string
//	    Price *float64
//	}
//
//	products := []Product{
//	    {Name: "Cheap", Price: ref.Ptr(10.0)},
//	    {Name: "Expensive", Price: ref.Ptr(1000.0)},
//	    {Name: "Affordable", Price: ref.Ptr(50.0)},
//	}
//
//	for _, p := range products {
//	    // Only products under $100
//	    affordablePrice := ref.Filter(p.Price, func(price float64) bool {
//	        return price < 100
//	    })
//	    if ref.IsNotNil(affordablePrice) {
//	        fmt.Printf("%s: $%.2f\n", p.Name, *affordablePrice)
//	    }
//	}
//	// Output:
//	// Cheap: $10.00
//	// Affordable: $50.00
//
// Example - Filtering with multiple conditions:
//
//	username := ref.Ptr("ab")
//
//	validUsername := ref.Filter(username, func(u string) bool {
//	    return len(u) >= 3 && len(u) <= 20
//	})
//	fmt.Println(validUsername) // nil (too short)
func Filter[T any](ptr *T, predicate func(T) bool) *T {
	if ptr == nil {
		return nil
	}
	if predicate(*ptr) {
		return ptr
	}
	return nil
}

// OrElse returns the pointer if it's not nil, otherwise returns the alternative pointer.
// This is useful for providing fallback pointer values.
//
// Example - Configuration fallback:
//
//	type Config struct {
//	    Timeout *time.Duration
//	}
//
//	userConfig := Config{} // Timeout is nil
//	defaultConfig := Config{
//	    Timeout: ref.Ptr(30 * time.Second),
//	}
//
//	timeout := ref.OrElse(userConfig.Timeout, defaultConfig.Timeout)
//	fmt.Println(*timeout) // 30s
//
// Example - HTTP client options:
//
//	type HTTPOptions struct {
//	    UserAgent *string
//	    Timeout   *time.Duration
//	}
//
//	opts := HTTPOptions{
//	    UserAgent: ref.Ptr("custom-agent"),
//	}
//
//	defaultOpts := HTTPOptions{
//	    UserAgent: ref.Ptr("default-agent"),
//	    Timeout:   ref.Ptr(10 * time.Second),
//	}
//
//	ua := ref.OrElse(opts.UserAgent, defaultOpts.UserAgent)
//	to := ref.OrElse(opts.Timeout, defaultOpts.Timeout)
//
// Example - Response data with cache:
//
//	liveData := fetchFromAPI()  // May return nil
//	cachedData := getFromCache() // Fallback
//
//	data := ref.OrElse(liveData, cachedData)
func OrElse[T any](ptr, alternative *T) *T {
	if ptr != nil {
		return ptr
	}
	return alternative
}

// OrElseGet returns the pointer if it's not nil, otherwise calls the supplier function.
// This is useful when the fallback value is expensive to compute and should only be
// computed if needed (lazy evaluation).
//
// Example - Lazy database query:
//
//	func GetUserName(userID int) *string {
//	    // Try cache first
//	    cached := getFromCache(userID)
//
//	    // Only query database if cache miss
//	    return ref.OrElseGet(cached, func() *string {
//	        return queryDatabase(userID) // Expensive operation
//	    })
//	}
//
// Example - Lazy configuration loading:
//
//	type Config struct {
//	    APIKey *string
//	}
//
//	envConfig := loadFromEnv() // Fast
//
//	config := ref.OrElseGet(envConfig.APIKey, func() *string {
//	    // Only load from file if env var not set
//	    cfg := loadFromFile() // Slow
//	    return cfg.APIKey
//	})
//
// Example - Lazy default value generation:
//
//	sessionID := ref.OrElseGet(existingSession, func() *string {
//	    // Only generate new UUID if no existing session
//	    return ref.Ptr(uuid.New().String())
//	})
//
// Example - Conditional API call:
//
//	cachedPrice := getCachedPrice(productID)
//
//	price := ref.OrElseGet(cachedPrice, func() *float64 {
//	    // Only call pricing API if cache miss
//	    return fetchPriceFromAPI(productID)
//	})
func OrElseGet[T any](ptr *T, supplier func() *T) *T {
	if ptr != nil {
		return ptr
	}
	return supplier()
}

// If executes the consumer function if the pointer is not nil.
// Returns true if the consumer was executed, false otherwise.
// This is useful for side effects on non-nil values.
//
// Example - Logging optional fields:
//
//	type User struct {
//	    Name  string
//	    Email *string
//	    Phone *string
//	}
//
//	user := User{
//	    Name:  "John",
//	    Email: ref.Ptr("john@example.com"),
//	}
//
//	ref.If(user.Email, func(email string) {
//	    log.Printf("Email: %s", email)
//	})
//
//	ref.If(user.Phone, func(phone string) {
//	    log.Printf("Phone: %s", phone) // Not executed
//	})
//
// Example - Sending notifications:
//
//	type Notification struct {
//	    Email *string
//	    SMS   *string
//	    Push  *string
//	}
//
//	notif := Notification{
//	    Email: ref.Ptr("user@example.com"),
//	    SMS:   ref.Ptr("+1234567890"),
//	}
//
//	ref.If(notif.Email, func(email string) {
//	    sendEmail(email, "Hello!")
//	})
//
//	ref.If(notif.SMS, func(phone string) {
//	    sendSMS(phone, "Hello!")
//	})
//
//	ref.If(notif.Push, func(token string) {
//	    sendPush(token, "Hello!") // Not executed
//	})
//
// Example - Processing optional metadata:
//
//	type Event struct {
//	    Name     string
//	    Metadata *map[string]string
//	}
//
//	event := Event{
//	    Name:     "user.login",
//	    Metadata: ref.Ptr(map[string]string{"ip": "127.0.0.1"}),
//	}
//
//	ref.If(event.Metadata, func(meta map[string]string) {
//	    for k, v := range meta {
//	        log.Printf("  %s: %s", k, v)
//	    }
//	})
func If[T any](ptr *T, consumer func(T)) bool {
	if ptr != nil {
		consumer(*ptr)
		return true
	}
	return false
}

// IfElse executes onPresent if pointer is not nil, otherwise executes onAbsent.
// This is useful for branching logic based on pointer presence.
//
// Example - User greeting with fallback:
//
//	user := User{
//	    Name:          "John",
//	    PreferredName: ref.Ptr("Johnny"),
//	}
//
//	ref.IfElse(user.PreferredName,
//	    func(name string) {
//	        fmt.Printf("Hello, %s!", name)
//	    },
//	    func() {
//	        fmt.Printf("Hello, %s!", user.Name)
//	    },
//	)
//
// Example - Configuration with default behavior:
//
//	config := Config{
//	    CustomHandler: ref.Ptr(myHandler),
//	}
//
//	ref.IfElse(config.CustomHandler,
//	    func(handler Handler) {
//	        handler.Handle(request)
//	    },
//	    func() {
//	        defaultHandler.Handle(request)
//	    },
//	)
//
// Example - Caching strategy:
//
//	cached := getFromCache(key)
//
//	ref.IfElse(cached,
//	    func(data Data) {
//	        log.Println("Cache hit")
//	        processData(data)
//	    },
//	    func() {
//	        log.Println("Cache miss, fetching...")
//	        data := fetchFromDB(key)
//	        saveToCache(key, data)
//	        processData(data)
//	    },
//	)
func IfElse[T any](ptr *T, onPresent func(T), onAbsent func()) {
	if ptr != nil {
		onPresent(*ptr)
	} else {
		onAbsent()
	}
}

// Must dereferences the pointer and panics if it's nil.
// This should only be used when you're certain the pointer is not nil,
// such as after validation or in test code.
//
// Example - After validation:
//
//	func ProcessUser(req CreateUserRequest) error {
//	    if ref.IsNil(req.Name) {
//	        return errors.New("name is required")
//	    }
//	    if ref.IsNil(req.Email) {
//	        return errors.New("email is required")
//	    }
//
//	    // Safe to use Must after validation
//	    name := ref.Must(req.Name)
//	    email := ref.Must(req.Email)
//
//	    createUser(name, email)
//	    return nil
//	}
//
// Example - Test code:
//
//	func TestUserCreation(t *testing.T) {
//	    user := createTestUser()
//
//	    // In tests, we expect these to be non-nil
//	    name := ref.Must(user.Name)
//	    email := ref.Must(user.Email)
//
//	    assert.Equal(t, "Test User", name)
//	    assert.Equal(t, "test@example.com", email)
//	}
//
// Example - Configuration that must exist:
//
//	config := loadConfig()
//
//	// These are required, panic if missing
//	dbHost := ref.Must(config.Database.Host)
//	dbPort := ref.Must(config.Database.Port)
//	apiKey := ref.Must(config.API.Key)
//
//	connectDB(dbHost, dbPort, apiKey)
//
// Warning: Only use Must when you're absolutely certain the pointer is not nil.
// For production code, prefer using If, IfElse, or explicit nil checks.
func Must[T any](ptr *T) T {
	if ptr == nil {
		panic("wrapify: attempted to dereference nil pointer")
	}
	return *ptr
}

// MustPtr is like Must but returns a panic-safe pointer dereference.
// Panics with a custom message if the pointer is nil.
//
// Example - With custom error message:
//
//	func GetConfig() Config {
//	    cfg := loadConfigFromFile()
//
//	    return Config{
//	        DBHost: ref.MustPtr(cfg.DBHost, "database host is required"),
//	        DBPort: ref.MustPtr(cfg.DBPort, "database port is required"),
//	        APIKey: ref.MustPtr(cfg.APIKey, "API key is required"),
//	    }
//	}
//
// Example - Critical validation:
//
//	func ProcessPayment(req PaymentRequest) error {
//	    amount := ref.MustPtr(req.Amount, "payment amount is required")
//	    currency := ref.MustPtr(req.Currency, "currency is required")
//	    cardToken := ref.MustPtr(req.CardToken, "card token is required")
//
//	    return processPayment(amount, currency, cardToken)
//	}
func MustPtr[T any](ptr *T, message string) T {
	if ptr == nil {
		panic("wrapify: " + message)
	}
	return *ptr
}

// FlatMap applies a function that returns a pointer to the dereferenced value.
// If the input pointer is nil, returns nil without calling the function.
// This is useful for chaining operations that may return nil.
//
// Example - Nested optional access:
//
//	type Address struct {
//	    City    string
//	    ZipCode *string
//	}
//
//	type User struct {
//	    Name    string
//	    Address *Address
//	}
//
//	user := User{
//	    Name:    "John",
//	    Address: ref.Ptr(Address{City: "NYC"}),
//	}
//
//	zipCode := ref.FlatMap(user.Address, func(addr Address) *string {
//	    return addr.ZipCode
//	})
//	// Returns nil if Address is nil OR ZipCode is nil
//
// Example - Database lookup chain:
//
//	userID := ref.Ptr(123)
//
//	user := ref.FlatMap(userID, func(id int) *User {
//	    return findUserByID(id) // May return nil
//	})
//
//	address := ref.FlatMap(user, func(u User) *Address {
//	    return u.Address // May be nil
//	})
//
// Example - API response chain:
//
//	type APIResponse struct {
//	    Data *ResponseData
//	}
//
//	type ResponseData struct {
//	    User *User
//	}
//
//	resp := callAPI()
//
//	user := ref.FlatMap(resp.Data, func(data ResponseData) *User {
//	    return data.User
//	})
//
// Example - Safe navigation:
//
//	type Config struct {
//	    Database *DatabaseConfig
//	}
//
//	type DatabaseConfig struct {
//	    Connection *ConnectionConfig
//	}
//
//	config := loadConfig()
//
//	connString := ref.FlatMap(config.Database, func(db DatabaseConfig) *string {
//	    return ref.FlatMap(db.Connection, func(conn ConnectionConfig) *string {
//	        return ref.Ptr(conn.String())
//	    })
//	})
func FlatMap[T any, R any](ptr *T, fn func(T) *R) *R {
	if ptr == nil {
		return nil
	}
	return fn(*ptr)
}

// Validate returns the pointer if it passes validation, otherwise returns nil.
// The validator function should return an error if validation fails.
// This is useful for filtering values based on complex validation rules.
//
// Example - Email validation:
//
//	func isValidEmail(email string) error {
//	    if !strings.Contains(email, "@") {
//	        return errors.New("invalid email format")
//	    }
//	    return nil
//	}
//
//	email := ref.Ptr("user@example.com")
//	validEmail := ref.Validate(email, isValidEmail)
//	// Returns email pointer
//
//	badEmail := ref.Ptr("invalid-email")
//	invalidEmail := ref.Validate(badEmail, isValidEmail)
//	// Returns nil
//
// Example - Age validation:
//
//	func isAdult(age int) error {
//	    if age < 18 {
//	        return errors.New("must be 18 or older")
//	    }
//	    return nil
//	}
//
//	age := ref.Ptr(25)
//	validAge := ref.Validate(age, isAdult) // Returns age pointer
//
//	minorAge := ref.Ptr(15)
//	invalidAge := ref.Validate(minorAge, isAdult) // Returns nil
//
// Example - Password strength validation:
//
//	func isStrongPassword(pwd string) error {
//	    if len(pwd) < 8 {
//	        return errors.New("password too short")
//	    }
//	    return nil
//	}
//
//	password := ref.Ptr("StrongP@ss123")
//	validPassword := ref.Validate(password, isStrongPassword)
//
//	if ref.IsNotNil(validPassword) {
//	    // Password is valid, proceed
//	    hashAndStore(*validPassword)
//	}
//
// Example - URL validation:
//
//	func isValidURL(url string) error {
//	    _, err := url.Parse(url)
//	    return err
//	}
//
//	inputURL := ref.Ptr("https://example.com")
//	validURL := ref.Validate(inputURL, isValidURL)
func Validate[T any](ptr *T, validator func(T) error) *T {
	if ptr == nil {
		return nil
	}
	if err := validator(*ptr); err != nil {
		return nil
	}
	return ptr
}

// MapOr applies a function to the pointer value if not nil, otherwise returns default.
// This combines Map and DerefOr in a single operation.
//
// Example - String formatting with default:
//
//	name := ref.Ptr("john")
//
//	formatted := ref.MapOr(name,
//	    func(n string) string {
//	        return strings.ToUpper(n)
//	    },
//	    "ANONYMOUS",
//	)
//	// Returns "JOHN"
//
//	var nilName *string
//	formatted2 := ref.MapOr(nilName,
//	    func(n string) string {
//	        return strings.ToUpper(n)
//	    },
//	    "ANONYMOUS",
//	)
//	// Returns "ANONYMOUS"
//
// Example - Price calculation with default:
//
//	discount := ref.Ptr(0.1) // 10% discount
//
//	finalPrice := ref.MapOr(discount,
//	    func(d float64) float64 {
//	        return 100.0 * (1 - d)
//	    },
//	    100.0, // No discount
//	)
//	// Returns 90.0
//
// Example - Date formatting:
//
//	createdAt := ref.Ptr(time.Now())
//
//	formatted := ref.MapOr(createdAt,
//	    func(t time.Time) string {
//	        return t.Format("2006-01-02")
//	    },
//	    "N/A",
//	)
func MapOr[T any, R any](ptr *T, fn func(T) R, defaultValue R) R {
	if ptr == nil {
		return defaultValue
	}
	return fn(*ptr)
}

// All checks if all pointers in the slice are non-nil.
// Returns true if all pointers are non-nil, false otherwise.
// This is useful for validating that all required fields are present.
//
// Example - Validating required fields:
//
//	type CreateUserRequest struct {
//	    Name     *string
//	    Email    *string
//	    Password *string
//	}
//
//	req := CreateUserRequest{
//	    Name:     ref.Ptr("John"),
//	    Email:    ref.Ptr("john@example.com"),
//	    Password: ref.Ptr("secret"),
//	}
//
//	if ref.All(req.Name, req.Email, req.Password) {
//	    createUser(*req.Name, *req.Email, *req.Password)
//	} else {
//	    return errors.New("all fields are required")
//	}
//
// Example - Configuration validation:
//
//	type DatabaseConfig struct {
//	    Host     *string
//	    Port     *int
//	    Database *string
//	    User     *string
//	    Password *string
//	}
//
//	func ValidateConfig(cfg DatabaseConfig) error {
//	    if !ref.All(cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password) {
//	        return errors.New("all database config fields are required")
//	    }
//	    return nil
//	}
//
// Example - Payment validation:
//
//	type Payment struct {
//	    Amount   *float64
//	    Currency *string
//	    CardToken *string
//	}
//
//	payment := Payment{
//	    Amount:   ref.Ptr(100.0),
//	    Currency: ref.Ptr("USD"),
//	    // CardToken is nil
//	}
//
//	if !ref.All(payment.Amount, payment.Currency, payment.CardToken) {
//	    return errors.New("incomplete payment information")
//	}
func All[T any](ptrs ...*T) bool {
	for _, ptr := range ptrs {
		if ptr == nil {
			return false
		}
	}
	return true
}

// Any checks if any pointer in the slice is non-nil.
// Returns true if at least one pointer is non-nil, false otherwise.
// This is useful for checking if any optional field is provided.
//
// Example - Checking if any contact method provided:
//
//	type ContactInfo struct {
//	    Email *string
//	    Phone *string
//	    SMS   *string
//	}
//
//	contact := ContactInfo{
//	    Phone: ref.Ptr("+1234567890"),
//	}
//
//	if ref.Any(contact.Email, contact.Phone, contact.SMS) {
//	    fmt.Println("At least one contact method provided")
//	} else {
//	    return errors.New("at least one contact method is required")
//	}
//
// Example - Search filter validation:
//
//	type SearchFilter struct {
//	    Query     *string
//	    Category  *string
//	    MinPrice  *float64
//	    MaxPrice  *float64
//	}
//
//	filter := SearchFilter{} // All nil
//
//	if !ref.Any(filter.Query, filter.Category, filter.MinPrice, filter.MaxPrice) {
//	    return errors.New("at least one search criteria is required")
//	}
//
// Example - Update validation:
//
//	type UpdateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Bio   *string
//	}
//
//	req := UpdateUserRequest{} // All nil
//
//	if !ref.Any(req.Name, req.Email, req.Bio) {
//	    return errors.New("at least one field must be updated")
//	}
func Any[T any](ptrs ...*T) bool {
	for _, ptr := range ptrs {
		if ptr != nil {
			return true
		}
	}
	return false
}
