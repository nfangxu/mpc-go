package psi

import "math/big"

func BigInt2Bytes(v []*big.Int) [][]byte {
	d := make([][]byte, 0)
	for _, b := range v {
		d = append(d, b.Bytes())
	}
	return d
}

func Bytes2bigInt(v [][]byte) []*big.Int {
	d := make([]*big.Int, 0)
	for _, b := range v {
		d = append(d, new(big.Int).SetBytes(b))
	}
	return d
}
