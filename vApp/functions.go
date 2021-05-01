package vApp

import (
	"fmt"
	"strconv"
)

//--new methods and functions to sweep the entire database for bookings
func (b *BookingNode) returnBookingInfo() (string, string, int, int) {
	return b.Bookedby, b.BookedOn.String(), b.Book[0], b.Book[1]
}
func (l *LinkedList) traverseBookings(data *[][]string, month int, day int, vName string) {
	if l.head == nil {
		fmt.Println("EMPTY STRING!")
	} else {
		currentNode := l.head
		for currentNode != nil {
			name, bookedon, starttime, endtime := currentNode.returnBookingInfo()
			//[month,day,name,bookedon,starttime,endtime,venueName]
			temp := []string{strconv.Itoa(month), strconv.Itoa(day), name, bookedon, strconv.Itoa(starttime), strconv.Itoa(endtime), vName}
			mutex.Lock()
			*data = append(*data, temp)
			mutex.Unlock()
			// fmt.Println("INSIDE LLIST:", *data)
			currentNode = currentNode.Next
		}
	}
}
func (v *Venue) fetchBookings(bookings *[][]string) {
	//this can be improved by moving concurrency to this level
	monthMap := *v.Availability
	for i := 1; i <= 12; i++ {
		dayArray, ok := monthMap[i]
		if ok {
			for j := 1; j < 31; j++ {
				if dayArray[j] != nil {
					if (*((*v.Availability)[i])[j]).head != nil {
						// fmt.Println("oneLevel Up:", *bookings)
						(*((*v.Availability)[i])[j]).traverseBookings(bookings, i, j, v.Name)
					}
				}
			}
		}
	}
}
func (a *AVLTree) fetchTreeBookings(root *Node, bookings *[][]string) {
	if root == nil {
		return
	}
	// fmt.Println((*(root.item)).name, "-", (*(root.item)).location, "-", (*(root.item)).capacity)
	(*(root.item)).fetchBookings(bookings)
	a.fetchTreeBookings(root.left, bookings)
	a.fetchTreeBookings(root.right, bookings)

}
func (a *AVLTree) FetchTreeBookingsWrapper(bookings *[][]string) {
	a.fetchTreeBookings(a.Root, bookings)
	// fmt.Println("Tree Level:", *bookings)
}
