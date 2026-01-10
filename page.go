package replify

// WithPage sets the page number for the `pagination` instance.
//
// This function updates the `page` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the page number to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithPage(v int) *pagination {
	if v < 1 {
		v = 1
	}
	p.page = v
	p.calculate()
	return p
}

// WithPerPage sets the number of items per page for the `pagination` instance.
//
// This function updates the `perPage` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
// Validates that perPage is >= 1, defaults to 10 if invalid value.
//
// Parameters:
//   - `v`: An integer representing the number of items per page to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithPerPage(v int) *pagination {
	if v < 1 {
		v = 10
	}
	p.perPage = v
	p.calculate()
	return p
}

// WithTotalPages sets the total number of pages for the `pagination` instance.
// Ensure that totalPages is >= 0, defaults to 0 if invalid value.
//
// This function updates the `totalPages` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the total number of pages to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithTotalPages(v int) *pagination {
	if v < 0 {
		v = 0
	}
	p.totalPages = v
	return p
}

// WithTotalItems sets the total number of items for the `pagination` instance.
// Ensure that totalItems is >= 0, defaults to 0 if invalid value.
//
// This function updates the `totalItems` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the total number of items to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithTotalItems(v int) *pagination {
	if v < 0 {
		v = 0
	}
	p.totalItems = v
	p.calculate()
	return p
}

// WithIsLast sets whether this is the last page in the `pagination` instance.
//
// This function updates the `isLast` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: A boolean value indicating whether this is the last page.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithIsLast(v bool) *pagination {
	p.isLast = v
	return p
}

// Normalize adjusts the pagination fields to ensure consistency.
//
// This method recalculates the `totalPages` based on `totalItems` and `perPage`,
// ensuring that the pagination state is coherent. It also adjusts the `page` field
// to ensure it does not exceed `totalPages`, and sets the `isLast` field appropriately
// based on the current page position.
func (p *pagination) Normalize() *pagination {
	// Calculate total pages only if perPage is valid to avoid division by zero.
	if p.perPage > 0 {
		// Calculate total pages (ceiling division)
		p.totalPages = (p.totalItems + p.perPage - 1) / p.perPage
	} else {
		// If perPage is zero or negative, we cannot calculate totalPages
		p.totalPages = 0
	}

	// If there are no items, set totalPages to 0 and page to 1.
	if p.totalItems == 0 {
		p.totalPages = 0
		p.page = 1
		p.isLast = true
		return p
	}

	// Adjust the current page if it exceeds totalPages.
	if p.page > p.totalPages && p.totalPages > 0 {
		p.page = p.totalPages
	}

	// Determine if we are on the last page.
	p.isLast = p.page >= p.totalPages
	return p
}

// Respond generates a map representation of the `pagination` instance.
//
// This method collects various fields related to pagination (e.g., `page`, `per_page`, etc.)
// and organizes them into a key-value map. It ensures that only valid pagination details
// are included in the response.
//
// The following fields are included in the pagination response:
//   - `page`: The current page number.
//   - `per_page`: The number of items per page.
//   - `total_pages`: The total number of pages available.
//   - `total_items`: The total number of items available across all pages.
//   - `is_last`: A boolean indicating if this is the last page.
//
// Returns:
//   - A `map[string]interface{}` containing the structured pagination data.
func (p *pagination) Respond() map[string]any {
	m := make(map[string]any)
	if !p.Available() {
		return m
	}
	m["page"] = p.page
	m["per_page"] = p.perPage
	m["total_pages"] = p.totalPages
	m["total_items"] = p.totalItems
	m["is_last"] = p.isLast
	return m
}

// Json serializes the `pagination` instance into a compact JSON string.
//
// This function uses the `encoding.Json` utility to generate a JSON representation
// of the `pagination` instance. The output is a compact JSON string with no additional
// whitespace or formatting, providing a minimalistic view of the pagination data.
//
// Returns:
//   - A compact JSON string representation of the `pagination` instance.
func (p *pagination) Json() string {
	return jsonpass(p.Respond())
}

// JsonPretty serializes the `pagination` instance into a prettified JSON string.
//
// This function uses the `encoding.JsonPretty` utility to generate a JSON representation
// of the `pagination` instance. The output is a human-readable JSON string with
// proper indentation and formatting for better readability, which is helpful for
// inspecting pagination data during development or debugging.
//
// Returns:
//   - A prettified JSON string representation of the `pagination` instance.
func (p *pagination) JsonPretty() string {
	return jsonpretty(p.Respond())
}

// calculate computes the total pages and determines if the current page is the last one.
//
// This method performs calculations based on the `totalItems` and `perPage` fields
// to derive the `totalPages`. It uses ceiling division to ensure that any remaining
// items that don't fill a complete page are still counted as an additional page.
// Additionally, it checks if the current `page` is the last page by comparing it
// to `totalPages`, setting the `isLast` field accordingly.
func (p *pagination) calculate() {
	// Ensure page is at least 1.
	if p.page <= 0 {
		p.page = 1
	}

	// Only calculate if perPage is valid to avoid division by zero.
	if p.totalItems > 0 && p.perPage > 0 {
		// Calculate total pages (ceiling division)
		p.totalPages = (p.totalItems + p.perPage - 1) / p.perPage
	}

	// If there are no items, there is 1 page (page 1) which is empty,
	// or 0 pages. Usually, for API consistency, we treat 0 items as being on the last page.
	if p.totalItems == 0 {
		p.totalPages = 0
		p.isLast = true
		return
	}

	// Determine if we are on the last page.
	// We are on the last page if the current page is greater than or equal to totalPages.
	if p.totalPages > 0 {
		p.isLast = p.page >= p.totalPages
	}
}
