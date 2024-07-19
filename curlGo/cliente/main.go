package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/url"
    "os"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
       fmt.Println("Usage: go run main.go [-v] [-X <method>] [-d <data>] [-H <header>] <URL>")
       return
    }

    var verbose bool
    var method string = "GET"
    var rawURL string
    var data string
    headers := make(map[string]string)

    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
       switch args[i] {
       case "-v":
          verbose = true
       case "-X":
          if i+1 < len(args) {
             method = args[i+1]
             i++
          } else {
             fmt.Println("Usage: go run main.go [-v] [-X <method>] [-d <data>] [-H <header>] <URL>")
             return
          }
       case "-d":
          if i+1 < len(args) {
             data = args[i+1]
             i++
          } else {
             fmt.Println("Usage: go run main.go [-v] [-X <method>] [-d <data>] [-H <header>] <URL>")
             return
          }
       case "-H":
          if i+1 < len(args) {
             header := strings.SplitN(args[i+1], ":", 2)
             if len(header) == 2 {
                headers[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])
             } else {
                fmt.Println("Invalid header format. Use: Header-Name: Header-Value")
                return
             }
             i++
          } else {
             fmt.Println("Usage: go run main.go [-v] [-X <method>] [-d <data>] [-H <header>] <URL>")
             return
          }
       default:
          rawURL = args[i]
       }
    }

    if rawURL == "" {
       fmt.Println("Usage: go run main.go [-v] [-X <method>] [-d <data>] [-H <header>] <URL>")
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
       err := conn.Close()
       if err != nil {
          return
       }
       fmt.Println("Conexão encerrada")
    }()

    path := parsedURL.Path
    if path == "" {
       path = "/"
    }

    var requestBuffer bytes.Buffer
    requestBuffer.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path))
    requestBuffer.WriteString(fmt.Sprintf("Host: %s\r\n", parsedURL.Hostname()))
    requestBuffer.WriteString("Accept: */*\r\n")
    requestBuffer.WriteString("Connection: close\r\n")

    for key, value := range headers {
       requestBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
    }

    if data != "" {
       if _, ok := headers["Content-Type"]; !ok {
          requestBuffer.WriteString("Content-Type: application/x-www-form-urlencoded\r\n")
       }
       requestBuffer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data)))
       requestBuffer.WriteString("\r\n")
       requestBuffer.WriteString(data)
    } else {
       requestBuffer.WriteString("\r\n")
    }

    request := requestBuffer.String()

    if verbose {
       fmt.Printf("> %s", request)
    }

    fmt.Fprintf(conn, request)

    reader := bufio.NewReader(conn)

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

    headersResp := make(map[string]string)
    for {
       line, err := reader.ReadString('\n')
       if err != nil {
          fmt.Println("Error reading response:", err)
          return
       }
       if line == "\r\n" {
          break
       }
       if verbose {
          fmt.Printf("< %s", line)
       }
       parts := strings.SplitN(line, ":", 2)
       if len(parts) == 2 {
          headersResp[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
       }
    }

    var responseBody strings.Builder
    for {
       line, err := reader.ReadString('\n')
       if err == io.EOF {
          break
       }
       if err != nil {
          fmt.Println("Error reading response body:", err)
          return
       }
       responseBody.WriteString(line)
       if verbose {
          fmt.Printf("< %s", line)
       }
    }

    var jsonResponse map[string]interface{}
    err = json.Unmarshal([]byte(responseBody.String()), &jsonResponse)
    if err == nil {
       jsonData, err := json.MarshalIndent(jsonResponse, "", "  ")
       if err == nil {
          fmt.Println(string(jsonData))
       }
    }

    fmt.Println()
