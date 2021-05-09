// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

/*
import (
	"io"
)

// MsgCmpctBlock implements the Message interface and represents a bitcoin
// CmpctBlock message.  
//
type MsgCmpctBlock struct {
	Header BlockHeader
	Nonce int64
	ShortIdCount int64
	ShortIds[] int64
	PrefilledTxnCount int64
	PrefilledTxn[] Transaction
}
// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCmpctBlock) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return readElements(r, &msg.Announce, &msg.Version)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgCmpctBlock) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	return writeElements(w, msg.Announce, msg.Version)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgCmpctBlock) Command() string {
	return CmdCmpctBlock
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgCmpctBlock) MaxPayloadLength(pver uint32) uint32 {
	return 999 //?
}

// NewMsgCmpctBlock returns a new bitcoin CmpctBlock message that conforms to
// the Message interface.  See MsgCmpctBlock for details.
func NewMsgCmpctBlock(announce int8, version int64) *MsgCmpctBlock {
	return &MsgCmpctBlock{
		Announce: announce,
		Version: version,
	}
}
*/