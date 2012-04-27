package main

import (
	"flag"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"
)

var agent = "gonoip/0.1 costinc@gmail.com"
var ipfile = "/var/run/gonoip"

var infoLog, errLog *log.Logger
var forceUpdate bool
var addHost bool

func init() {
	flag.BoolVar(&forceUpdate, "f", false, "update noip even if IP hasn't changed")
	flag.BoolVar(&addHost, "a", false, "print config with added host to stdout")
	flag.Parse()

	if addHost {
		infoLog = log.New(os.Stdout, "", 0)
		errLog = log.New(os.Stdout, "", 0)
	} else {
		var err error
		infoLog, err = syslog.NewLogger(syslog.LOG_INFO, 0)
		checkErr(err)
		errLog, err = syslog.NewLogger(syslog.LOG_ERR, 0)
		checkErr(err)
	}

	err := readConfig()
	if addHost && os.IsNotExit(err) {
		return
	}
	checkErr(err)
}

func main() {
	if addHost {
		AddHost()
		return
	}

	if len(hosts) == 0 {
		errLog.Fatal("no hosts defined")
	}

	if forceUpdate {
		infoLog.Print("forced update")
		for _, h := range hosts {
			updateNoIP(h)
		}
		return
	}

	ip, oldip := getIP(), getOldIP()
	if ip != oldip {
		infoLog.Print("IP changed, new IP: " + ip)
		for _, h := range hosts {
			updateNoIP(h)
		}
	} else {
		infoLog.Print("same IP, no update required")
	}
}

func updateNoIP(h *host) {
	var client http.Client

	req, err := http.NewRequest(
		"GET",
		"http://dynupdate.no-ip.com/nic/update?hostname="+h.Name,
		nil)
	checkErr(err)

	req.Header.Add("Authorization", "Basic "+h.Auth)
	req.Header.Add("User-Agent", agent)

	resp, err := client.Do(req)
	checkErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	args := strings.Split(string(body), " ")
	handleNoIP(h, args)
}

func handleNoIP(h *host, args []string) {
	prefix := h.Name + ": noip: "

	switch args[0] {
	case "good":
		infoLog.Print(prefix + "IP changed (" + args[1] + ")")
		err := ioutil.WriteFile(ipfile, []byte(args[1]), 0644)
		checkErr(err)
	case "nochg":
		infoLog.Print(prefix + "no change required (" + args[1] + ")")
		err := ioutil.WriteFile(ipfile, []byte(args[1]), 0644)
		checkErr(err)
	case "nohost":
		errLog.Fatal(prefix + "noip: host does not exist")
	case "badauth":
		errLog.Fatal(prefix + "invalid username or password")
	case "badagent":
		errLog.Fatal(prefix + "bad user agent: client disabled")
	case "!donator":
		errLog.Fatal(prefix + "feature not available")
	case "abuse":
		errLog.Fatal(prefix + "user blocked due to abuse")
	case "911":
		infoLog.Print(prefix + "fatal server error")
	}
}

func getIP() string {
	resp, err := http.Get("http://automation.whatismyip.com/n09230945.asp")
	checkErr(err)

	ip, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return string(ip)
}

func getOldIP() string {
	ip, err := ioutil.ReadFile(ipfile)
	if os.IsNotExist(err) {
		return ""
	}
	checkErr(err)

	return string(ip)
}

func checkErr(err error) {
	if err != nil {
		errLog.Fatal(err)
	}
}
