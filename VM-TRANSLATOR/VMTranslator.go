package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/* General description of the Parser struct & functions: "Handles the parsing of a single .vm file,
and encapsulates access to the input code. It reads VM commands, parses them,
and provides convenient access to their components. In addition, it removes all white space and comments."*/

// List of  Command Types:
const (
	C_ARITHMETIC = iota // Arithmetic commands
	C_PUSH              // Push onto stack command
	C_POP               // "Pop the top stack value and store it in segment[index]"" command
	C_LABEL             // Label declaration command
	C_GOTO              // Unconditional branching command
	C_IF                // Conditional branching command
	C_FUNCTION          // Function declaration command, "specifying the number of the function's local variables"
	C_RETURN            // Function invocation command, "specifying the number of the function's arguments"
	C_CALL              // "Transfer control back to the calling function" command
)

type Parser struct {
	file           *os.File
	scanner        *bufio.Scanner
	currentCommand string
	args           []string
}

func initParser(file *os.File) *Parser {
	scanner := bufio.NewScanner(file)
	var parser Parser = Parser{file: file, scanner: scanner}
	return &parser
}

// Question: "Are there more commands in the input?"
func (parser *Parser) hasMoreCommands() bool {
	for parser.scanner.Scan() {
		line := parser.scanner.Text()

		if strings.HasPrefix(line, "//") { // Case: this line is a comment
			continue
		} else if strings.TrimSpace(line) == "" { // Case: this line is empty
			continue
		} else {
			return true
		}
	}
	return false
}

/* "Reads the next command from the input and makes it the current command.
Should be called only if hasMoreCommands() is true. Initially, there is no current command."*/
func (parser *Parser) advance() {
	var inputCommand string = parser.scanner.Text() // Read next command
	parser.currentCommand = strings.Split(inputCommand, "//")[0]
	parser.args = strings.Split(parser.currentCommand, " ")
}

// "Returns the type of the current VM command. C_ARITHMETIC is returned for all the arithmetic commands."
func (parser *Parser) commandType() int {
	switch parser.args[0] {
	case "add":
		return C_ARITHMETIC
	case "sub":
		return C_ARITHMETIC
	case "neg":
		return C_ARITHMETIC
	case "eq":
		return C_ARITHMETIC
	case "gt":
		return C_ARITHMETIC
	case "lt":
		return C_ARITHMETIC
	case "and":
		return C_ARITHMETIC
	case "or":
		return C_ARITHMETIC
	case "not":
		return C_ARITHMETIC
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IF
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	case "call":
		return C_CALL
	default:
		fmt.Println(parser.args[0], "is an unrecongised command.")
		return -1
	}
}

/* "Return the first argument of the current command. In the case of C_ARITHMETIC, the command itself (add, sub, etc.)
is returned. Should not be called if the current command is C_RETURN."*/
func (parser *Parser) arg1() string {
	if parser.commandType() == C_ARITHMETIC {
		return parser.args[0]
	} else {
		return parser.args[1]
	}
}

/* "Return the second argument of the current command Should be called only if the current command is C_PUSH, C_POP,
C_FUNCTION, or C_CALL." */
func (parser *Parser) arg2() int {
	arg2, _ := strconv.Atoi(parser.args[2])
	return arg2
}

// General description of the CodeWriter struct and functions: "Translates VM commands into Hack assembly code."
type CodeWriter struct {
	file            *os.File
	comparisonCount int
	callCount       int
	currentFileName string
	currentFunction string
}

const (
	pushDcommand string = "@SP\nA=M\nM=D\n@SP\nM=M+1"
	popDcommand  string = "@SP\nAM=M-1\nD=M"
)

func initCodeWriter(file *os.File) *CodeWriter {
	var filePath []string = strings.Split(file.Name(), "/")
	var fileName string = strings.Split(filePath[len(filePath)-1], ".")[0]
	return &CodeWriter{file: file, comparisonCount: 0, callCount: 0, currentFileName: fileName, currentFunction: "main"}
}

// "Informs the code writer that the translation of a new VM file is started."
func (codewriter *CodeWriter) setFileName(filename string) {
	codewriter.currentFileName = filename
}

