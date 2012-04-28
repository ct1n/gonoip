package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"os"
)

func readValue(key string) string {
	stdin := bufio.NewReader(os.Stdin)

	_, err := os.Stderr.Write([]byte(key + ": "))
	checkErr(err)
	line, _, err := stdin.ReadLine()
	checkErr(err)

	return string(line)
}

func AddHost() {
	name := readValue("Host")
	user := readValue("User")
	pass := readValue("Password")

	auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))

	hosts = append(hosts, &host{Name: name, Auth: auth})

	data, err := json.MarshalIndent(hosts, "", "\t")
	checkErr(err)
	os.Stdout.Write(data)
}
