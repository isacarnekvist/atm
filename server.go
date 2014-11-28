package main

import (
    "net"
    "fmt"
)

func main() {
    print("Starting server\n")
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        print(ln)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            print("Error accepting\n")
        }
        go handleConnection(conn)
    }
}

func handleConnection(c net.Conn) {
    print("Connection started\n")
    data := make([]byte, 10)
    n, err := c.Read(data)
    if err != nil {
        print(err)
    }
    fmt.Printf("n is %d\n", n)
    for i, d := range data {
        fmt.Printf("data at %d is: %c\n", i, d)
    }
    response := []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}
    c.Write(response)
    c.Close()
    print("Connection closed\n")
}