# Language

The fig language is a simple expression based language

## Basic features:

| Features   | Examples |
| ---------- | ------------- |
| Number     | `1.5` or `1_000_000`  |
| Strings    | `"hello"` or `"Wayne's world"` |
| Boolean    | `true` or `false` |
| Arithmetic | `1 + 2` or `x * 15` or `x / 10` or `x - y` (Arithmetic only works on numbers) |
| Comparison | `x < y` or `x > y` or `x <= y` or `x >= y` (comparison works on any pair of types, not nust numbers) |
| Logical    | `x && y` or `x \|\| y` or `!x` (only works on booleans)|
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
* A where clause can show up in any function call or closures


## Streams

### Creating a stream

```
v = sys.streams.new(s = any_initial_value)
```

All the fields of the underlying value are available but calling `replace` on it will naturally edit the whole stream.

### Composing streams

```
v = sys.streams.join(x = some_stream, y = some_other_stream)
```

### Editing stream definitions

```
v1 = ...,
v2 = ...,
v = sys.streams.join(x= v1, y= v2)
z = v.replace(.x = 42)
```

The example here propagates the edit to the underlying `x` stream
(i.e. v1).  Future changes of v1 will still get reflected on `z`.

To change the stream definition itself, use `sys.streams.replace`:

```
sys.streams.replace(s= v, .x= 42)
```

### Snapshotting streams

```
sys.streams.snapshot(s = some_stream)
```

### Reactive stream expressions 

An expression like `x + y` is effectively a stream if either x or y is a stream. The result is basically another stream.

### Stateful stream functions

Stateful reactive streams can be built using `sys.streams.transform` which calls a handler on each delta, allowing it to mutate the stream in response (this mutation will not show up again in the handler):

```
delta_count = {
 transformed.result
 transformed = sys.streams.transform(s, xform)
 s = sys.streams.join(x = it.x, y = it.y, result = 0)
 xform = { it.result.replace(it.result + 1) }
}
```

## Macros

The `macro` function is a bit special:

* `where(x^5 = macro(x.get(5))` effectively replaces all occurence of `x^5` with `x.get(5)`
* `where(xml = macro({transform(it)})` effectively takes any occcurence of `xml(expr)` and calls `transform` on the AST of the expression allowing the macro to rewrite the AST.  This allows elegant ways of doing things like templating or JSX

