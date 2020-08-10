consistent
==========
[![Go Report Card](https://goreportcard.com/badge/github.com/liangyaopei/consistent)](https://goreportcard.com/report/github.com/liangyaopei/consistent)
[![GoDoc](https://godoc.org/github.com/liangyaopei/consistent?status.svg)](http://godoc.org/github.com/liangyaopei/consistent)
[English Version](README.md)

概览
-----
这个package以哈希环(hash ring)的方式，实现了一致性哈希(consistent hash)算法。它是线程安全的，可以并发地使用。
用户可以增加带权重的节点，更新节点的权重和删除哈希环上的节点。
这个package不需要第三方的依赖。

安装
-----
```
go get -u -v github.com/liangyaopei/consistent
```

使用例子
------
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

哈希函数
------
为了使得哈希之后的key的分布更加均匀，`DefaultHashFn` 对key进行了两次hash
```go
func DefaultHashFn(data []byte) uint64 {
	fn := fnv.New64a()
	sum := fn.Sum(data)
	fn.Reset()
	fn.Write(sum)
	return fn.Sum64()
}
```