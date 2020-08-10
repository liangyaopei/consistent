package consistent

import "hash/fnv"

// HashFn is the definition of hashing byte array to uin64
type HashFn func([]byte) uint64

// DefaultHashFn return the fnv hash function
// hashed twice to make the distribution more
// uniformly
func DefaultHashFn(data []byte) uint64 {
	fn := fnv.New64a()
	sum := fn.Sum(data)
	fn.Reset()
	fn.Write(sum)
	return fn.Sum64()
}
