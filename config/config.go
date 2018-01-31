package config

import (
    "os"
    "flag"
    "fmt"
    "log"
    "errors"
    "path/filepath"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Listen string
    DBConnect string
    Command string
}

type dbConfig struct {
    Host string
    User string
    Password string
    Name string
    Port int
    SSLmode string    
    URL string    
}

type serverConfig struct {
    Host string
    Port int
}

type tmpConfig struct {
    Database dbConfig
    Server serverConfig
}

func (c *tmpConfig) merge(cfg *tmpConfig) {
    fmt.Println("merging configs...")
    fmt.Printf("%+v\n", c)
    fmt.Printf("%+v\n", cfg)
}

func (c *tmpConfig) loadToml() error {
    if _, err := os.Stat(".cgpipemonrc"); err == nil {
        if _, err1 := toml.DecodeFile(".cgpipemonrc", c); err1 != nil {
            log.Print(err1)
        }
    } else {
        ex, err1 := os.Executable()
        if err1 != nil {
            log.Panic(err1)
        }
        exPath := filepath.Dir(ex)
        if _, err2 := toml.DecodeFile(exPath+string(os.PathSeparator)+".cgpipemonrc", c); err2 != nil {
            log.Print(err2)
        }

    }
    return nil    
}


func ServerUsage() {
    fmt.Println("Usage: cgpipemon-server cmd {opts}")
    fmt.Println("")
    fmt.Println("Valid commands:")
    fmt.Println("    serve     Start the cgpipemon server")
    fmt.Println("    createdb  Initialize the database")
    fmt.Println("")
}

func Init() (*Config, error) {
    if len(os.Args) < 2 {
        return nil, errors.New("Missing command")
    }

    tmpCfg := tmpConfig { 
        Server: serverConfig { 
            Host: "127.0.0.1",
            Port: 3000 } }

    if err := tmpCfg.loadToml() ; err != nil {
        return nil, err
    }

    if err := tmpCfg.loadCmdArgs() ; err != nil {
        return nil, err
    }

    var connStr string
    if tmpCfg.Database.URL == "" {
        connStr = buildDBConnString(tmpCfg.Database.Host, tmpCfg.Database.Port, tmpCfg.Database.Name, tmpCfg.Database.User, tmpCfg.Database.Password, tmpCfg.Database.SSLmode)
    } else {
        connStr = tmpCfg.Database.URL
    }

    address := fmt.Sprintf("%s:%d", tmpCfg.Server.Host, tmpCfg.Server.Port)

    cfg := &Config{
        Command: os.Args[1],
        Listen: address,
        DBConnect: connStr}

    cfg.Command = os.Args[1]
    cfg.Listen = address
    cfg.DBConnect = connStr

    return cfg, nil

}


func addServerFlags(flg *flag.FlagSet, cfg *tmpConfig) {
    flg.StringVar(&cfg.Server.Host, "listen", cfg.Server.Host, "Server host/IP")
    flg.IntVar(&cfg.Server.Port, "port", cfg.Server.Port, "Server port")
}
func addDBFlags(flg *flag.FlagSet, cfg *tmpConfig) {
    flg.StringVar(&cfg.Database.Host, "db-host", cfg.Database.Host, "Database hostname")
    flg.IntVar(&cfg.Database.Port, "db-port", cfg.Database.Port, "Database port")
    flg.StringVar(&cfg.Database.Name, "db-name", cfg.Database.Name, "Database name")
    flg.StringVar(&cfg.Database.User, "db-user", cfg.Database.User, "Database username")
    flg.StringVar(&cfg.Database.Password, "db-pass", cfg.Database.Password, "Database password")
    flg.StringVar(&cfg.Database.SSLmode, "db-sslmode",cfg.Database.SSLmode, "Database SSL mode (valid: disable, require, verify-ca, verify-full)")
    flg.StringVar(&cfg.Database.URL, "db-url", cfg.Database.URL, "Database connection string / URL")
}

func (c *tmpConfig) loadCmdArgs() error {
    serveFlags := flag.NewFlagSet("serve", flag.ExitOnError)
    createDBFlags := flag.NewFlagSet("createdb", flag.ExitOnError)

    switch os.Args[1] {
    case "serve":
        addServerFlags(serveFlags, c)
        addDBFlags(serveFlags, c)
        serveFlags.Parse(os.Args[2:])
    case "createdb":
        addDBFlags(serveFlags, c)
        createDBFlags.Parse(os.Args[2:])
    default:
        return errors.New("Unknown command: "+os.Args[1])

    }
    return nil

}

func buildDBConnString(host string, port int, dbname string, user string, pass string, ssl string) (string) {
    // host=localhost user=cgpipe_mon password=thisismypassword dbname=cgpipe_mon port=5432 sslmode=disable
    s := ""

    if host != "" {
        s += "host='"+host+"' "
    }

    if port > 0 {
        s += fmt.Sprintf("port=%d ",port)
    }

    if dbname != "" {
        s += "dbname='"+dbname+"' "
    }

    if user != "" {
        s += "user='"+user+"' "
    }

    if pass != "" {
        s += "password='"+pass+"' "
    }

    if ssl != "" {
        s += "sslmode='"+ssl+"' "
    }
    return s
}


