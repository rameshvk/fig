# Language

The fig language is a simple expression based language

## Basic features:

| Features   | Examples |
| ---------- | ------------- |
| Number     | `1.5` or `1_000_000`  |
| Strings    | `"hello"` or `"Wayne's world"` |
| Boolean    | `true` or `false` |
| Arithmetic | `1 + 2` or `x * 15` or `x / 10` or `x - y` (Arithmetic only works on numbers) |
| Comparison | `x < y` or `x > y` or `x <= y` or `x >= y` (comparison works on all types, not nust numbers) |
| Logical    | `x && y` or `x || y` or `!x` (only works on booleans)|
| Equality   | `x == y` or `x != y` (works on all types)|
| Names      | `user` (names are either global context or scopes as defined later) |
| Fn call    | `f(x = 5)` or `g(y = 22)` (function args are always named) |
| Closure    | `list.filter(by = {it.field > 22})` (curly braces define closures; `it.field` refers to named `field` arg) |
| Scope      | `list.filter(by = {it.field > z}, where(z = 22))` (where introduces a local scope in any function, allowing any names used before to be defined |
| Fn #2      | `f(x)` can be used instead of `f(x = x)`. |
| Objects    | `Object(x = 1, y = 2).x` (The `Object` function takes arbitrary names) |

## Syntactic sugar

* Commas are optional at the end of a line.
* `x = y = z` is equivalent to `x = y, where (y = z)`
* `x.0` is equivalent to `x.get(idx=0)` (There is no native array support)
* There is no support for setting a field. Instead, most objects 
  are expected to support a `replace` method which returns a new value.
* `x.replace(.y.z = 5)` is equivalent to `x.replace(x.y.replace(x.y.z.replace(value=5))`.  
When x is a `stream`, both forms do the right thing (i.e. propagate changes)

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
v = sys.streams.join(x = v1, y = v2)
z = v.replace(.x = 42)
```

The example here propagates the edit to the underlying `x` stream
(i.e. v1).  Future changes of v1 will still get reflected on `z`.

To change the stream definition itself, use `sys.streams.replace`:

```
sys.streams.replace(s = v, .x = 42)
```

### Snapshotting streams

```
sys.streams.snapshot(s = some_stream)
```

### Stateful stream functions

The `concat(x, y)` function which concatenates elements from two array streams:

```
concat = {
  stateful.result
  stateful = sys.streams.stateful(
    x = it.x
    y = it.y
    state = Object(result=concat(xval, yval), x=it.x)
    where(xval = sys.streams.snapshot(x), yval = sys.streams.snapshot(y))
    next = {
      new_state,
      new_state = it.state.replace(.x = resultx, result=resulty)
      resultx = xdelta.applyTo(it.result)
      resulty = ydelta.applyTo(resultx)
      xdelta = it.xdelta.split(path="x").affected
      ydelta = it.xdelta.split(path="x").unaffected.shift(offset=resultx.count())
    }
  )
}
```

`stateful` takes `state` as the `initial state` and `next` as the function which applies
any `delta` on the provided stream args (in this case `x` and `y`) to the state.  Any 
updates of the state are propagated.

## Macros

The `macro(xform, code)` function is treated like aa preprocessor directive.
The AST of the `code` arg (which can be any expression) is transformed by the `closure`. 

The example below indicates how a JSX-like extension might work with the preprocessor support.

```
module(
  name = "my module"
  export = fig.macro(
    code = '<div> my xml template {hello} </div>'
    xform = import("github.com/rameshvk/fig/modules/xml").parser
    where(hello = "hello, world")
  )
)
```

## Code matches

Macros allow arbitrary transformations but for more Fig=>Fig transformatinos, a more direct system would be useful:

```
module(
  .....
  fig.replace(
     match = [AnyName](z=[AnyValue])
     replace = [AnyName](x=[AnyValue])
  ) 
```
