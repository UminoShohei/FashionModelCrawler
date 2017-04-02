package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	j := 1
	for i := 1; i < 20; i++ {
		doc, err := goquery.NewDocumentFromResponse(getImages(i))
		fmt.Println(doc.Url)
		if err != nil {
			fmt.Print("url scarapping failed")
		}
		doc.Find("img").Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("src")
			saveImage(url, j)
			j++
			fmt.Println(url)
		})
	}
}

func getImages(n int) *http.Response {
	url := fmt.Sprintf("http://gensun.org/?mode=ajax&q=小松菜奈&pid=2282199&page=%d&size=100&safe=on&sort=", n)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Host", "gensun.org")
	req.Header.Set("Referer", "http://gensun.org/wid/2282199") //ここ引数で取れるようにするべきかな
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Host", "gensun.org")
	client := new(http.Client)
	resp, _ := client.Do(req)
	return resp
}

func saveImage(url string, i int) {
	fmt.Println(url)

	if strings.Compare(url, "/img/loading.gif") == 0 {
		fmt.Println("AAAAAAAAAAAA")
		return
	}
	defer func() {
		fmt.Println("End")
		err := recover()
		if err != nil {
			fmt.Println("Recover!:", err)
		}
	}()
	response, err := http.Get(url)
	fmt.Println(i)
	if err != nil {
		fmt.Println(err)
	}
	imageName := fmt.Sprintf("%d.jpg", i)
	defer response.Body.Close()
	file, err := os.Create(imageName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	io.Copy(file, response.Body)
}
