package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"conf"
	"meta"
	"log"
	"os"
	"bytes"
)

var c *conf.Conf
var mode_dic map[string]string

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
    attname
    ,typname
    ,attnotnull
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

const SELECT_PRIMARY_KEY = `
SELECT
	pg_attribute.attname
FROM
	pg_index
JOIN
	pg_attribute ON pg_attribute.attrelid = pg_index.indrelid
    AND pg_attribute.attnum = ANY(pg_index.indkey)
WHERE
	pg_index.indrelid = $1::regclass
	AND pg_index.indisprimary;
`

func init() {
	c = conf.NewConf("./conf.json")
	mode_dic = make(map[string]string)
	mode_dic["false"] = "NULLABLE"
	mode_dic["true"] = "REQUIRED"
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
	var tables []meta.TableMeta
	for rows.Next() {
		var table_name string
		rows.Scan(&table_name)
		fmt.Printf("table=%v\n", table_name)
		tbl := meta.NewTableMeta(table_name)
		tbl.TableName = table_name

		stmt, err := conn.Prepare(SELECT_PRIMARY_KEY)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		defer stmt.Close()

		row, err := stmt.Query(table_name)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		var pk string
		for row.Next() {
			row.Scan(&pk)
			break
		}

		st2, err := conn.Prepare(SELECT_TABLE_META)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		defer st2.Close()
		rw2, err := st2.Query(table_name)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		var n string
		var t string
		var m string
		for rw2.Next() {
			rw2.Scan(&n, &t, &m)
			fmt.Printf("column=%v, type=%v, mode=%v\n", n, t, mode_dic[m])
			cm := meta.NewColumnMeta(n, t, m, pk)
			tbl.Columns = append(tbl.Columns, cm)
		}
		tables = append(tables, tbl)
		var buf bytes.Buffer
		tbl.Export(&buf)
		fmt.Println(buf.String())
	}
}

func main() {
	doExport()
	os.Exit(0)
}
