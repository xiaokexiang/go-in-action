package chapter9

import (
	"testing"
)

func IsPalindrome(s string) bool {
	for i := range s {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

// 功能测试函数，Test开头
func Test_Hello(t *testing.T) {
	if IsPalindrome("dadadada") {
		t.Error(`IsPalindrome("dadadada") is false`)
	}
}

// 基准测试行数，Benchmark开头
func Benchmark_Hello(b *testing.B) {
	if IsPalindrome("dadadada") {
		b.Error(`IsPalindrome("dadadada") is false`)
	}
}

func Example_hello() {

}