// Returns "the assembly code that is the translation of the given arithmetic command."
func (codewriter *CodeWriter) writeArithmetic(command string) string {
	var assembly string
	var incrementSP string = "@SP\nM=M+1\n"
	var binarycommand string = popDcommand + "\n@SP\nAM=M-1\n"
	var count string = strconv.Itoa(codewriter.comparisonCount)
	var booleancommand1 string = popDcommand + "\n@SP\nAM=M-1\nD=M-D\n" + "@BOOLEAN_" + count + "\n"
	var booleancommand2 string = "\nD=0\n@FINAL_BOOLEAN_" + count + "\n0;JEQ\n(BOOLEAN_" + count
	booleancommand2 += ")\nD=-1\n(FINAL_BOOLEAN_" + count + ")\n" + pushDcommand + "\n"
	codewriter.comparisonCount = codewriter.comparisonCount + 1

	switch command {
	case "add":
		assembly = "// add\n" + binarycommand + "M=D+M\n" + incrementSP
	case "sub":
		assembly = "// sub\n" + binarycommand + "M=M-D\n" + incrementSP
	case "and":
		assembly = "// and\n" + binarycommand + "M=D&M\n" + "M=D\n" + incrementSP
	case "or":
		assembly = "// or\n" + binarycommand + "M=D|M\n" + "M=D\n" + incrementSP
	case "not":
		assembly = "// not\n" + popDcommand + "\nM=!M\n" + incrementSP
	case "neg":
		assembly = "// neg\n" + popDcommand + "\nM=-M\n" + incrementSP
	case "eq":
		assembly = "// eq\n" + booleancommand1 + "D;JEQ" + booleancommand2
	case "gt":
		assembly = "// gt\n" + booleancommand1 + "D;JGT" + booleancommand2
	case "lt":
		assembly = "// lt\n" + booleancommand1 + "D;JLT" + booleancommand2
	}
	return assembly
}

// Returns "the assembly code that is the translation of the given command, where command is C_PUSH."
func (codewriter *CodeWriter) writePush(segment string, index int) string {
	var assembly string = "// C_PUSH " + segment + "[" + strconv.Itoa(index) + "]\n"
	var indexString string = strconv.Itoa(index)
	var memoryIndex string = "@" + indexString
	var pushString string = "\nD=M\n" + pushDcommand + "\n"

	switch segment {
	case "local":
		assembly += memoryIndex + "\nD=A\n@LCL\nA=D+M" + pushString
	case "argument":
		assembly += memoryIndex + "\nD=A\n@ARG\nA=D+M" + pushString
	case "this":
		assembly += memoryIndex + "\nD=A\n@THIS\nA=D+M" + pushString
	case "that":
		assembly += memoryIndex + "\nD=A\n@THAT\nA=D+M" + pushString
	case "pointer":
		if index == 0 {
			assembly += "@THIS" + pushString
		} else {
			assembly += "@THAT" + pushString
		}
	case "static":
		assembly += "@" + codewriter.currentFileName + "_" + indexString + pushString
	case "temp":
		assembly += memoryIndex + "\nD=A\n@R5\nA=D+A" + pushString
	case "constant":
		assembly += memoryIndex + "\nD=A\n" + pushDcommand + "\n"
	}
	return assembly
}

// Returns "the assembly code that is the translation of the given command, where command is C_POP."
func (codewriter *CodeWriter) writePop(segment string, index int) string {
	var assembly string = "// C_POP " + segment + "[" + strconv.Itoa(index) + "]\n"
	var indexString string = strconv.Itoa(index)
	var memoryIndex string = "@" + indexString
	var popString string = "\n@FRAME\nM=D\n" + popDcommand + "\n@FRAME\nA=M\nM=D\n"

	switch segment {
	case "local":
		assembly += memoryIndex + "\nD=A\n@LCL\nD=D+M" + popString
	case "argument":
		assembly += memoryIndex + "\nD=A\n@ARG\nD=D+M" + popString
	case "this":
		assembly += memoryIndex + "\nD=A\n@THIS\nD=D+M" + popString
	case "that":
		assembly += memoryIndex + "\nD=A\n@THAT\nD=D+M" + popString
	case "pointer":
		if index == 0 {
			assembly += popDcommand + "\n@THIS\nM=D\n"
		} else {
			assembly += popDcommand + "\n@THAT\nM=D\n"
		}
	case "static":
		assembly += popDcommand + "\n@" + codewriter.currentFileName + "_" + indexString + "\nM=D\n"
	case "temp":
		assembly += memoryIndex + "\nD=A\n@R5\nD=D+A" + popString
	}
	return assembly
}

