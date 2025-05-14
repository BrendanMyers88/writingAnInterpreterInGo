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
* Only one method called `NextToken()` so no need to buffer or save tokens for this exercise. This method will output the value of the next token, and when calling the lexer the `NextToken()` method will be called repeatedly until the EOF token is read.
### 1.4 - Extending our Token Set and Lexer
### 1.5 - Start of a REPL
`writingAnInterpreterInGo/repl/repl.go`

## 2 - Parsing
### 2.1 - Parsers
* A parser takes input data (typically text) and builds a data structure. It gives a structural representation of the input, making sure the syntax is correct. It’s usually preceded by a separate lexical analyzer which creates the tokens from the input characters.
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
* An Abstract Syntax Tree is "abstract" based on certain details being omitted in the AST. Examples include Semicolons, newlines, whitespace, comments, braces, brackets, and parentheses. They aren’t represented in the AST, but guiding the parser when constructing it.

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
* CFGs are a set of rules that describe how to form correct (valid syntactically) sentences in a language. The most common notational formats of CFGs are Backus-Naur Form (BNF) or the Extended Backus-Naur Form (EBNF).

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
    * **_Recursive descent parsing_** - This is the parser we’ll be writing for Monkey, and in particular it’s a **"top down operator precedence"** parser, aka a "Pratt parser" after Vaughan Pratt.
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

* Every node in the AST has to implement the `Node` interface. It has to provide a `TokenLiteral()` method that returns the literal value of the token it's associated with.
* The `TokenLiteral()` will only be used for debugging/testing.
* Identifiers in other parts of a Monkey program **do** produce values, e.g.: `let x = valueProducingIdentifier;` so we have the Identifier struct type implement the Expression interface for simplicity.
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
* The second `*ast.InfixExpression` needs to have the integer literals 1 and 2 as its `Left` and `Right` child nodes, respectively
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
  * `parseExpression` sees `p.peekToken` is not a SEMICOLON and `peekPrecedence` is higher than the arg passed to `parseExpression` (SUM > LOWEST)
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