package server

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

// NewRedisStore creates a new redis-based store
//
// The address is the redis server address in host:port format.
// The provided namespace is used to construct a prefix of the
// form "{fig_<ns>}" for all the keys
func NewRedisStore(address, ns string) Store {
	prefix := "{fig_" + ns + "}"
	return red{redis.NewClient(&redis.Options{Addr: address}), prefix}
}

type red struct {
	*redis.Client
	prefix string
}

func (r red) GetSince(version int) (int, map[string]string) {
	keys := []string{r.prefix}
	n, err := luaGetSince.Run(r.Client, keys, version).Result()
	if err != nil {
		panic(err)
	}
	pair := n.([]interface{})
	var config map[string]string
	if err := json.Unmarshal([]byte(pair[1].(string)), &config); err != nil {
		panic(err)
	}

	return int(pair[0].(int64)), config
}

func (r red) Set(key string, val string) {
	keys := []string{r.prefix}
	_, err := luaSet.Run(r.Client, keys, key, val).Result()
	if err != nil {
		panic(err)
	}
}

func (r red) History(key, epoch string) (string, []string) {
	keys := []string{r.prefix}
	n, err := luaHistory.Run(r.Client, keys, key, epoch).Result()
	if err != nil {
		panic(err)
	}
	pair := n.([]interface{})
	items := pair[1].([]interface{})
	result := make([]string, len(items))
	for kk := range items {
		result[kk] = items[kk].(string)
	}
	return pair[0].(string), result
}

var luaCommon = `
  local prefix = KEYS[1]
  local keyVersions = prefix.."_versions"
  local keyEntry = function(key) return prefix.."_key"..tostring(key) end
`

var luaGetSince = redis.NewScript(luaCommon + `
  local min = 1+tonumber(ARGV[1])
  local items = redis.call("ZREVRANGEBYSCORE", keyVersions, "+inf", min, "WITHSCORES")

  if not(items) or #items == 0 then
    return {min-1, "{}"}
  end

  local result = {}
  local key = ""
  for idx, val in pairs(items) do
    if idx % 2 == 1 then
      key = val
    else
      result[key] = redis.call("ZREVRANGE", keyEntry(key), 0, 0)[1]
    end
  end

  return {tonumber(items[2]), cjson.encode(result)}
`)

var luaSet = redis.NewScript(luaCommon + `
  local key, val = ARGV[1], ARGV[2]
  local ver = 1
  local items = redis.call("ZREVRANGE", keyVersions, 0, 0, "WITHSCORES")

  if items and #items > 0 then
    ver = items[2]+1
  end

  redis.call("ZADD", keyEntry(key), ver, val)
  redis.call("ZADD", keyVersions, ver, key)
  return 0
`)

var luaHistory = redis.NewScript(luaCommon + `
  local key, max = ARGV[1], ARGV[2]
  if max == "" then
    max = "+inf"
  end


  local items = redis.call("ZREVRANGEBYSCORE", keyEntry(key), max, 0, "WITHSCORES")
  local result = {}

  key = ""
  for idx, val in pairs(items) do
    if idx % 2 == 0 then
      key = val - 1
    else
      result[1+#result] = val
    end
  end

  return {tostring(key), result}
`)
