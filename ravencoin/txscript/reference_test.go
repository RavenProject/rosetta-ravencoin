// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/RavenProject/rosetta-ravencoin/ravencoin/btcec"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/chaincfg/chainhash"
	"github.com/RavenProject/rosetta-ravencoin/ravencoin/wire"
	"github.com/RavenProject/rosetta-ravencoin/ravenutil"
)

// scriptTestName returns a descriptive test name for the given reference script
// test data.
func scriptTestName(test []interface{}) (string, error) {
	// Account for any optional leading witness data.
	var witnessOffset int
	if _, ok := test[0].([]interface{}); ok {
		witnessOffset++
	}

	// In addition to the optional leading witness data, the test must
	// consist of at least a signature script, public key script, flags,
	// and expected error.  Finally, it may optionally contain a comment.
	if len(test) < witnessOffset+4 || len(test) > witnessOffset+5 {
		return "", fmt.Errorf("invalid test length %d", len(test))
	}

	// Use the comment for the test name if one is specified, otherwise,
	// construct the name based on the signature script, public key script,
	// and flags.
	var name string
	if len(test) == witnessOffset+5 {
		name = fmt.Sprintf("test (%s)", test[witnessOffset+4])
	} else {
		name = fmt.Sprintf("test ([%s, %s, %s])", test[witnessOffset],
			test[witnessOffset+1], test[witnessOffset+2])
	}
	return name, nil
}

// parse hex string into a []byte.
func parseHex(tok string) ([]byte, error) {
	if !strings.HasPrefix(tok, "0x") {
		return nil, errors.New("not a hex number")
	}
	return hex.DecodeString(tok[2:])
}

// parseWitnessStack parses a json array of witness items encoded as hex into a
// slice of witness elements.
func parseWitnessStack(elements []interface{}) ([][]byte, error) {
	witness := make([][]byte, len(elements))
	for i, e := range elements {
		witElement, err := hex.DecodeString(e.(string))
		if err != nil {
			return nil, err
		}

		witness[i] = witElement
	}

	return witness, nil
}

// shortFormOps holds a map of opcode names to values for use in short form
// parsing.  It is declared here so it only needs to be created once.
var shortFormOps map[string]byte

// parseShortForm parses a string as as used in the Bitcoin Core reference tests
// into the script it came from.
//
// The format used for these tests is pretty simple if ad-hoc:
//   - Opcodes other than the push opcodes and unknown are present as
//     either OP_NAME or just NAME
//   - Plain numbers are made into push operations
//   - Numbers beginning with 0x are inserted into the []byte as-is (so
//     0x14 is OP_DATA_20)
//   - Single quoted strings are pushed as data
//   - Anything else is an error
func parseShortForm(script string) ([]byte, error) {
	// Only create the short form opcode map once.
	if shortFormOps == nil {
		ops := make(map[string]byte)
		for opcodeName, opcodeValue := range OpcodeByName {
			if strings.Contains(opcodeName, "OP_UNKNOWN") {
				continue
			}
			ops[opcodeName] = opcodeValue

			// The opcodes named OP_# can't have the OP_ prefix
			// stripped or they would conflict with the plain
			// numbers.  Also, since OP_FALSE and OP_TRUE are
			// aliases for the OP_0, and OP_1, respectively, they
			// have the same value, so detect those by name and
			// allow them.
			if (opcodeName == "OP_FALSE" || opcodeName == "OP_TRUE") ||
				(opcodeValue != OP_0 && (opcodeValue < OP_1 ||
					opcodeValue > OP_16)) {

				ops[strings.TrimPrefix(opcodeName, "OP_")] = opcodeValue
			}
		}
		shortFormOps = ops
	}

	// Split only does one separator so convert all \n and tab into  space.
	script = strings.Replace(script, "\n", " ", -1)
	script = strings.Replace(script, "\t", " ", -1)
	tokens := strings.Split(script, " ")
	builder := NewScriptBuilder()

	for _, tok := range tokens {
		if len(tok) == 0 {
			continue
		}
		// if parses as a plain number
		if num, err := strconv.ParseInt(tok, 10, 64); err == nil {
			builder.AddInt64(num)
			continue
		} else if bts, err := parseHex(tok); err == nil {
			// Concatenate the bytes manually since the test code
			// intentionally creates scripts that are too large and
			// would cause the builder to error otherwise.
			if builder.err == nil {
				builder.script = append(builder.script, bts...)
			}
		} else if len(tok) >= 2 &&
			tok[0] == '\'' && tok[len(tok)-1] == '\'' {
			builder.AddFullData([]byte(tok[1 : len(tok)-1]))
		} else if opcode, ok := shortFormOps[tok]; ok {
			builder.AddOp(opcode)
		} else {
			return nil, fmt.Errorf("bad token %q", tok)
		}

	}
	return builder.Script()
}

