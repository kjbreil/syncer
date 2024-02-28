package control

import (
	"testing"
)

func TestKey_IsLastIndex(t *testing.T) {
	tests := []struct {
		name string
		key  *Key
		want bool
	}{
		{
			name: "empty key",
			key:  &Key{},
			want: true,
		},
		{
			name: "last index",
			key: &Key{
				Index:  NewObjects(NewObject(MakePtr("test"))),
				IndexI: 0,
			},
			want: true,
		},
		{
			name: "not last index",
			key: &Key{
				Index:  NewObjects(NewObject(MakePtr("test")), NewObject(MakePtr("test"))),
				IndexI: 0,
			},
			want: false,
		},
		{
			name: "last index multi index",
			key: &Key{
				Index:  NewObjects(NewObject(MakePtr("test")), NewObject(MakePtr("test"))),
				IndexI: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.IsLastIndex(); got != tt.want {
				t.Errorf("IsLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
