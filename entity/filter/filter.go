package filter

import (
	"sync"
)

// Registry manages the mapping between permission strings and integer IDs.
// It is thread-safe for reading and writing.
type registry struct {
	mu     sync.RWMutex
	permID map[string]int
	nextID int
}

var Registry = &registry{
	permID: make(map[string]int),
	nextID: 0,
}

// Register registers a permission string and returns its unique ID.
// If the permission is already registered, it returns the existing ID.
func (r *registry) Register(perm string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id, ok := r.permID[perm]; ok {
		return id
	}

	id := r.nextID
	r.permID[perm] = id
	r.nextID++
	return id
}

// GetID returns the ID for a given permission string.
func (r *registry) GetID(perm string) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.permID[perm]
	return id, ok
}

// BitSet is a simple implementation of a bitset to store permissions.
type BitSet struct {
	set []uint64
}

func NewBitSet() *BitSet {
	return &BitSet{
		set: make([]uint64, 1),
	}
}

// Set sets the bit at the given index.
func (b *BitSet) Set(idx int) {
	wordIdx := idx / 64
	bitIdx := idx % 64

	for len(b.set) <= wordIdx {
		b.set = append(b.set, 0)
	}

	b.set[wordIdx] |= 1 << bitIdx
}

// Has checks if the bit at the given index is set.
func (b *BitSet) Has(idx int) bool {
	wordIdx := idx / 64
	bitIdx := idx % 64

	if wordIdx >= len(b.set) {
		return false
	}

	return (b.set[wordIdx] & (1 << bitIdx)) != 0
}
