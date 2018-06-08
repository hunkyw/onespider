package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"golang.org/x/text/encoding"
	"golang.org/x/net/html/charset"
	"bufio"
	"regexp"
)

func main() {
	//获取页面
	resp , err := http.Get("http://www.zhenai.com/zhenghun")
	if err != nil  {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}

    //转为UTF8 字符
	e := determineEncoding(resp.Body)
	utf8Reader := transform.NewReader(resp.Body , e.NewDecoder())
		all , err := ioutil.ReadAll(utf8Reader)
		if err != nil {
			panic(err)
		}
		//
		printCityList(all)

}

//正则表达提取字符 获得城市列表
func printCityList(contents []byte) {
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(contents,-1)
	for _,m := range matches {
		fmt.Printf("City : %s , URL : %s \n",m[2], m[1])
	}
	fmt.Printf("Matches found : %d \n", len(matches))
}

//发现编码格式并转为utf-8
func determineEncoding(r io.Reader) encoding.Encoding {
	bytes , err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e , _ , _ := charset.DetermineEncoding(bytes , "")
	return e
}
