package main

import(
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

var database *sql.DB;

type User struct{
	ID_user int;
	Firstname string;
	Lastname string;
	Username string;
	Email string;
	Password string;
}

func OpenDB(user string, password string) *sql.DB{
	fmt.Println("Openning database 'ingesoft'");
	db, err:=sql.Open("mysql", user+":"+password+"@tcp(localhost:3306)/ingesoft");
	if err!=nil{
		fmt.Println("Ocurrio un error en la apertura de la base de datos");
	}else{
		fmt.Println("The database was opened correctly");
	}
	return db;
}

func GetUser(values string, condition string) User{
	var user User;
	
	if (values==""){
		data, err:=database.Query("SELECT id_user, firstname, lastname, username, email, password FROM users WHERE "+condition);
		if err!=nil{
			fmt.Println("An error ocurred during the query in GetUser function");
			log.Fatal(err);
		}else{
			var id_user int;
			var firstname, lastname, username, email, password string;
			for data.Next(){
				err2:=data.Scan(&id_user, &firstname, &lastname, &username, &email, &password);
				if(err2!=nil){
					fmt.Println("Error while scanning result from query in GetUser function");
					log.Fatal(err2);
				}else{
					user.ID_user=id_user;
					user.Firstname=firstname;
					user.Lastname=lastname;
					user.Username=username;
					user.Email=email;
					user.Password=password;
				}
			}
		}
	}else{

	}
	return user;
}

func main() {
	database=OpenDB("root", "Californication16");
	user:=GetUser("", "username='ruben4181'");
	fmt.Println(user);
}