package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type devInfo struct {
	DevKey string
	Value  string
}

// FileName determines what file to check for
const fileName = "Alaris.ALA"

// Port is assigned for client connection
const port = ":9159"

func createDevFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit- Create")

	var ret = FileExists(fileName)

	if !ret {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't ready body", http.StatusBadRequest)
			return
		}

		s := string(body)

		log.Println(s)

		data := devInfo{}

		json.Unmarshal([]byte(s), &data)

		log.Println(data.DevKey)

		// Saving To File
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		l, err := f.WriteString(data.DevKey)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}

		fmt.Println(l, "bytes written successfully")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Error: File Exists")
		fmt.Println("Delete manually on server to save new file")
	}
}

func retrieveDevFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit- Load")

	var ret = FileExists(fileName)

	if ret {
		data, err := ioutil.ReadFile(fileName)

		if err != nil {
			fmt.Println(err)
		}

		s := string(data)

		fmt.Print(s)

		json.NewEncoder(w).Encode(s)
	} else {
		fmt.Println("File does not exist")
	}
}

func checkDevFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit- Check Dev File")

	var ret = FileExists(fileName)

	var s string

	if ret {
		s = "true"
	} else {
		s = "false"
	}

	fmt.Println(s)

	json.NewEncoder(w).Encode(s)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/FileCheck", checkDevFile).Methods("GET")
	myRouter.HandleFunc("/", createDevFile).Methods("POST")
	myRouter.HandleFunc("/", retrieveDevFile).Methods("GET")
	log.Fatal(http.ListenAndServe(port, myRouter))
}

// FileExists checks if the File Exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}

