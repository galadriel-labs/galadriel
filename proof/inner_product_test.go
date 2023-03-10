package proof

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/pgc/curve"
	"github.com/pgc/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInnerProduct(t *testing.T) {
	cases := []struct {
		curve   elliptic.Curve
		bitsize int
	}{
		{
			curve:   curve.BN256(),
			bitsize: 64,
		},
		{
			curve:   curve.BN256(),
			bitsize: 32,
		},
		{
			curve:   curve.BN256(),
			bitsize: 16,
		},
		{
			curve:   curve.BN256(),
			bitsize: 8,
		},
		{
			curve:   curve.BN256(),
			bitsize: 4,
		},
		{
			curve:   curve.S256(),
			bitsize: 64,
		},
		{
			curve:   curve.S256(),
			bitsize: 32,
		},
		{
			curve:   curve.S256(),
			bitsize: 16,
		},
		{
			curve:   curve.S256(),
			bitsize: 8,
		},
		{
			curve:   curve.S256(),
			bitsize: 4,
		},
	}
	for _, c := range cases {
		testInnerProduct(t, c.curve, c.bitsize)
	}

}

func testInnerProduct(t *testing.T, curve elliptic.Curve, bitsize int) {
	params := NewInnerProductRandomParams(curve, bitsize)
	p, c, a, b := NewInnerProductRandomCommitments(params, curve.Params().N, bitsize)

	proof, err := GenIPProof(params, p, c, a, b)
	require.Nil(t, err, "generate inner product failed")

	assert.True(t, VerifyIPProof(params, p, c, proof), "normal verify failed")
	assert.True(t, OptimizedVerifyIPProof(params, p, c, proof), "optimized verify failed")

	// test for invalid inner product proof.
	proof.a.Sub(proof.a, new(big.Int).SetUint64(1))
	assert.Equal(t, false, VerifyIPProof(params, p, c, proof), "invalid proof pass normal verify")
	assert.Equal(t, false, OptimizedVerifyIPProof(params, p, c, proof), "invalid proof pass optimized verify")
}

func BenchmarkBN256InnerProduct(b *testing.B) {
	benchmarkTest(b, curve.BN256(), 32)
}
func BenchmarkBN256InnerProductVerifyNoraml(b *testing.B) {
	benchmarkVerify(b, curve.BN256(), 32, VerifyIPProof)
}

func BenchmarkBN256InnerProductVerifyOptimized(b *testing.B) {
	benchmarkVerify(b, curve.BN256(), 32, OptimizedVerifyIPProof)
}

func BenchmarkS256InnerProduct(b *testing.B) {
	benchmarkTest(b, curve.S256(), 32)
}

func BenchmarkS256InnerProductVerifyNoraml(b *testing.B) {
	benchmarkVerify(b, curve.S256(), 32, VerifyIPProof)
}

func BenchmarkS256InnerProductVerifyOptimized(b *testing.B) {
	benchmarkVerify(b, curve.S256(), 32, OptimizedVerifyIPProof)
}

type verifyFunc func(params IPParams, p *utils.ECPoint, c *big.Int, proof *IPProof) bool

func benchmarkVerify(b *testing.B, curve elliptic.Curve, bitsize int, vf verifyFunc) {
	b.StopTimer()
	params := NewInnerProductRandomParams(curve, bitsize)
	p, c, a, pb := NewInnerProductRandomCommitments(params, curve.Params().N, bitsize)
	proof, err := GenIPProof(params, p, c, a, pb)
	if err != nil {
		b.Fatal("generate proof failed")
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !vf(params, p, c, proof) {
			b.Fatal("verify failed")
		}
	}
}

func benchmarkTest(b *testing.B, curve elliptic.Curve, bitsize int) {
	b.StopTimer()
	params := NewInnerProductRandomParams(curve, bitsize)
	p, c, a, pb := NewInnerProductRandomCommitments(params, curve.Params().N, bitsize)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenIPProof(params, p, c, a, pb)
		if err != nil {
			b.Fatal("generate proof failed")
		}
	}
}
