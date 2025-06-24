# Writing an Interpreter in Go

## Interpreter features:

* Tokenize and parse Monkey source code in a REPL

1. The Lexer
2. The Parser
3. The Abstract Syntax Tree (AST)
4. The Internal Object System
5. The Evaluator

## 1 - Lexing

### 1.1 - Lexical Analysis

* Source Code → Tokens → Abstract Syntax Tree
* Lexical Analysis or "Lexing" is the transformation from Source Code to Tokens.
* This is done by a Lexer, aka Tokenizer or Scanner
* Tokens are small, easily categorizable data structures which can be fed to the parser
* Example input to lexer:
  `let x = 5 + 5;`
* Example output:

```
[
LET,
IDENTIFIER("x"),
EQUAL_SIGN,
INTEGER(5),
PLUS_SIGN,
INTEGER(5),
SEMICOLON
]
```

### 1.2 - Defining Our Tokens

`~/projects/writingAnInterpreterInGo/token/token.go`

### 1.3 - The Lexer

* Only one method called `NextToken()` so no need to buffer or save tokens for this exercise. This method will output
  the value of the next token, and when calling the lexer the `NextToken()` method will be called repeatedly until the
  EOF token is read.

### 1.4 - Extending our Token Set and Lexer

### 1.5 - Start of a REPL

`writingAnInterpreterInGo/repl/repl.go`

## 2 - Parsing

### 2.1 - Parsers

* A parser takes input data (typically text) and builds a data structure. It gives a structural representation of the
  input, making sure the syntax is correct. It’s usually preceded by a separate lexical analyzer which creates the
  tokens from the input characters.
* Types of data structures:
  * Parse tree
  * Abstract syntax tree
  * Other hierarchical structures
* Ex:

```
> var input = '{"name": "Testname", "age": 100}';
> var output = JSON.parse(input);
> output
{ name: "Testname", age: 100 }
> output.name
"Testname"
> output.age
100
```

* An Abstract Syntax Tree is "abstract" based on certain details being omitted in the AST. Examples include Semicolons,
  newlines, whitespace, comments, braces, brackets, and parentheses. They aren’t represented in the AST, but guiding the
  parser when constructing it.

#### Parser Input:

```javascript
if (3 * 5 > 10) {
  return "hello";
} else {
  return "goodbye";
}
```

#### Parser Output:

```
> var input = 'if (3 * 5 > 10) { return "hello"; } else { return "goodbye"; }';
> var tokens = MagicLexer.parse(input);
> MagicParser.parse(tokens);
{
  type: "if-statement",
  condition: {
    type: "operator-expression",
    operator: ">",
    left: {
      type: "operator-expression",
      operator: "*",
      left: {
        type: "integer-literal",
        value: 3
      },
      right: {
        type: "integer-literal,
        value: 5
      },
    },
    right: {
      type: "integer-literal",
      value: 10
    }
  },
  consequence: {
    type: "return-statement",
    returnValue: {
      type: "string-literal",
      value: "hello"
    }
  },
  alternative: {
    type: "return-statement",
    returnValue: {
      type: "string-literal",
      value: "goodbye"
    }
  }
}
```

* Parsers analyze the input, checking that it conforms to the expected data structure.
* Parsing is also known as syntactic analysis.

### 2.2 - Why not a parser generator?

