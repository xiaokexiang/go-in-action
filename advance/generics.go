package main

import "fmt"

/*
泛型
*/

type MyMap[K int | string, V int | string] map[K]V

/*
泛型的方法
*/

func (m MyMap[K, V]) Print() {
	for k, v := range m {
		fmt.Printf("key: %s, value: %s\n", k, v)
	}
}

type MyStruct[T int | string] struct {
	Name string
	Data T
}

type IPrintData[T int | float32 | string] interface {
	Print(data T)
}

type MyChan[T int | string] chan T

func main() {
	var m MyMap[string, string] = map[string]string{
		"a": "b",
	}
	m.Print()
}
