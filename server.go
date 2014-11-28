package main

import (
    "net"
    "fmt"
)

// Förslagsvis ha en map där ett namn mappar till denna struct:
type User struct {
    password    string 
    balance     int
    temp_code   []int      // Hur fungerar allokeringen här?
}

func main() {
    print("Starting server\n")
    ln, err := net.Listen("tcp", ":8080")                       // Lyssna på port 8080
    if err != nil {
        print(ln)
    }
    for {
        conn, err := ln.Accept()                                // Vänta tills en anslutning begärs
        if err != nil {
            print("Error accepting\n")
        }
        go handleConnection(conn)                               // Starta en tråd/gorutin för varje anslutning
    }
}

func handleConnection(c net.Conn) {
    print("Connection started\n")
    data := make([]byte, 10)
    n, err := c.Read(data)                                      // Läs bytes till variabeln "data"
    if err != nil {
        print(err)
    }
    fmt.Printf("n is %d\n", n)                                  // Vet inte vad n är än
    for i, d := range data {
        fmt.Printf("data at %d is: %c\n", i, d)
    }
    // Vi kan svara med data efter att någon kontaktat
    response := []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}  // Det finns ett paket "bytes" för att hantera
    c.Write(response)                                           // bytesträngar som nog blir smidigare
    c.Close()
    print("Connection closed\n")
}