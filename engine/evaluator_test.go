package engine

import (
	"math"
	"testing"
)

func newTestEvaluator() (*Evaluator, *DataModel) {
	dm := NewDataModel()
	return NewEvaluator(dm), dm
}

func TestEvalConcat(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("concat", []any{"hello", " ", "world"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello world" {
		t.Errorf("concat = %q, want 'hello world'", result)
	}
}

func TestEvalConcatWithPath(t *testing.T) {
	eval, dm := newTestEvaluator()
	dm.Set("/name", "Alice")
	result, err := eval.Eval("concat", []any{
		"Hello, ",
		map[string]any{"path": "/name"},
		"!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Hello, Alice!" {
		t.Errorf("concat = %q, want 'Hello, Alice!'", result)
	}
}

func TestEvalFormat(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("format", []any{"Hello, {0}! You are {1}.", "Alice", "great"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Hello, Alice! You are great." {
		t.Errorf("format = %q", result)
	}
}

func TestEvalToUpperCase(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("toUpperCase", []any{"hello"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "HELLO" {
		t.Errorf("toUpperCase = %q", result)
	}
}

func TestEvalToLowerCase(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("toLowerCase", []any{"HELLO"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("toLowerCase = %q", result)
	}
}

func TestEvalTrim(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("trim", []any{"  hello  "})
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("trim = %q", result)
	}
}

func TestEvalSubstring(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("substring", []any{"hello world", float64(0), float64(5)})
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("substring = %q, want 'hello'", result)
	}
}

func TestEvalSubstringNoEnd(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("substring", []any{"hello world", float64(6)})
	if err != nil {
		t.Fatal(err)
	}
	if result != "world" {
		t.Errorf("substring = %q, want 'world'", result)
	}
}

func TestEvalLength(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("length", []any{"hello"})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(5) {
		t.Errorf("length = %v, want 5", result)
	}
}

func TestEvalAdd(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("add", []any{float64(2), float64(3)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(5) {
		t.Errorf("add = %v, want 5", result)
	}
}

func TestEvalSubtract(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("subtract", []any{float64(10), float64(3)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(7) {
		t.Errorf("subtract = %v, want 7", result)
	}
}

func TestEvalMultiply(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("multiply", []any{float64(4), float64(5)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(20) {
		t.Errorf("multiply = %v, want 20", result)
	}
}

func TestEvalDivide(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("divide", []any{float64(10), float64(2)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(5) {
		t.Errorf("divide = %v, want 5", result)
	}
}

func TestEvalDivideByZero(t *testing.T) {
	eval, _ := newTestEvaluator()
	_, err := eval.Eval("divide", []any{float64(10), float64(0)})
	if err == nil {
		t.Error("expected division by zero error")
	}
}

func TestEvalEquals(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("equals", []any{"hello", "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("equals = %v, want true", result)
	}

	result, err = eval.Eval("equals", []any{"hello", "world"})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("equals = %v, want false", result)
	}
}

func TestEvalGreaterThan(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("greaterThan", []any{float64(5), float64(3)})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("greaterThan(5,3) = %v, want true", result)
	}

	result, err = eval.Eval("greaterThan", []any{float64(2), float64(3)})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("greaterThan(2,3) = %v, want false", result)
	}
}

func TestEvalNot(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("not", []any{true})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("not(true) = %v, want false", result)
	}

	result, err = eval.Eval("not", []any{false})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("not(false) = %v, want true", result)
	}
}

func TestEvalNestedFunctionCall(t *testing.T) {
	eval, _ := newTestEvaluator()
	// concat(toUpperCase("hello"), " ", "world")
	result, err := eval.Eval("concat", []any{
		map[string]any{
			"functionCall": map[string]any{
				"name": "toUpperCase",
				"args": []any{"hello"},
			},
		},
		" ",
		"world",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "HELLO world" {
		t.Errorf("nested = %q, want 'HELLO world'", result)
	}
}

func TestEvalNestedWithPath(t *testing.T) {
	eval, dm := newTestEvaluator()
	dm.Set("/name", "alice")
	// concat(toUpperCase({path: "/name"}), "!")
	result, err := eval.Eval("concat", []any{
		map[string]any{
			"functionCall": map[string]any{
				"name": "toUpperCase",
				"args": []any{
					map[string]any{"path": "/name"},
				},
			},
		},
		"!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "ALICE!" {
		t.Errorf("nested = %q, want 'ALICE!'", result)
	}
}

func TestEvalUnknownFunction(t *testing.T) {
	eval, _ := newTestEvaluator()
	_, err := eval.Eval("bogusFunc", []any{"a"})
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

func TestEvalMissingPath(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("concat", []any{
		"Hello, ",
		map[string]any{"path": "/nonexistent"},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Missing path resolves to empty string
	if result != "Hello, " {
		t.Errorf("missing path = %q, want 'Hello, '", result)
	}
}

func TestPathsInArgs(t *testing.T) {
	args := []any{
		"literal",
		map[string]any{"path": "/name"},
		map[string]any{
			"functionCall": map[string]any{
				"name": "toUpperCase",
				"args": []any{
					map[string]any{"path": "/title"},
				},
			},
		},
	}
	paths := PathsInArgs(args)
	if len(paths) != 2 {
		t.Fatalf("paths = %d, want 2", len(paths))
	}
	if paths[0] != "/name" || paths[1] != "/title" {
		t.Errorf("paths = %v, want [/name, /title]", paths)
	}
}

func TestEvalIf(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("if", []any{true, "yes", "no"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "yes" {
		t.Errorf("if(true) = %v, want 'yes'", result)
	}
	result, err = eval.Eval("if", []any{false, "yes", "no"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "no" {
		t.Errorf("if(false) = %v, want 'no'", result)
	}
}

func TestEvalOr(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("or", []any{false, false, true})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("or(false,false,true) = %v, want true", result)
	}
	result, err = eval.Eval("or", []any{false, false})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("or(false,false) = %v, want false", result)
	}
}

func TestEvalAnd(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("and", []any{true, true, true})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("and(true,true,true) = %v, want true", result)
	}
	result, err = eval.Eval("and", []any{true, false})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("and(true,false) = %v, want false", result)
	}
}

func TestEvalIfLazy(t *testing.T) {
	eval, _ := newTestEvaluator()
	// When condition is true, falseVal is NOT evaluated (contains error-producing calc("",0,0))
	brokenExpr := map[string]any{
		"functionCall": map[string]any{
			"name": "calc",
			"args": []any{"", float64(0), float64(0)},
		},
	}
	result, err := eval.Eval("if", []any{true, "yes", brokenExpr})
	if err != nil {
		t.Fatalf("if(true, 'yes', broken) should not error, got: %v", err)
	}
	if result != "yes" {
		t.Errorf("if(true, 'yes', broken) = %v, want 'yes'", result)
	}

	// When condition is false, trueVal is NOT evaluated
	result, err = eval.Eval("if", []any{false, brokenExpr, "no"})
	if err != nil {
		t.Fatalf("if(false, broken, 'no') should not error, got: %v", err)
	}
	if result != "no" {
		t.Errorf("if(false, broken, 'no') = %v, want 'no'", result)
	}
}

func TestEvalOrLazy(t *testing.T) {
	eval, _ := newTestEvaluator()
	// First arg is true → short-circuit, never evaluate broken second arg
	brokenExpr := map[string]any{
		"functionCall": map[string]any{
			"name": "calc",
			"args": []any{"", float64(0), float64(0)},
		},
	}
	result, err := eval.Eval("or", []any{true, brokenExpr})
	if err != nil {
		t.Fatalf("or(true, broken) should not error, got: %v", err)
	}
	if result != true {
		t.Errorf("or(true, broken) = %v, want true", result)
	}
}

func TestEvalAndLazy(t *testing.T) {
	eval, _ := newTestEvaluator()
	// First arg is false → short-circuit, never evaluate broken second arg
	brokenExpr := map[string]any{
		"functionCall": map[string]any{
			"name": "calc",
			"args": []any{"", float64(0), float64(0)},
		},
	}
	result, err := eval.Eval("and", []any{false, brokenExpr})
	if err != nil {
		t.Fatalf("and(false, broken) should not error, got: %v", err)
	}
	if result != false {
		t.Errorf("and(false, broken) = %v, want false", result)
	}
}

func TestEvalToNumber(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("toNumber", []any{"42"})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(42) {
		t.Errorf("toNumber('42') = %v, want 42", result)
	}
}

func TestEvalToString(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("toString", []any{float64(42)})
	if err != nil {
		t.Fatal(err)
	}
	if result != "42" {
		t.Errorf("toString(42) = %v, want '42'", result)
	}
}

func TestEvalCalc(t *testing.T) {
	eval, _ := newTestEvaluator()
	cases := []struct {
		op   string
		a, b float64
		want float64
	}{
		{"+", 2, 3, 5},
		{"-", 10, 3, 7},
		{"*", 4, 5, 20},
		{"/", 10, 2, 5},
	}
	for _, tc := range cases {
		result, err := eval.Eval("calc", []any{tc.op, tc.a, tc.b})
		if err != nil {
			t.Fatalf("calc(%s, %v, %v): %v", tc.op, tc.a, tc.b, err)
		}
		if result != tc.want {
			t.Errorf("calc(%s, %v, %v) = %v, want %v", tc.op, tc.a, tc.b, result, tc.want)
		}
	}
}

func TestEvalCalcDivZero(t *testing.T) {
	eval, _ := newTestEvaluator()
	_, err := eval.Eval("calc", []any{"/", float64(1), float64(0)})
	if err == nil {
		t.Error("expected division by zero error")
	}
}

func TestEvalContains(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("contains", []any{"hello world", "world"})
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Errorf("contains('hello world','world') = %v, want true", result)
	}
	result, err = eval.Eval("contains", []any{"hello", "xyz"})
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Errorf("contains('hello','xyz') = %v, want false", result)
	}
}

func TestEvalNegate(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("negate", []any{float64(42)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(-42) {
		t.Errorf("negate(42) = %v, want -42", result)
	}
	result, err = eval.Eval("negate", []any{float64(-5)})
	if err != nil {
		t.Fatal(err)
	}
	if result != float64(5) {
		t.Errorf("negate(-5) = %v, want 5", result)
	}
}

func TestEvalAddWithPath(t *testing.T) {
	eval, dm := newTestEvaluator()
	dm.Set("/x", float64(10))
	result, err := eval.Eval("add", []any{
		map[string]any{"path": "/x"},
		float64(5),
	})
	if err != nil {
		t.Fatal(err)
	}
	f, ok := result.(float64)
	if !ok || math.Abs(f-15) > 0.001 {
		t.Errorf("add = %v, want 15", result)
	}
}

func TestEvalAppendToArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b"}
	result, err := eval.Eval("append", []any{arr, "c"})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 3 {
		t.Fatalf("append result len = %d, want 3", len(r))
	}
	if r[2] != "c" {
		t.Errorf("append[2] = %v, want 'c'", r[2])
	}
}

func TestEvalAppendToEmptyArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{}
	result, err := eval.Eval("append", []any{arr, "x"})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 1 {
		t.Fatalf("append to empty = %d, want 1", len(r))
	}
	if r[0] != "x" {
		t.Errorf("append[0] = %v, want 'x'", r[0])
	}
}

func TestEvalAppendToNonArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("append", []any{"not-array", "elem"})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 1 {
		t.Fatalf("append to non-array = %d, want 1", len(r))
	}
	if r[0] != "elem" {
		t.Errorf("append[0] = %v, want 'elem'", r[0])
	}
}

func TestEvalRemoveLastFromArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b", "c"}
	result, err := eval.Eval("removeLast", []any{arr})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 2 {
		t.Fatalf("removeLast result len = %d, want 2", len(r))
	}
	if r[0] != "a" || r[1] != "b" {
		t.Errorf("removeLast = %v, want [a b]", r)
	}
}

func TestEvalRemoveLastFromSingleElement(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"only"}
	result, err := eval.Eval("removeLast", []any{arr})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 0 {
		t.Fatalf("removeLast single = %d, want 0", len(r))
	}
}

func TestEvalRemoveLastFromEmptyArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{}
	result, err := eval.Eval("removeLast", []any{arr})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 0 {
		t.Fatalf("removeLast empty = %d, want 0", len(r))
	}
}

func TestEvalSliceWithEnd(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b", "c", "d", "e"}
	result, err := eval.Eval("slice", []any{arr, float64(1), float64(4)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 3 {
		t.Fatalf("slice(1,4) len = %d, want 3", len(r))
	}
	if r[0] != "b" || r[1] != "c" || r[2] != "d" {
		t.Errorf("slice(1,4) = %v, want [b c d]", r)
	}
}

func TestEvalSliceNoEnd(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b", "c", "d"}
	result, err := eval.Eval("slice", []any{arr, float64(2)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 2 {
		t.Fatalf("slice(2) len = %d, want 2", len(r))
	}
	if r[0] != "c" || r[1] != "d" {
		t.Errorf("slice(2) = %v, want [c d]", r)
	}
}

func TestEvalSliceFromZero(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b", "c"}
	result, err := eval.Eval("slice", []any{arr, float64(0), float64(2)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 2 {
		t.Fatalf("slice(0,2) len = %d, want 2", len(r))
	}
	if r[0] != "a" || r[1] != "b" {
		t.Errorf("slice(0,2) = %v, want [a b]", r)
	}
}

func TestEvalSliceOutOfBounds(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{"a", "b"}
	result, err := eval.Eval("slice", []any{arr, float64(5)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 0 {
		t.Fatalf("slice(5) on 2-element array = %d, want 0", len(r))
	}
}

func TestEvalSliceEmptyArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	arr := []any{}
	result, err := eval.Eval("slice", []any{arr, float64(0)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 0 {
		t.Fatalf("slice(0) on empty = %d, want 0", len(r))
	}
}

func TestEvalSliceNonArray(t *testing.T) {
	eval, _ := newTestEvaluator()
	result, err := eval.Eval("slice", []any{"not-array", float64(0)})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := result.([]any)
	if !ok || len(r) != 0 {
		t.Fatalf("slice on non-array = %d, want 0", len(r))
	}
}
