package main

//API to perform CRUD operations on the User data present in the User table of Mysql database 'MyApp'

//get necessary database driver for MySQL and gorilla/mux library for routing purposes
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// interface to connect to a database
type DatabaseConnection interface {
	connect() (*sql.DB, error)
}

// type struct to implement connect method for MySQL
type MySQLConnection struct {
	ConnectionString string
}

// Implementing Connect method for MySQL
func (m *MySQLConnection) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", m.ConnectionString)

	if err != nil {
		return nil, err
	}
	return db, nil
}

// function to connect to the database based on the type of database connection passed to it
func UseDatabase(dbConn DatabaseConnection) *sql.DB {
	db, err := dbConn.connect()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil
	}

	fmt.Println("Connected to the database successfully!")
	return db
}

// User type structure to create User objects
type User struct {
	UserId   string
	Name     string
	Password string
	isActive []uint8
}

// dummy route to create the website homepage
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

// dummy route to test the route on postman app
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

// Handler function to get the user based on UserID passed in the route
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	mysqlConn := &MySQLConnection{
		ConnectionString: "root:Hrithik@/MyApp",
	}
	db := UseDatabase(mysqlConn)
	defer db.Close()

	vars := mux.Vars(r)
	userId := vars["UserId"]
	fmt.Println(userId)

	// Call the GetUser function to fetch the user data from the database
	user, err := GetUser(db, userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Convert the user object to JSON and send it in the response
	//Set the header option to content type json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// function to implement the get user method from the database
func GetUser(db *sql.DB, id string) (*User, error) {
	query := "SELECT * FROM user WHERE UserId = ?"
	row := db.QueryRow(query, id)
	fmt.Println(row)
	user := &User{}
	err := row.Scan(&(user.UserId), &(user.Name), (&user.Password), &(user.isActive))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, nil
}

// Handler function to update the user based on the given UserID in the route
func updateUserHandler(w http.ResponseWriter, r *http.Request) {

	mysqlConn := &MySQLConnection{
		ConnectionString: "root:Hrithik@/MyApp",
	}
	db := UseDatabase(mysqlConn)
	defer db.Close()

	vars := mux.Vars(r)
	userId := vars["UserId"]
	fmt.Println(userId)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	// Call the UpdateUser function to update the data of a particular UserId
	UpdateUser(db, userId, user.Name, user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	io.WriteString(w, "User is updated successfully")
	fmt.Fprintln(w, "User updated successfully")
}

// method to update the user from the database
func UpdateUser(db *sql.DB, UserId, Name, Password string) error {
	query := "UPDATE user SET Name = ?, Password= ? WHERE Userid = ?"
	_, err := db.Exec(query, Name, Password, UserId)
	if err != nil {
		return err
	}
	return nil
}

// method to delete the user from the database
func DeleteUser(db *sql.DB, name string) error {
	query := "Delete from user where name = ?"
	_, err := db.Exec(query, name)
	fmt.Println(err)
	if err != nil {
		return err
	}

	return nil
}

// handler functino to delete the user based on name given in the route
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	mysqlConn := &MySQLConnection{
		ConnectionString: "root:Hrithik@/MyApp",
	}
	db := UseDatabase(mysqlConn)
	defer db.Close()
	vars := mux.Vars(r)
	nameStr := vars["Name"]
	prob := DeleteUser(db, nameStr)
	if prob != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	fmt.Println("user deleted successfully")
	io.WriteString(w, "User is deleted successfully\n")
}

// handler function to create the user taking details from request body
func createUserHandler(w http.ResponseWriter, r *http.Request) {

	mysqlConn := &MySQLConnection{
		ConnectionString: "root:Hrithik@/MyApp",
	}
	db := UseDatabase(mysqlConn)
	defer db.Close()

	// Parse JSON data from the request body
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	err := CreateUser(db, user.Name, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
}

// function to create a new user in the database
func CreateUser(db *sql.DB, name, password string) error {
	query := "INSERT INTO user (Name, Password) VALUES (?, ?)"
	_, err := db.Exec(query, name, password)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Go MySQL Tutorial")
	r := mux.NewRouter()

	//list of routes to handle CRUD operations
	r.HandleFunc("/", getRoot)
	r.HandleFunc("/hello", getHello)
	r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/delete/{Name}", deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/user/{UserId}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{UserId}", getUserHandler).Methods("GET")
	log.Println("Server lisening on :3300")

	//Enable the server on port 3300 using HTTP
	log.Fatal(http.ListenAndServe(":3300", r))
}
