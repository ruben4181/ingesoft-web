package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var database *sql.DB

//Definicion de las estructuras de datos

type User struct {
	ID_user   int
	Firstname string
	Lastname  string
	Username  string
	Email     string
	Password  string
}

type Program struct {
	ID_program          int
	Program_name        string
	Program_description string
	ID_user             int
}
type Post struct {
	ID_post       int
	Post_title    string
	Post_abstract string
	Post_body     string
	ID_user       int
	ID_program    int
}

type Event struct {
	ID_event           int
	Event_title        string
	Event_abstract     string
	Event_body         string
	Event_date_relased string
	Event_time_relased string
	Event_date         string
	Event_time         string
	ID_user            int
	ID_program         int
}

type InternalResponse struct {
	Status      int
	Description string
}

//Funciones que interactuan con la base de datos

const OK_STATUS = 0
const USER_PASS_NOT_MATCH = 1
const ERR_STATUS = -1
const EXCEPTION_STATUS = -2

/*
	EXAMPLE QUERY TO ADD NEW EVENT TO DATABASE

	INSERT INTO events(event_title, event_abstract, event_body, event_date_time_relased, event_date_time, id_user, id_program) VALUES('Primer evento Javeriano', 'El primer evento Javeriano que se realizara a finalizar este semestre contara con presentaciones artisticas y deportivas para todos los gustos', 'El contenido del evento', NOW(), '2018-11-26 07:00:00', 4, 1);
*/

func GetEvents(condition string) []Event {
	var events []Event
	data, err := database.Query("SELECT id_event, event_title, event_abstract, event_body, DATE(event_date_time_relased) as event_date_relased, " +
		"TIME(event_date_time_relased) as event_time_relased, DATE(event_date_time) as event_date, TIME(event_date_time) as event_time, " +
		"id_user, id_program from events WHERE " + condition)
	if err != nil {
		fmt.Println("An error ocurred during executing query in GetEvents function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_event, id_user, id_program int
			var event_title, event_abstract, event_body, event_date_relased, event_time_relased,
				event_date, event_time string
			err2 := data.Scan(&id_event, &event_title, &event_abstract, &event_body, &event_date_relased, &event_time_relased,
				&event_date, &event_time, &id_user, &id_program)
			if err2 != nil {
				fmt.Println("Error while scanning result from query in event")
				log.Fatal(err2)
			} else {
				events = append(events, Event{ID_event: id_event, Event_title: event_title, Event_abstract: event_abstract,
					Event_body: event_body, Event_date_relased: event_date_relased, Event_time_relased: event_time_relased,
					Event_date: event_date, Event_time: event_time, ID_user: id_user, ID_program: id_program})
			}
		}
	}
	return events
}

func GetPosts(condition string) []Post {
	var posts []Post
	data, err := database.Query("SELECT id_post, post_title, post_abstract, post_body, id_user, id_program FROM posts WHERE " + condition)
	if err != nil {
		fmt.Println("An error ocurred during executing query in GetPosts function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_post, id_user, id_program int
			var post_title, post_abstract, post_body string
			err2 := data.Scan(&id_post, &post_title, &post_abstract, &post_body, &id_user, &id_program)
			if err2 != nil {
				fmt.Println("Error while scanning result from query in post")
				log.Fatal(err)
			} else {
				posts = append(posts, Post{ID_post: id_post, Post_title: post_title, Post_abstract: post_abstract,
					Post_body: post_body, ID_user: id_user, ID_program: id_program})
			}
		}
	}
	return posts
}

