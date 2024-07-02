package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [-v] [-X <method>] <URL>")
		return
	}

	var verbose bool
	var method string = "GET"
	var rawURL string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v":
			verbose = true
		case "-X":
			if i+1 < len(args) {
				method = args[i+1]
				i++ // skip the next argument as it is the method
			} else {
				fmt.Println("Usage: go run main.go [-v] [-X <method>] <URL>")
				return
			}
		default:
			rawURL = args[i]
		}
	}

	if rawURL == "" {
		fmt.Println("Usage: go run main.go [-v] [-X <method>] <URL>")
		return
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	fmt.Printf("Connecting to %s\n", rawURL)

	conn, err := net.Dial("tcp", parsedURL.Hostname()+":"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer func() {
		conn.Close()
		fmt.Println("Conexão encerrada")
	}()

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	request := fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path)
	request += fmt.Sprintf("Host: %s\r\n", parsedURL.Hostname())
	request += "Accept: */*\r\n"
	request += "Connection: close\r\n\r\n"

	if verbose {
		fmt.Printf("> %s", request)
	}

	fmt.Fprintf(conn, request)

	reader := bufio.NewReader(conn)

	// Read the status line
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if verbose {
		fmt.Printf("< %s", statusLine)
	} else {
		fmt.Printf("Response from server: %s", statusLine)
	}

	// Read the headers
	headers := ""
	if verbose {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			if line == "\r\n" {
				break
			}
			headers += line
			fmt.Printf("< %s", line)
		}
	}

	// Read the body
	if method == "DELETE" && verbose {
		body := ""
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			body += line
			fmt.Printf("< %s", line)
		}
	}

	fmt.Println() // Print a new line for better formatting
}
