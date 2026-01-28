package control

import "testing"

// TestNewEntryAndRemoveEntry ensure constructors setup key slice length and flags.
func TestNewEntryAndRemoveEntry(t *testing.T) {
    e := NewEntry(3, "value")
    if len(e.Key) != 0 {
        t.Errorf("expected empty key slice, got %d", len(e.Key))
    }
    if e.Remove {
        t.Errorf("expected Remove false on NewEntry")
    }
    if e.Value == nil || e.Value.GetString_() != "value" {
        t.Errorf("unexpected value field")
    }

    r := NewRemoveEntry(2)
    if !r.Remove {
        t.Errorf("expected Remove true on NewRemoveEntry")
    }
    if r.Value != nil {
        t.Errorf("expected nil Value on remove entry")
    }
}

// helper to build Entries quickly.
func makeEntry(path []string, value any) *Entry {
    e := NewEntry(len(path), value)
    for _, p := range path {
        e.Key = append(e.Key, &Key{Key: p})
    }
    return e
}

// TestAddKeyAndIndex verify Entries.AddKey and AddIndex operations.
func TestAddKeyAndIndex(t *testing.T) {
    // Prepare two entries with blank top-level key
    ent := Entries{
        &Entry{Key: []*Key{{}}},
        &Entry{Key: []*Key{{}}},
    }
    ent.AddKey("root")
    for _, e := range ent {
        if e.Key[0].Key != "root" {
            t.Fatalf("AddKey failed; want root, got %s", e.Key[0].Key)
        }
    }

    ent2 := Entries{
        &Entry{Key: []*Key{{}}},
    }
    ent2.AddIndex(5)
    if len(ent2[0].Key[0].Index) == 0 || ent2[0].Key[0].Index[0].GetInt64() != 5 {
        t.Fatalf("AddIndex failed; expected index 5")
    }
}

// TestEntriesEqualsAndDiff checks Equals and Diff behaviour.
func TestEntriesEqualsAndDiff(t *testing.T) {
    a1 := makeEntry([]string{"root", "a"}, 1)
    a2 := makeEntry([]string{"root", "b"}, 2)
    set1 := Entries{a1, a2}
    set2 := Entries{a1, a2}

    if !set1.Equals(set2) {
        t.Fatalf("identical entries should be equal")
    }

    // modify
    b1 := makeEntry([]string{"root", "a"}, 3)
    set3 := Entries{b1, a2}
    diff := set1.Diff(set3)
    if len(diff) != 1 {
        t.Fatalf("expected single diff entry, got %d", len(diff))
    }
    if !diff[0].Equals(a1) {
        t.Fatalf("diff did not return expected entry")
    }
}
