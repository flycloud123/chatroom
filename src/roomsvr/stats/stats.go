package stats

import "strings"

func preProcessMsg(msg []byte) []byte {
	var ret []byte
	for _, ch := range msg {
		if ch >= 'A' && ch <= 'Z' {
			ch += 32
		} else if ! (ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9') {
			ch = ' '
		}

		ret = append(ret, ch)
	}
	return ret
}

func FindPopularWord(messages []string) string {
	if messages == nil {
		return ""
	}
	word2count := make(map[string]int)
	for _, msg := range messages {
		msgpro := string(preProcessMsg([]byte(msg)))
		msgSplits := strings.Split(msgpro, " ")
		for _, one := range msgSplits {
			if one == "" {
				continue
			}
			if _, ok := word2count[one]; ok {
				word2count[one] = word2count[one] + 1
			} else {
				word2count[one] = 1
			}
		}
	}

	maxCount := 0
	retWord := ""
	for word, count := range word2count {
		if maxCount < count {
			maxCount = count
			retWord = word
		}
	}
	return retWord
}