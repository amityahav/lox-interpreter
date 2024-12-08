package main

import (
	"fmt"
	"reflect"
	"time"
)

type Expression interface {
	Eval(env *Environment) (interface{}, error)
}

type NilExpr struct{}

func (ne *NilExpr) Eval(_ *Environment) (interface{}, error) { return nil, nil }

func (ne *NilExpr) String() string { return "nil" }

type LiteralExpr struct {
	Literal interface{}
	Line    int
}

func (le *LiteralExpr) Eval(_ *Environment) (interface{}, error) {
	return le.Literal, nil
}

func (le *LiteralExpr) String() string {
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

func (ue *UnaryExpr) Eval(env *Environment) (interface{}, error) {
	val, err := ue.Expr.Eval(env)
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
	return fmt.Sprintf("(%s %v)", ue.Unary, ue.Expr)
}

type BinaryExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
	Line      int
}

func (be *BinaryExpr) Eval(env *Environment) (interface{}, error) {
	leftVal, err := be.LeftExpr.Eval(env)
	if err != nil {
		return nil, err
	}

	rightVal, err := be.RightExpr.Eval(env)
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

func (le *LogicalExpr) Eval(env *Environment) (interface{}, error) {
	switch TokenType(le.Operator) {
	case OR:
		lv, err := le.LeftExpr.Eval(env)
		if err != nil {
			return nil, err
		}

		if isTrue(lv) {
			return lv, nil
		}

		rv, err := le.RightExpr.Eval(env)
		if err != nil {
			return nil, err
		}

		return rv, nil
	case AND:
		lv, err := le.LeftExpr.Eval(env)
		if err != nil {
			return nil, err
		}

		if !isTrue(lv) {
			return lv, nil
		}

		rv, err := le.RightExpr.Eval(env)
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

func (ge *GroupingExpr) Eval(env *Environment) (interface{}, error) {
	return ge.Expr.Eval(env)
}

func (ge *GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.Expr)
}

type IdentifierExpr struct {
	Name string
	Line int
}

func (id *IdentifierExpr) Eval(env *Environment) (interface{}, error) {
	varEnv, ok := env.Lookup(id.Name)
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.\n[line %d]", id.Name, id.Line)
	}

	return varEnv.Bindings[id.Name], nil
}

func (id *IdentifierExpr) String() string {
	return id.Name
}

type AssignmentExpr struct {
	Name string
	Expr Expression
	Line int
}

func (as *AssignmentExpr) Eval(env *Environment) (interface{}, error) {
	varEnv, ok := env.Lookup(as.Name)
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.\n[line %d]", as.Name, as.Line)
	}

	val, err := as.Expr.Eval(env)
	if err != nil {
		return nil, err
	}

	varEnv.SetBinding(as.Name, val)

	return val, nil
}

func (as *AssignmentExpr) String() string {
	return as.Name
}

type CallExpr struct {
	Callee Expression
	Args   []Expression
	Line   int
}

func (c *CallExpr) Eval(env *Environment) (interface{}, error) {
	val, err := c.Callee.Eval(env)
	if err != nil {
		return nil, err
	}

	caller, ok := val.(Caller)
	if !ok {
		return nil, fmt.Errorf("Can only call functions and classes.\n[line %d]", c.Line)
	}

	arity := caller.Arity()
	if arity != len(c.Args) {
		return nil, fmt.Errorf("Expected %d arguments but got %d.\n[line %d]", arity, len(c.Args), c.Line)
	}

	var as []interface{}

	for _, arg := range c.Args {
		v, err := arg.Eval(env)
		if err != nil {
			return nil, err
		}

		as = append(as, v)
	}

	return caller.Call(as...)
}

type Statement interface {
	Execute(env *Environment) (interface{}, error)
}

type NilStmt struct{}

func (ns *NilStmt) Execute(_ *Environment) (interface{}, error) { return nil, nil }

type FunDeclStmt struct {
	Name   string
	Params []IdentifierExpr
	Body   Statement
}

func (f *FunDeclStmt) Execute(env *Environment) (interface{}, error) {
	fc := FunCaller{
		Name:    f.Name,
		Params:  f.Params,
		Body:    f.Body,
		closure: env,
	}

	env.SetBinding(f.Name, &fc)

	return nil, nil
}

type VarDeclStmt struct {
	Name string
	Expr Expression
}

func (v *VarDeclStmt) Execute(env *Environment) (interface{}, error) {
	val, err := v.Expr.Eval(env)
	if err != nil {
		return nil, err
	}

	env.SetBinding(v.Name, val)

	return nil, nil
}

type ExprStmt struct {
	Expr Expression
}

func (es *ExprStmt) Execute(env *Environment) (interface{}, error) {
	return es.Expr.Eval(env)
}

type PrintStmt struct {
	Expr Expression
}

func (ps *PrintStmt) Execute(env *Environment) (interface{}, error) {
	val, err := ps.Expr.Eval(env)
	if err != nil {
		return nil, err
	}

	fmt.Println(val)

	return nil, nil
}

type BlockStatement struct {
	Stmts []Statement
}

func (b *BlockStatement) Execute(env *Environment) (interface{}, error) {
	localEnv := ExpandEnv(env)

	for _, stmt := range b.Stmts {
		_, err := stmt.Execute(localEnv)
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

func (is *IfStmt) Execute(env *Environment) (interface{}, error) {
	cond, err := is.Condition.Eval(env)
	if err != nil {
		return nil, err
	}

	if isTrue(cond) {
		return is.Then.Execute(env)
	}

	return is.Else.Execute(env)
}

type WhileStmt struct {
	Condition Expression
	Body      Statement
}

func (ws *WhileStmt) Execute(env *Environment) (interface{}, error) {
	for {
		expr, err := ws.Condition.Eval(env)
		if err != nil {
			return nil, err
		}

		if !isTrue(expr) {
			return nil, nil
		}

		_, err = ws.Body.Execute(env)
		if err != nil {
			return nil, err
		}
	}
}

type ReturnValue any

type ReturnStmt struct {
	Expr Expression
}

func (rs *ReturnStmt) Execute(env *Environment) (interface{}, error) {
	val, err := rs.Expr.Eval(env)
	if err != nil {
		return nil, err
	}

	// panic is used here in order to quickly unwind the interpreter back to the code
	// that started executing the body.
	panic(ReturnValue(val))
}

type Caller interface {
	Call(args ...interface{}) (interface{}, error)
	Arity() int
}

type NativeClock struct{}

func (nc *NativeClock) Call(_ ...interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (nc *NativeClock) Arity() int { return 0 }

func (nc *NativeClock) String() string {
	return "<native fn>"
}

type FunCaller struct {
	Name   string
	Params []IdentifierExpr
	Body   Statement

	closure *Environment
}

func (fc *FunCaller) Call(args ...interface{}) (ret interface{}, err error) {
	localEnv := ExpandEnv(fc.closure)

	for i := 0; i < len(fc.Params); i++ {
		localEnv.SetBinding(fc.Params[i].Name, args[i])
	}

	defer func() {
		if res := recover(); res != nil {
			if rv, ok := res.(ReturnValue); ok {
				ret = rv
				return
			}

			panic(res)
		}
	}()

	_, err = fc.Body.Execute(localEnv)
	return
}

func (fc *FunCaller) Arity() int { return len(fc.Params) }

func (fc *FunCaller) String() string {
	return fmt.Sprintf("<fn %s>", fc.Name)
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
