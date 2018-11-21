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

type Teacher struct {
	Username           string
	Email              string
	Firstname          string
	Lastname           string
	ID_teacher         int
	ID_user            int
	ID_department      int
	ID_program         int
	Degrees            []Degree
	Achievements       []Achievement
	Teacher_department string
}
type Achievement struct {
	ID_achievement          int
	ID_teacher              int
	Achievement_name        string
	Achievement_description string
	Achievement_year        string
}
type Degree struct {
	ID_degree         int
	ID_teacher        int
	Degree_name       string
	Degree_college    string
	Degree_city       string
	Degree_year       string
	Degree_extra_info string
}
type Department struct {
	ID_department   int
	Department_name string
}

type Course struct {
	ID_course           int
	Course_name         string
	Course_description  string
	Course_n_credits    int
	Course_requirements string
	ID_program          int
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

func GetCourses(condition string) []Course {
	var courses []Course
	queryString := "SELECT id_course, course_name, course_description, course_n_credits, course_requirements, id_program FROM courses WHERE " + condition
	data, err := database.Query(queryString)
	if err != nil {
		fmt.Println("Error during executing query in GetCourses function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_course, id_program, course_n_credits int
			var course_name, course_description, course_requirements string
			err2 := data.Scan(&id_course, &course_name, &course_description, &course_n_credits, &course_requirements, &id_program)
			if err2 != nil {
				fmt.Println("Error while scanning data from courses")
				log.Fatal(err2)
			} else {
				courses = append(courses, Course{ID_course: id_course, Course_name: course_name, Course_description: course_description,
					Course_n_credits: course_n_credits, Course_requirements: course_requirements, ID_program: id_program})
			}
		}
	}
	return courses
}

func GetEvents(condition string) []Event {
	var events []Event
	data, err := database.Query("SELECT id_event, event_title, event_abstract, event_body, DATE(event_date_time_relased) as event_date_relased, " +
		"TIME(event_date_time_relased) as event_time_relased, DATE(event_date_time) as event_date, TIME(event_date_time) as event_time, " +
		"id_user, id_program from events WHERE " + condition + " ORDER BY event_date_time")
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

func GetTeachers(condition string) []Teacher {
	var teachers []Teacher
	values := "id_teacher, id_user, id_program, id_department, firstname, lastname, username, email"
	queryString := "SELECT " + values + " FROM teachers NATURAL JOIN users WHERE " + condition
	data, err := database.Query(queryString)
	if err != nil {
		fmt.Println("An error ocurred executing query in GetTeachers function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_teacher, id_user, id_program, id_department int
			var firstname, lastname, username, email string
			err2 := data.Scan(&id_teacher, &id_user, &id_program, &id_department, &firstname, &lastname, &username, &email)
			var mtDegree []Degree
			var mtAchievements []Achievement
			if err2 != nil {
				fmt.Println("Error while scanning result from query in teachers")
				log.Fatal(err2)
			} else {
				teachers = append(teachers, Teacher{Username: username, Firstname: firstname, Lastname: lastname,
					ID_teacher: id_teacher, ID_user: id_user, ID_department: id_department, ID_program: id_program, Degrees: mtDegree, Achievements: mtAchievements})
			}
		}
		for i := 0; i < len(teachers); i++ {
			teachers[i].Teacher_department = GetDepartment(teachers[i].ID_department)
			teachers[i].Degrees = GetDegrees(teachers[i])
			teachers[i].Achievements = GetAchivements(teachers[i])
		}
	}
	return teachers
}

func GetDepartment(ID_department int) string {
	var value string
	data, err := database.Query("SELECT department_name from departments WHERE id_department=" + strconv.Itoa(ID_department))
	if err != nil {
		log.Fatal(err)
	} else {
		for data.Next() {
			err2 := data.Scan(&value)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
	return value
}

func GetDegrees(teacher Teacher) []Degree {
	queryString := "SELECT  id_degree, degree_name, degree_college, degree_city, degree_year FROM degrees WHERE id_teacher=" + strconv.Itoa(teacher.ID_teacher)
	data, err := database.Query(queryString)
	if err != nil {
		fmt.Println("Error during executing query in GetDegrees function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_degree int
			var degree_name, degree_college, degree_city, degree_year string
			err2 := data.Scan(&id_degree, &degree_name, &degree_college, &degree_city, &degree_year)
			if err2 != nil {
				fmt.Println("Error while scanning data from query in degrees")
				log.Fatal(err2)
			} else {
				teacher.Degrees = append(teacher.Degrees, Degree{ID_degree: id_degree, Degree_name: degree_name,
					Degree_college: degree_college, Degree_city: degree_city, Degree_year: degree_year})
			}
		}
	}
	return teacher.Degrees
}

func GetAchivements(teacher Teacher) []Achievement {
	queryString := "SELECT id_achievement, achievement_name, achievement_description, achievement_year FROM achievements WHERE id_teacher=" + strconv.Itoa(teacher.ID_teacher)
	data, err := database.Query(queryString)
	if err != nil {
		fmt.Println("Error during executing query in GetAchievements function")
		log.Fatal(err)
	} else {
		for data.Next() {
			var id_achievement int
			var achievement_name, achievement_description, achievement_year string
			err2 := data.Scan(&id_achievement, &achievement_name, &achievement_description, &achievement_year)
			if err2 != nil {
				fmt.Println("Error while scanning data from achievements")
				log.Fatal(err2)
			} else {
				teacher.Achievements = append(teacher.Achievements, Achievement{ID_achievement: id_achievement,
					ID_teacher: 1, Achievement_name: achievement_name, Achievement_description: achievement_description,
					Achievement_year: achievement_year})
			}
		}
	}
	return teacher.Achievements
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

func NewCourse(course Course) {
	queryString := "insert into courses(course_name, course_description, course_n_credits, course_requirements, id_program) VALUES('" +
		course.Course_name + "', '" + course.Course_description + "', " + strconv.Itoa(course.Course_n_credits) + ", '" +
		course.Course_requirements + "', " + strconv.Itoa(course.ID_program) + ")"
	_, err := database.Query(queryString)
	if err != nil {
		fmt.Println("Error during executing query in NewCourse function")
		log.Fatal(err)
	}
}

func NewTeacher(teacher Teacher) {
	var next_index int
	q, q_err := database.Query("SELECT 'auto_increment' FROM INFORMATION_SCHEMA.TABLES WHERE table_name='teachers'")
	if q_err != nil {
		log.Fatal(q_err)
	} else {
		for q.Next() {
			q_err2 := q.Scan(&next_index)
			if q_err2 != nil {
				log.Fatal(q_err2)
			}
		}
	}
	queryString := "INSERT INTO teachers(id_user, id_department, id_program) VALUES(" +
		strconv.Itoa(teacher.ID_user) + ", " + strconv.Itoa(teacher.ID_department) + ", " + strconv.Itoa(teacher.ID_program) + ")"
	_, err := database.Query(queryString)
	if err != nil {
		fmt.Println("An error ocurred while insert in events table")
		fmt.Println("The next query: \n" + queryString)
		log.Fatal(err)
	} else {
		deegres_lenght := len(teacher.Degrees)
		for i := 0; i < deegres_lenght; i++ {
			tmpString := "insert into degrees(id_teacher, degree_name, degree_college, degree_city, degree_year, degree_extra_info) VALUES(" +
				strconv.Itoa(next_index) + ", '" + teacher.Degrees[i].Degree_name + "', '" + teacher.Degrees[i].Degree_college +
				"', '" + teacher.Degrees[i].Degree_city + "', '" + teacher.Degrees[i].Degree_year + "', '" + teacher.Degrees[i].Degree_extra_info + "')"
			_, err2 := database.Query(tmpString)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
		achievements_lenght := len(teacher.Achievements)
		for i := 0; i < achievements_lenght; i++ {
			tmpString := "insert into achievements(id_teacher, achievement_name, achievement_description, achievement_year) VALUES(" +
				strconv.Itoa(next_index) + ", '" + teacher.Achievements[i].Achievement_name + "', '" +
				teacher.Achievements[i].Achievement_description + "', '" + teacher.Achievements[i].Achievement_year + "')"
			_, err2 := database.Query(tmpString)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
}

func NewEvent(event Event) {
	queryString := "INSERT INTO events(event_title, event_abstract, event_body, event_date_time_relased, event_date_time, id_user, id_program)" +
		"VALUES('" + event.Event_title + "', '" + event.Event_abstract + "', '" + event.Event_body + "', NOW(), '" + event.Event_date + " " + event.Event_time + "', " +
		strconv.Itoa(event.ID_user) + ", " + strconv.Itoa(event.ID_program) + ")"
	_, err := database.Query(queryString)
	if err != nil {
		fmt.Println("An error ocurred while insert in events table")
		fmt.Println("The next query: \n" + queryString)
		log.Fatal(err)
	} else {
		fmt.Println("New Event has been added")
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

func DelEvent(condition string) {
	_, err := database.Query("DELETE from events WHERE " + condition)
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

func UpdateEvent(event Event, condition string) {
	queryString := "UPDATE events SET event_title='" + event.Event_title + "', event_abstract='" + event.Event_abstract + "', " +
		"event_body='" + event.Event_body + "', event_date_time='" + event.Event_date + " " + event.Event_time + "', id_user='" + strconv.Itoa(event.ID_user) + "' " +
		"WHERE " + condition
	_, err := database.Query(queryString)
	if err != nil {
		fmt.Println("Error while executing query in UpdateEvent")
		log.Fatal(err)
	}
}

func UpdatePost(post Post, condition string) {
	queryString := "UPDATE posts SET post_title='" + post.Post_title + "', post_abstract='" + post.Post_abstract + "', " +
		"post_body='" + post.Post_body + "', id_user=" + strconv.Itoa(post.ID_user) + " WHERE " + condition
	_, err := database.Query(queryString)
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

func GetTeachersEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	params := mux.Vars(req)
	teachers := GetTeachers("id_program='" + params["id_program"] + "'")
	json.NewEncoder(w).Encode(teachers)
}

func GetCoursesEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	params := mux.Vars(req)
	courses := GetCourses("id_program='" + params["id_program"] + "'")
	json.NewEncoder(w).Encode(courses)
}

func NewCourseEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var course Course
	_ = json.NewDecoder(req.Body).Decode(&course)
	NewCourse(course)
	json.NewEncoder(w).Encode(course)
}

func NewTeacherEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var teacher Teacher
	_ = json.NewDecoder(req.Body).Decode(&teacher)
	NewTeacher(teacher)
	json.NewEncoder(w).Encode(teacher)
}

func NewEventEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var event Event
	_ = json.NewDecoder(req.Body).Decode(&event)
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

func UpdateEventEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var event Event
	_ = json.NewDecoder(req.Body).Decode(&event)
	UpdateEvent(event, "id_event="+strconv.Itoa(event.ID_event))
	json.NewEncoder(w).Encode(event)
}

func DelEventEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	var event Event
	_ = json.NewDecoder(req.Body).Decode(&event)
	DelEvent("id_event=" + strconv.Itoa(event.ID_event))
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

func GetEventEP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS,*")

	params := mux.Vars(req)
	var event Event
	events := GetEvents("id_event=" + params["id_event"])
	if len(events) > 0 {
		event = events[0]
	}

	json.NewEncoder(w).Encode(event)
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
	router.HandleFunc("/newEvent", NewEventEP).Methods("POST")
	router.HandleFunc("/updateEvent", UpdateEventEP).Methods("POST")
	router.HandleFunc("/delEvent", DelEventEP).Methods("POST")
	router.HandleFunc("/getEvent/{id_event}", GetEventEP).Methods("GET")
	router.HandleFunc("/getTeachers/{id_program}", GetTeachersEP).Methods("GET")
	router.HandleFunc("/getCourses/{id_program}", GetCoursesEP).Methods("GET")
	router.HandleFunc("/newCourse", NewCourseEP).Methods("POST")
	router.HandleFunc("/newTeacher", NewTeacherEP).Methods("POST")
	http.ListenAndServe(":8080", router)
}

/*

 */
