package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//Token bla bla bla
type Token struct {
	identifiers []string
	tokens      []Token
}

const (
	ERROR_OUT_OF_BOUNDS = iota + 1
	ERROR_WHITESPACE_IN_NAME
	ERROR_UNIDENTIFIED_TYPE
	ERROR_WRONG_ARGUMENT_FORMAT
	ERROR_FUNCTION_ALREADY_EXSITS
	ERROR_FUNCTION_DOESNT_EXIST
	ERROR_WRONG_TYPE_FORMAT
)

//Tokenizer bla bla bla
type Tokenizer struct {
	code      []string
	tokens    []Token
	functions map[string]bool
	index     int
	typeBytes map[string]bool
}

func (tokenizer Tokenizer) peak() string {
	if len(tokenizer.code) > tokenizer.index+1 {
		return tokenizer.code[tokenizer.index+1]
	}
	return ""
}

func (tokenizer *Tokenizer) next() string {
	if len(tokenizer.code) > tokenizer.index+1 {
		tokenizer.index++
		return tokenizer.code[tokenizer.index] //Return token
	}
	tokenizer.error(ERROR_OUT_OF_BOUNDS)
	return ""
}

func (tokenizer *Tokenizer) tokenizeCode() Token {
	//Is function
	if len(tokenizer.code)-1 < tokenizer.index {
		return Token{}
	}

	if tokenizer.isFunction() {
		//tokenize function
		functionToken := tokenizer.tokenizeFunction()
		tokenizer.next()
		tokenizer.next()
		//go to next line and Tokenize the body of the function
		for true {
			if tokenizer.code[tokenizer.index] == "]" {
				break
			}
			functionToken.tokens = append(functionToken.tokens, tokenizer.tokenizeCode())
		}
		if len(tokenizer.code)-1 > tokenizer.index {
			tokenizer.next()
		}
		return functionToken
	} else if tokenizer.isTypeByte() {
		token := tokenizer.tokenizeType()
		tokenizer.next()
		return token

	} else if tokenizer.isFunctionCall() {
		token := tokenizer.tokenizeFunctionCall()
		tokenizer.next()
		return token
	} else {
		tokenizer.error(ERROR_UNIDENTIFIED_TYPE)
	}

	return Token{}
}

//Type Check
func (tokenizer *Tokenizer) isFunction() bool {
	return tokenizer.peak() == "["
}
func (tokenizer *Tokenizer) isFunctionCall() bool {
	functionName, _ := tokenizer.getFunctionIdentifiers()
	//Check if function exist
	if tokenizer.functions[functionName] {
		return tokenizer.functions[functionName]
	}
	//Error if function does not exist
	tokenizer.error(ERROR_FUNCTION_DOESNT_EXIST)
	return false
}
func (tokenizer *Tokenizer) isTypeByte() bool {

	currentLine := tokenizer.code[tokenizer.index]
	if len(currentLine) < 2 {
		tokenizer.error(ERROR_WRONG_TYPE_FORMAT)
	}
	typeByte := currentLine[0]
	//Checks if the the length of the string is greater than 1 and validates that it's a single byte
	validatorByte := (map[bool]bool{true: (currentLine[1] == ' '), false: false})[len(currentLine) > 1]
	//Checks if it's a valid typeByte
	if tokenizer.typeBytes[string(typeByte)] != false && validatorByte == true {
		return true
	}
	return false
}

//Tokenizers
func (tokenizer *Tokenizer) tokenizeFunction() Token {
	functionName, UnformatedArguments := tokenizer.getFunctionIdentifiers()
	//Check if function allready exists
	if tokenizer.functionNameExists(functionName) {
		tokenizer.error(ERROR_FUNCTION_ALREADY_EXSITS) //Exits with error
	}
	//appends function to function list
	tokenizer.functions[functionName] = true
	argumentNames := strings.Split(strings.Replace(UnformatedArguments, " ", "", -1), ",") //...
	//Clean argument names + Check for whitespaces in argument names and whitespaces in function name
	argumentNames = tokenizer.cleanArgumentNames(argumentNames)
	tokenizer.whitespaceInName(functionName)
	//Creates the function identifiers
	var indentifiers []string
	indentifiers = append(indentifiers, functionName)
	//First element is the function name
	indentifiers = append(indentifiers, argumentNames...)
	//Function arguments [arg1,arg2...]
	token := Token{indentifiers, make([]Token, 0)}
	return token
}

