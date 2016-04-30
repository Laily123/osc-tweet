package utils

import (
	"github.com/gogather/com"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	// "github.com/gogather/com/log"
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

func (this *Http) Post(urlstr string, parm url.Values) (string, error) {
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

	// post
	client := http.Client{nil, nil, jar, 0}
	resp, err := client.PostForm(urlstr, parm)

	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	// store cookie
	cookieMap := jar.Cookies(u)
	length := len(cookieMap)
	// log.Greenln(length)
	if length > 0 {
		co, err := com.JsonEncode(cookieMap[length-1])
		if err != nil {
			return "", err
		}

		com.WriteFile(pathOscid, co)
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
		json,_ := com.ReadFile(pathOscid)
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
	cookieMap := jar.Cookies(u)
	length := len(cookieMap)
	// log.Greenln(length)
	if length > 0 {
		co, err := com.JsonEncode(cookieMap[length-1])
		if err != nil {
			return "", err
		}

		com.WriteFile(pathOscid, co)
	}

	return string(b), err
}
