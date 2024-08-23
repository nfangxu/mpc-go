package ecdh

import (
	"crypto/elliptic"
	"crypto/sha256"
	"github.com/pkg/errors"
	"math/big"
)

type Ecdh struct {
	Curve elliptic.Curve
}

// GetPoints 获取曲线上的点集合
func (e *Ecdh) GetPoints(data [][]byte) ([]*big.Int, []*big.Int, error) {
	hashes := e.sha256(data)

	xs, ys := make([]*big.Int, 0), make([]*big.Int, 0)

	for _, hash := range hashes {
		val := new(big.Int).SetBytes(hash)
		x, y, err := e.point(e.Curve, val)
		if err != nil {
			return nil, nil, err
		}
		xs = append(xs, x)
		ys = append(ys, y)
	}
	return xs, ys, nil
}

// Exp 数据与Key相乘，返回结果
func (e *Ecdh) Exp(xs, ys []*big.Int, key *big.Int) ([]*big.Int, []*big.Int) {
	return e.exp(e.Curve, xs, ys, key)
}

// Intersection 求交获取数据索引
func (e *Ecdh) Intersection(xs, _xs []*big.Int) []int {
	hashSet := make(map[string]bool)
	for _, x := range _xs {
		hashSet[string(x.Bytes())] = true
	}

	idx := make([]int, 0)
	for i, x := range xs {
		if _, ok := hashSet[string(x.Bytes())]; ok {
			idx = append(idx, i)
		}
	}
	return idx
}

func (e *Ecdh) sha256(data [][]byte) [][]byte {
	hashes := make([][]byte, 0)
	for _, datum := range data {
		sha := sha256.Sum256(datum)
		hashes = append(hashes, sha[:])
	}
	return hashes
}

// polynomial returns x³ - 3x + b.
func (e *Ecdh) polynomial(curve elliptic.Curve, x *big.Int) *big.Int {
	z := new(big.Int).Mul(x, x)
	z.Mul(z, x)

	_z := new(big.Int).Lsh(x, 1)
	_z.Add(_z, x)

	z.Sub(z, _z)
	z.Add(z, curve.Params().B)
	z.Mod(z, curve.Params().P)

	return z
}

func (e *Ecdh) calcY(curve elliptic.Curve, x *big.Int, y *big.Int) bool {
	res := y.ModSqrt(e.polynomial(curve, x), curve.Params().P)
	return res != nil
}

func (e *Ecdh) point(curve elliptic.Curve, x *big.Int) (*big.Int, *big.Int, error) {
	px := new(big.Int).Set(x)
	py := new(big.Int)
	one := big.NewInt(1)

	px.Lsh(px, 8)
	px.Mod(px, curve.Params().P)
	for i := 0; i < (2 << 8); i++ {
		if e.calcY(curve, px, py) {
			if !curve.IsOnCurve(px, py) {
				return nil, nil, errors.New("not on curve")
			}
			return px, py, nil
		}
		px.Add(px, one)
	}
	return nil, nil, errors.New("invalid")
}

func (e *Ecdh) exp(curve elliptic.Curve, xs, ys []*big.Int, key *big.Int) ([]*big.Int, []*big.Int) {
	_xs, _ys := make([]*big.Int, len(xs)), make([]*big.Int, len(xs))
	for i := 0; i < len(xs); i++ {
		x, y := curve.ScalarMult(xs[i], ys[i], key.Bytes())
		_xs[i] = x
		_ys[i] = y
	}
	return _xs, _ys
}
