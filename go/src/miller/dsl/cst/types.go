package cst

import (
	"miller/lib"
)

// ================================================================
// Main type definitions for CST build/execute
// ================================================================

// ----------------------------------------------------------------
// There are three CST roots: begin-block, body-block, and end-block.
//
// Next-level items are:
// * srec assignments
// * oosvar assignments
// * localvar assignments
// * emit et al.
// * bare-boolean
// * break/continue/return
// * statement block (if-body, for-body, etc)
// ----------------------------------------------------------------

// ----------------------------------------------------------------
// AST nodes (TNodeType) at the moment:
//
// NodeTypeStringLiteral
// NodeTypeIntLiteral
// NodeTypeFloatLiteral
// NodeTypeBoolLiteral
//
// NodeTypeDirectFieldName
// NodeTypeIndirectFieldName
//
// NodeTypeStatementBlock
// NodeTypeAssignment
// NodeTypeOperator
// NodeTypeContextVariable
// ----------------------------------------------------------------

// ----------------------------------------------------------------
type State struct {
	Inrec   *lib.Lrec
	Context *lib.Context
	// oosvars too
	// stack frames will go into individual statement-block nodes
}

func NewState(
	inrec *lib.Lrec,
	context *lib.Context,
) *State {
	return &State{
		Inrec:   inrec,
		Context: context,
	}
}

// ----------------------------------------------------------------
type IExecutable interface {
	Execute(state *State)
}

// ----------------------------------------------------------------
type Root struct {
	// array of statements/blocks
	executables []IExecutable
}

// ----------------------------------------------------------------
type SrecDirectFieldAssignmentNode struct {
	lhsFieldName string
	rhs          IEvaluable
}

type IndirectSrecFieldAssignmentNode struct {
	lhsFieldName IEvaluable
	rhs          IEvaluable
}

// xxx implement IExecutable

type StatementBlockNode struct {
	// list of statement
}

// xxx implement IExecutable

// ================================================================
type IEvaluable interface {
	Evaluate(state *State) lib.Mlrval
}
