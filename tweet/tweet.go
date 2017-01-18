package tweet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"osc-tweet/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/bitly/go-simplejson"
	"github.com/gogather/com"
)

var UA int
var Devmode = false

const (
	DEV_TWEET_URL  = "http://www.oschina.io/action/apiv2/tweet"
	PROD_TWEET_URL = "https://www.oschina.net/action/apiv2/tweet"
)

type OneResult struct {
	Result   string
	HpEntity HpEntityStc
}

type HpEntityStc struct {
	StrAutor   string
	StrContent string
}

func Tweet(message string) {
	if message == "" {
		message = "^_^"
	}
	// read info
	// var data interface{}
	var err error

	home := utils.GetHome()
	pathUserInfo := filepath.Join(home, ".osc", "userinfo")
	if com.FileExist(pathUserInfo) {
		// json, _ := com.ReadFile(pathUserInfo)
		// data, err = com.JsonDecode(json)
		if err != nil {
			log.Errorln("[Error]", "Parse userinfo file failed")
			return
		}
	} else {
		log.Warnln("login first")
		return
	}

	http := &utils.Http{}
	var response string
	if Devmode {
		response, err = http.Post(DEV_TWEET_URL, fmt.Sprintf("content=%s", message), false, 0)
	} else {
		response, err = http.Post(PROD_TWEET_URL, fmt.Sprintf("content=%s", message), false, 0)
	}

	if err != nil {
		log.Warnln("[Error]", err)
	}

	json, err := simplejson.NewJson([]byte(response))
	if err != nil {
		log.Errorln("发送失败")
		log.Errorln(response)
		return
	}
	code, _ := json.Get("code").Int()
	if code == 0 {
		msg, _ := json.Get("message").String()
		log.Errorln("发送失败：", msg)
		return
	}
	log.Infoln("发送成功")

}

func Joke() {
	api := `http://www.tuling123.com/openapi/api?key=380abd77ba6541dd1dee43220c42776b&info=%E8%AE%B2%E4%B8%AA%E7%AC%91%E8%AF%9D`
	http := &utils.Http{}
	msg, err := http.Get(api)
	if err != nil {
		log.Errorln(err)
	}

	data, err := com.JsonDecode(msg)
	if err != nil {
		log.Errorln(err)
	}

	json := data.(map[string]interface{})
	msg = json["text"].(string)

	reg := regexp.MustCompile(`<[\d\D]+>`)
	msg = reg.ReplaceAllString(msg, "")

	msg = com.SubString(msg, 0, 190)

	Tweet(msg)
}

func Weather(location string) {
	api := `http://www.tuling123.com/openapi/api?key=380abd77ba6541dd1dee43220c42776b&info=%E4%BB%8A%E5%A4%A9` + location + `%E5%A4%A9%E6%B0%94`
	http := &utils.Http{}
	msg, err := http.Get(api)
	if err != nil {
		log.Errorln(err)
	}

	data, err := com.JsonDecode(msg)
	if err != nil {
		log.Errorln(err)
	}

	json := data.(map[string]interface{})
	msg = json["text"].(string)

	reg := regexp.MustCompile(`<[\d\D]+>`)
	msg = reg.ReplaceAllString(msg, "")

	msg = com.SubString(msg, 0, 190)

	Tweet(msg)
}

func One() {
	t := time.Now()
	date := t.Format("2006-01-02")
	var res OneResult
	url := "http://211.152.49.184:7001/OneForWeb/one/getHpinfo?strDate=" + date
	resp := httpGet(url)
	json.Unmarshal([]byte(resp), &res)
	Tweet(res.HpEntity.StrContent)
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("wrong", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}
