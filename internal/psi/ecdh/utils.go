package ecdh

import (
	"crypto/elliptic"
	"math/big"
	"math/rand"
)

var defaultCurve = elliptic.P256()
var key *big.Int

// Curve 设置默认曲线
func Curve(curve elliptic.Curve) {
	defaultCurve = curve
}

func Key() *big.Int {
	if key == nil {
		s := big.NewInt(rand.Int63n(10) + 5)
		key = s
	}
	return key
}

// GetPoints 第一步：获取曲线上的点集合
func GetPoints(data [][]byte) ([]*big.Int, []*big.Int, error) {
	ecdh := &Ecdh{Curve: defaultCurve}
	return ecdh.GetPoints(data)
}

// Exp 第二步：数据与Key相乘，返回结果
func Exp(xs, ys []*big.Int, key *big.Int) ([]*big.Int, []*big.Int) {
	ecdh := &Ecdh{Curve: defaultCurve}
	return ecdh.Exp(xs, ys, key)
}

// Intersection 第三步：求交获取数据索引
func Intersection(xs, _xs []*big.Int) []int {
	ecdh := &Ecdh{Curve: defaultCurve}
	return ecdh.Intersection(xs, _xs)
}
