package proof

import (
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"math/big"

	log "github.com/inconshreveable/log15"
	"github.com/pgc/utils"
)

var (
	// ErrInvalidVectorSize means params of vector is not the same.
	ErrInvalidVectorSize = errors.New("gv, hv, a, b size not equal")
)

// IPParams contains the parameters of an inner product proof.
type IPParams interface {
	GV() *utils.GeneratorVector
	HV() *utils.GeneratorVector
	U() *utils.ECPoint
}

type ipparams struct {
	gv, hv *utils.GeneratorVector
	u      *utils.ECPoint
}

// MarshalJSON defines a custom way to marshal.
func (ipp *ipparams) MarshalJSON() ([]byte, error) {
	newJSON := struct {
		GV, HV *utils.GeneratorVector
		U      *utils.ECPoint
	}{
		GV: ipp.gv,
		HV: ipp.hv,
		U:  ipp.u,
	}

	return json.Marshal(&newJSON)
}

// NewIPParams returns a new instance of ipprams.
func NewIPParams(gv, hv *utils.GeneratorVector, u *utils.ECPoint) IPParams {
	return &ipparams{
		gv: gv,
		hv: hv,
		u:  u,
	}
}

// NewInnerProductRandomParams generates random params for test purpose
func NewInnerProductRandomParams(curve elliptic.Curve, size int) IPParams {
	params := ipparams{}
	params.gv = utils.NewRandomGeneratorVector(curve, size)
	params.hv = utils.NewRandomGeneratorVector(curve, size)
	params.u = utils.NewRandomECPoint(curve)

	return &params
}

func (p *ipparams) GV() *utils.GeneratorVector {
	return p.gv
}

func (p *ipparams) HV() *utils.GeneratorVector {
	return p.hv
}

func (p *ipparams) U() *utils.ECPoint {
	return p.u
}

// IPProof represents inner product proof.
type IPProof struct {
	l, r []*utils.ECPoint

	a, b *big.Int
}

// IPProofInput represents inner product proof solidity input format.
type IPProofInput struct {
	l, r []*big.Int
	a, b *big.Int
}

// MarshalJSON defines a custom way to marshal.
func (ip *IPProof) MarshalJSON() ([]byte, error) {
	newJSON := struct {
		L, R []*utils.ECPoint
		A, B string
	}{
		L: ip.l,
		R: ip.r,
		A: ip.a.String(),
		B: ip.b.String(),
	}

	return json.Marshal(&newJSON)
}

// UnmarshalJSON defines a custom way to unmarshal
func (ip *IPProof) UnmarshalJSON(bz []byte) error {
	var data struct {
		L, R []*utils.ECPoint
		A, B string
	}

	if err := json.Unmarshal(bz, &data); err != nil {
		return err
	}

	a, ok := new(big.Int).SetString(data.A, 10)
	if !ok {
		return errors.New("Invalid a in ip proof")
	}
	b, ok := new(big.Int).SetString(data.B, 10)
	if !ok {
		return errors.New("Invalid b in ip proof")
	}

	*ip = IPProof{
		l: data.L,
		r: data.R,
		a: a,
		b: b,
	}
	return nil
}

// ToSolidityInput change format to solidity contract input.
func (ipproof *IPProof) ToSolidityInput() *IPProofInput {
	input := IPProofInput{}
	input.l = make([]*big.Int, 0)
	input.r = make([]*big.Int, 0)

	for i := 0; i < len(ipproof.l); i++ {
		input.l = append(input.l, ipproof.l[i].X)
		input.l = append(input.l, ipproof.l[i].Y)

		input.r = append(input.r, ipproof.r[i].X)
		input.r = append(input.r, ipproof.r[i].Y)
	}

	input.a = new(big.Int).Set(ipproof.a)
	input.b = new(big.Int).Set(ipproof.b)

	return &input
}

