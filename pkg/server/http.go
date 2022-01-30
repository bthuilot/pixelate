package server

import (
	"log"
	"net/http"
)

func Init() {
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}