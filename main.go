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
	prevhash  []byte   //이전 블럭의 헤더 해쉬
	merkle    []byte   //Merkle
	timestamp int      //타임스탬프, 블럭이 생성된 대략적인 시간을 표시
	nonce     int      //nonce
	target    *big.Int //타겟 난이도
}

func (b *Block) get_header_hash() []byte {
	//블럭의 PoW 알고리즘을 위해 현재 블럭의 헤더 정보를 해쉬화함.
	header_hash := bytes.Join([][]byte{
		b.prevhash,
		b.merkle,
		[]byte(strconv.Itoa(b.timestamp)),
		[]byte(strconv.Itoa(b.nonce)),
	}, []byte{})

	return header_hash
}

func (b *Block) make_target(target int) {
	//난이도 값 설정.
	x := big.NewInt(1)
	x.Lsh(x, (uint)(256-target))

	b.target = x
}

func (b *Block) mine() (int, []byte) {
	//PoW 알고리즘 시행.
	var fin_hash [32]byte //최종 결과값을 담을 변수.

	for b.nonce < math.MaxInt64 { //Nonce값이 정수형의 최대가 될 때까지 실행, 만일 이상 넘어 갈 시 overflow.
		var hash [32]byte = sha256.Sum256(b.get_header_hash()) //헤더 정보를 SHA256으로 해쉬화.

		var hash_to_int big.Int
		hash_to_int.SetBytes(hash[:]) //해쉬화 된 값을 변환.

		if hash_to_int.Cmp(b.target) == -1 { //변환된 해쉬 값과 난이도 값을 비교. hash가 target보다 작을 경우.
			fin_hash = hash //해쉬를 최종화.
			fmt.Println("Current hash: ", fin_hash)
			break //루프 탈출.
		}

		b.nonce++ //Nonce값 증가.
	}

	return b.nonce, fin_hash[:]
}

func (b *Block) vaildate() bool {
	//블럭의 무결성 검증.
	var raw_hash = b.get_header_hash()
	var hashed = sha256.Sum256(raw_hash) //블럭 데이터 해쉬화

	var hash_to_int big.Int
	hash_to_int.SetBytes(hashed[:])

	return hash_to_int.Cmp(b.target) == -1 //블럭의 해쉬를 통해 무결성을 검증함
}

func build_random(length int) []byte {
	//랜덤한 바이트 생성.
	rand.Seed(time.Now().Unix())
	seed := "abcedfghijklmnopqrstunvwxyzABCEDFGHIJKLMNOPQRSTUNVWXYZ"
	buffer := make([]byte, length)
	for i := range buffer {
		buffer[i] = seed[rand.Intn(len(seed))]
	}

	return buffer
}

func gen_genesis_block() Block {
	//제네시스 블럭(초기 블럭) 생성.
	block := Block{
		prevhash:  []byte{},
		merkle:    build_random(32),
		timestamp: int(time.Now().Unix()),
		nonce:     0,
	}

	return block
}

func gen_block(b *Block) Block {
	//블럭 생성.
	block := Block{
		prevhash:  hash_block(b),
		merkle:    build_random(32),
		timestamp: int(time.Now().Unix()),
		nonce:     0,
	}

	return block
}

func hash_block(b *Block) []byte {
	//블럭 데이터를 해쉬화.
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

	var target_diff int = 24 //현재 난이도.
	var prev_block *Block    //기존 블럭을 저장.
	var time_total []time.Duration
	fmt.Println("Stating.. in count: ", count)
	fmt.Print("\n\n")

	for i := 0; i < count; i++ {
		started_time := time.Now() //현재시간.
		if i <= 0 {                //현재 생성해야 할 블럭이 제네시스 블럭인지 확인. 만일 loop count가 1보다 작거나 큰가? (첫 실행인가?)
			tblock := gen_genesis_block() //제네시스 블럭 생성.
			fmt.Println("Current Header: ", tblock.get_header_hash())
			tblock.make_target(target_diff) //현재 난이도 설정.
			fmt.Println("== Mining Start ==")
			tblock.mine() //PoW 알고리즘 시행.
			fmt.Println("== Finished ==")
			fmt.Println("Nonce: ", tblock.nonce)
			fmt.Println("Valid Status", tblock.vaildate()) //결과값을 검증하고, 검증 결과를 출력
			prev_block = &tblock                           //다음 블럭을 위해 블럭 데이터를 저장.

		}

		vblock := gen_block(prev_block) //일반 블럭 생성. 제네시스 블럭이 아닐 경우엔 기존 블럭의 정보가 필요함.
		fmt.Println("Current Header: ", vblock.get_header_hash())
		vblock.make_target(target_diff) //현재 난이도 설정.
		fmt.Println("== Mining Start ==")
		vblock.mine() //PoW 알고리즘 시행.
		fmt.Println("== Finished ==")
		fmt.Println("Nonce: ", vblock.nonce)
		fmt.Println("Valid Status", vblock.vaildate()) //결과값을 검증하고, 검증 결과를 출력.
		prev_block = &vblock                           //다음 블럭을 위해 블럭 데이터를 저장.

		eltime := time.Since(started_time) //알고리즘 시행 시간을 측정
		//난이도 책정
		if eltime.Microseconds() < 1800000 {
			target_diff = target_diff + 2
		} else if eltime.Microseconds() > 1800000 {
			target_diff = target_diff - 2
		}

		fmt.Println("[Diff] ", target_diff)
		fmt.Println("[Cnt] ", i)
		fmt.Print("\n\n")
		time_total = append(time_total, time.Duration(eltime.Milliseconds()))
	}

	fmt.Println("Finished, eltime: ", time_total)
	//버그픽스중
	var time_all_total int64 = 0
	for i := range time_total {
		time_all_total = time_all_total + time_total[i].Milliseconds()
	}
	eltime_total := time_all_total / int64(len(time_total))

	fmt.Println("\n[All-Total] ", eltime_total)

}