// NewInnerProductRandomCommitments generates random commitments for test
func NewInnerProductRandomCommitments(params IPParams, n *big.Int, size int) (*utils.ECPoint, *big.Int, *utils.FieldVector, *utils.FieldVector) {
	a := utils.NewRandomFieldVector(n, size)
	b := utils.NewRandomFieldVector(n, size)
	c := a.InnerProduct(b)
	pa := params.GV().SubVector(0, size).Commit(a.GetVector())
	pb := params.HV().SubVector(0, size).Commit(b.GetVector())
	p := new(utils.ECPoint).Add(pa, pb)

	return p, c, a, b
}

// NewIPProof creates instance of inner product proof.
func NewIPProof(l, r []*utils.ECPoint, a, b *big.Int) *IPProof {
	proof := IPProof{}
	proof.l = l
	proof.r = r
	proof.a = new(big.Int).Set(a)
	proof.b = new(big.Int).Set(b)

	return &proof
}

// GenIPProof generates proof using protocol 1 in bulletproof inner-product argument to prove that
// prover knows two vector a, b and p = g * a + h *b; and c = <a, b>;
// the g, h, p, h is public known by verifier.
// g and h are generator vectors.
func GenIPProof(params IPParams, p *utils.ECPoint, c *big.Int, a, b *utils.FieldVector) (*IPProof, error) {

	// test.
	if c.Cmp(a.InnerProduct(b)) != 0 {
		log.Error("not equal c, a inner b")
	}
	//
	gv := params.GV()
	hv := params.HV()
	// todo: only check in test if nessary.
	gva := gv.Commit(a.GetVector())
	hvb := hv.Commit(b.GetVector())
	gvahvb := new(utils.ECPoint).Add(gva, hvb)
	if !gvahvb.Equal(p) {
		log.Info("hp", "x", hv.Get(0).X)
	}
	u := params.U()
	if gv.Size() != hv.Size() || hv.Size() != a.Size() || a.Size() != b.Size() {
		return nil, ErrInvalidVectorSize
	}

	ue, np, err := cal(p, u, c)
	if err != nil {
		return nil, err
	}
	l, r := make([]*utils.ECPoint, 0), make([]*utils.ECPoint, 0)

	// call protocol 2.
	return genIPProofInternal(gv, hv, ue, np, a, b, l, r)
}

// cal calculates new p point on protocol 2.
func cal(p, u *utils.ECPoint, c *big.Int) (*utils.ECPoint, *utils.ECPoint, error) {
	e, err := utils.ComputeChallenge(p.Curve.Params().N, c)
	if err != nil {
		return nil, nil, err
	}
	// compute new p point.
	// p' = p + u * (e * c)
	ue := new(utils.ECPoint).ScalarMult(u, e)
	np := new(utils.ECPoint).ScalarMult(ue, c)
	np.Add(np, p)

	return ue, np, nil
}

