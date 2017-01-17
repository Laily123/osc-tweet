package login

import (
	"net/url"
	"osc-tweet/utils"
	"path/filepath"
	"regexp"

	"github.com/bitly/go-simplejson"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
)

const (
	DEV_URL  = "http://www.oschina.com/action/apiv2/login_validate"
	PROD_URL = "https://www.oschina.net/action/apiv2/login_validate"
)

func Login(username string, password string) {
	home := utils.GetHome()
	pathPwd := filepath.Join(home, ".osc", "password")

	// password = utils.SHA1(password)
	com.WriteFile(pathPwd, password)

	pathUsr := filepath.Join(home, ".osc", "username")
	com.WriteFile(pathUsr, username)

	http := &utils.Http{}
	response, err := http.Post(DEV_URL, url.Values{
		"username": {username},
		"pwd":      {password},
	}, true)

	if err != nil {
		log.Warnln("请检查网络")
		log.Redln(err)
		return
	}

	json, err := simplejson.NewJson([]byte(response))
	if err != nil {
		log.Redln("登陆失败：", err)
	}
	code, _ := json.Get("code").Int()
	if code == 1 {
		log.Greenln("登录成功")
	} else {
		msg, _ := json.Get("message").String()
		log.Redln("登录失败: ", msg)
	}
}

// get user_code
func getUserCode() {
	http := &utils.Http{}
	response, err := http.Get("https://www.oschina.net")
	if err != nil {
		log.Redln("[Error]", err)
		return
	}

	regex1 := `(^[\d\D]*)(name='user_code' value=')([\d\D][^\/]+)('\/>)([\d\D]*$)`
	reg := regexp.MustCompile(regex1)
	userCode := reg.ReplaceAllString(response, "$3")

	regex2 := `(^[\d\D]*)(<input type='hidden' name='user' value=')([\d][^']+)('\/>)([\d\D]*$)`
	reg = regexp.MustCompile(regex2)
	userId := reg.ReplaceAllString(response, "$3")

	content, _ := com.JsonEncode(map[string]interface{}{
		"user":      userId,
		"user_code": userCode,
	})

	home := utils.GetHome()
	pathUserCode := filepath.Join(home, ".osc", "userinfo")
	com.WriteFile(pathUserCode, content)
}
