package main

import (
	"fmt"
	"reflect"

	"github.com/martin2250/taxi_software_go/pkg/fpga"
)

func iterateType(offset uintptr, typ reflect.Type) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "_" {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			fmt.Printf("----- %s\n", field.Name)
			iterateType(field.Offset, field.Type)
			fmt.Printf("-----\n")
		} else {
			fmt.Printf("[0x%04X] %14s %s size=%d, align=%d\n",
				field.Offset+offset, field.Name, field.Type, field.Type.Size(), field.Type.Align(),
			)
		}
	}
}

func main() {
	typ := reflect.TypeOf(fpga.Fpga{})
	iterateType(0, typ)
}
