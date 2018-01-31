package main

import (
    "os"
    "os/signal"
    "syscall"
    // "bufio"
    // "flag"
    "fmt"
    "net/http"
    "strings"
    "log"
    // "golang.org/x/crypto/ssh/terminal"
    //"github.com/BurntSushi/toml"


    "github.com/compgen-io/cgpipemon/config"
    "github.com/compgen-io/cgpipemon/db"
    // "github.com/compgen-io/cgpipemon/model"
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





func startServer(address string, dbConn string) {
    pgdb, err1 := db.InitDb(dbConn)
    if err1 != nil {
        log.Fatal("DB Init: ", err1)
    }

    fmt.Println("Listening on: "+address)
    srv := &http.Server{Addr: address}

    http.HandleFunc("/", sayhelloName) // set router
    http.HandleFunc("/api/v1/ping", apiv1.Ping) // set router
    http.HandleFunc("/api/v1/auth", apiv1.Auth) // set router

    // capture ^C
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("")
        fmt.Println("Ctrl-C... shutting down server")
        srv.Shutdown(nil)
    }()

    err2 := srv.ListenAndServe() // set listen port
    if err2 != nil {
        log.Printf("Httpserver: ListenAndServe() error: %s", err2)
    }

    // cleanup
    pgdb.Close()

}

func createDB(dbConn string) {
    pgdb, err1 := db.InitDb(dbConn)
    if err1 != nil {
        log.Fatal("DB Init: ", err1)
    }

    pgdb.CreateDB()
    pgdb.Close()
}

func main() {
    cfg, err := config.Init()

    if err != nil {
        config.ServerUsage()
        log.Fatal(err)
    }

    switch cfg.Command {
    case "serve":
        startServer(cfg.Listen, cfg.DBConnect)
    case "createdb":
        createDB(cfg.DBConnect)
    default:
        fmt.Println("Unknown command: "+cfg.Command)
    }
}