package pgc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"

	log "github.com/inconshreveable/log15"
)

// SigmaProof includes items to prove two encryption of same message under two different public key.
type SigmaProof struct {
	//
	z1, z2, z3     *big.Int
	Pk1, Pk2       *ECPoint
	A1, A2, B1, B2 *ECPoint
	// encrypt ct point.
	ct1, ct2 *CTEncPoint
}

// SigmaProver .
type SigmaProver struct {
	// TestFlag means generation is for a test purpose(if set true).
	TestFlag bool
}

// GenerateProof generates proof to prove that two encryptions(ct1, ct2) of the same
// message under pk1 and pk2.
func (prover *SigmaProver) GenerateProof(pk1, pk2 *ecdsa.PublicKey, ct1, ct2 *TwistedELGamalCT) (proof *SigmaProof, err error) {
	curve := pk1.Curve
	n := curve.Params().N

	// check msg encoded is same.
	// testFlag == true means generate a fake proof for test purpose.
	if !prover.TestFlag && !bytes.Equal(ct1.EncMsg, ct2.EncMsg) {
		err = errors.New("msg encoded in two ct isn't same")
		return
	}

	// pick random number.
	a1, err := rand.Int(rand.Reader, n)
	if err != nil {
		return
	}
	a2, err := rand.Int(rand.Reader, n)
	if err != nil {
		return
	}
	b, err := rand.Int(rand.Reader, n)
	if err != nil {
		return
	}

	// compute proof.
	h := Params().ElgGenerator
	//
	A1 := new(ECPoint).SetFromPublicKey(pk1)
	A1.ScalarMult(A1, a1)
	A2 := new(ECPoint).SetFromPublicKey(pk2)
	A2.ScalarMult(A2, a2)

	// compute g * a1 + h * b
	hb := new(ECPoint).ScalarMult(h, b)
	ga1 := NewEmptyECPoint(curve).ScalarBaseMult(a1)
	B1 := new(ECPoint).Add(hb, ga1)
	// compute g * a2 + h * b
	ga2 := NewEmptyECPoint(curve).ScalarBaseMult(a2)
	B2 := new(ECPoint).Add(ga2, hb)

	// compute challenge e.
	e, err := SigmaFiatShamir(A1, A2, B1, B2)
	if err != nil {
		return
	}

	// compute proof.
	// z1.
	z1 := new(big.Int).Mul(e, ct1.R)
	z1 = z1.Mod(z1, n)
	z1 = z1.Add(z1, a1)
	z1 = z1.Mod(z1, n)
	// z2.
	z2 := new(big.Int).Mul(e, ct2.R)
	z2 = z2.Mod(z2, n)
	z2 = z2.Add(z2, a2)
	z2 = z2.Mod(z2, n)
	// z3.
	z3 := new(big.Int).Mul(e, new(big.Int).SetBytes(ct1.EncMsg))
	z3 = z3.Mod(z3, n)
	z3 = z3.Add(z3, b)
	z3 = z3.Mod(z3, n)

	// set proof.
	proof = new(SigmaProof)
	proof.ct1 = ct1.CopyPublicPoint()
	proof.ct2 = ct2.CopyPublicPoint()
	// warning: change to proof.pk1 will alse change pk1.
	proof.Pk1 = new(ECPoint).SetFromPublicKey(pk1)
	proof.Pk2 = new(ECPoint).SetFromPublicKey(pk2)
	proof.z1 = z1
	proof.z2 = z2
	proof.z3 = z3
	proof.A1 = A1
	proof.A2 = A2
	proof.B1 = B1
	proof.B2 = B2

	return
}

// SigmaFiatShamir calculate challenge e.
// Don't change A1, A2, B1, B2
func SigmaFiatShamir(A1, A2, B1, B2 *ECPoint) (*big.Int, error) {
	return ComputeChallenge(A1.Curve.Params().N, A1.X, A1.Y, A2.X, A2.Y, B1.X, B1.Y, B2.X, B2.Y)
}

// VerifySigmaProof validates proof.
func VerifySigmaProof(proof *SigmaProof) bool {
	e, err := SigmaFiatShamir(proof.A1, proof.A2, proof.B1, proof.B2)
	if err != nil {
		return false
	}
	// check pk1 * z1 == A1 + X1 * e.
	if !checkSigmaStep1(proof.Pk1, proof.A1, proof.ct1.X, proof.z1, e) {
		return false
	}

	// check pk2 * z2 == A2 + X2 * e.
	if !checkSigmaStep1(proof.Pk2, proof.A2, proof.ct2.X, proof.z2, e) {
		return false
	}

	h := Params().ElgGenerator
	// Check G * z1 + h * z3 == B1 + Y1 * e.
	if !checkSigmaStep2(h, proof.B1, proof.ct1.Y, proof.z1, proof.z3, e) {
		return false
	}

	// check g * z2 + h * z3 == B2 + Y2 * e.
	if !checkSigmaStep2(h, proof.B2, proof.ct2.Y, proof.z2, proof.z3, e) {
		return false
	}
	return true
}

