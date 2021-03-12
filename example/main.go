// Copyright (c) 2021 mineruniter969

package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	ethash "github.com/mineruniter969/ethash1.1-go"
)

func main() {
	blockNumber := uint64(12026746)
	difficulty := big.NewInt(1000000)

	sealHash, err := hex.DecodeString("a37444aa425e5dec89d8f9afe6bedd140f3893d4f0b60528aa5c212b0b20feb5")
	if err != nil {
		panic(err)
	}

	mixDigest, err := hex.DecodeString("1effb9fd7339bbd36db2effedbd7869d93a3d3d9b8080f2e8d8e767832db0331")
	if err != nil {
		panic(err)
	}

	nonce := uint64(0x8866cb397932d01b)

	hasher := ethash.NewHasher()

	err = hasher.Verify(blockNumber, sealHash, mixDigest, nonce, difficulty)
	if err != nil {
		fmt.Println("verification failed", err)
	} else {
		fmt.Println("verification passed")
	}
}
