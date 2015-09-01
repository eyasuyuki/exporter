package meta

import (
	"text/template"
	"io"
	"regexp"
	"bytes"
)

const STRUCT_TMPL = `
type {{ .TableTitle }} struct { {{ range $index, $column := .Columns }}
	{{$column.Title}} {{$column.GoType}}	{{$column.Tags}}{{ end }}
}`

var TMPL *template.Template

var TYPE_DIC map[string]string
var MODE_DIC map[string]string

func init() {
	TYPE_DIC = map[string]string{
		"int":"int64",
		"int4":"int64",
		"int8":"int64",
		"number":"float64",
		"text":"string",
		"char":"string",
		"varchar":"string",
		"timestamp":"time.Time",
		"timestamptz":"time.Time",
	}
	MODE_DIC = map[string]string{
		"false":"NULLABLE",
		"true":"REQIRED",
	}
	TMPL = template.Must(template.New("struct").Parse(STRUCT_TMPL))
}


type ColumnMeta struct {
	Title	string	`json:"title"`
	Name	string	`json:"name"`
	Type	string	`json:"type"`
	Primary	bool	`json:"primary"`
	GoType	string	`json:"go_type"`
	Mode	string	`json:"mode"`
	Tags	string
}

type TableMeta struct {
	TableName string `json:"table_name"`
	TableTitle string `json:"table_title"`
	Columns []ColumnMeta
}

func NewTableMeta(table_name string) TableMeta {
	var t TableMeta
	t.TableName = table_name
	t.TableTitle = Camel(table_name)
	return t
}

func NewColumnMeta(name, typ, mode, pk string) ColumnMeta {
	var cm ColumnMeta
	cm.Name = name
	cm.Title = Camel(name)
	cm.Type = typ
	cm.Primary = name == pk
	cm.Mode = MODE_DIC[mode]
	if cm.Mode == "NULLABLE" {
		cm.GoType = "*"+TYPE_DIC[typ]
		cm.Tags = "`json:\""+name+",omitempty\"	column:\""+name+"\""
	} else {
		cm.GoType = TYPE_DIC[typ]
		cm.Tags = "`json:\""+name+"\"	column:\""+name+"\""
	}
	if cm.Primary {
		cm.Tags = cm.Tags+"	db:\"pk\""
	}
	cm.Tags = cm.Tags+"`"
	return cm
}

func (t TableMeta)Export(w io.Writer) error {
	return TMPL.Execute(w, t)
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func Camel(src string)(string){
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
	}
	return string(bytes.Join(chunks, nil))
}
