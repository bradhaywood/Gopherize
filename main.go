package main

import(
    "fmt"
    "os"
    "strconv"
    "gopherize"
)

func main() {
    var host string
    var port int
    var err error
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <host> [<port>]\n", os.Args[0])
        os.Exit(1)
    }

    host, port = os.Args[1], 80

    if len(os.Args) > 2 {
        port, err = strconv.Atoi(os.Args[2])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid port\n")
            os.Exit(1)
        }
    }

    g := gopherize.Connect(host, port)
    
    if g.Err != nil {
        fmt.Fprintf(os.Stderr, "Problem connecting: %s\n", g.Err)
        os.Exit(1)
    }
    
    fmt.Println("Connected to " + g.Host)

    if status := g.Get("/asdsd/"); status == false {
        fmt.Fprintf(os.Stderr, "%s\n", g.Err)
        os.Exit(1)
    }

    fmt.Println("Found page " + g.Page)
}
