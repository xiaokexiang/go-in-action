package chapter2

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func Main() {
	fmt.Println("=============> bit operation <=============")
	bit()
	fmt.Println("=============> float operation <=============")
	float()
	fmt.Println("=============> string operation <=============")
	str()
	fmt.Println("=============> const operation <=============")
	constant()
}

func bit() {
	/*
			 1<<1: 00000001 -> 00000010
			 1<<5: 00000001 -> 00100000
			 或运算：有1为1否则为0
			  00000010
			| 00100000
			————————————
		      00100010
	*/
	var x uint8 = 1<<1 | 1<<5
	fmt.Printf("%08b\n", x) // 00100010
	fmt.Printf("%d\n", x)   // 34
}

func float() {
	var f32 float32 = math.MaxFloat32
	var f64 float64 = math.MaxFloat64

	fmt.Printf("f32: %8.1f\n", f32) // 8表示输出的值不足8个字符，则会用空格填充 1表示保留小数点后1位
	fmt.Println(f64)
}

func bool() {
	a := 0
	if !(a == 0 && a == 1) {
		fmt.Println(a)
	}
}

func str() {
	s := "你好世界"         // 默认按照utf8进行编码
	fmt.Println(len(s)) // 输出字符数量

	// 获取utf8下字符的编码值
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:]) // size表示按照utf8的编码占用的字节数
		fmt.Printf("Character: %c, Code: %d, size: %d\n", r, r, size)
		i += size
	}

	/*
		获取utf8下字符的编码值 字符 = n个字节（n>0）
		Character: 你, Code: [228 189 160]
		Character: 好, Code: [229 165 189]
		Character: 世, Code: [228 184 150]
		Character: 界, Code: [231 149 140]
		如果只包含ascii字符，那么直接通过s[i]来获取
	*/
	for _, char := range s {
		fmt.Printf("Character: %c, bytes: %v\n", char, []byte(string(char)))

	}

	fmt.Println(s[0], s[7])     // 通过下标访问，返回的是uint8类型的字节编码，即228和184
	fmt.Println(s[0:6])         // 默认会输出字符
	fmt.Println([]byte(s[0:5])) // 输出rune字节数组

	// 字符串的不可变性，a和b公用了底层字节数组
	a := "a"
	b := a
	a += "b"
	fmt.Println(a)
	fmt.Println(b)

	x := "abc"
	y := []byte(x)
	z := string(y)
	fmt.Printf("x: %s, y: %v, z: %s\n", x, y, z)

}

func constant() {
	const hello string = "hello" // 定义类型
	const (
		a int = iota // 从0依次曾增加1
		b
		c
		d
	)
	fmt.Println(a, b, c, d) // 0 1 2 3

	const t = 10                                // 无类型常量，会根据上下文确认转为不同的类型
	fmt.Printf("type: %T, v: %d\n", t+1, t+1)   // 此时t是int类型
	fmt.Printf("type: %T, v: %f", t+0.1, t+0.1) // 此时t是float64类型
}
