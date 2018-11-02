package easyrouter

import (
	"testing"
)

func TestInsertFind(t *testing.T) {
	// create and populate tree
	values := []string{
		"cat",
		"cats",
		"hey",
		"hello",
		"dog",
		"foo/*/bar",
	}

	tree := NewTree()
	for i, value := range values {
		tree.Insert(value, i)
	}

	if tree == nil {
		t.Error("expected tree to not be nil")
	}

	// test Find
	tests := []struct {
		value     string
		find      bool
		leafValue int
	}{
		{"cat", true, 0},
		{"cats", true, 1},
		{"car", false, -1},
		{"ca", false, -1},
		{"hey", true, 2},
		{"heys", false, -1},
		{"he", false, -1},
		{"hello dog", false, -1},
		{"foo/123/bar", true, 5},
		{"foo/123/bars", false, -1},
	}

	for i, test := range tests {
		found, value := tree.Find(test.value)
		if found != test.find {
			t.Errorf("expected value %s to be %t, but got %t", test.value, test.find, found)
		}
		if found {
			if value != test.leafValue {
				t.Errorf("expected value %d to be %d on test %d", value, test.leafValue, i)
			}
		}
	}
}

//car
//cats
//  c
//  |
//  a
//  /\
// r  t
//     |
//      s (leaf)
