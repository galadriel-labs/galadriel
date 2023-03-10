package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/pgc/curve"
)

// ECPoint represents a point on elliptic curve.
type ECPoint struct {
	X     *big.Int
	Y     *big.Int
	Curve elliptic.Curve
}

func (ec *ECPoint) checkNil() {
	if ec.X == nil && ec.Y == nil {
		ec.X = new(big.Int).SetUint64(0)
		ec.Y = new(big.Int).SetUint64(0)
	}
}

// MarshalJSON defines custom way to marshal json.
func (ec *ECPoint) MarshalJSON() ([]byte, error) {
	newJSON := struct {
		X     string
		Y     string
		Curve string
	}{
		X:     ec.X.String(),
		Y:     ec.Y.String(),
		Curve: ec.Curve.Params().Name,
	}

	return json.Marshal(&newJSON)
}

// UnmarshalJSON defines custom way to unmarshal json.
func (ec *ECPoint) UnmarshalJSON(bz []byte) error {
	var data struct {
		X     string
		Y     string
		Curve string
	}

	if err := json.Unmarshal(bz, &data); err != nil {
		return err
	}

	cv, err := curve.GetCurve(data.Curve)
	if err != nil {
		return err
	}

	x, ok := new(big.Int).SetString(data.X, 10)
	if !ok {
		return errors.New("Invalid ec point: x incorrect")
	}
	y, ok := new(big.Int).SetString(data.Y, 10)
	if !ok {
		return errors.New("Invalid ec point: y incorrect")
	}

	*ec = ECPoint{
		X:     x,
		Y:     y,
		Curve: cv,
	}

	return nil
}

// NewECPoint returns instance of ec point.
func NewECPoint(x, y *big.Int, curve elliptic.Curve) *ECPoint {
	ec := ECPoint{}
	ec.Curve = curve
	ec.X = new(big.Int).Set(x)
	ec.Y = new(big.Int).Set(y)

	return &ec
}

// NewRandomECPoint creates a ec point by randomness.
func NewRandomECPoint(curve elliptic.Curve) *ECPoint {
	h, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		panic(err)
	}

	hx, hy := curve.ScalarBaseMult(h.Bytes())

	return NewECPoint(hx, hy, curve)
}

// NewECPointByString takes a string as input and returns a ecpoint on curve.
func NewECPointByString(s string, curve elliptic.Curve) *ECPoint {
	data := Keccak256([]byte(s))
	scalar := new(big.Int).SetBytes(data)
	scalar.Mod(scalar, curve.Params().N)

	return NewEmptyECPoint(curve).ScalarBaseMult(scalar)
}

// NewECPointByBytes takes bytes as input and returns a ecpoint on curve.
func NewECPointByBytes(data []byte, curve elliptic.Curve) *ECPoint {
	return NewEmptyECPoint(curve).ScalarBaseMult(new(big.Int).SetBytes(data))
}

// NewEmptyECPoint creates instance of ec point without x or y point.
// But set the curve.(X, Y will be set to 0)
func NewEmptyECPoint(curve elliptic.Curve) *ECPoint {
	ec := ECPoint{}
	ec.Curve = curve
	ec.X = new(big.Int).SetUint64(0)
	ec.Y = new(big.Int).SetUint64(0)

	return &ec
}

// Add adds two ec point and set ec to new points.
func (ec *ECPoint) Add(first, second *ECPoint) *ECPoint {
	ec.checkNil()
	if first.X.Uint64() == 0 && first.Y.Uint64() == 0 {
		ec.Apply(second)
		return ec
	}
	if second.X.Uint64() == 0 && second.Y.Uint64() == 0 {
		ec.Apply(first)
		return ec
	}
	ec.X, ec.Y = first.Curve.Add(first.X, first.Y, second.X, second.Y)
	ec.Curve = first.Curve

	return ec
}

// Sub returns/set res of first - second.
// first/second unchanged.
func (ec *ECPoint) Sub(first, second *ECPoint) *ECPoint {
	ec.checkNil()
	if first.Equal(second) {
		ec.X = big.NewInt(0)
		ec.Y = big.NewInt(0)
		ec.Curve = first.Curve
		return ec
	}
	negation := new(ECPoint).Negation(second)
	ec.X, ec.Y = first.Curve.Add(negation.X, negation.Y, first.X, first.Y)
	ec.Curve = first.Curve

	return ec
}

