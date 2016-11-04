package main

import (
//	"bufio"
	"fmt"
	"log"
//	"sync"

//	"github.com/droyo/tailpipe"
)

func main() {
	defer log.Println("DEBUG: exit...")
	fileList := treeDir(".", ".log\\d")
	for _, name := range fileList {
		fmt.Println(name)
	}
	pool := NewLogPool()
	pool.AddPath(".", ".log\\d")
	pool.StartAll()
	pool.WaitAll()
	return
}
