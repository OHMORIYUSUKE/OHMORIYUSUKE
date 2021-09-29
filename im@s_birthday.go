package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	now := time.Now()

	nowUTC := now.UTC()

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	nowJST := nowUTC.In(jst)

	nowJST.Format("2006-01-02 15:04:05")

	//---

	//t := time.Now()
	month := int(nowJST.Month())
	day := nowJST.Day()

	var dayStr string
	dayStr = strconv.Itoa(day)

	var monthStr string
	monthStr = strconv.Itoa(month)

	url := "https://sparql.crssnky.xyz/spql/imas/query?output=json&force-accept=text%2Fplain&query=PREFIX%20schema%3A%20%3Chttp%3A%2F%2Fschema.org%2F%3E%0APREFIX%20rdfs%3A%20%20%3Chttp%3A%2F%2Fwww.w3.org%2F2000%2F01%2Frdf-schema%23%3E%0A%0ASELECT%20(sample(%3Fo)%20as%20%3Fdate)%20(sample(%3Fn)%20as%20%3Fname)%0AWHERE%20%7B%20%0A%20%20%3Fs%20schema%3AbirthDate%20%3Fo%3B%0A%20%20rdfs%3Alabel%20%3Fn%3B%0A%20%20FILTER(regex(str(%3Fo)%2C%20%22" + monthStr + "-" + dayStr + "%22)).%0A%7D%0Agroup%20by(%3Fn)%0Aorder%20by(%3Fname)"

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	client := new(http.Client)
	resp, err := client.Do(req)

	// URLがnilだったり、Timeoutが発生した場合にエラーを返す模様。
	// サーバーからのレスポンスとなる 401 Unauthroized Error などはResponseをチェックする。
	// サーバーとの疎通が開始する前の動作のよう。
	if err != nil {
		fmt.Println("Error Request:", err)
		return
	}
	// resp.Bodyはクローズすること。クローズしないとTCPコネクションを開きっぱなしになる。
	defer resp.Body.Close()

	// 200 OK 以外の場合はエラーメッセージを表示して終了
	if resp.StatusCode != 200 {
		fmt.Println("Error Response:", resp.Status)
		return
	}

	// Response Body を読み取り
	body, _ := ioutil.ReadAll(resp.Body)

	// JSONを構造体にエンコード
	var Contents Root
	json.Unmarshal(body, &Contents)

	fmt.Printf("%+v\n", Contents)

	joinedString := "<ul>"
	for i, data := range Contents.Value.Binding {
		fmt.Printf("index: %d,Name: %s\n", i, data.Name.Value)
		joinedString = joinedString + "<li><h2><a href="+ "https://www.google.com/search?q=" + data.Name.Value + "&tbm=isch&oq=" + data.Name.Value + "&sclient=img" +">" + data.Name.Value + "</a></h2></li>"
	}

	if len(Contents.Value.Binding) == 0 {
		joinedString = joinedString + "<li><h2>" + nowJST.Format("01月02日") +"誕生日の人おめでとう!" + "</h2></li>"
	}

	if month == 11 && day == 8 {
		joinedString = joinedString + "<li><h2>" + "大森裕介" + "</h2></li>"
	}
	joinedString = joinedString + "</ul><!--" + nowJST.Format("2006-01-02 15:04:05") + "-->"

	f, err := os.Open("README.md")
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)
	
	str := []byte(string(b))
	assigned := regexp.MustCompile("<!--imats-birthday-->\n\n(.*)\n\n<!--imats-birthday-->")
	group := assigned.FindSubmatch(str)

	replacedMd := strings.Replace(string(b), string(group[1]), joinedString, 1)
	
	file, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err) 
	}
	defer file.Close()

	file.Write(([]byte)(replacedMd))
}

type Root struct {
	Value Contents `json:"results"`
}

type Contents struct {
	Binding []Bindings `json:"bindings"`
}

type Bindings struct {
	Name Names `json:"name"`
}

type Names struct {
	Value string `json:value`
}
