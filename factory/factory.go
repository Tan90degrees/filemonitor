package factory

import (
	"fileMon/work"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
)

var Workers = runtime.NumCPU()

func InitWorkers(wchan *chan int) {
	for i := 0; i < Workers; i++ {
		*wchan <- 1
	}
}

func SeekWorker(wchan *chan int) bool {
	select {
	case <-*wchan:
		return true
	default:
		return false
	}
}

func Work(src string, dst string, wchan *chan int, wg *sync.WaitGroup) {
	dp, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range dp {
		iName := i.Name()
		if i.IsDir() {
			dName := path.Join(dst, iName)
			os.MkdirAll(dName, 0777)
			if SeekWorker(wchan) {
				wg.Add(1)
				go Work(path.Join(src, iName), dName, wchan, wg)
			} else {
				DoWork(path.Join(src, iName), dName, wchan, wg)
			}
		} else {
			work.Move(path.Join(src, iName), path.Join(dst, iName))
		}
	}
	*wchan <- 1
	wg.Done()
}

func DoWork(src string, dst string, wchan *chan int, wg *sync.WaitGroup) {
	dp, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range dp {
		iName := i.Name()
		if i.IsDir() {
			dName := path.Join(dst, iName)
			os.MkdirAll(dName, 0777)
			if SeekWorker(wchan) {
				wg.Add(1)
				go Work(path.Join(src, iName), dName, wchan, wg)
			} else {
				DoWork(path.Join(src, iName), dName, wchan, wg)
			}
		} else {
			work.Move(path.Join(src, iName), path.Join(dst, iName))
		}
	}
}
