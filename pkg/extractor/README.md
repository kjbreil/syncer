## Extractor
Extractor is initialized with a struct and then fed in the same struct again to find the changes. The initial load of
the struct brings in a blank version of the struct so the first run will return the current state of the struct. From
that point forward only changes will be returned.

```go

type data struct {
	String string
}

func main() {
	t := data {
		String: "test",
    }
	ext, err := extractor.New(&t)
	if err != nil {
		panic(err)
    }
    // gets full structure (in the case only String = "test"
    entries, err := ext.Entries(tt.structure)
    if err != nil {
        panic(err)
    }
    t.String = "new test"
    // gets changes
    entries, err = ext.Entries(tt.structure)
    if err != nil {
        panic(err)
    }
	
    ... 
}

```