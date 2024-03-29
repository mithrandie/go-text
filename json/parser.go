// Code generated by goyacc -o json/parser.go -v json/parser.output json/parser.y. DO NOT EDIT.

//line json/parser.y:2
package json

import __yyfmt__ "fmt"

//line json/parser.y:2

import "strconv"

//line json/parser.y:7
type yySymType struct {
	yys            int
	structure      Structure
	structures     []Structure
	object_member  ObjectMember
	object_members []ObjectMember
	token          Token
}

const NUMBER = 57346
const STRING = 57347
const BOOLEAN = 57348
const NULL = 57349
const FLOAT = 57350
const INTEGER = 57351

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUMBER",
	"STRING",
	"BOOLEAN",
	"NULL",
	"FLOAT",
	"INTEGER",
	"':'",
	"','",
	"'{'",
	"'}'",
	"'['",
	"']'",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line json/parser.y:110

func ParseJson(src string, useInteger bool) (Structure, EscapeType, error) {
	l := new(Lexer)
	l.Init(src, useInteger)
	yyParse(l)
	return l.structure, l.EscapeType(), l.err
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 31

var yyAct = [...]int{
	6, 5, 9, 10, 7, 8, 14, 11, 3, 19,
	4, 15, 2, 16, 20, 17, 18, 13, 12, 1,
	0, 0, 0, 0, 0, 21, 0, 23, 0, 0,
	22,
}

var yyPact = [...]int{
	-4, -1000, -1000, 12, -4, -1000, -1000, -1000, -1000, -1000,
	-1000, 0, 4, 6, -6, 3, -1000, 12, -4, -1000,
	-4, -1000, -1000, -1000,
}

var yyPgo = [...]int{
	0, 19, 18, 7, 6, 11,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 3, 3, 3, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5,
}

var yyR2 = [...]int{
	0, 0, 1, 3, 0, 1, 3, 0, 1, 3,
	3, 3, 1, 1, 1, 1, 1, 1,
}

var yyChk = [...]int{
	-1000, -1, -5, 12, 14, 5, 4, 8, 9, 6,
	7, -3, -2, 5, -4, -5, 13, 11, 10, 15,
	11, -3, -5, -4,
}

var yyDef = [...]int{
	1, -2, 2, 4, 7, 12, 13, 14, 15, 16,
	17, 0, 5, 0, 0, 8, 10, 4, 0, 11,
	7, 6, 3, 9,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 11, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 10, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 14, 3, 15, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 12, 3, 13,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
//line json/parser.y:28
		{
			yyVAL.structure = nil
			yylex.(*Lexer).structure = yyVAL.structure
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:33
		{
			yyVAL.structure = yyDollar[1].structure
			yylex.(*Lexer).structure = yyVAL.structure
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
//line json/parser.y:40
		{
			yyVAL.object_member = ObjectMember{Key: yyDollar[1].token.Literal, Value: yyDollar[3].structure}
		}
	case 4:
		yyDollar = yyS[yypt-0 : yypt+1]
//line json/parser.y:46
		{
			yyVAL.object_members = nil
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:50
		{
			yyVAL.object_members = []ObjectMember{yyDollar[1].object_member}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
//line json/parser.y:54
		{
			yyVAL.object_members = append([]ObjectMember{yyDollar[1].object_member}, yyDollar[3].object_members...)
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
//line json/parser.y:60
		{
			yyVAL.structures = []Structure{}
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:64
		{
			yyVAL.structures = []Structure{yyDollar[1].structure}
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line json/parser.y:68
		{
			yyVAL.structures = append([]Structure{yyDollar[1].structure}, yyDollar[3].structures...)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line json/parser.y:74
		{
			yyVAL.structure = Object{Members: yyDollar[2].object_members}
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line json/parser.y:78
		{
			yyVAL.structure = Array(yyDollar[2].structures)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:82
		{
			yyVAL.structure = String(yyDollar[1].token.Literal)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:86
		{
			f, _ := strconv.ParseFloat(yyDollar[1].token.Literal, 64)
			yyVAL.structure = Number(f)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:91
		{
			f, _ := strconv.ParseFloat(yyDollar[1].token.Literal, 64)
			yyVAL.structure = Float(f)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:96
		{
			i, _ := strconv.ParseInt(yyDollar[1].token.Literal, 10, 64)
			yyVAL.structure = Integer(i)
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:101
		{
			b, _ := strconv.ParseBool(yyDollar[1].token.Literal)
			yyVAL.structure = Boolean(b)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line json/parser.y:106
		{
			yyVAL.structure = Null{}
		}
	}
	goto yystack /* stack new state and value */
}
