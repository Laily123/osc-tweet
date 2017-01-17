package initial

import (
	"os"
	"osc-tweet/login"
	"osc-tweet/tweet"
	"osc-tweet/utils"
	"path/filepath"

	"flag"

	"github.com/gogather/com"
	"github.com/gogather/com/log"
)

const (
	VERSION = "0.0.0"
)

var option string
var username string
var password string
var message string
var phone int
var configpath string

func init() {
	flag.StringVar(&option, "o", "", "what operation you want,login or tweet")
	flag.StringVar(&username, "name", "", "account username")
	flag.StringVar(&password, "pwd", "", "account password")
	flag.StringVar(&message, "message", "", "tweet message")
	flag.StringVar(&configpath, "c", "", "config file path")
	flag.IntVar(&phone, "ua", 0, "ua, 0 is iphone, 1 is android, default is 0")
	flag.Parse()
}
func Run() {
	initProfileDir()

	if option == "" {
		flag.Usage()
		return
	}
	utils.UA = phone
	switch option {
	// 登陆
	case "login":
		if username == "" && password == "" {
			log.Dangerln("Invalid command, please use")
			log.Warnln("  -o login -name username -p password")
		} else {
			login.Login(username, password)
		}
		// 直接发送动弹
	case "tweet":
		if message == "" {
			log.Dangerln("Invalid command, please use")
			log.Warnln(" -o tweet -m message")
		} else {
			tweet.Tweet(message)
		}
	case "status":
		login.GetStatus()

	case "joke":
		tweet.Joke()
	case "weather":
		location := "%E6%B7%B1%E5%9C%B3"
		if len(os.Args) >= 3 {
			location = os.Args[2]
		}
		tweet.Weather(location)
	case "one":
		tweet.One()
	case "auto":
		Config(configpath)
	case "help":
		flag.Usage()
	default:
		log.Dangerln("Invalid command, please use")
		flag.Usage()
	}
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
