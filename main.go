package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// 获取公网 IPv6 地址
func getPublicIPv6() (string, error) {
	resp, err := http.Get("http://6.ipw.cn")
	if err != nil {
		return "", fmt.Errorf("error fetching IPv6 address: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应的 Content-Type，确保是文本
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		return "", fmt.Errorf("unexpected content type: %s", contentType)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// 去除多余的空白符号
	ip := strings.TrimSpace(string(body))
	return ip, nil
}

// 发送邮件
func sendEmail(subject, body string) error {
	host := "smtp.xx.com"        //根据自己需求改
	port := 465                  //根据自己需求改
	userName := "xxxxxx0@xx.com" //根据自己需求改
	password := "xxxxx"          ////根据自己需求改
	m := gomail.NewMessage()
	m.SetHeader("From", userName)
	m.SetHeader("To", "xxxxxx1@xx.com") //根据自己需求改
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, userName, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}

func initializeIPv6() string {
	// 示例内网 IPv6 地址 (ULA, Unique Local Address)，可以自定义
	// 例如，fd00:: 开头为内网保留地址
	return "fd00::1"
}

func main() {
	var initialIPv6 string

	initialIPv6 = initializeIPv6()
	log.Print("Initializing IPv6 address...", initialIPv6)
	// 获取初始的公网IPv6地址

	// 每隔一定时间检查公网IPv6地址变化
	ticker := time.NewTicker(2 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentIPv6, err := getPublicIPv6()
			if err != nil {
				log.Println("Error fetching current public IPv6 address:", err)
				continue
			}

			log.Println("Current Public IPv6 address:", currentIPv6)

			// 比较当前地址和初始地址是否不同
			if currentIPv6 != initialIPv6 {
				log.Println("Public IPv6 address has changed!")
				body := fmt.Sprintf("The public IPv6 address has changed to: %s", currentIPv6)

				// 公网IPv6地址发生变化，发送邮件
				err := sendEmail("Public IPv6 Address Changed", body)
				if err != nil {
					log.Println("Error sending email:", err)
				}

				// 更新初始IPv6地址
				initialIPv6 = currentIPv6
			}
		}
	}
}
