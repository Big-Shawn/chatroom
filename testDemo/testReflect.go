package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Myint int

func main() {

	sets := make(map[string]interface{})
	unit := make(map[string]interface{})
	unit["data"] = "Hello Map Type"
	sets["dd"] = 2
	sets["dd2"] = "hello"
	sets["dd3"] = unit

	var mi Myint = 10

	fmt.Println(reflect.TypeOf(mi).String())
	fmt.Println(reflect.TypeOf(mi).Kind().String())
	fmt.Println(reflect.ValueOf(mi).Kind().String())

	// typeOf get variable information
	s := reflect.TypeOf(sets)
	t := reflect.TypeOf(unit)
	fmt.Println("variable sets type is : " + s.String())
	fmt.Println("variable sets Kind is : " + s.Kind().String())
	fmt.Println("variable sets Key is : " + s.Key().String())
	fmt.Println("variable unit type is : " + t.String())

	// valueOf get variable information
	sV := reflect.ValueOf(sets)
	//s2 := sV.Interface().(string)
	fmt.Println(sV)
	fmt.Println("Variable sets value is : " + sV.String())
	fmt.Println("Variable sets Type is : " + sV.Type().String())
	fmt.Println("Variable sets Kind is : " + sV.Kind().String())

	mapRange := sV.MapRange()
	for mapRange.Next() {
		fmt.Printf("map key is %v, map value is %v \n", mapRange.Key(), mapRange.Value())
	}
	/**
	map key is dd, map value is 2
	map key is dd2, map value is hello
	map key is dd3, map value is map[data:Hello Map Type]
	*/

	// 使用 mapKeys 获取值得索引
	for k, v := range sV.MapKeys() {
		fmt.Printf("map key is %d, map value is %s \n", k, v)
	}

	// 使用golang 自带的数据类型进行处理
	// sV.kind 描述的是底层的基础数据类型，而不是自定义类型(the underlying type, not the static type.)
	// 如果希望返回自定义的数据类型应该是使用typeOf
	fmt.Println(sV.Kind() == reflect.Map)

	container := make(map[string]interface{})
	marshal, err := json.Marshal(sets)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))

	err = json.Unmarshal(marshal, &container)
	if err != nil {
		panic(err)
	}
	fmt.Println(container)

}
