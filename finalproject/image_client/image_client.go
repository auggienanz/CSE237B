package main

import "fmt"
import "net/http"
import "io/ioutil"
import "time"
import "strconv"
import "golang.org/x/exp/mmap"
import "os"
import "strings"
func main() {
    // Read time offset from file
    reader, err := mmap.Open("../offset.txt")
    if err != nil {
        fmt.Printf("Error opening mmap'd file")
        os.Exit(1)
    }
    data := make([]byte, 50) 
    n,_ := reader.ReadAt(data,0)
    if n == 0 {
        fmt.Printf("Read no data from file")
        os.Exit(1)
    }
    times := strings.Split(string(data[:n]), ",")
    // Time when synchronization occurred
    t_sync,_ := strconv.ParseInt(times[0],10,64)
    // Offset between client and server clock
    t_offset,_ := strconv.ParseInt(strings.TrimSpace(times[1]),10,64)

    // First, check if the synchronization result is from last 5 minutes
    if (time.Now().UnixNano() - t_sync > 300000000000) {
        fmt.Println("Synchronization is too old. Please rerun sync_client")
        os.Exit(1)
    }
    // Have a valid time synchronization!

    // Create a file to write data to
    f, _ := os.Create("../latency.csv")
    defer f.Close()

    t_start := time.Now().UnixNano()
    // Collect data for 5 minutes
    for time.Now().UnixNano() - t_start < 300000000000 {
        resp, err := http.Get("http://192.168.1.95:8081/red")
        if err != nil {
            // handle error
        }
        defer resp.Body.Close()

        // Record arrival time
        t_rcv := time.Now().UnixNano()

        // Extract transmit time from response
        body, _ := ioutil.ReadAll(resp.Body)        
        t_send_s := string(body[50:70])
        t_send, _ := strconv.ParseInt(t_send_s,10,64)

        latency := t_rcv + t_offset - t_send;

        // Write the results to a file
        f.Write([]byte(strconv.FormatInt(t_rcv, 10) +"," + strconv.FormatInt(latency, 10) + "\n"))
    }}
}