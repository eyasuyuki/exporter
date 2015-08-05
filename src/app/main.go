package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"text/template"
	"bytes"
	"fmt"
	"conf"
	"log"
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

const SELECT_ALL = `
SELECT * FROM {{.TableName}}
`
var TMPL *template.Template

type ColumnMeta struct {
	ColumnName string `json:"column_name"`
	ColumnType string `json:"column_type"`
}

type TableMeta struct {
	TableName string `json:"table_name"`
	Columns []ColumnMeta
}

func (t *TableMeta)MakeQuery() (string, error) {
	var query bytes.Buffer
	err := TMPL.Execute(&query, t)
	return query.String(), err
}

func init() {
	c = conf.NewConf("./conf.json")
	TMPL = template.Must(template.New("query").Parse(SELECT_ALL))
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
		meta := new(TableMeta)
		meta.TableName = table_name
		query, _ := meta.MakeQuery()
		fmt.Print(query)
		res, err := conn.Query(query)
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		defer res.Close()
		columns, err := res.Columns()
		if err != nil {
			log.Fatalf("error: %v", err);
		}
		fmt.Printf("columns=%d", len(columns))
		for c := range columns {
			fmt.Printf("c = %V", c)
		}
	}

}
