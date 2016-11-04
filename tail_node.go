package main

import (
	"bufio"
	"errors"
	"log"
	"fmt"
	"sync"

	"github.com/droyo/tailpipe"
)

// TailNode : Each TailNode will handle one of the log files.
type TailNode struct {
	path    string
	tf      *tailpipe.File
	enabled int
	running bool
	quit    chan bool
	done    chan bool
}

// Start the tail process of this node.
func (slf *TailNode) Start(wg *sync.WaitGroup) error {
	var (
		err error
	)
	if slf.enabled == 0 {
		err = errors.New("This node is not been enabled.")
		return err
	}
	slf.tf, err = tailpipe.Open(slf.path)
	if err != nil {
		return err
	}
	wg.Add(1)
	slf.running = true
	go func() {
		defer func() {
			// DEBUG:
			log.Println(slf.path + " exit...")
			slf.tf.Close()
			slf.running = false
			slf.done <- true
			wg.Done()
		}()
		scanner := bufio.NewScanner(slf.tf)
		for scanner.Scan() {
			select {
			case <-slf.quit:
				// TODO: Do something before quit.
				return
			default:
				fmt.Println(string(scanner.Bytes()))
			}
		}
	}()

	return nil
}

// Stop the tail process of this node.
func (slf *TailNode) Stop() {
	slf.quit <- true
	<-slf.done
	return
}

// TODO: delete thos code.
// NewTailNode create new instance of TailNode.
/*
func NewTailNode(name string) *TailNode {
	var (
		tailNode = new(TailNode)
	)
	tailNode.path = name
	tailNode.enabled = 1
	tailNode.running = false
	tailNode.quit = make(chan bool)
	return tailNode
}
*/

