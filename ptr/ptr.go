package ptr

// Pointer return pointer on any value.
func Pointer[T any](v T) *T {
	return &v
}

// Value return value of pointer, or default
func Value[T any](v *T) T {
	if v == nil {
		v = new(T)
	}
	return *v
}