// Returns "assembly code that effects the label command."
func (codewriter *CodeWriter) writeLabel(label string) string {
	var assembly string = "// label " + label + "\n"
	assembly += "(" + label + ")\n"
	return assembly
}

// Returns "assembly code that effects the goto command"
func (codewriter *CodeWriter) writeGoto(label string) string {
	var assembly string = "// goto " + label + "\n"
	assembly += "@" + label + "\n0;JMP\n"
	return assembly
}

// Returns "assembly code that effects the if-goto command"
func (codewriter *CodeWriter) writeIf(label string) string {
	var assembly string = "// if-goto " + label + "\n"

	// "The stack's topmost value is popped"
	assembly += popDcommand + "\n"
	/* "If the value is not zero, execution continues from the location marked by the label; otherwise, execution continues
	from the  next command in the program. The jump destination must be located in the same function."*/
	assembly += "@" + label + "\nD;JNE\n"
	return assembly
}

// Returns "assembly code that effects the call command"
func (codewriter *CodeWriter) writeCall(functionName string, numArgs int) string {
	var assembly string = "// call " + functionName + " " + strconv.Itoa(numArgs) + "\n"
	var return_address = "return_address_" + strconv.Itoa(codewriter.callCount)
	var pushString string = "\nD=M\n" + pushDcommand + "\n"

	// Step 1: "push return-address, (Using the label declared below)"
	assembly += "@" + return_address + "\nD=A\n" + pushDcommand + "\n"
	// Step 2: "push LCL, Save LCL of the calling function"
	assembly += "@LCL" + pushString
	// Step 3: "push ARG, Save ARG of the calling function"
	assembly += "@ARG" + pushString
	// Step 4: "push THIS, Save THIS of the calling function"
	assembly += "@THIS" + pushString
	// Step 5: "push THAT, Save THAT of the calling function"
	assembly += "@THAT" + pushString
	// Step 6: "ARG = SP-n-5, Reposition ARG (n = number of args.)"
	assembly += "\nD=M\n@" + strconv.Itoa(numArgs+5) + "\nD=D-A\n@ARG\nM=D\n"
	// Step 7: "LCL = SP, Reposition LCL"
	assembly += "\n@SP\nD=M\n@LCL\nM=D\n"
	// Step 8: "goto f, Transfer Control"
	assembly += "@" + functionName + "\n0;JMP\n"
	// Step 9: "(return-address), Declare a label for the return-address"
	assembly += "(" + return_address + ")\n"

	codewriter.callCount = codewriter.callCount + 1
	return assembly
}

/* Returns "assembly code that effects the VM initialization, also called bootstrap code.
This code must be placed at the beginning of the output file."*/
func (codewriter *CodeWriter) writeInit() string {
	var assembly string = "// Boostrap code\n@256\nD=A\n@SP\nM=D\n" // "Initialize the stack pointer to 0x0000"
	assembly += codewriter.writeCall("Sys.init", 0) + "0;JMP\n"     // "Start executing (the translated code of) Sys.init"
	return assembly
}

// Returns "assembly code that effects the return command"
func (codewriter *CodeWriter) writeReturn() string {
	var assembly string = "// return\n"
	var FRAME string = "\n@FRAME\n"
	var RET string = "\n@RET\n"

	// Step 1: "FRAME = LCL, FRAME is a temporary variable"
	assembly += "@LCL\nD=M" + FRAME + "M=D"
	// Step 2: "RET = *(FRAME-5)m Put the return-address in a temp. var."
	assembly += FRAME + "D=M\n@5\nD=D-A\nA=D\nD=M" + RET + "M=D\n"
	// Step 3: "*ARG = pop(), Reposition the return value for the caller"
	assembly += popDcommand + "\n@ARG\nA=M\nM=D\n"
	// Step 4: "SP = ARG + 1, Restore SP of the caller"
	assembly += "@ARG\nD=M+1\n@SP\nM=D"
	// Step 5: "THAT = *(FRAME-1), Restore THAT of the caller"
	assembly += FRAME + "D=M\n@1\nD=D-A\nA=D\nD=M\n@THAT\nM=D"
	// Step 6: "THIS = *(FRAME-2), Restore THIS of the caller"
	assembly += FRAME + "D=M\n@2\nD=D-A\nA=D\nD=M\n@THIS\nM=D"
	// Step 7: "ARG = *(FRAME-3), Restore ARG of the caller"
	assembly += FRAME + "D=M\n@3\nD=D-A\nA=D\nD=M\n@ARG\nM=D"
	// Step 8: "LCL = *(FRAME-4), Restore LCL of the caller"
	assembly += FRAME + "D=M\n@4\nD=D-A\nA=D\nD=M\n@LCL\nM=D"
	// Step 9: "goto RET, GOto return-address (in the caller's code)"
	assembly += RET + "A=M\n0;JMP\n"
	return assembly
}

