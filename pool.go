package bencode

import (
	"sync"
)

const strSliceLen = 20

var strslicePool = sync.Pool{
	New: func() any {
		var j [strSliceLen]string
		return &j
	},
}

func getStrArray() *[strSliceLen]string {
	return strslicePool.Get().(*[strSliceLen]string)
}

func putStrArray(ss *[strSliceLen]string) {
	strslicePool.Put(ss)
}
