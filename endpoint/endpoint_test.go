package endpoint

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/kjbreil/syncer/endpoint/settings"
)

type testStruct struct {
	String         string
	Int            int
	Slice          []int
	SliceStruct    []sd
	SlicePtr       []*int
	SlicePtrStruct []*sd
	Map            map[string]int
	MapStruct      map[string]testStruct
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*testStruct
	Sub            testSub
	SubPtr         *testStruct
}

type testSub struct {
	String string
}

type sd struct {
	name string
	data string
}

func TestEndpoint_Run(t *testing.T) {
	peersOne := []net.TCPAddr{
		{
			IP:   net.ParseIP("10.0.2.2"),
			Port: 45014,
		},
	}

	peersTwo := []net.TCPAddr{
		{
			IP:   net.ParseIP("10.0.2.2"),
			Port: 45014,
		},
	}

	portOne := 45014
	portTwo := 45014

	dataOne := testStruct{
		String: "String1",
	}

	dataTwo := testStruct{
		String: "",
	}

	endpointOne, err := New(&dataOne, &settings.Settings{
		Port:       portOne,
		Peers:      peersOne,
		AutoUpdate: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	endpointTwo, err := New(&dataTwo, &settings.Settings{
		Port:       portTwo,
		Peers:      peersTwo,
		AutoUpdate: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	endpointOne.Run(false)
	for !endpointOne.Running() {
		time.Sleep(time.Second)
	}
	endpointTwo.Run(false)
	for !endpointTwo.Running() {
		time.Sleep(time.Second)
	}

	if dataOne.String != dataTwo.String {
		t.Fatal("dataOne.String != dataTwo.String")
	}

	endpointTwo.client.Init()

	dataOne.String = "String2"
	time.Sleep(time.Second)

	if dataOne.String != dataTwo.String {
		t.Fatal("dataOne.String != dataTwo.String")
	}

	endpointOne.Stop()
	time.Sleep(10 * time.Second)
	endpointOne, err = New(&dataOne, &settings.Settings{
		Port:       portOne,
		Peers:      peersOne,
		AutoUpdate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	endpointOne.Run(false)

	dataTwo.String = "String3"
	fmt.Println(dataOne.String)
	fmt.Println(dataTwo.String)
	endpointOne.client.Changes()

	fmt.Println(dataOne.String)
	fmt.Println(dataTwo.String)
	// endpointTwo.Client.Init()
	fmt.Println(dataOne.String)
	fmt.Println(dataTwo.String)
	endpointOne.Stop()
	endpointTwo.Stop()
	endpointOne.Wait()
	endpointTwo.Wait()
}

func Test_randomInt(t *testing.T) {
	type args struct {
		l int
		h int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				l: 100,
				h: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomInt(tt.args.l, tt.args.h); got > tt.args.h || got < tt.args.l {
				t.Errorf("randomInt() = %v, which is outside %d and %d", got, tt.args.l, tt.args.h)
			}
		})
	}
}
