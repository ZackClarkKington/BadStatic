package main

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"reflect"
)

type Node struct {
	Type     string
	ID       string
	Children []Node
}

type Op struct {
	Operand Node
	Type    string
	Label   string
}

var nodes []Node

func main() {
	src := `
		var b = document.getElementById('a').value;
		eval(b);
`
	program, err := parser.ParseFile(nil, "", src, 0)
	nodes = make([]Node, len(program.Body))
	for i := 0; i < len(program.Body); i++ {
		//fmt.Println(reflect.TypeOf(program.Body[i]))
		nodes[i] = Walk(program.Body[i])
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Nodes:")
	for i := 0; i < len(nodes); i++ {
		fmt.Println(nodes[i])
	}
}

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
		node = Node{Type: "VariableExpression", ID: v.Name}
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Initializer)
		break
	case reflect.TypeOf(&ast.DotExpression{}):
		v, _ := n.(*ast.DotExpression)
		node = Node{Type: "DotExpression", ID: v.Identifier.Name}
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Left)
		break
	case reflect.TypeOf(&ast.CallExpression{}):
		v, _ := n.(*ast.CallExpression)
		node = Node{Type: "CallExpression"}
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
	case reflect.TypeOf(&ast.ExpressionStatement{}):
		v, _ := n.(*ast.ExpressionStatement)
		node = Node{Type: "ExpressionStatement"}
		node.Children = make([]Node, 1)
		node.Children[0] = Walk(v.Expression)
	default:
		fmt.Println(reflect.TypeOf(n))
		node = Node{Type: "Unknown"}
	}
	return node
}
