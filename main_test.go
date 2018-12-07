package main

import (
	"reflect"
	"testing"
)

func TestMergeStrArrays(t *testing.T){
	a := make([]string, 2)
	a[0] = "a string"
	a[1] = "array"
	b := make([]string, 2)
	b[0] = "has"
	b[1] = "been merged"

	result := MergeStrArrays(a,b)

	expectedResult := []string{"a string", "array", "has", "been merged"}

	if(!reflect.DeepEqual(result, expectedResult)){
		t.Errorf("Merged array did not match sum of input arrays, expected %v, got %v", expectedResult, result)
	}
}

func TestContainsStr(t *testing.T) {
	inputArr := []string{"foo", "bar", "lorem", "ipsum"}
	a := "bar"

	if(!ContainsStr(inputArr, a)){
		t.Errorf("String was not found in array, despite array containing string, expected true, got false")
	}

	b := "not_in_array"

	if(ContainsStr(inputArr, b)){
		t.Errorf("String was found in array, but array does not contain string, expected false, got true")
	}
}

func TestContainsIdentifier(t *testing.T) {
	inputArr := []string{"foo", "bar", "lorem", "ipsum"}
	a := "bar"

	if(!ContainsIdentifier(inputArr, a)){
		t.Errorf("Identifier was not found in array, despite array containing identifier, expected true, got false")
	}

	b := "not_in_array"

	if(ContainsIdentifier(inputArr, b)){
		t.Errorf("Identifier was found in array, but array does not contain identfier, expected false, got true")
	}

	c := "*"

	if(!ContainsIdentifier(inputArr, c)){
		t.Errorf("ContainsIdentifier did not return true when wildcard used, expected true, got false")
	}
}

func TestRuleApplies(t *testing.T) {
	testExpressionRule := Rule{Type:"Expression", ID:"eval"}
	testExpression := Node{Class:"Expression", Type:"ExpressionStatement"}
	evalIdentifier := Node{Type:"Identifier", ID: "eval"}
	notEvalIdentifier := Node{Type:"Identifier", ID:"not_eval"}
	testExpression.Children = make([]Node, 1)
	testExpression.Children[0] = evalIdentifier

	if(!RuleApplies(testExpressionRule, testExpression)){
		t.Errorf("Rule should apply, expected true, got false")
	}

	testExpression.Children[0] = notEvalIdentifier
	if(RuleApplies(testExpressionRule, testExpression)){
		t.Errorf("Rule should not apply, expected false, got true")
	}

	testExpressionRule.ID = "*"

	if(!RuleApplies(testExpressionRule, testExpression)){
		t.Errorf("Rule should always apply when id is wildcard, expected true, got false")
	}

	testPropertyDoesNotExistRule := Rule{Type:"PropertyDoesNotExist"}
	variables = make(map[string]Variable)
	variables["test"] = Variable{Properties:make([]string, 1)}
	variables["test"].Properties[0] = "test_prop"
	testNode := Node{Type:"DotExpression", Children:make([]Node, 1), ID: "test_prop"}
	testNode.Children[0] = Node{Type:"Identifier", ID:"test"}

	if(RuleApplies(testPropertyDoesNotExistRule, testNode)){
		t.Errorf("Rule should not apply as property exists, expected false, got true")
	}

	testNode.ID = "does_not_exist"

	if(!RuleApplies(testPropertyDoesNotExistRule, testNode)){
		t.Errorf("Rule should apply as property does not exist, expected true, got false")
	}
}