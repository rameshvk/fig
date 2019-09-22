# Language

The fig language is a simple expression based language

## Basic features:

| Features   | Examples |
| ---------- | ------------- |
| Number     | `1.5` or `1_000_000`  |
| Strings    | `"hello"` or `"Wayne's world"` or `"A \"quote\""`|
| Boolean    | `true` or `false` |
| Arithmetic | `1 + 2` or `x * 15` or `x / 10` or `x - y` (Arithmetic only works on numbers) |
| Comparison | `x < y` or `x > y` or `x <= y` or `x >= y` (comparison works on any pair of types, not nust numbers) |
| Logical    | `x & y` or `x \| y` or `!x` (only works on booleans)|
| Equality   | `x == y` or `x != y` (works on all types with comparison based on value, not reference)|
| Names      | `user` (names are either global context or scopes as defined later) |
| Fn call    | `f(x = 5)` or `g(y = 22)` (function args are named) |
| Closure    | `list.filter(by = {it.field > 22})` (curly braces define closures; `it.field` refers to named `field` arg) |
| Scope      | `list.filter(by = {it.field > z}, where(z = 22))` (where introduces a local scope in any function, allowing any names used before to be defined |
| Objects    | `object(x = 1, y = 2).x` (The `object` function takes arbitrary names) |
| Lists      | `list(1, 3)` (The `list` function is special) |

## Special characters in names

Names can have single quotes, back-quotes, square brackets and any non-whitespace unicode letter and all but the first unicode character can also be a unicode digit.

A single quoted name can include any character (just like a double-quote string).  Similarly with square brackets (which have to match, though they can be nested).  

Names can also have super-scripts and subscripts `x^5` is a superscript (which should be rendered in UI as x⁵).  Subscripts are done with underscores `x_5`.  Both super-scripts and sub-scripts are part of the name (so `x^1`  and `x^2` are not the same).   Names can also be primed: `x'` or `x''`. 

## Syntactic sugar

* Commas are optional at the end of a line.  Dangling operators are not allowed, i..e a line cannot end with `+`
* `x = y = z` is equivalent to `x = y, where (y= z)`
* There is no support for setting a field. Instead, most objects 
  are expected to support a `replace` method which returns a new value: `x.replace(5)`
* `x.replace(.y.z = 5)` is equivalent to `x.replace(x.y.replace(x.y.z.replace(5))`.  
  When x is a `stream`, both forms do the right thing (i.e. propagate changes)
* `List(x, y, z)` is shorthand for `List([0]= x, [1]= y, [2]= z)`
* Function calls with single args: `f(x)` is equivalent to `f(it = x)`
* Multi named arg shorthand #2: `f(x, y)` is equivalent to `f(x = x, y = y)` 
* A where clause can show up in any function call or in any closure


## Streams

### Creating a stream

```
v = sys.streams.new(s = any_initial_value)
```

All the fields of the underlying value are available but calling `replace` on it will naturally edit the whole stream.

If the input to `sys.streams.new` is itself a stream, the original stream is returned.

### Composing streams

All standard operations on streams (such as `x + y`) just return streams. The computation is treated as a reactive computation.

For example, `Object(x = stream1, y = stream2)` results in a stream whose objects have fields `x` and `y` that track the input stream.

### Editing stream definitions

Editing stream values cause back-propagation where it is meaningful.

For example, `z = Object(x = stream1, y = stream2), z' = z.replace(.x = 5)` effectively propagates the change upstream to `stream1` if that were possible.  If that isn't possible, the stream definition of `z` is changed so that its `x` field is replaced by a constant.

Explicit edits of a stream definition is possible using `sys.streams.replace` instead of `.replace`.

### Snapshotting streams

A single snapshot of a value (i.e. a non-stream fixed value) is obtained via `sys.streams.snapshot(s)`.

### Readonly streams

Readonly streams can be obtained by `sys.streams.readonly(s)` -- any changes are not propagated upstream in this (though it is not clear if edits will cause errors or simply cause the stream definition to change)
```

### Stateful stream functions

Stateful reactive streams can be built using `sys.streams.transform` which calls a handler on each delta, allowing it to mutate the stream in response (this mutation will not show up again in the handler).

The following example is a function that returns a stream which tracks number of deltas in its two input streams.

```
delta_count = {
 transformed.result
 transformed = sys.streams.transform(s, xform)
 s = Object(x = it.x, y = it.y, result = sys.streams.new(0))
 xform = { it.result.replace(it.result + 1) }
}
```

## Macros

The `macro` function is a bit special:

* `where(x^5 = macro(x.get(5))` effectively replaces all occurence of `x^5` with `x.get(5)`
* `where(xml = macro({transform(it)})` effectively takes any occcurence of `xml(expr)` and calls `transform` on the AST of the expression allowing the macro to rewrite the AST.  This allows elegant ways of doing things like templating or JSX

