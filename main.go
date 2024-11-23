package main

import (
	"MayCache/mcache"
	"MayCache/server"
	"flag"
	"fmt"
	"log"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *mcache.Group {
	return mcache.NewGroup("scores", 2<<10, mcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Maycache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	may := createGroup()
	if api {
		may.AddFilter()
		go server.StartAPIServer(apiAddr, may)
	}
	server.StartCacheServer(addrMap[port], addrs, may)
}
