package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"strings"
	"os"
	"github.com/levigross/grequests"
	"time"
)

var wg sync.WaitGroup

func main()  {
	now := time.Now()
	url:="http://wxgirllive.com/forum-134-5.html"
	doc, err := goquery.NewDocument(url)
	if err!=nil{
		fmt.Errorf("下载错误：%#v", err)
		os.Exit(-1)
	}

	//fmt.Println(doc)
	//doc.Find("#normalthread_*").Each(func(i int, s *goquery.Selection) {
	//	src:= s.Find(".new .xst").Text()
	//	fmt.Println(string(src))
	//})

	//fmt.Println(doc.Find("#normalthread_8359170 .new .xst").Text())
	doc.Find("[id^='normalthread_']").Each(func(i int,s *goquery.Selection) {
		title := s.Find(".new .xst").Text()
		sub_link, exists_ := s.Find("tr > th > a.s.xst").Attr("href")
		fmt.Println(sub_link)
		fmt.Println("开始-->", title)
		link := "http://wxgirllive.com/" + sub_link
		fmt.Println(link)
		if exists_ {
			wg.Add(1)
			go func(link_ string) {
				defer wg.Done()
				doc, err := goquery.NewDocument(link_)
				if err != nil {
					fmt.Errorf("下载错误：%#v", err)
					os.Exit(-1)
				}

				doc.Find("[id^='postmessage_'] > img").Each(func(d int, s *goquery.Selection) {
					file_path, exist := s.Attr("file")
					if exist == false {
					} else {
						fmt.Println(file_path)
						dirname := "gallery_5/"+sub_link
						if _, err := os.Stat(dirname); err != nil{
							fmt.Printf("创建下载文件夹：%s\n", dirname)
							os.MkdirAll(dirname, 0666)
						}
						res, _ := grequests.Get(file_path, &grequests.RequestOptions{
							// 结构体可以对指定的类型给值，而不一定都赋值
							// 赋值Headers 请求头
							Headers:map[string]string{
								"Referer":"http://wxgirllive.com",
								"User-Agent":"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
							},
						})
						ff := strings.Split(file_path,"/")
						file_name := ff[len(ff)-1]
						filename := dirname+"/"+file_name
						res.DownloadToFile(filename)
					}

				})
			}(link)

		}
	})
	wg.Wait()
	fmt.Printf("下载任务完成，耗时:%#v\n", time.Now().Sub(now))

}