func (tokenizer *Tokenizer) tokenizeType() Token {
	currentLine := tokenizer.code[tokenizer.index]
	typeByte := currentLine[0]
	var indentifiers []string
	switch typeByte {
	case 'T':
		indentifiers = append(indentifiers, string(typeByte)) //Append type
		indentifiers = append(indentifiers, currentLine[2:])  //currentLine[2:] =T |hello
		break

	default:
		indentifiers = append(indentifiers, string(typeByte))                       //Append type
		indentifiers = append(indentifiers, strings.Split(currentLine[2:], " ")...) //Append Arguments
	}
	token := Token{indentifiers, make([]Token, 0)}
	return token
}

func (tokenizer *Tokenizer) tokenizeFunctionCall() Token {
	functionName, UnformatedArguments := tokenizer.getFunctionIdentifiers()
	//Check if function exist and if it dosent error out
	// '\"' ]
	//ARRAY
	//TEMPSTR
	//","
	if !tokenizer.functionNameExists(functionName) {
		tokenizer.error(ERROR_FUNCTION_DOESNT_EXIST)
	}
	var arguments []string
	arguments = append(arguments, functionName)
	argument := ""
	for _, char := range UnformatedArguments {
		if char == ',' {
			arguments = append(arguments, argument)
			argument = ""
		} else {
			argument = argument + string(char)
		}
	}
	if argument != "" {
		arguments = append(arguments, argument)
	}

	token := Token{arguments, make([]Token, 0)}
	return token
}

//Tokenize

