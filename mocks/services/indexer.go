// Code generated by mockery v2.7.4. DO NOT EDIT.

package services

import (
	context "context"

	ravencoin "github.com/RavenProject/rosetta-ravencoin/ravencoin"
	mock "github.com/stretchr/testify/mock"

	types "github.com/coinbase/rosetta-sdk-go/types"
)

// Indexer is an autogenerated mock type for the Indexer type
type Indexer struct {
	mock.Mock
}

// GetBalance provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Indexer) GetBalance(_a0 context.Context, _a1 *types.AccountIdentifier, _a2 *types.Currency, _a3 *types.PartialBlockIdentifier) (*types.Amount, *types.BlockIdentifier, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *types.Amount
	if rf, ok := ret.Get(0).(func(context.Context, *types.AccountIdentifier, *types.Currency, *types.PartialBlockIdentifier) *types.Amount); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Amount)
		}
	}

	var r1 *types.BlockIdentifier
	if rf, ok := ret.Get(1).(func(context.Context, *types.AccountIdentifier, *types.Currency, *types.PartialBlockIdentifier) *types.BlockIdentifier); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*types.BlockIdentifier)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *types.AccountIdentifier, *types.Currency, *types.PartialBlockIdentifier) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetBlockLazy provides a mock function with given fields: _a0, _a1
func (_m *Indexer) GetBlockLazy(_a0 context.Context, _a1 *types.PartialBlockIdentifier) (*types.BlockResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.BlockResponse
	if rf, ok := ret.Get(0).(func(context.Context, *types.PartialBlockIdentifier) *types.BlockResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BlockResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.PartialBlockIdentifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockTransaction provides a mock function with given fields: _a0, _a1, _a2
func (_m *Indexer) GetBlockTransaction(_a0 context.Context, _a1 *types.BlockIdentifier, _a2 *types.TransactionIdentifier) (*types.Transaction, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, *types.BlockIdentifier, *types.TransactionIdentifier) *types.Transaction); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.BlockIdentifier, *types.TransactionIdentifier) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCoins provides a mock function with given fields: _a0, _a1
func (_m *Indexer) GetCoins(_a0 context.Context, _a1 *types.AccountIdentifier) ([]*types.Coin, *types.BlockIdentifier, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*types.Coin
	if rf, ok := ret.Get(0).(func(context.Context, *types.AccountIdentifier) []*types.Coin); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Coin)
		}
	}

	var r1 *types.BlockIdentifier
	if rf, ok := ret.Get(1).(func(context.Context, *types.AccountIdentifier) *types.BlockIdentifier); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*types.BlockIdentifier)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *types.AccountIdentifier) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetScriptPubKeys provides a mock function with given fields: _a0, _a1
func (_m *Indexer) GetScriptPubKeys(_a0 context.Context, _a1 []*types.Coin) ([]*ravencoin.ScriptPubKey, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*ravencoin.ScriptPubKey
	if rf, ok := ret.Get(0).(func(context.Context, []*types.Coin) []*ravencoin.ScriptPubKey); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ravencoin.ScriptPubKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*types.Coin) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
