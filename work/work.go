package work

import (
	"io"
	"log"
	"os"
	"time"
)

func Move(src string, dst string) {
	sp, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	dp, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(dp, sp)
	if err != nil {
		log.Fatal(err)
	}
	sp.Close()
	dp.Close()
	err = os.RemoveAll(src)
	if err != nil {
		log.Fatal(err)
	}
}

func Mon(min int) {
	time.Sleep(time.Minute * time.Duration(min))
}
