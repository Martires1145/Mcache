package server

import (
	"MayCache/mcache"
	"log"
	"net/http"
)

func StartCacheServer(addr string, addrs []string, may *mcache.Group) {
	peers := mcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	may.RegisterPeerPicker(peers)
	log.Println("mcache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}
