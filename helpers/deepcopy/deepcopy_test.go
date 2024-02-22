package deepcopy

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
	"testing"
)

type IFace interface {
	String() string
}

type IFaceImpl struct {
	S string
}

func (t *IFaceImpl) String() string {
	return t.S
}

func Test_copyValue(t *testing.T) {
	type args struct {
		dst any
		src any
	}
	tests := []struct {
		name   string
		dst    any
		src    any
		wantFn func(src, dst any) (bool, string)
	}{
		{
			name: "int",
			dst:  3,
			src:  1,
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), fmt.Sprintf("src: %v, dst: %v", src, dst)
			},
		},

		{
			name: "int pointer",
			dst:  control.MakePtr(3),
			src:  control.MakePtr(1),
			wantFn: func(src, dst any) (bool, string) {
				if src == dst {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "map ptr val",
			dst:  nil,
			src: map[string]*int{
				"test": control.MakePtr(1),
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(map[string]*int)
				d := dst.(map[string]*int)

				if s["test"] == d["test"] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "slice ptr val",
			dst:  nil,
			src: []*int{
				control.MakePtr(1),
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.([]*int)
				d := dst.([]*int)

				if s[0] == d[0] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "interface",
			dst:  nil,
			src: map[string]IFace{
				"test": &IFaceImpl{S: "test"},
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(map[string]IFace)
				d := dst.(map[string]IFace)

				if s["test"] == d["test"] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "unexported type",
			dst:  nil,
			src: struct {
				String     string
				unexported string
			}{
				String:     "String",
				unexported: "unexported",
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(struct {
					String     string
					unexported string
				})
				d := dst.(struct {
					String     string
					unexported string
				})
				if s.String != d.String {
					return false, fmt.Sprintf("s String: %v != d String: %v", s.String, d.String)
				}

				if s.unexported == d.unexported {
					return false, fmt.Sprintf("s unexported: %v == d unexported: %v", s.unexported, d.unexported)
				}

				return true, ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcV := reflect.ValueOf(tt.src)
			tt.dst = DeepCopy(srcV).Interface()
			if ok, errS := tt.wantFn(tt.src, tt.dst); !ok {
				t.Errorf(errS)
			}
		})
	}
}
