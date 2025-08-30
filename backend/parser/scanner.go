package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota
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

var keywordsReverse = map[TokenType]string{
	TOKEN_SELECT:      "SELECT",
	TOKEN_FROM:        "FROM",
	TOKEN_WHERE:       "WHERE",
	TOKEN_INSERT:      "INSERT",
	TOKEN_INTO:        "INTO",
	TOKEN_VALUES:      "VALUES",
	TOKEN_UPDATE:      "UPDATE",
	TOKEN_SET:         "SET",
	TOKEN_DELETE:      "DELETE",
	TOKEN_CREATE:      "CREATE",
	TOKEN_TABLE:       "TABLE",
	TOKEN_DROP:        "DROP",
	TOKEN_INDEX:       "INDEX",
	TOKEN_VIEW:        "VIEW",
	TOKEN_DATABASE:    "DATABASE",
	TOKEN_SCHEMA:      "SCHEMA",
	TOKEN_FUNCTION:    "FUNCTION",
	TOKEN_PROCEDURE:   "PROCEDURE",
	TOKEN_TRIGGER:     "TRIGGER",
	TOKEN_BEGIN:       "BEGIN",
	TOKEN_END:         "END",
	TOKEN_COMMIT:      "COMMIT",
	TOKEN_ROLLBACK:    "ROLLBACK",
	TOKEN_TRANSACTION: "TRANSACTION",
	TOKEN_AS:          "AS",
	TOKEN_ON:          "ON",
	TOKEN_INNER:       "INNER",
	TOKEN_LEFT:        "LEFT",
	TOKEN_RIGHT:       "RIGHT",
	TOKEN_FULL:        "FULL",
	TOKEN_OUTER:       "OUTER",
	TOKEN_JOIN:        "JOIN",
	TOKEN_UNION:       "UNION",
	TOKEN_INTERSECT:   "INTERSECT",
	TOKEN_EXCEPT:      "EXCEPT",
	TOKEN_GROUP:       "GROUP",
	TOKEN_BY:          "BY",
	TOKEN_ORDER:       "ORDER",
	TOKEN_HAVING:      "HAVING",
	TOKEN_LIMIT:       "LIMIT",
	TOKEN_OFFSET:      "OFFSET",
	TOKEN_DISTINCT:    "DISTINCT",
	TOKEN_ALL:         "ALL",
	TOKEN_AND:         "AND",
	TOKEN_OR:          "OR",
	TOKEN_NOT:         "NOT",
	TOKEN_NULL:        "NULL",
	TOKEN_IS:          "IS",
	TOKEN_IN:          "IN",
	TOKEN_EXISTS:      "EXISTS",
	TOKEN_BETWEEN:     "BETWEEN",
	TOKEN_LIKE:        "LIKE",
	TOKEN_ILIKE:       "ILIKE",
	TOKEN_SIMILAR:     "SIMILAR",
	TOKEN_CASE:        "CASE",
	TOKEN_WHEN:        "WHEN",
	TOKEN_THEN:        "THEN",
	TOKEN_ELSE:        "ELSE",
	TOKEN_CAST:        "CAST",
	TOKEN_EXTRACT:     "EXTRACT",
	TOKEN_SUBSTRING:   "SUBSTRING",
	TOKEN_POSITION:    "POSITION",
	TOKEN_OVERLAY:     "OVERLAY",
	TOKEN_TRIM:        "TRIM",
	TOKEN_COALESCE:    "COALESCE",
	TOKEN_NULLIF:      "NULLIF",
	TOKEN_GREATEST:    "GREATEST",
	TOKEN_LEAST:       "LEAST",
	TOKEN_TRUE:        "TRUE",
	TOKEN_FALSE:       "FALSE",
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

func (s *Scanner) peekChar() rune {
	return s.peek
}

func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.current) {
		s.readChar()
	}
}

// This scan an identifier or keyword
func (s *Scanner) scanIdentifier() Token {
	start := s.location - 1

	var builder strings.Builder

	builder.WriteRune(s.current)
	s.readChar()

	//Subsequent characters (letter ,digits, underscore)
	for unicode.IsLetter(s.current) || unicode.IsDigit(s.current) || s.current == '_' {
		builder.WriteRune(s.current)
		s.readChar()
	}
	value := builder.String()

	if tokenType, isKeyword := keywords[strings.ToUpper(value)]; isKeyword {
		return Token{
			Type:     tokenType,
			Value:    strings.ToUpper(value),
			Location: start,
		}
	}

	return Token{
		Type:     TOKEN_IDENT,
		Value:    value,
		Location: start,
	}

}

