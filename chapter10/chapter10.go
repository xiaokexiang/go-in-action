package chapter10

import (
	"fmt"
	"reflect"
	"strings"
)

type People interface {
	GetName() string
}

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"addr"`
}

// GetName 结构体方法 值接收者类型
func (p Person) GetName() string {
	return p.Name
}

// GetName2 结构体方法 指针接收者类型
func (p *Person) GetName2(in Person) string {
	return in.Name
}

// GetName3 函数
func GetName3(in Person) string {
	return in.Name
}

func Main() {
	testReflect()
}

// reflect.TypeOf() 获取变量的类型
// Elem(): 用于获取反射值的元素类型，用于指针的类型时，会返回指针指向的底层值的类型；用于非引用类型，会导致panic
func testReflect() {
	fmt.Println("-------------------------> reflect.TypeOf() <-------------------------")
	type1 := reflect.TypeOf(1)
	type2 := reflect.TypeOf("hello world")
	type3 := reflect.TypeOf(Person{})                        // 非指针类型
	type4 := reflect.TypeOf(&Person{})                       // 指针类型
	fmt.Printf("type of `%d` is %s\n", 1, type1)             // int
	fmt.Printf("type of `%s` is %s\n", "hello world", type2) // string
	fmt.Printf("type of `%s` is %s\n", "Person{}", type3)    // Person
	fmt.Printf("type of `%s` is %s\n", "&Person{}", type4)   // *Person
	fmt.Printf("获取指针指向的底层值类型 -> type4.Elem().equals(type3): %#v\n", type4.Elem() == type3)

	fmt.Println("-------------------------> 获取结构体的成员变量 <-------------------------")
	p := reflect.TypeOf(Person{Name: "jack", Age: 48, Address: "地球村"}) // 不能用指针
	numField := p.NumField()                                           // 获取变量数量
	for i := 0; i < numField; i++ {
		field := p.Field(i)
		fmt.Printf("字段名称: %s, 字段内存偏移量: %d, 是否为匿名成员: %t, 数据类型: %v, 包外是否可见: %t, 字段tag: %s\n",
			field.Name, field.Offset, field.Anonymous, field.Type, field.IsExported(), field.Tag.Get("json"))
	}
	// 通过字段名获取字段
	if name, ok := p.FieldByName("Name"); ok {
		fmt.Printf("Name is exported: %t\n", name.IsExported())
	}
	fmt.Println("-------------------------> 获取结构体的成员方法 <-------------------------")
	/*
		1. reflect.TypeOf()传入指针类型，那么返回了指针指向的类型的所有方法（包括值接收者方法和指针接收者方法）；如果传入值类型那么只会包含该类型的值接收者方法（设计如此）
		2. 方法名必须是exported，否则反射无法获取到
	*/
	p2 := reflect.TypeOf(&Person{})
	for i := 0; i < p2.NumMethod(); i++ {
		method2 := p2.Method(i)
		fmt.Printf("方法名: %s, 方法类型： %s, 包外是否可见: %t\n", method2.Name, method2.Type, method2.IsExported())
	}
	fmt.Println("-------------------------> 获取函数信息 <-------------------------")
	typeFunc := reflect.TypeOf(GetName3)
	var in, out []string
	for i := 0; i < typeFunc.NumIn(); i++ {
		in = append(in, typeFunc.In(i).String())
	}
	for i := 0; i < typeFunc.NumOut(); i++ {
		out = append(out, typeFunc.Out(i).String())
	}
	fmt.Printf("函数名称: %s, 函数入参: %s, 函数出参: %s\n", typeFunc.Name(), strings.Join(in, ","), strings.Join(out, ","))
	fmt.Println("-------------------------> 判断是否实现接口 <-------------------------")
	typePeople := reflect.TypeOf((*People)(nil)).Elem() // people是接口无法创建，需要通过nil强制转换
	p1 := reflect.TypeOf(Person{})
	fmt.Printf("People是否是接口: %t, Person是否实现接口People： %t\n", typePeople.Kind() == reflect.Interface, p1.Implements(typePeople))

	fmt.Printf("-------------------------> reflect.ValueOf() <-------------------------")
	value := reflect.ValueOf(1)
	fmt.Println(value)
	t := value.Type()                                                                                          // value转type
	fmt.Printf("value: %#v, type: %s, kind: %s, cast value: %#v", value, t, t.Kind(), value.Interface().(int)) // 强制转换数据

	person := &Person{Name: "aio", Address: "china", Age: 18}
	pv := reflect.ValueOf(person)       // 指针value
	pv2 := pv.Elem()                    // 指针转为非指针value
	pv2.Addr()                          // 非指针转为指针value
	p3 := pv.Interface().(*Person).Name // 获取初始值
	fmt.Println(p3)
	var p4 interface{}
	v := reflect.ValueOf(p4) // 没有指向具体的值
	var p5 *Person = nil
	v1 := reflect.ValueOf(p5)
	fmt.Printf("v invaild? %t, v1 is nil? %t\n", v.Kind() == reflect.Invalid, v1.IsNil())

	p6 := &Person{
		Name:    "jack",
		Age:     18,
		Address: "china",
	}
	v6 := reflect.ValueOf(p6)                         // 必须传递指针才能修改原来的值
	if ok := v6.Elem().FieldByName("").CanSet(); ok { // CanSet()用来判断成员是否导出，只有可以导出的成员变量才能被修改
		v6.Elem().FieldByName("Name").SetString("lucy") // 必须将指针转为非指针后才能修改值
	}
	fmt.Printf("new value: %#v\n", v6)
	v7 := reflect.ValueOf(make([]int, 1, 5))
	if v7.Len() > 0 {
		v7.Index(0).SetInt(10)
	}
	fmt.Printf("v7 value: %#v\n", v7)

	fmt.Printf("-------------------------> reflect.ValueOf()调用函数 <-------------------------\n")
	valueFunc := reflect.ValueOf(GetName3)
	slice1 := make([]reflect.Value, 1)
	slice1[0] = reflect.ValueOf(Person{Name: "zhangsan", Age: 18, Address: "usa"})
	resultValue := valueFunc.Call(slice1)
	fmt.Printf("resultValue: %#v\n", resultValue[0].String())

	fmt.Printf("-------------------------> reflect.ValueOf()调用成员方法 <-------------------------\n")
	v8 := reflect.ValueOf(&Person{
		Name:    "jack",
		Age:     18,
		Address: "China",
	})
	method := v8.MethodByName("GetName")
	resultValue1 := method.Call([]reflect.Value{})
	fmt.Printf("resultValue: %#v\n", resultValue1[0].String())
}
