package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Product struct {
	Id    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	Price int    `json:"harga"`
}

type User struct {
	Id       uint   `json:"id" gorm: "primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
	Level    int    `json:"level"`
}

type Result struct {
	Code    int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@tcp(localhost)/belajargo")
	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection Connect", err)
	}

	db.AutoMigrate(&Product{})
	db.AutoMigrate(&User{})

	handleRequest()

}

func handleRequest() {
	log.Println("Start at 8000")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/products", products).Methods("POST")
	myRouter.HandleFunc("/api/register", register).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func products(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(product)
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var user User
	json.Unmarshal(reqBody, &user)

	db.Create(&user)

	res := Result{Code: 200, Data: user, Message: "Berhasil Registrasi"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
