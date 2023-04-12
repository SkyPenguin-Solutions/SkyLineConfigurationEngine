package Rem

import (
	"fmt"
)

var builtins = map[string]*ObjectBUILTINFUNCTION{
	"len": &ObjectBUILTINFUNCTION{Function: func(args ...SLC_Object) SLC_Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *ObjectArray:
			return &ObjectInteger{Value: int64(len(arg.Elements))}
		case *ObjectString:
			return &ObjectInteger{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].ObjectDataType())
		}
	},
	},
	"system": &ObjectBUILTINFUNCTION{
		Function: func(args ...SLC_Object) SLC_Object {
			for _, arg := range args {
				switch argt := arg.(type) {
				case *ObjectString:
					if argt.ObjectInspectFunc() == "errors" {
						return &ObjectString{Value: "errors"}
					}
				default:
					return &ObjectString{Value: "0x01"}
				}
			}
			return &ObjectString{Value: "0x01"}
		},
	},
	"puts": &ObjectBUILTINFUNCTION{
		Function: func(args ...SLC_Object) SLC_Object {
			for _, arg := range args {
				fmt.Println(arg.ObjectInspectFunc())
			}

			return NULL
		},
	},
}
