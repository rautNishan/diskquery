package parser

import (
	"fmt"
	"unicode/utf8"
)

type TokenType int

const (
	TOKEN_EFOF TokenType = iota
	TOKEN_ERROR
	TOKEN_SCONST // String constant
	TOKEN_ICONST // Integer constant
	TOKEN_FCONST // Float constant
	TOKEN_IDENT  // Identifier
	TOKEN_PARAM  // Parameter ($1, $2, etc.)

	// Operators
	TOKEN_EQ       // =
	TOKEN_NE       // != or <>
	TOKEN_LT       // <
	TOKEN_LE       // <=
	TOKEN_GT       // >
	TOKEN_GE       // >=
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_MULTIPLY // *
	TOKEN_DIVIDE   // /
	TOKEN_MODULO   // %
	TOKEN_POWER    // ^
	TOKEN_CONCAT   // ||
	TOKEN_ASSIGN   // :=
	TOKEN_DOT      // .
	TOKEN_DOTDOT   // ..

	// Punctuation
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_COMMA     // ,
	TOKEN_SEMICOLON // ;
	TOKEN_LBRACKET  // [
	TOKEN_RBRACKET  // ]
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }

	// Keywords (start from 100 to avoid conflicts)
	TOKEN_SELECT = 100 + iota
	TOKEN_FROM
	TOKEN_WHERE
	TOKEN_INSERT
	TOKEN_INTO
	TOKEN_VALUES
	TOKEN_UPDATE
	TOKEN_SET
	TOKEN_DELETE
	TOKEN_CREATE
	TOKEN_TABLE
	TOKEN_DROP
	TOKEN_INDEX
	TOKEN_VIEW
	TOKEN_DATABASE
	TOKEN_SCHEMA
	TOKEN_FUNCTION
	TOKEN_PROCEDURE
	TOKEN_TRIGGER
	TOKEN_BEGIN
	TOKEN_END
	TOKEN_COMMIT
	TOKEN_ROLLBACK
	TOKEN_TRANSACTION
	TOKEN_AS
	TOKEN_ON
	TOKEN_INNER
	TOKEN_LEFT
	TOKEN_RIGHT
	TOKEN_FULL
	TOKEN_OUTER
	TOKEN_JOIN
	TOKEN_UNION
	TOKEN_INTERSECT
	TOKEN_EXCEPT
	TOKEN_GROUP
	TOKEN_BY
	TOKEN_ORDER
	TOKEN_HAVING
	TOKEN_LIMIT
	TOKEN_OFFSET
	TOKEN_DISTINCT
	TOKEN_ALL
	TOKEN_AND
	TOKEN_OR
	TOKEN_NOT
	TOKEN_NULL
	TOKEN_IS
	TOKEN_IN
	TOKEN_EXISTS
	TOKEN_BETWEEN
	TOKEN_LIKE
	TOKEN_ILIKE
	TOKEN_SIMILAR
	TOKEN_CASE
	TOKEN_WHEN
	TOKEN_THEN
	TOKEN_ELSE
	TOKEN_CAST
	TOKEN_EXTRACT
	TOKEN_SUBSTRING
	TOKEN_POSITION
	TOKEN_OVERLAY
	TOKEN_TRIM
	TOKEN_COALESCE
	TOKEN_NULLIF
	TOKEN_GREATEST
	TOKEN_LEAST
	TOKEN_TRUE
	TOKEN_FALSE
)

// Lexical token
type Token struct {
	Type     TokenType
	Value    string
	IntVal   int64
	FloatVal float64
	Location int
}

type ScannerState int

const (
	SCANNER_NORMAL ScannerState = iota // 0
	SCANNER_TYPE_NAME
	SCANNER_EXPR
	SCANNER_ASSIGN
)

//Scanner handles lexical analysis

type Scanner struct {
	query    string
	position int //Where we are in bytes
	current  rune
	peek     rune
	location int
	state    ScannerState
}

