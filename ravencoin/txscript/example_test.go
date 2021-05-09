// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript_test

import (
	"encoding/hex"
	"fmt"

	"github.com/RavenProject/rosetta-ravencoin/ravencoin/btcec"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/chaincfg"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/chaincfg/chainhash"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/txscript"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/wire"
	"github.com/RavenProject/rosetta-ravencoin/ravenutil"
)

// This example demonstrates creating a p2pk script
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func ExamplePayToAddrScript_PayToPublicKey() {
	// Parse the address to send the coins to into a ravenutil.Address
	// which is useful to ensure the accuracy of the address and determine
	// the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	//
	addressStr := "RACQw3N1wsHE3bEziCr9RZhEBJSDVdbAxb"
	address, err := ravenutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a public key script that pays to the address.
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)

	// Output:
	// Script Hex: 76a9140a10c69ec1243abcccfc6bafae1860dbfe2c80d788ac
	// Script Disassembly: OP_DUP OP_HASH160 0a10c69ec1243abcccfc6bafae1860dbfe2c80d7 OP_EQUALVERIFY OP_CHECKSIG
}

// This example demonstrates creating a p2h script
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func ExamplePayToAddrScript_PayToHash() {
	// Parse the address to send the coins to into a ravenutil.Address
	// which is useful to ensure the accuracy of the address and determine
	// the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	addressStr := "RACQw3N1wsHE3bEziCr9RZhEBJSDVdbAxb"
	address, err := ravenutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a public key script that pays to the address.
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)

	// Output:
	// Script Hex: 76a9140a10c69ec1243abcccfc6bafae1860dbfe2c80d788ac
	// Script Disassembly: OP_DUP OP_HASH160 0a10c69ec1243abcccfc6bafae1860dbfe2c80d7 OP_EQUALVERIFY OP_CHECKSIG
}

// This example demonstrates creating a p2pk script with replay protection.
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func ExamplePayToAddrReplayOutScript_PayToPublicKey() {
	// Parse the address to send the coins to into a ravenutil.Address
	// which is useful to ensure the accuracy of the address and determine
	// the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	addressStr := "RACQw3N1wsHE3bEziCr9RZhEBJSDVdbAxb"
	blockHash, errDecode := hex.DecodeString("47f1273bab0e66e76a5c2dd8fed808a23b2c08a22fcc46c00c78000400000000")
	blockHeight := int64(811742)
	address, err := ravenutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)

	if errDecode != nil {
		fmt.Println(errDecode)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a public key script that pays to the address.
	script, err := txscript.PayToAddrReplayOutScript(address, blockHash, blockHeight)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)

	// Output:
	// Script Hex: 76a9140a10c69ec1243abcccfc6bafae1860dbfe2c80d788ac20000000000400780cc046cc2fa2082c3ba208d8fed82d5c6ae7660eab3b27f14703de620cb4
	// Script Disassembly: OP_DUP OP_HASH160 0a10c69ec1243abcccfc6bafae1860dbfe2c80d7 OP_EQUALVERIFY OP_CHECKSIG 000000000400780cc046cc2fa2082c3ba208d8fed82d5c6ae7660eab3b27f147 de620c OP_CHECKBLOCKATHEIGHT
}

// This example demonstrates creating a p2h script with replay protection.
// It also prints the created script hex and uses the DisasmString function to
// display the disassembled script.
func ExamplePayToAddrReplayOutScript_PayToHash() {
	// Parse the address to send the coins to into a ravenutil.Address
	// which is useful to ensure the accuracy of the address and determine
	// the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	addressStr := "RACQw3N1wsHE3bEziCr9RZhEBJSDVdbAxb"
	blockHash, errDecode := hex.DecodeString("00000001cf4e27ce1dd8028408ed0a48edd445ba388170c9468ba0d42fff3052")
	blockHeight := int64(142091)
	address, err := ravenutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)

	if errDecode != nil {
		fmt.Println(errDecode)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a public key script that pays to the address.
	script, err := txscript.PayToAddrReplayOutScript(address, blockHash, blockHeight)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)

	// Output:
	// Script Hex: 76a9140a10c69ec1243abcccfc6bafae1860dbfe2c80d788ac205230ff2fd4a08b46c9708138ba45d4ed480aed088402d81dce274ecf01000000030b2b02b4
	// Script Disassembly: OP_DUP OP_HASH160 0a10c69ec1243abcccfc6bafae1860dbfe2c80d7 OP_EQUALVERIFY OP_CHECKSIG 5230ff2fd4a08b46c9708138ba45d4ed480aed088402d81dce274ecf01000000 0b2b02 OP_CHECKBLOCKATHEIGHT
}

// This example demonstrates extracting information from a standard p2pk script
func ExampleExtractPkScriptAddrs() {
	// Start with a standard pay-to-pubkey-hash script.
	scriptHex := "76a914567f61e0dd2ff69bc019f744b7f780ec81e87c8a88ac"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
		script, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addresses)
	fmt.Println("Required Signatures:", reqSigs)

	// Output:
	// Script Class: pubkeyhash
	// Addresses: [RHAYnhAztfvnx5f3yu56wce3poxTnwPJBQ]
	// Required Signatures: 1
}

