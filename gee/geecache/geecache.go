package geecache

import (
	"fmt"
	"log"
	"sync"
)

//这样设计 既可以传入实现了Get方法的结构体
//又可以只传入函数
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name string
	getter Getter
	mainCache cache

	peers  PeerPicker
}

var (
	mutex sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group{
	if (getter == nil) {
		panic("nil getter")
	}
	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mutex.RLock()
	g := groups[name]
	mutex.RUnlock()
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}


func (g *Group) Get(key string) (ByteView, error) {
	if (key == "") {
		return ByteView{}, fmt.Errorf("key is empty") 
	}

	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	return  g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {

			log.Printf("getter no peer") //报错

			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, err
			}
		}
	}
	


	return g.getLocally(key);
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b : cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)	
}

var _ PeerPicker = (*HTTPPool)(nil)