package vApp

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex

//AVL Tree Implementation
type Node struct {
	val    int
	item   *Venue
	left   *Node
	right  *Node
	height int
}

type AVLTree struct {
	Root *Node
}

// wrapper to use on the front end
func (a *AVLTree) Add(val int, item *Venue) {
	a.Root = a.insert(a.Root, val, item)
}
func (a *AVLTree) remove(val int) {
	a.Root = a.delete(a.Root, val)
}
func (a *AVLTree) update(oldVal int, newVal int) {
	a.Root = a.delete(a.Root, oldVal)
	a.Root = a.insert(a.Root, newVal, nil)
}

// AVL function implementation
func (a *AVLTree) insert(root *Node, val int, item *Venue) *Node {
	if root == nil {
		return &Node{val, item, nil, nil, 1}
	} else if val < root.val {
		root.left = a.insert(root.left, val, item)
	} else {
		root.right = a.insert(root.right, val, item)
	}
	root.height = 1 + a.max(a.height(root.left), a.height(root.right))
	balance := a.getBalance(root)

	if balance > 1 && val < root.left.val {
		return a.rightRotate(root)
	}
	if balance < -1 && val > root.right.val {
		return a.leftRotate(root)
	}
	if balance > 1 && val > root.left.val {
		root.left = a.leftRotate(root.left)
		return a.rightRotate(root)
	}
	if balance < -1 && val < root.right.val {
		root.right = a.rightRotate(root.right)
		return a.leftRotate(root)
	}
	return root
}

func (a *AVLTree) delete(root *Node, val int) *Node {
	if root == nil {
		return root
	} else if val < root.val {
		root.left = a.delete(root.left, val)
	} else if val > root.val {
		root.right = a.delete(root.right, val)
	} else {
		if root.left == nil {
			temp := root.right
			root = nil
			return temp
		} else if root.right == nil {
			temp := root.left
			root = nil
			return temp
		}
		temp := a.minValue(root.right)
		root.val = temp.val
		root.right = a.delete(root.right, temp.val)

	}
	if root == nil {
		return root
	}
	root.height = 1 + a.max(a.height(root.left), a.height(root.right))
	balance := a.getBalance(root)
	if balance > 1 && a.getBalance(root.left) >= 0 {
		return a.rightRotate(root)
	}
	if balance < -1 && a.getBalance(root.right) <= 0 {
		return a.leftRotate(root)
	}
	if balance > 1 && a.getBalance(root.left) < 0 {
		root.left = a.leftRotate(root.left)
		return a.rightRotate(root)
	}
	if balance < -1 && a.getBalance(root.right) > 0 {
		root.right = a.rightRotate(root.right)
		return a.leftRotate(root)
	}
	return root
}

func (a *AVLTree) minValue(root *Node) *Node {
	if root == nil || root.left == nil {
		return root
	}
	return a.minValue(root.left)
}
func (a *AVLTree) max(c int, b int) int {
	if c > b {
		return c
	} else {
		return b
	}
}

func (a *AVLTree) height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.height
}

func (a *AVLTree) getBalance(root *Node) int {
	if root == nil {
		return 0
	}
	return a.height(root.left) - a.height(root.right)
}

func (a *AVLTree) rightRotate(n *Node) *Node {
	y := n.left
	t3 := y.right
	y.right = n
	n.left = t3
	n.height = 1 + a.max(a.height(n.left), a.height(n.right))
	y.height = 1 + a.max(a.height(y.left), a.height(y.right))
	return y
}

func (a *AVLTree) leftRotate(n *Node) *Node {
	y := n.right
	t2 := y.left
	y.left = n
	n.right = t2
	n.height = 1 + a.max(a.height(n.left), a.height(n.right))
	y.height = 1 + a.max(a.height(y.left), a.height(y.right))
	return y
}
func (a *AVLTree) searchCapacity(root *Node, target int) {
	if root == nil {
		return
	}
	if target < root.val {
		a.searchCapacity(root.left, target)
	} else if target > root.val {
		a.searchCapacity(root.right, target)
	} else {
		fmt.Println((*(root.item)).Name, "-", (*(root.item)).Location, "-", (*(root.item)).Capacity)
		a.searchCapacity(root.right, target)
	}
}
func (a *AVLTree) SearchByName(root *Node, venueName string, target int) *Venue {
	if root == nil {
		return nil
	}
	if target < root.val {
		return a.SearchByName(root.left, venueName, target)
	} else if target > root.val {
		return a.SearchByName(root.right, venueName, target)
	} else {
		if (*(root.item)).Name != venueName {
			return a.SearchByName(root.right, venueName, target)
		} else {
			return root.item
		}
	}
	return root.item
}
func (a *AVLTree) SearchByDateAvailableWrapper(month int, day int, array *[]*Venue) {
	a.searchByDateAvailable(a.Root, month, day, array)
}
func (a *AVLTree) searchByDateAvailable(root *Node, month int, day int, array *[]*Venue) {
	if root == nil {
		return
	}
	if (*(root.item)).availableOnDate(month, day) {
		// fmt.Println((*(root.item)).Name, "-", (*(root.item)).Location, "-", (*(root.item)).Capacity)
		arrayReader(array, &(*(root.item)))
	}
	a.searchByDateAvailable(root.left, month, day, array)
	a.searchByDateAvailable(root.right, month, day, array)

}
func (a *AVLTree) printPreOrder() {
	a.preOrder(a.Root)
}

//print tree for verification of tree structure
func (a *AVLTree) preOrder(root *Node) {
	if root == nil {
		fmt.Println("EMPTY")
		return
	}
	fmt.Println(root.val, " - ", (*(root.item)).Name)
	fmt.Println("ON THE LEFT")
	a.preOrder(root.left)
	fmt.Println("ON THE RIGHT")
	a.preOrder(root.right)

}
func (a *AVLTree) printTree() {
	a.printTreeRec(a.Root)
}
func (a *AVLTree) printTreeRec(root *Node) {
	if root == nil {
		return
	}
	fmt.Println((*(root.item)).Name, "-", (*(root.item)).Location, "-", (*(root.item)).Capacity)
	a.printTreeRec(root.left)
	a.printTreeRec(root.right)

}
func (a *AVLTree) ReturnTree(results *[]*Venue) {
	a.returnTreeRec(a.Root, results)
}
func (a *AVLTree) returnTreeRec(root *Node, array *[]*Venue) {
	if root == nil {
		return
	}
	// fmt.Println((*(root.item)).name, "-", (*(root.item)).location, "-", (*(root.item)).capacity)
	a.returnTreeRec(root.left, array)
	a.returnTreeRec(root.right, array)
	// fmt.Println("DOUBLE CHECKING", root.item)
	arrayReader(array, &(*(root.item)))

}
func arrayReader(finalList *[]*Venue, input *Venue) {
	mutex.Lock()
	*finalList = append(*finalList, input)
	mutex.Unlock()
}
func (a *AVLTree) ReturnCapacity(target int, results *[]*Venue) {
	a.returnSearchCapacity(a.Root, target, results)
}
func (a *AVLTree) returnSearchCapacity(root *Node, target int, array *[]*Venue) {
	if root == nil {
		return
	}
	if target < root.val {
		a.returnSearchCapacity(root.left, target, array)
	} else if target > root.val {
		a.returnSearchCapacity(root.right, target, array)
	} else {
		arrayReader(array, &(*(root.item)))
		a.returnSearchCapacity(root.right, target, array)
	}
}
