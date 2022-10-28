package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/* General description of the JackTokenizer struct and functions: "Removes all comments and white space from the input stream
and breaks it down into Jack language tokens, as specified by the Jack grammar.*/

// List of token types:
const (
	KEYWORD      = 0
	SYMBOL       = 1
	IDENTIFIER   = 2
	INT_CONST    = 3
	STRING_CONST = 4
)

// List of keyword:
const (
	CLASS = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

var kerywordlist = map[string]int{
	"class":       CLASS,
	"method":      METHOD,
	"function":    FUNCTION,
	"constructor": CONSTRUCTOR,
	"int":         INT,
	"boolean":     BOOLEAN,
	"char":        CHAR,
	"void":        VOID,
	"var":         VAR,
	"static":      STATIC,
	"field":       FIELD,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
}

// Regular expression
const (
	keywordRegularExpression    = "class|method|function|constructor|int|boolean|char|void|var|static|field|let|do|if|else|while|return|true|false|null|this"
	symbolRegularExpression     = "[\\{|\\}|\\(|\\)|\\[|\\]|\\.|\\,|\\.|\\;|\\+|\\-|\\*|\\/|\\&|\\||\\<|\\>|\\=|\\~]"
	identifierRegularExpression = "[\\w_]+"
	integerRegularExpression    = "[0-9]+"
	stringRegularExpression     = "\"[^\"\n]*\""
)

type JackTokenizer struct {
	file         *os.File
	scanner      *bufio.Scanner
	tokenPointer int
	tokenList    []string
	writerT      *bufio.Writer
}

func initJackTokenizer(file *os.File, writerT *bufio.Writer) *JackTokenizer {
	symbolREGEX, _ := regexp.Compile(symbolRegularExpression)
	scanner := bufio.NewScanner(file)
	var jackTokenizer JackTokenizer = JackTokenizer{file: file, scanner: scanner, tokenPointer: 0, writerT: writerT}
	var longcomment bool = false

	for jackTokenizer.scanner.Scan() {
		var line string = jackTokenizer.scanner.Text()
		if strings.HasPrefix(line, "//") { // Case: this line is a comment
			continue
		} else if strings.TrimSpace(line) == "" { // Case: this line is empty
			continue
		} else {
			var preprocessedCommand string = strings.Split(line, "//")[0]
			preprocessedCommand = strings.TrimLeft(preprocessedCommand, " ")
			if strings.HasPrefix(preprocessedCommand, "/*") && strings.HasSuffix(preprocessedCommand, "*/") {
				continue
			}
			if strings.HasPrefix(preprocessedCommand, "/*") && !strings.HasSuffix(preprocessedCommand, "*/") {
				longcomment = true
				continue
			}
			if longcomment && strings.HasSuffix(preprocessedCommand, "*/") {
				longcomment = false
				continue
			} else if longcomment {
				continue
			}
			preprocessedCommand = strings.TrimSpace(preprocessedCommand)
			chars := []rune(preprocessedCommand)
			var inputstring string = ""
			var isDoubleQuote bool = false
			for i := 0; i < len(chars); i++ {
				char := string(chars[i])
				if !isDoubleQuote {
					if char == " " {
						if inputstring != "" {
							jackTokenizer.appendToken(inputstring)
							inputstring = ""
						}
						continue
					} else if symbolREGEX.MatchString(char) {
						if inputstring != "" {
							jackTokenizer.appendToken(inputstring)
							inputstring = ""
						}
						jackTokenizer.appendToken(char)
					} else if char == "\"" {
						isDoubleQuote = true
					} else {
						inputstring += char
					}
				} else {
					if char != "\"" {
						inputstring += char
					} else {
						isDoubleQuote = false
						inputstring = "\"" + inputstring + "\""
						jackTokenizer.appendToken(inputstring)
						inputstring = ""
					}
				}
			}
		}
	}
	jackTokenizer.writeTXML()
	return &jackTokenizer
}

func (jackTokenizer *JackTokenizer) appendToken(currentToken string) {
	jackTokenizer.tokenList = append(jackTokenizer.tokenList, currentToken)
	jackTokenizer.tokenPointer = jackTokenizer.tokenPointer + 1
}

func (jackTokenizer *JackTokenizer) getCurrentToken() string {
	if jackTokenizer.tokenPointer > len(jackTokenizer.tokenList) {
		jackTokenizer.tokenPointer = jackTokenizer.tokenPointer - len(jackTokenizer.tokenList)
	}
	return jackTokenizer.tokenList[jackTokenizer.tokenPointer]
}

// Question: "Do we have more tokens in the input?"
func (jackTokenizer *JackTokenizer) hasMoreTokens() bool {
	if jackTokenizer.tokenPointer < len(jackTokenizer.tokenList) {
		return true
	} else {
		return false
	}
}

// "Gets the next token from the input and makes it the current token."
func (jackTokenizer *JackTokenizer) advance() {
	if jackTokenizer.hasMoreTokens() {
		jackTokenizer.tokenPointer = jackTokenizer.tokenPointer + 1
	}
}

// "Returns the type of the current token."
func (jackTokenizer *JackTokenizer) tokenType() int {
	keywordREGEX, _ := regexp.Compile(keywordRegularExpression)
	symbolREGEX, _ := regexp.Compile(symbolRegularExpression)
	identifierREGEX, _ := regexp.Compile(identifierRegularExpression)
	integerREGEX, _ := regexp.Compile(integerRegularExpression)
	stringREGEX, _ := regexp.Compile(stringRegularExpression)

	if integerREGEX.MatchString(jackTokenizer.getCurrentToken()) {
		return INT_CONST
	} else if stringREGEX.MatchString(jackTokenizer.getCurrentToken()) {
		return STRING_CONST
	} else if keywordREGEX.MatchString(jackTokenizer.getCurrentToken()) {
		return KEYWORD
	} else if symbolREGEX.MatchString(jackTokenizer.getCurrentToken()) {
		return SYMBOL
	} else if identifierREGEX.MatchString(jackTokenizer.getCurrentToken()) {
		return IDENTIFIER
	} else {
		return -1
	}
}

// "Returns the keyword which is the current token. Should be called only when tokenType() is KEYWORD."
func (jackTokenizer *JackTokenizer) keyWord() int {
	return kerywordlist[jackTokenizer.getCurrentToken()]
}

// "Returns the character which is the current token. Should be called only when tokenType() is SYMBOL."
func (jackTokenizer *JackTokenizer) symbol() string {
	return jackTokenizer.getCurrentToken()
}

// "Returns the identifier which is the current token. Should be called only when tokenType() is IDENTIFIER."
func (jackTokenizer *JackTokenizer) identifier() string {
	return jackTokenizer.getCurrentToken()
}

// "Returns the integer value of the current token. Should be called only when tokenType() is INT_CONST."
func (jackTokenizer *JackTokenizer) intVal() int {
	currentINT, _ := strconv.Atoi(jackTokenizer.getCurrentToken())
	return currentINT
}

/*
"Returns the string value of the current token, without the double quotes.
Should be called only when tokenType() is STRING_CONST."
*/
func (jackTokenizer *JackTokenizer) stringVal() string {
	var currentSTRING string = jackTokenizer.getCurrentToken()
	return currentSTRING[1 : len(currentSTRING)-1]
}

func (jackTokenizer *JackTokenizer) writeTXML() {
	jackTokenizer.writerT.WriteString("<tokens>\n")
	jackTokenizer.tokenPointer = 0
	for i := 0; i < len(jackTokenizer.tokenList); i++ {
		tokenType := jackTokenizer.tokenType()
		switch tokenType {
		case INT_CONST:
			jackTokenizer.writerT.WriteString("<integerConstant> " + jackTokenizer.getCurrentToken() + " </integerConstant>\n")
		case STRING_CONST:
			jackTokenizer.writerT.WriteString("<stringConstant> " + jackTokenizer.getCurrentToken() + " </stringConstant>\n")
		case KEYWORD:
			jackTokenizer.writerT.WriteString("<keyword> " + jackTokenizer.getCurrentToken() + " </keyword>\n")
		case SYMBOL:
			jackTokenizer.writerT.WriteString("<symbol> " + jackTokenizer.getCurrentToken() + " </symbol>\n")
		case IDENTIFIER:
			jackTokenizer.writerT.WriteString("<identifier> " + jackTokenizer.getCurrentToken() + " </identifier>\n")
		}
		jackTokenizer.tokenPointer += 1
	}
	jackTokenizer.writerT.WriteString("</tokens>\n")
	jackTokenizer.tokenPointer = 0
}

/*
General description of the CompilationEngine struct and functions: "Effects the actual compilation output.
Gets its input from a JackTokenizer and emits its parsed structure into an output file/stream.
The output is generated by a series of compilexxx() routines, one for every syntactic element xxx of the Jack grammar.
The contract between these routines is that each compilexxx() routine should read the syntactic construct xxx from the input,
advance() the tokenizer exactly beyond xxx, and output the parsing of xxx.
Thus, compilexxx() may only be called if indeed xxx is the next syntactic element of the input.
This module emits structured printout of the code, wrapped in XML tags.
*/
type CompilationEngine struct {
	jackTokenizer *JackTokenizer
	file          *os.File
	writer        *bufio.Writer
}

func initCompilationEngine(jackTokenizer *JackTokenizer, file *os.File, writer *bufio.Writer) *CompilationEngine {
	compilationEngine := &CompilationEngine{jackTokenizer: jackTokenizer, file: file, writer: writer}
	return compilationEngine
}

// "Compiles a complete class."
func (ce *CompilationEngine) CompileClass() {
	ce.writer.WriteString("<class>\n")
	ce.CompileKeyWord()
	ce.IncrementPointer()
	ce.CompileIdentifier()
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is {
	ce.CompileClassVarDec()
	ce.CompileSubroutine()
	ce.CompileSymbol() // Here symbol is }
	ce.writer.WriteString("</class>\n")
}

// "Compiles a static declaration or a field declaration."
func (ce *CompilationEngine) CompileClassVarDec() {
	ce.IncrementPointer()
	for ce.jackTokenizer.tokenType() == KEYWORD &&
		(ce.jackTokenizer.keyWord() == STATIC || ce.jackTokenizer.keyWord() == FIELD) {
		ce.writer.WriteString("<classVarDec>\n")
		ce.CompileKeyWord()
		ce.CompileVarDec()
		ce.writer.WriteString("</classVarDec>\n")
		ce.IncrementPointer()
	}
	ce.DecrementPointer()
}

// "Compiles a complete method, function, or constructor."
func (ce *CompilationEngine) CompileSubroutine() {
	ce.IncrementPointer()
	for ce.jackTokenizer.tokenType() == KEYWORD &&
		(ce.jackTokenizer.keyWord() == CONSTRUCTOR || ce.jackTokenizer.keyWord() == FUNCTION ||
			ce.jackTokenizer.keyWord() == METHOD) {
		ce.writer.WriteString("<subroutineDec>\n")
		ce.CompileKeyWord()
		ce.IncrementPointer()
		if ce.jackTokenizer.tokenType() == KEYWORD && (ce.jackTokenizer.keyWord() == BOOLEAN ||
			ce.jackTokenizer.keyWord() == CHAR || ce.jackTokenizer.keyWord() == INT || ce.jackTokenizer.keyWord() == VOID) {
			ce.CompileKeyWord()
		} else if ce.jackTokenizer.tokenType() == IDENTIFIER {
			ce.CompileIdentifier()
		}
		ce.IncrementPointer()
		ce.CompileIdentifier()
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is (
		ce.writer.WriteString("<parameterList>\n")
		ce.CompileParameterList()
		ce.writer.WriteString("</parameterList>\n")
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is )
		ce.writer.WriteString("<subroutineBody>\n")
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is {
		ce.IncrementPointer()
		for ce.jackTokenizer.tokenType() == KEYWORD && ce.jackTokenizer.keyWord() == VAR {
			ce.writer.WriteString("<varDec>\n")
			ce.writer.WriteString("<keyword> var </keyword>\n")
			ce.CompileVarDec()
			ce.writer.WriteString("</varDec>\n")
			ce.IncrementPointer()
		}
		ce.DecrementPointer()
		ce.CompileStatements()
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is }
		ce.writer.WriteString("</subroutineBody>\n</subroutineDec>\n")
		ce.IncrementPointer()
	}
}

// "Compiles a (possibly empty) parameter list, not including the enclosing '()'"
func (ce *CompilationEngine) CompileParameterList() {
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == KEYWORD && (ce.jackTokenizer.keyWord() == BOOLEAN ||
		ce.jackTokenizer.keyWord() == CHAR || ce.jackTokenizer.keyWord() == INT || ce.jackTokenizer.keyWord() == VOID) {
		ce.CompileKeyWord()
	} else if ce.jackTokenizer.tokenType() == IDENTIFIER {
		ce.CompileIdentifier()
	} else {
		ce.DecrementPointer()
		return
	}
	ce.IncrementPointer()
	ce.CompileIdentifier()
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "," {
		ce.CompileSymbol() // Here symbol is ,
		ce.CompileParameterList()
	} else {
		ce.DecrementPointer()
	}
}

// "Compiles a var declaration"
func (ce *CompilationEngine) CompileVarDec() {
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == KEYWORD && (ce.jackTokenizer.keyWord() == BOOLEAN ||
		ce.jackTokenizer.keyWord() == CHAR || ce.jackTokenizer.keyWord() == INT || ce.jackTokenizer.keyWord() == VOID) {
		ce.CompileKeyWord()
	} else if ce.jackTokenizer.tokenType() == IDENTIFIER {
		ce.CompileIdentifier()
	}
	ce.IncrementPointer()
	ce.CompileIdentifier()
	ce.IncrementPointer()
	for ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "," {
		ce.CompileSymbol() // Here symbol is ,
		ce.IncrementPointer()
		ce.CompileIdentifier()
		ce.IncrementPointer()
	}
	ce.CompileSymbol() // Here symbol is ;
}

// "Compiles a sequence of statements, not including the enclosing '{}'"
func (ce *CompilationEngine) CompileStatements() {
	ce.writer.WriteString("<statements>\n")
	ce.IncrementPointer()
	for ce.jackTokenizer.tokenType() == KEYWORD {
		switch ce.jackTokenizer.keyWord() {
		case DO:
			ce.CompileDo()
		case LET:
			ce.CompileLet()
		case WHILE:
			ce.CompileWhile()
		case RETURN:
			ce.CompileReturn()
		case IF:
			ce.CompileIf()
		}
		ce.IncrementPointer()
	}
	ce.DecrementPointer()
	ce.writer.WriteString("</statements>\n")
}

// "Compiles a do statement"
func (ce *CompilationEngine) CompileDo() {
	ce.writer.WriteString("<doStatement>\n")
	ce.CompileKeyWord() // Here keyword is do
	ce.IncrementPointer()
	ce.CompileIdentifier()
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is ( or .
	if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "." {
		ce.IncrementPointer()
		ce.CompileIdentifier()
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is (
	}
	ce.CompileExpressionList()
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is )
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is ;
	ce.writer.WriteString("</doStatement>\n")
}

// "Compiles a let statement"
func (ce *CompilationEngine) CompileLet() {
	ce.writer.WriteString("<letStatement>\n")
	ce.CompileKeyWord() // Here keyword is let
	ce.IncrementPointer()
	ce.CompileIdentifier()
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "[" {
		ce.CompileSymbol() // Here symbol is [
		ce.CompileExpression()
		ce.CompileSymbol() // Here symbol is ]
		ce.IncrementPointer()
	}
	ce.CompileSymbol() // Here symbol is =
	ce.CompileExpression()
	ce.CompileSymbol() // Here symbol is ;
	ce.writer.WriteString("</letStatement>\n")
}

// "Compiles a while statement"
func (ce *CompilationEngine) CompileWhile() {
	ce.writer.WriteString("<whileStatement>\n")
	ce.CompileKeyWord() // Here keyword is while
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is (
	ce.CompileExpression()
	ce.CompileSymbol() // Here symbol is )
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is {
	ce.CompileStatements()
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is }
	ce.writer.WriteString("</whileStatement>\n")
}

// "Compiles a return statement"
func (ce *CompilationEngine) CompileReturn() {
	ce.writer.WriteString("<returnStatement>\n")
	ce.CompileKeyWord() // Here keyword is return
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == ";" {
		ce.CompileSymbol() // Here symbol is ;
		ce.writer.WriteString("</returnStatement>\n")
		return
	} else {
		ce.DecrementPointer()
		ce.CompileExpression()
		ce.CompileSymbol() // Here symbol is ;
		ce.writer.WriteString("</returnStatement>\n")
	}
}

// "Compiles a if statement, possibly with a trailing else clause."
func (ce *CompilationEngine) CompileIf() {
	ce.writer.WriteString("<ifStatement>\n")
	ce.CompileKeyWord() // Here keyword is if
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is (
	ce.CompileExpression()
	ce.CompileSymbol() // Here symbol is )
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is {
	ce.CompileStatements()
	ce.IncrementPointer()
	ce.CompileSymbol() // Here symbol is }
	ce.IncrementPointer()
	if ce.jackTokenizer.tokenType() == KEYWORD && ce.jackTokenizer.keyWord() == ELSE {
		ce.CompileKeyWord() // Here keyword is else
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is {
		ce.CompileStatements()
		ce.IncrementPointer()
		ce.CompileSymbol() // Here symbol is }
	} else {
		ce.DecrementPointer()
	}
	ce.writer.WriteString("</ifStatement>\n")
}

// "Compiles an expression"
func (ce *CompilationEngine) CompileExpression() {
	var output string = ce.CompileExpressionOutput()
	if output != "" {
		ce.writer.WriteString("<expression>\n" + output + "</expression>\n")
	}
}

func (ce *CompilationEngine) CompileExpressionOutput() string {
	var output string = ""
	ce.IncrementPointer()
	output += ce.CompileTerm()
	ce.IncrementPointer()

	var symbolRegularExpression2 string = "[\\+|\\-|\\*|\\/|\\&|\\||\\<|\\>|\\=]"
	symbolREGEX2, _ := regexp.Compile(symbolRegularExpression2)

	for ce.jackTokenizer.tokenType() == SYMBOL && symbolREGEX2.MatchString(ce.jackTokenizer.symbol()) {
		switch ce.jackTokenizer.symbol() {
		case ">":
			output += "<symbol> &gt; </symbol>\n"
		case "<":
			output += "<symbol> &lt; </symbol>\n"
		case "&":
			output += "<symbol> &amp; </symbol>\n"
		default:
			output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n"
		}
		ce.IncrementPointer()
		output += ce.CompileTerm()
		ce.IncrementPointer()
	}
	return output
}

/*
"Compiles a term.
This routine is faced with a slight difficulty when trying to decide between some of the alternative parsing rules.
Specifically, if the current token is an identifier, the routine must distinguish between a variable,
an array entry, and a subroutine call. A single look-ahead token, which may be one of '[', '(', '.'
suffcies to distinguish between the three possibilites.
Any other token is not part of this term and should not be advanced over."
*/
func (ce *CompilationEngine) CompileTerm() string {
	var output string = ce.CompileTermOutput()
	if output != "" {
		return "<term>\n" + output + "</term>\n"
	}
	return ""
}

func (ce *CompilationEngine) CompileTermOutput() string {
	var output string = ""
	if ce.jackTokenizer.tokenType() == STRING_CONST {
		output += "<stringConstant> " + ce.jackTokenizer.stringVal() + " </stringConstant>\n"
	} else if ce.jackTokenizer.tokenType() == INT_CONST {
		output += "<integerConstant> " + strconv.Itoa(ce.jackTokenizer.intVal()) + " </integerConstant>\n"
	} else if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "(" {
		output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is (
		output += "<expression>\n" + ce.CompileExpressionOutput() + "</expression>\n"
		output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is )
	} else if ce.jackTokenizer.tokenType() == KEYWORD && (ce.jackTokenizer.keyWord() == THIS ||
		ce.jackTokenizer.keyWord() == NULL || ce.jackTokenizer.keyWord() == TRUE || ce.jackTokenizer.keyWord() == FALSE) {
		output += "<keyword> " + ce.jackTokenizer.getCurrentToken() + " </keyword>\n"
	} else if ce.jackTokenizer.tokenType() == IDENTIFIER {
		output += "<identifier> " + ce.jackTokenizer.identifier() + " </identifier>\n"
		ce.IncrementPointer()
		if ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "[" {
			output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is [
			output += "<expression>\n" + ce.CompileExpressionOutput() + "</expression>\n"
			output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is ]
		} else if ce.jackTokenizer.tokenType() == SYMBOL && (ce.jackTokenizer.symbol() == "(" || ce.jackTokenizer.symbol() == ".") {
			output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n"
			if ce.jackTokenizer.symbol() == "." {
				ce.IncrementPointer()
				output += "<identifier> " + ce.jackTokenizer.identifier() + " </identifier>\n"
				ce.IncrementPointer()
				output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is (
			}
			output += "<expressionList>\n"
			var expressionOutput string = ce.CompileExpressionOutput()
			if expressionOutput != "" {
				output += "<expression>\n" + expressionOutput + "</expression>\n"
			}
			for ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "," {
				output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is ,
				expressionOutput = ce.CompileExpressionOutput()
				if expressionOutput != "" {
					output += "<expression>\n" + expressionOutput + "</expression>\n"
				} else {
					ce.DecrementPointer()
				}
			}
			output += "</expressionList>\n"
			output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is )
		} else {
			ce.DecrementPointer()
		}
	} else if ce.jackTokenizer.tokenType() == SYMBOL && (ce.jackTokenizer.symbol() == "-" || ce.jackTokenizer.symbol() == "~") {
		output += "<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n" // Here symbol is - or ~
		ce.IncrementPointer()
		output += ce.CompileTerm()
	} else {
		ce.DecrementPointer()
	}
	return output
}

// "Compiles a (possibly empty) comma-separated list of expressions."
func (ce *CompilationEngine) CompileExpressionList() {
	ce.writer.WriteString("<expressionList>\n")
	ce.CompileExpression()
	for ce.jackTokenizer.tokenType() == SYMBOL && ce.jackTokenizer.symbol() == "," {
		ce.CompileSymbol() // Here symbol is ,
		ce.CompileExpression()
	}
	ce.DecrementPointer()
	ce.writer.WriteString("</expressionList>\n")
}

func (ce *CompilationEngine) CompileSymbol() {
	ce.writer.WriteString("<symbol> " + ce.jackTokenizer.symbol() + " </symbol>\n")
}

func (ce *CompilationEngine) CompileIdentifier() {
	ce.writer.WriteString("<identifier> " + ce.jackTokenizer.identifier() + " </identifier>\n")
}

func (ce *CompilationEngine) CompileKeyWord() {
	ce.writer.WriteString("<keyword> " + ce.jackTokenizer.getCurrentToken() + " </keyword>\n")
}

func (ce *CompilationEngine) IncrementPointer() {
	if ce.jackTokenizer.hasMoreTokens() {
		ce.jackTokenizer.advance()
	}
}

func (ce *CompilationEngine) DecrementPointer() {
	ce.jackTokenizer.tokenPointer -= 1
}

// Jack Analyzer
func main() {
	filePath := os.Args[1]
	var jackTokenizer *JackTokenizer
	var compilationEngine *CompilationEngine
	var outputFileName string
	var outputFileT string
	files, _ := ioutil.ReadDir(filePath) // os.ReadDir not supported in Coursera
	for _, file := range files {
		if filepath.Ext(filepath.Join(filePath, file.Name())) == ".jack" {
			outputFileName = strings.TrimSuffix(filepath.Join(filePath, file.Name()), ".jack") + ".xml"
			output, _ := os.Create(outputFileName)
			defer output.Close()
			var writer *bufio.Writer = bufio.NewWriter(output)
			defer writer.Flush()
			outputFileT = strings.TrimSuffix(filepath.Join(filePath, file.Name()), ".jack") + "T.xml"
			outputT, _ := os.Create(outputFileT)
			defer outputT.Close()
			var writerT *bufio.Writer = bufio.NewWriter(outputT)
			defer writerT.Flush()
			file, _ := os.Open(filepath.Join(filePath, file.Name()))
			defer file.Close()
			jackTokenizer = initJackTokenizer(file, writerT)
			compilationEngine = initCompilationEngine(jackTokenizer, file, writer)
			compilationEngine.jackTokenizer.scanner = bufio.NewScanner(compilationEngine.jackTokenizer.file)
			compilationEngine.CompileClass()
		}
	}
}
