package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	body           map[string]interface{}
	hashID         string // the hash of a block is its identifier generated using cryptography
	previousHashID string
	ts             time.Time
	proofOfWork    int
}

type Blockchain struct {
	genesis    Block // first block of chain
	chain      []Block
	difficulty int
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.body)
	blockData := b.previousHashID + string(data) + b.ts.String() + strconv.Itoa(b.proofOfWork)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hashID, strings.Repeat("0", difficulty)) {
		b.proofOfWork++
		b.hashID = b.calculateHash()
	}
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hashID: "0",
		ts:     time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

func (b *Blockchain) addBlock(data interface{}) {
	blockData := map[string]interface{}{
		"data": data,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		body:           blockData,
		previousHashID: lastBlock.hashID,
		ts:             time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b *Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hashID != currentBlock.calculateHash() || currentBlock.previousHashID != previousBlock.hashID {
			return false
		}
	}
	return true
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain
	blockchain.addBlock(map[string]string{
		"from":  "ABC",
		"to":    "DEF",
		"sends": "2",
	})
	blockchain.addBlock(map[string]string{
		"from":  "DEF",
		"to":    "XGF",
		"sends": "2",
	})

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
}
