package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"conf"
	"log"
)

var c *conf.Conf

const SELECT_TABLE_NAMES = `
SELECT table_name
  FROM information_schema.tables
 WHERE table_schema='public'
   AND table_type='BASE TABLE'
`

func init() {
	c = conf.NewConf("./conf.json")
}

func main() {
	fmt.Printf("%v:%v, %v\n", c.Host, c.Port, c.Database)

	var spec = "user=postgres host="+c.Host+" port="+c.Port+" dbname="+c.Database+" sslmode=disable"
	fmt.Println(spec)

	conn, err := sql.Open("postgres", spec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer conn.Close()

	rows, err := conn.Query(SELECT_TABLE_NAMES)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table_name string
		rows.Scan(&table_name)
		fmt.Println(table_name)
	}

}
