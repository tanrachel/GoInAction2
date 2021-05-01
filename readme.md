# Venue Booking System 

Go In Action 2 

## what it contains 

- main.go: where the server is located along with handler functions.
- initializer.go: seeding the application with test data.
- functionHelper.go: all other functions used in the handler.
- userDB.go: functions to implement a mysql db for the user database.
- vApp: module for holding all the different data structures to hold the venue booking system.
- templates: html templates.

## running the application 

go run *.go  

## query used for the mysql database 
 CREATE database VenueDB;
 USE VenueDB;
 CREATE TABLE Users (UserName VARCHAR(30) NOT NULL PRIMARY KEY, Password VARCHAR(256));