package ss

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

func SetReadTimeout(c net.Conn) {
	c.SetReadDeadline(time.Now().Add(time.Second * 5))
}

// PipeThenClose copies data from src to dst, closes dst when done.
func PipeThenClose(src, dst net.Conn, addFlow func(int)) {
	defer dst.Close()
	buf := leakyBuf.Get()
	defer leakyBuf.Put(buf)
	for {
		SetReadTimeout(src)
		n, err := src.Read(buf)
		if addFlow != nil {
			addFlow(n)
		}
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			// Note: avoid overwrite err returned by Read.
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Println("write:", err)
				break
			}
		}
		if err != nil {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			/*
				if bool(log.) && err != io.EOF {
					log.Println("read:", err)
				}
			*/
			break
		}
	}
	return
}

// PipeThenClose copies data from src to dst, closes dst when done, with ota verification.
func PipeThenCloseOta(src *Conn, dst net.Conn, addFlow func(int)) {
	const (
		dataLenLen  = 2
		hmacSha1Len = 10
		idxData0    = dataLenLen + hmacSha1Len
	)

	defer func() {
		dst.Close()
	}()
	// sometimes it have to fill large block
	buf := leakyBuf.Get()
	defer leakyBuf.Put(buf)
	for i := 1; ; i += 1 {
		SetReadTimeout(src)
		if n, err := io.ReadFull(src, buf[:dataLenLen+hmacSha1Len]); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("conn=%p #%v read header error n=%v: %v", src, i, n, err)
			break
		}
		dataLen := binary.BigEndian.Uint16(buf[:dataLenLen])
		expectedHmacSha1 := buf[dataLenLen:idxData0]

		var dataBuf []byte
		if len(buf) < int(idxData0+dataLen) {
			dataBuf = make([]byte, dataLen)
		} else {
			dataBuf = buf[idxData0 : idxData0+dataLen]
		}
		if n, err := io.ReadFull(src, dataBuf); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("conn=%p #%v read data error n=%v: %v", src, i, n, err)
			break
		}
		addFlow(int(dataLen))
		chunkIdBytes := make([]byte, 4)
		chunkId := src.GetAndIncrChunkId()
		binary.BigEndian.PutUint32(chunkIdBytes, chunkId)
		actualHmacSha1 := HmacSha1(append(src.GetIv(), chunkIdBytes...), dataBuf)
		if !bytes.Equal(expectedHmacSha1, actualHmacSha1) {
			log.Printf("conn=%p #%v read data hmac-sha1 mismatch, iv=%v chunkId=%v src=%v dst=%v len=%v expeced=%v actual=%v", src, i, src.GetIv(), chunkId, src.RemoteAddr(), dst.RemoteAddr(), dataLen, expectedHmacSha1, actualHmacSha1)
			break
		}
		if n, err := dst.Write(dataBuf); err != nil {
			log.Printf("conn=%p #%v write data error n=%v: %v", dst, i, n, err)
			break
		}
	}
	return
}
