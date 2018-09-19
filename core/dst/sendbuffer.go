// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

import (
	"fmt"
	"runtime/debug"

	"sync"

	"github.com/juju/ratelimit"
)

/*
	                          sendWindow
	                         v
		[S|S|S|S|Q|Q|Q|Q| | | | | | | | | ]
		         ^       ^writeSlot
		          sendSlot
*/
type sendBuffer struct {
	mux       *Mux              // we send packets here
	scheduler *ratelimit.Bucket // sets send rate for packets

	sendWindow int // maximum number of outstanding non-acked packets
	packetRate int // target pps

	send     packetList // buffered packets
	sendSlot int        // buffer slot from which to send next packet

	lost     packetList // list of packets reported lost by timeout
	lostSlot int        // next lost packet to resend

	closed  bool
	closing bool
	mut     sync.Mutex
	cond    *sync.Cond
}

const (
	schedulerRate     = 1e6
	schedulerCapacity = schedulerRate / 40
)

// newSendBuffer creates a new send buffer with a zero window.
// SetRateAndWindow() must be called to set an initial packet rate and send
// window before using.
func newSendBuffer(m *Mux) *sendBuffer {
	b := &sendBuffer{
		mux:       m,
		scheduler: ratelimit.NewBucketWithRate(schedulerRate, schedulerCapacity),
	}
	b.cond = sync.NewCond(&b.mut)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
			}
		}()
		b.writerLoop()
	}()
	return b
}

// Write puts a new packet in send buffer and schedules a send. Blocks when
// the window size is or would be exceeded.
func (b *sendBuffer) Write(pkt packet) error {
	b.mut.Lock()
	defer b.mut.Unlock()

	for b.send.Full() || b.send.Len() >= b.sendWindow {
		if b.closing {
			return ErrClosedConn
		}
		if debugConnection {
			log.Println(b, "Write blocked")
		}
		b.cond.Wait()
	}
	if !b.send.Append(pkt) {
		panic("bug: append failed")
	}
	b.cond.Broadcast()
	return nil
}

// Acknowledge removes packets with lower sequence numbers from the loss list
// or send buffer.
func (b *sendBuffer) Acknowledge(seq sequenceNo) {
	b.mut.Lock()

	if cut := b.lost.CutLessSeq(seq); cut > 0 {
		if debugConnection {
			log.Println(b, "cut", cut, "from loss list")
		}
		// Next resend should always start with the first packet, regardless
		// of what we might already have resent previously.
		b.lostSlot = 0
		b.cond.Broadcast()
	}

	if cut := b.send.CutLessSeq(seq); cut > 0 {
		if debugConnection {
			log.Println(b, "cut", cut, "from send list")
		}
		b.sendSlot -= cut
		b.cond.Broadcast()
	}

	b.mut.Unlock()
}

func (b *sendBuffer) NegativeAck(seq sequenceNo) {
	b.mut.Lock()

	pkts := b.send.PopSequence(seq)
	if cut := len(pkts); cut > 0 {
		b.lost.AppendAll(pkts)
		if debugConnection {
			log.Println(b, "cut", cut, "from send list, adding to loss list")
			log.Println(seq, pkts)
		}
		b.sendSlot -= cut
		b.lostSlot = 0
		b.cond.Broadcast()
	}

	b.mut.Unlock()
}

// ScheduleResend arranges for a resend of all currently unacknowledged
// packets.
func (b *sendBuffer) ScheduleResend() {
	b.mut.Lock()

	if b.sendSlot > 0 {
		// There are packets that have been sent but not acked. Move them from
		// the send buffer to the loss list for retransmission.
		if debugConnection {
			log.Println(b, "scheduled resend from send list", b.sendSlot)
		}

		// Append the packets to the loss list and rewind the send buffer
		b.lost.AppendAll(b.send.All()[:b.sendSlot])
		b.send.Cut(b.sendSlot)
		b.sendSlot = 0
		b.cond.Broadcast()
	}

	if b.lostSlot > 0 {
		// Also resend whatever was already in the loss list
		if debugConnection {
			log.Println(b, "scheduled resend from loss list", b.lostSlot)
		}
		b.lostSlot = 0
		b.cond.Broadcast()
	}

	b.mut.Unlock()
}

// SetWindowAndRate sets the window size (in packets) and packet rate (in
// packets per second) to use when sending.
func (b *sendBuffer) SetWindowAndRate(sendWindow, packetRate int) {
	b.mut.Lock()
	if debugConnection {
		log.Println(b, "new window & rate", sendWindow, packetRate)
	}
	b.packetRate = packetRate
	b.sendWindow = sendWindow
	if b.sendWindow > b.send.Cap() {
		b.send.Resize(b.sendWindow)
		b.cond.Broadcast()
	}
	b.mut.Unlock()
}

// Stop stops the send buffer from any doing further sending, but waits for
// the current buffers to be drained.
func (b *sendBuffer) Stop() {
	b.mut.Lock()

	if b.closed || b.closing {
		return
	}

	b.closing = true
	for b.lost.Len() > 0 || b.send.Len() > 0 {
		b.cond.Wait()
	}

	b.closed = true
	b.cond.Broadcast()
	b.mut.Unlock()
}

// CrashStop stops the send buffer from any doing further sending, without
// waiting for buffers to drain.
func (b *sendBuffer) CrashStop() {
	b.mut.Lock()

	if b.closed || b.closing {
		return
	}

	b.closing = true
	b.closed = true
	b.cond.Broadcast()
	b.mut.Unlock()
}

func (b *sendBuffer) String() string {
	return fmt.Sprintf("sendBuffer@%p", b)
}

func (b *sendBuffer) writerLoop() {
	if debugConnection {
		log.Println(b, "writer() starting")
		defer log.Println(b, "writer() exiting")
	}

	b.scheduler.Take(schedulerCapacity)
	for {
		var pkt packet
		b.mut.Lock()
		for b.lostSlot >= b.sendWindow ||
			(b.sendSlot == b.send.Len() && b.lostSlot == b.lost.Len()) {
			if b.closed {
				b.mut.Unlock()
				return
			}

			if debugConnection {
				log.Println(b, "writer() paused", b.lostSlot, b.sendSlot, b.sendWindow, b.lost.Len())
			}
			b.cond.Wait()
		}

		if b.lostSlot < b.lost.Len() {
			pkt = b.lost.All()[b.lostSlot]
			pkt.hdr.timestamp = timestampMicros()
			b.lostSlot++

			if debugConnection {
				log.Println(b, "resend", b.lostSlot, b.lost.Len(), b.sendWindow, pkt.hdr.connID, pkt.hdr.sequenceNo)
			}
		} else if b.sendSlot < b.send.Len() {
			pkt = b.send.All()[b.sendSlot]
			pkt.hdr.timestamp = timestampMicros()
			b.sendSlot++

			if debugConnection {
				log.Println(b, "send", b.sendSlot, b.send.Len(), b.sendWindow, pkt.hdr.connID, pkt.hdr.sequenceNo)
			}
		}

		b.cond.Broadcast()
		packetRate := b.packetRate
		b.mut.Unlock()

		if pkt.dst != nil {
			b.scheduler.Wait(schedulerRate / int64(packetRate))
			b.mux.write(pkt)
		}
	}
}
