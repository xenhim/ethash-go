// Copyright (c) 2021 mineruniter969

package ethash

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"
)

type Hasher struct {
	mu sync.Mutex

	caches map[int][]uint32
}

func NewHasher() *Hasher {
	return &Hasher{
		caches: make(map[int][]uint32),
	}
}

func (h *Hasher) getCache(epoch int) []uint32 {
	h.mu.Lock()
	defer h.mu.Unlock()

	cache, ok := h.caches[epoch]
	if ok {
		return cache
	}

	// cache does not exist, generating

	cache = make([]uint32, CalcCacheSize(epoch)/4)

	GenerateCache(cache, uint64(epoch), CalcSeedHash(epoch))

	h.caches[epoch] = cache

	return cache
}

// Verify verifies block seal. This method returns nil if the block is correct, and a non-nil error if the block is incorrect.
func (h *Hasher) Verify(number uint64, sealHash []byte, mixDigest []byte, nonce uint64, difficulty *big.Int) error {
	epoch := int(number / EpochLength)
	if epoch >= maxEpoch {
		// epoch too high
		return fmt.Errorf("epoch too high")
	}

	cache := h.getCache(epoch)

	mixDigestComputed, powHash := HashimotoLight(CalcDatasetSize(epoch), cache, sealHash, nonce)
	if bytes.Equal(mixDigest, mixDigestComputed) {
		// invalid mix digest
		return fmt.Errorf("invalid mix digest")
	}

	computedTarget := new(big.Int).SetBytes(powHash)
	computedDifficulty := new(big.Int).Div(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), computedTarget)

	if difficulty.Cmp(computedDifficulty) < 0 {
		// valid
		return nil
	}

	return fmt.Errorf("proof-of-work solution is invalid")
}

// VerifyWithoutMix accomplishes the same task as Verify does, but without the mix digest verification
func (h *Hasher) VerifyWithoutMix(number uint64, sealHash []byte, nonce uint64, difficulty *big.Int) error {
	epoch := int(number / EpochLength)
	if epoch >= maxEpoch {
		// epoch too high
		return fmt.Errorf("epoch too high")
	}

	cache := h.getCache(epoch)

	_, powHash := HashimotoLight(CalcDatasetSize(epoch), cache, sealHash, nonce)

	computedTarget := new(big.Int).SetBytes(powHash)
	computedDifficulty := new(big.Int).Div(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), computedTarget)

	if difficulty.Cmp(computedDifficulty) < 0 {
		// valid
		return nil
	}

	return fmt.Errorf("proof-of-work solution is invalid")
}