// Negation returns negation of other point.
// set ec to -other.
func (ec *ECPoint) Negation(other *ECPoint) *ECPoint {
	ec.checkNil()
	ec.Curve = other.Curve
	ec.X = new(big.Int).Set(other.X)
	ec.Y = new(big.Int).Neg(other.Y)
	ec.Y = ec.Y.Mod(ec.Y, other.Curve.Params().P)

	return ec
}

// Compress compress the point.
func (ec *ECPoint) Compress() []byte {
	ec.checkNil()

	compressed := make([]byte, 33)
	format := byte(0x2)
	// odd
	if ec.Y.Bit(0) == 1 {
		format |= 0x1
	}

	compressed[0] = format
	xbytes := ec.X.Bytes()
	copy(compressed[33-len(xbytes):], ec.X.Bytes())

	return compressed
}

// MurmurHashToUint64 calcute murmur hash
func (ec *ECPoint) MurmurHashToUint64() uint64 {
	return MurmurHash64AWithFixedSalt(ec.Compress())
}

func DecompressPointBytes(cv elliptic.Curve, data []byte) (*ECPoint, error) {
	if len(data) != 33 {
		return nil, errors.New("invalid point len, expect 33")
	}

	ybit := data[0]&0x1 == 1
	x := new(big.Int).SetBytes(data[1:])
	return DecompressPoint(cv, x, ybit)
}

// DecompressPoint decompress point
func DecompressPoint(cv elliptic.Curve, x *big.Int, ybit bool) (*ECPoint, error) {
	// y^2 = x^3 + ax + b
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)

	a, err := curve.GetCurveA(cv.Params().Name)
	if err != nil {
		return nil, err
	}
	ax := new(big.Int).Mul(a, x)

	x3.Add(x3, ax)
	x3.Add(x3, cv.Params().B)
	x3.Mod(x3, cv.Params().P)

	// get y
	y := new(big.Int).ModSqrt(x3, cv.Params().P)
	if isOdd(y) != ybit {
		y.Sub(cv.Params().P, y)
	}
	// check
	{
		y2 := new(big.Int).Mul(y, y)
		y2.Mod(y2, cv.Params().P)
		if x3.Cmp(y2) != 0 {
			return new(ECPoint), errors.New("invalid point, not on curve")
		}
	}

	pk := new(ECPoint)
	pk.Curve = cv
	pk.X = x
	pk.Y = y

	return pk, nil
}

// ScalarMult returns ec * scalar.
// set ec to new point.
func (ec *ECPoint) ScalarMult(base *ECPoint, k *big.Int) *ECPoint {
	ec.X, ec.Y = base.Curve.ScalarMult(base.X, base.Y, k.Bytes())
	ec.Curve = base.Curve
	ec.checkNil()

	return ec
}

// ScalarBaseMult returns g * scalar.(curve must be set)
// set ec to new point.
func (ec *ECPoint) ScalarBaseMult(k *big.Int) *ECPoint {
	ec.checkNil()
	ec.X, ec.Y = ec.Curve.ScalarBaseMult(k.Bytes())

	return ec
}

// SetFromPublicKey set current point to another ec point by publickey.
func (ec *ECPoint) SetFromPublicKey(other *ecdsa.PublicKey) *ECPoint {
	ec.X = new(big.Int).Set(other.X)
	ec.Y = new(big.Int).Set(other.Y)
	ec.Curve = other.Curve

	return ec
}

// ToPublicKey returns a public key of ecpoint.
func (ec *ECPoint) ToPublicKey() *ecdsa.PublicKey {
	res := new(ecdsa.PublicKey)
	res.X = ec.X
	res.Y = ec.Y
	res.Curve = ec.Curve

	return res
}

// Copy returns a new instance of current ec point.
func (ec *ECPoint) Copy() *ECPoint {
	n := new(ECPoint)
	n.X = new(big.Int).Set(ec.X)
	n.Y = new(big.Int).Set(ec.Y)
	n.Curve = ec.Curve

	return n
}

// Equal check two point equal or not.
// return res.
func (ec *ECPoint) Equal(other *ECPoint) bool {
	return ec.X.Cmp(other.X) == 0 && ec.Y.Cmp(other.Y) == 0 && ec.Curve.Params().Name == other.Curve.Params().Name
}

//
func (ec *ECPoint) Apply(other *ECPoint) {
	ec.X = other.X
	ec.Y = other.Y
	ec.Curve = other.Curve
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}
