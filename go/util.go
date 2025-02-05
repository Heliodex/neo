package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/actor"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/gas"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/nep17"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

func getTxs(c *rpcclient.Client, acc *wallet.Account) {
	start, stop := uint64(0), uint64(time.Now().UnixMilli())
	xfers, err := c.GetNEP17Transfers(acc.Contract.ScriptHash(), &start, &stop, nil, nil)
	if err != nil {
		panic(err)
	}

	received := xfers.Received
	sent := xfers.Sent

	fmt.Println("Received:", len(received))
	for _, t := range received {
		logTransfer(t)
	}

	fmt.Println("Sent:", len(sent))
	for _, t := range sent {
		logTransfer(t)
	}
}

// actor and token (token is gas for now (forever))
func NewToken(c *rpcclient.Client, acc *wallet.Account) (act *actor.Actor, token *nep17.Token, err error) {
	signer := transaction.Signer{
		Account: acc.Contract.ScriptHash(),
		Scopes:  transaction.CalledByEntry,
	}

	a, err := actor.New(c, []actor.SignerAccount{{
		Signer:  signer,
		Account: acc,
	}})
	if err != nil {
		return
	}

	return a, gas.New(a), nil
}

type Balance struct {
	Amount      string
	LastUpdated uint32
}

func GetBalance(c *rpcclient.Client, acc *wallet.Account) (amount string, updated uint32, err error) {
	bal, err := c.GetNEP17Balances(acc.ScriptHash())
	if err != nil {
		return
	}

	for _, b := range bal.Balances {
		if b.Asset == gas.Hash {
			return b.Amount, b.LastUpdated, nil
		}
	}

	return "", 0, fmt.Errorf("no token balance found")
}

func Tx(g *nep17.Token, from, to *wallet.Account, amount uint64) (*transaction.Transaction, error) {
	return g.TransferTransaction(from.Contract.ScriptHash(), to.Contract.ScriptHash(), big.NewInt(int64(amount)), nil)
}

func AccountFromWIF(wif string) (*wallet.Account, error) {
	sk, err := keys.NewPrivateKeyFromWIF(wif)
	if err != nil {
		panic(err)
	}

	return wallet.NewAccountFromPrivateKey(sk), nil
}
