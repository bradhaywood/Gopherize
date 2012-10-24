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
}

func Connect(host string, port int) (Gopherize, error) {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    var gopherize = Gopherize{conn, "/", host, port }
    return gopherize, err
}

func (g *Gopherize) Get(page string) (bool, error) {
    buf := make([]byte, 1024)
    foundOk := false
    var returnError error
    fmt.Fprintf(g.Conn, "HEAD %s HTTP/1.1\r\n", page)
    fmt.Fprintf(g.Conn, "Host: %s:%d\r\n\r\n", g.Host, g.Port)

    re404, err404 := regexp.Compile(`HTTP/1.1 404`)
    re200, err200 := regexp.Compile(`HTTP/1.1 200`)
    re30x, err30x := regexp.Compile(`HTTP/1.1 30(\d+)`)
    if err404 != nil {
        return false, err404
    }

    if err200 != nil {
        return false, err200
    }

    if err30x != nil {
        return false, err30x
    }

    for {
        readlen, err := g.Conn.Read(buf)

        if err != nil {
            returnError = errors.New(fmt.Sprintf("Problem reading data (%s)", err))
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
            returnError = errors.New("Page Not Found")
            break
        }

        if re30x.Match(buf) {
            number := re30x.FindStringSubmatch(string(buf))
            returnError = errors.New(fmt.Sprintf("30%s Redirect", number[1]))
            break
        }

        //fmt.Printf("%s\n", buf)
    }

    if foundOk == true {
        g.Page = page
        return true, nil
    }

    return false, returnError
}

func (g *Gopherize) Links() ([]string, error) {
    buf := make([]byte, 1024)
    fmt.Fprintf(g.Conn, "GET %s HTTP/1.1\r\n", g.Page)
    fmt.Fprintf(g.Conn, "Host: %s:%d\r\n\r\n", g.Host, g.Port)

    var links []string
    var readlen int
    linkRegex, err := regexp.Compile(`<a\s[^>]*href\s*=\s*\"([^\"]*)\"[^>]*>(.*?)</a>`)
    if err != nil {
        return links, err
    }

    for {
        readlen, err = g.Conn.Read(buf)
        if readlen == 0 {
            break
        }

        if err != nil {
            return links, err
        }

        //fmt.Printf("%s\n", buf)
        if linkRegex.Match(buf) {
            server := linkRegex.FindStringSubmatch(string(buf))
            links = append(links, server[1])
        }
    }

    return links, nil
}

