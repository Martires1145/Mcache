package bloomfilter

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/rand"
	"time"
)

type HashFunc func(string) uint64

type BloomFilter struct {
	bitsPerKey uint64
	n          uint64
	bitList    BitList
	hashList   []HashFunc
}

func New(n uint64, bitsPerKey uint64, bitList BitList, hashFunc ...HashFunc) *BloomFilter {
	bits := n * bitsPerKey
	if bits < 64 {
		bits = 64
	}
	bloom := &BloomFilter{
		bitsPerKey: bitsPerKey,
		n:          n,
		bitList:    bitList,
		hashList:   nil,
	}
	bloom.bitList.SetCapacity(bits)

	if len(hashFunc) == 0 {
		rand.Seed(time.Now().UnixNano())
		bloom.hashList = defaultHashFunc(int(bitsPerKey))
	} else {
		bloom.hashList = hashFunc
	}
	return bloom
}

func defaultHashFunc(k int) []HashFunc {
	var hashFunc []HashFunc
	if k < 1 {
		k = 1
	} else if k > 30 {
		k = 30
	}
	for i := 0; i < k; i++ {
		salt := uint64(rand.Int63())
		hashFunc = append(hashFunc, func(data string) uint64 {
			h := sha256.New()
			h.Write([]byte(data))
			sum := h.Sum(nil)
			// 使用SHA-256的前8个字节并添加盐值
			return binary.BigEndian.Uint64(sum[:8]) + salt
		})
	}
	return hashFunc
}

func (bf *BloomFilter) Add(key string) {
	for _, hash := range bf.hashList {
		index := hash(key) % bf.bitList.Len()
		err := bf.bitList.Set(index)
		if err != nil {
			log.Fatal("add error")
		}
	}
}

func (bf *BloomFilter) MightContain(key string) bool {
	for _, hash := range bf.hashList {
		index := hash(key) % bf.bitList.Len()
		ok, err := bf.bitList.Check(index)
		if !ok || err != nil {
			return false
		}
	}
	return true
}