* Context-Free Grammar (CFG) is the input for the majority of parser generators (ex. Yacc, Bison, or ANTLR.
* CFGs are a set of rules that describe how to form correct (valid syntactically) sentences in a language. The most
  common notational formats of CFGs are Backus-Naur Form (BNF) or the Extended Backus-Naur Form (EBNF).

EmcaScript syntax in BNF:

```
PrimaryExpression ::= "this"
  | ObjectLiteral
  | ( "(" Expression ")" )
  | Identifier
  | ArrayLiteral
  | Literal
Literal ::= ( <DECIMAL_LITERAL>
  | <HEX_INTEGER_LITERAL>
  | <STRING_LITERAL>
  | <BOOLEAN_LITERAL>
  | <NULL_LITERAL>
  | <REGULAR_EXPRESSION_LITERAL> )
Identifier ::= <IDENTIFIER_NAME>
ArrayLiteral ::= "[" ( ( Elision )? "]"
  | ElementList Elision "]"
  | ( ElementList )? "]" )
ElementList ::= ( Elision )? AssignmentExpression
  ( Elision AssignmentExpression )*
Elision ::= ( "," )+
ObjectLiteral ::= "{" ( PropertyNameAndValueList  )? "}"
PropertyNameAndValueList ::= PropertyNameAndValue ( "," PropertyNameAndValue
  | "," )*
PropertyNameAndValue ::= PropertyName ":" AssignmentExpression
PropertyName ::= Identifier
  | <STRING_LITERAL>
  | <DECIMEL_LITERAL>
```

* A parser generator would take something like the above and turn it into compilable C code.

### 2.3 - Writing a Parser for the Monkey Programming Language

Two main strategies when parsing a programming language:

* Top-down parsing
  * **_Recursive descent parsing_** - This is the parser we’ll be writing for Monkey, and in particular it’s a **"top
    down operator precedence"** parser, aka a "Pratt parser" after Vaughan Pratt.
  * Early parsing
  * Predictive parsing
  * Begins with constructing root node of the AST and then descends
  * Recommended for newcomers to parsers
* Bottom-up parsing

1. Parsing statements (let and return)
2. Parsing expressions
3. Extend parser to be capable of parsing a large subset of Monkey programming language.
4. As we go along, build up necessary structures for our AST

### 2.4 - Parser's first steps: parsing let statements

* Variable bindings are statements of the following form (let statements):

```javascript
let x = 5;
let y = 10;
let foobar = add(5, 5);
let boofar = 5 * 5 / 10 + 18 - add(5, 5) + multiply(124);
let anotherName = barfoo;
```

#### Valid Monkey program:

```
let x = 10;
let y = 15;
let add = fn(a, b) {
  return a + b;
};
```

#### Variable binding base-parts:

```
let <identifier> = <expression>;
```

* Expressions produce values:
  * `5`
  * `add(5, 5)`
* Statements do not:
  * `let x = 5`
  * `return 5`

* Every node in the AST has to implement the `Node` interface. It has to provide a `TokenLiteral()` method that returns
  the literal value of the token it's associated with.
* The `TokenLiteral()` will only be used for debugging/testing.
* Identifiers in other parts of a Monkey program **do** produce values, e.g.: `let x = valueProducingIdentifier;` so we
  have the Identifier struct type implement the Expression interface for simplicity.

#### AST representation of `let x = 5;` in Monkey programming language

```
            [ *ast.Program ]
            [  Statements  ]
                    |
                    v
          [ *ast.LetStatement ]
          [       Name        ]
          [       Value       ]
            /                \
          v                    v
[ *ast.Identifier ]    [ *ast.Expression ]
```

#### parseProgram pseudocode:

```
function parseProgram() {
  program = newProgramASTNode()
  advanceTokens()
  
  for (currentToken() != EOF_TOKEN) {
    statement = null
    if (currentToken() == LET_TOKEN {
      statement = parseLetStatement()
    } else if (currentToken() == RETURN_TOKEN) {
      statement = parseReturnStatement()
    } else if (currentToken() == IF_TOKEN) {
      statement = parseIfStatement()
    }

    if (statement != null) {
      program.Statements.push(statement)
    }

    advanceTokens()
  }
  
  return program    
}

function parseLetStatement() {
  advanceTokens()
  identifier = parseIdentifier()
  advanceTokens()
  if currentToken() != EQUAL_TOKEN {
    parseError("no equal sign!")
    
    return null
  }
  advanceTokens()
  
  value = parseExpression()
  
  variableStatement = newVariableStatementASTNode()
  variableStatement.identifier = identifier
  variableStatement.value = value
  
  return variableStatement
}

function parseIdentifier() {
  identifier = newIdentifierASTNode()
  identifier.token = currentToken()
  
  return identifier
}

function parseExpression() {
  if (currentToken() == INTEGER_TOKEN) {
    if (nextToken() == PLUS_TOKEN) {
      return parseOperatorExpression()
    } else if (nextToken() == SEMICOLON_TOKEN) {
      return parseIntegerLiteral()
    }
  } else if (currentToken() == LEFT_PAREN) {
    return parseGroupedExpression()
  }
}

function parseOperatorExpression() {
  operatorExpression = newOperatorExpression()
  operatorExpression.left = parseIntegerLiteral()
  advanceTokens()
  operatorExpression.operator = currentToken()
  advanceTokens()
  operatorExpression.right = parseExpression()
  
  return operatorExpression
}
```

### 2.5 - Parsing Return Statements

#### Example return statements in Monkey:

```javascript
return 5;
return 10;
return add(15);
```

#### Return Statement structure:

```
return <expression>;
```

### 2.6 - Parsing Expressions

* In Monkey, everything except `let` and `return` statements are expressions

#### Examples w/ Prefix Operators

```
-5
!true
!false
```

#### Examples w/ Infix Operators

```
5 + 5
5 - 5
5 / 5
5 * 5
```

#### Examples of Comparison Operators

```
foo == bar
foo != bar
foo < bar
foo > bar
```

#### Examples Grouped Expressions and Order of Evaluation Influence:

```
5 * (5 + 5)
((5 + 5) * 5) * 5
```

#### Examples of call expressions:

```
add(2, 3)
add(add(2, 3), add(5, 10))
max(5, add(5, (5 * 5)))
```

#### Examples of Identifiers as expressions:

```
foo * bar / foobar
add(foo, bar)
```

#### Function literals are expressions too:

```
let add = fn(x, y) { return x + y }
```

#### Examples of Function literal in place of identifier

```
fn(x, y) { return x + y }(5, 5)
(fn(x) { return x }(5) + 10 ) * 10
```

#### Examples of if expressions:

```
let result = if (10 > 5) { true } else { false };
result // => true
```

###### Terminology

* Prefix operator - An operator in front of its operand, Ex:
  * `--5`
    * Operator: `--` (Decrement)
    * Operand: `5`
* Postfix operator - An operator after its operand, Ex:
  * `foobar++`
    * Operator: `++` (Increment)
    * Operand: `foobar`
  * Will _not_ be included in this book.
* Infix operator - An operator which sits between two operands
  * `5 * 8`
    * Operator: `*`
    * Operands: `5`, `8`
* Binary expressions - Operator has two operands
* Operator precedence/Order of operations - Priority which different operators have.
  * `5 + 5 * 10`

### 2.7 - How Pratt Parsing Works

* Suppose we're parsing the following expression statement:
  * `1 + 2 + 3;`
* The goal isn't to represent all operators and operands in the resulting AST, but to nest the nodes correctly.
* We want the AST (serialized as a string) to look like this:
  * `((1 + 2 ) + 3)`
* The AST needs 2 `*ast.InfixExpression` node
* The `*ast.InfixExpression` higher in the tree should have the integer literal 3 as it's `Right` child node
* the `Left` child node needs to be the second `*ast.InfixExpression`.
* The second `*ast.InfixExpression` needs to have the integer literals 1 and 2 as its `Left` and `Right` child nodes,
  respectively
* See 2.7, page 170/522 in Kindle for AST map representation or:
* `*ast.InfixExpression`
  * (`Left Node`) `*ast.InfixExpression`
    * (`Left Node`) `*ast.IntegerLiteral`
      * `1`
    * (`Right Node`) `*ast.IntegerLiteral`
      * `2`
  * (`Right Node`) `*ast.IntegerLiteral`
    * `3`

* `1 + 2 + 3 ;`
  * `1 = p.curToken`
  * `+ = p.peekToken`
  * `parseExpression` checks for `prefixParseFn` associated with current `p.curToken`
    * `token.INT` -> parseIntegerLiteral -> `*ast.IntegerLiteral` assigned to `leftExp`
  * `parseExpression` sees `p.peekToken` is not a SEMICOLON and `peekPrecedence` is higher than the arg passed to
    `parseExpression` (SUM > LOWEST)
    * fetch `infixParseFn` for `p.peekToken` -> `parseInfixExpression` -> returns to `leftExp` to advance the token
  * `+ = p.curToken`
  * `2 = p.peekToken`
  * loops here...

### 2.8 - Extending the Parser

* Capital first letter = Public; lower first letter = Private for method names
* Adding function literals in this chapter

### 2.9 - Read-Parse-Print-Loop

* Until now the REPL was more of a RLPL (Read, Lex, Parse, Loop)
* Next we're going to change it to an RPPL (Read, Parse, Print, Loop)

## 3 - Evaluation

### 3.1 - Giving Meaning to Symbols

* We'll be implementing the E (Evaluation) in our REPL in this chapter.
* Without the evaluator, we don't get the output of `1 + 2 = 3`, `4 < 3 = false`, etc.
* The evaluation process of an interpreter defines how the programming language being interpreted works.

```javascript
let num = 5;
if (num) {
  return a;
} else {
  return b;
}
```

* In some languages, this would return `a`, in other languages `b`. This could be a boolean return as well.
* This will also determine the order in which to return outputs and evaluation.

### 3.2 - Strategies of Evaluation

* Evaluation is where interpreter implementations diverge the most.
* Tree-walk interpreters evaluate as they traverse the AST on the fly.
  * ie, printing a string, adding two numbers, or executing a function's body.
* Sometimes the evaluation step is preceded by small optimizations that re-write the AST or convert it into another
  intermediate representation.
* Others traverse the AST but instead of interpreting the AST itself, then convert to byte-code first. These can't be
  interpreted by the system like assembly, but instead are interpreted by a virtual machine.
  * This can be beneficial for performance.
* JIT ("just in time") interpreters compile byte-code to native machine code right before it's executed.
* Other JIT interpreters skip byte-code and go straight from the AST to native machine code.
* Ruby (through v1.8) was a tree-walking interpreter, evaluating the AST while traversing it.
* Ruby (v1.9+) switched to a virtual machine architecture.
  * Interpreter parses source code
  * Builds an AST
  * Compiles AST into bytecode
  * Bytecode gets executed in a virtual machine, improving performance.
* Lua switching to LuaJIT improved benchmarks by up to 50x performance.

### 3.3 - A Tree-Walking Interpreter

* We'll be building a tree-walking interpreter.
* We'll take the AST the parser builds and interpret in "on the fly" without any preprocessing or compilation step.
  * Similar to Lisp interpreter.
  * Inspired by the interpreter in "The Structure and Interpretation of Computer Programs" (SICP)
  * Easiest way to get started, understand, and extend on later.
* We need two things:
  * Tree-waking Evaluator
  * A way to represent Monkey values in our host language of Go.
* Pseudocode:

```
function eval(astNode) {
  if (astNode is integerLiteral) {
    return astNode.integerValue
  } else if (astNode is booleanLiteral) {
    return astNode.booleanValue
  } else if (astNode is infixExpression) {
    leftEvaluated = eval(astNode.Left)
    rightEvaluated = eval(astNode.Right)
    
    if astNode.Operator == "+" {
      return leftEvaluated + rightEvaluated
    } else if ast.Operator == "-" {
      return leftEvaluated - rightEvaluated
    }
  }   
```

* `eval` is recursive.
* When `astNode is infixExpression` is true, `eval` calls itself again twice to evaluated the left and right operands of
  the infix expression.
* This can cascade to another infix expression or integer literal or boolean literal, or identifier, or etc.
* Same concept as the recursion in the AST, except we're evaluating the tree, not building it.

### 3.4 - Representing Objects

* Not object-oriented, but we need an "object representation" or a "value system" that defines what our `eval` function
  returns.
* Different ways of representing objects:
  * Native types (integers, booleans, etc.) of the host language
  * Values/objects only represented as pointers
  * Native types and pointers are mixed
  * etc.
* Why the variety?
  * Host languages differ, i.e., an interpreter written in Ruby can't represent values the same way as an interpreter
    written in C.
  * Languages being interpreted differ, i.e., some only need representations of primitive data like integers,
    characters, or bytes whereas others you'll have lists, dictionaries, functions, or compound data-types.
  * Resulting execution speed and memory consumption while evaluating programs differ depends on the choice.
  * For high execution speed, you can't use a slow and bloated object system
  * If writing a garbage collector, you need to think about how it'll keep track of values in the system.
  * If performance doesn't matter, then keeping things simple and understandable can be advantageous.

#### Foundation of our Object System

* We'll represent every value we encounter in Monkey as an `Object`.
* All values will be wrapped inside a struct, which fulfills the `Object` interface.
* In this chapter, we represented Integer, Boolean, and Null ObjectTypes in our `object.go` file.

### 3.5 - Evaluating Expressions

* In this chapter, we'll be writing our `eval`
* The first version will looks like:

```
func Eval(node ast.Node) object.Object
```

* Eval takes an ast.Node and returns object.Object

#### Integer Literals

* Given an `*ast.IntegerLiteral` our `Eval` function should return an `*object.Integer` whose `Value` field contains the
  same integer as `*ast.IntegerLiteral.Value`

#### Completing the REPL
* We'll now be updating the RPPL to an REPL in `repl.go`

#### Prefix Expressions
#### Infix Expressions
* The eight infix operators in Monkey are: `+, -, *, /, >, <, ==, !=`

### 3.6 - Conditionals
* Hardest part of conditional evaluation is deciding when to evaluate what.

### 3.7 - Return Statements
* Return statements stop the evaluation for a series of statements and leave behind the value their expression has evaluated to.

### 3.8 - Abort! Abort! There's been a mistake!, or: Error Handling
* Instead of returning nil values, we'll instead change those to handle errors in this chapter
* This _won't_ be user-defined exceptions, but internal error handling for:
  * Wrong operators
  * Unsupported operations
  * Other user or internal errors that may arise during execution.
* This will be handled very similarly to the return statements of 3.7.
  * This is because Errors and return statements both stop the eval for a series of statements.

### 3.9 - Bindings & The Environment
* Add bindings to the interpreter by adding support for `let` statements
* We also need to support evaluation of identifiers

### 3.10 - Functions & Function Calls
* We'll add support for functions and function calls to the interpreter in this chapter

### 3.11 - Who's taking the trash out?
* In Monkey, we're re-using Go's garbage-collector rather than handmaking one ourselves up to this point.
* Go's garbage collector keeps track of which `object.Integer` are still reachable to us and which are not, saving a great deal of system memory.
* The GC would need to:
  * Keep track of object allocations
  * Keep track of references to objects
  * Make enough memory available for future object allocations
  * Give memory back when it's not necessary anymore (without this, we'd have a memory leak)
