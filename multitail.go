package main

import (
//	"bufio"
	"fmt"
	"log"
//	"sync"

//	"github.com/droyo/tailpipe"
)

func main() {
	/*
	var (
		logFiles = make(map[string]TailNode)
		wg       sync.WaitGroup
		err      error
	)
	logFiles["pacman.log1"] = TailNode{enabled: true}
	logFiles["pacman.log2"] = TailNode{enabled: true}

	for name, logFile := range logFiles {
		if logFile.enabled == false {
			continue
		}
		logFile.tf, err = tailpipe.Open(name)
		if err != nil {
			log.Println(err)
		}
		wg.Add(1)
		tf := logFile.tf
		go func() {
			defer func() {
				tf.Close()
				wg.Done()
			}()
			scanner := bufio.NewScanner(tf)
			for scanner.Scan() {
				fmt.Println(string(scanner.Bytes()))
				//os.Stdout.Write(scanner.Bytes())
			}
		}()
	}

	wg.Wait()
	*/
	fileList := treeDir(".", ".go$")
	for _, name := range fileList {
		fmt.Println(name)
	}
	log.Println("DEBUG: exit...")
	return
}
