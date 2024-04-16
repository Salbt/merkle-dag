package merkledag

import (
	"encoding/json"
	"hash"
	"log"
)

const (
	BLOB       = "blob"
	LIST       = "list"
	TREE       = "tree"
	BLOCK_SIZE = 256 * 1024
)

type Link struct {
	Name string // name or alias of this link
	Hash []byte // cryptographic hash of target
	Size uint64 // total size of target
}

type Object struct {
	Links []Link // array of links
	Data  []byte // opaque content data
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	var MerkleRoot []byte
	// 检查node类型
	if node.Type() == FILE {
		MerkleRoot = storeFile(store, node.(File), h)
	}
	if node.Type() == DIR {
		MerkleRoot = storeDir(store, node.(Dir), h)
	}
	return MerkleRoot
}

// 当存储的node是文件时，dag的存储逻辑
func storeFile(store KVStore, file File, h hash.Hash) []byte {
	var key []byte
	// 检查文件大小, 文件大小小于256KB直接以blob形式存储在kv中
	// 文件大小大于256KB，以list存储在kv中，key为blob的hash.sum, vaule为list的json
	if file.Size() <= BLOCK_SIZE {
		// 构建blob对象
		blob := Object{
			Data: file.Bytes(),
		}
		key = PutStore(store, blob, h)

	} else {
		// 构建list对象
		list := Object{
			Links: make([]Link, 0),
			Data:  make([]byte, 0),
		}
		for start := uint64(0); start <= file.Size(); start += BLOCK_SIZE {
			// 构建blob的link
			var link Link

			// blob的大小
			size := uint64(BLOCK_SIZE)
			if file.Size()-start < BLOCK_SIZE {
				size = file.Size() - start
			}

			// blob数据
			blob := Object{
				Data: file.Bytes()[start : start+size-1],
			}

			_hash := PutStore(store, blob, h)

			link = Link{
				Size: size,
				Hash: _hash,
			}

			// 将link和data装入list
			list.Links = append(list.Links, link)
			list.Data = append(list.Data, BLOB...)

			// 计算list hash
			key = h.Sum(key)
		}

		PutStore(store, list, h)

	}
	return key
}

// 当存储文件类型为list, dag的存储逻辑
func storeDir(store KVStore, dir Dir, h hash.Hash) []byte {
	var key []byte
	// 使用非递归方式实现目录的遍历
	var dirStack []Dir     // 遍历节点使用的栈
	var treeStack []Object // 构建tree对象使用的栈

	dirStack = append(dirStack, dir) // 将根目录压入栈
	treeStack = append(treeStack, Object{
		Links: make([]Link, 0),
		Data:  make([]byte, 0),
	}) // 将根目录对象压入栈

	var it DirIterator
	var top int
	for len(dirStack) != 0 {

		top = len(dirStack) - 1 // 栈顶指针
		it = dirStack[top].It() // 获取迭代器
		node := it.Node()       // 用迭代器获得文件
		// 判断node类型
		// 如果是FILE就将其构建成link加入当前目录
		if node.Type() == FILE {
			file := node.(File)

			// 构建blob/list对象
			data := BLOB
			if file.Size() > BLOCK_SIZE {
				data = LIST
			}
			key = storeFile(store, file, h)

			treeStack[top].Links = append(treeStack[top].Links, Link{
				Name: file.Name(),
				Hash: key,
				Size: file.Size(),
			})
			treeStack[top].Data = append(treeStack[top].Data, data...)

			// 当前目录下没有文件就将目录弹出栈
			for !it.Next() && len(dirStack) != 0 {
				// 将构建完成的tree json化存储入kv中
				_dir := dirStack[top]
				key = PutStore(store, treeStack[top], h)
				// pop stack
				dirStack = dirStack[:top]
				treeStack = treeStack[:top]

				if top = len(dirStack) - 1; top < 0 {
					return key
				} // 更新top, 如果top < 0 ,说明栈中只有根目录

				// 将构建好的tree添加到其父tree对象中
				treeStack[top].Data = append(treeStack[top].Data, TREE...)
				treeStack[top].Links = append(treeStack[top].Links, Link{
					Name: _dir.Name(),
					Hash: key,
					Size: _dir.Size(),
				})
				it = dirStack[top].It() // 更新迭代器
			}
		}
		if node.Type() == DIR {
			// 将文件压栈，构建一个新tree并压栈

			dirStack = append(dirStack, node.(Dir))
			treeStack = append(treeStack, Object{
				Links: make([]Link, 0),
				Data:  make([]byte, 0),
			})
		}

	}

	return key
}

func PutStore(store KVStore, data Object, h hash.Hash) (key []byte) {
	value, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Vaule json Marshal failed")
	}
	h.Reset()
	h.Write(value)
	key = h.Sum(nil)
	if err = store.Put(key, value); err != nil {
		log.Fatal("Put kv failed")
	}
	return
}

func GetStore(store KVStore, hash []byte) (tree *Object) {
	//	使用hash获取，文件json对象
	value, err := store.Get(hash)
	if err != nil {
		log.Fatal("get kv failed")
		return nil
	}

	err = json.Unmarshal(value, &tree)
	if err != nil {
		log.Fatal("tree json failed")
	}

	return
}
