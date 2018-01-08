package main

import (
	"fmt"
	"io"
	"os"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func playMP3(file string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return
	}
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return
	}
}
