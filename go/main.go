package main

import (
	"context"
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/gas"
)

const (
	wif1 = "L557bcBh8mXTp68ABb8Bt4VzyMPAgQeyXib3s2KRot1TCproLGDV"
	wif2 = "L1L32AoWKitygC3HFDv3SWt6ndq78jLJiD7PYsARPZUFmBoRdLge"
)

func logTransfer(t result.NEP17Transfer) {
	if t.Asset != gas.Hash {
		return
	}
	var addr string
	if t.Address == "" {
		addr = "<NIL> "
	} else {
		addr = t.Address[1:7]
	}

	fmt.Println(addr, t.Amount, t.Index, t.Timestamp, t.TxHash)
}

func main() {
	const endpoint = "http://localhost:20331"

	c, err := rpcclient.New(context.TODO(), endpoint, rpcclient.Options{})
	if err != nil {
		panic(err)
	} else if err = c.Init(); err != nil {
		panic(err)
	} else if err = c.Ping(); err != nil {
		panic(err)
	}

	acc1, err := AccountFromWIF(wif1)
	if err != nil {
		panic(err)
	}

	acc2, err := AccountFromWIF(wif2)
	if err != nil {
		panic(err)
	}

	fmt.Println("1:", acc1.Address)
	fmt.Println("2:", acc2.Address)

	// Get the balance of the account
	bal1, updated1, err := GetBalance(c, acc1)
	if err != nil {
		panic(err)
	}

	fmt.Println("1 Amount:", bal1)
	fmt.Println("1 Updated:", updated1)
	fmt.Println()

	bal2, updated2, err := GetBalance(c, acc2)
	if err != nil {
		panic(err)
	}

	fmt.Println("2 Amount:", bal2)
	fmt.Println("2 Updated:", updated2)
	fmt.Println()

	// send 1 GAS to myself
	// create a transaction
	_, g, err := NewToken(c, acc1)
	if err != nil {
		panic(err)
	}

	s, err := g.TotalSupply()
	if err != nil {
		panic(err)
	}

	fmt.Println("Total Supply:", s)

	tx, err := Tx(g, acc1, acc2, 50)
	if err != nil {
		panic(err)
	}

	fee := tx.NetworkFee + tx.SystemFee
	fmt.Println("Fee:", fee)

	hash, err := c.SendRawTransaction(tx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hash:", hash)
	fmt.Println("ValidUntil:", tx.ValidUntilBlock)

	// get mempool
	mempool, err := c.GetRawMemPool()
	if err != nil {
		panic(err)
	}

	fmt.Println("Mempool:", len(mempool))
	for _, tx := range mempool {
		fmt.Println(tx)
	}

	// get transactions
	getTxs(c, acc1)
	getTxs(c, acc2)
}
