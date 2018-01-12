package main

import (
	"context"
	"io"
	"log"
	"os"
	"sync"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var cgoMutex = sync.Mutex{} // attempt to fix cgo breaking...

func playMP3(ctx context.Context, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Println(err)
		return
	}

	cgoMutex.Lock()
	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	cgoMutex.Unlock()
	if err != nil {
		log.Println(err)
		return
	}

	doneCopy := make(chan bool)
	go func() {
		io.Copy(p, d)
		doneCopy <- true
	}()
L:
	for {
		select {
		case <-doneCopy:
			cgoMutex.Lock()
			f.Close()
			p.Close()
			d.Close()
			cgoMutex.Unlock()
			break L
		case <-ctx.Done():
			cgoMutex.Lock()
			f.Close()
			p.Close()
			d.Close()
			cgoMutex.Unlock()
			break L
		}
	}
}
