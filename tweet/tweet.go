package tweet

import (
	"github.com/Laily123/osc-tweet/utils"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"net/url"
	"path/filepath"
	"regexp"
    "time"
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
)

type OneResult struct{
    Result string
    HpEntity HpEntityStc
}

type HpEntityStc struct{
    StrAutor string
    StrContent string
}

func Tweet(message string) {
	// read info
	var data interface{}
	var err error

	home := utils.GetHome()
	pathUserInfo := filepath.Join(home, ".osc", "userinfo")
	if com.FileExist(pathUserInfo) {
		json,_ := com.ReadFile(pathUserInfo)
		data, err = com.JsonDecode(json)
		if err != nil {
			log.Redln("[Error]", "Parse userinfo file failed")
			return
		}
	} else {
		log.Warnln("login first")
		return
	}

	jsonData, ok := data.(map[string]interface{})
	if !ok {
		log.Redln("[Error]", "illeage data")
		return
	}

	userId, ok := jsonData["user"].(string)
	if !ok {
		log.Redln("[Error]", "get user id failed")
		return
	}

	userCode, ok := jsonData["user_code"].(string)
	if !ok {
		log.Redln("[Error]", "get user code failed")
		return
	}

	http := &utils.Http{}
	response, err := http.Post("http://www.oschina.net/action/tweet/pub", url.Values{
		"user":      {userId},
		"user_code": {userCode},
		"msg":       {message},
	})

	if err != nil {
		log.Warnln("[Error]", err)
	}

	tweetResult, err := com.JsonDecode(response)
	if err != nil {
		log.Redln("发送失败")
		return
	} else {
		_, ok := tweetResult.(map[string]interface{})["error"].(float64)
		if ok {
			log.Warnln(tweetResult.(map[string]interface{})["msg"])
			return
		}

		id, ok := tweetResult.(map[string]interface{})["log"].(float64)
		if ok {
			log.Greenf("开源中国第[%d]条动态发送成功\n", int64(id))
			log.Blueln(message)
		}
	}

}

func Joke() {
	api := `http://www.tuling123.com/openapi/api?key=380abd77ba6541dd1dee43220c42776b&info=%E8%AE%B2%E4%B8%AA%E7%AC%91%E8%AF%9D`
	http := &utils.Http{}
	msg, err := http.Get(api)
	if err != nil {
		log.Redln(err)
	}

	data, err := com.JsonDecode(msg)
	if err != nil {
		log.Redln(err)
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
		log.Redln(err)
	}

	data, err := com.JsonDecode(msg)
	if err != nil {
		log.Redln(err)
	}

	json := data.(map[string]interface{})
	msg = json["text"].(string)

	reg := regexp.MustCompile(`<[\d\D]+>`)
	msg = reg.ReplaceAllString(msg, "")

	msg = com.SubString(msg, 0, 190)

	Tweet(msg)
}

func One(){
    t := time.Now()
    date := t.Format("2006-01-02")
    var res OneResult
    url := "http://211.152.49.184:7001/OneForWeb/one/getHpinfo?strDate="+date
    resp := httpGet(url)
    json.Unmarshal([]byte(resp),&res)
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


