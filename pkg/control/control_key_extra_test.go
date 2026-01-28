package control

import "testing"

func TestKey_HasNoIndex(t *testing.T) {
    tests := []struct {
        name string
        key  *Key
        want bool
    }{
        {name: "nil index", key: &Key{}, want: true},
        {name: "empty slice", key: &Key{Index: []*Object{}}, want: false},
        {name: "non-empty slice", key: &Key{Index: NewObjects(NewObject(1))}, want: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.key.HasNoIndex(); got != tt.want {
                t.Errorf("HasNoIndex() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestKey_GetCurrentIndex(t *testing.T) {
    idxObj1 := NewObject(1)
    idxObj2 := NewObject(2)

    tests := []struct {
        name string
        key  *Key
        want *Object
    }{
        {name: "nil index", key: &Key{}, want: nil},
        {name: "first index", key: &Key{Index: NewObjects(idxObj1, idxObj2), IndexI: 0}, want: idxObj1},
        {name: "second index", key: &Key{Index: NewObjects(idxObj1, idxObj2), IndexI: 1}, want: idxObj2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.key.GetCurrentIndex()
            if (got == nil) != (tt.want == nil) {
                t.Fatalf("nil mismatch: got %v want %v", got, tt.want)
            }
            if got != nil && !got.Equals(tt.want) {
                t.Errorf("GetCurrentIndex() = %v, want %v", got, tt.want)
            }
        })
    }
}
