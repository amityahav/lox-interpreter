package main

import (
	"fmt"
	"reflect"
)

type Expression interface {
	Eval() (interface{}, error)
	String() string
}

type LiteralExpr struct {
	Literal interface{}
	Line    int
}

func (le *LiteralExpr) Eval() (interface{}, error) {
	return le.Literal, nil
}

func (le *LiteralExpr) String() string {
	if le.Literal == nil {
		return "nil"
	}

	if v, ok := le.Literal.(float64); ok {
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.1f", v)
		}
	}

	return fmt.Sprintf("%v", le.Literal)
}

type UnaryExpr struct {
	Unary string
	Expr  Expression
	Line  int
}

func (ue *UnaryExpr) Eval() (interface{}, error) {
	val, err := ue.Expr.Eval()
	if err != nil {
		return nil, err
	}

	switch TokenType(ue.Unary) {
	case MINUS:
		v, ok := val.(float64)
		if !ok {
			return nil, fmt.Errorf("Operand must be a number.\n[line %d]", ue.Line)
		}

		return -v, nil
	case BANG:
		return !isTrue(val), nil
	}

	// unreachable
	return nil, nil
}

func (ue *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", ue.Unary, ue.Expr.String())
}

type BinaryExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
	Line      int
}

func (be *BinaryExpr) Eval() (interface{}, error) {
	leftVal, err := be.LeftExpr.Eval()
	if err != nil {
		return nil, err
	}

	rightVal, err := be.RightExpr.Eval()
	if err != nil {
		return nil, err
	}

	switch TokenType(be.Operator) {
	case SLASH, STAR, MINUS, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			return nil, fmt.Errorf("Operands must be numbers.\n[line %d]", be.Line)
		}

		switch TokenType(be.Operator) {
		case SLASH:
			return lv / rv, nil
		case STAR:
			return lv * rv, nil
		case MINUS:
			return lv - rv, nil
		case LESS:
			return lv < rv, nil
		case LESS_EQUAL:
			return lv <= rv, nil
		case GREATER:
			return lv > rv, nil
		case GREATER_EQUAL:
			return lv >= rv, nil
		}
	case PLUS:
		t1 := reflect.TypeOf(leftVal).Kind()
		t2 := reflect.TypeOf(rightVal).Kind()

		if (t1 != reflect.Float64 && t1 != reflect.String) ||
			(t2 != reflect.Float64 && t2 != reflect.String) ||
			(t1 != t2) {
			return nil, fmt.Errorf("Operands must be two numbers or two strings.\n[line %d]", be.Line)
		}

		lv, ok := leftVal.(float64)
		if ok {
			rv, _ := rightVal.(float64)
			return lv + rv, nil
		}

		lvs, ok := leftVal.(string)
		if ok {
			rvs, _ := rightVal.(string)
			return lvs + rvs, nil
		}
	case EQUAL_EQUAL:
		return leftVal == rightVal, nil
	case BANG_EQUAL:
		return leftVal != rightVal, nil
	}

	// unreachable
	return nil, nil
}

func (be *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Operator, be.LeftExpr, be.RightExpr)
}

type GroupingExpr struct {
	Expr Expression
	Line int
}

func (ge *GroupingExpr) Eval() (interface{}, error) {
	return ge.Expr.Eval()
}

func (ge *GroupingExpr) String() string {
	return fmt.Sprintf("(group %s)", ge.Expr.String())
}

func isTrue(val interface{}) bool {
	if val == nil {
		return false
	}

	if v, ok := val.(bool); ok {
		return v
	}

	return true
}