// Helper functions
func (tokenizer *Tokenizer) error(errorCode int) {
	switch errorCode {
	case ERROR_OUT_OF_BOUNDS:
		fmt.Println("ERROR_OUT_OF_BOUNDS:", "INDEX", tokenizer.index, "OUT OF", len(tokenizer.code)-1)
		os.Exit(errorCode)
		break

	case ERROR_UNIDENTIFIED_TYPE:
		fmt.Println("ERROR_UNIDENTIFIED_TYPE:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WHITESPACE_IN_NAME:
		fmt.Println("ERROR_WHITESPACE_IN_NAME:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WRONG_ARGUMENT_FORMAT:
		fmt.Println("ERROR_WRONG_ARGUMENT_FORMAT:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_FUNCTION_ALREADY_EXSITS:
		fmt.Println("ERROR_FUNCTION_ALREADY_EXSITS:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break

	case ERROR_FUNCTION_DOESNT_EXIST:
		fmt.Println("ERROR_FUNCTION_DOESNT_EXIST:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WRONG_TYPE_FORMAT:
		fmt.Println("ERROR_WRONG_TYPE_FORMAT:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break

	default:
		fmt.Println("ERROR_UNKOWN_ERROR:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	}
}

func (tokenizer *Tokenizer) whitespaceInName(name string) {
	isWhiteSpace := func(c rune) bool {
		return c == '\t' || c == '\r' || c == ' '
	}
	if strings.IndexFunc(name, isWhiteSpace) != -1 {
		tokenizer.error(ERROR_WHITESPACE_IN_NAME)
	}

}

func (tokenizer *Tokenizer) functionNameExists(functionName string) bool {
	return tokenizer.functions[functionName]
}

func (tokenizer *Tokenizer) cleanArgumentNames(names []string) []string {
	var cleanNames []string
	for _, name := range names {
		cleanNames = append(cleanNames, strings.TrimSpace(name))
		tokenizer.whitespaceInName(cleanNames[len(cleanNames)-1]) //Stops tokinization if true and throws error
	}
	return cleanNames
}

func (tokenizer *Tokenizer) getFunctionIdentifiers() (string, string) {
	currentLine := tokenizer.code[tokenizer.index]
	UnformatedIdentifiers := strings.Split(currentLine, "(") //box(arg1,arg2) =[box, argstring]
	functionName := strings.TrimSpace(UnformatedIdentifiers[0])
	arguments := strings.Split(UnformatedIdentifiers[len(UnformatedIdentifiers)-1], ")")
	return functionName, arguments[0]
}

func main() {
	b, err := ioutil.ReadFile("file.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	sourceCode := string(b)
	cleanCode := sanitize(sourceCode)
	parsedCode := parseCode(cleanCode)

	tokenizer := Tokenizer{parsedCode, make([]Token, 0), map[string]bool{}, 0, map[string]bool{"T": true}}
	tokens := make([]Token, 0)
	for tokenizer.code[tokenizer.index] != "." {
		tokens = append(tokens, tokenizer.tokenizeCode())
	}
	printTokens(tokens, 0)
}

func sanitize(str string) []string {
	tempCode := strings.Split(str, "\n")

	for index := range tempCode {
		tempCode[index] = strings.Replace(tempCode[index], string(rune(13)), "", -1) // REMOVE 13/CR

		if tempCode[index] != "" {
			i := 0
			for string(tempCode[index][i]) == " " {
				i++
				if i >= len(tempCode[index]) {
					break
				}
			}

			tempCode[index] = tempCode[index][i:]
		}
	}

	var array []string
	for _, _string := range tempCode {
		if _string != "" {
			array = append(array, _string)
		}

	}
	return array
}

func parseCode(code []string) []string {
	breaks := make(map[string]int)
	breaks["["] = 1 //"()"
	breaks["]"] = 1
	tokenString := ""
	var tempCode []string
	for _, _string := range code {
		for _, char := range _string { //"a,"
			if breaks[string(char)] == 1 {
				if tokenString != "" {
					tempCode = append(tempCode, tokenString)
				}
				tempCode = append(tempCode, string(char))
				tokenString = ""
			} else {
				tokenString += string(char)
			}
		}
		if tokenString != "" {
			tempCode = append(tempCode, tokenString)
			tokenString = ""
		}
	}
	if tokenString != "" {
		tempCode = append(tempCode, tokenString)
	}
	tempCode = append(tempCode, ".") //EOF
	return tempCode
}

func tokenizeCode(code []string) []Token {

	tokenizer := Tokenizer{code, make([]Token, 0), map[string]bool{}, -1, map[string]bool{}}
	tokens := make([]Token, 0)
	tempToken := Token{make([]string, 0), make([]Token, 0)}

	for tokenizer.index < len(code)-1 {
		element := tokenizer.next()

		if tokenizer.peak() == "[" {
			brackLeft := 0
			brackRight := -1
			i := 2
			for brackLeft != brackRight {
				if code[i+tokenizer.index] == "[" {
					brackLeft++
				}
				if code[i+tokenizer.index] == "]" {
					brackRight++
				}
				i++
			}
			tempToken.identifiers = strings.Split(element, " ")
			tempToken.tokens = tokenizeCode(code[tokenizer.index+2 : i+tokenizer.index]) //CALL FUNCTION

			tokenizer.index += i              //SKIP BEYOND FUNCTION
			if tokenizer.index >= len(code) { //NO MORE TOKENS RETURN LAST
				tempToken.identifiers = nil
				tempToken.identifiers = append(tempToken.identifiers, strings.Split(element, " ")...)
				tokens = append(tokens, tempToken)
				return tokens
			}

			if string(code[tokenizer.index][0]) == "(" { //I

				tempToken.identifiers = append(tempToken.identifiers, append(append(make([]string, 0), "()"), strings.Split(code[tokenizer.index][1:len(code[tokenizer.index])-1], ",")...)...) //REMOVE ( and ), AND SPLIT STRING WITH ","
			}
			tokens = append(tokens, tempToken)
			tempToken.identifiers = nil
			tokenizer.index-- //GO ONE BACK FOR NEXT LOOP
		} else {
			//MAKE SURE THAT IT IS IN A FUNCTION
			if element != "]" && element != "[" && string(element[0]) != "(" {
				if len(tempToken.identifiers) > 0 {
					if tempToken.identifiers[0] == "T" { // T FOR TEXT
						tokens = append(tokens, Token{append(make([]string, 0), element), make([]Token, 0)}) //TEXT DONT SPLIT INTO ARRAY
					} else {
						tokens = append(tokens, Token{strings.Split(element, " "), make([]Token, 0)}) //NOT TEXT SPLIT INTO ARRAY
					}
				} else {
					if strings.Split(element, " ")[0] == "T" {
						tokens = append(tokens, Token{append(make([]string, 0), element), make([]Token, 0)}) //TEXT DONT SPLIT INTO ARRAY
					} else {
						tokens = append(tokens, Token{strings.Split(element, " "), make([]Token, 0)}) //IF NOT IN FUNCTION SPLIT INTO ARRAY}
					}
				}

			}
		}
	}

	return tokens
}

func printTokens(tokens []Token, liftOff int) {
	for _, element := range tokens {
		fmt.Println("├" + strings.Repeat("─", liftOff*4) + " " + strings.Join(element.identifiers, " "))
		if len(element.tokens) > 0 { //Check if token has sub tokens
			printTokens(element.tokens, liftOff+1)
		}
	}
}

func printCode(code []string) {
	for _, _string := range code {
		fmt.Println(_string)
	}
}
