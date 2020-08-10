package _examples

import (
	"testing"

	"github.com/liangyaopei/consistent"
)

func BenchmarkLocateKey(b *testing.B) {
	nodes := []string{
		"185.199.110.153",
		"185.199.110.154",
		"185.199.110.155",
	}
	ring := consistent.New(nodes, consistent.DefaultHashFn)
	keys := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			ring.LocateKeyStr(key)
		}
	}
}
