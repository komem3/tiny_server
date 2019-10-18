package tiny_http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Response struct {
	Version       string
	StatusCode    int
	StatusMessage string
	Headers       map[string]string
	Body          string
}

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func HandleConnection(con net.Conn) {
	defer con.Close()

	request, err := parseRequest(con)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	response := Response{
		Version: "Http/1.1",
		Headers: map[string]string{},
	}

	ServeHTTP(&response, &request)

	if err := parseResponse(con, response); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("[%s]: %s %d %s\n", request.Method, request.Path, response.StatusCode, response.StatusMessage)
}

func ServeHTTP(w *Response, r *Request) {
	w.StatusCode = 200
	w.StatusMessage = "OK"
	w.Headers["Server"] = "tiny server"
	w.Body = "Hello World"
}

func parseResponse(con net.Conn, response Response) error {
	writeString := fmt.Sprintf("%s %d %s\r\n", response.Version, response.StatusCode, response.StatusMessage)
	for k, v := range response.Headers {
		writeString += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	writeString += fmt.Sprintf("\r\n%s", response.Body)

	_, err := con.Write([]byte(writeString))
	return err
}

func parseRequest(con net.Conn) (request Request, err error) {
	b := bufio.NewReader(con)

	line, _, err := b.ReadLine()
	if err != nil {
		return
	}

	// parse method route version
	contents := strings.Split(string(line), " ")
	if len(contents) != 3 {
		err = fmt.Errorf("Request is not http protocol\n")
		return
	}
	request = Request{
		Method:  contents[0],
		Path:    contents[1],
		Version: contents[2],
		Headers: map[string]string{},
	}

	// parse header
	for {
		line, _, err = b.ReadLine()
		if err != nil {
			return
		}
		if len(line) < 2 || (line[0] == '\r' && line[1] == '\n') {
			break
		}
		keyValue := strings.Split(string(line), ":")
		if len(keyValue) < 2 {
			err = fmt.Errorf("Header is wrong format\n")
			return
		}
		request.Headers[keyValue[0]] = strings.Join(keyValue[1:], "")
	}

	for err != nil || len(line) != 0 {
		line, _, err = b.ReadLine()
		if len(line) == 0 || (err != nil && err != io.EOF) {
			return
		}
		request.Body = fmt.Sprintf("%s%s", request.Body, string(line))
	}

	return
}
