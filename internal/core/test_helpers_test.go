package core_test

func ptr[T any](v T) *T {
	return &v
}
