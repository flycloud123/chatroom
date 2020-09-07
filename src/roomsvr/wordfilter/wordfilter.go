package wordfilter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var trie *Trie


type Trie struct {
	next map[rune]*Trie
	isWord bool
}

func CreateTire() *Trie {
	root := new(Trie)
	root.next = make(map[rune]*Trie)
	root.isWord = false
	return root
}

func (self *Trie) Insert(word string) {
	if len(word) == 0 {
		return
	}

	for _, v := range word {
		if v >= 'A' && v <= 'Z' {
			v += 32
		}

		if self.next[v] == nil {
			node := new(Trie)
			node.next = make(map[rune]*Trie)
			node.isWord = false
			self.next[v] = node
		}
		self = self.next[v]
	}

	self.isWord = true
}

func (self *Trie) preorder() {
	if self.isWord {
		return
	}
}

func (self *Trie) Search(word string) bool {
	for _, v := range word {
		if self.next[v] == nil {
			return false
		}
		self = self.next[v]
	}
	return self.isWord
}

func (self *Trie) SearchIn(word string) (bool, int) {
	count := 0
	for _, v := range word {
		if v >= 'A' && v <= 'Z' {
			v += 32
		}
		if self.next[v] == nil {
			return false, 0
		}
		count++
		self = self.next[v]
		if self.isWord {
			return true, count
		}
	}
	return false, 0
}

func LoadDirtyFromFile(path string) {
	trie = CreateTire()

	fi, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		trie.Insert(string(line))
	}
}

func ReplaceDirty(word []byte) []byte {
	length := len(word)
	for i := 0; i < length;  {
		if isin, count := trie.SearchIn(string(word[i:])); isin {
			for j := i; j < i + count; j++ {
				word[j] = '*'
			}
			i += count
		}
		i++
	}
	return word
}


