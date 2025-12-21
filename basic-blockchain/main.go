package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := NewGenesisBlock()
	return &BlockChain{[]*Block{genesisBlock}}
}

func (bc *BlockChain) Add(data []byte) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	b := NewBlock(data, prevHash)
	bc.Blocks = append(bc.Blocks, b)
}

func NewBlock(data, prevBlockHash []byte) *Block {
	b := &Block{time.Now().UTC(), data, prevBlockHash, []byte{}}
	b.SetHash()
	return b
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"), []byte{})
}

type Block struct {
	Timestamp     time.Time
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	ts := []byte(strconv.FormatInt(b.Timestamp.Unix(), 10))
	bts := bytes.Join([][]byte{ts, b.Data, b.PrevBlockHash}, []byte{})
	hash := sha256.Sum256(bts)

	b.Hash = hash[:]
}

func main() {
	bc := NewBlockChain()
	bc.Add([]byte("First Block"))

	for i, b := range bc.Blocks {
		fmt.Println("data", string(b.Data))
		fmt.Println("Hash", string(b.Hash))
		fmt.Println("PrevBlockHash", string(b.PrevBlockHash))
		fmt.Println("==end==", i)
	}
}