func (s *Scanner) scanQuotedIdentifier() Token {
	start := s.location - 1
	var builder strings.Builder

	s.readChar()

	for s.current != '"' && s.current != 0 {
		if s.current == '"' && s.peekChar() == '"' {
			builder.WriteRune('"')
			s.readChar() //Skip first quote
			s.readChar() //Skip second quote
		} else {
			builder.WriteRune(s.current)
			s.readChar()
		}
	}

	if s.current == '"' {
		s.readChar() //Skip closing quote
	}

	return Token{
		Type:     TOKEN_IDENT,
		Value:    builder.String(),
		Location: start,
	}
}

func (s *Scanner) scanNumber() Token {
	start := s.location - 1
	var builder strings.Builder

	isFloat := false

	for unicode.IsDigit(s.current) {
		builder.WriteRune(s.current)
		s.readChar()
	}

	if s.current == '.' && unicode.IsDigit(s.peekChar()) {
		isFloat = true
		builder.WriteRune(s.current)
		s.readChar()

		for unicode.IsDigit(s.current) {
			builder.WriteRune(s.current)
			s.readChar()
		}
	}

	//Exponent part
	if s.current == 'e' || s.current == 'E' {
		isFloat = true
		builder.WriteRune(s.current)
		s.readChar()

		if s.current == '+' || s.current == '-' {
			builder.WriteRune(s.current)
			s.readChar()
		}

		for unicode.IsDigit(s.current) {
			builder.WriteRune(s.current)
			s.readChar()
		}
	}

	value := builder.String()

	if isFloat {
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Token{Type: TOKEN_ERROR, Value: "Invalid float: " + value, Location: start}
		}

		return Token{
			Type:     TOKEN_FCONST,
			Value:    value,
			FloatVal: floatVal,
			Location: start,
		}
	}

	intVal, err := strconv.ParseInt(value, 10, 60)
	if err != nil {
		return Token{
			Type:     TOKEN_ERROR,
			Value:    "Invalid integer: " + value,
			Location: start,
		}
	}

	return Token{
		Type:     TOKEN_ICONST,
		Value:    value,
		IntVal:   intVal,
		Location: start,
	}
}

func (s *Scanner) scanString(quote rune) Token {
	start := s.location - 1
	var builder strings.Builder
	s.readChar()

	for s.current != quote && s.current != 0 {
		if s.current == quote && s.peekChar() == quote {
			builder.WriteRune(s.current)
			s.readChar() //Skip first quote
			s.readChar() //Skip second quote
		} else if s.current == '\\' {
			s.readChar()
			switch s.current {
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			case 'r':
				builder.WriteRune('\r')
			case '\\':
				builder.WriteRune('\\')
			case '\'':
				builder.WriteRune('\'')
			case '"':
				builder.WriteRune('"')
			default:
				builder.WriteRune(s.current)
			}
			s.readChar()
		} else {
			builder.WriteRune(s.current)
			s.readChar()
		}
	}
	if s.current == quote {
		s.readChar() //Skip closing quote
	}

	return Token{
		Type:     TOKEN_SCONST,
		Value:    builder.String(),
		Location: start,
	}
}

// scanParameter scans a paramater marker ($1, $2, etc)
func (s *Scanner) scanParameter() Token {
	start := s.location - 1
	var builder strings.Builder

	builder.WriteRune(s.current)
	s.readChar()

	if !unicode.IsDigit(s.current) {
		return Token{Type: TOKEN_ERROR, Value: "Invalid parameter", Location: start}
	}

	for unicode.IsDigit(s.current) {
		builder.WriteRune(s.current)
		s.readChar()
	}

	return Token{
		Type:     TOKEN_PARAM,
		Value:    builder.String(),
		Location: start,
	}
}

// scanComment scans a comment and returns the next token
func (s *Scanner) scanComment() Token {
	if s.current == '-' && s.peekChar() == '-' {
		// Line comment
		for s.current != '\n' && s.current != 0 {
			s.readChar()
		}
	} else if s.current == '/' && s.peekChar() == '*' {
		// Block comment
		s.readChar() // skip '/'
		s.readChar() // skip '*'

		for {
			if s.current == 0 {
				break
			}
			if s.current == '*' && s.peekChar() == '/' {
				s.readChar() // skip '*'
				s.readChar() // skip '/'
				break
			}
			s.readChar()
		}
	}

	// Return next token after comment
	return s.NextToken()
}

