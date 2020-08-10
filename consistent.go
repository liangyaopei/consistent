package consistent

import (
	"sort"
	"strconv"
	"sync"
)

// New returns a hash ring
// with customized hash function
func New(nodes []string, hashFn HashFn) *HashRing {
	ring := &HashRing{
		nodeWeights: make(map[string]int),
		ring:        make(map[uint64]string),
		sortedKey:   make([]uint64, 0),
		hashFn:      hashFn,
	}
	ring.lock = &sync.RWMutex{}
	ring.generateRing(nodes)
	return ring
}

// NewWithWeight returns a hash ring
// with customized weight
func NewWithWeight(nodeWeight map[string]int, hashFn HashFn) *HashRing {
	ring := &HashRing{
		nodeWeights: nodeWeight,
		ring:        make(map[uint64]string),
		sortedKey:   make([]uint64, 0),
		hashFn:      hashFn,
	}
	ring.lock = &sync.RWMutex{}
	nodes := make([]string, 0, len(nodeWeight))
	for node := range nodeWeight {
		nodes = append(nodes, node)
	}
	ring.generateRing(nodes)
	return ring
}

// HashRing represents hash ring
type HashRing struct {
	lock *sync.RWMutex // lock

	nodeWeights map[string]int    // stores node's weight
	ring        map[uint64]string // hash ring stores key-node pair
	sortedKey   []uint64          // stores node's replication hash val in ascending order
	hashFn      HashFn            // hash function to generate byte array's uint64 hash value
}

// AddNode add new node with weight 1
func (r *HashRing) AddNode(node string) *HashRing {
	return r.AddNodeWeight(node, 1)
}

// AddNodeWeight add new node with weight
func (r *HashRing) AddNodeWeight(node string, weight int) *HashRing {
	if weight <= 0 {
		return r
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.nodeWeights[node]; ok {
		return r
	}
	// add new node
	r.addNodeWeightIdx(node, weight, 0)
	r.nodeWeights[node] = weight
	return r
}

func (r *HashRing) addNodeWeightIdx(node string, weight int, idx int) {
	// add new node with weight
	for i := idx; i < weight; i++ {
		nodeKey := r.genNodeKeyIdxStr(node, i)
		key := r.hashFn(str2Bytes(nodeKey))
		//idx := sort.Search(len(r.sortedKey), func(i int) bool {
		//	return r.sortedKey[i] >= key
		//})
		r.sortedKey = append(r.sortedKey, key)
		// copy(r.sortedKey[idx+1:], r.sortedKey[idx:])
		// r.sortedKey[idx] = key
		r.ring[key] = node
	}
	// sort node in ascending order
	sort.Slice(r.sortedKey, func(i, j int) bool {
		return r.sortedKey[i] < r.sortedKey[j]
	})
}

// DelNode deletes node in hash ring
func (r *HashRing) DelNode(node string) *HashRing {
	r.lock.Lock()
	defer r.lock.Unlock()

	weight, ok := r.nodeWeights[node]
	if !ok {
		return r
	}
	r.delNodeIdx(node, weight, 0)
	delete(r.nodeWeights, node)
	return r
}

func (r *HashRing) delNodeIdx(node string, weight int, idx int) {
	for i := idx; i < weight; i++ {
		nodeKey := r.genNodeKeyIdxStr(node, i)
		key := r.hashFn(str2Bytes(nodeKey))
		// binary search find the key
		idx := sort.Search(len(r.sortedKey), func(i int) bool {
			return r.sortedKey[i] >= key
		})
		// key found
		if idx < len(r.sortedKey) && key == r.sortedKey[idx] {
			// delete the key
			copy(r.sortedKey[idx:], r.sortedKey[idx+1:])
			r.sortedKey = r.sortedKey[:len(r.sortedKey)-1]
		}
		// delete key in ring
		delete(r.ring, key)
	}
}

// UpdateNodeWeight update node' weight
// it can increase or decrease the
// node's replication
func (r *HashRing) UpdateNodeWeight(node string, weight int) *HashRing {
	if weight < 0 {
		return r
	}
	r.lock.Lock()
	defer r.lock.Unlock()

	oldWeight, ok := r.nodeWeights[node]
	if ok || oldWeight == weight {
		return r
	}
	// add node replication
	if oldWeight < weight {
		r.addNodeWeightIdx(node, weight, oldWeight)
	}

	// del node replication
	if oldWeight > weight {
		r.delNodeIdx(node, weight, oldWeight)
	}
	// update weight
	r.nodeWeights[node] = weight
	return r
}

func (r *HashRing) generateRing(nodes []string) {
	for _, node := range nodes {

		weight, ok := r.nodeWeights[node]
		if !ok {
			weight = 1
			r.nodeWeights[node] = weight
		}

		for i := 0; i < weight; i++ {
			nodeKey := r.genNodeKeyIdxStr(node, i)
			key := r.hashFn(str2Bytes(nodeKey))
			r.ring[key] = node
			r.sortedKey = append(r.sortedKey, key)
		}
	}
	// sort node in ascending order
	sort.Slice(r.sortedKey, func(i, j int) bool {
		return r.sortedKey[i] < r.sortedKey[j]
	})
}

func (r *HashRing) genNodeKeyIdxStr(node string, idx int) string {
	return node + "%%" + strconv.Itoa(idx)
}

// LocateKey locates the position of
// byte array data in hash ring
func (r *HashRing) LocateKey(data []byte) string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	dataKey := r.hashFn(data)
	pos := sort.Search(len(r.sortedKey), func(i int) bool {
		return r.sortedKey[i] >= dataKey
	})
	if pos == len(r.sortedKey) {
		pos = 0
	}
	return r.ring[r.sortedKey[pos]]
}

// LocateKeyStr locates the position of
// string data in hash ring
func (r *HashRing) LocateKeyStr(data string) string {
	return r.LocateKey(str2Bytes(data))
}

// GetNodes returns all nodes
// in hash ring
func (r *HashRing) GetNodes() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	res := make([]string, 0, len(r.nodeWeights))
	for node := range r.nodeWeights {
		res = append(res, node)
	}
	return res
}

// GetNodeWeight returns all node and
// corresponding weight in hash ring
func (r *HashRing) GetNodeWeight() map[string]int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	res := make(map[string]int, len(r.nodeWeights))
	for k, v := range r.nodeWeights {
		res[k] = v
	}
	return res
}

// GetHashRing return hash ring
// string is node,
// []uint64 is corresponding key
func (r *HashRing) GetHashRing() map[string][]uint64 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	res := make(map[string][]uint64)
	for key, node := range r.ring {
		list, ok := res[node]
		if !ok {
			list = make([]uint64, 0)
		}
		list = append(list, key)
		res[node] = list
	}
	return res
}
