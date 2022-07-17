package btree

import "fmt"

/**
* B树的特性
* 每一个节点最多有 m 个子节点
* 每一个非叶子节点（除根节点）最少有 ⌈m/2⌉ 个子节点
* 如果根节点不是叶子节点，那么它至少有两个子节点
* 有 k 个子节点的非叶子节点拥有 k − 1 个键
* 所有的叶子节点都在同一层
**/

/**
* 3阶B树举例
	     k1     |      k2
		/       |       \
	  l1|l12    l3|l32   l4|l42

**/

// 非并发安全
type BTree struct {
	m    int // 阶数
	root *BTreeNode
}

type BTreeNode struct {
	isLeaf   bool  // 是否为叶子结点
	keyNum   int   // key个数
	keyList  []int // key 有序
	leafList []*BTreeNode
}

func createNode(m int, keyNum int, isLeaf bool) *BTreeNode {
	return &BTreeNode{
		isLeaf:   isLeaf,
		keyNum:   keyNum,
		keyList:  make([]int, m),          // 这里多创建了一个
		leafList: make([]*BTreeNode, m+1), // 这里多创建了一个
	}
}

/**
* 查找key, 返回key所在结点及下标
* 采用递归的方法
**/
func (bTreeNode *BTreeNode) search(key int) (*BTreeNode, int) {
	if bTreeNode.isLeaf || bTreeNode.keyNum == 0 {
		return nil, -1
	}

	// 找到第一个不比key小的，注意leaf的数量比key数量多1
	idx := 0
	for idx < bTreeNode.keyNum && key > bTreeNode.keyList[idx] { // 这里可以采用二分查找提高效率
		idx++
	}

	// 在当前节点找到了key
	if idx < bTreeNode.keyNum && bTreeNode.keyList[idx] == key {
		return bTreeNode, idx
	}

	// 是叶子结点 则找不到了
	if bTreeNode.leafList[idx].isLeaf {
		return nil, -1
	}

	return bTreeNode.leafList[idx].search(key)
}

/**
* 插入key
* 时间复杂度为O(logn)
**/
func (bTreeNode *BTreeNode) intert(key int, m int) ([]*BTreeNode, int) {
	// 找到第一个不比key小的，注意leaf的数量比key数量多1
	idx := 0
	for idx < bTreeNode.keyNum && key > bTreeNode.keyList[idx] {
		idx++
	}

	// BTreeNode已有该key
	if idx < bTreeNode.keyNum && bTreeNode.keyList[idx] == key {
		return nil, 0
	}

	if bTreeNode.isLeaf { // 叶子结点
		// idx及idx后面元素往后移动
		if idx < m-1 {
			copy(bTreeNode.keyList[idx+1:], bTreeNode.keyList[idx:])
		}
		bTreeNode.keyList[idx] = key
		bTreeNode.keyNum += 1
	} else { // 非叶子结点
		// 判断结点是否要满了
		// 先往上提
		children, midKey := bTreeNode.leafList[idx].intert(key, m)
		if children != nil {
			// idx是midKey的插入位置,idx后面元素往后移动
			if idx < m-1 {
				copy(bTreeNode.keyList[idx+1:], bTreeNode.keyList[idx:])
				copy(bTreeNode.leafList[idx+2:], bTreeNode.leafList[idx+1:])
			}
			bTreeNode.keyList[idx] = midKey
			bTreeNode.leafList[idx] = children[0]
			bTreeNode.leafList[idx+1] = children[1]
			bTreeNode.keyNum += 1
		}
	}

	if bTreeNode.keyNum == m { // 节点key数量超了
		// 从中间将该结点一分为2
		mid := bTreeNode.keyNum / 2
		midKey := bTreeNode.keyList[mid]

		// 左儿子
		leftNode := createNode(m, mid, bTreeNode.isLeaf)
		copy(leftNode.keyList, bTreeNode.keyList[:mid])
		copy(leftNode.leafList, bTreeNode.leafList[:mid+1])

		// 右边儿子
		rightNode := createNode(m, m-mid-1, bTreeNode.isLeaf)
		copy(rightNode.keyList, bTreeNode.keyList[mid+1:])
		copy(rightNode.leafList, bTreeNode.leafList[mid+1:])

		return []*BTreeNode{leftNode, rightNode}, midKey
	}
	return nil, 0
}

