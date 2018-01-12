package main

import (
	"context"
	"fmt"
	"io"
	"os"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func playMP3(ctx context.Context, file string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return
	}

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return
	}

	fmt.Printf("Length: %d[bytes]\n", d.Length())
	doneCopy := make(chan bool)
	go func() {
		io.Copy(p, d)
		doneCopy <- true
	}()
L:
	for {
		select {
		case <-doneCopy:
			p.Close()
			d.Close()
			break L
		case <-ctx.Done():
			p.Close()
			d.Close()
			break L
		}
	}
}
