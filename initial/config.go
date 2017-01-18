package initial

import (
	"fmt"
	"osc-tweet/login"
	"osc-tweet/tweet"
	"osc-tweet/utils"

	"github.com/go-ini/ini"
	"github.com/gogather/com/log"
)

// 读取配置文件
func Config(configpath string) {
	if configpath == "" {
		log.Redln("配置文件不存在")
		return
	}

	cfg, err := ini.InsensitiveLoad(configpath)
	printErr(err, "读取配置文件失败")

	devmode := cfg.Section("config").Key("devmode").MustBool(false)
	login.Devmode = devmode
	tweet.Devmode = devmode
	ua := cfg.Section("config").Key("ua").MustInt(0)
	utils.UA = ua

	username := cfg.Section("user").Key("name").String()
	pwd := cfg.Section("user").Key("pwd").String()
	if username == "" || pwd == "" {
		log.Redln("用户名和密码必须配置")
	}
	login.Login(username, pwd)

	step := cfg.Section("config").Key("iterator").MustInt(0)
	var content string
	if step == 0 {
		content = "^_^"
	} else {
		content = cfg.Section("content").Key(fmt.Sprintf("#%d", step)).String()
	}
	if content == "" {
		content = "^_^"
	}
	tweet.Tweet(content)
	step++
	cfg.Section("config").Key("iterator").SetValue(fmt.Sprintf("%d", step))
	cfg.SaveTo(configpath)
}

func printErr(err error, msg string) {
	if err != nil {
		log.Redln(msg, err)
	}
}
