package initial

import (
	"fmt"
	"os"
	"osc-tweet/login"
	"osc-tweet/tweet"
	"osc-tweet/utils"
	"path/filepath"

	"github.com/gogather/com"
	"github.com/gogather/com/log"
)

const (
	VERSION = "0.0.0"
)

func Run() {
	initProfileDir()

	length := len(os.Args)

	username := ""
	password := ""
	message := ""

	if length < 2 {
		showHelp()
	} else {
		if os.Args[1] == "login" {
			if length < 4 {
				log.Dangerln("Invalid command, please use")
				log.Warnln("    osc login username password")
			} else {
				username = os.Args[2]
				password = os.Args[3]

				login.Login(username, password)
			}
		} else if os.Args[1] == "tweet" {
			if length < 3 {
				log.Dangerln("Invalid command, please use")
				log.Warnln("    osc tweet message")
			} else {
				message = os.Args[2]

				tweet.Tweet(message)
			}
		} else if os.Args[1] == "status" {
			login.GetStatus()
		} else if os.Args[1] == "joke" {
			tweet.Joke()
		} else if os.Args[1] == "weather" {
			location := "%E6%B7%B1%E5%9C%B3"
			if len(os.Args) >= 3 {
				location = os.Args[2]
			}
			tweet.Weather(location)
		} else if os.Args[1] == "one" {
			tweet.One()
		} else if os.Args[1] == "help" {
			showHelp()
		} else {
			log.Dangerln("Invalid command, please use")
			log.Warnln("    osc help")
		}
	}
}

func showHelp() {
	log.Warnln("oschina command line tool")
	fmt.Println("Usage:")
	fmt.Println("    login user password")
	fmt.Println("    status")
	fmt.Println("    tweet message")
	log.Blueln("version", VERSION)
}

func initProfileDir() {
	home := utils.GetHome()
	path := filepath.Join(home, ".osc")

	if !com.FileExist(path) {
		err := com.Mkdir(path)
		if err != nil {
			log.Fatalln("Create profile directory failed!")
		}
	}
}

func readConfig() {

}
