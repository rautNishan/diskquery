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
	fmt.Printf("Scanner: %+v\n", scanner)
}
