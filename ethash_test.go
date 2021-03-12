// Copyright (c) 2021 mineruniter969

package ethash

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestEthash(t *testing.T) {
	hasher := NewHasher()

	sealHash := hexutil.MustDecode("0x25185aec81538c09691832370742d89b4474a9454f684c979902f5f3c8de9131")
	difficulty := big.NewInt(100000)

	number := 12000000

	nonces := []uint64{
		0x6e661e92759dceef,
		0x1a02070f16a2009a,
		0xd5104dc7669b1634,
		0x68255aaf95eeda74,
	}

	for _, nonce := range nonces {
		err := hasher.VerifyWithoutMix(uint64(number), sealHash, nonce, difficulty)
		if err != nil {
			t.Errorf("verification failure at nonce %v: %v", hexutil.EncodeUint64(nonce), err)
		}
	}
}