// parseScriptFlags parses the provided flags string from the format used in the
// reference tests into ScriptFlags suitable for use in the script engine.
func parseScriptFlags(flagStr string) (ScriptFlags, error) {
	var flags ScriptFlags

	sFlags := strings.Split(flagStr, ",")
	for _, flag := range sFlags {
		switch flag {
		case "":
			// Nothing.
		case "CHECKLOCKTIMEVERIFY":
			flags |= ScriptVerifyCheckLockTimeVerify
		case "CHECKSEQUENCEVERIFY":
			flags |= ScriptVerifyCheckSequenceVerify
		case "CLEANSTACK":
			flags |= ScriptVerifyCleanStack
		case "DERSIG":
			flags |= ScriptVerifyDERSignatures
		case "DISCOURAGE_UPGRADABLE_NOPS":
			flags |= ScriptDiscourageUpgradableNops
		case "LOW_S":
			flags |= ScriptVerifyLowS
		case "MINIMALDATA":
			flags |= ScriptVerifyMinimalData
		case "NONE":
			// Nothing.
		case "NULLDUMMY":
			flags |= ScriptStrictMultiSig
		case "NULLFAIL":
			flags |= ScriptVerifyNullFail
		case "P2SH":
			flags |= ScriptBip16
		case "SIGPUSHONLY":
			flags |= ScriptVerifySigPushOnly
		case "STRICTENC":
			flags |= ScriptVerifyStrictEncoding
		case "WITNESS":
			flags |= ScriptVerifyWitness
		case "DISCOURAGE_UPGRADABLE_WITNESS_PROGRAM":
			flags |= ScriptVerifyDiscourageUpgradeableWitnessProgram
		case "MINIMALIF":
			flags |= ScriptVerifyMinimalIf
		case "WITNESS_PUBKEYTYPE":
			flags |= ScriptVerifyWitnessPubKeyType
		default:
			return flags, fmt.Errorf("invalid flag: %s", flag)
		}
	}
	return flags, nil
}

