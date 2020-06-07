# Key-path Notation

## Introduction

This document is the primary reference for Key-path notation as it relates to selectr. This document is intended to informally describe each language construct and their use with selectr.

The basic lexical structure is:

```
KEY_PATH = (IDENTIFIER)? (INDEX_EXPRESSION | ATTRIBUTE_EXPRESSION)*
```

Key-path notation is interpreted as a chain of attribute or index expressions as a means of traversing data. As such, there are only two types of literals: String literals and Integer literals. These value types are meant to be used in index expressions.

If the first expression is an attribute expression, the dot can be ommited (leaving just the identifier).

#### Examples:

```
attribute.nestedAttribute

.attribute

['my string subscript']

[0]

['foo'][0]["bar"].object.array[1].value

object.array[1].object.nestedObject.array[0].someValue
```

## Input format

Key-path notation is formatted as a sequence of Unicode points encoded in UTF-8.

## Identifiers

```
[a-zA-Z][a-zA-Z0-9_]*
```

An identifier is any nonempty ASCII string where:

- The first character is a letter.
- The remaining characters are alphanumeric or `_`.

Identifiers are to be used in dot notation within attribute expressions.

## Whitespace

Whitespace is any non-empty string containing the following characters:

- `U+0009` (horizontal tab, `'\t'`)
- `U+000A` (line feed, `'\n'`)
- `U+000D` (carriage return, `'\r'`)
- `U+0020` (space, `'\s'`)

All forms of whitespace serve only to separate tokens in the grammar and have no semantic significance.

## Literal Expressions

```
LITERAL_EXPRESSION = INTEGER_LITERAL | STRING_LITERAL
```

### Integer literals

```
INTEGER_LITERAL = [0-9]*
```

Integer literals are a sequence of digits.

#### Examples:

```
1

1000

123456789

0
```

### String literals

```
QUOTE = ' | "

QUOTE_ESCAPE = \' | \"

ASCII_ESCAPE = (
    `\a` | `\b` | `\e` | `\f` | `\n` | `\r` | `\t` | `\v` | `\\` | `\?`
)

CHARACTER = [^"']

STRING_LITERAL = QUOTE (
    QUOTE_ESCAPE |
    ASCII_ESCAPE |
    CHARACTER
)* QUOTE
```

String literals are a sequence of characters contained between a pair of quotes. Both single and double quotes can be used in a string literal expression.

#### Examples: 

```
"A string"

'Another string'

"abc123$#%"

"\"\n\t\s"
```

## Attribute Expressions

```
DOT = .

ATTRIBUTE_EXPRESSION = DOT IDENTIFER
```

Attribute expressions denote a reference to an attribute of the subject.

## Index Expressions

```
LBRACKET = [

RBRACKET = ]

INDEX_EXPRESSION = LBRACKET LITERAL_EXPRESSION RBRACKET
```

Index expressions denote a reference to an attribute of an object or an element of an array.
