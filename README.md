```
  |   _______ _____  ______ _______ _______    |
  |   |         |   |_____/ |       |_____|    |
  |   |_____  __|__ |    \_ |_____  |     |    |
  |                                            |
```

Fast, smart and simple sidecar caching proxy
 - Simple configuration with json config or managment api
 - Cache stampede problem solution
 - Ratelimits
 - RequestID append and check
 - Indepotency checker
 - Failover cache
 - Circut Breaker protection
 - Prometeus metrics

Why
===

In microservice world you may need cache for service communication. Usually developers implement this cache in microservice itself but it is very common issue. What if we can just user cache proxy in front of microservice for that? The same with Circuit Breaker, RequestID generation and logging, rate limits and so on. You can try to configure nginx as side car for that, but it is not easy to configure and you may need to know hot to write in lua script. `Circa` need only a simple json config with rules for you enpdoints that helps you to do all that staff from the box in 5 minutes. Also it can modify requests, glue responces, 

Configuration
=============

Sipliest way to configure `circa` - prepare json configuration file 
First we need to configure `storages` - backends that used for storing cache. You may use a few storages for different goals.

```json
{
  "storages": {
    "main_ha": "redis://redis.default.svc.cluster.local/5?pool_size=30&timeout=5ms",
    "secondary": "redis://redis_secondary.default.svc.cluster.local/10",  // pool_size=10 timeout=30ms by default
    "fast": "mem://?size=5000"
  },
  // ...
}  
```
Specify target host and timeout for forwarding requests
```json
  // storages section
  "options": {
    "host": "http://localhost:8888",
    "timeout": "1m30s"
  },
  // ...
```

 Than you must define routes, three types:
 1) exact - `/my/path`
 2) with parameter - `/my/path/{parameter}`
 3) prefix or proxying - `/my/*`
 4) mixed `/my/{parameter}/path/*`

 And need to define rules for our route: each rule heve `kind` and options related to this kind of rule, 

 ```json
 {
  // ...
  "rules": {
    "/posts/{id}": [  // route
      {
        "kind": "cache",  // rule kind "cache" - simple cache strategy
        "ttl": "10s",  // cache ttl
        "key": "post-{id}",  // cache key template, can contain url parameter from route
        "storage": "main_ha"  // storage name from "storages" section that will be used for this rule
      }
      // {... other rule}
    ]
    // other route
  }
 }
 ```
 That is, you can run `circa` and it will cache requests to GET /posts/{id} for 10 seconds and proxy all other requests 

 There are a list of rules kinds and parameters for each:

1. Cache with ttl
```json
{
    "kind": "cache",  // by default
    "ttl": "10m50s", // cache time to live  Valid time units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”.
    "key": "....."
}
```

2. Cache with early expiration - 
```json
{
    "kind": "early",
    "ttl": "10h", // cache time to live
    "early_ttl": "1h",  // time for pre invalidation
    "key": "....."
}
```

3. Cache with hit expiration - 
```json
{
    "kind": "hit",
    "ttl": "10h", // cache time to live
    "hits": 100,  // number of hits before cache will be invalidated
    "update_after": 10, // optional; number of hits for pre invalidation
    "key": "....."
}
```

4. Failover cache - 
```json
{
    "kind": "fail",
    "ttl": "10h", // cache time to live
    "key": "....."
}
```

5. Rate limit - 
```json
{
    "kind": "rate-limit",
    "ttl": "10m", // limit period
    "hits": 100, //  number of hits  ( read as 100 request per 10 min)
    "key": "....."
}
```

6. Retry - 
```json
{
    "kind": "retry",
    "methods": ["GET", "HEAD"],
    "backoff": "5s",
    "count": 5 // retry attempts
}
```


7. request_id -  `check_return` - check that backend return response with tha same request ID
```json
{
    "kind": "request_id",
    "methods": ["GET", "POST", "DELETE", "PUT", "PATCH"],
    "check_return": false // by default
}
```

8. idempotency
```json
{
    "kind": "idempotency",
    "ttl": "1m"
    // "key": "{R:body|hash}",
}
```

9. invalidate
```json
{
    "kind": "invalidate",
    "methods": ["POST", "DELETE", "PUT", "PATCH"],
    "key": "posts-{id}:key"  // the key template that will be deleted after success request
}
```

10. skip - will skip all rules see "Routing"
```json
{
    "kind": "skip"
}
```

11. proxy - rule to change method or target host for proxing 
```json
{
    "kind": "proxy",
    "methods": ["GET"],  
    "path": "/posts/",  // optional
    "method": "POST",  // optional
    "target": "https://google.com" // optional
}
```
Will proxy all get methods as post request to the https://google.com/posts/

12. Glue
```json
{
  "kind": "glue",
  "model": "single", // "list" 
  "calls": {
    "user_info": "/users/{id}/info",
    "user_meta": "/users/{id}/meta",
    "posts": "/posts/?user={id}"
  },
  "mapping": {
    "posts": "posts"
  },  // {..info, ..meta, "posts": posts}
  "keys": {
    "posts": "user_id"
  }
}

TODO:
 
 - rule disable (manage api)

 - rate-limit with a sliding window
 - circuit-breaker 
 - retry by header Retry-After
 - block by condition
 - lock rule
 - hot cache

 - glue[advance]
 - transform Request
 - transform Responce 
 
 - trasing metrics
 - sentry
 - audit rule 

 - AUTH: jwt 
 
 - proxy filebody
 - proxy Location - bug
 - unix socket as target

 - control panel
