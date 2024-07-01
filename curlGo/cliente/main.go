package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <URL>")
		return
	}

	rawURL := os.Args[1]
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	fmt.Printf("Connecting to %s\n", rawURL)

	conn, err := net.Dial("tcp", parsedURL.Hostname()+":80")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	request := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	request += fmt.Sprintf("Host: %s\r\n", parsedURL.Hostname())
	request += "Accept: */*\r\n"
	request += "Connection: close\r\n\r\n"

	fmt.Fprintf(conn, request)

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Printf("Sending request GET %s HTTP/1.1\n", path)
	fmt.Printf("Host: %s\n", parsedURL.Hostname())
	fmt.Println("Accept: */*")

	fmt.Printf("Response from server: %s\n", statusLine)

}
