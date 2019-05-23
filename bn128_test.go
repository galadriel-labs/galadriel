package pgc

import (
	"crypto/rand"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestG1Point(t *testing.T) {
	// get a random g1 point.
	_, p, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}

	x, y, err := G1ToPoint(p)
	if err != nil {
		t.Error(err)
		return
	}

	newP, err := PointToG1(x, y)
	if err != nil {
		t.Error(err)
		return
	}

	if p.String() != newP.String() {
		t.Error("convert/unconvert from/to g1 point failed")
	}

	c := BN128{}
	params := c.Params()
	gx, gy := params.Gx, params.Gy
	r, err := rand.Int(rand.Reader, params.N)
	if err != nil {
		t.Error(err)
		return
	}

	a, b := c.ScalarBaseMult(r.Bytes())
	aa, bb := c.ScalarMult(gx, gy, r.Bytes())
	if a.Cmp(aa) != 0 || b.Cmp(bb) != 0 {
		t.Error("curve error")
	}

	// get a neg of a point.
	pp := MustPointToG1(a, b)
	ppNeg := new(bn256.G1).Neg(pp)
	a, bNeg := MustG1ToPoint(ppNeg)

	pOrder := new(big.Int).Add(bNeg, b)
	if pOrder.Cmp(params.P) != 0 {
		t.Error("curve p error")
	}
}