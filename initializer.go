/*
---initializer.go----
seeding the application with some test venues and booking
*/
package main

import (
	vApp "AssignmentV2/vApp"
	"time"
)

func initializer(m *map[string]*vApp.AVLTree) {
	venue1 := &vApp.Venue{"Shangri-La", "Bukit Timah", "hotel", 10, &map[int][]*vApp.LinkedList{}}
	venue11 := &vApp.Venue{"Pegasus", "Bukit Timah", "hotel", 50, &map[int][]*vApp.LinkedList{}}
	venue12 := &vApp.Venue{"Gallop Hill", "Bukit Timah", "hotel", 100, &map[int][]*vApp.LinkedList{}}
	venue13 := &vApp.Venue{"VIP", "Bukit Timah", "hotel", 70, &map[int][]*vApp.LinkedList{}}
	venue14 := &vApp.Venue{"Guesthouse", "Bukit Timah", "hotel", 70, &map[int][]*vApp.LinkedList{}}
	venue15 := &vApp.Venue{"D'Hotel", "Bukit Timah", "hotel", 70, &map[int][]*vApp.LinkedList{}}

	tempBooking1 := &vApp.BookingNode{"Rachel", time.Now().Local(), []int{6, 8}, nil}
	tempBooking2 := &vApp.BookingNode{"Jessica", time.Now().Local(), []int{9, 10}, nil}
	tempBooking3 := &vApp.BookingNode{"Rachel", time.Now().Local(), []int{12, 13}, nil}
	tempBooking4 := &vApp.BookingNode{"John", time.Now().Local(), []int{15, 17}, nil}
	tempBooking5 := &vApp.BookingNode{"Cam", time.Now().Local(), []int{20, 22}, nil}

	venue1.AddBooking(tempBooking1, 5, 1)
	venue1.AddBooking(tempBooking2, 5, 1)
	venue1.AddBooking(tempBooking3, 5, 1)
	venue1.AddBooking(tempBooking4, 6, 1)
	venue1.AddBooking(tempBooking5, 6, 1)

	addNewVenue(venue1)
	addNewVenue(venue11)
	addNewVenue(venue12)
	addNewVenue(venue13)
	addNewVenue(venue14)
	addNewVenue(venue15)

	venue2 := &vApp.Venue{"La Braceria", "Bukit Timah", "restaurant", 20, &map[int][]*vApp.LinkedList{}}
	venue3 := &vApp.Venue{"Caruso", "Bukit Timah", "restaurant", 30, &map[int][]*vApp.LinkedList{}}
	venue4 := &vApp.Venue{"LOMBA", "Bukit Timah", "restaurant", 70, &map[int][]*vApp.LinkedList{}}
	venue5 := &vApp.Venue{"Fratini La Trattoria", "Bukit Timah", "restaurant", 50, &map[int][]*vApp.LinkedList{}}
	venue6 := &vApp.Venue{"Ampang", "Bukit Timah", "hotel", 70, &map[int][]*vApp.LinkedList{}}
	venue7 := &vApp.Venue{"Pasta Fresca", "Bukit Timah", "restaurant", 75, &map[int][]*vApp.LinkedList{}}
	tempBooking11 := &vApp.BookingNode{"Rachel", time.Now().Local(), []int{6, 8}, nil}
	tempBooking22 := &vApp.BookingNode{"Jessica", time.Now().Local(), []int{9, 10}, nil}
	tempBooking33 := &vApp.BookingNode{"Rachel", time.Now().Local(), []int{12, 13}, nil}
	tempBooking44 := &vApp.BookingNode{"John", time.Now().Local(), []int{15, 17}, nil}
	tempBooking55 := &vApp.BookingNode{"Cam", time.Now().Local(), []int{20, 22}, nil}
	venue2.AddBooking(tempBooking11, 5, 1)
	venue2.AddBooking(tempBooking22, 5, 1)
	venue2.AddBooking(tempBooking33, 5, 1)
	venue2.AddBooking(tempBooking44, 6, 1)
	venue2.AddBooking(tempBooking55, 6, 1)
	((*m)["restaurant"]) = &vApp.AVLTree{}
	(*((*m)["restaurant"])).Add(venue2.Capacity, venue2)
	(*((*m)["restaurant"])).Add(venue3.Capacity, venue3)
	(*((*m)["restaurant"])).Add(venue4.Capacity, venue4)
	(*((*m)["restaurant"])).Add(venue5.Capacity, venue5)
	(*((*m)["restaurant"])).Add(venue6.Capacity, venue6)
	(*((*m)["restaurant"])).Add(venue7.Capacity, venue7)

}
