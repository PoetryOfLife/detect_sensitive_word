package DFA

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSensitiveTrie(t *testing.T) {

	file, err := os.ReadFile("./../min_word.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
	sensitiveWords := strings.Split(string(file), "，")

	st := NewSensitiveTrie()
	st.AddSensitiveWords(sensitiveWords)

	testTxt, err := os.ReadFile("./../test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}

	start := time.Now()
	fmt.Println(st.Match(string(testTxt)))
	fmt.Println(time.Since(start))
}
