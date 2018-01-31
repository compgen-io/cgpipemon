package apiv1

import (
    "fmt"
    "net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK") // send data to client side
}


func Auth(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK") // send data to client side
}
