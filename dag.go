package merkledag

import (
	"encoding/json"
	"fmt"
	"hash"
)

const (
	BLOCK_SIZE = 256 * 1024
)

type Link struct {
	Name string // name or alias of this link
	Hash []byte // cryptographic hash of target
	Size uint64 // total size of target
}

type Object struct {
	Links []Link // array of links
	Data []byte // opaque content data
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	var tree []byte
	var err error
	if node.Type() == FILE {
		file := node.(File)
		// 返回hash
		tree , err = StoreFile(store,file, h)
		if err != nil {
			_ = fmt.Errorf("Error storing file: %s", err)
		}
	}

	if node.Type() == DIR {
		dir := node.(Dir)
		tree ,err = StoreDir(store,dir, h)
		if err != nil {
			_ = fmt.Errorf("Error storing directory: %s", err)
		}
	}
	h.Reset()
	h.Write(tree)
	return h.Sum(nil)
}

// StoreFile 将文件存储在KV中
func StoreFile(store KVStore, file File, h hash.Hash) ([]byte,error) {
	h.Reset()
	h.Write(file.Bytes())
	var tree = Object{
		Links: []Link{
            Link{
                Name: file.Name(),
                Hash: h.Sum(nil),
                Size: file.Size(),
            },
        },
		Data: make([]byte, 0),
	}

	// 判断File的对象类型
	if file.Size() > BLOCK_SIZE {
		// 对其进行分块
		for i := uint64(0); i < file.Size(); i += BLOCK_SIZE{
			// 创建blob对象
			blob := Object{
				Data: file.Bytes()[i : i+BLOCK_SIZE-1],
			}
			data, _ := json.Marshal(blob)
			tree.Data = append(tree.Data, data...)
		}
	} else {
		tree.Data = append(tree.Data, file.Bytes()...)
	}
	TreeJson,_ := json.Marshal(tree)
	_ = store.Put(h.Sum(nil), TreeJson)
	return TreeJson,nil
}
// StoreDir 将目录存储在KV中
func StoreDir(store KVStore, dir Dir, h hash.Hash) ([]byte, error) {
	// 定义栈结构
	var stack = make([]Node, 0)
	// 将第一个node压入栈
	stack = append(stack, dir)
	// 需要构建的树-文件对象
	var tree Object
	it := dir.It()
	// 栈不为空或者节点还有孩子时就继续遍历
	for len(stack) != 0 {
		// 拿到子节点，并将子节点入栈
		node := it.Node()
		// 如果node是文件类型就将文件存入tree中
		if node.Type() == FILE {
			 var object Object
			 json0Object,_:=StoreFile(store, node.(File), h)
			 _ = json.Unmarshal(json0Object, &object)
			 tree.Links = append(tree.Links, object.Links...)
			 tree.Data = append(tree.Data, json0Object...)
		}
		if node.Type() == DIR {
			// 更新迭代器
			it = node.(Dir).It()
			// 将该node压入栈中
			stack = append(stack, node.(Dir))
		}
		// 迭代node,如果没有下一个node就弹栈
		if !it.Next() {
			stack = stack[:len(stack)-1]
		}
	}

	TreeJson,_ := json.Marshal(tree)
	return TreeJson, nil

}
