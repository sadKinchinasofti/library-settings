package handlers

import (
	"log"
)

func insertBook(l *Library) error {
	e := NewLibRepo()
	if errInsertion := e.InsertBook(l); errInsertion != nil {
		log.Panic("error in insertion", errInsertion)
		return errInsertion
	}
	return nil
}

func getBook(id string) []GetLibrary {
	getb := NewLibRepo()
	if len(id) == 0 {
		lib := getb.GetBook()
		return lib
	}
	libone := getb.GetBookByID(id)
	return libone

}

func deleteBook(l *GetLibrary) error {
	delb := NewLibRepo()
	if errdeletion := delb.DeleteBook(l); errdeletion != nil {
		return errdeletion
	}
	return nil
}

func updateBook(l *GetLibrary) error {
	updateb := NewLibRepo()
	if errUpdation := updateb.UpdateBook(l); errUpdation != nil {
		log.Panic("error in updation", errUpdation)
		return errUpdation
	}
	return nil
}
