# fig

Fig is a simple config server written in Go.

* Configuration in fig is a collection of key-value pairs. That is,
  this is flat and not hierarchical
* Fig has a full audit trail
* The actual values are expected to be *code* rather than viewed as
  data.  For example, the configuration for REDIS_URL could be
  something like: `if (env == "production", "prod.redis.io",
  "test.redis.io")`
* The variables and functions used in the configuration value can
  allow clear implementation of both feature flags and dynamic
  configuration.  Their actual definition is application-provided.
* The actual storage format is a simple list (much like LISP). This
  allow for new clients to implement the evaluation very easily
* The UI for fig encodes raw code into this format.  It also provides
  vieewers and differs.

## Status

* [X] REDIS backed store
* [ ] Basic API server
* [ ] API key auth
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

