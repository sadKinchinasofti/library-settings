package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("method", r.Method)
	fmt.Println("URL", r.URL)
	switch method := r.Method; method {
	case "POST":
		var l Library
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			panic(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err := insertBook(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	case "GET":
		id, ok := r.URL.Query()["BookID"]
		fmt.Println("helo....")
		//name, nameok := r.URL.Query()["BookName"]
		fmt.Println(id)
		//fmt.Println(name)
		if !ok || len(id[0]) < 1 {
			if libdata := getBook(""); libdata == nil {
				json.NewEncoder(w).Encode("Error getting the data")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(libdata)
				w.WriteHeader(http.StatusOK)
			}
		} else {
			fmt.Printf("type of id %T\n", id[0])
			//fmt.Printf("type of name %T\n", name[0])
			if libdata := getBook(id[0]); libdata == nil {
				json.NewEncoder(w).Encode("Error getting the data")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(libdata)
				w.WriteHeader(http.StatusOK)
			}
		}
	case "DELETE":
		id, ok := r.URL.Query()["BookID"]
		//name, nameok := r.URL.Query()["BookName"]
		var l GetLibrary
		fmt.Println(id)
		if !ok || len(id[0]) < 1 {
			json.NewEncoder(w).Encode("please provide the id")
		} else {
			l = GetLibrary{
				BookID: id[0],
				//BookName: name[0],
			}
			if errdel := deleteBook(&l); errdel != nil {
				json.NewEncoder(w).Encode("error in delete operation")
			} else {
				json.NewEncoder(w).Encode("Delete successfully")
			}

		}
	case "PUT":
		//information sent in message body
		var l GetLibrary
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			panic(err)
		}
		BookID := l.BookID
		BookName := l.BookName
		BookAuthor := l.BookAuthor
		fmt.Println(BookID)
		fmt.Println(BookName)
		fmt.Println(BookAuthor)
		fmt.Printf("type of BookID %T\n", BookID)
		fmt.Printf("type of BookName %T\n", BookName)
		fmt.Printf("type of BookAuthor %T\n", BookAuthor)
		l1 := GetLibrary{
			BookID:     BookID,
			BookName:   BookName,
			BookAuthor: BookAuthor,
		}
		if errupd := updateBook(&l1); errupd != nil {
			json.NewEncoder(w).Encode("error in update operation")
		} else {
			json.NewEncoder(w).Encode("Update successfully")
		}

	}
}
