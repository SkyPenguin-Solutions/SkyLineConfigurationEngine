package SkyLine_Configuration_Engine_Backend_Source

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
// Contains -> This file contains functions that are dedicated to checking, prepping and setting a variable before it can be registered into the environment
//

import (
	"fmt"
	"log"
	"os"
)

func (Env *Engine_Environment_Of_Environment) Engine_Set_Varname(NameObj string, Value SLC_Object) SLC_Object {
	Current := Env.StoreObj[NameObj]
	if Current != nil && Env.ReadOnly[NameObj] {
		log.Fatalf("Stop trying to modify a constant value, this is just insane, just stop -> %s is a CONSTANT VALUE AND CAN NOT BE CHANGED", NameObj)
	}
	if len(Env.PermitMod) > 0 {
		for _, valv := range Env.PermitMod {
			if valv == NameObj {
				Env.StoreObj[NameObj] = Value
				return Value
			}
		}
		if Env.EngineOuter != nil {
			return Env.EngineOuter.Engine_Set_Varname(NameObj, Value)
		}
		fmt.Println("Had a scope error -> ISSUE_DEV_LOG")
		os.Exit(5)
	}
	Env.StoreObj[NameObj] = Value
	return Value
}