// Keywords mapping - case insensitive
var keywords = map[string]TokenType{
	"SELECT":      TOKEN_SELECT,
	"FROM":        TOKEN_FROM,
	"WHERE":       TOKEN_WHERE,
	"INSERT":      TOKEN_INSERT,
	"INTO":        TOKEN_INTO,
	"VALUES":      TOKEN_VALUES,
	"UPDATE":      TOKEN_UPDATE,
	"SET":         TOKEN_SET,
	"DELETE":      TOKEN_DELETE,
	"CREATE":      TOKEN_CREATE,
	"TABLE":       TOKEN_TABLE,
	"DROP":        TOKEN_DROP,
	"INDEX":       TOKEN_INDEX,
	"VIEW":        TOKEN_VIEW,
	"DATABASE":    TOKEN_DATABASE,
	"SCHEMA":      TOKEN_SCHEMA,
	"FUNCTION":    TOKEN_FUNCTION,
	"PROCEDURE":   TOKEN_PROCEDURE,
	"TRIGGER":     TOKEN_TRIGGER,
	"BEGIN":       TOKEN_BEGIN,
	"END":         TOKEN_END,
	"COMMIT":      TOKEN_COMMIT,
	"ROLLBACK":    TOKEN_ROLLBACK,
	"TRANSACTION": TOKEN_TRANSACTION,
	"AS":          TOKEN_AS,
	"ON":          TOKEN_ON,
	"INNER":       TOKEN_INNER,
	"LEFT":        TOKEN_LEFT,
	"RIGHT":       TOKEN_RIGHT,
	"FULL":        TOKEN_FULL,
	"OUTER":       TOKEN_OUTER,
	"JOIN":        TOKEN_JOIN,
	"UNION":       TOKEN_UNION,
	"INTERSECT":   TOKEN_INTERSECT,
	"EXCEPT":      TOKEN_EXCEPT,
	"GROUP":       TOKEN_GROUP,
	"BY":          TOKEN_BY,
	"ORDER":       TOKEN_ORDER,
	"HAVING":      TOKEN_HAVING,
	"LIMIT":       TOKEN_LIMIT,
	"OFFSET":      TOKEN_OFFSET,
	"DISTINCT":    TOKEN_DISTINCT,
	"ALL":         TOKEN_ALL,
	"AND":         TOKEN_AND,
	"OR":          TOKEN_OR,
	"NOT":         TOKEN_NOT,
	"NULL":        TOKEN_NULL,
	"IS":          TOKEN_IS,
	"IN":          TOKEN_IN,
	"EXISTS":      TOKEN_EXISTS,
	"BETWEEN":     TOKEN_BETWEEN,
	"LIKE":        TOKEN_LIKE,
	"ILIKE":       TOKEN_ILIKE,
	"SIMILAR":     TOKEN_SIMILAR,
	"CASE":        TOKEN_CASE,
	"WHEN":        TOKEN_WHEN,
	"THEN":        TOKEN_THEN,
	"ELSE":        TOKEN_ELSE,
	"CAST":        TOKEN_CAST,
	"EXTRACT":     TOKEN_EXTRACT,
	"SUBSTRING":   TOKEN_SUBSTRING,
	"POSITION":    TOKEN_POSITION,
	"OVERLAY":     TOKEN_OVERLAY,
	"TRIM":        TOKEN_TRIM,
	"COALESCE":    TOKEN_COALESCE,
	"NULLIF":      TOKEN_NULLIF,
	"GREATEST":    TOKEN_GREATEST,
	"LEAST":       TOKEN_LEAST,
	"TRUE":        TOKEN_TRUE,
	"FALSE":       TOKEN_FALSE,
}

func NewScanner(query string, state ScannerState) *Scanner {
	s := &Scanner{
		query:    query,
		position: 0,
		location: 0,
		state:    state,
	}
	s.readChar()
	return s
}

func (s *Scanner) readChar() {
	// Check if we have reached the end of the query string
	if s.position >= len(s.query) {
		// End of input: mark current and peek as 0 (EOF)
		s.current = 0
		s.peek = 0
	} else {
		// Shift previous peek to current
		s.current = s.peek
		// Update peek: the next character after the current position
		if s.position+1 >= len(s.query) {
			// If we are at the last character, peek is 0 (EOF)
			s.peek = 0
		} else {
			// Decode the rune at the next position safely (UTF-8 aware)
			r, size := utf8.DecodeRuneInString(s.query[s.position+1:])
			s.peek = r
			_ = size // size is ignored because we don't need it here
		}

		// Special handling for the very first character in the query

		if s.position == 0 {
			// Decode the rune at the current position
			r, _ := utf8.DecodeRuneInString(s.query[s.position:])
			s.current = r
		}

		r, size := utf8.DecodeRuneInString(s.query[s.position:])
		fmt.Printf("Size: %d\n", size)
		s.current = r
		s.position += size
		s.location++
	}
}
