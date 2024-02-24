package control

import (
	"testing"
)

func TestEntry_IsLastIndex(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
		want  bool
	}{
		{
			name:  "empty entry",
			entry: Entry{},
			want:  true,
		},
		{
			name: "last index",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
				},
			},
			want: true,
		},
		{
			name: "not last index",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test")), NewObject(MakePtr("test"))),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.entry.IsLastIndex(); got != tt.want {
				t.Errorf("IsLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_IsLastKeyIndex(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
		want  bool
	}{
		{
			name:  "empty entry",
			entry: Entry{},
			want:  true,
		},
		{
			name: "last key",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
				},
			},
			want: true,
		},
		{
			name: "not last key",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test")), NewObject(MakePtr("test"))),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.entry.IsLastKeyIndex(); got != tt.want {
				t.Errorf("IsLastKeyIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
