package control

func (e *Entry) Advance() *Entry {
	e.Key = e.Key[1:]
	return e
}
