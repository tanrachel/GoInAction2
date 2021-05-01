package vApp

import (
	"fmt"
)

//Venue Struct - holds pointer to LinkedList to hold BookingNodes
type Venue struct {
	Name         string
	Location     string
	VenueType    string
	Capacity     int
	Availability *map[int][]*LinkedList
}

func (v *Venue) printVenue() {
	fmt.Println("Name: ", v.Name)
	fmt.Println("Location: ", v.Location)
	fmt.Println("Type: ", v.VenueType)
	fmt.Println("Booked: ", v.Availability)
	fmt.Println(" ")
}

func (v *Venue) printSchedule(month int, day int) {
	monthMap := *v.Availability
	dayArray, ok := monthMap[month]
	if !ok {
		fmt.Println("Available All day!")
		return
	} else {
		if dayArray[day] == nil {
			fmt.Println("Available All day!")
			return
		} else {
			if (*((*v.Availability)[month])[day]).head == nil {
				fmt.Println("Available All day!")
				return
			} else {
				(*((*v.Availability)[month])[day]).printList()
			}
		}
	}
}

//if key doesn't exist, assume available all day, if day is nil assume available all day
func (v *Venue) availableOnDate(month int, day int) bool {
	monthMap := *v.Availability
	_, monthExist := monthMap[month]
	if !monthExist {
		return true
	}
	dayArray := monthMap[month]
	if dayArray[day] == nil {
		return true
	}
	return false

}

//add booking given a BookingNode is already created
func (v *Venue) AddBooking(b *BookingNode, month int, day int) bool {
	monthMap := *v.Availability
	_, monthExist := monthMap[month]
	if !monthExist {
		dayArray := make([]*LinkedList, 31)
		monthMap[month] = dayArray
	}
	dayArray := monthMap[month]
	if dayArray[day] == nil {
		dayArray[day] = &LinkedList{}
	}
	return (*((*v.Availability)[month])[day]).add(b)
}

//delete node matching on name and time
func (v *Venue) RemoveBooking(removeName string, bookTime int, month int, day int) (bool, *BookingNode) {
	monthMap := *v.Availability
	_, monthExist := monthMap[month]
	if !monthExist {
		fmt.Println("There is no booking in this month!")
		return false, nil
	}
	dayArray := monthMap[month]
	if dayArray[day] == nil {
		fmt.Println("There are no bookings on this day!")
		return false, nil
	}
	removedBooking := (*((*v.Availability)[month])[day]).delete(removeName, bookTime)
	if removedBooking == nil {
		return false, nil
	} else {
		return true, removedBooking
	}

}
