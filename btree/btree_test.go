// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package btree

import (
	"testing"
)

type testpair struct {
	insN int

	//      root
	//     0/  \1
	//    a      b
	//   0/    0/ \1
	//  c      e   f
	//  1\    1\
	//    d     g
	//
	// data g's path is 101
	path []byte
}

var pairs = []testpair{
	{6, nil},
	{-1, []byte("0")},
	{11, []byte("1")},
	{4, []byte("01")},
	{15, []byte("11")},
	{13, []byte("110")},
	{4, nil}, // duplicate data 4
	{2, []byte("010")},
	{15, nil}, // duplicate data 15
}

var (
	InOrderOut   = []int{-1, 2, 4, 6, 11, 13, 15}
	PreOrderOut  = []int{6, -1, 4, 2, 11, 15, 13}
	PostOrderOut = []int{2, 4, -1, 13, 15, 11, 6}

	depthPairs    = 4
	depthcmpPairs = 5
)

var cmpPairs = []testpair{
	{11, nil},
	{-1, []byte("0")},
	{6, []byte("01")},
	{4, []byte("010")},
	{15, []byte("1")},
	{13, []byte("10")},
	{2, []byte("0111")},
}

func Travel(root *TreeNode, t *testing.T) {
	t.Helper() // Golang 1.9+
	it := NewTreeIter(root, InOrder)
	if it == nil {
		t.Fatalf("can not construct tree iterator")
	}
	for it.Next() {
		t.Logf("value: %d\n", it.Value())
	}
}

func ConstructTree(pairs []testpair) *TreeNode {
	var root *TreeNode
	for i, v := range pairs {
		if i == 0 {
			root = NewTree(v.insN)
		} else {
			root.AddTreeNode(v.insN)
		}
	}
	return root
}

func TestAddTreeNode(t *testing.T) {
	var root *TreeNode
	root.AddTreeNode(0)
	if root != nil {
		t.Fatal("nil recevier test case failed")
	}
	root = ConstructTree(pairs)

	Travel(root, t)

	var cur = root
	for i, v := range pairs {
		cur = root
		for _, j := range v.path {
			if cur == nil {
				t.Fatalf("Test data pairs is not correct at case %d\n", i)
			}
			if j == '0' {
				cur = cur.left
			} else {
				cur = cur.right
			}
		}
		if cur == nil {
			t.Fatalf("Test data pairs is not correct at case %d\n", i)
		}

		if cur.data != v.insN {
			if i == 0 { // duplicate data not to check, but root need check
				t.Fatalf("case %d TreeNode data = %d, want %d\n", i, cur.data, v.insN)
			}
		} else {
			t.Logf("case %d (depth: %d, value:%d) pass\n", i, len(v.path), cur.data)
		}
	}

	// after adding duplicate data, verify tree data
	var validTree *TreeNode
	for i, v := range pairs {
		if i == 0 {
			validTree = NewTree(v.insN)
		} else if v.path != nil { // only add uniq data to tree
			validTree.AddTreeNode(v.insN)
		}
	}
	Travel(validTree, t)
	if !validTree.CompareTree(root, InOrder) {
		t.Logf("The tree has duplicate node data\n")
	}
}

func TestCompareTree(t *testing.T) {
	ta := ConstructTree(pairs)
	tb := ConstructTree(cmpPairs)

	if !ta.CompareTree(tb, InOrder) {
		t.Fatal("tree compare error, should equal")
	}

	ta.AddTreeNode(18)
	if ta.CompareTree(tb, InOrder) {
		t.Fatal("tree compare error, should not equal")
	}

	tb.AddTreeNode(18)
	if !ta.CompareTree(tb, InOrder) {
		t.Fatal("tree compare error, should equal")
	}
}

func TestTreeDepth(t *testing.T) {
	ta := ConstructTree(pairs)
	tb := ConstructTree(cmpPairs)

	da := ta.Depth()
	db := tb.Depth()
	if da != depthPairs {
		t.Fatalf("Depth = %d, want %d\n", da, depthPairs)
	}
	if db != depthcmpPairs {
		t.Fatalf("Depth = %d, want %d\n", db, depthcmpPairs)
	}
}

func BenchmarkTreeDepth(b *testing.B) {
	t := ConstructTree(cmpPairs)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Depth()
	}
	b.StopTimer()
}
