package main

import "fmt"
import "net"
import "time"
import "strconv"
import "bufio"
import "strings"
import "os"

func main() {
    

    f, _ := os.Create("./result.csv")
    defer f.Close()

    for i := 0; i < 1000; i++ {
        conn, err := net.Dial("tcp", "192.168.137.140:8080")
        if err != nil {
            // handle error
        }   
        // Send current time
        fmt.Fprintf(conn, strconv.FormatInt(time.Now().UnixNano(),10)+"\n")
        // Read response from server
        status, _ := bufio.NewReader(conn).ReadString('\n')
        // Record response arrival time
        t_dst := time.Now().UnixNano()

        // Parse server response
        result := strings.Split(status, ",")
        t_org,_ := strconv.ParseInt(result[0],10,64)
        t_rec,_ := strconv.ParseInt(result[1],10,64)
        t_xmt,_ := strconv.ParseInt(strings.TrimSpace(result[2]),10,64)

        // Calculate round trip time and offset
        rtt := (t_dst - t_org) - (t_xmt - t_rec)
        offset := 0.5 * float64((t_rec - t_org) + (t_xmt - t_dst))
        fmt.Println("RTT: " + strconv.FormatInt(rtt, 10))
        fmt.Println("Offset: " + strconv.FormatFloat(offset, 'f',-1,64))

        
        f.Write([]byte(strconv.FormatInt(rtt, 10) +"," + strconv.FormatFloat(offset, 'f',-1,64) + "\n"))

    }
    
}