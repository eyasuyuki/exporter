package meta

import (
	_ "github.com/lib/pq"
	"database/sql"
	"conf"
	"testing"
	"fmt"
	"app"
	"meta"
)

const CREATE = `
DROP TABLE IF EXISTS test01;
CREATE TABLE test01 (
	col01 bigint
);
`

var CONF *conf.Conf

func init() {
	CONF = conf.NewConf("../../conf.json")
}

func TestMeta(t *testing.T) {
	var spec = "user=postgres host="+CONF.Host+" port="+CONF.Port+" dbname="+CONF.Database+" sslmode=disable"
	fmt.Println(spec)
	conn, err := sql.Open("postgres", spec)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	defer conn.Close()

	_, err = conn.Exec(CREATE)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	stmt, err := conn.Prepare(main.SELECT_TABLE_META)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	defer stmt.Close()

	row, err := stmt.Query("test01")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	var m meta.ColumnMeta
	var n,tt,mm string
	for row.Next() {
		row.Scan(&n, &tt, &mm)
		m = meta.NewColumnMeta(n, tt, mm, "")
		if n == "col01" {
			break
		}
		
	}

	if m.Name != "col01" {
		t.Errorf("m.Name=%v, expected=%v", m.Name, "col01")
	} else if m.Type != "int8" {
		t.Errorf("m.Type=%v, expected=%v", m.Type, "int8")
	} else if m.GoType != "*int64" {
		t.Errorf("m.GoType=%v, expected=%v", m.GoType, "*int64")
	}

	
	row.Close()
	
}