// parseExpectedResult parses the provided expected result string into allowed
// script error codes.  An error is returned if the expected result string is
// not supported.
func parseExpectedResult(expected string) ([]ErrorCode, error) {
	switch expected {
	case "OK":
		return nil, nil
	case "UNKNOWN_ERROR":
		return []ErrorCode{ErrNumberTooBig, ErrMinimalData}, nil
	case "PUBKEYTYPE":
		return []ErrorCode{ErrPubKeyType}, nil
	case "SIG_DER":
		return []ErrorCode{ErrSigTooShort, ErrSigTooLong,
			ErrSigInvalidSeqID, ErrSigInvalidDataLen, ErrSigMissingSTypeID,
			ErrSigMissingSLen, ErrSigInvalidSLen,
			ErrSigInvalidRIntID, ErrSigZeroRLen, ErrSigNegativeR,
			ErrSigTooMuchRPadding, ErrSigInvalidSIntID,
			ErrSigZeroSLen, ErrSigNegativeS, ErrSigTooMuchSPadding,
			ErrInvalidSigHashType}, nil
	case "EVAL_FALSE":
		return []ErrorCode{ErrEvalFalse, ErrEmptyStack}, nil
	case "EQUALVERIFY":
		return []ErrorCode{ErrEqualVerify}, nil
	case "NULLFAIL":
		return []ErrorCode{ErrNullFail}, nil
	case "SIG_HIGH_S":
		return []ErrorCode{ErrSigHighS}, nil
	case "SIG_HASHTYPE":
		return []ErrorCode{ErrInvalidSigHashType}, nil
	case "SIG_NULLDUMMY":
		return []ErrorCode{ErrSigNullDummy}, nil
	case "SIG_PUSHONLY":
		return []ErrorCode{ErrNotPushOnly}, nil
	case "CLEANSTACK":
		return []ErrorCode{ErrCleanStack}, nil
	case "BAD_OPCODE":
		return []ErrorCode{ErrReservedOpcode, ErrMalformedPush}, nil
	case "UNBALANCED_CONDITIONAL":
		return []ErrorCode{ErrUnbalancedConditional,
			ErrInvalidStackOperation}, nil
	case "OP_RETURN":
		return []ErrorCode{ErrEarlyReturn}, nil
	case "VERIFY":
		return []ErrorCode{ErrVerify}, nil
	case "INVALID_STACK_OPERATION", "INVALID_ALTSTACK_OPERATION":
		return []ErrorCode{ErrInvalidStackOperation}, nil
	case "DISABLED_OPCODE":
		return []ErrorCode{ErrDisabledOpcode}, nil
	case "DISCOURAGE_UPGRADABLE_NOPS":
		return []ErrorCode{ErrDiscourageUpgradableNOPs}, nil
	case "PUSH_SIZE":
		return []ErrorCode{ErrElementTooBig}, nil
	case "OP_COUNT":
		return []ErrorCode{ErrTooManyOperations}, nil
	case "STACK_SIZE":
		return []ErrorCode{ErrStackOverflow}, nil
	case "SCRIPT_SIZE":
		return []ErrorCode{ErrScriptTooBig}, nil
	case "PUBKEY_COUNT":
		return []ErrorCode{ErrInvalidPubKeyCount}, nil
	case "SIG_COUNT":
		return []ErrorCode{ErrInvalidSignatureCount}, nil
	case "MINIMALDATA":
		return []ErrorCode{ErrMinimalData}, nil
	case "NEGATIVE_LOCKTIME":
		return []ErrorCode{ErrNegativeLockTime}, nil
	case "UNSATISFIED_LOCKTIME":
		return []ErrorCode{ErrUnsatisfiedLockTime}, nil
	case "MINIMALIF":
		return []ErrorCode{ErrMinimalIf}, nil
	case "DISCOURAGE_UPGRADABLE_WITNESS_PROGRAM":
		return []ErrorCode{ErrDiscourageUpgradableWitnessProgram}, nil
	case "WITNESS_PROGRAM_WRONG_LENGTH":
		return []ErrorCode{ErrWitnessProgramWrongLength}, nil
	case "WITNESS_PROGRAM_WITNESS_EMPTY":
		return []ErrorCode{ErrWitnessProgramEmpty}, nil
	case "WITNESS_PROGRAM_MISMATCH":
		return []ErrorCode{ErrWitnessProgramMismatch}, nil
	case "WITNESS_MALLEATED":
		return []ErrorCode{ErrWitnessMalleated}, nil
	case "WITNESS_MALLEATED_P2SH":
		return []ErrorCode{ErrWitnessMalleatedP2SH}, nil
	case "WITNESS_UNEXPECTED":
		return []ErrorCode{ErrWitnessUnexpected}, nil
	case "WITNESS_PUBKEYTYPE":
		return []ErrorCode{ErrWitnessPubKeyType}, nil
	}

	return nil, fmt.Errorf("unrecognized expected result in test data: %v",
		expected)
}

