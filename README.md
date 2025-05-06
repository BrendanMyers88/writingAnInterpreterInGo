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
```go
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
