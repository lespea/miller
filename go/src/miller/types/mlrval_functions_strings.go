package types

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ================================================================
func MlrvalStrlen(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromInt64(int64(utf8.RuneCountInString(ma.printrep)))
}

// ================================================================
func MlrvalToString(ma *Mlrval) Mlrval {
	return MlrvalFromString(ma.String())
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

func dot_s_xx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(ma.String() + mb.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*VOID   */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDot(ma, mb *Mlrval) Mlrval {
	return dot_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// substr(s,m,n) gives substring of s from 1-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func MlrvalSubstr(ma, mb, mc *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsInt() {
		return MlrvalFromError()
	}
	if !mc.IsInt() {
		return MlrvalFromError()
	}
	strlen := int64(len(ma.printrep))

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := unaliasArrayLengthIndex(strlen, mb.intval)
	n, nok := unaliasArrayLengthIndex(strlen, mc.intval)

	if !mok || !nok {
		return MlrvalFromString("")
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive.
		return MlrvalFromString(ma.printrep[m : n+1])
	}
}

// ================================================================
func MlrvalSsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mc.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		strings.Replace(ma.printrep, mb.printrep, mc.printrep, 1),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalGsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mc.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := regexp.MustCompile(mb.printrep)
	return MlrvalFromString(
		re.ReplaceAllString(ma.printrep, mc.printrep),
	)
}

// ================================================================
func MlrvalTruncate(ma, mb *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsInt() {
		return MlrvalFromError()
	}
	if mb.intval < 0 {
		return MlrvalFromError()
	}

	oldLength := int64(len(ma.printrep))
	maxLength := mb.intval
	if oldLength <= maxLength {
		return *ma
	} else {
		return MlrvalFromString(ma.printrep[0:maxLength])
	}
}

// ================================================================
func MlrvalLStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimLeft(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

func MlrvalRStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimRight(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

func MlrvalStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.Trim(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalCollapseWhitespace(ma *Mlrval) Mlrval {
	return MlrvalCollapseWhitespaceRegexp(ma, WhitespaceRegexp())
}

func MlrvalCollapseWhitespaceRegexp(ma *Mlrval, whitespaceRegexp *regexp.Regexp) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(whitespaceRegexp.ReplaceAllString(ma.printrep, " "))
	} else {
		return *ma
	}
}

func WhitespaceRegexp() *regexp.Regexp {
	return regexp.MustCompile("\\s+")
}

// ================================================================
func MlrvalToUpper(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToUpper(ma.printrep))
	} else if ma.mvtype == MT_VOID {
		return *ma
	} else {
		return *ma
	}
}

func MlrvalToLower(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToLower(ma.printrep))
	} else if ma.mvtype == MT_VOID {
		return *ma
	} else {
		return *ma
	}
}

func MlrvalCapitalize(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		if ma.printrep == "" {
			return *ma
		} else {
			runes := []rune(ma.printrep)
			rfirst := runes[0]
			rrest := runes[1:]
			sfirst := strings.ToUpper(string(rfirst))
			srest := string(rrest)
			return MlrvalFromString(sfirst + srest)
		}
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalCleanWhitespace(ma *Mlrval) Mlrval {
	temp := MlrvalCollapseWhitespaceRegexp(ma, WhitespaceRegexp())
	return MlrvalStrip(&temp)
}

// ================================================================
func MlrvalHexfmt(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_INT {
		return MlrvalFromString("0x" + strconv.FormatUint(uint64(ma.intval), 16))
	} else {
		return *ma
	}
}
