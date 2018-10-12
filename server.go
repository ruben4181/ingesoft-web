package main

import(
		"fmt"
		"net/http"
		"github.com/gorilla/mux"
		"encoding/json"
)

type Test struct{
	Nombre string;
	Apellido string;
	Edad int;
	Numeros []string;
}


func TestFunc(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Println("Esta en TestFunc!!!")
	t:=Test{Nombre:"Ruben", Apellido:"Vargas", Edad:20, Numeros:[]string{"3116021602", "3146880001", "3233041731"}}
	json.NewEncoder(w).Encode(t)
}

func main(){
		router:=mux.NewRouter();

		router.HandleFunc("/test", TestFunc).Methods("GET");
		http.ListenAndServe(":8080", router);
}