func NewEvent(event Event){
	queryString:="INSERT INTO events(event_title, event_abstract, event_body, event_date_time_relased, event_date_time, id_user, id_program)"+ 
	"VALUES('"+event.Event_title+"', '"+event.Event_abstract+"', '"+event.Event_body+"', NOW(), '"+event.Event_date+" "+event.Event_time+"', "+
	strconv.Itoa(event.ID_user)+", "+strconv.Itoa(event.ID_program)")"
	_, err:=database.Query(queryString);
	if err!=nil{
		fmt.Println("An error ocurred while insert in events table");
		fmt.Println("The next query: \n"+queryString);
		log.Fatal(err);
	}else{
		fmt.Println("New Event has been added");
	}
}

func NewPost(post Post) {
	queryString := "INSERT INTO posts(post_title, post_abstract, post_body, id_user, id_program) VALUES('" +
		post.Post_title + "', '" + post.Post_abstract + "', '" + post.Post_body + "', " + strconv.Itoa(post.ID_user) + ", " + strconv.Itoa(post.ID_program) + ")"
	_, err := database.Query(queryString)
	if err != nil {
		fmt.Println("An error ocurred while insert in posts' table")
		fmt.Println("The next query: \n" + queryString)
		log.Fatal(err)
	} else {
		fmt.Println("New Post has been added")
	}
}

func OpenDB(user string, password string) *sql.DB {
	fmt.Println("Openning database 'ingesoft'")
	db, err := sql.Open("mysql", user+":"+password+"@tcp(localhost:3306)/ingesoft")
	if err != nil {
		fmt.Println("Ocurrio un error en la apertura de la base de datos")
	} else {
		fmt.Println("The database was opened correctly")
	}
	return db
}

func DelPost(condition string) {
	_, err := database.Query("DELETE from posts WHERE " + condition)
	if err != nil {
		fmt.Println("An error ocurred during executing Query in DelPost function")
		log.Fatal(err)
	}
}

func DelEvent(condition string){
	_, err:= database.Query("DELETE from events WHERE "+condition)
	if err != nil {
		fmt.Println("An error ocurred during executing Query in DelEvent function")
		log.Fatal(err)
	}
}

func GetUser(condition string) User {
	var user User

	data, err := database.Query("SELECT id_user, firstname, lastname, username, email, password FROM users WHERE " + condition)
	if err != nil {
		fmt.Println("An error ocurred during executing query in GetUser function")
		log.Fatal(err)
	} else {
		var id_user int
		var firstname, lastname, username, email, password string
		for data.Next() {
			err2 := data.Scan(&id_user, &firstname, &lastname, &username, &email, &password)
			if err2 != nil {
				fmt.Println("Error while scanning result from query in GetUser function")
				log.Fatal(err2)
			} else {
				user.ID_user = id_user
				user.Firstname = firstname
				user.Lastname = lastname
				user.Username = username
				user.Email = email
				user.Password = password
			}
		}
	}
	return user
}

func AuthenticateUser(username string, password string, action string) InternalResponse {
	user := GetUser("username='" + username + "' AND password='" + password + "'")
	if user.Username != "" {
		if user.Username == username {
			return InternalResponse{Status: OK_STATUS, Description: "Well done!"}
		} else {
			return InternalResponse{Status: EXCEPTION_STATUS, Description: "Something went wrong"}
		}
	}
	return InternalResponse{Status: USER_PASS_NOT_MATCH, Description: "Username and/or Password Don't match"}
}

func GetPrograms(condition string) []Program {
	var programs []Program

	data, err := database.Query("SELECT id_program, program_name, Program_description, id_user FROM programs WHERE " + condition)
	if err != nil {
		fmt.Println("An error ocurred during query in GetPrograms function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_program, id_user int
			var program_name string
			var program_description string
			err2 := data.Scan(&id_program, &program_name, &program_description, &id_user)
			if err2 != nil {
				fmt.Println("Error while scanning result from query in GetPrograms function")
				log.Fatal(err2)
			} else {
				programs = append(programs, Program{ID_program: id_program, Program_name: program_name, Program_description: program_description, ID_user: id_user})
			}
		}
	}
	return programs
}

