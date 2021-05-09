// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"strconv"
	"strings"
)

// CHECK Ravencoin

const (
	// ProtocolVersion is the latest protocol version this package supports.
	ProtocolVersion uint32 = 70028

	//! initial proto version, to be increased after version/verack negotiation
	INIT_PROTO_VERSION uint32 = 209;
	
	//! In this version, 'getheaders' was introduced.
	GETHEADERS_VERSION uint32 = 31800;
	
	//! assetdata network request is allowed for this version
	ASSETDATA_VERSION uint32 = 70017;
	
	//! getassetdata reutrn asstnotfound, and assetdata doesn't have blockhash in the data
	X16RV2_VERSION uint32 = 70025;
	
	//! getassetdata reutrn asstnotfound, and assetdata doesn't have blockhash in the data
	KAWPOW_VERSION uint32 = 70027;
	
	//! disconnect from peers older than this proto version
	//!!! Anytime this value is changed please also update the "MY_VERSION" value to match in the
	//!!! ./test/functional/test_framework/mininode.py file. Not doing so will cause verack to fail!
	MIN_PEER_PROTO_VERSION uint32 = X16RV2_VERSION;
	
	//! nTime field added to CAddress, starting with this version;
	//! if possible, avoid requesting addresses nodes older than this
	CADDR_TIME_VERSION uint32 = 31402;
	
	//! BIP 0031, pong message, is enabled for all versions AFTER this one
	BIP0031_VERSION uint32 = 60000;
	
	//! "filter*" commands are disabled without NODE_BLOOM after and including this version
	NO_BLOOM_VERSION uint32 = 70011;
	
	//! "sendheaders" command and announcing blocks with headers starts with this version
	SENDHEADERS_VERSION uint32 = 70012;
	
	//! "feefilter" tells peers to filter invs to you by fee starts with this version
	FEEFILTER_VERSION uint32 = 70013;
	
	//! short-id-based block download starts with this version
	SHORT_IDS_BLOCKS_VERSION uint32 = 70014;
	
	//! not banning for invalid compact blocks starts with this version
	INVALID_CB_NO_BAN_VERSION uint32 = 70015;
	
	//! getassetdata reutrn asstnotfound, and assetdata doesn't have blockhash in the data
	ASSETDATA_VERSION_UPDATED uint32 = 70020;
	
	//! In this version, 'rip5 (messaging and restricted assets)' was introduced
	MESSAGING_RESTRICTED_ASSETS_VERSION uint32 = 70026;
)

// ServiceFlag identifies services supported by a bitcoin peer.
type ServiceFlag uint64

const (
	// SFNodeNetwork is a flag used to indicate a peer is a full node.
	SFNodeNetwork ServiceFlag = 1 << iota

	// SFNodeGetUTXO is a flag used to indicate a peer supports the
	// getutxos and utxos commands (BIP0064).
	SFNodeGetUTXO

	// SFNodeBloom is a flag used to indicate a peer supports bloom
	// filtering.
	SFNodeBloom

	// SFNodeWitness is a flag used to indicate a peer supports blocks
	// and transactions including witness data (BIP0144).
	SFNodeWitness

	// SFNodeXthin is a flag used to indicate a peer supports xthin blocks.
	SFNodeXthin

	// SFNodeBit5 is a flag used to indicate a peer supports a service
	// defined by bit 5.
	SFNodeBit5

	// SFNodeCF is a flag used to indicate a peer supports committed
	// filters (CFs).
	SFNodeCF

	// SFNode2X is a flag used to indicate a peer is running the Segwit2X
	// software.
	SFNode2X
)

// Map of service flags back to their constant names for pretty printing.
var sfStrings = map[ServiceFlag]string{
	SFNodeNetwork: "SFNodeNetwork",
	SFNodeGetUTXO: "SFNodeGetUTXO",
	SFNodeBloom:   "SFNodeBloom",
	SFNodeWitness: "SFNodeWitness",
	SFNodeXthin:   "SFNodeXthin",
	SFNodeBit5:    "SFNodeBit5",
	SFNodeCF:      "SFNodeCF",
	SFNode2X:      "SFNode2X",
}

// orderedSFStrings is an ordered list of service flags from highest to
// lowest.
var orderedSFStrings = []ServiceFlag{
	SFNodeNetwork,
	SFNodeGetUTXO,
	SFNodeBloom,
	SFNodeWitness,
	SFNodeXthin,
	SFNodeBit5,
	SFNodeCF,
	SFNode2X,
}

// String returns the ServiceFlag in human-readable form.
func (f ServiceFlag) String() string {
	// No flags are set.
	if f == 0 {
		return "0x0"
	}

	// Add individual bit flags.
	s := ""
	for _, flag := range orderedSFStrings {
		if f&flag == flag {
			s += sfStrings[flag] + "|"
			f -= flag
		}
	}

	// Add any remaining flags which aren't accounted for as hex.
	s = strings.TrimRight(s, "|")
	if f != 0 {
		s += "|0x" + strconv.FormatUint(uint64(f), 16)
	}
	s = strings.TrimLeft(s, "|")
	return s
}

// BitcoinNet represents which bitcoin network a message belongs to.
type BitcoinNet uint32

// Constants used to indicate the message Raven network.  They can also be
// used to seek to the next message when a stream's state is unknown, but
// this package does not provide that functionality since it's generally a
// better idea to simply disconnect clients that are misbehaving over TCP.
const (
	// MainNet represents the main Raven network.
	MainNet BitcoinNet = 0x63617368 // CHECK Ravencoin

	// TestNet represents the testnet network.
	TestNet BitcoinNet = 0xbff2cde6 // CHECK Ravencoin

	// Regtest represents the regtest network.
	Regtest BitcoinNet = 0x2f54cc9d // CHECK Ravencoin
)

// bnStrings is a map of bitcoin networks back to their constant names for
// pretty printing.
var bnStrings = map[BitcoinNet]string{
	MainNet: "MainNet",
	TestNet: "TestNet",
	Regtest: "Regtest",
}

// String returns the BitcoinNet in human-readable form.
func (n BitcoinNet) String() string {
	if s, ok := bnStrings[n]; ok {
		return s
	}

	return fmt.Sprintf("Unknown BitcoinNet (%d)", uint32(n))
}
