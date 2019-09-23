# fig

[![Status](https://travis-ci.com/rameshvk/fig.svg?branch=master)](https://travis-ci.com/rameshvk/fig?branch=master)
[![GoDoc](https://godoc.org/github.com/rameshvk/fig?status.svg)](https://godoc.org/github.com/rameshvk/fig/pkg/fig)
[![codecov](https://codecov.io/gh/rameshvk/fig/branch/master/graph/badge.svg)](https://codecov.io/gh/rameshvk/fig)
[![Go Report Card](https://goreportcard.com/badge/github.com/rameshvk/fig)](https://goreportcard.com/report/github.com/rameshvk/fig)

Fig is a simple config server written in Go.

Fig takes the idea of **config as code** one step further by
considering every configuration or rollout flag a `function`.

Such functions are parameterized by things like `current user` or
`environment`.  For example, rolling out to 10 percent of users can
be thought of as a boolean function `strings.hash(user.email) < 0.1`

This moves away from traditional configuration (seen as JSON) into a
richer, more dynamic but still structured form.

Unlike a general purpose programming language, the configuration
language is meant to be rendered in UI well.  The same rollout
example above, if written as
`strings.hash(user.email) < ui.slider.percent(0.1)` would automatically
render a slider for the configuration.

## Architecture

1. The data is stored in a persistent store with a simple interface
allowing multiple backend implementations
2. The store is a effectively a key value store with the value being
the parsed functions.  Parsing results in a simple JSON
list-expression that can be easily evaluated on any client.
3. The server exposes a simple HTTP api to fetch the current value on
a key.  The server can also evaluate expressions if the parameters are
provided in JSON. Any app specific global extensions would only work
if the code for the server is modified to host these.
3. Client libraries cache the fetched values but interpret the
functions on each call.  Client libraries allow arbitrary extensions
(injected as global variables) in the native language.  A proxy server
is a server that uses a client for the backend (and so gets in-memory
caching).
4. Both the regular server and proxy server support authentication
modes (github initially).
5. Both the server and proxy server come with UI for administering the
key, value store. The UI for any particular config involves
interpreting the same config code but with a UI runtime (where
functions get rendered).  Custom renderers (as with the `ui.slider`
example) have access to the AST of the node and can freely mutate
them.  The current plan is to only support custom renderers within the
Fig language itself.  For example, the app can expose the
`ui.myslider` by simply using the `ui.myslider` key to store the
underlying UI code.

## Status

* [X] REDIS backed store
* [X] Basic API server
* [X] API key auth
* [ ] Basic edit UI
* [ ] UI auth
* [ ] API GH auth
* [ ] UI code parser
* [ ] UI code viewer
* [ ] UI history
* [ ] UI differ
* [ ] js client
* [ ] rb clienta
* [ ] py client
* [ ] proxy server

