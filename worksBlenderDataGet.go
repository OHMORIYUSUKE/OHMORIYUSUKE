package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Start!")

	baseUrl := os.Getenv("APIURL")
	authHeaderValue := os.Getenv("APIKEY")

	url := baseUrl + "?filters=type[contains]" + "Blender" + "&limit=10000"

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

	for i, data := range Contents.Value {
		fmt.Printf("index: %d,Id: %s, Title: %s,CreatedAt: %s\n", i, data.Id, data.Title, data.CreatedAt)
	}

	var titles []string
	var images []string
	var urls []string

	joinedString := "<table>"
	length := len(Contents.Value)
	fmt.Printf("%v\n", length)
	for i, data := range Contents.Value {
		if i < 2 {
			titles = append(titles, data.Title)
			images = append(images, data.Image.Url)
			urls = append(urls, data.Url)
			if i%2 == 1 {
				joinedString = joinedString + "<tr>" + "<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>" + "<th><a href=" + urls[1] + ">" + "<img src=" + images[1] + "></a></th>" + "</tr>" + "<tr>" + "<td>" + titles[0] + "</td>" + "<td>" + titles[1] + "</td>" + "</tr>"
				titles = nil
				images = nil
				urls = nil
			}
			if length%2 == 1 && length-1 == i {
				joinedString = joinedString + "<tr>" + "<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>" + "<th></th>" + "</tr>" + "<tr>" + "<td>" + titles[0] + "</td>" + "<td></td>" + "</tr>"
			}
		}
	}
	joinedString = joinedString + "</table>"
	joinedString = joinedString + "<details><summary><h3>もっと見る...</h3></summary><table>"
	for i, data := range Contents.Value {
		if i >= 2 {
			titles = append(titles, data.Title)
			images = append(images, data.Image.Url)
			urls = append(urls, data.Url)
			if i%2 == 1 {
				joinedString = joinedString + "<tr>" + "<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>" + "<th><a href=" + urls[1] + ">" + "<img src=" + images[1] + "></a></th>" + "</tr>" + "<tr>" + "<td>" + titles[0] + "</td>" + "<td>" + titles[1] + "</td>" + "</tr>"
				titles = nil
				images = nil
				urls = nil
			}
			if length%2 == 1 && length-1 == i {
				joinedString = joinedString + "<tr>" + "<th><a href=" + urls[0] + ">" + "<img src=" + images[0] + "></a></th>" + "<th></th>" + "</tr>" + "<tr>" + "<td>" + titles[0] + "</td>" + "<td></td>" + "</tr>"
			}
		}
	}
	joinedString = joinedString + "</table></details>"

	f, err := os.Open("README.md")
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)
	str := []byte(string(b))
	assigned := regexp.MustCompile("<!--works-Blender-->\n\n(.*)\n\n<!--works-Blender-->")
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
	Value []Contents `json:"contents"`
}

type Contents struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Image       Images `json:"image"`
}

type Images struct {
	Url    string `json:url`
	Height int    `json:height`
	Width  int    `json:width`
}
