package kcpcfg

import kcp "github.com/xtaci/kcp-go"

type KCPConfigArgs struct {
	Key          *string
	Crypt        *string
	Mode         *string
	MTU          *int
	SndWnd       *int
	RcvWnd       *int
	DataShard    *int
	ParityShard  *int
	DSCP         *int
	NoComp       *bool
	AckNodelay   *bool
	NoDelay      *int
	Interval     *int
	Resend       *int
	NoCongestion *int
	SockBuf      *int
	KeepAlive    *int
	Block        kcp.BlockCrypt
}
