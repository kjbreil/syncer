package control

func (k *Key) IsLastIndex() bool {
	return k.GetIndex() == nil || len(k.GetIndex())-1 == int(k.IndexI)
}

func (k *Key) HasNoIndex() bool {
	return k.GetIndex() == nil
}

func (k *Key) GetCurrentIndex() *Object {
	if k.GetIndex() == nil {
		return nil
	}
	return k.GetIndex()[int(k.IndexI)]
}