func UpdatedEvent(event Event, condition string){
	queryString:="UPDATE events SET event_title='"+event.Event_title+"', event_abstract='"+event.Event_abstract+"', "+
	"event_body='"+event.Event_body+"', event_date_time='"+event.Event_date+" "+event.Event_time+"', id_user='"event.ID_user"' "+
	"WHERE "+condition
	_, err:=database.Query(queryString);
	if err!=nil{
		fmt.Println("Error while executing query in UpdateEvent")
		log.Fatal(err)
	}
}

func UpdatePost(post Post, condition string) {
	queryString := "UPDATE posts SET post_title='" + post.Post_title + "', post_abstract='" + post.Post_abstract + "', " +
		"post_body='" + post.Post_body + "', id_user=" + strconv.Itoa(post.ID_user) + " WHERE " + condition
	_, err := database.Query(query)
	if err != nil {
		fmt.Println("Error while executing query in UpdatePost")
		log.Fatal(err)
	}
}

//Funciones que interactuan con el Frontend

func GetProgramsEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Println("Sending programs to client")
	programs := GetPrograms("1")
	json.NewEncoder(w).Encode(programs)
}

func GetEventsEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	params := mux.Vars(req)
	events := GetEvents("id_program='" + params["id_program"] + "'")
	json.NewEncoder(w).Encode(events)
}

func GetPostsEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	params := mux.Vars(req)
	posts := GetPosts("id_program='" + params["id_program"] + "'")
	json.NewEncoder(w).Encode(posts)
}

func NewEventEP(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")
	
	var event Event
	_=json.NewDecoder(req.Body).Decode(&event)
	NewEvent(event)
	json.NewEncoder(w).Encode(event)
}

func NewPostEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")
	fmt.Println("Creating new posts")

	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)
	fmt.Println(post)
	NewPost(post)
	json.NewEncoder(w).Encode(post)
}

func UpdatePostEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)
	UpdatePost(post, "id_post="+strconv.Itoa(post.ID_post))
	json.NewEncoder(w).Encode(post)
}

func DelEventEP(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")
	
	var event Event
	_=json.NewDecoder(req.Body).Decode(&event)
	DelEvent("id_post="+strconv.Itoa(event.ID_post))
	json.NewEncoder(w).Encode(event)
	fmt.Println("Event deleted")
}

func DelPostEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")
	fmt.Println("Deleting Post")

	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)
	DelPost("id_post=" + strconv.Itoa(post.ID_post))
	json.NewEncoder(w).Encode(post)
	fmt.Println("Post deleted")
}


func GetPostEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	params := mux.Vars(req)
	var post Post
	posts := GetPosts("id_post=" + params["id_post"])
	if len(posts) > 0 {
		post = posts[0]
	}

	json.NewEncoder(w).Encode(post)
}

//Main

func main() {
	database = OpenDB("root", "Californication16")
	user := GetUser("username='ruben4181'")
	fmt.Println(user)
	//NewPost(Post{ID_post:0, Post_title:"Titulo Post 4", Post_abstract:"-Abstract Post 4", Post_body:"Body Post en algun formato", ID_user:3, ID_program:1});
	fmt.Println(AuthenticateUser("ruben4181", "Dadada", "None"))

	//Todo lo concerniente a http
	router := mux.NewRouter()
	router.HandleFunc("/getPrograms", GetProgramsEP).Methods("GET")
	router.HandleFunc("/getPosts/{id_program}", GetPostsEP).Methods("GET")
	router.HandleFunc("/newPost", NewPostEP).Methods("POST")
	router.HandleFunc("/delPost", DelPostEP).Methods("POST")
	router.HandleFunc("/getPost/{id_post}", GetPostEP).Methods("GET")
	router.HandleFunc("/updatePost", UpdatePostEP).Methods("POST")
	router.HandleFunc("/getEvents/{id_program}", GetEventsEP).Methods("GET")

	http.ListenAndServe(":8080", router)
}
