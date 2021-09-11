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

	for i, data := range Contents.Value {
		fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
	}

	//-------------------------------------
	joinedString := ""
	for _, data := range Contents.Value {
		//fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
		joinedString = joinedString + "<a href=" + data.Url + ">" + "<img src=" + data.Image.Url + "></a><br />" + "### " + data.Title + "<br />"
	}

	//-------------------------------------
	f, err := os.Open("README2.md")
    if err != nil{
        fmt.Println("error")
    }
    defer f.Close()

    // 一気に全部読み取り
    b, err := ioutil.ReadAll(f)
    // 出力
    //fmt.Println(string(b))
  //-------------------------------------
  str := []byte(string(b))
  assigned := regexp.MustCompile("<!--status-->\r\n\r\n(.*)\r\n\r\n<!--status-->")
  group := assigned.FindSubmatch(str)
  fmt.Println(string(group[1]))

  replacedMd := strings.Replace(string(b), string(group[1]), joinedString, 1)
  //fmt.Println(replacedMd)

  file, err := os.Create("README2.md")
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
