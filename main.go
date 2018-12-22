package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"io/ioutil"
	"os"
	"reflect"
)

type Node struct {
	Type     string
	ID       string
	Children []Node
	Class    string
}

type Rules struct {
	All []Rule
}

type Rule struct {
	Type string
	Action RuleAction
	ID string
}

type RuleAction struct {
	Type string
	Info string
}

type Variable struct {
	Properties []string
}

type BasicBlock struct {
	Statements []Node
	Parameters []Node
}

type DAGNode struct {
	Parent *DAGNode
	Left *DAGNode
	Right *DAGNode
	Labels []string
}
var nodes []Node
var basicBlocks map[int]BasicBlock
var numBlocks int = 0
var rules Rules
var variables map[string]Variable
var lastVExpr Node = Node{}

func main() {
	dat, fError := ioutil.ReadFile(os.Args[1])
	CheckErr(fError)
	src := string(dat)
	program, err := parser.ParseFile(nil, "", src, 0)
	CheckErr(err)

	nodes = make([]Node, len(program.Body))
	basicBlocks = make(map[int]BasicBlock)
	variables = make(map[string]Variable)
	fmt.Println("Unknowns:")
	for i := 0; i < len(program.Body); i++ {
		nodes[i] = Walk(program.Body[i])
	}

	fmt.Println("Nodes:")
	for i := 0; i < len(nodes); i++ {
		fmt.Println(nodes[i])
	}

	fmt.Println("Variables:")
	fmt.Println(variables)

	fmt.Println("Rules:")
	rules = ParseRules(os.Args[2])
	fmt.Println(rules)

	for i := 0; i < len(nodes); i++ {
		CheckNode(nodes[i])
	}

	fmt.Println("Basic Blocks:")
	for i := 0; i < len(basicBlocks); i++ {
		fmt.Println(basicBlocks[i])
	}
	fmt.Println("A basic block:")
	ConstructDAG(basicBlocks[0])
}

// Apply action takes the corresponding action for a rule and executes it
// If the action is to fail, then the analyzer will print the error and exit,
// if the action is to warn, then the analyzer will print a warning and continue.
func ApplyAction(action RuleAction){
	switch action.Type {
	case "fail":
		fmt.Println(fmt.Errorf("\nError:\n\t%s\n", action.Info))
		os.Exit(1)
		break
	case "warn":
		fmt.Println(fmt.Errorf("\nWarning:\n\t%s\n", action.Info))
		break
	default:
		fmt.Println("Action not implemented")
		break
	}
}

// MergeStrArrays takes two arrays of strings and merges them together,
// before returning the result.
func MergeStrArrays(a []string, b []string) []string{
	var mergedArr []string = make([]string, len(a) + len(b))
	var mPos int = 0

	for i := 0; i < len(a); i++ {
		mergedArr[mPos] = a[i]
		mPos++
	}

	for i := 0; i < len(b); i++ {
		mergedArr[mPos] = b[i]
		mPos++
	}

	return mergedArr
}

// GetNodeIdentifiers recurses through the children of a node, adding the id of each
// node that corresponds to the Identifier type to an array of identifiers which is
// then returned.
func GetNodeIdentifiers(node Node) []string {
	if node.Type == "Identifier" {
		identifiers := make([]string, 1)
		identifiers[0] = node.ID
		return identifiers
	}

	var identifiers []string = make([]string, len(node.Children))
	for i := 0; i < len(node.Children); i++ {
		if node.Children[i].Type == "Identifier" {
			identifiers = MergeStrArrays(identifiers, GetNodeIdentifiers(node.Children[i]))
		}
	}

	return identifiers
}

// ContainsStr simply checks if a string is contained within an array of strings.
func ContainsStr(strings []string, str string) bool {
	for i := 0; i < len(strings); i++ {
		if strings[i] == str {
			return true
		}
	}

	return false
}

func ContainsIdentifier(identifiers []string, id string) bool {
	if id == "*" {
		return true
	}

	return ContainsStr(identifiers, id)
}

// RuleApplies checks if the specific conditions of a particular rule apply to the node in question.
func RuleApplies(rule Rule, node Node) bool{
	switch rule.Type {
	case "Expression":
		return node.Class == "Expression" && ContainsIdentifier(GetNodeIdentifiers(node), rule.ID)
		break
	case "PropertyDoesNotExist":
		if node.Type != "DotExpression" {
			return false
		}

		key := node.Children[0].ID

		pVar, ok := variables[key]
		return ok && !ContainsStr(pVar.Properties, node.ID)
		break
	}
	return false
}

