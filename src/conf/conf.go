package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Conf struct {
	Host		string	`json:"host"`
	Port		string	`json:"port"`
	Database	string `json"database"`
}

func NewConf(config_file string) *Conf {
	c := new(Conf)
	file, err := ioutil.ReadFile(config_file)
	if err != nil {
		log.Printf("%v\n", err)
	}

	json.Unmarshal(file, c)

	return c
}
