package main 

import (
    "net"
    "fmt"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8080")          // Anslut till ip:port
    if err != nil {
        print("Error connecting")
        return
        print(conn)
    }

    // Skicka data
    data := []byte{0x41, 0x42}
    conn.Write(data)

    // Ta emot!
    resp := make([]byte, 8)
    conn.Read(resp)

    // Skriv ut
    for _, d := range resp {
        fmt.Printf("%x", d)
    }
    print("\n")

    // Stäng
    conn.Close()
}