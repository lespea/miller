// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"miller/parsing/token"
)

const (
	NoState    = -1
	NumStates  = 192
	NumSymbols = 316
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = new(token.Token)
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '"'
1: '"'
2: '0'
3: 'x'
4: '0'
5: 'b'
6: '.'
7: '-'
8: '.'
9: '.'
10: '-'
11: '.'
12: '.'
13: '-'
14: 'M'
15: '_'
16: 'P'
17: 'I'
18: 'M'
19: '_'
20: 'E'
21: 'I'
22: 'P'
23: 'S'
24: 'I'
25: 'F'
26: 'S'
27: 'I'
28: 'R'
29: 'S'
30: 'O'
31: 'P'
32: 'S'
33: 'O'
34: 'F'
35: 'S'
36: 'O'
37: 'R'
38: 'S'
39: 'N'
40: 'F'
41: 'N'
42: 'R'
43: 'F'
44: 'N'
45: 'R'
46: 'F'
47: 'I'
48: 'L'
49: 'E'
50: 'N'
51: 'A'
52: 'M'
53: 'E'
54: 'F'
55: 'I'
56: 'L'
57: 'E'
58: 'N'
59: 'U'
60: 'M'
61: 'b'
62: 'e'
63: 'g'
64: 'i'
65: 'n'
66: 'e'
67: 'n'
68: 'd'
69: 'f'
70: 'i'
71: 'l'
72: 't'
73: 'e'
74: 'r'
75: 'i'
76: 'n'
77: 't'
78: 'f'
79: 'l'
80: 'o'
81: 'a'
82: 't'
83: '$'
84: '$'
85: '{'
86: '}'
87: '$'
88: '*'
89: '@'
90: '@'
91: '{'
92: '}'
93: '@'
94: '*'
95: '%'
96: '%'
97: '%'
98: 'p'
99: 'a'
100: 'n'
101: 'i'
102: 'c'
103: '%'
104: '%'
105: '%'
106: ';'
107: '{'
108: '}'
109: '='
110: '['
111: ']'
112: '$'
113: '['
114: '@'
115: '['
116: '|'
117: '|'
118: '='
119: '^'
120: '^'
121: '='
122: '&'
123: '&'
124: '='
125: '|'
126: '='
127: '^'
128: '='
129: '<'
130: '<'
131: '='
132: '>'
133: '>'
134: '='
135: '>'
136: '>'
137: '>'
138: '='
139: '+'
140: '='
141: '.'
142: '='
143: '-'
144: '='
145: '*'
146: '='
147: '/'
148: '='
149: '/'
150: '/'
151: '='
152: '%'
153: '='
154: '*'
155: '*'
156: '='
157: '?'
158: ':'
159: '|'
160: '|'
161: '^'
162: '^'
163: '&'
164: '&'
165: '='
166: '~'
167: '!'
168: '='
169: '~'
170: '='
171: '='
172: '!'
173: '='
174: '>'
175: '>'
176: '='
177: '<'
178: '<'
179: '='
180: '|'
181: '^'
182: '&'
183: '<'
184: '<'
185: '>'
186: '>'
187: '>'
188: '>'
189: '>'
190: '+'
191: '-'
192: '.'
193: '+'
194: '.'
195: '-'
196: '.'
197: '*'
198: '/'
199: '/'
200: '/'
201: '%'
202: '.'
203: '*'
204: '.'
205: '/'
206: '.'
207: '/'
208: '/'
209: '!'
210: '~'
211: '*'
212: '*'
213: '('
214: ')'
215: ','
216: '_'
217: ' '
218: '!'
219: '#'
220: '$'
221: '%'
222: '&'
223: '''
224: '\'
225: '('
226: ')'
227: '*'
228: '+'
229: ','
230: '-'
231: '.'
232: '/'
233: ':'
234: ';'
235: '<'
236: '='
237: '>'
238: '?'
239: '@'
240: '['
241: ']'
242: '^'
243: '_'
244: '`'
245: '{'
246: '|'
247: '}'
248: '~'
249: '\'
250: '"'
251: 'e'
252: 'E'
253: 't'
254: 'r'
255: 'u'
256: 'e'
257: 'f'
258: 'a'
259: 'l'
260: 's'
261: 'e'
262: ' '
263: '!'
264: '#'
265: '$'
266: '%'
267: '&'
268: '''
269: '\'
270: '('
271: ')'
272: '*'
273: '+'
274: ','
275: '-'
276: '.'
277: '/'
278: ':'
279: ';'
280: '<'
281: '='
282: '>'
283: '?'
284: '@'
285: '['
286: ']'
287: '^'
288: '_'
289: '`'
290: '|'
291: '~'
292: '\'
293: '{'
294: '\'
295: '}'
296: ' '
297: '\t'
298: '\n'
299: '\r'
300: 'a'-'z'
301: 'A'-'Z'
302: '0'-'9'
303: '0'-'9'
304: 'a'-'f'
305: 'A'-'F'
306: '0'-'1'
307: 'A'-'Z'
308: 'a'-'z'
309: '0'-'9'
310: \u0100-\U0010ffff
311: 'A'-'Z'
312: 'a'-'z'
313: '0'-'9'
314: \u0100-\U0010ffff
315: .
*/
