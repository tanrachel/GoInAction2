/*
-------Main.go-------
host the main function
global variables are declared here
server is declared here as well as handlers that serve the pages
*/
package main

import (
	vApp "AssignmentV2/vApp"
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	bcrypt "golang.org/x/crypto/bcrypt"
)

var (
	userInput   string
	myMap       = make(map[string]*vApp.AVLTree)
	capacityMap = make(map[int]bool)
	reader      = bufio.NewReader(os.Stdin)
	wg          sync.WaitGroup
	mutex       sync.Mutex
	tpl         *template.Template
	// mapUsers    = map[string]user{}
	mapUsers    *sql.DB
	mapSessions = map[string]string{}
)

type user struct {
	UserName string
	Password string
}

func main() {
	var err error
	mapUsers, err = sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/VenueDB")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
	adminAdder()
	defer mapUsers.Close()
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/browse", browse)
	http.HandleFunc("/searchMenu", searchMenu)
	http.HandleFunc("/book", book)
	http.HandleFunc("/showbook", showbookings)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/delete", deleteBooking)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/admin/addvenue", addVenue)
	http.HandleFunc("/admin/deleteSessions", deleteSessions)
	http.HandleFunc("/admin/deleteUsers", deleteUsers)
	http.ListenAndServe(":8080", nil)
}

