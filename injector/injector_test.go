package injector

// func TestInjector_Add(t *testing.T) {
//
// 	type testStruct struct {
// 		String string
// 		Int    int
// 		Slice  []string
// 		Map    map[int]string
// 		Sub    *testStruct
// 	}
//
// 	ts := testStruct{
// 		String: "Test",
// 		Int:    1,
// 		Slice:  []string{"S0"},
// 		Map: map[int]string{
// 			1: "M1",
// 		},
// 		Sub: &testStruct{
// 			String: "SubTest",
// 			Slice:  nil,
// 			Map:    nil,
// 		},
// 	}
//
// 	// Create an injector that is invalid since the passed struct is not a pointer
// 	inj, err := New(ts)
// 	if err == nil {
// 		t.Fatal("passed non pointer to injector")
// 	}
//
// 	inj, err = New(&ts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = inj.Add(parse("testStruct.String", "Test2", Add))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ts.String != "Test2" {
// 		t.Fatal("failed to add string value at base level")
// 	}
//
// 	err = inj.Add(parse("testStruct.Slice[0]", "S0N", Add))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ts.Slice[0] != "S0N" {
// 		t.Fatal("failed to add slice value at base level")
// 	}
//
// 	err = inj.Add(parse("testStruct.Slice[5]", "S5", Add))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ts.Slice[5] != "S5" {
// 		t.Fatal("failed to add slice value at base level")
// 	}
//
// 	err = inj.Add(parse("testStruct.Map[1]", "M1N", Add))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ts.Map[1] != "M1N" {
// 		t.Fatal("failed to add map value at base level")
// 	}
//
// 	err = inj.Add(parse("testStruct.Map[2]", "M2", Add))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if ts.Map[2] != "M2" {
// 		t.Fatal("failed to add map value at base level")
// 	}
//
// }
