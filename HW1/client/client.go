package main

import "fmt"
import "net"
import "time"
import "strconv"
import "bufio"
import "strings"
import "os"
import "sort"
import "math"

type timing_endpoint struct {
    time int64
    kind int
}

type my_sample struct {
    time int64
    lambda int64
}

type timing_endpoints []timing_endpoint

func (slice timing_endpoints) Len() int {
    return len(slice)
}

func (slice timing_endpoints) Less(i, j int) bool {
    return slice[i].time < slice[j].time;
}

func (slice timing_endpoints) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func main() {
    

    f, _ := os.Create("./result.csv")
    defer f.Close()

    const num_samples = 100
    var samples [3*num_samples]timing_endpoint
    var measurements [num_samples]my_sample
    for i := 0; i < num_samples; i++ {
        conn, err := net.Dial("tcp", "100.81.2.162:8080")
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
        offsetf := 0.5 * float64((t_rec - t_org) + (t_xmt - t_dst))
        offset := int64(offsetf)
        //fmt.Println("RTT: " + strconv.FormatInt(rtt, 10))
        //  fmt.Println("Offset: " + strconv.FormatInt(offset, 10))
        samples[3*i] = timing_endpoint{offset - rtt/2, -1}
        samples[3*i + 1] = timing_endpoint{offset, 0}
        samples[3*i + 2] = timing_endpoint{offset + rtt/2, 1}

        measurements[i] = my_sample{offset, rtt/2}

        
        f.Write([]byte(strconv.FormatInt(rtt, 10) +"," + strconv.FormatInt(offset, 10) + "\n"))

    }

    l, u := selection_alg(samples[:])
    fmt.Println("Lower bound: " + strconv.FormatInt(l,10))
    fmt.Println("Upper Bound: " + strconv.FormatInt(u,10))
    clustered_samps := cluster_algorithm(measurements[:],l,u)
    final_offset := combining_algorithm(clustered_samps[:])
    fmt.Println("Final Offset: ", final_offset)
    
}

func selection_alg(samples timing_endpoints) (int64, int64) {
    var c, d, f int
    // Select the low point, midpoint, and high point of these intervals
    // Sort these values in a list from lowest to highest.
    sort.Sort(samples)
    f = 0
    n := len(samples)
    m := n/3
    var l, u int64

    for f < m/2 {
        // Set the number of midpoints d = 0. Set c = 0
        d = 0
        c = 0
        // Scan from lowest endpoint to highest.
        for i := 0; i < n; i++ {
            // Add one to c for every low point, subtract one for every high point
            c -= samples[i].kind
            // Add one to d for every midpoint
            if samples[i].kind == 0 {
                d++
            }
            // if c >= m - f, stop and set l = current low point
            if c >= (m - f) {
                fmt.Println("Setting l to ",samples[i].time)
                // Set l = current low point
                l = samples[i].time
                break
            }
        }

        // Set c = 0. Scan from highest endpoint to lowest.
        c = 0
        for i := n-1; i > -1; i-- {
            // Add 1 to c for every high point, subtract 1 for every low point
            c += samples[i].kind
            // Add 1 to d for every midpoint
            if samples[i].kind == 0 {
                d++
            }
            // If c >= m - f, stop and set u = current high point
            fmt.Println("C: ",c)
            if c >= (m - f) {
                // Set u = current high position
                u = samples[i].time
                break
            }
        }

        // Is d <= f and l < u?
        if (d <= f && l < u) {
            // Yes => SUCCESS
            // intersection interval is [l,u]
            fmt.Println("f:",f)
            return l, u
            break
        } else {
            // Add 1 to f
            f = f + 1
            // if f >= m/2, then failure condition
            if f >= m/2 {
                fmt.Println("Failure condition :(")
                fmt.Println("l:",l)
                fmt.Println("u:",u)
                fmt.Println("f:",f)
                fmt.Println("m:",m)
                fmt.Println("m/2: ",m/2)
                return int64(0), int64(0)
            }
        }
    }
    return 0, 0
}

func cluster_algorithm(samples []my_sample, l int64, u int64) []my_sample {
    // Remove any samples that are not in the correct interval
    for i := 0; i < len(samples); i++ {
        if samples[i].time < l || samples[i].time > u {
            // this sample is outside the correct interval
            samples[i] = samples[len(samples) - 1]
            samples = samples[:len(samples)-1]
            fmt.Println("Removed a bad sample")
        }
    }

    min_samples := 40
    m := len(samples)
    phi := make([]int64,m,m)
    for m > min_samples {
        for i := 0; i < m; i++ {
            for j:= 0; j < m; j++ {
                phi[i] += int64(math.Pow(float64(samples[j].time - samples[i].time),2))
            }
            phi[i] = int64(math.Sqrt(1/float64((m - 1)) * float64(phi[i])))
        }
        // Now, remove the sample with the largest phi
        var max_phi int64
        var max_idx int
        for i := 0; i < m; i++ {
            if phi[i] > max_phi {
                max_phi = phi[i]
                max_idx = i
            }
        }
        samples[max_idx] = samples[m - 1]
        samples = samples[:m-1]
        m = len(samples)
        fmt.Println("m: ",m)
    }
    return samples
}

func combining_algorithm(samples []my_sample) int64 {
    var y, z float64
    for i := 0; i < len(samples); i++ {
        y += 1/float64(samples[i].lambda)
        z += float64(samples[i].time)/float64(samples[i].lambda)
    }
    return int64(z/y)
}