package bloomfilter

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

type HashFunc func(string) uint64

type BloomFilter struct {
	bitsPerKey uint64
	n          uint64
	bitsList   []uint64
	hashList   []HashFunc
}

func New(n uint64, bitsPerKey uint64, hashFunc ...HashFunc) *BloomFilter {
	bits := n * bitsPerKey
	if bits < 64 {
		bits = 64
	}
	bloom := &BloomFilter{
		bitsPerKey: bitsPerKey,
		n:          n,
		bitsList:   make([]uint64, (bits+63)/64),
		hashList:   nil,
	}
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
		index := hash(key) % uint64(len(bf.bitsList))
		bf.bitsList[index] |= 1 << (index % 64)
	}
}

func (bf *BloomFilter) MightContain(key string) bool {
	for _, hash := range bf.hashList {
		index := hash(key) % uint64(len(bf.bitsList))
		if bf.bitsList[index]&(1<<(index%64)) == 0 {
			return false
		}
	}
	return true
}
