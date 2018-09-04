// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

import (
	"fmt"
	"io"
	"os"
	"time"
)

type windowCC struct {
	minWindow     int
	maxWindow     int
	currentWindow int
	minRate       int
	maxRate       int
	currentRate   int
	targetRate    int

	curRTT time.Duration
	minRTT time.Duration

	statsFile io.WriteCloser
	start     time.Time
}

func newWindowCC() *windowCC {
	var statsFile io.WriteCloser

	if debugCC {
		statsFile, _ = os.Create(fmt.Sprintf("cc-log-%d.csv", time.Now().Unix()))
		fmt.Fprintf(statsFile, "ms,minWin,maxWin,curWin,minRate,maxRate,curRate,minRTT,curRTT\n")
	}

	return &windowCC{
		minWindow:     1, // Packets
		maxWindow:     16 << 10,
		currentWindow: 1,

		minRate:     100,  // PPS
		maxRate:     80e3, // Roughly 1 Gbps at 1500 bytes per packet
		currentRate: 100,
		targetRate:  1000,

		minRTT:    10 * time.Second,
		statsFile: statsFile,
		start:     time.Now(),
	}
}

func (w *windowCC) Ack() {
	if w.curRTT > w.minRTT+100*time.Millisecond {
		return
	}

	changed := false

	if w.currentWindow < w.maxWindow {
		w.currentWindow++
		changed = true
	}

	if w.currentRate != w.targetRate {
		w.currentRate = (w.currentRate*7 + w.targetRate) / 8
		changed = true
	}

	if changed && debugCC {
		w.log()
		log.Println("Ack", w.currentWindow, w.currentRate)
	}
}

func (w *windowCC) NegAck() {
	if w.currentWindow > w.minWindow {
		w.currentWindow /= 2
	}
	if w.currentRate > w.minRate {
		w.currentRate /= 2
	}
	if debugCC {
		w.log()
		log.Println("NegAck", w.currentWindow, w.currentRate)
	}
}

func (w *windowCC) Exp() {
	w.currentWindow = w.minWindow
	if debugCC {
		w.log()
		log.Println("Exp", w.currentWindow, w.currentRate)
	}
}

func (w *windowCC) SendWindow() int {
	if w.currentWindow < w.minWindow {
		return w.minWindow
	}
	if w.currentWindow > w.maxWindow {
		return w.maxWindow
	}
	return w.currentWindow
}

func (w *windowCC) PacketRate() int {
	if w.currentRate < w.minRate {
		return w.minRate
	}
	if w.currentRate > w.maxRate {
		return w.maxRate
	}
	return w.currentRate
}

func (w *windowCC) UpdateRTT(rtt time.Duration) {
	w.curRTT = rtt
	if w.curRTT < w.minRTT {
		w.minRTT = w.curRTT
		if debugCC {
			log.Println("Min RTT", w.minRTT)
		}
	}

	if w.curRTT > w.minRTT+200*time.Millisecond && w.targetRate > 2*w.minRate {
		w.targetRate -= w.minRate
	} else if w.curRTT < w.minRTT+20*time.Millisecond && w.targetRate < w.maxRate {
		w.targetRate += w.minRate
	}

	if debugCC {
		w.log()
		log.Println("RTT", w.curRTT, "target rate", w.targetRate, "current rate", w.currentRate, "current window", w.currentWindow)
	}
}

func (w *windowCC) log() {
	if w.statsFile == nil {
		return
	}
	fmt.Fprintf(w.statsFile, "%.02f,%d,%d,%d,%d,%d,%d,%.02f,%.02f\n", time.Since(w.start).Seconds()*1000, w.minWindow, w.maxWindow, w.currentWindow, w.minRate, w.maxRate, w.currentRate, w.minRTT.Seconds()*1000, w.curRTT.Seconds()*1000)
}
