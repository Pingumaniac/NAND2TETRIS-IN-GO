package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// General description: "Translates Hack assembly language mnemonics into binary codes."
var code_dest = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var code_comp = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"M":   "1110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"!M":  "1110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"-M":  "1110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"M+1": "1110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"M-1": "1110010",
	"D+A": "0000010",
	"D+M": "1000010",
	"D-A": "0010011",
	"D-M": "1010011",
	"A-D": "0000111",
	"M-D": "1000111",
	"D&A": "0000000",
	"D&M": "1000000",
	"D|A": "0010101",
	"D|M": "1010101",
}

var code_jump = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

// Returns the binary code of the dest mnemonic
func dest(mnemonic string) string {
	return code_dest[mnemonic]
}

// Returns the binary code of the comp mnemonic
func comp(mnemonic string) string {
	return code_comp[mnemonic]
}

// Returns the binary code of the jump mnemonic
func jump(mnemonic string) string {
	return code_jump[mnemonic]
}

// General description:  "Keeps a correspondence between symbolic labels and numeric addresses."
type SymbolTable struct {
	symbols map[string]int
}

func initSymbolTable() *SymbolTable {
	table := new(SymbolTable)
	table.symbols = map[string]int{
		// The Hack language features 23 predefined symbols:
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
	}
	return table
}

// "Adds the pair (symbol, address) to the SymbolTable
func (table *SymbolTable) addEntry(symbol string, address int) {
	table.symbols[symbol] = address
}

// "Does the symbol table contain the given symbol?"
func (table *SymbolTable) contains(symbol string) bool {
	_, ok := table.symbols[symbol]
	return ok
}

// "Returns the address associated with the symbols"
func (table *SymbolTable) GetAddress(symbol string) int {
	return table.symbols[symbol]
}

/* General description: "Encapsulates access to the input code.
Reads an assembly language command, parses it, and provides convenient access to command's components (fields and symbols).
Removes all white space and comments" */
const (
	A_COMMAND = 0
	C_COMMAND = 1
	L_COMMAND = 2
)

type Parser struct {
	file           *os.File
	scanner        *bufio.Scanner
	currentCommand string
	ramAddress     int
}

