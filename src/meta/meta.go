package meta

import (
	"github.com/etgryphon/stringUp"
	"crypto/x509/pkix"
	"sync/atomic"
)

const TYPE_TMPL =
	"type {{ .TableTitle }} struct {"
	+"{{ $column := range .Columns }}"
	+"	{{$column.Title}} {{$column.GoType}}	`{{$tag := range $column.Tags}}{{$tag.Name}}:\"{{$tag.Value}}\"	{{end}}`"
	+"{{ end }}"
	+"}"

var TYPE_DIC map[string]string
var MODE_DIC map[string]string

func init() {
	TYPE_DIC = map[string]string{
		"int":"int64",
		"int4":"int64",
		"number":"float64",
		"timestamp":"time.Time",
		"timestamptz":"time.Time",
	}
	MODE_DIC = map[string]string{
		"false":"NULLABLE",
		"true":"REQIRED",
	}
}

type Tag struct {
	Name	string
	Value	string
}

type ColumnMeta struct {
	Title	string	`json:"title"`
	Name	string	`json:"name"`
	Type	string	`json:"type"`
	Primary	bool	`json:"primary"`
	GoType	string	`json:"go_type"`
	Mode	string	`json:"mode"`
	Tags	[]Tag
}

type TableMeta struct {
	TableName string `json:"table_name"`
	TableTitle string `json:"table_title"`
	Columns []ColumnMeta
}

func NewTableMeta(table_name string) *TableMeta {
	t := new(TableMeta)
	t.TableName = table_name
	t.TableTitle = stringUp.CamelCase(table_name)
	return t
}

func (t TableMeta)AddColumnMeta(name, typ, mode string) TableMeta {
	var cm ColumnMeta
	cm.Name = name
	cm.Type = typ
	cm.GoType = TYPE_DIC[typ]
	cm.Mode = MODE_DIC[mode]
	return append(t.Columns, cm)
}
