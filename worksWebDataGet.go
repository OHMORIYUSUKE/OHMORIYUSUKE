package main

import (
	"os"
    "encoding/json"
    "fmt"
    "net/http"
	"regexp"
	"strings"
	"log"
	"io/ioutil"
    "time"
)

func main() {
    fmt.Println("Start!")

	baseUrl := os.Getenv("APIURL")
	authHeaderValue := os.Getenv("APIKEY")

	url := baseUrl + "?filters=type[contains]" + "Web" + "&limit=10000"

    authHeaderName := "X-API-KEY"

    req, _ := http.NewRequest(http.MethodGet, url, nil)
    req.Header.Set(authHeaderName, authHeaderValue)

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

	//---
    now := time.Now()
    fmt.Println(now.Format(time.RFC3339))

    nowUTC := now.UTC() 
    fmt.Println(nowUTC.Format(time.RFC3339))

    jst := time.FixedZone("Asia/Tokyo", 9*60*60)

    nowJST := nowUTC.In(jst)                        
    fmt.Println(nowJST.Format(time.RFC3339))
    //---

	for i, data := range Contents.Value {
		fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
	}

	//-------------------------------------
	var titles[] string 
    var images[] string 
    var urls[] string 
    //---
    length := len(Contents.Value)

	joinedString := "<table>"
	// for _, data := range Contents.Value {
	// 	//fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
	// 	joinedString = joinedString + "<a href=" + data.Url + ">" + "<img src=" + data.Image.Url + "></a><br />" + "### " + data.Title + "<br />"
	// }
    for i, data := range Contents.Value {
		//fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
        titles = append(titles, data.Title)
        images = append(images, data.Image.Url)
        urls = append(urls, data.Url)
        if i % 2 == 1 {
            joinedString = joinedString + "<tr>"+"<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>"+"<th><a href=" + urls[1] + ">" + "<img src=" + images[1] + "></a></th>"+"</tr>" + "<tr>"+"<td>" + titles[0] + "</td>"+"<td>" + titles[1] + "</td>"+"</tr>"
            titles = nil
            images = nil
            urls = nil
        }
        if length % 2 == 1 && length - 1 == i {
            joinedString = joinedString + "<tr>"+"<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>"+"<th></th>"+"</tr>" + "<tr>"+"<td>" + titles[0] + "</td>"+"<td></td>"+"</tr>"
        }
	}
    joinedString = joinedString + "</table>"
	//-------------------------------------
	f, err := os.Open("README.md")
    if err != nil{
        fmt.Println("error")
    }
    defer f.Close()

    // 一気に全部読み取り
    b, err := ioutil.ReadAll(f)
    // 出力
    //fmt.Println(string(b))
  //-------------------------------------
  // TODO 改行コードを変える
  str := []byte(string(b))
  assigned := regexp.MustCompile("<!--works-Web-->\n\n(.*)\n\n<!--works-Web-->")
  group := assigned.FindSubmatch(str)
  fmt.Println(string(group[1]))

  replacedMd := strings.Replace(string(b), string(group[1]), joinedString, 1)
  //fmt.Println(replacedMd)

  file, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)  //ファイルが開けなかったときエラー出力
	}
	defer file.Close()

	file.Write(([]byte)(replacedMd))
	//-----------------------------------

}

type Root struct {
    Value []Contents `json:"contents"`
}

type Contents struct {
    Id  string  `json:"id"`
    CreatedAt  string `json:"createdAt"`
	Title string `json:"title"`
	Description string `json:"description"`
	Url string `json:"url"`
	Image Images `json:"image"`
}

type Images struct {
	Url string `json:url`
	Height int `json:height`
	Width int  `json:width`
}
