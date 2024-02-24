package control

import (
	"testing"
)

func TestEntry_Advance(t *testing.T) {
	tests := []struct {
		name     string
		entry    Entry
		expected Entry
	}{
		{
			name: "empty entry",
			entry: Entry{
				Key: []*Key{},
			},
			expected: Entry{
				Key: []*Key{},
			},
		},
		{
			name: "advance index",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
				},
			},
			expected: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
				},
			},
		},
		{
			name: "advance key",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
					{
						Key:   "test2",
						Index: NewObjects(NewObject(MakePtr("test2"))),
					},
				},
				KeyI: 0,
			},
			expected: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
					{
						Key:   "test2",
						Index: NewObjects(NewObject(MakePtr("test2"))),
					},
				},
				KeyI: 1,
			},
		},
		{
			name: "advance key",
			entry: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test")), NewObject(MakePtr("test"))),
					},
					{
						Key:   "test2",
						Index: NewObjects(NewObject(MakePtr("test2"))),
					},
				},
				KeyI: 0,
			},
			expected: Entry{
				Key: []*Key{
					{
						Key:   "test",
						Index: NewObjects(NewObject(MakePtr("test"))),
					},
					{
						Key:   "test2",
						Index: NewObjects(NewObject(MakePtr("test2"))),
					},
				},
				KeyI: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.entry.Advance()
			if !tt.entry.Equals(&tt.expected) {
				t.Errorf("Advance() = %v, want %v", tt.entry, tt.expected)
			}
		})
	}
}
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
