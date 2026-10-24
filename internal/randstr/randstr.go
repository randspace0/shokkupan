package randstr

import (
	"crypto/rand"
	"math/big"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// String returns n random letters from a CSPRNG; it is used for setup keys,
// recovery codes and ids, so it must not use math/rand.
func String(n int) string {
	max := big.NewInt(int64(len(letters)))
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b)
}
