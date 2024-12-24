package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSplitExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantTokens  []string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Простое выражение",
			input:      "2+2",
			wantTokens: []string{"2", "+", "2"},
		},
		{
			name:       "Скобки и пробелы",
			input:      " (  3 + 4 ) * 2 ",
			wantTokens: []string{"(", "3", "+", "4", ")", "*", "2"},
		},
		{
			name:        "Недопустимый символ",
			input:       "2+2^5",
			wantErr:     true,
			errContains: "недопустимый символ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokens, err := splitExpression(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("splitExpression(%q) error = %v, wantErr = %v",
					tt.input, err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.errContains) {
				t.Errorf("splitExpression(%q) error = %v, want substring %q",
					tt.input, err, tt.errContains)
			}
			if !tt.wantErr && !reflect.DeepEqual(gotTokens, tt.wantTokens) {
				t.Errorf("splitExpression(%q) = %v, want %v",
					tt.input, gotTokens, tt.wantTokens)
			}
		})
	}
}

func TestInfixToPostfix(t *testing.T) {
	tests := []struct {
		name        string
		tokens      []string
		want        []string
		wantErr     bool
		errContains string
	}{
		{
			name:   "Простое выражение",
			tokens: []string{"2", "+", "2"},
			want:   []string{"2", "2", "+"},
		},
		{
			name:   "Скобки",
			tokens: []string{"(", "3", "+", "4", ")", "*", "2"},
			want:   []string{"3", "4", "+", "2", "*"},
		},
		{
			name:        "Несбалансированные скобки",
			tokens:      []string{"(", "2", "+"},
			wantErr:     true,
			errContains: "несбалансированные скобки",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := infixToPostfix(tt.tokens)
			if (err != nil) != tt.wantErr {
				t.Fatalf("infixToPostfix(%v) error = %v, wantErr = %v",
					tt.tokens, err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.errContains) {
				t.Errorf("infixToPostfix(%v) error = %v, want substring %q",
					tt.tokens, err, tt.errContains)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("infixToPostfix(%v) = %v, want %v",
					tt.tokens, got, tt.want)
			}
		})
	}
}

func TestComputePostfix(t *testing.T) {
	tests := []struct {
		name    string
		postfix []string
		want    float64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "2 2 + -> 4",
			postfix: []string{"2", "2", "+"},
			want:    4,
		},
		{
			name:    "3 4 + 2 * -> 14",
			postfix: []string{"3", "4", "+", "2", "*"},
			want:    14,
		},
		{
			name:    "Деление на ноль",
			postfix: []string{"2", "0", "/"},
			wantErr: true,
			errMsg:  "деление на ноль",
		},
		{
			name:    "Недостаточно операндов",
			postfix: []string{"+", "2", "3"},
			wantErr: true,
			errMsg:  "недостаточно операндов",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computePostfix(tt.postfix)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("computePostfix(%v) expected error, got nil", tt.postfix)
				}
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("computePostfix(%v) error = %v, want substring %q",
						tt.postfix, err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("computePostfix(%v) got error %v, want nil", tt.postfix, err)
				}
				if got != tt.want {
					t.Errorf("computePostfix(%v) = %v, want %v",
						tt.postfix, got, tt.want)
				}
			}
		})
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		want        float64
		wantErr     bool
		errContains string
	}{
		{
			name:       "Простая арифметика",
			expression: "2+2",
			want:       4,
		},
		{
			name:       "Со скобками",
			expression: "(3+4)*2",
			want:       14,
		},
		{
			name:        "Деление на ноль",
			expression:  "2/0",
			wantErr:     true,
			errContains: "деление на ноль",
		},
		{
			name:        "Несбалансированные скобки",
			expression:  "2+(3*2",
			wantErr:     true,
			errContains: "несбалансированные скобки",
		},
		{
			name:        "Недопустимый символ",
			expression:  "2^3",
			wantErr:     true,
			errContains: "недопустимый символ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.expression)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Calc(%q) expected error, got nil", tt.expression)
				}
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Calc(%q) error = %v, want substring %q",
						tt.expression, err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Fatalf("Calc(%q) got error %v, want nil", tt.expression, err)
				}
				if got != tt.want {
					t.Errorf("Calc(%q) = %v, want %v", tt.expression, got, tt.want)
				}
			}
		})
	}
}

func TestCalculateHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(calculateHandler))
	defer srv.Close()
	postJSON := func(url string, body []byte) (*http.Response, []byte, error) {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, err
		}
		respBody := make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(respBody)
		resp.Body.Close()
		return resp, respBody, nil
	}
	tests := []struct {
		name           string
		expression     string
		wantStatusCode int
		wantField      string
		wantValuePart  string
	}{
		{
			name:           "Успешное вычисление",
			expression:     "2+2*2",
			wantStatusCode: http.StatusOK,
			wantField:      "result",
			wantValuePart:  "6",
		},
		{
			name:           "Ошибка валидации (несбалансированные скобки)",
			expression:     "2+(3*2",
			wantStatusCode: http.StatusUnprocessableEntity,
			wantField:      "error",
			wantValuePart:  "Expression is not valid",
		},
		{
			name:           "Деление на ноль",
			expression:     "2/0",
			wantStatusCode: http.StatusUnprocessableEntity,
			wantField:      "error",
			wantValuePart:  "Expression is not valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]string{
				"expression": tt.expression,
			})
			resp, respBody, err := postJSON(srv.URL, reqBody)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d", resp.StatusCode, tt.wantStatusCode)
			}
			var result map[string]string
			if err := json.Unmarshal(respBody, &result); err != nil {
				t.Fatalf("invalid JSON in response: %s", string(respBody))
			}
			gotValue, ok := result[tt.wantField]
			if !ok {
				t.Fatalf("response does not contain field %q", tt.wantField)
			}
			if !contains(gotValue, tt.wantValuePart) {
				t.Errorf("field %q = %q, want substring %q",
					tt.wantField, gotValue, tt.wantValuePart)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && bytes.Contains([]byte(s), []byte(substr)))
}
