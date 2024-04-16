package merkledag

import (
	"encoding/json"
	"log"
	"strings"
)

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 	根据hash和path， 返回对应的文件, hash对应的类型是tree
	file := GetStore(store, hash)
	// 	处理路径，将路径顺序存储在数组中
	paths := strings.Split(path, "/")
	//  根据路径返回文件json
	for _, fileName := range paths {
		for _, link := range file.Links {
			if link.Name == fileName {
				file = GetStore(store, link.Hash)
				break
			} else {
				file = nil
			}
		}
		if file == nil {
			log.Fatal("the path does not exist")
		}
	}
	if len(file.Links) == 0 {
		return file.Data
	} else {
		// 返回list对象，而非list数据
		data, err := json.Marshal(file)
		if err != nil {
			log.Fatal("json failed")
		}
		return data
	}
}
