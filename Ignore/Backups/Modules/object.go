package Rem

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//
//  													 _____ _       __    _
// 														|   __| |_ _ _|  |  |_|___ ___
// 													    |__   | '_| | |  |__| |   | -_|
// 														|_____|_,_|_  |_____|_|_|_|___|
//		   														  |___|
//
//
// The SkyLine configuration language is a language and engine designed to act as a modification language to the SkyLine programming language. This language is
// very minimal and contains a regex base lexer, a very basic parser, a few keywords, a base interpreter and that is all as well as some backend engine code. This
// language is purely modified to be an extension to the SkyLine programming language, something that can be a pre processor language post processing for the main
// SkyLine script. Below is more technical information for the language
//
// Lexer       : Regex based lexer with minimal tokens and base syntax
// Parser      : Base parser with minimal tokens and base syntax with simple error systems
// REPL        : Does not exist
// Environment : Extremely minimal
// Types       : String, Boolean, Integer
// Statements  : set, import, use, errors, output, system, constant/const
//
//
// File contains -> This file contains all necessary functions for the objects within SLC allowing us to view the objects type

func (ObjInteger *ObjectInteger) ObjectDataType() ObjectsDataType  { return OBJECT_INTEGER } // Return integer object
func (ObjBoolean *ObjectBoolean) ObjectDataType() ObjectsDataType  { return OBJECT_BOOLEAN } // Return boolean object
func (ObjString *ObjectString) ObjectDataType() ObjectsDataType    { return OBJECT_STRING }  // Return string object
func (ObjArray *ObjectArray) ObjectDataType() ObjectsDataType      { return OBJECT_ARRAY }   // Return array object
func (ObjNullField *ObjectNULL) ObjectDataType() ObjectsDataType   { return OBJECT_NULL }    // Return null object
func (ObjErrorField *ObjectERROR) ObjectDataType() ObjectsDataType { return OBJECT_ERROR }   // Return Error object
func (ObjBuiltIn *ObjectBUILTINFUNCTION) ObjectDataType() ObjectsDataType {
	return OBJECT_BUILT_IN_FUNCTION
} // Return built in function object

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//
//  													 _____ _       __    _
// 														|   __| |_ _ _|  |  |_|___ ___
// 													    |__   | '_| | |  |__| |   | -_|
// 														|_____|_,_|_  |_____|_|_|_|___|
//		   														  |___|
//
//
// The SkyLine configuration language is a language and engine designed to act as a modification language to the SkyLine programming language. This language is
// very minimal and contains a regex base lexer, a very basic parser, a few keywords, a base interpreter and that is all as well as some backend engine code. This
// language is purely modified to be an extension to the SkyLine programming language, something that can be a pre processor language post processing for the main
// SkyLine script. Below is more technical information for the language
//
// Lexer       : Regex based lexer with minimal tokens and base syntax
// Parser      : Base parser with minimal tokens and base syntax with simple error systems
// REPL        : Does not exist
// Environment : Extremely minimal
// Types       : String, Boolean, Integer
// Statements  : set, import, use, errors, output, system, constant/const
//
//
// File contains -> This file hashes the SkyLine configuration language data types and allows for valid parsing

func (ObjInteger *ObjectInteger) GrabHashKey() HashingKey {
	return HashingKey{
		Type:  ObjInteger.ObjectDataType(),
		Value: uint64(ObjInteger.Value),
	}
}

func (ObjBoolean *ObjectBoolean) GrabHashKey() HashingKey {
	var val uint64
	if ObjBoolean.Value {
		val = 1 // true:BINARY
	} else {
		val = 0 // false:BINARY
	}
	return HashingKey{
		Type:  ObjBoolean.ObjectDataType(),
		Value: val,
	}
}

func (ObjString *ObjectString) GrabHashKey() HashingKey {
	STR_HASH := fnv.New64a()
	STR_HASH.Write([]byte(ObjString.Value))
	return HashingKey{
		Value: STR_HASH.Sum64(),
		Type:  ObjString.ObjectDataType(),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//
//  													 _____ _       __    _
// 														|   __| |_ _ _|  |  |_|___ ___
// 													    |__   | '_| | |  |__| |   | -_|
// 														|_____|_,_|_  |_____|_|_|_|___|
//		   														  |___|
//
//
// The SkyLine configuration language is a language and engine designed to act as a modification language to the SkyLine programming language. This language is
// very minimal and contains a regex base lexer, a very basic parser, a few keywords, a base interpreter and that is all as well as some backend engine code. This
// language is purely modified to be an extension to the SkyLine programming language, something that can be a pre processor language post processing for the main
// SkyLine script. Below is more technical information for the language
//
// Lexer       : Regex based lexer with minimal tokens and base syntax
// Parser      : Base parser with minimal tokens and base syntax with simple error systems
// REPL        : Does not exist
// Environment : Extremely minimal
// Types       : String, Boolean, Integer
// Statements  : set, import, use, errors, output, system, constant/const
//
//
// File contains -> This file contains information on the objects and functions that allow the interpreter or engine to access the current values of the object

func (ObjInteger *ObjectInteger) ObjectInspectFunc() string {
	return fmt.Sprintf("%d", ObjInteger.Value)
} // Grab the integer value

func (ObjBoolean *ObjectBoolean) ObjectInspectFunc() string {
	return fmt.Sprintf("%t", ObjBoolean.Value)
} // Grab the boolean value

func (ObjString *ObjectString) ObjectInspectFunc() string {
	return ObjString.Value
} // Grab the string value

func (ObjNULL *ObjectNULL) ObjectInspectFunc() string {
	return "NULL"
} // Grab the null field value

func (ObjERROR *ObjectERROR) ObjectInspectFunc() string {
	return "Error: " + ObjERROR.Message
} // Grab the error field and message

func (ObjBuiltIn *ObjectBUILTINFUNCTION) ObjectInspectFunc() string {
	return "Object<BUILT-IN-FUNCTION>"
} // Grab the function value

func (ObjArray *ObjectArray) ObjectInspectFunc() string {
	var Out bytes.Buffer
	ArrayElems := []string{}
	for _, element := range ObjArray.Elements {
		ArrayElems = append(ArrayElems, element.ObjectInspectFunc())
	}
	Out.WriteString("[")
	Out.WriteString(strings.Join(ArrayElems, ", "))
	Out.WriteString("];")
	return Out.String()
}
