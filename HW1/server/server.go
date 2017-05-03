package main

import "fmt"
import "net"
import "time"
import "strconv"
import "strings"
import "bufio"

func main() {
    fmt.Println("Starting server")
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        // handle error
        fmt.Println("Error creating listener")
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting")
            // handle error
        }
        fmt.Println("Opened new connection")
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    for {
        // Get time that client sent
        t_org, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            return
        }
        // Record arrival time
        t_rec := time.Now().UnixNano()
        // Respond with client send time, arrival time, and server xmit time
        conn.Write([]byte(strings.TrimSpace(t_org) + "," + strconv.FormatInt(t_rec, 10) + ","+ strconv.FormatInt(time.Now().UnixNano(),10) + "\n"))
        // fmt.Println("Sent response to the client")
    }
}
