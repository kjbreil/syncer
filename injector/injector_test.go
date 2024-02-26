package injector

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/helpers/equal"
	. "github.com/kjbreil/syncer/helpers/test"
	"testing"
)

//nolint:gocognit
func TestInjector_Add(t *testing.T) {
	ts := TestStruct{}
	inj, err := New(&ts)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		entries control.Entries
		preFn   func()
		wantErr bool
		wantFn  func() error
	}{
		{
			name: "Test Sub",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStruct",
						},
						{
							Key: "S",
						},
					},
					Value: &control.Object{String_: control.MakePtr("change string")},
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.SubStruct.S != "change string" {
					return fmt.Errorf("string %s should be \"change string\"", ts.String)
				}
				return nil
			},
		},
		{
			name: "change string",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "String",
						},
					},
					Value: &control.Object{String_: control.MakePtr("change string")},
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.String != "change string" {
					return fmt.Errorf("string %s should be \"change string\"", ts.String)
				}
				return nil
			},
		},
		{
			name: "Add To Slice",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(0),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if len(ts.Slice) == 0 || ts.Slice[0] != 1 {
					return fmt.Errorf("ts.Slice is length %d, should be 1", len(ts.Slice))
				}
				return nil
			},
		},
		{
			name: "Nil Slice",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Slice",
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Slice != nil {
					return fmt.Errorf("ts.Slice should be nil")
				}
				return nil
			},
		},
		{
			name: "Add To SliceSlice",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceSlice",
							Index: control.NewObjects(control.NewObject(0), control.NewObject(1)),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.SliceSlice[0][1] != 1 {
					return fmt.Errorf("ts.SliceSlice value not changed")
				}
				return nil
			},
		},
		{
			name: "Add To SliceStruct",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Data 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Data 2")),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if len(ts.SliceStruct) != 2 {
					return fmt.Errorf("ts.SliceStruct not two objets")
				}
				return nil
			},
		},
		{
			name: "Add To Map",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects("test"),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.Map["test"]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "Add To MapStruct",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects("test"),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.Map["test"]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "Nil Map",
			preFn: func() {
				ts.Map = map[string]int{"test": 1}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Map",
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Map != nil {
					return fmt.Errorf("ts.Map should be nil")
				}
				return nil
			},
		},
		{
			name: "Remove Map Key",
			preFn: func() {
				ts.Map = map[string]int{"Base First": 1}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Map != nil {
					if _, ok := ts.Map["Base First"]; ok {
						return fmt.Errorf("ts.Map[\"Base First\"] should be nil")
					}
				}
				return nil
			},
		},
		{
			name: "Remove Not there Map",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if _, ok := ts.Map["Base First"]; ok {
					return fmt.Errorf("ts.Map[\"Base First\"] should not be there")
				}
				return nil
			},
		},
		{
			name: "Add to Empty Map",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.Map["Base First"]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"Base First\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "Add to Empty MapStruct",
			preFn: func() {
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base First Test Struct")),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.MapStruct != nil {
					if _, ok := ts.MapStruct["Base First"]; !ok {
						return fmt.Errorf("ts.MapStruct[\"Base First\"] should exist")
					}
				}
				return nil
			},
		},
		{
			name: "Add to Empty MapMapType",
			preFn: func() {
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMapType",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First")), control.NewObject(control.MakePtr("Base Second"))),
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("inside")),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.MapMapType != nil {
					if ms, ok := ts.MapMapType["Base First"]; !ok {
						return fmt.Errorf("ts.MapStruct[\"Base First\"] should exist")
					} else {
						if ss, ok := ms["Base Second"]; !ok {
							return fmt.Errorf("ts.MapStruct[\"Base Second\"] should exist")
						} else {
							if ss.S != "inside" {
								return fmt.Errorf("ts.MapStruct[\"Base Second\"].S should be inside")
							}
						}
					}
				}
				return nil
			},
		},
		{
			name: "MapKeyBool",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyBool",
							Index: control.NewObjects(true),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.MapKeyBool[true]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "MapKeyUint",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyUint",
							Index: control.NewObjects(uint(1)),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.MapKeyUint[1]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "MapKeyFloat",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyFloat",
							Index: control.NewObjects(float64(1)),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.MapKeyFloat[1]; !ok || v != 1 {
					return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
				}
				return nil
			},
		},
		{
			name: "Add To MapMap",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMap",
							Index: control.NewObjects("top", control.NewObject("bottom")),
						},
					},
					Value: control.NewObject(1),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.MapMap["top"]; !ok {
					return fmt.Errorf("MapMap[\"top\"] was not found")
				} else {
					if v, ok := v["bottom"]; !ok || v != 1 {
						return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
					}
				}
				return nil
			},
		},
		{
			name: "Add Two To MapMap",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMap",
							Index: control.NewObjects("top", control.NewObject("bottom1")),
						},
					},
					Value: control.NewObject(1),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMap",
							Index: control.NewObjects("top", control.NewObject("bottom2")),
						},
					},
					Value: control.NewObject(2),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if v, ok := ts.MapMap["top"]; !ok {
					return fmt.Errorf("MapMap[\"top\"] was not found")
				} else {
					if v, ok := v["bottom2"]; !ok || v != 2 {
						return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
					}
				}
				return nil
			},
		},
		{
			name: "Add Pointer Struct",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStructPtr",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("SubStructPtr String")),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.SubStructPtr != nil && ts.SubStructPtr.String == "SubStructPtr String" {
					return nil
				}
				return fmt.Errorf("ts.SubStructPtr.String should be SubStructPtr String")
			},
		},
		{
			name: "Remove Pointer Struct",
			preFn: func() {
				ts.SubStructPtr = &TestStruct{
					String: "SubStructPtr String",
				}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStructPtr",
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.SubStructPtr != nil {
					return fmt.Errorf("ts.SubStructPtr should not exist")
				}
				return nil
			},
		},
		{
			name: "Add Interface",
			preFn: func() {
				ts.Interface = &TestInterfaceImpl{}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Interface",
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Interface == nil {
					return fmt.Errorf("ts.Interface should exist")
				}
				if ts.Interface.String() != "TestInterface" {
					return fmt.Errorf("ts.Interface.String() should be TestInterface")
				}
				return nil
			},
		},
		{
			name: "Remove Interface",
			preFn: func() {
				ts.Interface = &TestInterfaceImpl{S: "TestInterface"}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Interface",
						},
					},
					Remove: true,
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Interface != nil {
					return fmt.Errorf("ts.Interface should not exist")
				}
				return nil
			},
		},
		{
			name: "Add to Array",
			preFn: func() {
				ts.Array = [10]int{1, 2, 3, 4}
			},
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(4)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(5))),
				},
			},
			wantErr: false,
			wantFn: func() error {
				if ts.Array[4] != 5 {
					return fmt.Errorf("ts.Array[4] should be 5, is %d", ts.Array[4])
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFn != nil {
				tt.preFn()
			}

			for _, e := range tt.entries {
				err = inj.Add(e)
				if (err != nil) != tt.wantErr {
					t.Errorf("Injector.Add() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.wantFn != nil {
				if err = tt.wantFn(); err != nil {
					t.Errorf("Injector.Add() = %v", err)
				}
			}
		})
	}
}