// This example demonstrates extracting information from a standard p2pk script with replay portection.
func ExampleExtractPkScriptReplayOutAddrs() {
	// Start with a standard pay-to-pubkey-hash script with replay protection.
	scriptHex := "76a91473fff1a04bd772b89ed146d8dce862e8580579dc88ac2047f1273bab0e66e76a5c2dd8fed808a23b2c08a22fcc46c00c7800040000000003de620cb4"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
		script, &chaincfg.MainNetParams)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addresses)
	fmt.Println("Required Signatures:", reqSigs)

	// Output:
	// Script Class: pubkeyhashreplayout
	// Addresses: [RKrYPsFRHNbueoyRSqvEo8SLT8HqkfrQx7]
	// Required Signatures: 1
}

// This example demonstrates extracting information from a standard p2h script
func ExampleExtractP2sAddrs() {
	// Start with a standard p2h script
	scriptHex := "a9140ae425efe434874e2a090d44ae853caa73c703de87"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
		script, &chaincfg.MainNetParams)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addresses)
	fmt.Println("Required Signatures:", reqSigs)

	// Output:
	// Script Class: scripthash
	// Addresses: [r7EBCBTJs3bbNq6AF5jjnZaMsTT4CD9g1h]
	// Required Signatures: 1
}

// This example demonstrates extracting information from a standard p2h script with replay protection
func ExampleExtractP2sReplayOutAddrs() {
	// Start with a standard p2h script
	scriptHex := "a914cdd4b6749cf4b394c1e48c3ede7fe483512af9db872000000001cf4e27ce1dd8028408ed0a48edd445ba388170c9468ba0d42fff3052030b2b02b4"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
		script, &chaincfg.MainNetParams)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addresses)
	fmt.Println("Required Signatures:", reqSigs)

	// Output:
	// Script Class: scripthashreplayout
	// Addresses: [rQzvU7ULHTLHGDdyD5UXy77kYH3PtTaFdh]
	// Required Signatures: 1
}

// This example demonstrates manually creating and signing a redeem transaction.
func ExampleSignTxOutput() {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyHash := ravenutil.Hash160(pubKey.SerializeCompressed())
	addr, err := ravenutil.NewAddressPubKeyHash(pubKeyHash,
		&chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0})
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Sign the redeeming transaction.
	lookupKey := func(a ravenutil.Address) (*btcec.PrivateKey, bool, error) {
		// Ordinarily this function would involve looking up the private
		// key for the provided address, but since the only thing being
		// signed in this example uses the address associated with the
		// private key from above, simply return it with the compressed
		// flag set since the address is using the associated compressed
		// public key.
		//
		// NOTE: If you want to prove the code is actually signing the
		// transaction properly, uncomment the following line which
		// intentionally returns an invalid key to sign with, which in
		// turn will result in a failure during the script execution
		// when verifying the signature.
		//
		// privKey.D.SetInt64(12345)
		//
		return privKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(&chaincfg.MainNetParams,
		redeemTx, 0, originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	redeemTx.TxIn[0].SignatureScript = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Transaction successfully signed")

	// Output:
	// Transaction successfully signed
}

// This example demonstrates manually creating and signing a redeem transaction.
func ExampleRavenSignTxOutput() {
	blockHash, err := hex.DecodeString("0da5ee723b7923feb580518541c6f098206330dbc711a6678922c11f2ccf1abb")
	if err != nil {
		fmt.Println(err)
		return
	}
	blockHeight := int64(0)

	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("2c3a48576fe6e8a466e78cd2957c9dc62128135540bbea0685d7c4a23ea35a6c")
	if err != nil {
		fmt.Println(err)
		return
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyHash := ravenutil.Hash160(pubKey.SerializeCompressed())
	addr, err := ravenutil.NewAddressPubKeyHash(pubKeyHash,
		&chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0})
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrReplayOutScript(addr, blockHash, blockHeight)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Sign the redeeming transaction.
	lookupKey := func(a ravenutil.Address) (*btcec.PrivateKey, bool, error) {
		// Ordinarily this function would involve looking up the private
		// key for the provided address, but since the only thing being
		// signed in this example uses the address associated with the
		// private key from above, simply return it with the compressed
		// flag set since the address is using the associated compressed
		// public key.
		//
		// NOTE: If you want to prove the code is actually signing the
		// transaction properly, uncomment the following line which
		// intentionally returns an invalid key to sign with, which in
		// turn will result in a failure during the script execution
		// when verifying the signature.
		//
		// privKey.D.SetInt64(12345)
		//
		return privKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(&chaincfg.MainNetParams,
		redeemTx, 0, originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	redeemTx.TxIn[0].SignatureScript = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, nil, nil, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Transaction successfully signed")

	// Output:
	// Transaction successfully signed
}
