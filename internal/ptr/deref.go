package ptr

func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func PtrOrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func SlicePtrOrNil[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}
	return &s
}
