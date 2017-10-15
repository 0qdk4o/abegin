// Copyright (c) 2017, 0qdk4o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package btree

// Method indicate tree traversal method, pre-order, in-order, post-order
type Method int

// show tree traversal method
const (
	PreOrder Method = iota
	InOrder
	PostOrder
)

// Iter is a tree iterator, Recursively traversal
type Iter struct {
	nodeChain []*TreeNode
	pleft     *TreeNode
	pright    *TreeNode
	order     Method
	value     int
}

// NewTreeIter construct a new Tree Iterator
func NewTreeIter(t *TreeNode, m Method) *Iter {
	if t == nil {
		return nil
	}
	it := &Iter{order: m}
	it.fillNodes(t)
	return it
}

func (it *Iter) fillNodes(t *TreeNode) {
	if t == nil {
		return
	}

	if it.order == PreOrder {
		it.nodeChain = append(it.nodeChain, t)
		it.fillNodes(t.left)
		it.fillNodes(t.right)
	} else if it.order == InOrder {
		it.fillNodes(t.left)
		it.nodeChain = append(it.nodeChain, t)
		it.fillNodes(t.right)
	} else if it.order == PostOrder {
		it.fillNodes(t.left)
		it.fillNodes(t.right)
		it.nodeChain = append(it.nodeChain, t)
	}
}

// Next return true if next tree node is exist
func (it *Iter) Next() bool {
	if s := len(it.nodeChain); s > 0 {
		it.value = it.nodeChain[0].data
		it.nodeChain = it.nodeChain[1:s]
		return true
	}
	return false
}

// Value return current node value
func (it *Iter) Value() int {
	return it.value
}

// Iter2 is a iterator implemented with channel, it's similar with yield in python
type Iter2 struct {
	ch    chan int
	order Method
	value int
}

// NewTreeIter2 construct new iterator with t
func NewTreeIter2(t *TreeNode, m Method) *Iter2 {
	if t == nil {
		return nil
	}
	ch := make(chan int)
	it2 := &Iter2{ch: ch, order: m}
	go func() {
		defer close(it2.ch)
		it2.travelTree(t)
	}()
	return it2
}

func (it2 *Iter2) travelTree(tn *TreeNode) {
	if tn == nil {
		return
	}
	if it2.order == PreOrder {
		it2.ch <- tn.data
		it2.travelTree(tn.left)
		it2.travelTree(tn.right)
	} else if it2.order == InOrder {
		it2.travelTree(tn.left)
		it2.ch <- tn.data
		it2.travelTree(tn.right)
	} else if it2.order == PostOrder {
		it2.travelTree(tn.left)
		it2.travelTree(tn.right)
		it2.ch <- tn.data
	}
}

// Next return true if there is still valid data exist
func (it2 *Iter2) Next() bool {
	for v := range it2.ch {
		it2.value = v
		return true
	}
	return false
}

// Value return current node data
func (it2 *Iter2) Value() int {
	return it2.value
}

// Iter3 is a tree iterator, use for loop other than recursively traversal
type Iter3 struct {
	pnode     *TreeNode
	prenode   *TreeNode
	nodeChain []*TreeNode
	order     Method
}

// NewTreeIter3 construct a new Tree Iterator
func NewTreeIter3(t *TreeNode, m Method) *Iter3 {
	it := &Iter3{pnode: t, order: m}
	return it
}

// Next return true if there is still valid data exist
func (it *Iter3) Next() bool {
	if it.order == PreOrder {
		return it.hasNextPreOrder()
	} else if it.order == InOrder {
		return it.hasNextInOrder()
	} else if it.order == PostOrder {
		return it.hasNextPostOrder()
	}

	return false
}

func (it *Iter3) hasNextPreOrder() bool {
	for {
		if it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.prenode = it.pnode
			it.pnode = it.pnode.left
			return true
		}
		l := len(it.nodeChain)
		if l > 0 {
			it.pnode = it.nodeChain[l-1].right
			it.nodeChain = it.nodeChain[:l-1] // pop
		}
		if it.pnode == nil && l == 0 {
			return false
		}
	}
}

func (it *Iter3) hasNextInOrder() bool {
	for {
		for it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.pnode = it.pnode.left
		}

		l := len(it.nodeChain)
		if l > 0 {
			it.pnode = it.nodeChain[l-1].right
			it.prenode = it.nodeChain[l-1]
			it.nodeChain = it.nodeChain[:l-1] // pop
			return true
		}
		if it.pnode == nil && l == 0 {
			return false
		}
	}
}

func (it *Iter3) hasNextPostOrder() bool {
	for {
		for it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.pnode = it.pnode.left
		}

		l := len(it.nodeChain)
		if l > 0 {
			it.pnode = it.nodeChain[l-1].right
			if it.pnode == nil || it.prenode == it.pnode {
				it.prenode = it.nodeChain[l-1]
				it.nodeChain = it.nodeChain[:l-1] // pop
				it.pnode = nil
				return true
			}
		}
		if it.pnode == nil && l == 0 {
			return false
		}
	}
}

// Value return current node data
func (it *Iter3) Value() int {
	return it.prenode.data
}

// Iter4 is a tree iterator, use chan and loop
type Iter4 struct {
	pnode     *TreeNode
	prenode   *TreeNode
	nodeChain []*TreeNode
	order     Method
	ch        chan int
	value     int
}

// NewTreeIter4 construct a new Tree Iterator
func NewTreeIter4(t *TreeNode, m Method) *Iter4 {
	var c = make(chan int)
	it := &Iter4{pnode: t, order: m, ch: c}
	go func() {
		if it.order == PreOrder {
			it.hasNextPreOrder()
		} else if it.order == InOrder {
			it.hasNextInOrder()
		} else if it.order == PostOrder {
			it.hasNextPostOrder()
		}
		close(it.ch)
	}()
	return it
}

// Next return true if there is still valid data exist
func (it *Iter4) Next() bool {
	for v := range it.ch {
		it.value = v
		return true
	}
	return false
}

func (it *Iter4) hasNextPreOrder() {
	for {
		if it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.ch <- it.pnode.data
			it.pnode = it.pnode.left
			continue
		}
		l := len(it.nodeChain)
		if l > 0 {
			it.pnode = it.nodeChain[l-1].right
			it.nodeChain = it.nodeChain[:l-1] // pop
		}
		if it.pnode == nil && l == 0 {
			return
		}
	}
}

func (it *Iter4) hasNextInOrder() {
	for {
		for it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.pnode = it.pnode.left
		}

		l := len(it.nodeChain)
		if l > 0 {
			it.ch <- it.nodeChain[l-1].data
			it.pnode = it.nodeChain[l-1].right
			it.nodeChain = it.nodeChain[:l-1] // pop
			continue
		}
		if it.pnode == nil && l == 0 {
			return
		}
	}
}

func (it *Iter4) hasNextPostOrder() {
	for {
		for it.pnode != nil {
			it.nodeChain = append(it.nodeChain, it.pnode) // push
			it.pnode = it.pnode.left
		}

		l := len(it.nodeChain)
		if l > 0 {
			it.pnode = it.nodeChain[l-1].right
			if it.pnode == nil || it.prenode == it.pnode {
				it.ch <- it.nodeChain[l-1].data
				it.prenode = it.nodeChain[l-1]
				it.nodeChain = it.nodeChain[:l-1] // pop
				it.pnode = nil
				continue
			}
		}
		if it.pnode == nil && l == 0 {
			return
		}
	}
}

// Value return current node data
func (it *Iter4) Value() int {
	return it.value
}
