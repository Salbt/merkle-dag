package merkledag

import (
	"encoding/json"
	"fmt"
)

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据hash和path， 返回对应的文件, hash对应的类型是tree
	// 文件作为object存储在store中，先根据hash值找到对应的value
	var data []byte // 文件数据
	var tree Object
	h := hp.Get()
	h.Write(hash)
	value, err := store.Get(h.Sum(nil))
	if err != nil {
		_ = fmt.Errorf("get value from store failed, err: %s", err)
		return nil
	}

	// 得到tree的json,将其转化为tree类型
	err = json.Unmarshal(value, &tree)
	if err != nil {
		_ = fmt.Errorf("unmarshal value to tree failed, err: %s", err)
		return nil
	}

	// 遍历links中的link，link.name比较path文件获得对应的数据所在位置
	var start uint64 // 文件字节初始位置
	var link Link
	// 遍历tree对象，找到path相同的文件对象名称，根据文件位置和大小拿到文件对象
	for link = range tree.Links {
		if link.Name == path {
			data = tree.Data[start : start+link.Size]
		}
		start += link.Size
	}
	return data
}
