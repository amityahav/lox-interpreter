package main

import (
	"fmt"
	"reflect"
)

type Expression interface {
	Eval() (interface{}, error)
	String() string
}

type NoopExpr struct{}

func (ne *NoopExpr) Eval() (interface{}, error) { return nil, nil }

func (ne *NoopExpr) String() string { return "" }

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

type LogicalExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
}

func (le *LogicalExpr) Eval() (interface{}, error) {
	switch TokenType(le.Operator) {
	case OR:
		lv, err := le.LeftExpr.Eval()
		if err != nil {
			return nil, err
		}

		if isTrue(lv) {
			return lv, nil
		}

		rv, err := le.RightExpr.Eval()
		if err != nil {
			return nil, err
		}

		return rv, nil
	case AND:
		lv, err := le.LeftExpr.Eval()
		if err != nil {
			return nil, err
		}

		if !isTrue(lv) {
			return lv, nil
		}

		rv, err := le.RightExpr.Eval()
		if err != nil {
			return nil, err
		}

		return rv, nil
	}

	// unreachable
	return nil, nil
}

func (le *LogicalExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", le.LeftExpr, le.Operator, le.RightExpr)

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

type IdentifierExpr struct {
	Name string
	Line int

	state *Environment
}

func (id *IdentifierExpr) Eval() (interface{}, error) {
	scope, ok := id.state.GetScopeFor(id.Name)
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.\n[line %d]", id.Name, id.Line)
	}

	return scope.Bindings[id.Name], nil
}

func (id *IdentifierExpr) String() string {
	return id.Name
}

type AssignmentExpr struct {
	Name string
	Expr Expression
	Line int

	state *Environment
}

func (as *AssignmentExpr) Eval() (interface{}, error) {
	scope, ok := as.state.GetScopeFor(as.Name)
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.\n[line %d]", as.Name, as.Line)
	}

	val, err := as.Expr.Eval()
	if err != nil {
		return nil, err
	}

	scope.SetBinding(as.Name, val)

	return val, nil
}

func (as *AssignmentExpr) String() string {
	return as.Name
}

type CallExpr struct {
	Callee Expression
	Args   []Expression
}

func (c *CallExpr) Eval() (interface{}, error) {
	val, err := c.Callee.Eval()
	if err != nil {
		return nil, err
	}

	caller, ok := val.(Caller)
	if !ok {
		panic("not a function")
	}

	var as []interface{}

	for _, arg := range c.Args {
		v, err := arg.Eval()
		if err != nil {
			return nil, err
		}

		as = append(as, v)
	}

	return caller.Call(as...)
}

func (c *CallExpr) String() string {
	return ""
}

type Statement interface {
	Execute() (interface{}, error)
}

type NoopStmt struct{}

func (ns *NoopStmt) Execute() (interface{}, error) { return nil, nil }

type VarDeclStmt struct {
	Name string
	Expr Expression

	state *Environment
}

func (v *VarDeclStmt) Execute() (interface{}, error) {
	scope := v.state.GetInnermostScope()

	val, err := v.Expr.Eval()
	if err != nil {
		return nil, err
	}

	scope.SetBinding(v.Name, val)

	return nil, nil
}

type ExprStmt struct {
	Expr Expression
}

func (es *ExprStmt) Execute() (interface{}, error) {
	return es.Expr.Eval()
}

type PrintStmt struct {
	Expr Expression
}

func (ps *PrintStmt) Execute() (interface{}, error) {
	val, err := ps.Expr.Eval()
	if err != nil {
		return nil, err
	}

	if val == nil {
		fmt.Println("nil")
		return nil, nil
	}

	fmt.Println(val)

	return nil, nil
}

type BlockStatement struct {
	Stmts []Statement

	state *Environment
}

func (b *BlockStatement) Execute() (interface{}, error) {
	// start a new scope
	b.state.GrowScopes()
	defer func() {
		// close scope
		b.state.CloseInnermostScope()
	}()

	for _, stmt := range b.Stmts {
		_, err := stmt.Execute()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type IfStmt struct {
	Condition Expression
	Then      Statement
	Else      Statement
}

func (is *IfStmt) Execute() (interface{}, error) {
	cond, err := is.Condition.Eval()
	if err != nil {
		return nil, err
	}

	if isTrue(cond) {
		return is.Then.Execute()
	}

	return is.Else.Execute()
}

type WhileStmt struct {
	Condition Expression
	Body      Statement
}

func (ws *WhileStmt) Execute() (interface{}, error) {
	for {
		expr, err := ws.Condition.Eval()
		if err != nil {
			return nil, err
		}

		if !isTrue(expr) {
			return nil, nil
		}

		_, err = ws.Body.Execute()
		if err != nil {
			return nil, err
		}
	}
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
