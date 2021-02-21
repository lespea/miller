package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"miller/src/lib"
)

// ================================================================
func MlrvalStrlen(output, input1 *Mlrval) {
	if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromInt(int(utf8.RuneCountInString(input1.printrep)))
	}
}

// ================================================================
func MlrvalToString(output, input1 *Mlrval) {
	output.SetFromString(input1.String())
}

// ================================================================
// Dot operator, with loose typecasting.
//
// For most operations, I don't like loose typecasting -- for example, in PHP
// "10" + 2 is the number 12 and in JavaScript it's the string "102", and I
// find both of those horrid and error-prone. In Miller, "10"+2 is MT_ERROR, by
// design, unless intentional casting is done like '$x=int("10")+2'.
//
// However, for dotting, in practice I tipped over and allowed dotting of
// strings and ints: so while "10" + 2 is an error in Miller, '"10". 2' is
// "102". Unlike with "+", with "." there is no ambiguity about what the output
// should be: always the string concatenation of the string representations of
// the two arguments. So, we do the string-cast for the user.

func dot_s_xx(input1, input2 *Mlrval) Mlrval {
	return MlrvalFromString(input1.String() + input2.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID   STRING    INT       FLOAT     BOOL      ARRAY     MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*VOID   */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*ARRAY  */ {_absn, _absn, _absn, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*MAP    */ {_absn, _absn, _absn, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
}

func MlrvalDot(input1, input2 *Mlrval) Mlrval {
	return dot_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// substr(s,m,n) gives substring of s from 1-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func MlrvalSubstr(output, input1, input2, input3 *Mlrval) {
	if !input1.IsStringOrVoid() {
		if input1.IsNumeric() {
			// JIT-stringify, if not already (e.g. intval scanned from string
			// in input-file data)
			input1.setPrintRep()
		} else {
			output.SetFromError()
			return
		}
	}
	// TODO: fix this with regard to UTF-8 and runes.
	strlen := int(len(input1.printrep))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !input2.IsInt() {
		output.SetFromError()
		return
	}
	lowerMindex := input2.intval

	upperMindex := strlen
	if input3.IsEmpty() {
		// Keep strlen
	} else if !input3.IsInt() {
		output.SetFromError()
		return
	} else {
		upperMindex = input3.intval
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := UnaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := UnaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		output.SetFromString("")
	} else if m > n {
		output.SetFromError()
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		output.SetFromString(input1.printrep[m : n+1])
	}
}

// ================================================================
func MlrvalTruncate(input1, input2 *Mlrval) Mlrval {
	if input1.IsErrorOrAbsent() {
		return *input1
	}
	if input2.IsErrorOrAbsent() {
		return *input2
	}
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input2.IsInt() {
		return MlrvalFromError()
	}
	if input2.intval < 0 {
		return MlrvalFromError()
	}

	oldLength := int(len(input1.printrep))
	maxLength := input2.intval
	if oldLength <= maxLength {
		return *input1
	} else {
		return MlrvalFromString(input1.printrep[0:maxLength])
	}
}

// ================================================================
func MlrvalLStrip(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(strings.TrimLeft(input1.printrep, " \t"))
	} else {
		output.CopyFrom(input1)
	}
}

func MlrvalRStrip(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(strings.TrimRight(input1.printrep, " \t"))
	} else {
		output.CopyFrom(input1)
	}
}

func MlrvalStrip(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(strings.Trim(input1.printrep, " \t"))
	} else {
		output.CopyFrom(input1)
	}
}

// ----------------------------------------------------------------
func MlrvalCollapseWhitespace(output, input1 *Mlrval) {
	MlrvalCollapseWhitespaceRegexp(output, input1, WhitespaceRegexp())
}

func MlrvalCollapseWhitespaceRegexp(output, input1 *Mlrval, whitespaceRegexp *regexp.Regexp) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(whitespaceRegexp.ReplaceAllString(input1.printrep, " "))
	} else {
		output.CopyFrom(input1)
	}
}

func WhitespaceRegexp() *regexp.Regexp {
	return regexp.MustCompile("\\s+")
}

// ================================================================
func MlrvalToUpper(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(strings.ToUpper(input1.printrep))
	} else if input1.mvtype == MT_VOID {
		output.CopyFrom(input1)
	} else {
		output.CopyFrom(input1)
	}
}

func MlrvalToLower(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		output.SetFromString(strings.ToLower(input1.printrep))
	} else if input1.mvtype == MT_VOID {
		output.CopyFrom(input1)
	} else {
		output.CopyFrom(input1)
	}
}

func MlrvalCapitalize(output, input1 *Mlrval) {
	if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			output.CopyFrom(input1)
		} else {
			runes := []rune(input1.printrep)
			rfirst := runes[0]
			rrest := runes[1:]
			sfirst := strings.ToUpper(string(rfirst))
			srest := string(rrest)
			output.SetFromString(sfirst + srest)
		}
	} else {
		output.CopyFrom(input1)
	}
}

// ----------------------------------------------------------------
func MlrvalCleanWhitespace(output, input1 *Mlrval) {
	MlrvalCollapseWhitespaceRegexp(output, input1, WhitespaceRegexp())
	MlrvalStrip(output, output)
}

// ================================================================
func MlrvalHexfmt(output, input1 *Mlrval) {
	if input1.mvtype == MT_INT {
		output.SetFromString("0x" + strconv.FormatUint(uint64(input1.intval), 16))
	} else {
		output.CopyFrom(input1)
	}
}

// ----------------------------------------------------------------
func fmtnum_is(input1, input2 *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			input2.printrep,
			input1.intval,
		),
	)
}

func fmtnum_fs(input1, input2 *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			input2.printrep,
			input1.floatval,
		),
	)
}

func fmtnum_bs(input1, input2 *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			input2.printrep,
			lib.BoolToInt(input1.boolval),
		),
	)
}

var fmtnum_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _erro, _absn, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*STRING */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _absn, _erro, fmtnum_is, _erro, _erro, _erro, _erro, _erro},
	/*FLOAT  */ {_erro, _absn, _erro, fmtnum_fs, _erro, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _absn, _erro, fmtnum_bs, _erro, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*MAP    */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalFmtNum(input1, input2 *Mlrval) Mlrval {
	return fmtnum_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
