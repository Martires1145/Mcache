package bloomfilter

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := New(10000, 10, &Int64Slice{})

	bf.Add("apple")
	bf.Add("banana")
	bf.Add("cherry")

	if !bf.MightContain("apple") {
		t.Fatalf("should be true")
	}

	if !bf.MightContain("banana") {
		t.Fatalf("should be true")
	}

	if bf.MightContain("data") {
		t.Fatalf("should be false")
	}
}
