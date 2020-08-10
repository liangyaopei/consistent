consistent
==========
[![Go Report Card](https://goreportcard.com/badge/github.com/liangyaopei/consistent)](https://goreportcard.com/report/github.com/liangyaopei/consistent)
[![GoDoc](https://godoc.org/github.com/liangyaopei/consistent?status.svg)](http://godoc.org/github.com/liangyaopei/consistent)
[中文版](README_zh.md)

Overview
--------
This package implements consistent hash algorithm using hash ring, it is thread-safe, and can be used
concurrently. 
You can add node with weight, update node's weight and delete node in hash ring.
This package doesn't require third-party dependency.

Install
-------
```
go get -u -v github.com/liangyaopei/consistent
```

Example
-------
```go
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
```

Hash Function
-------
To make hashed key distribute more uniformly, `DefaultHashFn` hash the key twice.
```go
func DefaultHashFn(data []byte) uint64 {
	fn := fnv.New64a()
	sum := fn.Sum(data)
	fn.Reset()
	fn.Write(sum)
	return fn.Sum64()
}
```