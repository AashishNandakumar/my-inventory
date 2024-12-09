package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {

	err := a.Initialize("localhost", "5432", "admin", "admin", "inventory_test")
	if err != nil {
		log.Fatal("Could not initialize app:", err)
	}
	createTable()

	m.Run() // Run the tests
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INT,
		price NUMERIC(10,2)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE FROM products")
	_, err = a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1") // Reset the sequence
	if err != nil {
		log.Fatal("Could not clear table:", err)
	}
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal("Could not add product:", err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 5000.00)
	request, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code %v but got %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name": "keyboard", "quantity": 100, "price": 5000.00}`)
	request, err := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "keyboard" {
		t.Errorf("Expected product name to be 'keyboard'. Got '%v'", m["name"])
	}
	if m["quantity"] != 100.0 {
		t.Errorf("Expected product quantity to be 100. Got %v", m["quantity"])
	}
	if m["price"] != 5000.00 {
		t.Errorf("Expected product price to be 5000.00. Got %v", m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()

	addProduct("keyboard", 100, 5000.00)

	req, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatal("Error in creating request: ", err)
	}
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, err = http.NewRequest("DELETE", "/products/1", nil)
	if err != nil {
		t.Fatal("Error in creating request: ", err)
	}
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, err = http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatal("Error in creating request: ", err)
	}
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()

	addProduct("keyboard", 100, 5000.00)

	req, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatal("Error in creating request: ", err)
	}
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	payload := []byte(`{"name": "keyboard", "quantity": 200, "price": 6000.00}`)
	req, err = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal("Error in creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id to be %v. Got %v", oldValue["id"], newValue["id"])
	}
	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected name to be %v. Got %v", oldValue["name"], newValue["name"])
	}
	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("Expected quantity to be %v. Got %v", oldValue["quantity"], newValue["quantity"])
	}
	if oldValue["price"] == newValue["price"] {
		t.Errorf("Expected price to be %v. Got %v", oldValue["price"], newValue["price"])
	}
}