func TestInjector_AddAll(t *testing.T) {
	ts := TestStruct{
		Interface: &TestInterfaceImpl{S: ""},
	}
	inj, err := New(&ts)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	entries := []*control.Entry{
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base String")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key: "Int",
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key: "Interface",
				},
				{
					Key: "S",
				},
			},
			Value: control.NewObject(control.MakePtr("TestInterface")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Slice",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Slice",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Slice",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(2)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(3))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Name 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Data 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Name 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Data 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SlicePtr",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SlicePtr",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SlicePtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SlicePtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceInterface",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "S",
				},
			},
			Value: control.NewObject(control.MakePtr("TestInterface")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "SliceSlice",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Array",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Name 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Data 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Name 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SliceStruct Data 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtr",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtr",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Data 1")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Name",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayPtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
				{
					Key: "Data",
				},
			},
			Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Data 2")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayInterface",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
				},
				{
					Key: "S",
				},
			},
			Value: control.NewObject(control.MakePtr("TestInterface")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayArray",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayArray",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayArray",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1))), control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(3))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "ArrayArray",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1))), control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(4))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Map",
					Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "Map",
					Index: control.NewObjects(control.NewObject(control.MakePtr("Base Second"))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapKeyInt",
					Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapKeyBool",
					Index: control.NewObjects(control.NewObject(control.MakePtr(true))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapKeyUint",
					Index: control.NewObjects(control.NewObject(control.MakePtr(uint64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapKeyFloat",
					Index: control.NewObjects(control.NewObject(control.MakePtr(float64(1.000000)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base First Test Struct")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapPtr",
					Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapPtrStruct",
					Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base First Test Struct")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapInterface",
					Index: control.NewObjects(control.NewObject(control.MakePtr("one"))),
				},
				{
					Key: "S",
				},
			},
			Value: control.NewObject(control.MakePtr("TestInterface")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapMap",
					Index: control.NewObjects(control.NewObject(control.MakePtr("top")), control.NewObject(control.MakePtr("second"))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(3))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapSlice",
					Index: control.NewObjects(control.NewObject(control.MakePtr("one")), control.NewObject(control.MakePtr(int64(0)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(1))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key:   "MapSlice",
					Index: control.NewObjects(control.NewObject(control.MakePtr("one")), control.NewObject(control.MakePtr(int64(1)))),
				},
			},
			Value: control.NewObject(control.MakePtr(int64(2))),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key: "SubStruct",
				},
				{
					Key: "S",
				},
			},
			Value: control.NewObject(control.MakePtr("Base Test Sub")),
		},
		{
			Key: []*control.Key{
				{
					Key: "TestStruct",
				},
				{
					Key: "SubStructPtr",
				},
				{
					Key: "String",
				},
			},
			Value: control.NewObject(control.MakePtr("Base Sub Test Struct")),
		},
	}

	err = inj.AddAll(entries)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	bs := MakeBaseTestStruct()
	if !equal.Any(ts, bs) {
		ext, _ := extractor.New(ts)
		_, _ = ext.Diff(&ts)
		diffEntries := ext.Entries(&bs)
		fmt.Println(diffEntries)
		t.Fatalf("ts does not equal MakeBaseTestStruct()")
	}

}
