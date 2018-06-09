package main

import (
	"database/sql"
	"encoding/json"
	"esports-league-manager/Backend/Server/databaseAccess"
	"esports-league-manager/Backend/Server/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"os"
)

//TODO: put this in its own file preferrably in own package
func openDB(db *sql.DB, user, pass, dbName string) {

}

type configuration struct {
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
	DbName string `json:"dbName"`
	Port   string `json:"port"`
}

func newApp(db *sql.DB) *gin.Engine {
	app := gin.Default()

	routes.UsersDAO = databaseAccess.CreateUsersDao(db)

	//create session manager with secure cookies
	//cookieName := "elmsession"
	//hashKey := securecookie.GenerateRandomKey(24)
	//blockKey := securecookie.GenerateRandomKey(24)
	//secureCookie := securecookie.New(hashKey, blockKey)
	//
	//elmSessions := sessions.New(sessions.Config{
	//	Cookie: cookieName,
	//	Encode: secureCookie.Encode,
	//	Decode: secureCookie.Decode,
	//})
	//routesTest.RegisterLoginHandlers(app.Party("/login"), &db, elmSessions)
	routes.RegisterUserHandlers(app.Group("/api/users"))
	//routesTest.RegisterLeagueHandlers(app.Party("/api/leagues"), &db, elmSessions)

	return app
}

func main() {
	//get program config
	file, err := os.Open("Backend/conf.json")
	if err != nil {
		log.Fatal("error opening config: ", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config configuration
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("error decoding config: ", err)
		return
	}

	//create database connection
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPass, config.DbName)
	db, err := sql.Open("pgx", connStr)
	defer db.Close()

	if err != nil {
		log.Fatal("error opening db: ", err)
	}

	//start router/webapp
	app := newApp(db)
	app.Run(config.Port)
}
