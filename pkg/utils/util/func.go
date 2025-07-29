package util

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

func Ptr[T any](v T) *T {
	return &v
}
