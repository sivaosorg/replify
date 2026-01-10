package coll

// HashMap is a generic hash map data structure that maps keys of type `K` to values of type `V`.
// `K` and `V` must be comparable types, meaning that they support comparison operators (like == and !=).
// The `items` field stores the actual map, which is used to store key-value pairs.
type HashMap[K, V comparable] struct {
	items map[K]V
}

// NewHashMap is a constructor function that initializes and returns a pointer to a new, empty `HashMap`.
// It creates a new map with keys of type `K` and values of type `V`.
//
// Generics are used to make this function flexible for any types of keys and values as long as they
// are comparable. The function ensures that the map is properly initialized before returning it.
//
// Example usage:
//
//	hashMap := NewHashMap[string, int]() // Creates a HashMap with string keys and int values.
func NewHashMap[K, V comparable]() *HashMap[K, V] {
	hash := &HashMap[K, V]{
		items: make(map[K]V),
	}
	return hash
}

// Put adds a new key-value pair to the HashMap. If the key already exists, its value is updated.
// Parameters:
//   - `key`: The key to be added or updated in the map.
//   - `value`: The value to be associated with the key.
//
// Example:
//
//	hashMap.Put("apple", 5) // Inserts or updates the value of "apple" to 5.
func (hash *HashMap[K, V]) Put(key K, value V) {
	hash.items[key] = value
}

// Get retrieves the value associated with the given key from the HashMap.
// If the key exists, it returns the corresponding value. If the key does not exist, it returns the zero value for the type `V`.
// Parameters:
//   - `key`: The key whose associated value is to be returned.
//
// Returns:
//   - The value associated with the key, or the zero value for `V` if the key is not found.
//
// Example:
//
//	value := hashMap.Get("apple") // Retrieves the value associated with "apple".
func (hash *HashMap[K, V]) Get(key K) (value V) {
	return hash.items[key]
}

// Remove deletes the key-value pair from the HashMap for the specified key.
// If the key does not exist, no action is taken.
// Parameters:
//   - `key`: The key to be removed from the map.
//
// Example:
//
//	hashMap.Remove("apple") // Removes the "apple" key and its associated value.
func (hash *HashMap[K, V]) Remove(key K) {
	delete(hash.items, key)
}

// Clear removes all key-value pairs from the HashMap, effectively resetting it to an empty map.
//
// Example:
//
//	hashMap.Clear() // Clears the map, removing all entries.
func (hash *HashMap[K, V]) Clear() {
	hash.items = make(map[K]V)
}

// Size returns the number of key-value pairs currently stored in the HashMap.
// Returns:
//   - The number of key-value pairs in the map.
//
// Example:
//
//	size := hashMap.Size() // Gets the size of the map.
func (hash *HashMap[K, V]) Size() int {
	return len(hash.items)
}

// IsEmpty checks if the HashMap is empty (i.e., contains no key-value pairs).
// Returns:
//   - `true` if the map is empty, `false` otherwise.
//
// Example:
//
//	isEmpty := hashMap.IsEmpty() // Returns true if the map is empty.
func (hash *HashMap[K, V]) IsEmpty() bool {
	return len(hash.items) == 0
}

// ContainsKey checks whether the specified key exists in the HashMap.
// Parameters:
//   - `key`: The key to check for existence in the map.
//
// Returns:
//   - `true` if the key exists in the map, `false` otherwise.
//
// Example:
//
//	exists := hashMap.ContainsKey("apple") // Checks if the key "apple" exists in the map.
func (hash *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := hash.items[key]
	return ok
}

// KeySet returns a slice of all keys currently stored in the HashMap.
// Returns:
//   - A slice containing all keys in the map.
//
// Example:
//
//	keys := hashMap.KeySet() // Returns all the keys in the map.
func (hash *HashMap[K, V]) KeySet() []K {
	keys := make([]K, hash.Size())
	i := 0
	for key := range hash.items {
		keys[i] = key
	}
	return keys
}