// checkSigmaStep1 checks pk * z == A + X * e or not.
func checkSigmaStep1(pk, A, X *ECPoint, z, e *big.Int) bool {

	// compute x * e + A
	actual := new(ECPoint).ScalarMult(X, e)
	actual.Add(actual, A)
	// compute pk * z.
	expect := new(ECPoint).ScalarMult(pk, z)

	if !expect.Equals(actual) {
		log.Warn("pk * z != A + X * e", "expect x", expect.X, "expect y", expect.Y, "actual x", actual.X, "actual y", actual.Y)
		return false
	}

	return true
}

// checkSigmaStep2 checks g * z + h * z' == B + Y * e.
func checkSigmaStep2(h, B, Y *ECPoint, za, zb, e *big.Int) bool {
	// compute g * za + h * zb.
	gza := NewEmptyECPoint(h.Curve).ScalarBaseMult(za)
	hzb := new(ECPoint).ScalarMult(h, zb)
	actual := new(ECPoint).Add(gza, hzb)
	// compute B + Y * e.
	ye := new(ECPoint).ScalarMult(Y, e)
	expect := new(ECPoint).Add(ye, B)
	if !expect.Equals(actual) {
		log.Warn("g * z + h * z' != B + Y * e", "expect x", expect.X, "expect y", expect.Y, "actual x", actual.X, "actual y", actual.Y)
		return false
	}

	return true
}

// DLESigmaProver holds method to prove two ciphertexts encrypt the same value under same public key.
type DLESigmaProver struct {
	//
}

// DLESigmaProof includes items of zero-knowledge proof.
type DLESigmaProof struct {
	A1, A2 *ECPoint

	Z *big.Int

	g1, h1, g2, h2 *ECPoint
}

// GenerateProof generates zero knowledge proof to prove two ciphertexts encrypt the same value under same public key.
func (dleProver *DLESigmaProver) GenerateProof(ori, refresh *TwistedELGamalCT, sk *ecdsa.PrivateKey) (*DLESigmaProof, error) {
	// g1 = Y(fresh) - Y(ori)
	g1 := new(ECPoint).Sub(refresh.Y, ori.Y)
	// h1 = X(fresh) - X(ori)
	h1 := new(ECPoint).Sub(refresh.X, ori.X)
	// g2 = G base point.
	g2 := NewECPoint(sk.Curve.Params().Gx, sk.Curve.Params().Gy, sk.Curve)
	// h2 = pk.
	h2 := new(ECPoint).SetFromPublicKey(&sk.PublicKey)
	// witness = sk.
	w := new(big.Int).Set(sk.D)
	return dleProver.generateProof(g1, h1, g2, h2, w)
}

// generateProof generates items to prove g1 ^ w == h1; g2 ^ w == h2; w==w.
func (dleProver *DLESigmaProver) generateProof(g1, h1, g2, h2 *ECPoint, w *big.Int) (*DLESigmaProof, error) {
	//
	curve := g1.Curve
	n := curve.Params().N
	a, err := rand.Int(rand.Reader, n)
	if err != nil {
		return nil, err
	}

	// A1 = g1 * a; A2 = g2 * a;
	A1 := new(ECPoint).ScalarMult(g1, a)
	A2 := new(ECPoint).ScalarMult(g2, a)
	// compute challenge e.
	e, err := ComputeChallenge(n, A1.X, A1.Y, A2.X, A2.Y)
	if err != nil {
		return nil, err
	}

	// compute z = a + e * w.
	z := new(big.Int).Mul(e, w)
	z = z.Mod(z, n)
	z = z.Add(z, a)
	z = z.Mod(z, n)

	// set proof
	proof := new(DLESigmaProof)
	proof.A1 = A1
	proof.A2 = A2
	proof.Z = z
	proof.g1 = g1
	proof.h1 = h1
	proof.g2 = g2
	proof.h2 = h2

	return proof, nil
}

// VerifyDLESigmaProof verify the proof is valid or not.
func VerifyDLESigmaProof(proof *DLESigmaProof) bool {
	curve := proof.A1.Curve

	e, err := ComputeChallenge(curve.Params().N, proof.A1.X, proof.A1.Y, proof.A2.X, proof.A2.Y)
	if err != nil {
		return false
	}

	// check g1 * z == A1 + h1 * e.
	if !checkDLESigmaProof(proof.g1, proof.A1, proof.h1, proof.Z, e) {
		return false
	}
	// check g2 * z == A2 + h2 * e.
	if !checkDLESigmaProof(proof.g2, proof.A2, proof.h2, proof.Z, e) {
		return false
	}

	return true
}

// checkDLESigmaProof checks g * z == A + h * e.
func checkDLESigmaProof(g, A, H *ECPoint, z, e *big.Int) bool {
	// g * z.
	gz := new(ECPoint).ScalarMult(g, z)
	// h * e + A.
	he := new(ECPoint).ScalarMult(H, e)
	expect := new(ECPoint).Add(A, he)

	if !expect.Equals(gz) {
		log.Warn("g * z != A + h * e", "expect x", expect.X, "expect y", expect.Y, "actual x", gz.X, "actual y", gz.Y)
		return false
	}

	return true
}