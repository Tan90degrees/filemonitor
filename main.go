package main

import (
	"fileMon/factory"
	"fileMon/work"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var err error
	var monTime int
	lenArg := len(os.Args)
	if lenArg == 4 {
		monTime, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		monTime = 120
	}
	fmt.Printf("Time is setted at %d mins.\n", monTime)
	for {
		timeStart := time.Now()
		if lenArg == 3 || lenArg == 4 {
			src := path.Clean(os.Args[1])
			dst := path.Clean(os.Args[2])
			wchan := make(chan int, factory.Workers)
			factory.InitWorkers(&wchan)
			dp, err := os.ReadDir(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range dp {
				iName := i.Name()
				if i.IsDir() {
					dName := path.Join(dst, iName)
					os.MkdirAll(dName, 0777)
					if factory.SeekWorker(&wchan) {
						wg.Add(1)
						go factory.Work(path.Join(src, iName), dName, &wchan, &wg)
					} else {
						factory.DoWork(path.Join(src, iName), dName, &wchan, &wg)
					}
				} else {
					work.Move(path.Join(src, iName), path.Join(dst, iName))
				}
			}
			wg.Wait()
			close(wchan)
			for _, i := range dp {
				os.RemoveAll(path.Join(src, i.Name()))
			}
			fmt.Println(time.Since(timeStart).Seconds(), "s")
		} else {
			fmt.Println("Wrong parameters.")
			os.Exit(1)
		}
		work.Mon(monTime)
	}
}
