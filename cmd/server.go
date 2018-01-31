package main

import (
    "os"
    "syscall"
    "bufio"
    "fmt"
    "net/http"
    "strings"
    "log"
    "golang.org/x/crypto/ssh/terminal"

    "github.com/compgen-io/cgpipemon/config"
    "github.com/compgen-io/cgpipemon/db"
    "github.com/compgen-io/cgpipemon/model"
    "github.com/compgen-io/cgpipemon/apiv1"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  // parse arguments, you have to call this by yourself
    fmt.Println(r.Form)  // print form information in server side
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello!") // send data to client side
}

func main() {

    cfg := config.Init()

    fmt.Println(cfg.Listen)
    fmt.Println(cfg.DBConnect)
    fmt.Println(cfg.Command)

    pgdb, err1 := db.InitDb(cfg)
    if err1 != nil {
        log.Fatal("DB Init: ", err1)
    }

    if cfg.Command == "createdb" {
        pgdb.CreateDB(cfg)
    } else if cfg.Command == "testlogin" {
        reader := bufio.NewReader(os.Stdin)

        fmt.Print("Enter Username: ")
        username, _ := reader.ReadString('\n')

        fmt.Print("Enter Password: ")
        passwd, err := terminal.ReadPassword(int(syscall.Stdin))
        fmt.Println("")
        if err == nil && model.CheckPass(pgdb.DbConn, strings.TrimSpace(username), strings.TrimSpace(string(passwd))) {
            fmt.Println("Success!")
        } else {
            fmt.Println("Error")
        }

    } else {
        
        http.HandleFunc("/", sayhelloName) // set router
        http.HandleFunc("/api/v1/ping", apiv1.Ping) // set router
        http.HandleFunc("/api/v1/auth", apiv1.Auth) // set router

        err := http.ListenAndServe(cfg.Listen, nil) // set listen port
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    }
}
