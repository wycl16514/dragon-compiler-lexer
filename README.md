编译原理是非常基础，难度也大，而且是比较晦涩难懂的基础理论。也难怪说世界上能做编译器的人总共坐不满一个教室。虽然不简单，但我们也不要随意放弃，尽自己最大努力，看看能不能在教室里挤出个座位。

在我看来编译原理的复杂主要在于语言难以描述上，因此我们必须辅助其他手段对其进行理解和掌握，这里我们就得依靠代码的方式来辅助我们在言语上的理解不足。例如在词法解析中有个概念叫token，它其实是将源代码中的字符串进行分类，并给每种分类赋予一个整数值，描述起来不好理解，我们看看代码实现就明白了，首先完成如下定义，token.go:
```
package lexer

type Tag uint32

const (
	AND Tag = iota + 256
	BASIC
	BREAK
	DO
	EQ
	FALSE
	GE
	ID
	IF
	ELSE
	INDEX
	LE
	INT 
	FLOAT 
	MINUS
	PLUS
	NE
	NUM
	OR
	REAL
	TEMP
	TRUE
	WHILE
	LEFT_BRACE
	RIGHT_BRACE
	AND_OPERATOR
	OR_OPERATOR
	ASSIGN_OPERATOR
	NEGATE_OPERATOR
	LESS_OPERATOR
	GREATER_OPERATOR

	EOF

	ERROR
)

var token_map = make(map[Tag]string)

func init() {

	token_map[AND] = "&&"
	token_map[BASIC] = "BASIC"
	token_map[DO] = "do"
	token_map[ELSE] = "else"
	token_map[EQ] = "EQ"
	token_map[FALSE] = "FALSE"
	token_map[GE] = "GE"
	token_map[ID] = "ID"
	token_map[IF] = "if"
	token_map[INT] = "int"
	token_map[FLOAT] = "float"
	token_map[INDEX] = "INDEX"
	token_map[LE] = "<="
	token_map[MINUS] = "-"
	token_map[PLUS] = "+"
	token_map[NE] = "!="
	token_map[NUM] = "NUM"
	token_map[OR] = "OR"
	token_map[REAL] = "REAL"
	token_map[TEMP] = "TEMP"
	token_map[TRUE] = "TRUE"
	token_map[WHILE] = "while"
	token_map[AND_OPERATOR] = "&"
	token_map[OR_OPERATOR] = "|"
	token_map[ASSIGN_OPERATOR] = "="
	token_map[NEGATE_OPERATOR] = "!"
	token_map[LESS_OPERATOR] = "<"
	token_map[GREATER_OPERATOR] = ">"
	token_map[LEFT_BRACE] = "{"
	token_map[RIGHT_BRACE] = "}"
	token_map[EOF] = "EOF"
	token_map[ERROR] = "ERROR"
}

type Token struct {
	Tag Tag
}

func (t *Token) ToString() string {
	return token_map[t.Tag]
}

func NewToken(tag Tag) Token {
	return Token{
		Tag: tag,
	}
}

```
从上面代码看，所谓token其实就是简单的枚举类型，我们接下来看看关键字和数字对应的token实现，增加word_token.go,其内容如下：
```
package lexer

type Word struct {
	lexeme string
	Tag    Token
}

func NewWordToken(s string, tag Tag) Word {
	return Word{
		lexeme: s,
		Tag:    NewToken(tag),
	}
}

func (w *Word) ToString() string {
	return w.lexeme
}

func GetKeyWords() []Word {
	key_words := []Word{}
	key_words = append(key_words, NewWordToken("&&", AND))
	key_words = append(key_words, NewWordToken("||", OR))
	key_words = append(key_words, NewWordToken("==", EQ))
	key_words = append(key_words, NewWordToken("!=", NE))
	key_words = append(key_words, NewWordToken("<=", LE))
	key_words = append(key_words, NewWordToken(">=", GE))
	key_words = append(key_words, NewWordToken("minus", MINUS))
	key_words = append(key_words, NewWordToken("true", TRUE))
	key_words = append(key_words, NewWordToken("false", FALSE))
	key_words = append(key_words, NewWordToken("t", TEMP))
	key_words = append(key_words, NewWordToken("if", IF))
	key_words = append(key_words, NewWordToken("else", ELSE))

	return key_words
}

```
可以看到workd_token其实就是将源代码中特定字符串赋予特定的整数值，特别是对关键字字符串进行留存以便以后使用，同样的我们增加number_token.go，用来给数字字符串赋予特定整数值：
```
package lexer 

import (
	"strconv"
	"fmt"
)

type Num struct {
	Tag Token
	value int 
}

func NewNumToken(val int) Num {
	return Num {
		Tag: NewToken(NUM),
		value: val,
	}
}

func (n *Num) ToString() string {
	return  strconv.Itoa(n.value)
}

type Real struct {
	Tag Token 
	value float64 
}

func NewRealToken(val float64) Real{
	return Real {
		value: val,
		Tag: NewToken(REAL),
	}
}

func (r *Real) ToString() (string) {
    return fmt.Sprintf("%.7f", r.value)
}
```
最后我们看看词法扫描本质在做什么，词法扫描就是依次读入源代码字符，然后看所读到的字符到底属于那种类型的字符串，然后将字符串与给定类型的token对象联系起来，增加lexer.go然后我们看看其实现：
```
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

```

最后我们增加main.go，然后将上面的代码运行起来：
```
package main

import (
	"fmt"
	"lexer"
)

func main() {
	source := "if a >= 100.34"
	my_lexer := lexer.NewLexer(source)
	for {
		token, err := my_lexer.Scan()
		if err != nil {
			fmt.Println("lexer error: ", err)
			break
		}

		if token.Tag == lexer.EOF {
			break
		} else {
			fmt.Println("read token: ", token)
		}
	}
}

```
对于以上代码的具体逻辑，大家可以在我的公众号或是在b站上搜索"lexer“就可以看到具体的讲解和调试演示过程。
