# Fig list expressions

This document describes the intermediate JSON format used for Fig.

## Simple values

The three basic types are `bool`, `number` and `string`.  These are represeented like so:


| Type    | JSON |
| ------- | ------------- |
| Bool    | `["bool:m:n", true}` or `["bool:m:n", false]`  |
| Number  | `["number:m:n", 1.5]`  |
| String  | `["string:m:n", "hello"]`|

In all the examples above, `m` and `n` are integer values  which refer
to the location in the source code where the corresponding tokens are
from.  This is optional.

## Names

Names are either global  or local identifiers.  Something like `x` in
source code would be mapped to `["name:m:n", "x"]`

## Fields

Field access is via the dot operator. `x.y` maps to `[".:m:n", ["name:m:n"], "x"], ["string:m:n", "y"]]`.

Note that the `x` part is mapped to use `name` while the actual field
accessed is mapped to `string`.  Fig allows the second  part to be
name as well: `x.(y)` would map to `[".:m:n", ["name:m:n"], "x"], ["name:m:n", "y"]]`.

In fact, fig allows `(some_expression).(some_other_expression)` and
these would get mapped appropriately.

## Standard binary operators

All the standard binary operators like `+`, `-`, `*`, `/`, `&`, `|`,
`==`, `!=`, `<`, `<=`, `>`, `>=` behave as one would expect.  These
always have exactly two args: `["op:m:n", left, right]`.

The unary operators `!` and the unary `-` only have one arg.

## Equals

The `=` operator is only used within functions and is a bit special in
that `x = y` will have the `x` part mapped to  `["string:m:n", "x"]`
instead  of `["name:m:n", "x"]` because this is defining `x`.

The equals operator can only appear in function args and closures and
is right associative.

## Function calls and closure

A function call `f(x, y)` gets mapped to `["call:m:n", ["name:m:n", "f"], ["name:m:n", "x"], ["name:m:n"], "y"]`

The list contains the function being called as the first entry with
all args following it later. All parts can have more involved
expressions.

Closures are similar: `{x}` maps to `["{}:m:n", ["name:m:n", "x"]]`.
Closures like `{x.y, x = z}` get treated much like the args list for
functions.

