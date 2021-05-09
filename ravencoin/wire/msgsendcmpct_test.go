// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

/*
import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)
// TestSendCmpctLatest tests the MsgSendCmpct API against the latest protocol version.
func TestSendCmpctLatest(t *testing.T) {
	pver := ProtocolVersion

	t.Errorf("Not converted yet...");

	minfee := rand.Int63()
	msg := NewMsgSendCmpct(minfee)
	if msg.MinFee != minfee {
		t.Errorf("NewMsgSendCmpct: wrong minfee - got %v, want %v",
			msg.MinFee, minfee)
	}

	// Ensure the command is expected value.
	wantCmd := "SendCmpct"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgSendCmpct: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(8)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Test encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.BtcEncode(&buf, pver, BaseEncoding)
	if err != nil {
		t.Errorf("encode of MsgSendCmpct failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := NewMsgSendCmpct(0)
	err = readmsg.BtcDecode(&buf, pver, BaseEncoding)
	if err != nil {
		t.Errorf("decode of MsgSendCmpct failed [%v] err <%v>", buf, err)
	}

	// Ensure minfee is the same.
	if msg.MinFee != readmsg.MinFee {
		t.Errorf("Should get same minfee for protocol version %d", pver)
	}
}

// TestSendCmpctWire tests the MsgSendCmpct wire encode and decode for various protocol
// versions.
func TestSendCmpctWire(t *testing.T) {
	tests := []struct {
		in   MsgSendCmpct // Message to encode
		out  MsgSendCmpct // Expected decoded message
		buf  []byte       // Wire encoding
		pver uint32       // Protocol version for wire encoding
	}{
		// Latest protocol version.
		{
			MsgSendCmpct{MinFee: 123123}, // 0x1e0f3
			MsgSendCmpct{MinFee: 123123}, // 0x1e0f3
			[]byte{0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
			ProtocolVersion,
		},

		// Protocol version SendCmpctVersion
		{
			MsgSendCmpct{MinFee: 456456}, // 0x6f708
			MsgSendCmpct{MinFee: 456456}, // 0x6f708
			[]byte{0x08, 0xf7, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00},
			SendCmpctVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		err := test.in.BtcEncode(&buf, test.pver, BaseEncoding)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from wire format.
		var msg MsgSendCmpct
		rbuf := bytes.NewReader(test.buf)
		err = msg.BtcDecode(rbuf, test.pver, BaseEncoding)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(msg, test.out) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(msg), spew.Sdump(test.out))
			continue
		}
	}
}

// TestSendCmpctWireErrors performs negative tests against wire encode and decode
// of MsgSendCmpct to confirm error paths work correctly.
func TestSendCmpctWireErrors(t *testing.T) {
	pver := ProtocolVersion
	pverNoSendCmpct := SendCmpctVersion - 1
	wireErr := &MessageError{}

	baseSendCmpct := NewMsgSendCmpct(123123) // 0x1e0f3
	baseSendCmpctEncoded := []byte{
		0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	tests := []struct {
		in       *MsgSendCmpct // Value to encode
		buf      []byte        // Wire encoding
		pver     uint32        // Protocol version for wire encoding
		max      int           // Max size of fixed buffer to induce errors
		writeErr error         // Expected write error
		readErr  error         // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force error in minfee.
		{baseSendCmpct, baseSendCmpctEncoded, pver, 0, io.ErrShortWrite, io.EOF},
		// Force error due to unsupported protocol version.
		{baseSendCmpct, baseSendCmpctEncoded, pverNoSendCmpct, 4, wireErr, wireErr},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		w := newFixedWriter(test.max)
		err := test.in.BtcEncode(w, test.pver, BaseEncoding)
		if reflect.TypeOf(err) != reflect.TypeOf(test.writeErr) {
			t.Errorf("BtcEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.writeErr {
				t.Errorf("BtcEncode #%d wrong error got: %v, "+
					"want: %v", i, err, test.writeErr)
				continue
			}
		}

		// Decode from wire format.
		var msg MsgSendCmpct
		r := newFixedReader(test.max, test.buf)
		err = msg.BtcDecode(r, test.pver, BaseEncoding)
		if reflect.TypeOf(err) != reflect.TypeOf(test.readErr) {
			t.Errorf("BtcDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}

		// For errors which are not of type MessageError, check them for
		// equality.
		if _, ok := err.(*MessageError); !ok {
			if err != test.readErr {
				t.Errorf("BtcDecode #%d wrong error got: %v, "+
					"want: %v", i, err, test.readErr)
				continue
			}
		}

	}
}
*/