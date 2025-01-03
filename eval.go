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

type ObjectGetExpr struct {
	Object Expression
	Prop   string
	Line   int
}

func (o *ObjectGetExpr) Eval(env *Environment) (interface{}, error) {
	val, err := o.Object.Eval(env)
	if err != nil {
		return nil, err
	}

	obj, ok := val.(*ClassInstance)
	if !ok {
		return nil, fmt.Errorf("Invalid operation, %v not an instance of an object.\n[line %d]", val, o.Line)
	}

	m, ok := obj.findMethod(o.Prop)
	if ok {
		return m, nil
	}

	m, ok = obj.Properties[o.Prop]
	if !ok {
		return nil, fmt.Errorf("Object %s has no property called %s\n[line %d]", obj.Name, o.Prop, o.Line)
	}

	return m, nil
}

type ObjectSetExpr struct {
	Object Expression
	Prop   string
	Expr   Expression
	Line   int
}

func (o *ObjectSetExpr) Eval(env *Environment) (interface{}, error) {
	val, err := o.Object.Eval(env)
	if err != nil {
		return nil, err
	}

	obj, ok := val.(*ClassInstance)
	if !ok {
		return nil, fmt.Errorf("Invalid operation, %v not an instance of an object.\n[line %d]", val, o.Line)
	}

	_, found := obj.findMethod(o.Prop)
	if found {
		return nil, fmt.Errorf("Invalid operation, cant set a method %s of object %s\n[line %d]", o.Prop, obj.Name, o.Line)
	}

	val, err = o.Expr.Eval(env)
	if err != nil {
		return nil, err
	}

	obj.Properties[o.Prop] = val

	return nil, nil
}

type Statement interface {
	Execute(env *Environment) (interface{}, error)
}

type NilStmt struct{}

func (ns *NilStmt) Execute(_ *Environment) (interface{}, error) { return nil, nil }

type ClassDeclStmt struct {
	Name       string
	SuperClass *IdentifierExpr
	Methods    []*FunDeclStmt
}

func (c *ClassDeclStmt) Execute(env *Environment) (interface{}, error) {
	cc := ClassCaller{
		Name:    c.Name,
		Methods: c.Methods,
		closure: env,
	}

	if c.SuperClass != nil {
		sc, err := c.SuperClass.Eval(env)
		if err != nil {
			return nil, err
		}

		v, ok := sc.(*ClassCaller)
		if !ok {
			return nil, fmt.Errorf("%s must be of class type.\n[line %d]", c.SuperClass.Name, c.SuperClass.Line)
		}

		cc.SuperClass = v
	}

	for _, m := range cc.Methods {
		if m.Name == "init" {
			cc.arity = len(m.Params)
			break
		}
	}

	env.SetBinding(c.Name, &cc)

	return nil, nil
}

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

	return &fc, nil
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

type BlockStmt struct {
	Stmts []Statement
}

func (b *BlockStmt) Execute(env *Environment) (interface{}, error) {
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

type ReturnValue struct {
	Value interface{}
}

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
	panic(&ReturnValue{Value: val})
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

type ClassInstance struct {
	Name               string
	SuperClassInstance *ClassInstance
	Methods            map[string]interface{}
	Properties         map[string]interface{}
}

func (ci *ClassInstance) findMethod(name string) (interface{}, bool) {
	m, ok := ci.Methods[name]
	if ok {
		return m, true
	}

	if ci.SuperClassInstance != nil {
		return ci.SuperClassInstance.findMethod(name)
	}

	return nil, false
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("%s instance", ci.Name)
}

type ClassCaller struct {
	Name       string
	SuperClass *ClassCaller
	Methods    []*FunDeclStmt

	arity   int
	closure *Environment
}

func (ci *ClassInstance) initSuperClass(env *Environment, superClass *ClassCaller) (*ClassInstance, error) {
	if superClass == nil {
		return nil, nil
	}

	sci := ClassInstance{
		Name:    superClass.Name,
		Methods: make(map[string]interface{}),
	}

	i, err := sci.initSuperClass(env, superClass.SuperClass)
	if err != nil {
		return nil, err
	}

	localEnv := ExpandEnv(env)
	localEnv.SetBinding("super", nil)
	if i != nil {
		localEnv.SetBinding("super", i)
		sci.SuperClassInstance = i
	}

	for _, m := range superClass.Methods {
		val, err := m.Execute(localEnv)
		if err != nil {
			return nil, err
		}

		sci.Methods[m.Name] = val
	}

	return &sci, nil
}

func (cc *ClassCaller) Call(args ...interface{}) (interface{}, error) {
	ci := ClassInstance{
		Name:       cc.Name,
		Methods:    map[string]interface{}{},
		Properties: make(map[string]interface{}),
	}

	localEnv := ExpandEnv(cc.closure)

	// TODO: make sure params cannot shadow "this" keyword
	localEnv.SetBinding("this", &ci)

	sci, err := ci.initSuperClass(localEnv, cc.SuperClass)
	if err != nil {
		return nil, err
	}

	localEnv.SetBinding("super", nil)
	if sci != nil {
		localEnv.SetBinding("super", sci)
		ci.SuperClassInstance = sci
	}

	for _, m := range cc.Methods {
		val, err := m.Execute(localEnv)
		if err != nil {
			return nil, err
		}

		ci.Methods[m.Name] = val
	}

	if v, ok := ci.Methods["init"]; ok {
		initializer := v.(Caller)
		_, err := initializer.Call(args...)
		if err != nil {
			return nil, err
		}
	}

	return &ci, nil
}

func (cc *ClassCaller) Arity() int { return cc.arity }

func (cc *ClassCaller) String() string {
	return fmt.Sprintf("%s instance", cc.Name)
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
			if rv, ok := res.(*ReturnValue); ok {
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
