package pgcsys

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/pgc/proof"
	"github.com/pgc/utils"
)

// Account represents the account used in pgc system.
type Account struct {
	name    string
	sk      *ecdsa.PrivateKey
	balance *proof.CTEncPoint
	// balance of account. for speed up purpose only.
	m      *big.Int
	nonce  uint64
	params proof.CTParams
}

func NewAccount(name string,
	sk *ecdsa.PrivateKey,
	balance *proof.CTEncPoint,
	m *big.Int,
	params proof.CTParams,
	nonce uint64) *Account {
	return &Account{
		name:    name,
		sk:      sk,
		balance: balance,
		m:       m,
		params:  params,
		nonce:   nonce,
	}
}

// CreateTestAccount creates account for test pupose.
func CreateTestAccount(params proof.CTParams, name string, balance *big.Int) *Account {
	a := Account{}
	a.name = name
	a.m = new(big.Int).Set(balance)
	a.nonce = 0
	a.params = params
	key := proof.MustGenerateKey(params)

	a.sk = key
	ct, err := proof.EncryptOnChain(params, &a.sk.PublicKey, balance.Bytes())
	if err != nil {
		panic(err)
	}
	a.balance = ct.CopyPublicPoint()

	return &a
}

func (a *Account) Balance() *proof.CTEncPoint {
	return a.balance
}

// UpdateBalance updates user's encrypted balance and nonce.
func (a *Account) UpdateBalance(nonce *big.Int, ct [4]*big.Int) {
	a.balance.X = utils.NewECPoint(ct[0], ct[1], a.sk.Curve)
	a.balance.Y = utils.NewECPoint(ct[2], ct[3], a.sk.Curve)
	a.nonce = nonce.Uint64()

	// decrypt balance
	balance, _ := proof.Decrypt(a.params, a.sk, a.balance)
	a.m = new(big.Int).SetBytes(balance)
}

func (a *Account) UpdateBalanceInter(nonce *big.Int, nbalace *proof.CTEncPoint) {
	a.balance = nbalace
	a.nonce = nonce.Uint64()

	// decrypt balance
	balance, _ := proof.Decrypt(a.params, a.sk, a.balance)
	a.m = new(big.Int).SetBytes(balance)
}

func (a *Account) M() *big.Int {
	return a.m
}

func (a *Account) Priv() *ecdsa.PrivateKey {
	return a.sk
}
