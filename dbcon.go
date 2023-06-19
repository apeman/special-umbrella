package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
)
var db = dbConnect()

func dbConnect() *sql.DB{
dbDriver := "mysql"
    dbUser := "root"
    dbPass := "apeman"
    dbName := "imgbd"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
   if err != nil {
    log.Fatal(err)
  }
  return db
}

