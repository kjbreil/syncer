package control

func MakePtr[V any](v V) *V {
	return &v
}
