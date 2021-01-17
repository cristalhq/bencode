package bencode

import "sync"

var strslicePool = sync.Pool{
	New: func() interface{} {
		var j [20]string
		return &j
	},
}

func getStrArray() *[20]string {
	return strslicePool.Get().(*[20]string)
}

func putStrArray(ss *[20]string) {
	strslicePool.Put(ss)
}
