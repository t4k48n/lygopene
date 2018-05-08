package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	maxMin = 6 * 60
)

func main() {
	if len(os.Args) != 3 {
		printUsage()
		os.Exit(1)
	}
	wm, werr := strconv.Atoi(os.Args[1])
	rm, rerr := strconv.Atoi(os.Args[2])
	if werr != nil || wm < 1 || wm > maxMin || rerr != nil || rm < 1 || rm > maxMin {
		printUsage()
		os.Exit(1)
	}
	workMin, restMin := time.Duration(wm)*time.Minute, time.Duration(rm)*time.Minute

	if termbox.Init() != nil {
		printInitializationError()
		os.Exit(1)
	}
	defer termbox.Close()

	c_c_ch := make(chan bool)
	res_ch := make(chan bool)
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
				c_c_ch <- true
			}
			if ev.Type == termbox.EventResize {
				res_ch <- true
			}
		}
	}()

	t1, t2 := time.NewTicker(workMin), &time.Ticker{}
	ccolor := termbox.ColorRed
	drawAll(ccolor)
	termbox.Flush()
	for {
		select {
		case <-t1.C:
			t1.Stop()
			t2 = time.NewTicker(restMin)
			ccolor = termbox.ColorGreen
			drawAll(ccolor)
			termbox.Flush()
		case <-t2.C:
			t2.Stop()
			t1 = time.NewTicker(workMin)
			ccolor = termbox.ColorRed
			drawAll(ccolor)
			termbox.Flush()
		case <-res_ch:
			drawAll(ccolor)
			termbox.Flush()
		case <-c_c_ch:
			return
		}
	}
}

func printUsage() {
	fmt.Printf("Simplest timer for Pomodoro Technique.\n\nusage: %v <work minute> <rest minute>\n\n", os.Args[0])
}

func printInitializationError() {
	fmt.Errorf("initialization failed\n\n")
}

func drawAll(c termbox.Attribute) {
	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			termbox.SetCell(x, y, ' ', c, c)
		}
	}
}