// CheckNode looks to see whether any of the rules parsed apply to a node or its' children.
func CheckNode(node Node){
	for i := 0; i < len(rules.All); i++ {
		if RuleApplies(rules.All[i], node) {
			ApplyAction(rules.All[i].Action)
		}
	}

	for i := 0; i < len(node.Children); i++ {
		CheckNode(node.Children[i])
	}
}

// CheckErr is a utility function, it just panics if an error is not nil.
func CheckErr(e error){
	if e != nil {
		panic(e)
	}
}

// ParseRules reads the rules file (should be in JSON format) and attempts to parse it into a set of rules.
func ParseRules(fileName string) Rules{
	dat, fError := ioutil.ReadFile(fileName)
	CheckErr(fError)
	rulesStr := string(dat)
	var rules Rules
	rules.All = make([]Rule, 0)
	err := json.Unmarshal([]byte(rulesStr), &rules.All)
	CheckErr(err)
	return rules
}

// CreateVariableProp adds a named property to a named variable stored within the global
// variables table, if the named variable exists within the table.
func CreateVariableProp(node Node) {
	var key string = node.Children[0].ID

	if node.Type == "ObjectProperty" {
		key = lastVExpr.ID
	}

	parentVar, ok := variables[key]

	if ok {
		parentVar.Properties = append(parentVar.Properties, node.ID)
		variables[key] = parentVar
	}
}

