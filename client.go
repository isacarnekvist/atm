package main 

import (
    "net"
    "fmt"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8080")
    if err != nil {
        print("Error connecting")
        return
        print(conn)
    }
    // Send something!
    data := []byte{0x41, 0x42}
    conn.Write(data)
    resp := make([]byte, 8)
    conn.Read(resp)
    for _, d := range resp {
        fmt.Printf("%x", d)
    }
    print("\n")
    conn.Close()
}