/**
* 删除结点内的某个key
**/
func (bTreeNode *BTreeNode) removeKey(idx int) {
	if bTreeNode.isLeaf { // 叶子结点
		copy(bTreeNode.keyList[idx:], bTreeNode.keyList[idx+1:])
		bTreeNode.keyNum -= 1
	} else {
		// 在叶子结点找一个元素替换掉被删除结点
		leftChild := bTreeNode.leafList[idx]
		rightChild := bTreeNode.leafList[idx+1]
		var tmp *BTreeNode
		if leftChild != nil {
			// 从左儿子结点中找到最大的,替换掉被删除的key
			tmp = leftChild
			for !tmp.isLeaf {
				tmp = tmp.leafList[bTreeNode.keyNum]
			}
			bTreeNode.keyList[idx] = tmp.keyList[tmp.keyNum-1]
			tmp.keyNum -= 1
		} else if rightChild != nil {
			// 从右儿子结点找到最小的,替换掉被删除的key
			tmp = rightChild
			for !tmp.isLeaf {
				tmp = tmp.leafList[0]
			}
			bTreeNode.keyList[idx] = tmp.keyList[0]
			tmp.keyNum -= 1
		} else {
			fmt.Println("wrong!!!!")
		}
	}
}

/**
* 左旋
**/
func (bTreeNode *BTreeNode) leftShift(idx int) {
	pos := idx
	if idx == 0 {
		pos = 1
	}
	curNode := bTreeNode.leafList[idx]
	rightBro := bTreeNode.leafList[pos]
	// 将父节点的分隔值复制到缺少元素节点的最后,（分隔值被移下来；缺少元素的节点现在有最小数量的元素）
	curNode.keyList[bTreeNode.keyNum] = bTreeNode.keyList[idx]
	curNode.keyNum += 1
	// 兄弟结点的左儿子 作为 缺少元素节点 的右儿子
	curNode.leafList[bTreeNode.keyNum] = rightBro.leafList[0]

	// 将父节点的分隔值替换为右兄弟的第一个元素
	curNode.keyList[idx] = rightBro.keyList[0]

	// 右兄弟结点少了个元素
	copy(rightBro.keyList, rightBro.keyList[1:])
	copy(rightBro.leafList, rightBro.leafList[1:])
	rightBro.keyNum -= 1
}

/**
* 右旋
**/
func (bTreeNode *BTreeNode) rightShift(idx int) {
	pos := idx
	if pos == 0 {
		pos = 1
	}
	curNode := bTreeNode.leafList[idx]
	leftBro := bTreeNode.leafList[pos-1]
	// 当前结点数据往后一个位置
	copy(bTreeNode.leafList[1:], bTreeNode.leafList)
	copy(bTreeNode.keyList[1:], bTreeNode.keyList)
	// 将父节点的分隔值复制到缺少元素节点的第一个节点
	curNode.keyList[0] = bTreeNode.keyList[idx]
	curNode.keyNum += 1
	curNode.leafList[0] = leftBro.leafList[leftBro.keyNum] // 左兄弟的最后一个儿子放到当前结点的第一个儿子的位置

	curNode.keyList[idx] = leftBro.keyList[leftBro.keyNum-1] // 将父节点的分隔值替换为左兄弟的最后一个元素
	leftBro.keyNum -= 1
}

/**
* 合并
* 将分隔值及左右儿子结点合并
**/
func (bTreeNode *BTreeNode) merge(curNode *BTreeNode, idx int) {
	pos := idx
	if pos == 0 {
		pos = 1
	}
	leftBro := bTreeNode.leafList[pos-1]
	rightBro := bTreeNode.leafList[pos]
	// 将分隔值复制到左边的节点（左边的节点可以是缺少元素的节点或者拥有最小数量元素的兄弟节点）
	leftBro.keyList[leftBro.keyNum] = bTreeNode.keyList[pos-1]
	leftBro.keyNum += 1
	// 将分隔值的右边节点中所有的元素移动到左边节点（左边节点现在拥有最大数量的元素，右边节点为空）
	if rightBro != nil {
		copy(leftBro.leafList[curNode.keyNum:], rightBro.leafList)
		copy(leftBro.keyList[curNode.keyNum:], rightBro.keyList)
		leftBro.keyNum += rightBro.keyNum
	}

	// 将父节点中的分隔值和空的右子树移除（父节点失去了一个元素）
	copy(bTreeNode.keyList[pos-1:], bTreeNode.keyList[pos:])
	copy(bTreeNode.leafList[pos:], bTreeNode.leafList[pos+1:])
	bTreeNode.keyNum -= 1
}

