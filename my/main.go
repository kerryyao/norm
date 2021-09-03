/*
Author: suguo.yao(ysg@myschools.me)
Description: 用于从myschools.me中下载标准组件
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const GITEA = "https://myschools.me"
const PATH = "/suguo/norm/src/branch/master"

var path string

func init() {
	flag.StringVar(&path, "dl", "", "下载组件名")
}

func main() {
	flag.Parse()
	if path == "" {
		panic("下载组件名指定")
	}
	var client http.Client
	var wg sync.WaitGroup
	start := time.Now()
	dl(client, path, &wg)
	wg.Wait()
	fmt.Printf("total time: %.2f s\n", float64(time.Since(start))/float64(time.Second))
}

// get all file link and download it
func dl(client http.Client, path string, wg *sync.WaitGroup) {
	if !isExist(path) {
		os.MkdirAll(path, 0775)
	}

	url := fmt.Sprintf("%s%s/%s", GITEA, PATH, path)
	html, err := getHtml(client, url)
	if err != nil {
		fmt.Printf("get html error: %s", err.Error())
		return
	}
	urlPattern := regexp.MustCompile(fmt.Sprintf(`%s/%s/\S*go`, PATH, path))
	links := urlPattern.FindAllSubmatch(html, -1)
	for _, link := range links {
		tmp := strings.Split(string(link[0]), "/")
		filename := tmp[len(tmp)-1]
		wg.Add(1)
		go downloadFile(client, path, filename, wg)

	}
}

// download file
func downloadFile(client http.Client, path, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("start to download: ", filename)
	fileURL := fmt.Sprintf("%s/suguo/norm/raw/branch/master/%s/%s", GITEA, path, filename)
	resp, err := client.Get(fileURL)
	if err != nil {
		fmt.Printf("download file %s failed due to: %s\n", filename, err.Error())
		return
	}
	defer resp.Body.Close()
	var buff [1024]byte
	// 创建文件
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		fmt.Printf("create file: %s error\n", filename)
		return
	}
	defer file.Close()
	// 写入文件
	for {
		n, err := resp.Body.Read(buff[:])
		if err != nil {
			if err == io.EOF {
				file.Write(buff[:n])
				break
			}
			fmt.Println("error: ", err)
			os.Remove(filepath.Join(path, filename))
			return
		}
		file.Write(buff[:n])
	}
	fmt.Println("finish download:", filename)
}

// get html source
func getHtml(client http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// if file or directory exits
func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