// add in concurrency here
func showbookings(res http.ResponseWriter, req *http.Request) {
	m := make(map[string]interface{})
	headerPasser(res, req, &m)
	dataDump := make(map[string][][]string)
	////-start concurrency
	// for _, tree := range myMap {
	// 	gotoTree := *tree
	// 	bookings := [][]string{}
	// 	gotoTree.fetchTreeBookingsWrapper(&bookings)
	// 	for _, booking := range bookings {
	// 		vName := booking[6]
	// 		if _, ok := dataDump[vName]; ok {
	// 			dataDump[vName] = append(dataDump[vName], booking)
	// 		} else {
	// 			temp := [][]string{booking}
	// 			dataDump[vName] = temp
	// 		}
	// 	}
	// }
	////-end concurrency
	typeSlice := []string{}
	for key, _ := range myMap {
		typeSlice = append(typeSlice, key)
	}
	taskLoad := len(typeSlice)
	tasks := make(chan string, taskLoad)
	wg.Add(4)

	for gr := 1; gr <= 4; gr++ {
		//concurrency go go deeper to
		go showBookingsConc(tasks, &dataDump)

	}
	for _, venueType := range typeSlice {
		tasks <- venueType
	}
	close(tasks)
	wg.Wait()

	/////
	if len(dataDump) == 0 {
		m["no-data"] = "No bookings in the system!"
	} else {
		m["data"] = dataDump
	}
	tpl.ExecuteTemplate(res, "showBookings.gohtml", m)
}
func admin(res http.ResponseWriter, req *http.Request) {
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if userInfo.UserName == "admin" {
			tpl.ExecuteTemplate(res, "admin.gohtml", m)
		} else {
			m["notLoggedInError"] = "Sorry you have to be an admin to access!"
			tpl.ExecuteTemplate(res, "index.gohtml", m)
		}
	}
}
func addVenue(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if userInfo.UserName == "admin" {
			//addVenue here
			if req.Method == http.MethodPost {
				venueName := req.Form["venue"][0]
				location := req.Form["location"][0]
				capacity, _ := strconv.Atoi(req.Form["capacity"][0])
				venueType := req.Form["venuetype"][0]
				// venueName := req.Form["venue"]
				// location := req.Form["location"]
				// capacity := req.Form["capacity"]
				// venueType := req.Form["venuetype"]
				fmt.Println(venueName, location, venueType, capacity)
				tempV := &vApp.Venue{venueName, location, venueType, capacity, &map[int][]*vApp.LinkedList{}}
				mutex.Lock()
				addNewVenue(tempV)
				mutex.Unlock()
				m["successAction"] = "Successfully added venue!"
				tpl.ExecuteTemplate(res, "admin.gohtml", m)
			} else {
				tpl.ExecuteTemplate(res, "addVenue.gohtml", m)
			}
		} else {
			m["notLoggedInError"] = "Sorry you have to be an admin to access!"
			tpl.ExecuteTemplate(res, "index.gohtml", m)
		}
	}
}
func deleteSessions(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if userInfo.UserName == "admin" {

			if req.Method == http.MethodPost {
				chosenSession := req.Form["chosenSession"]
				if len(chosenSession) > 0 {
					mutex.Lock()
					delete(mapSessions, chosenSession[0])
					mutex.Unlock()
					c := &http.Cookie{
						Name:   "session",
						Value:  "",
						MaxAge: -1,
					}
					http.SetCookie(res, c)
					m["successAction"] = "Session successfully deleted!"
					headerPasser(res, req, &m)
					tpl.ExecuteTemplate(res, "index.gohtml", m)
				} else {
					m["error"] = "Error in submission!"
					tpl.ExecuteTemplate(res, "index.gohtml", m)
				}

			} else {
				tempList := [][]string{}
				for key, value := range mapSessions {
					temp := []string{key, value}
					tempList = append(tempList, temp)
				}
				m["sessionData"] = tempList
				tpl.ExecuteTemplate(res, "deleteSession.gohtml", m)
			}
		} else {
			m["notLoggedInError"] = "Sorry you have to be an admin to access!"
			tpl.ExecuteTemplate(res, "index.gohtml", m)
		}
	}
}
func deleteUsers(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if userInfo.UserName == "admin" {

			if req.Method == http.MethodPost {
				chosenUser := req.Form["chosenUser"]
				if len(chosenUser) > 0 {
					if chosenUser[0] == "admin" {
						m["notLoggedInError"] = "Admin cannot be deleted!"
						tpl.ExecuteTemplate(res, "index.gohtml", m)
						return
					} else {
						mutex.Lock()
						// delete(mapUsers, chosenUser[0])
						deleteUserDB(mapUsers, chosenUser[0])
						mutex.Unlock()
						m["successAction"] = "User successfully deleted!"
						headerPasser(res, req, &m)
						tpl.ExecuteTemplate(res, "index.gohtml", m)
					}
				} else {
					m["error"] = "Error in selection!"
					tpl.ExecuteTemplate(res, "index.gohtml", m)
				}
			} else {
				// tempList := []string{}
				// for key, _ := range mapUsers {
				// 	tempList = append(tempList, key)
				// }
				allUsers, _ := mapUsers.Query("SELECT * FROM VenueDB.Users")
				tempList := []string{}
				for allUsers.Next() {
					var person user
					allUsers.Scan(&person.UserName, &person.Password)
					tempList = append(tempList, person.UserName)
				}
				m["userData"] = tempList
				tpl.ExecuteTemplate(res, "deleteUsers.gohtml", m)
			}
		} else {
			m["notLoggedInError"] = "Sorry you have to be an admin to access!"
			tpl.ExecuteTemplate(res, "index.gohtml", m)
		}
	}
}
func init() {
	funcMap := template.FuncMap{
		"deref": func(i *vApp.Venue) vApp.Venue { return *i },
	}
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
	initializer(&myMap)
}
func login(res http.ResponseWriter, req *http.Request) {
	if loggedInBool(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	m := make(map[string]interface{})
	if req.Method == http.MethodPost {
		var u user
		un := req.FormValue("username")
		pw := req.FormValue("password")
		// u, ok := mapUsers[un]
		u, ok := queryUser(mapUsers, un)
		if !ok {
			m["invalidlogin"] = "Username doesn't exist!"
			tpl.ExecuteTemplate(res, "userForm.gohtml", m)
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
		if err != nil {
			m["invalidlogin"] = "Username and password don't match!!"
			tpl.ExecuteTemplate(res, "userForm.gohtml", m)
			return
		}
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(res, c)
		mutex.Lock()
		mapSessions[c.Value] = un
		mutex.Unlock()
		m["logInStatus"] = "1"
		m["user"] = u.UserName
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	if req.Method == http.MethodGet {
		m["login"] = "1"
	}
	tpl.ExecuteTemplate(res, "userForm.gohtml", m)
}
func signup(res http.ResponseWriter, req *http.Request) {
	if loggedInBool(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// var u user
	m := make(map[string]interface{})
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pw := req.FormValue("password")
		if len(pw) == 0 || len(un) == 0 {
			m["invalid"] = "Invalid Username or Password!"
			tpl.ExecuteTemplate(res, "userForm.gohtml", m)
			return
		}
		// if _, ok := mapUsers[un]; ok {
		// 	m["invalid"] = "Username taken! Please try another!"
		// 	tpl.ExecuteTemplate(res, "userForm.gohtml", m)
		// 	return
		// }
		if _, ok := queryUser(mapUsers, un); ok {
			m["invalid"] = "Username taken! Please try another!"
			tpl.ExecuteTemplate(res, "userForm.gohtml", m)
			return
		}
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(res, c)
		mapSessions[c.Value] = un
		bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
		if err != nil {
			m["invalid"] = "Sorry, there was error on my end, please try again!"
			tpl.ExecuteTemplate(res, "userForm.gohtml", m)
			return
		}
		// u = user{un, string(bs)}
		mutex.Lock()
		// mapUsers[un] = u
		insertRecord(mapUsers, un, string(bs))
		mutex.Unlock()
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	m["signup"] = "this is signup"
	tpl.ExecuteTemplate(res, "userForm.gohtml", m)
}
func logout(res http.ResponseWriter, req *http.Request) {
	if !loggedInBool(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	mutex.Lock()
	delete(mapSessions, c.Value)
	mutex.Unlock()
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, c)
	m := make(map[string]interface{})
	m["logInStatus"] = "0"
	tpl.ExecuteTemplate(res, "index.gohtml", m)
}

func index(res http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	headerPasser(res, req, &data)
	fmt.Println("CHECKING INDEX")
	for key, value := range data {
		fmt.Println(key, value)
	}

	err := tpl.ExecuteTemplate(res, "index.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}

}
func browse(res http.ResponseWriter, req *http.Request) {
	m := make(map[string]interface{})
	m["data"] = browseVenues()
	headerPasser(res, req, &m)
	err := tpl.ExecuteTemplate(res, "browsingResult.gohtml", m)
	if err != nil {
		log.Fatalln(err)
	}

}

func searchMenu(res http.ResponseWriter, req *http.Request) {
	m := make(map[string]interface{})
	headerPasser(res, req, &m)
	req.ParseForm()
	choice := req.Form["choice"]
	choice2 := req.Form["choice2"]
	choice3 := req.Form["choice3"]

	fmt.Println("IM HERE!", len(choice2), choice, choice2, choice3)
	if len(choice2) > 0 && len(choice3) > 0 {
		//search for date availability
		i, err1 := strconv.Atoi(choice2[0])
		j, err2 := strconv.Atoi(choice3[0])
		if err1 != nil || err2 != nil {
			m["error"] = "There was an error! Please try again!"
		} else {
			dateCheckErr := dateChecker(i, j)
			if dateCheckErr != nil {
				m["error"] = "Please enter a valid date!"
			} else {
				typeSlice := []string{}
				for key, _ := range myMap {
					typeSlice = append(typeSlice, key)
				}
				data := []*vApp.Venue{}
				taskLoad := len(typeSlice)
				tasks := make(chan string, taskLoad)
				wg.Add(4)

				for gr := 1; gr <= 4; gr++ {
					go searchForAvailabilityConc(tasks, i, j, &data)

				}
				for _, venueType := range typeSlice {
					tasks <- venueType
				}
				close(tasks)
				wg.Wait()
				m["data"] = data
			}
		}

	} else if len(choice) > 0 && len(choice2) == 0 {
		m["choice"] = choice[0]
		fmt.Println("OVER HERE", m["choice"])
		if choice[0] == "capacity" {
			capacitySlice := []int{}
			for key, _ := range capacityMap {
				capacitySlice = append(capacitySlice, key)
			}
			m["capacity"] = capacitySlice
		} else if choice[0] == "types" {
			typeSlice := []string{}
			for key, _ := range myMap {
				typeSlice = append(typeSlice, key)
			}
			m["types"] = typeSlice
		}
		// err := tpl.ExecuteTemplate(res, "searchResult.gohtml", m)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
	} else if len(choice2) > 0 && len(choice3) == 0 {
		a := strings.Split(choice2[0], " ")
		if len(a) == 1 {
			i, err := strconv.Atoi(a[0])
			//this is a type - search through tree by type
			if err != nil {
				typeSlice := []string{}
				for key, _ := range myMap {
					typeSlice = append(typeSlice, key)
				}
				err := printByTypeCheck(a[0], typeSlice)
				if err != nil {
					m["error"] = "There was an error! Please try again!"
				} else {
					data := []*vApp.Venue{}
					(*(myMap[strings.TrimSpace(a[0])])).ReturnTree(&data)
					m["data"] = data
				}
			} else {
				//this is if it's an integer - so capacity
				capacitySlice := []int{}
				for key, _ := range capacityMap {
					capacitySlice = append(capacitySlice, key)
				}
				err := printByCapacityCheck(i)
				if err != nil {
					m["error"] = "There was error! Please try again!"
				} else {
					typeSlice := []string{}
					for key, _ := range myMap {
						typeSlice = append(typeSlice, key)
					}
					data := []*vApp.Venue{}
					taskLoad := len(typeSlice)
					tasks := make(chan string, taskLoad)
					wg.Add(4)
					for gr := 1; gr <= 4; gr++ {
						go searchTreeConc(tasks, i, &data)
					}
					for _, venueType := range typeSlice {
						tasks <- venueType
					}
					close(tasks)
					wg.Wait()
					m["data"] = data
				}
			}
		}
	}
	for key, value := range m {
		fmt.Println(key, value)
	}
	err := tpl.ExecuteTemplate(res, "searchResult.gohtml", m)
	if err != nil {
		log.Fatalln(err)
	}
}

func book(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if req.Method == "GET" {
			err := tpl.ExecuteTemplate(res, "booking.gohtml", nil)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			venueName := req.Form["venue"][0]
			capacity, _ := strconv.Atoi(req.Form["capacity"][0])
			venueType := req.Form["venuetype"][0]
			month, _ := strconv.Atoi(req.Form["month"][0])
			day, _ := strconv.Atoi(req.Form["day"][0])
			starttime, _ := strconv.Atoi(req.Form["starttime"][0])
			endtime, _ := strconv.Atoi(req.Form["endtime"][0])
			foundVenue, err := findVenue(venueName, capacity, venueType)
			dateCheckErr := dateChecker(month, day)
			timeCheck := timeChecker(starttime, endtime)
			if err != nil || dateCheckErr != nil || timeCheck != true {
				m["error"] = "Error occurred! Please double check your inputs!"
			} else {
				tempBooking := &vApp.BookingNode{userInfo.UserName, time.Now().Local(), []int{starttime, endtime}, nil}
				mutex.Lock()
				sucess := foundVenue.AddBooking(tempBooking, month, day)
				mutex.Unlock()
				if sucess {
					m["success"] = "Successfully Added!"
				} else {
					m["error"] = "Error occurred! Please double check your inputs!"
				}
			}
			err = tpl.ExecuteTemplate(res, "results.gohtml", m)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func edit(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if req.Method == "GET" {
			err := tpl.ExecuteTemplate(res, "edit.gohtml", nil)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			if len(req.Form["month"]) == 2 && len(req.Form["day"]) == 2 && len(req.Form["starttime"]) == 2 && len(req.Form["endtime"]) == 2 {
				venueName := req.Form["venue"][0]
				capacity, _ := strconv.Atoi(req.Form["capacity"][0])
				venueType := req.Form["venuetype"][0]
				oldmonth, _ := strconv.Atoi(req.Form["month"][0])
				oldday, _ := strconv.Atoi(req.Form["day"][0])
				oldstarttime, _ := strconv.Atoi(req.Form["starttime"][0])
				// oldendtime, _ := strconv.Atoi(req.Form["endtime"][0])
				newmonth, _ := strconv.Atoi(req.Form["month"][1])
				newday, _ := strconv.Atoi(req.Form["day"][1])
				newstarttime, _ := strconv.Atoi(req.Form["starttime"][1])
				newendtime, _ := strconv.Atoi(req.Form["endtime"][1])
				foundVenue, err := findVenue(venueName, capacity, venueType)
				if err != nil {
					m["error"] = "Venue not found!!"
					fmt.Println("error here")
				} else {
					dateCheckErr := dateChecker(newmonth, newday)
					timeCheck := timeChecker(newstarttime, newendtime)
					mutex.Lock()
					success, _ := foundVenue.RemoveBooking(userInfo.UserName, oldstarttime, oldmonth, oldday)
					mutex.Unlock()
					if err != nil || dateCheckErr != nil || timeCheck != true {
						m["error"] = "Error occurred! Please double check your inputs!"
					} else if success == false {
						m["error"] = "Booking not found!"
					} else {
						tempBooking := &vApp.BookingNode{userInfo.UserName, time.Now().Local(), []int{newstarttime, newendtime}, nil}
						mutex.Lock()
						sucess := foundVenue.AddBooking(tempBooking, newmonth, newday)
						mutex.Unlock()
						if sucess {
							m["success"] = "Successfully Updated!!"
						} else {
							m["error"] = "Error occurred! Please double check your inputs!"
						}
					}
				}
			}
		}
		fmt.Println("hitting here")
		err := tpl.ExecuteTemplate(res, "results.gohtml", m)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
func deleteBooking(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	m := make(map[string]interface{})
	logIn, userInfo := headerPasser(res, req, &m)
	if !logIn {
		m["notLoggedInError"] = "Please login to do this!"
		tpl.ExecuteTemplate(res, "index.gohtml", m)
	} else {
		if req.Method == "GET" {
			err := tpl.ExecuteTemplate(res, "delete.gohtml", nil)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			venueName := req.Form["venue"][0]
			capacity, _ := strconv.Atoi(req.Form["capacity"][0])
			venueType := req.Form["venuetype"][0]
			oldmonth, _ := strconv.Atoi(req.Form["month"][0])
			oldday, _ := strconv.Atoi(req.Form["day"][0])
			oldstarttime, _ := strconv.Atoi(req.Form["starttime"][0])
			oldendtime, _ := strconv.Atoi(req.Form["endtime"][0])
			foundVenue, err := findVenue(venueName, capacity, venueType)
			if err != nil {
				m["error"] = "Venue not found!"
			} else {
				dateCheckErr := dateChecker(oldmonth, oldday)
				timeCheck := timeChecker(oldstarttime, oldendtime)
				mutex.Lock()
				success, _ := foundVenue.RemoveBooking(userInfo.UserName, oldstarttime, oldmonth, oldday)
				mutex.Unlock()
				if err != nil || dateCheckErr != nil || timeCheck != true {
					m["error"] = "Error occurred! Please double check your inputs!"
				} else if success == false {
					m["error"] = "Booking not found!"
				} else {
					m["success"] = "Booking deleted!"
				}
			}
		}
		err := tpl.ExecuteTemplate(res, "results.gohtml", m)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