// createSpendTx generates a basic spending transaction given the passed
// signature, witness and public key scripts.
func createSpendingTx(sigScript, pkScript []byte,
	outputValue int64) *wire.MsgTx {

	coinbaseTx := wire.NewMsgTx(wire.TxVersion)

	outPoint := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(outPoint, []byte{OP_0, OP_0})
	txOut := wire.NewTxOut(outputValue, pkScript)
	coinbaseTx.AddTxIn(txIn)
	coinbaseTx.AddTxOut(txOut)

	spendingTx := wire.NewMsgTx(wire.TxVersion)
	coinbaseTxSha := coinbaseTx.TxHash()
	outPoint = wire.NewOutPoint(&coinbaseTxSha, 0)
	txIn = wire.NewTxIn(outPoint, sigScript)
	txOut = wire.NewTxOut(outputValue, nil)

	spendingTx.AddTxIn(txIn)
	spendingTx.AddTxOut(txOut)

	return spendingTx
}

// scriptWithInputVal wraps a target pkScript with the value of the output in
// which it is contained. The inputVal is necessary in order to properly
// validate inputs which spend nested, or native witness programs.
type scriptWithInputVal struct {
	inputVal int64
	pkScript []byte
}

// testScripts ensures all of the passed script tests execute with the expected
// results with or without using a signature cache, as specified by the
// parameter.
func testScripts(t *testing.T, tests [][]interface{}, useSigCache bool) {
	// Create a signature cache to use only if requested.
	var sigCache *SigCache
	if useSigCache {
		sigCache = NewSigCache(10)
	}

	for i, test := range tests {
		// "Format is: [[wit..., amount]?, scriptSig, scriptPubKey,
		//    flags, expected_scripterror, ... comments]"

		// Skip single line comments.
		if len(test) == 1 {
			continue
		}

		// Construct a name for the test based on the comment and test
		// data.
		name, err := scriptTestName(test)
		if err != nil {
			t.Errorf("TestScripts: invalid test #%d: %v", i, err)
			continue
		}

		var (
			inputAmt ravenutil.Amount
		)

		// When the first field of the test data is a slice it contains
		// witness data and everything else is offset by 1 as a result.
		witnessOffset := 0
		if witnessData, ok := test[0].([]interface{}); ok {
			witnessOffset++

			inputAmt, err = ravenutil.NewAmount(witnessData[len(witnessData)-1].(float64))
			if err != nil {
				t.Errorf("%s: can't parse input amt: %v",
					name, err)
				continue
			}

		}

		// Extract and parse the signature script from the test fields.
		scriptSigStr, ok := test[witnessOffset].(string)
		if !ok {
			t.Errorf("%s: signature script is not a string", name)
			continue
		}
		scriptSig, err := parseShortForm(scriptSigStr)
		if err != nil {
			t.Errorf("%s: can't parse signature script: %v", name,
				err)
			continue
		}

		// Extract and parse the public key script from the test fields.
		scriptPubKeyStr, ok := test[witnessOffset+1].(string)
		if !ok {
			t.Errorf("%s: public key script is not a string", name)
			continue
		}
		scriptPubKey, err := parseShortForm(scriptPubKeyStr)
		if err != nil {
			t.Errorf("%s: can't parse public key script: %v", name,
				err)
			continue
		}

		// Extract and parse the script flags from the test fields.
		flagsStr, ok := test[witnessOffset+2].(string)
		if !ok {
			t.Errorf("%s: flags field is not a string", name)
			continue
		}
		flags, err := parseScriptFlags(flagsStr)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}

		// Extract and parse the expected result from the test fields.
		//
		// Convert the expected result string into the allowed script
		// error codes.  This is necessary because txscript is more
		// fine grained with its errors than the reference test data, so
		// some of the reference test data errors map to more than one
		// possibility.
		resultStr, ok := test[witnessOffset+3].(string)
		if !ok {
			t.Errorf("%s: result field is not a string", name)
			continue
		}
		allowedErrorCodes, err := parseExpectedResult(resultStr)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}

		// Generate a transaction pair such that one spends from the
		// other and the provided signature and public key scripts are
		// used, then create a new engine to execute the scripts.
		tx := createSpendingTx(scriptSig, scriptPubKey, int64(inputAmt))
		vm, err := NewEngine(scriptPubKey, tx, 0, flags, sigCache, nil, int64(inputAmt))
		if err == nil {
			err = vm.Execute()
		}

		// Ensure there were no errors when the expected result is OK.
		if resultStr == "OK" {
			if err != nil {
				t.Errorf("%s failed to execute: %v", name, err)
			}
			continue
		}

		// At this point an error was expected so ensure the result of
		// the execution matches it.
		success := false
		for _, code := range allowedErrorCodes {
			if IsErrorCode(err, code) {
				success = true
				break
			}
		}
		if !success {
			if serr, ok := err.(Error); ok {
				t.Errorf("%s: want error codes %v, got %v", name,
					allowedErrorCodes, serr.ErrorCode)
				continue
			}
			t.Errorf("%s: want error codes %v, got err: %v (%T)",
				name, allowedErrorCodes, err, err)
			continue
		}
	}
}

