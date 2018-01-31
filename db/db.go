package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/compgen-io/cgpipemon/config"
	"github.com/compgen-io/cgpipemon/auth"
	"github.com/compgen-io/cgpipemon/model"

	_ "github.com/lib/pq"
)

func InitDb(cfg *config.Config) (*pgDb, error) {
	fmt.Println("InitDb")
	if dbConn, err := sql.Open("postgres", cfg.DBConnect); err != nil {
		fmt.Println("Connect failed: "+ cfg.DBConnect)
		return nil, err
	} else {
		fmt.Println("Good connection!")
		p := &pgDb{DbConn: dbConn}
		if err := p.DbConn.Ping(); err != nil {
			fmt.Println("No PING!")
			return nil, err
		}
		return p, nil
	}

}

type pgDb struct {
	DbConn *sql.DB
}

func (p *pgDb) CreateDB(cfg *config.Config) error {
	fmt.Println("initializing tables")
	initSql, err := cfg.Asset("db/schema.sql")
    if err != nil {
        log.Fatal("Missing schema.sql")
    }


	if _, err := p.DbConn.Query(string(initSql)); err != nil {
		log.Panic(err)
		return err
	}
	fmt.Println("Reloaded schema")

    pass := string(auth.GenerateRandom(16))
    fmt.Println("admin password: " + pass)

    _, err1 := model.NewUser(p.DbConn, "admin", pass, true)
    if err1 != nil {
		log.Panic(err1)
		return err1
	}    	
    	
	return nil
}
