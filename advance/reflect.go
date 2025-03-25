package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/*
变量结构: Pair<Type, Value>
Type:
	1. static type(int string ...)
	2. concrete type(interface指向的具体数据类型)
Value: 变量的值
ex: var a int = 0，即Pair<Type:int,Value:0>
*/

func baseReflect(arg any) {
	fmt.Println("type: ", reflect.TypeOf(arg))
	fmt.Println("value: ", reflect.ValueOf(arg))
}

/*
tag: `key:"value"` 类似java中的JsonFormat注解
*/

type User struct {
	Name string `json:"t_name"` // 不能忘记“”号
	Age  int    `json:"t_age"`
}

func (u *User) Print() {
	fmt.Printf("Name: %s\n", u.Name)
	fmt.Printf("Age: %d\n", u.Age)
}

/*
NumField()针对的是User，因为是struct中的field
NumMethod()针对的是*User，因为方法中定义的是*User
通过Elem()用于获取指针指向的元素类型
*/
func structReflect(arg any) {
	ot := reflect.TypeOf(arg)
	ov := reflect.ValueOf(arg)
	t := ot
	v := ov
	if t.Kind() == reflect.Pointer { // 解引用，获取引用的类型， 否则因为指针没有具体的字段，NumField()会报错
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("type: ", t)
	fmt.Println("value: ", v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		fmt.Printf("field: %v, tag: %v, value: %v\n", field.Name, field.Tag.Get("json"), value)
	}
	for i := 0; i < ot.NumMethod(); i++ { // 方法使用的是*User指针，所以这里的是ot即指针类型
		m := ot.Method(i)
		fmt.Printf("method: %v, value: %v\n", m.Name, m.Type)
	}
}

func JsonTag(arg User) {
	jsonStr, err := json.Marshal(arg)
	if err != nil {
		fmt.Println("json format error: ", err)
	}
	fmt.Printf("jsonStr: %s\n", jsonStr)
	u := &User{}
	if err := json.Unmarshal(jsonStr, u); err != nil {
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	fmt.Printf("json: %v\n", u)

}

func main() {
	//baseReflect(10)
	//fmt.Println("reflect struct ...")
	//structReflect(&User{
	//	Name: "zs",
	//	Age:  18,
	//})
	JsonTag(User{"ls", 28})
}
