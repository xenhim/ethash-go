// Copyright (c) 2021 mineruniter969

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethash "github.com/mineruniter969/ethash1.1-go"
)

var sealHashArg = flag.String("sealhash", "", "block header seal hash")
var epoch = flag.Int("epoch", 0, "epoch")
var difficultyArg = flag.Uint64("difficulty", 0, "difficulty")

func clear0x(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s[2:]
	}

	return s
}

func searcherThread(dataset []uint32, sealHash []byte, difficulty *big.Int, start uint64) {
	nonce := start
	for {
		mixDigest, powHash := ethash.HashimotoFull(dataset, sealHash, nonce)
		computedDifficulty := new(big.Int).Div(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), new(big.Int).SetBytes(powHash))
		if computedDifficulty.Cmp(difficulty) >= 0 {
			sol := make([]byte, 8)
			new(big.Int).SetUint64(nonce).FillBytes(sol)
			fmt.Println("solution", hexutil.Encode(sol), "difficulty", computedDifficulty, "powhash", hex.EncodeToString(powHash), "mixdigest", hex.EncodeToString(mixDigest))
		}
		nonce++
	}
}

func main() {
	flag.Parse()
	sealHash, err := hex.DecodeString(clear0x(*sealHashArg))
	if err != nil {
		fmt.Println("search: invalid sealhash:", err)
		os.Exit(1)
	}

	if len(sealHash) != 32 {
		fmt.Println("search: invalid sealhash length:", len(sealHash))
		os.Exit(1)
	}

	difficulty := new(big.Int).SetUint64(*difficultyArg)

	if difficulty.Sign() <= 0 {
		fmt.Println("search: invalid difficulty")
		os.Exit(1)
	}

	threads := runtime.GOMAXPROCS(0)

	cache := make([]uint32, ethash.CalcCacheSize(*epoch)/4)
	ethash.GenerateCache(cache, uint64(*epoch), ethash.CalcSeedHash(*epoch))
	dataset := make([]uint32, ethash.CalcDatasetSize(*epoch)/4)
	ethash.GenerateDataset(dataset, uint64(*epoch), cache)
	fmt.Println("generated dag")

	for i := 0; i < threads; i++ {
		go searcherThread(dataset, sealHash, difficulty, rand.Uint64())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
