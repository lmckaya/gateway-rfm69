package main
import (
    "log"
    "net/http"
)

func GetThings(w http.ResponseWriter, r *http.Request) {
   log.Println("GetThings: ",r);
}
func GetThing(w http.ResponseWriter, r *http.Request) {
   log.Println("GetThing: ",r);
}
