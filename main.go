package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

const (
	ConfigFile = "config.json"
)

type Config map[string]string

var config = Config{}

func main() {

	// get run path
	program, err := os.Executable()
	if err != nil {
		panic(err)
	}
	programPath := filepath.Dir(program)
	fmt.Printf("program path=%s\n", programPath)

	// read config
	configPath := filepath.Join(programPath, ConfigFile)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &config)

	// download
	res, err := http.Get(config["url"])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res := string(bytes)

		// check
		status := strings.Contains(res, "ATCButton")
		fmt.Printf("status=%v\n", status)
		sendMail(config["user"], config["user"], "gonike", strconv.FormatBool(status))

	} else {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

}

func sendMail(to string, from, title, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader(`From`, from)
	mail.SetHeader(`To`, to)
	mail.SetHeader(`Subject`, title)
	mail.SetBody(`text/html`, body)

	host := config["host"]
	port, _ := strconv.Atoi(config["port"])
	user := config["user"]
	password := config["password"]

	err := gomail.NewDialer(host, port, user, password).DialAndSend(mail)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
