package extractor

import (
	"github.com/kjbreil/syncer/control"
	"testing"
)

func TestExtractor_GetDiff(t *testing.T) {

	tests := []struct {
		name string
		// structure must be a &struct otherwise reflect only sees an interfaces when dereferencing the pointer
		structure any
		want      []*control.Entry
	}{
		{
			name: "string",
			structure: &struct {
				String string
			}{
				String: "test",
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "Stringss",
						},
					},
					Value: control.NewObject(control.MakePtr("test")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, err := New(tt.structure)
			if err != nil {
				t.Fatalf("could not create extractor: %v", err)
			}
			got, err := ext.GetDiff(tt.structure)
			if err != nil {
				t.Fatalf("could not get diff: %v", err)
			}
			if !got.Equals(tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
				// t.Logf("##########\n\n%s\n\n##########", got.Diff(tt.want).Struct())
				t.Logf("##########\n\n%s\n\n##########", got.Struct())
			}
		})
	}
}
