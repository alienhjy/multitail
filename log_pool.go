package main

import (
	"os"
	"log"
	"errors"
	"sync"
	"regexp"
	"path/filepath"

//	"github.com/droyo/tailpipe"
)

func treeDir(path string, reg string) []string {
	var fileList = make([]string, 0, 16)
	err := filepath.Walk(path,
		func(name string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			var matched bool
			matched, err = regexp.MatchString(reg, name)
			if matched == false || err != nil {
				return err
			}
			name, _ = filepath.Abs(name)
			fileList = append(fileList, name)
			return nil
		})
	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
	}
	return fileList
}

// LogPool : Manage all log "tail" node.
type LogPool struct {
	logFiles map[string]*TailNode
	logPaths []string
	namePattern string
	wg       sync.WaitGroup
}

// Init : Init log pool.
func (slf *LogPool) Init(namePattern string) {
	slf.logFiles = make(map[string]*TailNode)
	slf.logPaths = make([]string, 0, 16)
	slf.namePattern = namePattern
}

// AddOne add a tail_node to pool.
func (slf *LogPool) AddOne(name string) error {
	var (
		err error
	)
	if slf.logFiles[name] != nil {
		slf.logFiles[name].enabled++
		err = errors.New("Node " + name + " is exist.")
		return err
	}
	tailNode := new(TailNode)
	tailNode.path = name
	tailNode.enabled = 1
	tailNode.running = false
	tailNode.quit = make(chan bool)
	slf.logFiles[name] = tailNode
	return nil
}

// AddPath add a path to pool, and listen change of this path.
func (slf *LogPool) AddPath(path string) {
	path, _ = filepath.Abs(path)
	for _, log_path := range slf.logPaths {
		if path == log_path {
			return
		}
	}
	slf.logPaths = append(slf.logPaths, path)
	names := treeDir(name, slf.namePattern)
	for _, name := range names {
		slf.AddOne(name)
	}
}

// DeleteOne a tail_node.
func (slf *LogPool) DeleteOne(name string) {
	slf.logFiles[name].enabled--
	if slf.logFiles[name].enabled > 0 {
		return
	}
	if slf.logFiles[name].running {
		slf.logFiles[name].Stop()
	}
	delete(slf.logFiles, name)
	return
}

// DeletePath delete a path from LogPool,
// and try remove sub entry of it.
func (slf *LogPool) DeletePath(path string) {
	return
}

// StartAll tail process.
func (slf *LogPool) StartAll() {
	return
}

// WaitGroup wait all of node exit.
func (slf *LogPool) WaitGroup() {
	slf.wg.Wait()
	return
}

