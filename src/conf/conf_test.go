package conf

import (
	"testing"
	"conf"
)

const HOST = "192.168.59.103"
const PORT = "5432"
const DATABASE = "db01"

const CONF_JSON = "../../conf.json"

func TestConf(t *testing.T) {
	c := conf.NewConf(CONF_JSON)
	if c.Host != HOST {
		t.Errorf("error: c.Host=%v, expected=%v", c.Host, HOST)
	} else if c.Port != PORT {
		t.Errorf("error: c.Port=%v, expected=%v", c.Port, PORT)
	} else if c.Database != DATABASE {
		t.Errorf("error: c.Host=%v, expected=%v", c.Host, DATABASE)
	}
}
