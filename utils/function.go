package utils

import (
	"fmt"
    "crypto/sha1"
    "io"
	"github.com/gogather/com/log"
)

func GetHome() string {
	home, err := Home()
	if err!=nil {
		log.Fatalln("Can NOT find user path!")
	}
	return home
}

//对字符串进行SHA1哈希
func SHA1(data string) string {
    t := sha1.New();
    io.WriteString(t,data);
    return fmt.Sprintf("%x",t.Sum(nil));
}
