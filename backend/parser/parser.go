package parser

import "fmt"

type RawParseMode int

const (
	RAW_PARSE_DEFAULT RawParseMode = iota //0
	RAW_PARSE_TYPE_NAME
	RAW_PARSE_SQL_EXPR
	RAW_PARSE_SQL_ASSIGN1
	RAW_PARSE_SQL_ASSIGN2
	RAW_PARSE_SQL_ASSIGN3
)

func RawParse(query string, parseMode RawParseMode) {
	scanner := NewScanner(query, 0)
	tokens := scanner.GetTokens()
	for _, token := range tokens {
		fmt.Printf("Token raw type: %d\n", token.Type)
		fmt.Printf("TokeyType: %s, Value: %s\n", keywordsReverse[token.Type], token.Value)
		fmt.Printf("Token details: %v\n", token)
	}
}