/**
* 删除
**/
func (bTreeNode *BTreeNode) delete(key int, m int) {
	// 找到第一个不比key小的，注意leaf的数量比key数量多1
	idx := 0
	for idx < bTreeNode.keyNum && key > bTreeNode.keyList[idx] { // 这里可以采用二分查找提高效率
		idx++
	}
	curNode := bTreeNode.leafList[idx]

	if idx < bTreeNode.keyNum && bTreeNode.keyList[idx] == key {
		// 在当前节点找到了key
		bTreeNode.removeKey(idx)
	} else if curNode != nil {
		curNode.delete(key, m)
	}

	/**
	* 每一个非叶子节点（除根节点）最少有 ⌈m/2⌉ 个子节点
	* 有 k 个子节点的非叶子节点拥有 k − 1 个键
	* 因此非叶子结点最少有m/2个键
	**/
	// 叶子节点拥有的元素数量小于最低要求
	if curNode != nil && curNode.isLeaf && curNode.keyNum < m/2 { // 需要调整
		pos := idx
		if pos == 0 {
			pos = 1
		}
		leftBro := bTreeNode.leafList[pos-1]
		rightBro := bTreeNode.leafList[pos]
		if rightBro != nil && rightBro.keyNum > m/2 { // 父结点key并放到自己的最后 右兄弟的第一个key放到父结点
			// 如果缺少元素节点的右兄弟存在且拥有多余的元素，那么向左旋转
			bTreeNode.leftShift(idx)
		} else if leftBro != nil && leftBro.keyNum > m/2 { // 父结点key并放到自己的开头 左兄弟的最后一个key放到父结点
			// 如果缺少元素节点的左兄弟存在且拥有多余的元素，那么向右旋转
			bTreeNode.rightShift(idx)
		} else {
			// 如果它的两个直接兄弟节点都只有最小数量的元素，那么将它与一个直接兄弟节点以及父节点中它们的分隔值合并
			bTreeNode.merge(curNode, idx)
		}
	}
}

/**
* 创建m阶树
* 阶代表的最大子节点数目
* 每个结点key的最大数量为m-1
**/
func Create(m int) (*BTree, error) {
	if m < 2 {
		return nil, fmt.Errorf("m should be bigger than 2")
	}

	bTree := &BTree{
		m:    m,
		root: createNode(m, 0, true),
	}
	return bTree, nil
}

/** 查找key, 返回key所在结点及下标
* 时间复杂度为O(logn)
**/
func (bTree *BTree) Search(key int) (*BTreeNode, int) {
	return bTree.root.search(key)
}

/** 插入key
* 时间复杂度O(logn)
**/
func (bTree *BTree) Insert(key int) {
	children, midKey := bTree.root.intert(key, bTree.m)
	if children != nil {
		root := createNode(bTree.m, 1, false)
		root.keyList[0] = midKey
		root.leafList[0] = children[0]
		root.leafList[1] = children[1]
		bTree.root = root
	}
}

/**
* 删除key
* 时间复杂度O(logn)
**/
func (bTree *BTree) Delete(key int) {
	bTree.root.delete(key, bTree.m)
	// 父节点的元素数量小于最小值，重新平衡父节点
	if bTree.root.keyNum == 0 && bTree.root.leafList[0] != nil {
		bTree.root = bTree.root.leafList[0]
	}
}

// 层次遍历辅助结构体
type traversalHelper struct {
	father *BTreeNode
	idx    int
	child  *BTreeNode
}

/**
* 层次遍历
* 以[分隔值(父结点key)｜key数量] key... 格式打印
* 如[4|2] 5 6 表示当前结点的分隔值为4, key数量为2, 分别是5和6
**/
func (bTreeNode *BTreeNode) traversalLevel() {
	que := []*traversalHelper{}
	que = append(que, &traversalHelper{
		father: nil,
		child:  bTreeNode,
	})
	for len(que) > 0 {
		tmp := []*traversalHelper{}
		for _, helper := range que {
			if helper == nil || helper.child == nil {
				fmt.Println("kkkkkkkk ", len(que), helper == nil)
			}
			ftr := "root"
			if helper.father != nil {
				idx := helper.idx
				if idx == (helper.father.keyNum) {
					idx -= 1
				}
				ftr = fmt.Sprintf("%v", helper.father.keyList[idx])
			}
			fmt.Printf("[%v|%v] ", ftr, helper.child.keyNum)
			for i := 0; i <= helper.child.keyNum; i++ {
				// 注意leaf的数量比key数量多1
				if i < helper.child.keyNum {
					fmt.Print(helper.child.keyList[i], " ")
				}
				if !helper.child.isLeaf {
					tmp = append(tmp,
						&traversalHelper{
							father: helper.child,
							idx:    i,
							child:  helper.child.leafList[i],
						})
				}
			}
			fmt.Print("   ")
		}
		que = tmp
		fmt.Println()
	}
}

/**
* 层次遍历
* 以[分隔值(父结点key)｜key数量] key... 格式打印
* 如[4|2] 5 6 表示当前结点的分隔值为4, key数量为2, 分别是5和6
**/
func (bTree *BTree) TraversalLevel() {
	fmt.Printf("enter... m:%v mean max childer cnt is %v\n", bTree.m, bTree.m-1)
	bTree.root.traversalLevel()
}
