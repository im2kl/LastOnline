package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	router := gin.Default()
	router.POST("/status", func(c *gin.Context) {

		b := c.Request.Header.Get("X-Device-Key")
		if b == "" {
			c.String(http.StatusBadRequest, "")
		} else {
			go setStatus(b)
			c.String(http.StatusOK, "%s", b)
		}

	})
	router.GET("/status/:id", func(c *gin.Context) {
		// same than  c.Input.FromPath("id") in this context
		id := c.Param("id")
		if id == "" {
			c.String(http.StatusUnauthorized, "No id provided")
		}
		rt, err := getStatus(id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		println(id)
		c.String(http.StatusOK, rt)
		// ... read account with id
	})
	router.Run()
}

func init() {
	//os.Remove("./DeviceStatus.db")
	var err error
	db, err = sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	sqlStmt := `
		create table status(id string not null primary key, name integer);
		`
	/*
		sqlStmt := `
		create table status(id string , name integer);
		`*/
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	//statusUpdate()
	go test()
}

func setStatus(id string) error {

	now := time.Now()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// INSERT INTO table (id, name, age) VALUES(1, "A", 19)
	//stmt, err := tx.Prepare("insert into status(id, name) values(?, ?) ON DUPLICATE KEY UPDATE name=?")
	stmt, err := tx.Prepare("REPLACE into status(id, name) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, now.Unix())
	if err != nil {
		return err
	}

	tx.Commit()
	return nil

}

func getStatus(id string) (string, error) {
	//rows, err := db.Query(fmt.Sprintf("select id, name from status where id=%s", id))
	rows, err := db.Query(fmt.Sprintf("select id, name from status"))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name int
		err = rows.Scan(&id, &name)
		if err != nil {
			return "", err
		}
		fmt.Println(id, name)
		return fmt.Sprintf("{id:%s,time:%d}", id, name), nil
	}
	err = rows.Err()
	return "", err

}

func test() {
	for {
		time.Sleep(10 * time.Second)

		rows, err := db.Query("select id, name from status")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			var name int
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
