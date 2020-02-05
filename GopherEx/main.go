package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//Token bla bla bla
type Token struct {
	prefix []string
	array  []Token
}

//Tokenizer bla bla bla
type Tokenizer struct {
	code  []string
	index int
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
	fmt.Println("[ERROR] TRIED TO ACCES TOKEN BEYOND EOF")
	os.Exit(38) //ERROR_HANDLE_EOF
	return ""
}

func main() {
	b, err := ioutil.ReadFile("file.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	sourceCode := string(b)
	cleanCode := sanitize(sourceCode)
	parsedCode := parseCode(cleanCode)
	fmt.Println("")
	printTokens(tokenizeCode(parsedCode), 0)
	fmt.Println("")
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

	return tempCode
}

func tokenizeCode(code []string) []Token {

	tokenizer := Tokenizer{code, -1}
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
			tempToken.prefix = append(tempToken.prefix, strings.Split(element, " ")...)
			tempToken.array = tokenizeCode(code[tokenizer.index+2 : i+tokenizer.index]) //CALL FUNCTION

			tokenizer.index += i              //SKIP BEYOND FUNCTION
			if tokenizer.index >= len(code) { //NO MORE TOKENS RETURN LAST
				tempToken.prefix = nil
				tempToken.prefix = append(tempToken.prefix, strings.Split(element, " ")...)
				tokens = append(tokens, tempToken)
				return tokens
			}

			if string(code[tokenizer.index][0]) == "(" {

				tempToken.prefix = append(tempToken.prefix, append(append(make([]string, 0), "()"), strings.Split(code[tokenizer.index][1:len(code[tokenizer.index])-1], ",")...)...) //REMOVE ( and ), AND SPLIT STRING WITH ","
			}
			tokens = append(tokens, tempToken)
			tempToken.prefix = nil
			tokenizer.index-- //GO ONE BACK FOR NEXT LOOP
		} else {
			//MAKE SURE THAT IT IS IN A FUNCTION
			if element != "]" && element != "[" && string(element[0]) != "(" {
				if len(tempToken.prefix) > 0 {
					if tempToken.prefix[0] == "T" { // T FOR TEXT
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
		fmt.Println("├" + strings.Repeat("─", liftOff*4) + " " + strings.Join(element.prefix, " "))
		if len(element.array) > 0 {
			liftOff++
			printTokens(element.array, liftOff)
			liftOff--
		}
	}
}

func printCode(code []string) {
	for _, _string := range code {
		fmt.Println(_string)
	}
}
