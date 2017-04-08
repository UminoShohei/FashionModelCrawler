package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	readCSV()

}
func readCSV() {
	var fp *os.File
	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		var err error
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない！
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		j := 1
		for i := 1; i < 20; i++ {
			doc, err := goquery.NewDocumentFromResponse(getImages(record[0], record[1], i))
			fmt.Println(doc.Url)
			if err != nil {
				fmt.Print("url scarapping failed")
			}
			doc.Find("img").Each(func(_ int, s *goquery.Selection) {
				url, _ := s.Attr("src")
				saveImage(record[0], url, j)
				j++
				fmt.Println(url)
			})
		}
	}
}
func getImages(name, pid string, n int) *http.Response {
	url := fmt.Sprintf("http://gensun.org/?mode=ajax&q=%s&pid=%s&page=%d&size=100&safe=on&sort=", name, pid, n)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Host", "gensun.org")
	req.Header.Set("Referer", fmt.Sprintf("http://gensun.org/?q=%s", name))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Host", "gensun.org")
	client := new(http.Client)
	resp, _ := client.Do(req)
	return resp
}

func saveImage(name, url string, i int) {
	if strings.Compare(url, "/img/loading.gif") == 0 {
		return
	}
	defer func() {
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
	imageName := fmt.Sprintf("/Users/uminoshohei/Project/go/src/github.com/UminoShohei/get_model_images/hoge/%s_%d.jpg", name, i)
	defer response.Body.Close()
	file, err := os.Create(imageName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	io.Copy(file, response.Body)
}
