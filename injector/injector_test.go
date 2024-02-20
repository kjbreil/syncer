package injector

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"net"
	"testing"
)

type TestStruct struct {
	String         string
	Int            int
	Slice          []int
	SliceStruct    []SD
	SlicePtr       []*int
	SlicePtrStruct []*SD
	SliceSlice     [][]int
	Map            map[string]int
	MapStruct      map[string]TestStruct
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*TestStruct
	MapMap         map[string]map[string]int
	Sub            TestSub
	SubPtr         *TestStruct
	IP             net.IP
}
type TestSub struct {
	Sub2 string
}

type TestSub2 struct {
	String string
}

type SD struct {
	Name string
	Data string
}
type TestStruct2 struct {
	Reservations map[string]Reservation4
}
type Reservation4 struct {
	IP       net.IP           `json:"IP,omitempty"`
	MAC      net.HardwareAddr `json:"MAC,omitempty"`
	Hostname string           `json:"Hostname,omitempty"`
}

func TestInjector_AddS(t *testing.T) {
	ts := TestStruct2{

		// Reservations: map[string]Reservation4{
		// 	"sn": {
		// 		IP: net.ParseIP("192.168.1.1"),
		// 	},
		// },
	}
	inj, err := New(&ts)
	if err != nil {
		t.Fatal(err)
	}

	// entries := []*control.Entry{
	// 	{
	// 		Key: []*control.Key{
	// 			{
	// 				Key: "TestStruct2",
	// 			},
	// 			{
	// 				Key:   "SliceStruct",
	// 				Index: control.NewObjects(0),
	// 			},
	// 			{
	// 				Key: "Name",
	// 			},
	// 		},
	// 		Value: control.NewObject("oneone"),
	// 	},
	// }

	entries := []*control.Entry{
		{
			Key: []*control.Key{
				{
					Key: "TestStruct2",
				},
				{
					Key:   "Reservations",
					Index: control.NewObjects("sn"),
				},

				{
					Key:   "IP",
					Index: control.NewObjects(0),
				},
			},
			Value: control.NewObject(uint8(255)),
		},
	}
	for _, e := range entries {

		err = inj.Add(e)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println(ts)
}

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
			name: "change string",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Sub",
						},
						{
							Key: "Sub2",
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
			name: "Test Sub",
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
			name: "Add To Slice",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "IP",
							Index: control.NewObjects(0),
						},
					},
					Value: control.NewObject(255),
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
			name: "Add To SliceSlice",
			entries: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceSlice",
							Index: control.NewObjects(0, control.NewObject(1)),
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
					if v, ok := v["bottom"]; !ok || v != 1 {
						return fmt.Errorf("ts.Map[\"test\"] is %d, should be 1", v)
					}
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
