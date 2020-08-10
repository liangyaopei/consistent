package _examples

import (
	"testing"

	"github.com/liangyaopei/consistent"
)

func TestNew(t *testing.T) {
	nodes := []string{
		"185.199.110.153",
		"185.199.110.154",
		"185.199.110.155",
	}
	ring := consistent.New(nodes, consistent.DefaultHashFn)

	keys := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	nodeAdd := "185.199.110.156"
	nodeDel := "185.199.110.153"

	oriDis := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		oriDis[key] = node
	}

	ring.AddNodeWeight(nodeAdd, 1)
	addDes := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		addDes[key] = node
	}

	ring.DelNode(nodeDel)
	delDes := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		delDes[key] = node
	}

	t.Logf("adding node:%s,del node:%s", nodeAdd, nodeDel)

	for _, key := range keys {
		t.Logf("key:%15s,ori:%15s,add:%15s,del:%15s", key, oriDis[key], addDes[key], delDes[key])
	}

	for node, weight := range ring.GetNodeWeight() {
		t.Logf("node:%s,weight:%d", node, weight)
	}

	for node, keys := range ring.GetHashRing() {
		t.Logf("node:%s,keys:%v", node, keys)
	}
}

func TestNewWithWeight(t *testing.T) {
	nodes := map[string]int{
		"185.199.110.153": 3,
		"185.199.110.154": 2,
		"185.199.110.155": 1,
	}
	ring := consistent.NewWithWeight(nodes, consistent.DefaultHashFn)

	keys := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	nodeAdd := "185.199.110.156"
	nodeDel := "185.199.110.153"

	oriDis := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		oriDis[key] = node
	}

	ring.AddNodeWeight(nodeAdd, 1)
	addDes := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		addDes[key] = node
	}

	ring.DelNode(nodeDel)
	delDes := make(map[string]string)
	for _, key := range keys {
		node := ring.LocateKeyStr(key)
		delDes[key] = node
	}

	t.Logf("adding node:%s,del node:%s", nodeAdd, nodeDel)

	for _, key := range keys {
		t.Logf("key:%15s,ori:%15s,add:%15s,del:%15s", key, oriDis[key], addDes[key], delDes[key])
	}

	for node, weight := range ring.GetNodeWeight() {
		t.Logf("node:%s,weight:%d", node, weight)
	}

	for node, keys := range ring.GetHashRing() {
		t.Logf("node:%s,keys:%v", node, keys)
	}
}