// Walk walks the AST for the script and attempts to decorate each node so we can better
// analyse the program.
func Walk(n interface{}) Node {
	var node Node
	switch reflect.TypeOf(n) {
	case reflect.TypeOf(&ast.VariableStatement{}):
		v, _ := n.(*ast.VariableStatement)
		node = Node{Type: "VariableStatement"}
		node.Children = make([]Node, len(v.List))
		for i := 0; i < len(v.List); i++ {
			node.Children[i] = Walk(v.List[i])
		}
		break
	case reflect.TypeOf(&ast.VariableExpression{}):
		v, _ := n.(*ast.VariableExpression)
		node = Node{Type: "VariableExpression", ID: v.Name, Class: "Expression"}
		variables[v.Name] = Variable{Properties:make([]string, 1)}
		lastVExpr = node
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Initializer)
		break
	case reflect.TypeOf(&ast.DotExpression{}):
		v, _ := n.(*ast.DotExpression)
		node = Node{Type: "DotExpression", ID: v.Identifier.Name, Class: "Expression"}
		if v.Left != nil {
			node.Children = make([]Node, 1)
			node.Children[0] = Walk(v.Left)
		}
		break
	case reflect.TypeOf(&ast.CallExpression{}):
		v, _ := n.(*ast.CallExpression)
		node = Node{Type: "CallExpression", Class: "Expression"}
		node.Children = make([]Node, len(v.ArgumentList)+1)
		for i := 0; i < len(v.ArgumentList); i++ {
			node.Children[i] = Walk(v.ArgumentList[i])
		}
		node.Children[len(node.Children)-1] = Walk(v.Callee)
		break
	case reflect.TypeOf(&ast.StringLiteral{}):
		v, _ := n.(*ast.StringLiteral)
		node = Node{Type: "StringLiteral", ID: v.Value}
		break
	case reflect.TypeOf(&ast.Identifier{}):
		v, _ := n.(*ast.Identifier)
		node = Node{Type: "Identifier", ID: v.Name}
		break
	case reflect.TypeOf(&ast.ExpressionStatement{}):
		v, _ := n.(*ast.ExpressionStatement)
		node = Node{Type: "ExpressionStatement", Class: "Expression"}
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Expression)
		break
	case reflect.TypeOf(&ast.FunctionLiteral{}):
		v, _ := n.(*ast.FunctionLiteral)
		node = Node{Type: "FunctionLiteral"}
		if v.Name != nil {
			node.ID = v.Name.Name
		}
		node.Children = make([]Node, len(v.ParameterList.List) + len(v.DeclarationList))
		basicBlocks[numBlocks] = BasicBlock{Parameters:make([]Node, len(v.ParameterList.List)), Statements:make([]Node, len(v.DeclarationList))}
		var n int = 0
		for i := 0; i < len(node.Children); i++ {
			if i < len(node.Children) / 2 - 1  && i < len(v.ParameterList.List){
				paramNode := Walk(v.ParameterList.List[i])
				node.Children[i] = paramNode
				basicBlocks[numBlocks].Parameters[i] = paramNode

			} else {
				declNode := Walk(v.DeclarationList[n])
				node.Children[i] = declNode
				basicBlocks[numBlocks].Statements[n] = declNode
				n++
			}
		}
		numBlocks++
		break
	case reflect.TypeOf(&ast.VariableDeclaration{}):
		v, _ := n.(*ast.VariableDeclaration)
		node = Node{Type: "VariableDeclaration"}
		node.Children = make([]Node, len(v.List))
		for i := 0; i < len(v.List); i++ {
			node.Children[i] = Walk(v.List[i])
		}
		break
	case reflect.TypeOf(&ast.NumberLiteral{}):
		v, _ := n.(*ast.NumberLiteral)
		node = Node{Type: "NumberLiteral", ID: v.Literal}
		break
	case reflect.TypeOf(&ast.BinaryExpression{}):
		v, _ := n.(*ast.BinaryExpression)
		node = Node{Type: "BinaryExpression", ID: v.Operator.String(), Class: "Expression"}
		node.Children = make([]Node, 2)
		node.Children[0] = Walk(v.Left)
		node.Children[1] = Walk(v.Right)
		break
	case reflect.TypeOf(&ast.ObjectLiteral{}):
		v, _ := n.(*ast.ObjectLiteral)
		node = Node{Type: "ObjectLiteral"}
		node.Children = make([]Node, len(v.Value))
		for i := 0; i < len(v.Value); i++ {
			node.Children[i] = Walk(v.Value[i])
			if node.Children[i].Type == "ObjectProperty" {
				CreateVariableProp(node.Children[i])
			}
		}
		break
	case reflect.TypeOf(&ast.IfStatement{}):
		v, _ := n.(*ast.IfStatement)
		node = Node{Type: "IfStatement"}

		if v.Alternate != nil {
			node.Children = make([]Node, 2)
			node.Children[1] = Walk(v.Alternate)
		} else {
			node.Children = make([]Node, 1)
		}

		node.Children[0] = Walk(v.Consequent)
		break
	case reflect.TypeOf(&ast.BlockStatement{}):
		v, _ := n.(*ast.BlockStatement)
		node = Node{Type: "BlockStatement"}
		node.Children = make([]Node, len(v.List))
		for i:= 0; i < len(v.List); i++ {
			node.Children[i] = Walk(v.List[i])
		}
		break
	case reflect.TypeOf(ast.Property{}):
		v, _ := n.(ast.Property)
		node = Node{Type: "ObjectProperty", ID: v.Key}
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Value)
		break
	case reflect.TypeOf(&ast.AssignExpression{}):
		v, _ := n.(*ast.AssignExpression)
		node = Node{Type:"AssignExpression", ID:v.Operator.String(), Class:"Expression"}
		node.Children = make([]Node, 2)
		node.Children[0] = Walk(v.Left)
		node.Children[1] = Walk(v.Right)
		if node.Children[0].Type == "DotExpression" {
			CreateVariableProp(node.Children[0])
		}
		break
	case reflect.TypeOf(&ast.FunctionStatement{}):
		v, _ := n.(*ast.FunctionStatement)
		node = Node{Type:"FunctionStatement", Children:make([]Node, 1)}
		node.Children[0] = Walk(v.Function)
		break
	default:
		fmt.Println(reflect.TypeOf(n))
		node = Node{Type: "Unknown"}
	}
	return node
}

func GetIDs(stmt Node) []string {
	var toReturn []string = make([]string, 1)
	toReturn[0] = stmt.ID
	numChildren := len(stmt.Children)
	if(numChildren > 0){

		for i:=0; i < numChildren; i++ {
			toReturn = MergeStrArrays(toReturn, GetIDs(stmt.Children[i]))
		}
	}

	return toReturn
}

func ConstructDAG(block BasicBlock){
	for i:=0; i < len(block.Statements); i++ {
		fmt.Println(block.Statements[i])
		fmt.Println(GetIDs(block.Statements[i]))
	}
}