// TestScripts ensures all of the tests in script_tests.json execute with the
// expected results as defined in the test data.
func TestScripts(t *testing.T) {
	file, err := ioutil.ReadFile("data/script_tests.json")
	if err != nil {
		t.Fatalf("TestScripts: %v\n", err)
	}

	var tests [][]interface{}
	err = json.Unmarshal(file, &tests)
	if err != nil {
		t.Fatalf("TestScripts couldn't Unmarshal: %v", err)
	}

	// Run all script tests with and without the signature cache.
	testScripts(t, tests, true)
	testScripts(t, tests, false)
}

// testVecF64ToUint32 properly handles conversion of float64s read from the JSON
// test data to unsigned 32-bit integers.  This is necessary because some of the
// test data uses -1 as a shortcut to mean max uint32 and direct conversion of a
// negative float to an unsigned int is implementation dependent and therefore
// doesn't result in the expected value on all platforms.  This function woks
// around that limitation by converting to a 32-bit signed integer first and
// then to a 32-bit unsigned integer which results in the expected behavior on
// all platforms.
func testVecF64ToUint32(f float64) uint32 {
	return uint32(int32(f))
}

// TestTxInvalidTests ensures all of the tests in tx_invalid.json fail as
// expected.
func TestTxInvalidTests(t *testing.T) {
	file, err := ioutil.ReadFile("data/tx_invalid.json")
	if err != nil {
		t.Fatalf("TestTxInvalidTests: %v\n", err)
	}

	var tests [][]interface{}
	err = json.Unmarshal(file, &tests)
	if err != nil {
		t.Fatalf("TestTxInvalidTests couldn't Unmarshal: %v\n", err)
	}

	// form is either:
	//   ["this is a comment "]
	// or:
	//   [[[previous hash, previous index, previous scriptPubKey]...,]
	//	serializedTransaction, verifyFlags]
testloop:
	for i, test := range tests {
		inputs, ok := test[0].([]interface{})
		if !ok {
			continue
		}

		if len(test) != 3 {
			t.Errorf("bad test (bad length) %d: %v", i, test)
			continue

		}
		serializedhex, ok := test[1].(string)
		if !ok {
			t.Errorf("bad test (arg 2 not string) %d: %v", i, test)
			continue
		}
		serializedTx, err := hex.DecodeString(serializedhex)
		if err != nil {
			t.Errorf("bad test (arg 2 not hex %v) %d: %v", err, i,
				test)
			continue
		}

		tx, err := ravenutil.NewTxFromBytes(serializedTx)
		if err != nil {
			t.Errorf("bad test (arg 2 not msgtx %v) %d: %v", err,
				i, test)
			continue
		}

		verifyFlags, ok := test[2].(string)
		if !ok {
			t.Errorf("bad test (arg 3 not string) %d: %v", i, test)
			continue
		}

		flags, err := parseScriptFlags(verifyFlags)
		if err != nil {
			t.Errorf("bad test %d: %v", i, err)
			continue
		}

		prevOuts := make(map[wire.OutPoint]scriptWithInputVal)
		for j, iinput := range inputs {
			input, ok := iinput.([]interface{})
			if !ok {
				t.Errorf("bad test (%dth input not array)"+
					"%d: %v", j, i, test)
				continue testloop
			}

			if len(input) < 3 || len(input) > 4 {
				t.Errorf("bad test (%dth input wrong length)"+
					"%d: %v", j, i, test)
				continue testloop
			}

			previoustx, ok := input[0].(string)
			if !ok {
				t.Errorf("bad test (%dth input hash not string)"+
					"%d: %v", j, i, test)
				continue testloop
			}

			prevhash, err := chainhash.NewHashFromStr(previoustx)
			if err != nil {
				t.Errorf("bad test (%dth input hash not hash %v)"+
					"%d: %v", j, err, i, test)
				continue testloop
			}

			idxf, ok := input[1].(float64)
			if !ok {
				t.Errorf("bad test (%dth input idx not number)"+
					"%d: %v", j, i, test)
				continue testloop
			}
			idx := testVecF64ToUint32(idxf)

			oscript, ok := input[2].(string)
			if !ok {
				t.Errorf("bad test (%dth input script not "+
					"string) %d: %v", j, i, test)
				continue testloop
			}

			script, err := parseShortForm(oscript)
			if err != nil {
				t.Errorf("bad test (%dth input script doesn't "+
					"parse %v) %d: %v", j, err, i, test)
				continue testloop
			}

			var inputValue float64
			if len(input) == 4 {
				inputValue, ok = input[3].(float64)
				if !ok {
					t.Errorf("bad test (%dth input value not int) "+
						"%d: %v", j, i, test)
					continue
				}
			}

			v := scriptWithInputVal{
				inputVal: int64(inputValue),
				pkScript: script,
			}
			prevOuts[*wire.NewOutPoint(prevhash, idx)] = v
		}

		for k, txin := range tx.MsgTx().TxIn {
			prevOut, ok := prevOuts[txin.PreviousOutPoint]
			if !ok {
				t.Errorf("bad test (missing %dth input) %d:%v",
					k, i, test)
				continue testloop
			}
			// These are meant to fail, so as soon as the first
			// input fails the transaction has failed. (some of the
			// test txns have good inputs, too..
			vm, err := NewEngine(prevOut.pkScript, tx.MsgTx(), k,
				flags, nil, nil, prevOut.inputVal)
			if err != nil {
				continue testloop
			}

			err = vm.Execute()
			if err != nil {
				continue testloop
			}

		}
		t.Errorf("test (%d:%v) succeeded when should fail",
			i, test)
	}
}

