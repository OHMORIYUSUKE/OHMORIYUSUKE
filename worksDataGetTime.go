package main

import (
    "time"
	"os"
	"log"
	"regexp"
	"strings"
	"fmt"
	"io/ioutil"
)

func main(){
	now := time.Now()

    nowUTC := now.UTC() 

    jst := time.FixedZone("Asia/Tokyo", 9*60*60)

    nowJST := nowUTC.In(jst)

	nowJST.Format("2006-01-02 15:04:05")

	//-----------------------
	f, err := os.Open("README.md")
    if err != nil{
        fmt.Println("error")
    }
    defer f.Close()

    b, err := ioutil.ReadAll(f)

	str := []byte(string(b))
	assigned := regexp.MustCompile("<!--works-GetDtataTime-->\n\n(.*)\n\n<!--works-GetDtataTime-->")
	group := assigned.FindSubmatch(str)

	replacedMd := strings.Replace(string(b), string(group[1]), "最終更新 : " + nowJST.Format("2006-01-02 15:04:05"), 1)

	file, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err) 
	}
	defer file.Close()

	file.Write(([]byte)(replacedMd))
}