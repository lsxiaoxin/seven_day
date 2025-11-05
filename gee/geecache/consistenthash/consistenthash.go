package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash             //哈希函数 
	replicas int              //虚拟节点和真节点的倍数
	keys     []int            //保存虚拟节点
	hashMap  map[int]string   //映射虚拟节和真节点
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//添加节点，参数是真节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)	
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

//获取到真实的节点，参数是输入的key
func (m *Map) Get(key string) string{
	if (len(m.keys) == 0) {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	//二分查找，找出第一个下标, 满足 m.keys[idx] >= hash
	//查不到会返回一个大于等于 len(m.keys) 的索引，所以要 % len
	//m.keys[] 得到虚拟节点， hashmap得到真实真节点的名称
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	} )

	return m.hashMap[m.keys[idx % len(m.keys)]]
}