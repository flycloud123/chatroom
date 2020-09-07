package stats

import "testing"

func Test_FindPopularWord(t *testing.T) {
	var messages []string
	messages = append(messages, "hello    world")
	messages = append(messages, "hello world, @!# 123 *( we are ok")
	messages = append(messages, "hello worldChampion, we are ok")
	messages = append(messages, "hello worldChampion, we are ok")

	ret := FindPopularWord(messages)
	if ret != "hello" {
		t.Errorf("not ok")
	}
}