// genIPProofInternal generates proof for inner product proof.
// g, h are two public vector generator used in bullet proof.
// u = u * e(u represents a fix point in protocol);
// p = g * a + h * b + u * c.
func genIPProofInternal(g, h *utils.GeneratorVector, u, p *utils.ECPoint, a, b *utils.FieldVector, l, r []*utils.ECPoint) (*IPProof, error) {
	if g.Size() == 1 {
		return NewIPProof(l, r, a.First(), b.First()), nil
	}

	// todo: make sure size is power of 2.
	gLeft := g.HalfLeft()
	gRight := g.HalfRight()
	hLeft := h.HalfLeft()
	hRight := h.HalfRight()

	aLeft := a.HalfLeft()
	aRight := a.HalfRight()
	bLeft := b.HalfLeft()
	bRight := b.HalfRight()

	cL := aLeft.InnerProduct(bRight)
	cR := aRight.InnerProduct(bLeft)

	// compute L = gRight * aLeft + hLeft * bRight + u * cL.
	curve := u.Curve
	l1 := gRight.Commit(aLeft.GetVector())
	l2 := hLeft.Commit(bRight.GetVector())

	lp := new(utils.ECPoint).ScalarMult(u, cL)
	lp.Add(lp, l1)
	lp.Add(lp, l2)
	l = append(l, lp)

	// compute R = gLeft * aRight + hRight * bLeft + u * cR.
	r1 := gLeft.Commit(aRight.GetVector())
	r2 := hRight.Commit(bLeft.GetVector())
	rp := new(utils.ECPoint).ScalarMult(u, cR)
	rp.Add(rp, r1)
	rp.Add(rp, r2)
	r = append(r, rp)

	// compute challenge x base on l, r.
	// todo: may be change the challenge compution.
	n := curve.Params().N
	e, err := utils.ComputeChallenge(n, lp.X, lp.Y, rp.X, rp.Y)
	if err != nil {
		return nil, err
	}
	eInverse := new(big.Int).ModInverse(e, n)

	gPrime := gLeft.HadamardScalar(eInverse).AddGeneratorVector(gRight.HadamardScalar(e))
	hPrime := hLeft.HadamardScalar(e).AddGeneratorVector(hRight.HadamardScalar(eInverse))

	aPrime := aLeft.Times(e).AddFieldVector(aRight.Times(eInverse))
	bPrime := bLeft.Times(eInverse).AddFieldVector(bRight.Times(e))

	// compute e ^ 2.
	eSquare := new(big.Int).Mul(e, e)
	eSquare.Mod(eSquare, n)
	// compute e ^ -2 = eInverse ^ 2.
	eInverseSquare := new(big.Int).Mul(eInverse, eInverse)
	eInverseSquare.Mod(eInverseSquare, n)

	// compute p' = l * (x ^ 2) + p + r * (x ^ -2).
	newP := new(utils.ECPoint).ScalarMult(lp, eSquare)
	newP.Add(newP, p)
	rTmp := new(utils.ECPoint).ScalarMult(rp, eInverseSquare)
	newP.Add(newP, rTmp)

	return genIPProofInternal(gPrime, hPrime, u, newP, aPrime, bPrime, l, r)
}

// VerifyIPProof validates inner product proof.
func VerifyIPProof(params IPParams, p *utils.ECPoint, c *big.Int, proof *IPProof) bool {
	gv := params.GV()
	hv := params.HV()
	u := params.U()
	ue, np, err := basicCheckAndCal(gv, hv, u, p, c, proof)
	if err != nil {
		log.Debug("optimized inner product proof invalid", "err", err)
		return false
	}

	return verifyIPProof(gv, hv, ue, np, proof)
}

// verifyIPProof validates inner product proof.
func verifyIPProof(g, h *utils.GeneratorVector, u, p *utils.ECPoint, proof *IPProof) bool {
	curve := u.Curve
	n := curve.Params().N

	for i, l := range proof.l {
		gLeft := g.HalfLeft()
		gRight := g.HalfRight()
		hLeft := h.HalfLeft()
		hRight := h.HalfRight()

		r := proof.r[i]

		e, err := utils.ComputeChallenge(n, l.X, l.Y, r.X, r.Y)
		if err != nil {
			log.Warn("IPVerifier compute challenge failed in protocol 2", "error", err)
			return false
		}
		eInverse := new(big.Int).ModInverse(e, n)

		gPrime := gLeft.HadamardScalar(eInverse).AddGeneratorVector(gRight.HadamardScalar(e))
		hPrime := hLeft.HadamardScalar(e).AddGeneratorVector(hRight.HadamardScalar(eInverse))

		// Compute e ^ 2.
		eSquare := new(big.Int).Mul(e, e)
		eSquare.Mod(eSquare, n)
		// compute e ^ -2 = eInverse ^ 2.
		eInverseSquare := new(big.Int).Mul(eInverse, eInverse)
		eInverseSquare.Mod(eInverseSquare, n)

		// update params.
		np := new(utils.ECPoint).ScalarMult(l, eSquare)
		np.Add(np, p)
		rTmp := new(utils.ECPoint).ScalarMult(r, eInverseSquare)
		np.Add(np, rTmp)

		// set new params.
		p = np.Copy()
		g = gPrime
		h = hPrime
	}

	if g.Size() != 1 {
		log.Warn("IPVerifier g generator size != 1")
		return false
	}

	if h.Size() != 1 {
		log.Warn("IPVerifier h generator size != 1")
		return false
	}

	c := new(big.Int).Mul(proof.a, proof.b)
	c.Mod(c, n)

	// compute u * c.
	want := new(utils.ECPoint).ScalarMult(u, c)
	// compute g * a.
	ga := new(utils.ECPoint).ScalarMult(g.Get(0), proof.a)
	// compute h * b.
	hb := new(utils.ECPoint).ScalarMult(h.Get(0), proof.b)
	// compute g * a + h * b + u * c.
	want.Add(want, ga)
	want.Add(want, hb)

	if !p.Equal(want) {
		log.Warn("Verifier p != p1", "want x", p.X, "want y", p.Y, "actual x", want.X, "actual y", want.Y)
		return false
	}

	return true
}

