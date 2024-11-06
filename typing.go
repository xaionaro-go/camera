package camera

func typeAssertOrZero[T any](iface any) T {
	v, _ := iface.(T)
	return v
}
