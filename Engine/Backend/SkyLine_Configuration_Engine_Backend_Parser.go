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
//
// File contains -> This file contains every single part of the SkyLine configuration language parser
//
//

import (
	"fmt"
	"strconv"
)

func NewParser(Scanner *ScannerStructureRegister) *SLC_Parser_Structure {
	SLC_Parser := &SLC_Parser_Structure{
		Scanner:      Scanner,
		EngineErrors: []string{},
	}
	SLC_Parser.PrefixParseFunctions = map[TokenDataType]PrefixParserFunctions{
		IDENTIFIER:       SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_IDENTIFIER,
		OBJECT_INTEGER:   SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_INTEGER,
		OBJECT_STRING:    SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_STRING,
		BOOL_TRUE_Token:  SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_BOOLEAN_TRUTHY,
		BOOL_FALSE_Token: SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_BOOLEAN_TRUTHY,
		LPAREN_Token:     SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_GROUPED_LIST_EXPRESSION,
		LBRACKET_Token:   SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_ARRAY,
		ENGINE_Token:     SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_ENGINE_SECTOR,
	}
	SLC_Parser.InfixParseFunctions = map[TokenDataType]InfixParserFunctions{
		MODIFY_Token:   SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_INFIX_EXPRESSION,
		LPAREN_Token:   SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_CALL_EXPRESSION,
		LBRACKET_Token: SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_INDEX_EXPRESSION,
	}
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	return SLC_Parser
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN() {
	SLC_Parser.CurrentToken = SLC_Parser.PeekToken
	SLC_Parser.PeekToken = SLC_Parser.Scanner.NT()
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_CURRENT_TOKEN_IS_THIS(t TokenDataType) bool {
	return SLC_Parser.CurrentToken.TokenDataType == t
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_PEEK_TOKEN_IS_THIS(t TokenDataType) bool {
	return SLC_Parser.PeekToken.TokenDataType == t
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_EXPECT_PEEK(t TokenDataType) bool {
	if SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(t) {
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
		return true
	} else {
		SLC_Parser.peekError(t)
		return false
	}
}

func (SLC_Parser *SLC_Parser_Structure) Errors() []string {
	return SLC_Parser.EngineErrors
}

func (SLC_Parser *SLC_Parser_Structure) peekError(t TokenDataType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, SLC_Parser.PeekToken.TokenDataType)
	SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
}

func (SLC_Parser *SLC_Parser_Structure) noPrefixParseFnError(t TokenDataType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
}

func (SLC_Parser *SLC_Parser_Structure) ParseProgram() *Engine_Prog {
	program := &Engine_Prog{}
	program.Statements = []AbstractSyntaxTree_Statement{}

	for !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(END_OF_FILE) {
		stmt := SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_STATEMENT()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	}

	return program
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_STATEMENT() AbstractSyntaxTree_Statement {
	switch SLC_Parser.CurrentToken.TokenDataType {
	case SET_Token:
		return SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_VARIABLE_DECL()
	case CONSTANT_Token:
		return SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_CONSTANT_STATEMENT()
	default:
		return SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_EXPRESSION_STATEMENT()
	}
}

// parseConstStatement parses a constant declaration.
func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_CONSTANT_STATEMENT() *Constant_Statement_AbstractSyntaxTree {
	stmt := &Constant_Statement_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}
	if !SLC_Parser.ENGINE_EXPECT_PEEK(IDENTIFIER) {
		return nil
	}
	stmt.Name = &Identifier_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Value: SLC_Parser.CurrentToken.Literal}
	if !SLC_Parser.ENGINE_EXPECT_PEEK(ASSIGN_Token) {
		return nil
	}
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	stmt.Value = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)
	for !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(SEMICOLON_Token) {

		if SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(END_OF_FILE) {
			SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, "unterminated const statement")
			return nil
		}

		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	}
	return stmt
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_VARIABLE_DECL() *Assignment_Statement_AbstractSyntaxTree {
	stmt := &Assignment_Statement_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}

	if !SLC_Parser.ENGINE_EXPECT_PEEK(IDENTIFIER) {
		return nil
	}

	stmt.Name = &Identifier_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Value: SLC_Parser.CurrentToken.Literal}

	if !SLC_Parser.ENGINE_EXPECT_PEEK(ASSIGN_Token) {
		return nil
	}

	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()

	stmt.Value = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)

	if SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(SEMICOLON_Token) {
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	}

	return stmt
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_EXPRESSION_STATEMENT() *Expression_Statement_AbstractSyntaxTree {
	stmt := &Expression_Statement_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}

	stmt.Expression = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)

	if SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(SEMICOLON_Token) {
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	}

	return stmt
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(precedence int) AbstractSyntaxTree_Expression {
	prefix := SLC_Parser.PrefixParseFunctions[SLC_Parser.CurrentToken.TokenDataType]
	if prefix == nil {
		SLC_Parser.noPrefixParseFnError(SLC_Parser.CurrentToken.TokenDataType)
		return nil
	}
	leftExp := prefix()

	for !SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(SEMICOLON_Token) && precedence < SLC_Parser.peekPrecedence() {
		infix := SLC_Parser.InfixParseFunctions[SLC_Parser.PeekToken.TokenDataType]
		if infix == nil {
			return leftExp
		}

		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (SLC_Parser *SLC_Parser_Structure) peekPrecedence() int {
	if p, ok := ParserPrecedences[SLC_Parser.PeekToken.TokenDataType]; ok {
		return p
	}

	return LOWEST
}

func (SLC_Parser *SLC_Parser_Structure) curPrecedence() int {
	if p, ok := ParserPrecedences[SLC_Parser.CurrentToken.TokenDataType]; ok {
		return p
	}

	return LOWEST
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_IDENTIFIER() AbstractSyntaxTree_Expression {
	return &Identifier_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Value: SLC_Parser.CurrentToken.Literal}
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_INTEGER() AbstractSyntaxTree_Expression {
	lit := &IntegerDataType_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}

	value, err := strconv.ParseInt(SLC_Parser.CurrentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", SLC_Parser.CurrentToken.Literal)
		SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_STRING() AbstractSyntaxTree_Expression {
	return &StringDataType_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Value: SLC_Parser.CurrentToken.Literal}
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_PREFIX_EXPRESSION() AbstractSyntaxTree_Expression {
	expression := &PrefixExpression_Expression_AbstractSyntaxTree{
		TokenRegister: SLC_Parser.CurrentToken,
		Operator:      SLC_Parser.CurrentToken.Literal,
	}

	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()

	expression.Right = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(PREFIX)

	return expression
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_INFIX_EXPRESSION(left AbstractSyntaxTree_Expression) AbstractSyntaxTree_Expression {
	expression := &InfixExpression_Expression_AbstractSyntaxTree{
		TokenRegister: SLC_Parser.CurrentToken,
		Operator:      SLC_Parser.CurrentToken.Literal,
		Left:          left,
	}

	precedence := SLC_Parser.curPrecedence()
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	expression.Right = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(precedence)
	fmt.Println("Parsing expression -> ", expression.Left.TokenConstructToString()+expression.Right.TokenConstructToString())
	return expression
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_BOOLEAN_TRUTHY() AbstractSyntaxTree_Expression {
	return &BooleanDataType_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Value: SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(BOOL_TRUE_Token)}
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_GROUPED_LIST_EXPRESSION() AbstractSyntaxTree_Expression {
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()

	exp := SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)

	if !SLC_Parser.ENGINE_EXPECT_PEEK(RPAREN_Token) {
		return nil
	}

	return exp
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_ENGINE_UNIT_BLOCK_STATEMENTS() *BlockStatement_Statement_AbstractSyntaxTree {
	block := &BlockStatement_Statement_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}
	block.Statements = []AbstractSyntaxTree_Statement{}

	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()

	for !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(RBRACE_Token) && !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(END_OF_FILE) {
		stmt := SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_STATEMENT()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	}

	return block
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_CALL_EXPRESSION(function AbstractSyntaxTree_Expression) AbstractSyntaxTree_Expression {
	exp := &CallFunction_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Function: function}
	exp.Arguments = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_EXPRESSION_LIST(RPAREN_Token)
	return exp
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_EXPRESSION_LIST(end TokenDataType) []AbstractSyntaxTree_Expression {
	list := []AbstractSyntaxTree_Expression{}

	if SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(end) {
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
		return list
	}

	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	list = append(list, SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST))

	for SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(COMMA_Token) {
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
		list = append(list, SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST))
	}

	if !SLC_Parser.ENGINE_EXPECT_PEEK(end) {
		return nil
	}

	return list
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_ARRAY() AbstractSyntaxTree_Expression {
	array := &ArrayLiteral_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}

	array.Elements = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_EXPRESSION_LIST(RBRACKET_Token)

	return array
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_INDEX_EXPRESSION(left AbstractSyntaxTree_Expression) AbstractSyntaxTree_Expression {
	exp := &IndexLit_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken, Left: left}

	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	exp.Index = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)

	if !SLC_Parser.ENGINE_EXPECT_PEEK(RBRACKET_Token) {
		return nil
	}

	return exp
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_BRACKET_ENGINE_KEY_SECTOR() AbstractSyntaxTree_Expression {
	if !SLC_Parser.ENGINE_EXPECT_PEEK(LPAREN_Token) {
		msg := fmt.Sprintf("expected '(' but got %s", SLC_Parser.CurrentToken.Literal)
		SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
		return nil
	}
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	tmp := SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST)
	if tmp == nil {
		return nil
	}
	if !SLC_Parser.ENGINE_EXPECT_PEEK(RPAREN_Token) {
		msg := fmt.Sprintf("expected ')' but got %s", SLC_Parser.CurrentToken.Literal)
		SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
		return nil
	}
	return tmp
}

func (SLC_Parser *SLC_Parser_Structure) ENGINE_UNIT_CALL_PARSE_FUNCTION_ENGINE_SECTOR() AbstractSyntaxTree_Expression {
	expression := &ENGINE_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}
	expression.Value = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_BRACKET_ENGINE_KEY_SECTOR()
	if expression.Value == nil {
		return nil
	}
	if !SLC_Parser.ENGINE_EXPECT_PEEK(LBRACE_Token) {
		return nil
	}
	SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
	for !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(RBRACE_Token) {
		if SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(END_OF_FILE) {
			SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, "unterminated switch statement")
			return nil
		}
		tmp := &INIT_Expression_AbstractSyntaxTree{TokenRegister: SLC_Parser.CurrentToken}
		if SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(ENGINE_INITATE_Token) {
			SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
			tmp.Expression = append(tmp.Expression, SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST))
			for SLC_Parser.ENGINE_PEEK_TOKEN_IS_THIS(COMMA_Token) {
				SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
				SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
				tmp.Expression = append(tmp.Expression, SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_PARSE_EXPRESSION(LOWEST))
			}
		}
		if !SLC_Parser.ENGINE_EXPECT_PEEK(LBRACE_Token) {
			msg := fmt.Sprintf("expected token to be '{', got %s instead", SLC_Parser.CurrentToken.TokenDataType)
			SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
			fmt.Printf("error\n")
			return nil
		}
		tmp.Sub_UNIT = SLC_Parser.ENGINE_UNIT_CALL_PARSE_FUNCTION_ENGINE_UNIT_BLOCK_STATEMENTS()
		if !SLC_Parser.ENGINE_CURRENT_TOKEN_IS_THIS(RBRACE_Token) {
			msg := fmt.Sprintf("Syntax Error: expected token to be '}', got %s instead", SLC_Parser.CurrentToken.TokenDataType)
			SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
			fmt.Printf("error\n")
			return nil
		}
		SLC_Parser.ENGINE_PARSER_LOAD_NEXT_FAULT_TOKEN()
		expression.SubUnits = append(expression.SubUnits, tmp)
	}
	count := 0
	for _, c := range expression.SubUnits {
		if c.Default {
			count++
		}
	}
	if count > 1 {
		msg := "A switch-statement should only have one default block"
		SLC_Parser.EngineErrors = append(SLC_Parser.EngineErrors, msg)
		return nil

	}
	return expression

}
