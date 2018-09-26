package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/zeny-io/mboxparser"
    "bytes"
    "os"
    "strconv"
    "flag"
    "strings"
    "io/ioutil"
)

/* global flags */
var port = flag.Int("port", 9090, "Port to listen on")
var host = flag.String("host", "0.0.0.0", "IP to bind to")
var mailroot = flag.String("mailroot", "/var/mail", "Mail root")

func mboxServer(w http.ResponseWriter, r *http.Request) {
    fmt.Println("path", r.URL.Path)
    /* tmp debug */
    w.Header().Set("Access-Control-Allow-Origin", "*")

    if strings.Contains(r.URL.Path, "../") {
        http.Error(w, "{\"error\":\"invalid request\"}", 400)
        return
    }

    if r.URL.Path == "/" {
        files, err := ioutil.ReadDir(*mailroot)
        if err != nil {
            fmt.Fprint(w, "[]")
            return
        }
        fmt.Fprint(w, "[")
        for idx, f := range files {
            fmt.Fprint(w, strconv.QuoteToASCII(f.Name()))
            if idx != len(files)-1 {
                fmt.Fprint(w, ",")
            }
        }
        fmt.Fprint(w, "]")
        return
    }

    fp := *mailroot + r.URL.Path

    if _, err := os.Stat(fp); err != nil {
        if ! os.IsNotExist(err) {
            log.Printf("Failed to open user mailbox: %s", err)
        }
        http.Error(w, "{\"error\":\"mailbox not found\"}", 404)
        return
    }

    if mbox, err := mboxparser.ReadFile(fp); err == nil {
        fmt.Fprint(w, "[")
        for idx, mail := range mbox.Messages {
            fmt.Fprint(w, "{")
            for k, vs := range mail.Header {
                for _, v := range vs {
                    fmt.Fprintf(w, "%s: %s,", strconv.QuoteToASCII(k), strconv.QuoteToASCII(v))
                }
            }
            fmt.Fprint(w, "\"Content\": [")
            for idx2, body := range mail.Bodies {
                fmt.Fprint(w, "{")
                for k, vs := range body.Header {
                    for _, v := range vs {
                        fmt.Fprintf(w, "%s: %s,", strconv.QuoteToASCII(k), strconv.QuoteToASCII(v))
                    }
                }
                fmt.Fprint(w, "\"Data\":")
                buf := new(bytes.Buffer)
                buf.ReadFrom(body.Content)
                fmt.Fprintf(w, "%s", strconv.QuoteToASCII(buf.String()))
                fmt.Fprint(w, "}")
                if idx2 != len(mail.Bodies)-1 {
                    fmt.Fprint(w, ",")
                }
            }
            fmt.Fprint(w, "]}")
            if idx != len(mbox.Messages)-1 {
                fmt.Fprint(w, ",")
            }
        }
        fmt.Fprint(w, "]")
    } else {
        http.Error(w, strings.Join([]string{"{\"error\":", strconv.QuoteToASCII(err.Error()), "}"}, ""), 500)
    }
}

func main() {
    http.HandleFunc("/", mboxServer)
    err := http.ListenAndServe(strings.Join([]string{*host, ":", strconv.Itoa(*port)}, ""), nil)

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
