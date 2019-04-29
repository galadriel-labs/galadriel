package pgc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

// CTEncPoint respresents encrypted ct tx point on curve.
type CTEncPoint struct {
	X, Y *ecdsa.PublicKey
}

// TwistedELGamalCT respresents a encrypted transaction encoded in twisted-elgamal format.
type TwistedELGamalCT struct {
	// X, Y
	CTEncPoint
	// Random number used in encryption.
	R *big.Int
	// encrypted msg.
	EncMsg []byte
}

// CopyPublicPoint copys public encrypted ec point.
func (ct *TwistedELGamalCT) CopyPublicPoint() *CTEncPoint {
	ecPoints := CTEncPoint{}

	x := new(ecdsa.PublicKey)
	x.X = ct.X.X
	x.Y = ct.X.Y
	x.Curve = ct.X.Curve

	y := new(ecdsa.PublicKey)
	y.X = ct.Y.X
	y.Y = ct.Y.Y
	y.Curve = ct.X.Curve

	ecPoints.X = x
	ecPoints.Y = y
	return &ecPoints
}

// GenerateKey generates key pair using btcec s256 curve.
func GenerateKey() (key *ecdsa.PrivateKey, err error) {
	key, err = ecdsa.GenerateKey(S256(), rand.Reader)
	return
}

// Encrypt encrypts msg by in twisted elgamal way.
func Encrypt(pk *ecdsa.PublicKey, msg []byte) (*TwistedELGamalCT, error) {
	// Create ct instance.
	ct := new(TwistedELGamalCT)
	ct.X = new(ecdsa.PublicKey)
	ct.Y = new(ecdsa.PublicKey)

	// set curve
	curve := pk.Curve
	ct.X.Curve = curve
	ct.Y.Curve = curve

	// get random.
	r, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		return nil, err
	}
	// for sigma proof purpose.
	ct.R = new(big.Int).Set(r)
	ct.EncMsg = make([]byte, len(msg))
	copy(ct.EncMsg, msg)

	// compute pk * r.(pk ^ r)
	ct.X.X, ct.X.Y = curve.ScalarMult(pk.X, pk.Y, r.Bytes())
	// compute g * r.(g ^ r)
	s1X, s1Y := curve.ScalarBaseMult(r.Bytes())
	// compute h * m.(h ^ m)
	pubParams := Params()
	s2X, s2Y := curve.ScalarMult(pubParams.ElgGenerator.X, pubParams.ElgGenerator.Y, msg)
	// compute g * r + h * m.
	ct.Y.X, ct.Y.Y = curve.Add(s1X, s1Y, s2X, s2Y)

	return ct, nil
}

// Decrypt decrypts msg in twisted elgamal way.
func Decrypt(sk *ecdsa.PrivateKey, ct *TwistedELGamalCT) []byte {
	// get public curve.
	curve := sk.PublicKey.Curve
	// compute the inverse of sk.
	skInverse := new(big.Int).ModInverse(sk.D, curve.Params().N)
	// compute X * skInverse(X ^ skInverse) == g * r.
	ct.X.X, ct.X.Y = curve.ScalarMult(ct.X.X, ct.X.Y, skInverse.Bytes())
	// use ct.Y - ct.X Point to get encoded msg.
	// get Affine negation formulas of ct.Y.
	encodedMsg := SubECPoint(ct.Y, ct.X)
	_ = encodedMsg

	// todo: decrypt encoded msg.
	return []byte{}
}

// DecryptEncodedMsg decrypts and returns original bytes of msg.
func DecryptEncodedMsg() []byte {
	//
	return []byte{}
}

// SubECPoint computes x - y.
func SubECPoint(x, y *ecdsa.PublicKey) *ecdsa.PublicKey {
	curve := x.Curve

	// create instance of new point.
	newPoint := new(ecdsa.PublicKey)
	newPoint.Curve = curve

	// get negation of x.
	negation := negation(y)
	// calculate negation of new point.
	newPointNegaX, newPointNegaY := curve.Add(negation.X, negation.Y, x.X, x.Y)
	newPoint.X = newPointNegaX
	// set new point y to symmetry
	newPoint.Y = newPointNegaY

	return newPoint
}

// negation returns negation of ecpoint x on curve.
func negation(x *ecdsa.PublicKey) *ecdsa.PublicKey {
	negation := new(ecdsa.PublicKey)
	negation.Curve = x.Curve

	negation.X = new(big.Int).Set(x.X)
	negation.Y = new(big.Int).Set(x.Y)
	// y' = -y mod p.
	negation.Y = negation.Y.Neg(negation.Y)
	negation.Y = negation.Y.Mod(negation.Y, x.Curve.Params().P)
	//negation.Y = negation.Y.Mod(negation.Y, negation.Curve.Params().N)
	if !negation.Curve.IsOnCurve(negation.X, negation.Y) {
		panic("negation of x, y is not on curve")
	}

	return negation
}
