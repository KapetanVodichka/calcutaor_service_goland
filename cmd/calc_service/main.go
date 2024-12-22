package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

var operatorPrecedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

type requestBody struct {
	Expression string `json:"expression"`
}

type responseSuccess struct {
	Result string `json:"result"`
}

type responseError struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("Сервер запущен на порту :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	result, calcErr := Calc(req.Expression)
	if calcErr != nil {
		if isExpressionError(calcErr) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(responseError{Error: "Expression is not valid"})
			return
		}
		log.Println("Internal error:", calcErr)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseSuccess{Result: fmt.Sprintf("%v", result)})
}

func isExpressionError(err error) bool {
	switch err.Error() {
	case "деление на ноль",
		"некорректное выражение",
		"недопустимый символ",
		"несбалансированные скобки",
		"недостаточно операндов для операции":
		return true
	}
	if strings.Contains(err.Error(), "недопустимый символ") ||
		strings.Contains(err.Error(), "несбалансированные скобки") ||
		strings.Contains(err.Error(), "деление на ноль") ||
		strings.Contains(err.Error(), "некорректное выражение") {
		return true
	}
	return false
}

func Calc(expression string) (float64, error) {
	tokens, err := splitExpression(expression)
	if err != nil {
		return 0, err
	}
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	result, err := computePostfix(postfix)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func splitExpression(expr string) ([]string, error) {
	var tokens []string
	var numberBuilder strings.Builder
	for i, ch := range expr {
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsDigit(ch) || ch == '.' {
			numberBuilder.WriteRune(ch)
		} else {
			if numberBuilder.Len() > 0 {
				tokens = append(tokens, numberBuilder.String())
				numberBuilder.Reset()
			}
			if strings.ContainsRune("+-*/()", ch) {
				tokens = append(tokens, string(ch))
			} else {
				return nil, fmt.Errorf("недопустимый символ '%c' на позиции %d", ch, i)
			}
		}
	}
	if numberBuilder.Len() > 0 {
		tokens = append(tokens, numberBuilder.String())
	}
	return tokens, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []string
	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if isOperator(token) {
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) &&
				operatorPrecedence[stack[len(stack)-1]] >= operatorPrecedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, errors.New("несбалансированные скобки")
			}
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" || top == ")" {
			return nil, errors.New("несбалансированные скобки")
		}
		output = append(output, top)
	}
	return output, nil
}

func computePostfix(postfix []string) (float64, error) {
	var stack []float64
	for _, token := range postfix {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("некорректное число: %s", token)
			}
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("недостаточно операндов для операции")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				result = a / b
			}
			stack = append(stack, result)
		}
	}
	if len(stack) != 1 {
		return 0, errors.New("некорректное выражение")
	}
	return stack[0], nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	_, exists := operatorPrecedence[s]
	return exists
}