// Returns "assembly code that effects the function command"
func (codewriter *CodeWriter) writeFunction(functionName string, numLocals int) string {
	codewriter.currentFunction = functionName
	var assembly string = "// function " + functionName + strconv.Itoa(numLocals) + "\n"

	// Step 1: "(f), Declare a label for the function entry"
	assembly += "(" + functionName + ")\n@SP\nA=M\n"

	// Step 2: "repeat k times:, k = number of local variables"
	for i := 0; i < numLocals; i++ {
		// Step 3: "push 0, Initialize all of them to 0"
		assembly += "M=0\nA=A+1\n"
	}
	assembly += "D=A\n@SP\nM=D\n"
	return assembly
}

func generateASM(outputFileName string, writer *bufio.Writer, parser *Parser, codewriter *CodeWriter) {
	parser.scanner = bufio.NewScanner(parser.file)
	var assembly string

	for parser.hasMoreCommands() {
		parser.advance()
		switch parser.commandType() {
		case C_ARITHMETIC:
			assembly = codewriter.writeArithmetic(parser.arg1())
		case C_PUSH:
			assembly = codewriter.writePush(parser.arg1(), parser.arg2())
		case C_POP:
			assembly = codewriter.writePop(parser.arg1(), parser.arg2())
		case C_LABEL:
			assembly = codewriter.writeLabel(parser.arg1())
		case C_GOTO:
			assembly = codewriter.writeGoto(parser.arg1())
		case C_IF:
			assembly = codewriter.writeIf(parser.arg1())
		case C_FUNCTION:
			assembly = codewriter.writeFunction(parser.arg1(), parser.arg2())
		case C_RETURN:
			assembly = codewriter.writeReturn()
		case C_CALL:
			assembly = codewriter.writeCall(parser.arg1(), parser.arg2())
		}
		writer.WriteString(assembly)
	}
}

func main() {
	filePath := os.Args[1]
	var parser *Parser
	var codewriter *CodeWriter
	var fileName string = strings.TrimSuffix(filePath, ".vm")
	var outputFileName string

	if fileName != filePath { // Case: the input is a .vm file
		outputFileName = fileName + ".asm"
		file, err1 := os.Open(filePath)
		if err1 != nil {
			fmt.Println(err1)
		}
		defer file.Close()

		parser = initParser(file)
		codewriter = initCodeWriter(file)
		output, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(outputFileName, "successfully created.")
		}
		defer output.Close()

		var writer *bufio.Writer = bufio.NewWriter(output)
		defer writer.Flush()
		generateASM(outputFileName, writer, parser, codewriter)
	} else { // Case: the input is a directory
		files, err2 := ioutil.ReadDir(fileName) // os.ReadDir not supported in Coursera
		if err2 != nil {
			fmt.Println(err2)
		}

		outputFileName = filepath.Join(fileName, filepath.Base(fileName)) + ".asm"
		output, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(outputFileName, "successfully created.")
		}
		defer output.Close()
		var writer *bufio.Writer = bufio.NewWriter(output)
		defer writer.Flush()
		var init bool = false
		for _, file := range files {
			if filepath.Ext(filepath.Join(fileName, file.Name())) == ".vm" {
				file, err3 := os.Open(filepath.Join(fileName, file.Name()))
				if err3 != nil {
					fmt.Println(err3)
				}
				defer file.Close()
				parser = initParser(file)
				codewriter = initCodeWriter(file)

				// generate boostrap code
				if init == false {
					var assembly string = codewriter.writeInit()
					writer.WriteString(assembly)
					init = true
				}

				generateASM(outputFileName, writer, parser, codewriter)
			}
		}
	}
}
