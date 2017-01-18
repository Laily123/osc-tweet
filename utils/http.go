package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gogather/com"
)

var UA int

const (
	UA_IPHONE  = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_2 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13C75 Safari/601.1"
	UA_ANDRIOD = "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19"
)

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (this *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	this.lk.Lock()
	this.cookies[u.Host] = cookies
	this.lk.Unlock()
}

func (this *Jar) Cookies(u *url.URL) []*http.Cookie {
	return this.cookies[u.Host]
}

func (this *Jar) ParseCookies(json string) *http.Cookie {
	c := &http.Cookie{}
	cookiesObj, err := com.JsonDecode(json)
	cookies := cookiesObj.(map[string]interface{})

	if err == nil {
		if err == nil {
			// set cookie
			c.Name = cookies["Name"].(string)
			c.Value = cookies["Value"].(string)
			c.Path = cookies["Path"].(string)
			c.Domain = cookies["Domain"].(string)
			c.RawExpires = cookies["RawExpires"].(string)
		}
	}

	return c
}

type Http struct {
	cookies *Jar
}

func (this *Http) Post(urlstr string, parm string, storeCookies bool, ua int) (string, error) {
	home := GetHome()
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}

	pathOscid := filepath.Join(home, ".osc", "oscid")
	jar := NewJar()

	// read cookie
	if com.FileExist(pathOscid) {
		json, _ := com.ReadFile(pathOscid)
		c := jar.ParseCookies(json)
		jar.SetCookies(u, []*http.Cookie{c})
	}

	client := &http.Client{nil, nil, jar, 0}
	params := strings.NewReader(parm)
	req, err := http.NewRequest("POST", urlstr, params)
	if err != nil {
		log.Println("new req error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if UA == 1 {
		req.Header.Set("User-Agent", UA_ANDRIOD)
	} else {
		req.Header.Set("User-Agent", UA_IPHONE)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	// store cookie
	if storeCookies {
		cookieMap := jar.Cookies(u)
		length := len(cookieMap)
		if length > 0 {
			co, err := com.JsonEncode(cookieMap[length-1])
			if err != nil {
				return "", err
			}
			com.WriteFile(pathOscid, co)
		}
	}

	return string(b), err
}

func (this *Http) Get(urlstr string) (string, error) {
	home := GetHome()
	u, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}

	pathOscid := filepath.Join(home, ".osc", "oscid")
	jar := NewJar()

	// read cookie
	if com.FileExist(pathOscid) {
		json, _ := com.ReadFile(pathOscid)
		c := jar.ParseCookies(json)
		jar.SetCookies(u, []*http.Cookie{c})
	}

	// get
	client := http.Client{nil, nil, jar, 0}
	resp, err := client.Get(urlstr)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	// store cookie
	// cookieMap := jar.Cookies(u)
	// length := len(cookieMap)
	// // log.Greenln(length)
	// if length > 0 {
	// 	co, err := com.JsonEncode(cookieMap[length-1])
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	com.WriteFile(pathOscid, co)
	// }

	return string(b), err
}
