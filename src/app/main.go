package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"conf"
	"log"
	"os"
)

var c *conf.Conf

const SELECT_TABLE_NAMES = `
SELECT
	table_name
FROM
	information_schema.tables
WHERE
	table_schema='public'
	AND table_type='BASE TABLE'
`

const SELECT_TABLE_META = `
SELECT
    attname,typname
FROM
    pg_attribute
    ,pg_type
WHERE
    attrelid = $1::regclass
    AND pg_attribute.attnum > 0
    AND pg_attribute.atttypid=pg_type.oid
ORDER BY
	attnum ASC
`
type ColumnMeta struct {
	ColumnName string `json:"column_name"`
	ColumnType string `json:"column_type"`
}

type TableMeta struct {
	TableName string `json:"table_name"`
	Columns []ColumnMeta
}

func init() {
	c = conf.NewConf("./conf.json")
}

func doExport() {
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
		fmt.Printf("table=%v\n", table_name)
		stmt, err := conn.Prepare(SELECT_TABLE_META)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		defer stmt.Close()
		row, err := stmt.Query(table_name)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		var n string
		var t string
		for row.Next() {
			row.Scan(&n, &t)
			fmt.Printf("column=%v, type=%v\n", n, t)
		}
	}
}

func main() {
	doExport()
	os.Exit(0)
}
