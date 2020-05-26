package netutil

import (
	"fmt"
	"log"
	"os"
	"net"
	"time"
)

func demo() {
	if len(os.Args) <= 1 {
		fmt.Println("usage:go run clent.go need Your Content")
		return
	}

	log.Println("begin dial")

	conn, err := net.Dial("tcp", ":3306")

	if err != nil {
		log.Println("dail error", err)
		return
	}

	defer conn.Close()
	log.Println("dail ok")

	time.Sleep(time.Second * 2)

	data := os.Args[1]
	conn.Write([]byte(data))

	time.Sleep(time.Second * 5)
}
