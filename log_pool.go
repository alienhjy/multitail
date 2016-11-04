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
	logPaths map[string]bool
//	namePattern string
	wg       *sync.WaitGroup
}

// NewLogPool create a new pool and return its point.
func NewLogPool() *LogPool {
	pool := new(LogPool)
	pool.logFiles = make(map[string]*TailNode)
	pool.logPaths = make(map[string]bool)
	pool.wg = new(sync.WaitGroup)
	return pool
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
func (slf *LogPool) AddPath(path string, pattern string) {
	path, _ = filepath.Abs(path)
	if _, exist := slf.logPaths[path]; exist {
		log.Println(path + " is existed in path list.")
		return
	}
	slf.logPaths[path] = true

	names := treeDir(path, pattern)
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
	if _, exist := slf.logPaths[path]; !exist {
		log.Println(path + " is not existed in path list.")
		return
	}
	slf.logPaths[path] = false

	// delete all file node under this path.
	// TODO: this implement must be improved.
	for name, _ := range slf.logFiles {
		matched, err := regexp.MatchString(path, name)
		if matched == true || err == nil {
			slf.DeleteOne(name)
		}
	}
	return
}

// StartAll tail process.
func (slf *LogPool) StartAll() {
	for _, node := range slf.logFiles {
		node.Start(slf.wg)
	}
	return
}

// WaitGroup wait all of node exit.
func (slf *LogPool) WaitAll() {
	slf.wg.Wait()
	return
}

