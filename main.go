package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/robfig/cron/v3"
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

	// cron
	c := cron.New(cron.WithSeconds()) // add seconds to standard cron
	c.AddFunc(config["cron"], func() {
		fmt.Println("working...")

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
			if status {
				sendMail(config["from"], config["title"], strconv.FormatBool(status), config["to"])
			}

		} else {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
	})
	c.Start()
	select {} // keep main alive
}

func sendMail(from, title, body, to string) error {
	mail := gomail.NewMessage()
	mail.SetHeader(`From`, from)
	mail.SetHeader(`Subject`, title)
	mail.SetBody(`text/html`, body)
	mail.SetHeader(`To`, to)

	host := config["host"]
	port, _ := strconv.Atoi(config["port"])
	user := config["from"]
	password := config["password"]

	sender := gomail.NewDialer(host, port, user, password)
	sender.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := sender.DialAndSend(mail)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("mail sent to=%s\n", to)
	return err
}
