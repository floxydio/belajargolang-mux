package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Product struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@tcp(localhost)/belajargo?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection Connect", err)
	}

	db.AutoMigrate(&Product{})

	handleRequest()

}

func handleRequest() {
	log.Println("Start at 8000")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/getproduct", getProduct)
	myRouter.HandleFunc("/api/getproduct/{id}", returnSingleProduct)
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to homepage")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(reqBody, &product)
	db.Create(&product)

	fmt.Println("Endpoint hit: Creating")
	json.NewEncoder(w).Encode(product)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	myProduct := []Product{}

	db.Find(&myProduct)
	fmt.Println("Endpoint hit: returnAllBookings")
	json.NewEncoder(w).Encode(myProduct)
}

func returnSingleProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	myProduct := []Product{}
	db.Find(&myProduct)

	for _, Product := range myProduct {
		s, err := strconv.Atoi(key)
		if err == nil {
			if Product.Id == s {
				fmt.Println(Product)
				fmt.Println("Endpoint hit:", key)
				json.NewEncoder(w).Encode(Product)
			}
		}

	}
}
