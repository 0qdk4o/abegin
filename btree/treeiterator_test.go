// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package btree

import (
	"testing"
)

func TestTreeIter(t *testing.T) {
	root := ConstructTree(pairs)

	it := NewTreeIter(root, InOrder)
	if it == nil {
		t.Fatalf("can not construct iterator")
	}

	for i, v := range InOrderOut {
		if !it.Next() {
			t.Fatalf("[InOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
	}

	it = NewTreeIter(root, PreOrder)
	if it == nil {
		t.Fatalf("can not construct iterator")
	}

	for i, v := range PreOrderOut {
		if !it.Next() {
			t.Fatalf("[PreOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
	}

	it = NewTreeIter(root, PostOrder)
	if it == nil {
		t.Fatalf("can not construct iterator")
	}

	for i, v := range PostOrderOut {
		if !it.Next() {
			t.Fatalf("[PostOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
	}
}

func TestTreeIter2(t *testing.T) {
	root := ConstructTree(pairs)
	it := NewTreeIter2(root, InOrder)
	if it == nil {
		t.Fatalf("can not construct iterator2")
	}
	for i, v := range InOrderOut {
		if !it.Next() {
			t.Fatalf("[InOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("InOrder: %d\n", it.Value())
	}

	it = NewTreeIter2(root, PreOrder)
	if it == nil {
		t.Fatalf("can not construct iterator2")
	}
	for i, v := range PreOrderOut {
		if !it.Next() {
			t.Fatalf("[PreOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PreOrder: %d\n", it.Value())
	}

	it = NewTreeIter2(root, PostOrder)
	if it == nil {
		t.Fatalf("can not construct iterator2")
	}
	for i, v := range PostOrderOut {
		if !it.Next() {
			t.Fatalf("[PostOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PostOrder: %d\n", it.Value())
	}
}

func TestTreeIter3(t *testing.T) {
	root := ConstructTree(pairs)
	it := NewTreeIter3(root, PreOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range PreOrderOut {
		if !it.Next() {
			t.Fatalf("[PreOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PreOrder: %d\n", it.Value())
	}

	it = NewTreeIter3(root, InOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range InOrderOut {
		if !it.Next() {
			t.Fatalf("[InOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("InOrder: %d\n", it.Value())
	}

	it = NewTreeIter3(root, PostOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range PostOrderOut {
		if !it.Next() {
			t.Fatalf("[PostOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PostOrder: %d\n", it.Value())
	}
}

func TestTreeIter4(t *testing.T) {
	root := ConstructTree(pairs)
	it := NewTreeIter4(root, PreOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range PreOrderOut {
		if !it.Next() {
			t.Fatalf("[PreOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PreOrder: %d\n", it.Value())
	}

	it = NewTreeIter4(root, InOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range InOrderOut {
		if !it.Next() {
			t.Fatalf("[InOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("InOrder: %d\n", it.Value())
	}

	it = NewTreeIter4(root, PostOrder)
	if it == nil {
		t.Fatalf("can not construct iterator3")
	}
	for i, v := range PostOrderOut {
		if !it.Next() {
			t.Fatalf("[PostOrder]Item absent from index %d value %d\n", i, v)
		}
		if v != it.Value() {
			t.Fatalf("Iterator output wrong sequence, output=%d, want %d\n",
				it.Value(), v)
		}
		t.Logf("PostOrder: %d\n", it.Value())
	}
}

// Recursive traversal
func BenchmarkTreeIter(b *testing.B) {
	t := ConstructTree(pairs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := NewTreeIter(t, InOrder)
		for it.Next() {
			// _ = it.Value()
		}
	}
	b.StopTimer()
}

// Recursive traversal and channel
func BenchmarkTreeIter2(b *testing.B) {
	t := ConstructTree(pairs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := NewTreeIter2(t, InOrder)
		for it.Next() {
			// _ = it.Value()
		}
	}
	b.StopTimer()
}

// loop is best performance
func BenchmarkTreeIter3(b *testing.B) {
	t := ConstructTree(pairs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := NewTreeIter3(t, InOrder)
		for it.Next() {
			// _ = it.Value()
		}
	}
	b.StopTimer()
}

// worst performance
func BenchmarkTreeIter4(b *testing.B) {
	t := ConstructTree(pairs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := NewTreeIter4(t, InOrder)
		for it.Next() {
			// _ = it.Value()
		}
	}
	b.StopTimer()
}
