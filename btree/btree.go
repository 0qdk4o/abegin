// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package btree

// TreeNode represent binary search tree
type TreeNode struct {
	data  int
	left  *TreeNode
	right *TreeNode
}

// NewTree construct new tree
func NewTree(n int) *TreeNode {
	return &TreeNode{n, nil, nil}
}

// AddTreeNode add tree node to tree
func (t *TreeNode) AddTreeNode(n int) {
	if t == nil || t.data == n {
		return
	}

	if n < t.data {
		if t.left == nil {
			t.left = NewTree(n)
			return
		}
		t.left.AddTreeNode(n)
	} else {
		if t.right == nil {
			t.right = NewTree(n)
			return
		}
		t.right.AddTreeNode(n)
	}
}

// CompareTree tell whether tree t is equal to b, tree's structure may diff
func (t *TreeNode) CompareTree(b *TreeNode, order Method) bool {
	it := NewTreeIter(t, order)
	ib := NewTreeIter3(b, order)

	ift := it.Next()
	ifb := ib.Next()

	for ift && ifb {
		if it.Value() != ib.Value() {
			return false
		}
		ift = it.Next()
		ifb = ib.Next()
	}

	if ift || ifb {
		return false
	}

	return true
}

// IsSameTree return true if have the same tree structure and value
func (t *TreeNode) IsSameTree(b *TreeNode) bool {
	//TODO
	return true
}

// Search return true if there is a node that value equal to v in the tree
func (t *TreeNode) Search(v int) bool {
	it := NewTreeIter3(t, InOrder)
	if it == nil {
		return false
	}

	for it.Next() {
		if it.Value() == v {
			return true
		}
	}
	return false
}

// Depth return the depth of the tree t
func (t *TreeNode) Depth() int {
	ldepth := 0
	rdepth := 0
	if t.left != nil {
		ldepth = t.left.Depth()
	}

	if t.right != nil {
		rdepth = t.right.Depth()
	}

	if ldepth > rdepth {
		return ldepth + 1
	}

	return rdepth + 1
}
