// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

/*
import (
	"io"
)
// MsgSendCmpct implements the Message interface and represents a bitcoin
// sendcmpct message.  
//
type MsgSendCmpct struct {
	Announce int8
	Version int64
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgSendCmpct) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return readElement(r, &msg.Announce, &msg.Version)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgSendCmpct) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	return writeElement(w, msg.Announce, msg.Version)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgSendCmpct) Command() string {
	return CmdSendCmpct
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgSendCmpct) MaxPayloadLength(pver uint32) uint32 {
	return 9
}

// NewMsgSendCmpct returns a new bitcoin SendCmpct message that conforms to
// the Message interface.  See MsgSendCmpct for details.
func NewMsgSendCmpct(announce int8, version int64) *MsgSendCmpct {
	return &MsgSendCmpct{
		Announce: announce,
		Version: version,
	}
}
*/