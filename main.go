package main

//TODO:
/*
	Multi-Threading,
	Document
*/

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

type Block struct {
	prevhash  []byte
	merkle    []byte
	timestamp int
	nonce     int
	target    *big.Int
}

func (b *Block) get_header_hash() []byte {
	header_hash := bytes.Join([][]byte{
		b.prevhash,
		b.merkle,
		[]byte(strconv.Itoa(b.timestamp)),
		[]byte(strconv.Itoa(b.nonce)),
	}, []byte{})

	return header_hash
}

func (b *Block) make_target(target int) {
	x := big.NewInt(1)
	x.Lsh(x, (uint)(256-target))

	b.target = x
}

func (b *Block) mine() (int, []byte) {
	var fin_hash [32]byte

	for b.nonce < math.MaxInt64 {
		var hash [32]byte = sha256.Sum256(b.get_header_hash())

		var hash_to_int big.Int
		hash_to_int.SetBytes(hash[:])

		if hash_to_int.Cmp(b.target) == -1 {
			fin_hash = hash
			fmt.Println("Current hash: ", fin_hash)
			break
		}

		b.nonce++
	}

	return b.nonce, fin_hash[:]
}

func (b *Block) vaildate() bool {
	var raw_hash = b.get_header_hash()
	var hashed = sha256.Sum256(raw_hash)

	var hash_to_int big.Int
	hash_to_int.SetBytes(hashed[:])

	return hash_to_int.Cmp(b.target) == -1
}

func build_random(length int) []byte {
	seed := "abcedfghijklmnopqrstunvwxyzABCEDFGHIJKLMNOPQRSTUNVWXYZ"
	buffer := make([]byte, length)
	for i := range buffer {
		buffer[i] = seed[rand.Intn(len(seed))]
	}

	return buffer
}

func gen_genesis_block() Block {
	block := Block{
		prevhash:  []byte{},
		merkle:    build_random(32),
		timestamp: int(time.Now().Unix()),
		nonce:     0,
	}

	return block
}

func gen_block(b *Block) Block {
	block := Block{
		prevhash:  hash_block(b),
		merkle:    build_random(32),
		timestamp: int(time.Now().Unix()),
		nonce:     0,
	}

	return block
}

func hash_block(b *Block) []byte {
	header_hash := bytes.Join([][]byte{
		b.merkle,
		[]byte(strconv.Itoa(b.timestamp)),
		[]byte(strconv.Itoa(b.nonce)),
	}, []byte{})

	return header_hash
}

func main() {

	fmt.Println("Count?")
	var count int = 0
	fmt.Scanln(&count)

	var target_diff int = 24
	var prev_block *Block
	var time_total []time.Duration
	fmt.Println("Stating.. in count: ", count)
	fmt.Print("\n\n")

	for i := 0; i < count; i++ {
		started_time := time.Now()
		if i <= 0 {
			tblock := gen_genesis_block()
			fmt.Println("Current Header: ", tblock.get_header_hash())
			tblock.make_target(target_diff)
			fmt.Println("== Mining Start ==")
			tblock.mine()
			fmt.Println("== Finished ==")
			fmt.Println("Nonce: ", tblock.nonce)
			fmt.Println("Valid Status", tblock.vaildate())
			prev_block = &tblock
		}

		vblock := gen_block(prev_block)
		fmt.Println("Current Header: ", vblock.get_header_hash())
		vblock.make_target(target_diff)
		fmt.Println("== Mining Start ==")
		vblock.mine()
		fmt.Println("== Finished ==")
		fmt.Println("Nonce: ", vblock.nonce)
		fmt.Println("Valid Status", vblock.vaildate())
		prev_block = &vblock

		eltime := time.Since(started_time)
		if eltime.Nanoseconds() < 180000000 {
			target_diff = target_diff + 2
		} else if eltime.Nanoseconds() > 180000000 {
			target_diff = target_diff - 2
		}

		fmt.Println("[Diff] ", target_diff)
		fmt.Println("[Cnt] ", count)
		fmt.Print("\n\n")
		time_total = append(time_total, time.Duration(eltime.Milliseconds()))
	}

	fmt.Println("Finished, eltime: ", time_total)

	var time_all_total int64 = 0
	for i := range time_total {
		time_all_total = time_all_total + time_total[i].Milliseconds()
	}
	eltime_total := time_all_total / int64(len(time_total))

	fmt.Println("\n[All-Total] ", eltime_total)

}
