package main

import (
	"encoding/json"
	"io/ioutil"
)

var cfgFile = "/etc/gonoip.conf"

type host struct {
	Name string
	Auth string
}

var hosts []*host

func readConfig() error {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &hosts)
	return err
}
