package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type Product struct {
	Id     uint   `form:"id" json:"id" gorm:"primaryKey" `
	Title  string `form:"title" json:"title"`
	Price  int    `form:"harga" json:"harga"`
	Author string `form:"author" json:"author"`
}

type CustomValidator struct {
	validator *validator.Validate
}

type Result struct {
	Code    int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var validate *validator.Validate

func main() {

	db, err = gorm.Open("mysql", "root:@tcp(localhost)/belajargo")
	validate = validator.New()
	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection Connect")
	}

	db.AutoMigrate(&Product{})

	handleRequest()

}

func handleRequest() {

	log.Println("Start at 8000")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/products", products).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func products(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)
	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Berhasil Menambahkan Produk"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Products Get"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
