package main

import (
        "bytes"
        "fmt"
        "log"
        "os/exec"
        "strings"
)

func node_exporter() {
        cmd := exec.Command("/go/bin/node_exporter")
        cmd.Stdin = strings.NewReader("")
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("node_exporter: %q\n", out.String())
}

func auth() {
        cmd := exec.Command("/go/bin/auth")
        cmd.Stdin = strings.NewReader("")
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("auth: %q\n", out.String())
}

func main() {
	go node_exporter()
	auth()
}
