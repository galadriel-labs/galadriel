package pgc

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"io"
	"math/big"
)

const (
	// This is copied from https://github.com/golang/crypto/blob/master/openpgp/elgamal/elgamal_test.go
	// The rationality should be discussed before usage.
	p = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"
	g = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"
)

// The original implementation is in https://github.com/golang/crypto/blob/master/openpgp/elgamal/elgamal.go.
// The implementation here will change the encryption and decryption to implement twisting in pgc.
// https://eprint.iacr.org/2019/319.pdf

// PublicKey represents an ElGamal public key.
// Y represents public value y := g ** x mod p, and x is the private key.
type PublicKey struct {
	G, P, Y *big.Int
}

// PrivateKey represents an ElGamal private key.
// X is the private key in (1, q-1)
type PrivateKey struct {
	PublicKey
	X *big.Int
}

// Encrypt encrypts the given message to the given public key. The result is a
// pair of integers. Errors can result from reading random, or because msg is
// too large to be encrypted to the public key.
func Encrypt(random io.Reader, pub *PublicKey, msg []byte) (c1, c2 *big.Int, err error) {
	pLen := (pub.P.BitLen() + 7) / 8
	if len(msg) > pLen-11 {
		err = errors.New("elgamal: message too long")
		return
	}

	// EM = 0x02 || PS || 0x00 || M
	em := make([]byte, pLen-1)
	em[0] = 2
	ps, mm := em[1:len(em)-len(msg)-1], em[len(em)-len(msg):]
	err = nonZeroRandomBytes(ps, random)
	if err != nil {
		return
	}
	em[len(em)-len(msg)-1] = 0
	copy(mm, msg)

	m := new(big.Int).SetBytes(em)

	k, err := rand.Int(random, pub.P)
	if err != nil {
		return
	}

	c1 = new(big.Int).Exp(pub.G, k, pub.P)
	s := new(big.Int).Exp(pub.Y, k, pub.P)
	c2 = s.Mul(s, m)
	c2.Mod(c2, pub.P)

	return
}

// Decrypt takes two integers, resulting from an ElGamal encryption, and
// returns the plaintext of the message. An error can result only if the
// ciphertext is invalid. Users should keep in mind that this is a padding
// oracle and thus, if exposed to an adaptive chosen ciphertext attack, can
// be used to break the cryptosystem.  See ``Chosen Ciphertext Attacks
// Against Protocols Based on the RSA Encryption Standard PKCS #1'', Daniel
// Bleichenbacher, Advances in Cryptology (Crypto '98),
func Decrypt(priv *PrivateKey, c1, c2 *big.Int) (msg []byte, err error) {
	s := new(big.Int).Exp(c1, priv.X, priv.P)
	s.ModInverse(s, priv.P)
	s.Mul(s, c2)
	s.Mod(s, priv.P)
	em := s.Bytes()

	firstByteIsTwo := subtle.ConstantTimeByteEq(em[0], 2)

	// The remainder of the plaintext must be a string of non-zero random
	// octets, followed by a 0, followed by the message.
	//   lookingForIndex: 1 iff we are still looking for the zero.
	//   index: the offset of the first zero byte.
	var lookingForIndex, index int
	lookingForIndex = 1

	for i := 1; i < len(em); i++ {
		equals0 := subtle.ConstantTimeByteEq(em[i], 0)
		index = subtle.ConstantTimeSelect(lookingForIndex&equals0, i, index)
		lookingForIndex = subtle.ConstantTimeSelect(equals0, 0, lookingForIndex)
	}

	if firstByteIsTwo != 1 || lookingForIndex != 0 || index < 9 {
		return nil, errors.New("elgamal: decryption error")
	}
	return em[index+1:], nil
}

// nonZeroRandomBytes fills the given slice with non-zero random octets.
func nonZeroRandomBytes(s []byte, rand io.Reader) (err error) {
	_, err = io.ReadFull(rand, s)
	if err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		for s[i] == 0 {
			_, err = io.ReadFull(rand, s[i:i+1])
			if err != nil {
				return
			}
		}
	}

	return
}