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

    g, err := gopherize.Connect(host, port)
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Problem connecting: %s\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Connected to " + g.Host)

    if status, err := g.Get("/"); status == false {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }

    fmt.Println("Found page " + g.Page)
    
    links, err := g.Links()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting links: %s\n", err)
        os.Exit(1)
    }

    for _, link := range links {
        fmt.Println(link)
    }
}
