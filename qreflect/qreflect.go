package qreflect

import "reflect"

//获取所有json tag的 name
func TagNameJsonNames(structs interface{}) []string  {
	t :=reflect.TypeOf(structs)
	tagNames := []string{}
	for i:=0;i<t.NumField();i++ {
		tagName :=t.Field(i).Tag.Get("json")
		if tagName!="" {
			tagNames = append(tagNames,tagName)
		}
	}

	return tagNames
}
