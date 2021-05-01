package vApp

import (
	"fmt"
	"time"
)

type BookingNode struct {
	Bookedby string
	BookedOn time.Time
	Book     []int
	Next     *BookingNode
}

func (b *BookingNode) printBooking() {
	fmt.Println("Booked by:", b.Bookedby)
	fmt.Println("Booked on: ", b.BookedOn)
	fmt.Println("Booked time: ", b.Book[0], "-", b.Book[1])
}
