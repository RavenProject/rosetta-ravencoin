// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/RavenProject/rosetta-ravencoin/ravencoin"
	"github.com/RavenProject/rosetta-ravencoin/configuration"
	mocks "github.com/RavenProject/rosetta-ravencoin/mocks/services"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func forceHexDecode(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("could not decode hex %s", s)
	}

	return b
}

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := types.MarshalMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", types.PrintStruct(i))
	}

	return m
}

func TestConstructionService(t *testing.T) {
	networkIdentifier = &types.NetworkIdentifier{
		Network:    ravencoin.TestnetNetwork,
		Blockchain: ravencoin.Blockchain,
	}

	cfg := &configuration.Configuration{
		Mode:     configuration.Online,
		Network:  networkIdentifier,
		Params:   ravencoin.TestnetParams,
		Currency: ravencoin.TestnetCurrency,
	}

	mockIndexer := &mocks.Indexer{}
	mockClient := &mocks.Client{}
	servicer := NewConstructionAPIService(cfg, mockClient, mockIndexer)
	ctx := context.Background()

	// Test Derive
	publicKey := &types.PublicKey{
		Bytes: forceHexDecode(
			t,
			"03b0da749730dc9b4b1f4a14d6902877a92541f5368778853d9c4a0cb7802dcfb2",
		),
		CurveType: types.Secp256k1,
	}
	deriveResponse, err := servicer.ConstructionDerive(ctx, &types.ConstructionDeriveRequest{
		NetworkIdentifier: networkIdentifier,
		PublicKey:         publicKey,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
		},
	}, deriveResponse)

	// Test Preprocess
	ops := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: ravencoin.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "-1000000",
				Currency: ravencoin.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "b14157a5c50503c8cd202a173613dd27e0027343c3d50cf85852dd020bf59c7f:1",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			Type: ravencoin.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "954843",
				Currency: ravencoin.TestnetCurrency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 2,
			},
			Type: ravencoin.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "44657",
				Currency: ravencoin.TestnetCurrency,
			},
		},
	}
	feeMultiplier := float64(0.75)
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&types.ConstructionPreprocessRequest{
			NetworkIdentifier:      networkIdentifier,
			Operations:             ops,
			SuggestedFeeMultiplier: &feeMultiplier,
		},
	)
	assert.Nil(t, err)
	options := &preprocessOptions{
		Coins: []*types.Coin{
			{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "b14157a5c50503c8cd202a173613dd27e0027343c3d50cf85852dd020bf59c7f:1",
				},
				Amount: &types.Amount{
					Value:    "-1000000",
					Currency: ravencoin.TestnetCurrency,
				},
			},
		},
		EstimatedSize: 220,
		FeeMultiplier: &feeMultiplier,
	}
	assert.Equal(t, &types.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, options),
	}, preprocessResponse)

	// Test Metadata
	metadata := &constructionMetadata{
		ReplayBlockHash: "0000000000007602abdc55c41b27487abbaaf017495a9f5e329d1e9c9e957675",
		ReplayBlockHeight: 212,
		ScriptPubKeys: []*ravencoin.ScriptPubKey{
			{
				ASM:          "OP_DUP OP_HASH160 6295e79e575b12e5fbf5642eb79a004025a97334 OP_EQUALVERIFY OP_CHECKSIG",
				Hex:          "76a9146295e79e575b12e5fbf5642eb79a004025a9733488ac",
				RequiredSigs: 1,
				Type:         "pubkeyhash",
				Addresses: []string{
					"RJGTst37SwJqpRBiDex3qoVa4xAH6mekJg",
				},
			},
		},
	}

	// Normal Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		ravencoin.MinFeeRate*10,
		nil,
	).Once()
	mockClient.On(
		"GetBestBlock",
		ctx,
	).Return(
		int64(312),
		nil,
	).Twice()
	mockClient.On(
		"GetHashFromIndex",
		ctx,
		int64(212),
	).Return(
		"0000000000007602abdc55c41b27487abbaaf017495a9f5e329d1e9c9e957675",
		nil,
	).Twice()
	metadataResponse, err := servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "2474999", // 3,299,999 * 0.75
				Currency: ravencoin.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Low Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		ravencoin.MinFeeRate,
		nil,
	).Once()
	metadataResponse, err = servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "330000", // we don't go below minimum fee rate
				Currency: ravencoin.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Test Payloads
	unsignedRaw := "7b227472616e73616374696f6e223a2230313030303030303031376639636635306230326464353235386638306364356333343337333032653032376464313333363137326132306364633830333035633561353537343162313031303030303030303066666666666666663032646239313065303030303030303030303139373661393134356464316433613034383131396332376232383239333035363732346439353232663236643934353838616337316165303030303030303030303030313937366139313435646431643361303438313139633237623238323933303536373234643935323266323664393435383861633030303030303030222c227363726970745075624b657973223a5b7b2261736d223a224f505f445550204f505f484153483136302036323935653739653537356231326535666266353634326562373961303034303235613937333334204f505f455155414c564552494659204f505f434845434b534947222c22686578223a223736613931343632393565373965353735623132653566626635363432656237396130303430323561393733333438386163222c2272657153696773223a312c2274797065223a227075626b657968617368222c22616464726573736573223a5b22524a47547374333753774a717052426944657833716f566134784148366d656b4a67225d7d5d2c22696e7075745f616d6f756e7473223a5b222d31303030303030225d2c22696e7075745f616464726573736573223a5b226d70353256755866544b687a59707552336a4c76504559595543577438344a374435225d7d" // nolint
	payloadsResponse, err := servicer.ConstructionPayloads(ctx, &types.ConstructionPayloadsRequest{
		NetworkIdentifier: networkIdentifier,
		Operations:        ops,
		Metadata:          forceMarshalMap(t, metadata),
	})
	assert.Nil(t, err)
	val0 := int64(0)
	val1 := int64(1)
	parseOps := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        0,
				NetworkIndex: &val0,
			},
			Type: ravencoin.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "-1000000",
				Currency: ravencoin.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "b14157a5c50503c8cd202a173613dd27e0027343c3d50cf85852dd020bf59c7f:1",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        1,
				NetworkIndex: &val0,
			},
			Type: ravencoin.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "954843",
				Currency: ravencoin.TestnetCurrency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        2,
				NetworkIndex: &val1,
			},
			Type: ravencoin.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
			},
			Amount: &types.Amount{
				Value:    "44657",
				Currency: ravencoin.TestnetCurrency,
			},
		},
	}

	signingPayload := &types.SigningPayload{
		Bytes: forceHexDecode(
			t,
			"dd512e214c1a6cfd5e9e441ff1fdce9f54bf48c0482622538057a173b8d9e325",
		),
		AccountIdentifier: &types.AccountIdentifier{
			Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5",
		},
		SignatureType: types.Ecdsa,
	}
	assert.Equal(t, &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedRaw,
		Payloads:            []*types.SigningPayload{signingPayload},
	}, payloadsResponse)

	// Test Parse Unsigned
	parseUnsignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            false,
		Transaction:       unsignedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations:               parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{},
	}, parseUnsignedResponse)

	// Test Combine
	signedRaw := "7b227472616e73616374696f6e223a22303130303030303030313766396366353062303264643532353866383063643563333433373330326530323764643133333631373261323063646338303330356335613535373431623130313030303030303661343733303434303232303235383736656338623966353164333433613561353661633534396330633832383030356566343565626539646131363664623634356330393135373232336630323230346364303862373237386138383839613831313335393135626365313064316566336262393262323137663831613064653765373966666233646664366163353031323130336230646137343937333064633962346231663461313464363930323837376139323534316635333638373738383533643963346130636237383032646366623266666666666666663032646239313065303030303030303030303139373661393134356464316433613034383131396332376232383239333035363732346439353232663236643934353838616337316165303030303030303030303030313937366139313435646431643361303438313139633237623238323933303536373234643935323266323664393435383861633030303030303030222c22696e7075745f616d6f756e7473223a5b222d31303030303030225d7d" // nolint
	combineResponse, err := servicer.ConstructionCombine(ctx, &types.ConstructionCombineRequest{
		NetworkIdentifier:   networkIdentifier,
		UnsignedTransaction: unsignedRaw,
		Signatures: []*types.Signature{
			{
				Bytes: forceHexDecode(
					t,
					"25876ec8b9f51d343a5a56ac549c0c828005ef45ebe9da166db645c09157223f4cd08b7278a8889a81135915bce10d1ef3bb92b217f81a0de7e79ffb3dfd6ac5", // nolint
				),
				SigningPayload: signingPayload,
				PublicKey:      publicKey,
				SignatureType:  types.Ecdsa,
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionCombineResponse{
		SignedTransaction: signedRaw,
	}, combineResponse)

	// Test Parse Signed
	parseSignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            true,
		Transaction:       signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations: parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{
			{Address: "mp52VuXfTKhzYpuR3jLvPEYYUCWt84J7D5"},
		},
	}, parseSignedResponse)

	// Test Hash
	transactionIdentifier := &types.TransactionIdentifier{
		Hash: "2ec3d97e6c354ee919c04ede79abf4da8ce6b3289c05ec84f0a5b6f5381cf21d",
	}
	hashResponse, err := servicer.ConstructionHash(ctx, &types.ConstructionHashRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, hashResponse)

	// Test Submit
	ravencoinTransaction := "01000000017f9cf50b02dd5258f80cd5c3437302e027dd1336172a20cdc80305c5a55741b1010000006a473044022025876ec8b9f51d343a5a56ac549c0c828005ef45ebe9da166db645c09157223f02204cd08b7278a8889a81135915bce10d1ef3bb92b217f81a0de7e79ffb3dfd6ac5012103b0da749730dc9b4b1f4a14d6902877a92541f5368778853d9c4a0cb7802dcfb2ffffffff02db910e00000000001976a9145dd1d3a048119c27b28293056724d9522f26d94588ac71ae0000000000001976a9145dd1d3a048119c27b28293056724d9522f26d94588ac00000000" // nolint
	mockClient.On(
		"SendRawTransaction",
		ctx,
		ravencoinTransaction,
	).Return(
		transactionIdentifier.Hash,
		nil,
	)
	submitResponse, err := servicer.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, submitResponse)

	mockClient.AssertExpectations(t)
	mockIndexer.AssertExpectations(t)
}