// TestTxValidTests ensures all of the tests in tx_valid.json pass as expected.
func TestTxValidTests(t *testing.T) {
	file, err := ioutil.ReadFile("data/tx_valid.json")
	if err != nil {
		t.Fatalf("TestTxValidTests: %v\n", err)
	}

	var tests [][]interface{}
	err = json.Unmarshal(file, &tests)
	if err != nil {
		t.Fatalf("TestTxValidTests couldn't Unmarshal: %v\n", err)
	}

	// form is either:
	//   ["this is a comment "]
	// or:
	//   [[[previous hash, previous index, previous scriptPubKey, input value]...,]
	//	serializedTransaction, verifyFlags]
testloop:
	for i, test := range tests {
		inputs, ok := test[0].([]interface{})
		if !ok {
			continue
		}

		if len(test) != 3 {
			t.Errorf("bad test (bad length) %d: %v", i, test)
			continue
		}
		serializedhex, ok := test[1].(string)
		if !ok {
			t.Errorf("bad test (arg 2 not string) %d: %v", i, test)
			continue
		}
		serializedTx, err := hex.DecodeString(serializedhex)
		if err != nil {
			t.Errorf("bad test (arg 2 not hex %v) %d: %v", err, i,
				test)
			continue
		}

		tx, err := ravenutil.NewTxFromBytes(serializedTx)
		if err != nil {
			t.Errorf("bad test (arg 2 not msgtx %v) %d: %v", err,
				i, test)
			continue
		}

		verifyFlags, ok := test[2].(string)
		if !ok {
			t.Errorf("bad test (arg 3 not string) %d: %v", i, test)
			continue
		}

		flags, err := parseScriptFlags(verifyFlags)
		if err != nil {
			t.Errorf("bad test %d: %v", i, err)
			continue
		}

		prevOuts := make(map[wire.OutPoint]scriptWithInputVal)
		for j, iinput := range inputs {
			input, ok := iinput.([]interface{})
			if !ok {
				t.Errorf("bad test (%dth input not array)"+
					"%d: %v", j, i, test)
				continue
			}

			if len(input) < 3 || len(input) > 4 {
				t.Errorf("bad test (%dth input wrong length)"+
					"%d: %v", j, i, test)
				continue
			}

			previoustx, ok := input[0].(string)
			if !ok {
				t.Errorf("bad test (%dth input hash not string)"+
					"%d: %v", j, i, test)
				continue
			}

			prevhash, err := chainhash.NewHashFromStr(previoustx)
			if err != nil {
				t.Errorf("bad test (%dth input hash not hash %v)"+
					"%d: %v", j, err, i, test)
				continue
			}

			idxf, ok := input[1].(float64)
			if !ok {
				t.Errorf("bad test (%dth input idx not number)"+
					"%d: %v", j, i, test)
				continue
			}
			idx := testVecF64ToUint32(idxf)

			oscript, ok := input[2].(string)
			if !ok {
				t.Errorf("bad test (%dth input script not "+
					"string) %d: %v", j, i, test)
				continue
			}

			script, err := parseShortForm(oscript)
			if err != nil {
				t.Errorf("bad test (%dth input script doesn't "+
					"parse %v) %d: %v", j, err, i, test)
				continue
			}

			var inputValue float64
			if len(input) == 4 {
				inputValue, ok = input[3].(float64)
				if !ok {
					t.Errorf("bad test (%dth input value not int) "+
						"%d: %v", j, i, test)
					continue
				}
			}

			v := scriptWithInputVal{
				inputVal: int64(inputValue),
				pkScript: script,
			}
			prevOuts[*wire.NewOutPoint(prevhash, idx)] = v
		}

		for k, txin := range tx.MsgTx().TxIn {
			prevOut, ok := prevOuts[txin.PreviousOutPoint]
			if !ok {
				t.Errorf("bad test (missing %dth input) %d:%v",
					k, i, test)
				continue testloop
			}
			vm, err := NewEngine(prevOut.pkScript, tx.MsgTx(), k,
				flags, nil, nil, prevOut.inputVal)
			if err != nil {
				t.Errorf("test (%d:%v:%d) failed to create "+
					"script: %v", i, test, k, err)
				continue
			}

			err = vm.Execute()
			if err != nil {
				t.Errorf("test (%d:%v:%d) failed to execute: "+
					"%v", i, test, k, err)
				continue
			}
		}
	}
}

