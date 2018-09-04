// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

// Error represents the various dst-internal error conditions.
type Error struct {
	Err string
}

// Error returns a string representation of the error.
func (e Error) Error() string {
	return e.Err
}

var (
	ErrClosedConn       = &Error{"operation on closed connection"}
	ErrClosedMux        = &Error{"operation on closed mux"}
	ErrHandshakeTimeout = &Error{"handshake timeout"}
	ErrNotDST           = &Error{"network is not dst"}
	ErrNotImplemented   = &Error{"not implemented"}
)
