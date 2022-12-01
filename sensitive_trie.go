package main

import "sync"

//SensitiveTrie 敏感词前缀树
type SensitiveTrie struct {
	mutex sync.RWMutex
	root  *TrieNode
}

// TrieNode 敏感词前缀树节点
type TrieNode struct {
	childMap map[rune]*TrieNode
	Data     string
	End      bool
}

// NewSensitiveTrie 构造敏感词前缀树实例
func NewSensitiveTrie() *SensitiveTrie {
	return &SensitiveTrie{
		root: &TrieNode{
			End: false,
		},
	}
}

// AddChild 前缀树添加子节点
func (tn *TrieNode) AddChild(c rune) *TrieNode {
	if tn.childMap == nil {
		tn.childMap = make(map[rune]*TrieNode)
	}

	if trieNode, ok := tn.childMap[c]; ok {
		// 存在不添加了
		return trieNode
	} else {
		// 不存在
		tn.childMap[c] = &TrieNode{
			childMap: nil,
			End:      false,
		}
		return tn.childMap[c]
	}
}

// AddSensitiveWord 添加敏感字
func (st *SensitiveTrie) AddSensitiveWord(sensitiveWord string) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	trieNode := st.root
	sensitiveChars := []rune(sensitiveWord)
	for _, char := range sensitiveChars {
		trieNode = trieNode.AddChild(char)
	}
	trieNode.End = true
	trieNode.Data = sensitiveWord
}

//AddSensitiveWords 批量添加敏感字
func (st *SensitiveTrie) AddSensitiveWords(sensitiveWords []string) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	for _, sensitiveWord := range sensitiveWords {
		trieNode := st.root
		sensitiveChars := []rune(sensitiveWord)
		for _, char := range sensitiveChars {
			trieNode = trieNode.AddChild(char)
		}
		trieNode.End = true
		trieNode.Data = sensitiveWord
	}
}

//RefreshSensitiveTrie 刷新敏感字段数，获取最新的敏感字段
func (st *SensitiveTrie) RefreshSensitiveTrie(sensitiveWords []string) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.root = &TrieNode{
		End: false,
	}
	for _, sensitiveWord := range sensitiveWords {
		trieNode := st.root
		sensitiveChars := []rune(sensitiveWord)
		for _, char := range sensitiveChars {
			trieNode = trieNode.AddChild(char)
		}
		trieNode.End = true
		trieNode.Data = sensitiveWord
	}
}

// FindChild 前缀树寻找字节点
func (tn *TrieNode) FindChild(c rune) *TrieNode {
	if tn.childMap == nil {
		return nil
	}
	if trieNode, ok := tn.childMap[c]; ok {
		return trieNode
	}
	return nil
}

// Match 匹配是否包含敏感字，并且返回涉及到的敏感字
func (st *SensitiveTrie) Match(text string) (sensitiveWords []string) {
	if st.root == nil {
		return nil
	}

	sensitiveMap := make(map[string]*struct{}) // 利用map把相同的敏感词去重
	textChars := []rune(text)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	for i, textLen := 0, len(textChars); i < textLen; i++ {
		trieNode := st.root.FindChild(textChars[i])
		if trieNode == nil {
			continue
		}
		j := i + 1
		for ; j < textLen && trieNode != nil; j++ {
			if trieNode.End {
				if _, ok := sensitiveMap[trieNode.Data]; !ok {
					sensitiveWords = append(sensitiveWords, trieNode.Data)
				}
				sensitiveMap[trieNode.Data] = nil
			}
			trieNode = trieNode.FindChild(textChars[j])
		}
		if j == textLen && trieNode != nil && trieNode.End {
			if _, ok := sensitiveMap[trieNode.Data]; !ok {
				sensitiveWords = append(sensitiveWords, trieNode.Data)
			}
			sensitiveMap[trieNode.Data] = nil
		}
	}
	return sensitiveWords
}

func (st *SensitiveTrie) SensitiveNums() int {
	return nums(st.root)
}

func nums(node *TrieNode) int {
	var num int
	if node.End {
		num++
	}
	for _, trieNode := range node.childMap {
		num += nums(trieNode)
	}
	return num
}
