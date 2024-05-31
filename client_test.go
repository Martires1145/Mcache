package Mcache

import (
	"fmt"
	"github.com/Martires1145/Mcache/singleflight"
	"log"
	"sync"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestClient_Get(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := New(2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestSingleFlight(t *testing.T) {
	sg := &singleflight.Group{}
	cnt, all := 0, 1000
	wg1, wg2 := sync.WaitGroup{}, sync.WaitGroup{}
	wg1.Add(all)
	wg2.Add(all)
	for i := 0; i < all; i++ {
		go func() {
			wg1.Wait()
			sg.Do("key", func() (any, error) {
				cnt++
				return nil, nil
			})
			wg2.Done()
		}()
		wg1.Done()
	}
	wg2.Wait()

	if cnt != 1 {
		t.Fatalf("failed")
	}
}
