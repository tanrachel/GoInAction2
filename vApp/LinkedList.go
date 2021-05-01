package vApp

import (
	"fmt"
)

//LinkedList implementation
type LinkedList struct {
	head *BookingNode
	size int
}

// Sorted Add base on time book
func (l *LinkedList) add(b *BookingNode) bool {
	//if Linkedlist is empty, add the node to the head
	if l.head == nil {
		l.head = b
	} else {
		// evaluate whether b's time frame overlaps with other bookings in the list
		bStart := b.Book[0]
		bEnd := b.Book[1]
		currentNode := l.head
		previousNode := l.head
		//if the b is the earliest booking, add to the front
		if bStart < currentNode.Book[0] && bEnd <= currentNode.Book[1] {
			b.Next = l.head
			l.head = b
			return true
		}
		//start traversing through the list to look the node to insert in front
		for currentNode.Next != nil {
			currStart := currentNode.Book[0]
			//if currentNode is still less than b, keep traversing
			if bStart > currStart {
				previousNode = currentNode
				currentNode = currentNode.Next
			} else {
				//if we hit a node who's time is larger than ours, we check to see if both b's time does not overlap with current and previous bookings
				if previousNode.Book[1] <= bStart && currStart >= bEnd {
					previousNode.Next = b
					b.Next = currentNode
					return true
				} else {
					//return false if it does overlap
					return false
				}
			}

		}
		//take into account we might be at the last node during traversal
		// this is if b is the largest so far in the list, add to the end
		if bStart >= currentNode.Book[1] {
			currentNode.Next = b
			return true
		} else {
			//add infront of the last node
			if previousNode.Book[1] <= bStart && currentNode.Book[0] >= bEnd {
				previousNode.Next = b
				b.Next = currentNode
				return true
			} else {
				return false
			}
		}
	}
	l.size++
	return true
}

func (l *LinkedList) delete(targetName string, startTime int) *BookingNode {
	if l.head == nil {
		return nil
	} else {
		currentNode := l.head
		previousNode := l.head
		if currentNode.Bookedby == targetName && currentNode.Book[0] == startTime {
			l.head = currentNode.Next
			return currentNode
		}
		for currentNode.Next != nil {
			if currentNode.Bookedby == targetName && currentNode.Book[0] == startTime {
				previousNode.Next = currentNode.Next
				return currentNode
			} else if currentNode.Book[0] < startTime {
				previousNode = currentNode
				currentNode = currentNode.Next
			} else {
				return nil
			}
		}
		if currentNode.Bookedby == targetName && currentNode.Book[0] == startTime {
			previousNode.Next = nil
			return currentNode
		}
		return nil
	}
}
func (l *LinkedList) printList() {
	if l.head == nil {
		fmt.Println("EMPTY STRING!")
	} else {
		currentNode := l.head
		for currentNode != nil {
			currentNode.printBooking()
			currentNode = currentNode.Next
		}
	}
}
