package gopherize

import (
    "fmt"
    "net"
    "errors"
    "regexp"
)

type Gopherize struct {
    Conn net.Conn
    Page string
    Host string
    Port int
    Err  error
}

func Connect(host string, port int) Gopherize {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    var gopherize = Gopherize{conn, "/", host, port, err }
    return gopherize
}

func Get(g Gopherize, page string) Gopherize {
    buf := make([]byte, 1024)
    foundOk := false
    fmt.Fprintf(g.Conn, "HEAD %s HTTP/1.1\r\n", page)
    fmt.Fprintf(g.Conn, "Host: %s:%d\r\n\r\n", g.Host, g.Port)

    re404, err404 := regexp.Compile(`HTTP/1.1 404`)
    re200, err200 := regexp.Compile(`HTTP/1.1 200`)

    if err404 != nil {
        g.Err = err404 
        return g
    }

    if err200 != nil {
        g.Err = err200 
        return g
    }

    for {
        readlen, err := g.Conn.Read(buf)

        if err != nil {
            g.Err = err
            break
        }

        if readlen == 0 {
            break
        }

        if re200.Match(buf) {
            foundOk = true
            break
        }

        if re404.Match(buf) {
            g.Err = errors.New("Page Not Found")
            break
        }

        fmt.Printf("%s\n", buf)
    }

    if foundOk == true {
        g.Page = page
        g.Err  = nil
    }

    return g
}