// OptimizedVerifyIPProof verify inner product proof using multi-exponentiation.
func OptimizedVerifyIPProof(params IPParams, p *utils.ECPoint, c *big.Int, proof *IPProof) bool {
	gv := params.GV()
	hv := params.HV()
	u := params.U()
	ue, np, err := basicCheckAndCal(gv, hv, u, p, c, proof)
	if err != nil {
		log.Debug("optimized inner product proof invalid", "err", err)
		return false
	}

	return optmizedVerifyIPProof(gv, hv, ue, np, proof)
}

func basicCheckAndCal(gv, hv *utils.GeneratorVector, u, p *utils.ECPoint, c *big.Int, proof *IPProof) (*utils.ECPoint, *utils.ECPoint, error) {
	if gv.Size() != hv.Size() {
		return nil, nil, errors.New("invalid ip proof: gv, hv size not equal")
	}

	if len(proof.l) != len(proof.r) {
		return nil, nil, errors.New("invalid ip proof: l, r's lenght not equal")
	}

	// todo: //
	// check 2^len(proof.r) == len(gv)

	ue, np, err := cal(p, u, c)
	if err != nil {
		return nil, nil, err
	}

	return ue, np, nil
}

//
func optmizedVerifyIPProof(g, h *utils.GeneratorVector, u, p *utils.ECPoint, proof *IPProof) bool {
	curve := u.Curve
	n := curve.Params().N

	right := p.Copy()
	xjs := make([]*big.Int, 0)
	xjsInv := make([]*big.Int, 0)
	for i := 0; i < len(proof.l); i++ {
		xj, err := utils.ComputeChallenge(n, proof.l[i].X, proof.l[i].Y, proof.r[i].X, proof.r[i].Y)
		if err != nil {
			log.Warn("compute challenge for optimize inner product failed", "err", err)
			return false
		}
		xjInv := new(big.Int).ModInverse(xj, n)
		xjs = append(xjs, xj)
		xjsInv = append(xjsInv, xjInv)

		xj2 := new(big.Int).Mul(xj, xj)
		xj2.Mod(xj2, n)

		xj2Inv := new(big.Int).Mul(xjInv, xjInv)
		xj2Inv.Mod(xj2Inv, n)

		right.Add(right, new(utils.ECPoint).ScalarMult(proof.l[i], xj2))
		right.Add(right, new(utils.ECPoint).ScalarMult(proof.r[i], xj2Inv))
	}

	left := u.Copy()
	ab := new(big.Int).Mul(proof.a, proof.b)
	ab.Mod(ab, n)
	left.ScalarMult(left, ab)

	s := make([]*big.Int, g.Size())

	for i := 0; i < g.Size(); i++ {
		for j := 0; j < len(proof.l); j++ {
			tmp := new(big.Int)
			if utils.SmallParseBinary(i, j, len(proof.r)) {
				tmp.Set(xjs[j])
			} else {
				tmp.Set(xjsInv[j])
			}

			if j == 0 {
				s[i] = tmp
			} else {
				s[i].Mul(s[i], tmp)
				s[i].Mod(s[i], n)
			}
		}
	}

	as := utils.NewFieldVector(s, n).Times(proof.a).GetVector()
	left.Add(left, g.Commit(as))
	bsinv := utils.NewFieldVector(s, n).ModInverse().Times(proof.b).GetVector()
	left.Add(left, h.Commit(bsinv))

	return left.Equal(right)
}
