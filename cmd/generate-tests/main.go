// Copyright (c) 2021 mineruniter969

package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	ethash "github.com/mineruniter969/ethash1.1-go"
)

func main() {
	headerRLP := "0xf901f7a01bef91439a3e070a6586851c11e6fd79bbbea074b2b836727b8e75c7d4a6b698a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794ea3cb5f94fa2ddd52ec6dd6eb75cf824f4058ca1a00c6e51346be0670ce63ac5f05324e27d20b180146269c5aab844d09a2b108c64a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302004002832fefd880845511ed2a80a0e55d02c555a7969361cf74a9ec6211d8c14e4517930a00442f171bdb1698d17588307692cf71b12f6d"
	nonce := uint64(0x307692cf71b12f6d)
	var header *types.Header

	err := rlp.DecodeBytes(hexutil.MustDecode(headerRLP), &header)
	if err != nil {
		panic(err)
	}

	fmt.Println(header.Number)

	seedHash := ethash.SeedHash(header.Number.Uint64())
	cache := make([]uint32, ethash.CacheSize(header.Number.Uint64())/4)
	ethash.GenerateCache(cache, header.Number.Uint64()/ethash.EpochLength, seedHash)

	sealHash := ethash.SealHash(header)

	fmt.Println("seed hash:", hexutil.Encode(seedHash))
	fmt.Println("cache size:", ethash.CacheSize(header.Number.Uint64()))
	fmt.Println("dataset size:", ethash.DatasetSize(header.Number.Uint64()))
	fmt.Println("header hash:", sealHash)

	digest, powHash := ethash.HashimotoLight(ethash.DatasetSize(header.Number.Uint64()), cache, sealHash.Bytes(), nonce)

	fmt.Println("pow hash", hexutil.Encode(powHash))
	fmt.Println("mix digest", hexutil.Encode(digest))
}
