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
var devmode bool

func init() {
	flag.StringVar(&option, "o", "", "what operation you want,login or tweet")
	flag.StringVar(&username, "u", "", "account username")
	flag.StringVar(&password, "p", "", "account password")
	flag.StringVar(&message, "m", "", "tweet message")
	flag.StringVar(&configpath, "c", "", "config file path")
	flag.BoolVar(&devmode, "dev", false, "dev mode to visit the url: oschina.com")
	flag.IntVar(&phone, "ua", 0, "ua, 0 is iphone, 1 is android, default is 0")
	flag.Parse()
}
func Run() {
	initProfileDir()

	// check if config file is exist
	if configpath != "" {
		Config(configpath)
		return
	}

	utils.UA = phone
	login.Devmode = devmode
	tweet.Devmode = devmode

	if option == "" {
		flag.Usage()
		return
	}
	if username == "" || password == "" {
		log.Redln("用户名和密码必须有 使用 -u xxx -p xxxx -o xxx")
		return
	}

	login.Login(username, password)

	switch option {
	// 登陆
	case "login":
		login.Login(username, password)
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
