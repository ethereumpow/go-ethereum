// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"
	"runtime"

	"github.com/ethereum/go-ethereum/core/types"
)

type cachedTransaction struct {
	tx          *types.Transaction
	blockNumber *big.Int
}

// senderCacher is a concurrent transaction sender recoverer and cacher.
var senderCacher = newTxSenderCacher(runtime.NumCPU())

// txSenderCacherRequest is a request for recovering transaction senders with a
// specific signature scheme and caching it into the transactions themselves.
//
// The inc field defines the number of transactions to skip after each recovery,
// which is used to feed the same underlying input array to different threads but
// ensure they process the early transactions fast.
type txSenderCacherRequest struct {
	signer types.Signer
	txs    []*cachedTransaction
	inc    int
}

// txSenderCacher is a helper structure to concurrently ecrecover transaction
// senders from digital signatures on background threads.
type txSenderCacher struct {
	threads int
	tasks   chan *txSenderCacherRequest
}

// newTxSenderCacher creates a new transaction sender background cacher and starts
// as many processing goroutines as allowed by the GOMAXPROCS on construction.
func newTxSenderCacher(threads int) *txSenderCacher {
	cacher := &txSenderCacher{
		tasks:   make(chan *txSenderCacherRequest, threads),
		threads: threads,
	}
	for i := 0; i < threads; i++ {
		go cacher.cache()
	}
	return cacher
}

// cache is an infinite loop, caching transaction senders from various forms of
// data structures.
func (cacher *txSenderCacher) cache() {
	for task := range cacher.tasks {
		for i := 0; i < len(task.txs); i += task.inc {
			currentTx := task.txs[i]
			if currentTx != nil {
				types.Sender(task.signer, currentTx.tx, currentTx.blockNumber)
			}
		}
	}
}

// recover recovers the senders from a batch of transactions and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) recoverWithBlockNumber(signer types.Signer, txs []*types.Transaction, blockNumber *big.Int) {
	// If there's nothing to recover, abort
	if len(txs) == 0 {
		return
	}
	cachedTxs := make([]*cachedTransaction, len(txs))
	for _, tx := range txs {
		cachedTxs = append(cachedTxs, &cachedTransaction{tx: tx, blockNumber: blockNumber})
	}
	cacher.recover(signer, cachedTxs)
}

// recover recovers the senders from a batch of transactions and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) recover(signer types.Signer, txs []*cachedTransaction) {
	// If there's nothing to recover, abort
	if len(txs) == 0 {
		return
	}
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}

	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs[:i],
			inc:    tasks,
		}
	}
}

// recoverFromBlocks recovers the senders from a batch of blocks and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) recoverFromBlocks(signer types.Signer, blocks []*types.Block) {
	count := 0
	for _, block := range blocks {
		count += len(block.Transactions())
	}
	txs := make([]*cachedTransaction, 0, count)
	for _, block := range blocks {
		for _, tx := range block.Transactions() {
			txs = append(txs, &cachedTransaction{tx: tx, blockNumber: block.Number()})
		}
	}
	cacher.recover(signer, txs)
}