func TestRavenCalcSignatureHash_SigHashAll(t *testing.T) {
	unsignedSerializedTx := "0100000001383dedb49935f49f5ef93b6b007d468c2a337cfe6f5dc0af62a151a4192198590000000000ffffffff0140420f00000000003f76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac205230ff2fd4a08b46c9708138ba45d4ed480aed088402d81dce274ecf01000000030b2b02b400000000"
	scriptHex := "76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac20c243be1a6b3d319e40e89b159235a320a1cd50d35c2e52bc79e94b990100000003d92c02b4"

	var tx wire.MsgTx
	rawTx, _ := hex.DecodeString(unsignedSerializedTx)
	err := tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse transaction: %v", err)
	}

	subScript, _ := hex.DecodeString(scriptHex)
	parsedScript, err := parseScript(subScript)
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse sub-script: %v", err)
	}

	hash := calcSignatureHash(parsedScript, SigHashAll, &tx, 0)
	if hex.EncodeToString(hash) != "7c7bfbdc8f8592bc3b618a3c61e9e58b562f6d668307602a3edc88cadf6d1113" {
		t.Errorf("TestRavenCalcSignatureHash failed test: wrong hash")
	}
}

func TestRavenCalcSignatureHash_SigHashSingle(t *testing.T) {
	unsignedSerializedTx := "0100000001383dedb49935f49f5ef93b6b007d468c2a337cfe6f5dc0af62a151a4192198590000000000ffffffff0140420f00000000003f76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac205230ff2fd4a08b46c9708138ba45d4ed480aed088402d81dce274ecf01000000030b2b02b400000000"
	scriptHex := "76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac20c243be1a6b3d319e40e89b159235a320a1cd50d35c2e52bc79e94b990100000003d92c02b4"

	var tx wire.MsgTx
	rawTx, _ := hex.DecodeString(unsignedSerializedTx)
	err := tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse transaction: %v", err)
	}

	subScript, _ := hex.DecodeString(scriptHex)
	parsedScript, err := parseScript(subScript)
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse sub-script: %v", err)
	}

	hash := calcSignatureHash(parsedScript, SigHashSingle, &tx, 0)
	if hex.EncodeToString(hash) != "1fbdbdbeb088c68a5780b883e54f3d923704a15ddc886a8c1eeb710992be56a0" {
		t.Errorf("TestRavenCalcSignatureHash failed test: wrong hash")
	}
}

