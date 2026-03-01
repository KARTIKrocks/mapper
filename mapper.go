package mapper

import (
	"fmt"
	"reflect"
	"sync"
)

// Mapper is a thread-safe registry of type mapping functions.
type Mapper struct {
	mu    sync.RWMutex
	funcs map[typePair]func(any) any
}

// typePair uniquely identifies a source-to-destination type mapping.
type typePair struct {
	src, dst reflect.Type
}

// New creates a new empty Mapper.
func New() *Mapper {
	return &Mapper{funcs: make(map[typePair]func(any) any)}
}

// Reset removes all registered mappings.
func (m *Mapper) Reset() {
	m.mu.Lock()
	m.funcs = make(map[typePair]func(any) any)
	m.mu.Unlock()
}

// lookup retrieves a mapping function for the given type pair.
func (m *Mapper) lookup(key typePair) (func(any) any, bool) {
	m.mu.RLock()
	fn, ok := m.funcs[key]
	m.mu.RUnlock()
	return fn, ok
}

var global = New()

// Global returns the default Mapper instance used by package-level functions.
func Global() *Mapper { return global }

// --- Global API ---

// Register adds a mapping function from type S to type D.
// If a mapping for the same type pair already exists, it is replaced.
//
// The type pair is determined at compile time from the function signature.
// When using [Map], the source type is determined at runtime from the value
// passed, so it must match the registered source type exactly (e.g. if you
// register func(User) UserDTO, you must call Map with a User, not *User).
func Register[S, D any](fn func(S) D) {
	RegisterTo(global, fn)
}

// Map converts src to type D using a registered mapping function.
// It panics if no mapping is registered for the source type to D.
//
// The source type is determined at runtime via reflect.TypeOf(src),
// so it must match the type used in [Register] exactly.
func Map[D any](src any) D {
	return MapFrom[D](global, src)
}

// MapErr converts src to type D using a registered mapping function.
// It returns an error if no mapping is registered.
func MapErr[D any](src any) (D, error) {
	return MapErrFrom[D](global, src)
}

// MapSlice converts each element of src from type S to type D.
// It panics if no mapping is registered from S to D.
func MapSlice[S, D any](src []S) []D {
	return MapSliceFrom[S, D](global, src)
}

// MapSliceErr converts each element of src from type S to type D.
// It returns an error if no mapping is registered from S to D.
func MapSliceErr[S, D any](src []S) ([]D, error) {
	return MapSliceErrFrom[S, D](global, src)
}

// Has reports whether a mapping from type S to type D is registered.
func Has[S, D any]() bool {
	return HasIn[S, D](global)
}

// ResetGlobal removes all mappings from the global Mapper.
func ResetGlobal() { global.Reset() }

// --- Instance API ---

// RegisterTo adds a mapping function from type S to type D on the given Mapper.
func RegisterTo[S, D any](m *Mapper, fn func(S) D) {
	if fn == nil {
		panic("mapper: mapping function must not be nil")
	}
	key := typePair{
		src: reflect.TypeFor[S](),
		dst: reflect.TypeFor[D](),
	}
	wrapped := func(src any) any {
		return fn(src.(S))
	}
	m.mu.Lock()
	m.funcs[key] = wrapped
	m.mu.Unlock()
}

// MapFrom converts src to type D using the given Mapper.
// It panics if no mapping is registered for the source type to D.
//
// This is the optimized hot path — it avoids the error-return overhead
// of [MapErrFrom] and panics directly on lookup failure.
func MapFrom[D any](m *Mapper, src any) D {
	if src == nil {
		panic("mapper: source must not be nil")
	}
	key := typePair{
		src: reflect.TypeOf(src),
		dst: reflect.TypeFor[D](),
	}
	fn, ok := m.lookup(key)
	if !ok {
		panic(fmt.Sprintf("mapper: no mapping registered from %s to %s", key.src, key.dst))
	}
	return fn(src).(D)
}

// MapErrFrom converts src to type D using the given Mapper.
// It returns an error if no mapping is registered.
func MapErrFrom[D any](m *Mapper, src any) (D, error) {
	var zero D
	if src == nil {
		return zero, fmt.Errorf("mapper: source must not be nil")
	}
	key := typePair{
		src: reflect.TypeOf(src),
		dst: reflect.TypeFor[D](),
	}
	fn, ok := m.lookup(key)
	if !ok {
		return zero, fmt.Errorf("mapper: no mapping registered from %s to %s", key.src, key.dst)
	}
	return fn(src).(D), nil
}

// MapSliceFrom converts each element of src from S to D using the given Mapper.
// It panics if no mapping is registered from S to D.
func MapSliceFrom[S, D any](m *Mapper, src []S) []D {
	if src == nil {
		return nil
	}
	if len(src) == 0 {
		return []D{}
	}
	key := typePair{
		src: reflect.TypeFor[S](),
		dst: reflect.TypeFor[D](),
	}
	fn, ok := m.lookup(key)
	if !ok {
		panic(fmt.Sprintf("mapper: no mapping registered from %s to %s", key.src, key.dst))
	}
	result := make([]D, len(src))
	for i, s := range src {
		result[i] = fn(s).(D)
	}
	return result
}

// MapSliceErrFrom converts each element of src from S to D using the given Mapper.
// It returns an error if no mapping is registered from S to D.
func MapSliceErrFrom[S, D any](m *Mapper, src []S) ([]D, error) {
	if src == nil {
		return nil, nil
	}
	if len(src) == 0 {
		return []D{}, nil
	}
	key := typePair{
		src: reflect.TypeFor[S](),
		dst: reflect.TypeFor[D](),
	}
	fn, ok := m.lookup(key)
	if !ok {
		return nil, fmt.Errorf("mapper: no mapping registered from %s to %s", key.src, key.dst)
	}
	result := make([]D, len(src))
	for i, s := range src {
		result[i] = fn(s).(D)
	}
	return result, nil
}

// HasIn reports whether a mapping from S to D is registered in the given Mapper.
func HasIn[S, D any](m *Mapper) bool {
	key := typePair{
		src: reflect.TypeFor[S](),
		dst: reflect.TypeFor[D](),
	}
	_, ok := m.lookup(key)
	return ok
}
