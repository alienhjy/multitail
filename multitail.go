package main

import (
	"fmt"
	"log"
	"sync"
	"bufio"
	"github.com/droyo/tailpipe"
)

type TailNode struct {
	r *tailpipe.File
	enabled bool
}

func main() {
	var (
		logFiles map[string]TailNode = make(map[string]TailNode)
		wg sync.WaitGroup
		err error
	)
	logFiles["pacman.log1"] = TailNode{enabled: true}
	logFiles["pacman.log2"] = TailNode{enabled: true}

	for name, logFile := range logFiles {
		if logFile.enabled == false {
			continue
		}
		logFile.r, err = tailpipe.Open(name)
		if err != nil {
			log.Println(err)
		}
		wg.Add(1)
		go func() {
			defer func() {
				logFile.r.Close()
				wg.Done()
			}()
			scanner := bufio.NewScanner(logFile.r)
			for scanner.Scan() {
				fmt.Println(string(scanner.Bytes()))
				//os.Stdout.Write(scanner.Bytes())
			}
		}()
	}

	wg.Wait()
	log.Println("DEBUG: exit...")
	return
}
