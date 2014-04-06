/*
Copyright 2014 Gavin Bong.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific
language governing permissions and limitations under the
License.
*/

package progress

import (
    "net/http"
    "runtime"
    "fmt"
    "io"
    "time"
    "log"
    "bytes"
    "net"
    "strings"
    "strconv"

    "github.com/javouhey/seneca/util"
)

const (
    FFMPEG_USERAGENT = "Lavf"
)

/* 
 Typical fields submitted every seconds by ffmpeg.

   frame=0
   fps=0.0
   stream_0_0_q=0.0
   total_size=N/A
   out_time_ms=0
   out_time=00:00:00.000000
   dup_frames=0
   drop_frames=0
   progress=continue / end 
 */
type Status struct {
    frame int32
    drop_frames int32
    progress string
}

func (s *Status) parse(httpBody string) {
    lines := strings.Split(httpBody, "\n")
    for _, keyvaluepair := range lines {
        parts := strings.Split(keyvaluepair, "=")
        if len(parts) == 2 {
            switch parts[0] {
            case "frame":
                if r, err := strconv.ParseInt(parts[1], 10, 32); err == nil {
                    s.frame = int32(r)
                }

            case "progress":
                s.progress = parts[1]

            case "drop_frames":
                if r, err := strconv.ParseInt(parts[1], 10, 32); err == nil {
                    s.drop_frames = int32(r)
                }
            }
        }
    }
    return
}

type MyHandler struct {
    pings chan<- Status
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ua := r.Header.Get("User-Agent")
    //fmt.Println("ua:", ua)

    if !strings.HasPrefix(ua, FFMPEG_USERAGENT) {
        w.WriteHeader(http.StatusForbidden)
        w.Write([]byte("for internal use only\n\n"))
        w.(http.Flusher).Flush()
        return
    }

    reader := r.Body
    defer reader.Close()

    var (
        buffer bytes.Buffer
        n int 
        err error
    )
    for {
        n = 0; err = nil
        var tmp []byte

        tmp = make([]byte, 256)
        n, err = io.ReadAtLeast(reader, tmp, 132)
        if err != nil {
            switch err {
            case io.EOF:
                goto finish
            case io.ErrUnexpectedEOF:
                goto next
            default:
                continue
            }
        }
    next:
        if n>0 {
            tmp = tmp[0:n]
            buffer.Write(tmp)

            status := Status{}
            status.parse(buffer.String())
            //log.Printf("%#v\n", status)
            h.pings <- status
        }
        buffer.Reset()
    }
finish:
    w.WriteHeader(http.StatusNoContent)
    w.(http.Flusher).Flush()
}

// goroutine responsible for printing progress ticks
func StatusLogger(q <-chan Status) {
    for {
        stat, ok := <- q
        if !ok {
            break
        }

        if stat.progress == "continue" {
            switch {
            case stat.frame == 0: fmt.Printf(".")
            default: fmt.Printf(" %d", stat.frame)
            }
        } else {
            fmt.Printf(" Completed\n")
        }
        runtime.Gosched()
    }
}

// goroutine responsible for starting the webserver
func Progress(l net.Listener, q chan<- Status, port int) {
    //httpPort := strconv.Itoa(port)

    s := &http.Server{
      Addr: util.ToPort(port),
      Handler: MyHandler{q},
      //ReadTimeout: 60 * 60 * time.Second,
      WriteTimeout: 40 * time.Second,
      MaxHeaderBytes: 1 << 20,
    }
    //log.Printf("HTTP server listening on port %s\n", httpPort)
    log.Println(s.Serve(l))
}
