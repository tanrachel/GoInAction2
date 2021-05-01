/*
----functionHelper.go-----
functions that are used in handlers in main.go are placed here.
*/
package main

import (
	vApp "AssignmentV2/vApp"
	"errors"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func browseVenues() []*vApp.Venue {
	// resultArray := []*vApp.Venue{}
	// for _, val := range myMap {
	// 	tree := *val
	// 	tree.returnTree(&resultArray)
	// }
	// for _, value := range resultArray {
	// 	(*value).printVenue()
	// }
	// return resultArray

	typeSlice := []string{}
	for key, _ := range myMap {
		typeSlice = append(typeSlice, key)
	}
	data := []*vApp.Venue{}
	taskLoad := len(typeSlice)
	tasks := make(chan string, taskLoad)
	wg.Add(4)

	for gr := 1; gr <= 4; gr++ {
		go browseTree(tasks, &data)

	}
	for _, venueType := range typeSlice {
		tasks <- venueType
	}
	close(tasks)
	wg.Wait()
	return data
}
func browseTree(tasks chan string, array *[]*vApp.Venue) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}
		((*(&myMap))[task]).ReturnTree(array)
	}

}
func printByTypeCheck(target string, evaluate []string) error {
	for _, val := range evaluate {
		if val == target {
			return nil
		}
	}
	return errors.New("Invalid type entered!")
}
func printByCapacityCheck(target int) error {
	_, ok := capacityMap[target]
	if !ok {
		return errors.New("Invalid capacity entered!")
	}
	return nil
}

func searchTreeConc(tasks chan string, capacity int, array *[]*vApp.Venue) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}
		((*(&myMap))[task]).ReturnCapacity(capacity, array)
	}

}
func dateChecker(month int, day int) error {
	_, todaymonth, todayday := time.Now().Date()
	if month == int(todaymonth) && day >= todayday && day < 32 {
		return nil
	} else if month > int(todaymonth) && day > 0 && day < 32 {
		return nil
	} else {
		return errors.New("Sorry, invalid date! Please enter a month and day this year")
	}
}

func capacityMapper(v *vApp.Venue) {
	_, ok := capacityMap[v.Capacity]
	if !ok {
		capacityMap[v.Capacity] = true
	}
}

func searchForAvailabilityConc(tasks chan string, month int, day int, array *[]*vApp.Venue) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}
		((*(&myMap))[task]).SearchByDateAvailableWrapper(month, day, array)
	}
}
func searchForVenue(a *(vApp.AVLTree), target string, capacity int) (*(vApp.Venue), error) {
	found := a.SearchByName(a.Root, target, capacity)
	if found == nil {
		return nil, errors.New("Venue isn't found, are you sure you entered the correct information?")
	}
	return found, nil
}
func findVenue(venue string, capacity int, venueType string) (*vApp.Venue, error) {
	typeSlice := []string{}
	for key, _ := range myMap {
		typeSlice = append(typeSlice, key)
	}
	err1, err2 := printByTypeCheck(venueType, typeSlice), printByCapacityCheck(capacity)
	if err1 != nil && err2 != nil {
		return nil, errors.New("Venue not found!")
	} else {
		foundVenue, err := searchForVenue(myMap[venueType], venue, capacity)
		return foundVenue, err
	}
}
func timeChecker(timeStart, timeEnd int) bool {
	if timeEnd <= timeStart {
		return false
	} else if timeStart < 0 || timeEnd < 0 || timeStart > 24 || timeEnd > 24 {
		return false
	} else {
		return true
	}
}

func getUser(res http.ResponseWriter, req *http.Request) user {
	c, err := req.Cookie("session")
	if err != nil {
		id := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: id.String(),
		}
		http.SetCookie(res, c)
	}
	var u user
	if e, ok := mapSessions[c.Value]; ok {
		// u = mapUsers[e]
		u, _ = queryUser(mapUsers, e)
	}
	return u
}

func loggedInBool(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	e := mapSessions[c.Value]
	// _, ok := mapUsers[e]
	_, ok := queryUser(mapUsers, e)
	fmt.Println("loggedInBool:", e, ok, c.Value)
	return ok
}

func headerPasser(res http.ResponseWriter, req *http.Request, data *map[string]interface{}) (bool, user) {
	loggedIn := loggedInBool(req)
	thisUser := getUser(res, req)
	getMap := *data
	if loggedIn {
		getMap["logInStatus"] = "1"
		getMap["username"] = thisUser.UserName
	} else {
		getMap["logInStatus"] = "0"
	}
	return loggedIn, thisUser
}
func adminAdder() {
	// var u user
	un := "admin"
	pw := "1234"
	bs, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	// u = user{un, string(bs)}
	// mapUsers[un] = u
	insertRecord(mapUsers, un, string(bs))
}

func showBookingsConc(tasks chan string, data *map[string][][]string) {
	defer wg.Done()
	bookings := [][]string{}
	for {
		task, ok := <-tasks
		if !ok {
			return
		}
		((*(&myMap))[task]).FetchTreeBookingsWrapper(&bookings)
		for _, booking := range bookings {
			vName := booking[6]
			if _, ok := (*data)[vName]; ok {
				mutex.Lock()
				(*data)[vName] = append((*data)[vName], booking)
				mutex.Unlock()
			} else {
				temp := [][]string{booking}
				mutex.Lock()
				(*data)[vName] = temp
				mutex.Unlock()
			}
		}
	}
}

func addNewVenue(v *vApp.Venue) {
	_, ok := ((*(&myMap))[v.VenueType])
	if !ok {
		((*(&myMap))[v.VenueType]) = &vApp.AVLTree{}
	}
	(*((*(&myMap))[v.VenueType])).Add(v.Capacity, v)
	capacityMapper(v)
}