func TestRavenCalcSignatureHash_SigHashNone(t *testing.T) {
	unsignedSerializedTx := "0100000001383dedb49935f49f5ef93b6b007d468c2a337cfe6f5dc0af62a151a4192198590000000000ffffffff0140420f00000000003f76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac205230ff2fd4a08b46c9708138ba45d4ed480aed088402d81dce274ecf01000000030b2b02b400000000"
	scriptHex := "76a914da46f44467949ac9321b16402c32bbeede5e3e5f88ac20c243be1a6b3d319e40e89b159235a320a1cd50d35c2e52bc79e94b990100000003d92c02b4"

	var tx wire.MsgTx
	rawTx, _ := hex.DecodeString(unsignedSerializedTx)
	err := tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse transaction: %v", err)
	}

	subScript, _ := hex.DecodeString(scriptHex)
	parsedScript, err := parseScript(subScript)
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse sub-script: %v", err)
	}

	hash := calcSignatureHash(parsedScript, SigHashNone, &tx, 0)
	if hex.EncodeToString(hash) != "de80af8f03f41ac8c5027514b8ce86cb64b67a3d08210cc3b55e84c5f7f07cb6" {
		t.Errorf("TestRavenCalcSignatureHash failed test: wrong hash")
	}
}

func TestRavenRawTxInSignature_SigHashAll(t *testing.T) {
	unsignedSerializedTx := "0100000001dc57f92c814871a404ae667003c00da6fcdbe5a1e16610fec4337d74db23fe4f0100000000ffffffff0100ca9a3b000000003e76a91484754e80eeb0a77a3895eacc233d93acc8a689ff88ac202362d1e7b9ffcde357d1c83e60875c7ddbd46b92c3031abb77ab0b7c8c768d0c024101b400000000"
	scriptHex := "76a9145f3e58d7136483ac14c305c21f15267ac7c3233688ac20bb1acf2c1fc1228967a611c7db30632098f0c641855180b5fe23793b72eea50d00b4"
	expectedSign := "3045022100f6b966be25f737f7116ae270b8f00f6f1f9330587f6d3831d0d7109b5fb97e1f02204f65396c964d4f8f5e1b283e500fa0bf03f7f7ef47cd49e3497046b27cd795ae01"

	var tx wire.MsgTx
	rawTx, _ := hex.DecodeString(unsignedSerializedTx)
	err := tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		t.Errorf("TestRavenCalcSignatureHash failed test: "+
			"Failed to parse transaction: %v", err)
	}

	subScript, _ := hex.DecodeString(scriptHex)
	privKeyBytes, _ := hex.DecodeString("2791eaa5b6200f6654a947af9e64191644f1db55de87fdecd8ee4aa1074b7f8a")
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)

	sig, _ := RawTxInSignature(&tx, 0, subScript, SigHashAll, privKey)

	if hex.EncodeToString(sig) != expectedSign {
		t.Errorf("TestRavenCalcSignatureHash failed test: wrong hash")
	}

}
