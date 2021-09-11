package main

import (
	"os"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
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
    body, _ := io.ReadAll(resp.Body)

    // JSONを構造体にエンコード
    var Contents Root
    json.Unmarshal(body, &Contents)

    //fmt.Printf("%+v\n", Contents)

	//---

	for i, data := range Contents.Value {
		fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i,data.Id, data.Title, data.CreatedAt)
	}

}

type Root struct {
    Value []Contents `json:"contents"`
}

type Contents struct {
    Id  string  `json:"id"`
    CreatedAt  string `json:"createdAt"`
	Title string `json:"title"`
}
