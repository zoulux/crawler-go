package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	url := "https://www.baidu.com/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fetch url err:%v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("err status code:%v\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read body err:%v\n", err)
		return
	}

	fmt.Println("body:", string(body))

}