func initParser(file *os.File) *Parser {
	scanner := bufio.NewScanner(file)
	var parser Parser = Parser{file: file, scanner: scanner, ramAddress: 0}
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
Should be called only if hasMoreCommands() is true.
Initially there is no curent command."
*/
func (parser *Parser) advance() {
	var inputCommand string = parser.scanner.Text() // Read next command
	var actualCommand string = strings.Split(inputCommand, "//")[0]
	parser.currentCommand = strings.TrimSpace(actualCommand)
}

/* "Returns the type of the current command:
A_COMMAND for @Xxx whjere Xxx is either a symbol or a decimal number
C_COMMAND for dest=comp;jump
L_COMMAND (actually, pseudo-command) for (Xxx) where Xxx is a symbol"
*/
func (parser *Parser) commandType() int {
	if strings.HasPrefix(parser.currentCommand, "@") {
		return A_COMMAND
	} else if strings.HasPrefix(parser.currentCommand, "(") && strings.HasSuffix(parser.currentCommand, ")") {
		return L_COMMAND
	} else {
		return C_COMMAND
	}
}

/* "Returns the symbol or decimal Xxx of the current command @Xxx or (Xxx).
Should be called only when commandType() is A_COMMAND or L_COMMAND."
*/
func (parser *Parser) symbol() string {
	var commandtype int = parser.commandType()
	if commandtype == A_COMMAND {
		var val string = strings.Replace(parser.currentCommand, "@", "", 1)
		return val
	} else if commandtype == L_COMMAND {
		var val string = strings.Replace(parser.currentCommand, "(", "", 1)
		val = strings.Replace(val, ")", "", 1)
		return val
	} else { // C_COMMAND does not consist of a symbol
		return parser.currentCommand
	}

}

/* "Returns the dest mnemonic in the current C-command (8 possibilities)
Should be called only when commandType is C_COMMAND."
*/
func (parser *Parser) dest() string {
	if strings.Contains(parser.currentCommand, "=") {
		return strings.Split(parser.currentCommand, "=")[0]
	} else {
		return ""
	}
}

/* "Returns the comp mnemonic in the current C-command (28 possibilities).
Should be called only when commandType() is C_COMMAND."
*/
func (parser *Parser) comp() string {
	if strings.Contains(parser.currentCommand, "=") {
		return strings.Split(parser.currentCommand, "=")[1]
	} else {
		return strings.Split(parser.currentCommand, ";")[0]
	}
}

/* "Returns the jump mnemonic in the current C-Command (8 possiblities).
Should be called only when commandType() is C_COMMAND."
*/
func (parser *Parser) jump() string {
	if strings.Contains(parser.currentCommand, ";") {
		return strings.Split(parser.currentCommand, ";")[1]
	} else {
		return ""
	}
}

func decimalToBinary(decimal int) int {
	var binary int = 0
	var counter int = 1
	var remainder int = 0

	for decimal != 0 {
		remainder = decimal % 2
		decimal = decimal / 2
		binary = binary + remainder*counter
		counter = counter * 10
	}
	return binary
}

func addLCOMMAND(parser *Parser, symboltable *SymbolTable) *SymbolTable {
	for parser.hasMoreCommands() {
		parser.advance()
		if parser.commandType() == L_COMMAND {
			symboltable.addEntry(parser.symbol(), parser.ramAddress)
		} else {
			parser.ramAddress = parser.ramAddress + 1
		}
	}
	return symboltable
}

func generateHack(filepath string, parser *Parser, symboltable *SymbolTable) {
	output, err3 := os.Create(strings.TrimSuffix(filepath, ".asm") + ".hack")
	if err3 != nil {
		fmt.Println(err3)
	}
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	_, err4 := parser.file.Seek(0, 0)
	if err4 != nil {
		fmt.Println(err4)
	}
	parser.scanner = bufio.NewScanner(parser.file)
	parser.currentCommand = ""
	parser.ramAddress = 16

	var hackCommand string
	for parser.hasMoreCommands() {
		parser.advance()
		if parser.commandType() == A_COMMAND {
			symbol := parser.symbol()
			address, err5 := strconv.Atoi(symbol)
			if err5 != nil {
				if symboltable.contains(symbol) == false {
					symboltable.addEntry(symbol, parser.ramAddress)
					address = parser.ramAddress
					parser.ramAddress = parser.ramAddress + 1
				} else {
					address = symboltable.GetAddress(symbol)
				}
			}

			hackCommand = strconv.Itoa(decimalToBinary(address)) // convert address from decimal to binary
			for len(hackCommand) < 16 {
				hackCommand = "0" + hackCommand
			}
			hackCommand = hackCommand + "\n"
			writer.WriteString(hackCommand)
		}
		if parser.commandType() == C_COMMAND {
			hackCommand = "111" + comp(parser.comp()) + dest(parser.dest()) + jump(parser.jump()) + "\n"
			writer.WriteString(hackCommand)
		}
	}
	fmt.Println(strings.TrimSuffix(filepath, ".asm") + ".hack successfully created.")
}

func getChoice() bool {
	fmt.Println("Would you like to compile another file? Type 'y' for Yes, 'n' for No.")
	fmt.Print(">")
	var choice string
	_, err3 := fmt.Scanln(&choice)
	if err3 != nil {
		fmt.Println(err3)
		return false
	}
	switch strings.ToLower(choice) {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Println("Invalid input.")
		return false
	}
}

func main() {
	var symboltable *SymbolTable = initSymbolTable()

	for {
		fmt.Println("Enter the path of the .asm file to compile into .hack file:")
		fmt.Print(">")
		var filepath string
		_, err1 := fmt.Scanln(&filepath)
		if err1 != nil {
			fmt.Println(err1)
		}

		file, err2 := os.Open(filepath)
		if err2 != nil {
			fmt.Println(err2)
		}
		defer file.Close()

		var parser *Parser = initParser(file)
		symboltable = addLCOMMAND(parser, symboltable)
		generateHack(filepath, parser, symboltable)
		var choice bool = getChoice()
		if choice {
			continue
		} else {
			break
		}
	}
}