// NextToken returns the next token from the input
func (s *Scanner) NextToken() Token {
	s.skipWhitespace()

	if s.current == 0 {
		return Token{Type: TOKEN_EOF, Location: s.location}
	}
	location := s.location - 1

	switch {
	case unicode.IsLetter(s.current) || s.current == '_':
		return s.scanIdentifier()

	case s.current == '"':
		return s.scanQuotedIdentifier()

	case unicode.IsDigit(s.current):
		return s.scanNumber()

	case s.current == '\'' || s.current == '"':
		return s.scanString(s.current)

	case s.current == '$':
		return s.scanParameter()

	case s.current == '-' && s.peekChar() == '-':
		return s.scanComment()

	case s.current == '/' && s.peekChar() == '*':
		return s.scanComment()

	// Two-character operators
	case s.current == '<':
		s.readChar()
		if s.current == '=' {
			s.readChar()
			return Token{Type: TOKEN_LE, Value: "<=", Location: location}
		} else if s.current == '>' {
			s.readChar()
			return Token{Type: TOKEN_NE, Value: "<>", Location: location}
		}
		return Token{Type: TOKEN_LT, Value: "<", Location: location}

	case s.current == '>':
		s.readChar()
		if s.current == '=' {
			s.readChar()
			return Token{Type: TOKEN_GE, Value: ">=", Location: location}
		}
		return Token{Type: TOKEN_GT, Value: ">", Location: location}

	case s.current == '!':
		s.readChar()
		if s.current == '=' {
			s.readChar()
			return Token{Type: TOKEN_NE, Value: "!=", Location: location}
		}
		return Token{Type: TOKEN_ERROR, Value: "unexpected character: !", Location: location}

	case s.current == '|':
		s.readChar()
		if s.current == '|' {
			s.readChar()
			return Token{Type: TOKEN_CONCAT, Value: "||", Location: location}
		}
		return Token{Type: TOKEN_ERROR, Value: "unexpected character: |", Location: location}

	case s.current == ':':
		s.readChar()
		if s.current == '=' {
			s.readChar()
			return Token{Type: TOKEN_ASSIGN, Value: ":=", Location: location}
		}
		return Token{Type: TOKEN_ERROR, Value: "unexpected character: :", Location: location}

	case s.current == '.':
		s.readChar()
		if s.current == '.' {
			s.readChar()
			return Token{Type: TOKEN_DOTDOT, Value: "..", Location: location}
		}
		return Token{Type: TOKEN_DOT, Value: ".", Location: location}

	// Single-character tokens
	case s.current == '=':
		s.readChar()
		return Token{Type: TOKEN_EQ, Value: "=", Location: location}

	case s.current == '+':
		s.readChar()
		return Token{Type: TOKEN_PLUS, Value: "+", Location: location}

	case s.current == '-':
		s.readChar()
		return Token{Type: TOKEN_MINUS, Value: "-", Location: location}

	case s.current == '*':
		s.readChar()
		return Token{Type: TOKEN_MULTIPLY, Value: "*", Location: location}

	case s.current == '/':
		s.readChar()
		return Token{Type: TOKEN_DIVIDE, Value: "/", Location: location}

	case s.current == '%':
		s.readChar()
		return Token{Type: TOKEN_MODULO, Value: "%", Location: location}

	case s.current == '^':
		s.readChar()
		return Token{Type: TOKEN_POWER, Value: "^", Location: location}

	case s.current == '(':
		s.readChar()
		return Token{Type: TOKEN_LPAREN, Value: "(", Location: location}

	case s.current == ')':
		s.readChar()
		return Token{Type: TOKEN_RPAREN, Value: ")", Location: location}

	case s.current == ',':
		s.readChar()
		return Token{Type: TOKEN_COMMA, Value: ",", Location: location}

	case s.current == ';':
		s.readChar()
		return Token{Type: TOKEN_SEMICOLON, Value: ";", Location: location}

	case s.current == '[':
		s.readChar()
		return Token{Type: TOKEN_LBRACKET, Value: "[", Location: location}

	case s.current == ']':
		s.readChar()
		return Token{Type: TOKEN_RBRACKET, Value: "]", Location: location}

	case s.current == '{':
		s.readChar()
		return Token{Type: TOKEN_LBRACE, Value: "{", Location: location}

	case s.current == '}':
		s.readChar()
		return Token{Type: TOKEN_RBRACE, Value: "}", Location: location}

	default:
		char := s.current
		s.readChar()
		return Token{
			Type:     TOKEN_ERROR,
			Value:    fmt.Sprintf("unexpected character: %c", char),
			Location: location,
		}
	}
}

// SetState sets the scanner state for different parsing modes
func (s *Scanner) SetState(state ScannerState) {
	s.state = state
}

// GetTokens tokenizes the entire input and returns all tokens
func (s *Scanner) GetTokens() []Token {
	var tokens []Token

	for {
		token := s.NextToken()
		tokens = append(tokens, token)
		if token.Type == TOKEN_EOF || token.Type == TOKEN_ERROR {
			break
		}
	}

	return tokens
}
