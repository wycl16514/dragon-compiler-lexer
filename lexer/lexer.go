package lexer

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	peek      byte
	line      int
	reader    *bufio.Reader
	key_words map[string]Token
}

func NewLexer(source string) Lexer {
	str := strings.NewReader(source)
	source_reader := bufio.NewReaderSize(str, len(source))
	lexer := Lexer{
		line:      1,
		reader:    source_reader,
		key_words: make(map[string]Token),
	}

	lexer.reserve()

	return lexer
}

func (l *Lexer) reserve() {
	key_words := GetKeyWords()
	for _, key_word := range key_words {
		l.key_words[key_word.ToString()] = key_word.Tag
	}
}

func (l *Lexer) Readch() error {
	char, err := l.reader.ReadByte() //提前读取下一个字符
	l.peek = char
	return err
}

func (l *Lexer) ReadCharacter(c byte) (bool, error) {
	chars, err := l.reader.Peek(1)
	if err != nil {
		return false, err
	}

	peekChar := chars[0]
	if peekChar != c {
		return false, nil
	}

	l.Readch() //越过当前peek的字符
	return true, nil
}

func (l *Lexer) Scan() (Token, error) {
	for {
		err := l.Readch()
		if err != nil {
			return NewToken(ERROR), err
		}

		if l.peek == ' ' || l.peek == '\t' {
			continue
		} else if l.peek == '\n' {
			l.line = l.line + 1
		} else {
			break
		}
	}

	switch l.peek {
	case '{':
		return NewToken(LEFT_BRACE), nil
	case '}':
		return NewToken(RIGHT_BRACE), nil
	case '+':
		return NewToken(PLUS), nil
	case '-':
		return NewToken(MINUS), nil
	case '&':
		if ok, err := l.ReadCharacter('&'); ok {
			word := NewWordToken("&&", AND)
			return word.Tag, err
		} else {
			return NewToken(AND_OPERATOR), err
		}
	case '|':
		if ok, err := l.ReadCharacter('|'); ok {
			word := NewWordToken("||", OR)
			return word.Tag, err
		} else {
			return NewToken(OR_OPERATOR), err
		}

	case '=':
		if ok, err := l.ReadCharacter('='); ok {
			word := NewWordToken("==", EQ)
			return word.Tag, err
		} else {
			return NewToken(ASSIGN_OPERATOR), err
		}

	case '!':
		if ok, err := l.ReadCharacter('='); ok {
			word := NewWordToken("!=", NE)
			return word.Tag, err
		} else {
			return NewToken(NEGATE_OPERATOR), err
		}

	case '<':
		if ok, err := l.ReadCharacter('='); ok {
			word := NewWordToken("<=", LE)
			return word.Tag, err
		} else {
			return NewToken(LESS_OPERATOR), err
		}

	case '>':
		if ok, err := l.ReadCharacter('='); ok {
			word := NewWordToken(">=", GE)
			return word.Tag, err
		} else {
			return NewToken(GREATER_OPERATOR), err
		}

	}

	if unicode.IsNumber(rune(l.peek)) {
		var v int
		var err error
		for {
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				break
			}
			v = 10*v + num
			l.Readch()
		}

		if l.peek != '.' {
			return NewToken(NUM), err
		}

		x := float64(v)
		d := float64(10)
		for {
			l.Readch()
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				break
			}

			x = x + float64(num)/d
			d = d * 10
		}

		return NewToken(REAL), err
	}

	if unicode.IsLetter(rune(l.peek)) {
		var buffer []byte
		for {
			buffer = append(buffer, l.peek)
			l.Readch()
			if !unicode.IsLetter(rune(l.peek)) {
				break
			}
		}

		s := string(buffer)
		token, ok := l.key_words[s]
		if ok {
			return token, nil
		}

		return NewToken(ID), nil
	}

	return NewToken(EOF), nil
}