* Types of GC algorithms:
  * Mark and Sweep
  * Generational GC
  * Stop-the-world GC
  * Concurrent GC
  * Need to know how it's organizing memory and handing memory fragmentation
* Unfortunately, we can't take over Go's Garbage Collector easily as by default, Go prohibits exactly that

## 4 - Extending the Interpreter

### 4.1 - Data Types & Functions
* In 4.1 we'll add new data types to the interpreter, rather than just having boolean/integers.
  * This will include:
    * Adding new token types
    * Modifying the lexer
    * Extending the parser
    * Adding support for the data types to the evaluator and object system
  * These data types are already present in Go, therefor we'll just need to make them available in Monkey rather than making them from scratch.
  * We'll add some new functions to make the interpreter more powerful
    * Built-in functions
  * First thing we'll do is add the string data type.

### 4.2 - Strings
* First, we'll add support for string literals to the lexer. The basic structure is:
  * `"<sequence of characters>"`, i.e., a sequence of characters enclosed by double quotes.

### 4.3 - Built-in Functions
* In this section we'll add built-in functions, specifically inheriting from Go
* For example: A function that returns the current time
  * Normally does via asking the kernel, which is handled by system calls.
  * If a language doesn't offer the use of system calls, then the language implementation has to provide something to make these calls on behalf of the users.
  * The built-ins will be defined by us, the implementers of the interpreter.
  * The end-user can use them, but they are defined by us.
* The only restriction is they need to accept zero or more `object.Object` as arguments and return an `object.Object.`
* Builtin function's we'll add:
  * `len`: Return number of characters in a string.
  * `>> len("Hello World!")`
  * `12`

### 4.4 - Array
* Read through...
