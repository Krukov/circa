# circa

Fast, smart and simple cache proxy 

Request + Response (Message)         -|
↓                                     |
Storages - Redis, mem                -| 
↓                                     |  - {Config} 
Rules - simple, early, hit, ets      -|
↓                                     |
Handler - router rules and call them -|
↓
Server - Proxy server itself - Call handler

Configuration look like
```json
{
  "storages": {
    "redis": "redis://:password@localhost:9090/1?client_name=my&option",
    "fast": "mem://lru?size=1000"
  },
  "rules": {
    "{version}/users/{id}": {
      "{version|strip(v)}:users:{id}": {
        "ttl": "5h",
        "invalidation": {
          "/users/{id}/*": "!GET"
        },
        "storage": "redis",
        "cache_control": "public"
      },
      "users:{id}": {
        "ttl": "1h",
        "type": "fail",
        "storage": "localhost"
      },
      "options": {
        "timeout": "2s",
        "propagate": false
      }
    },
    "books": {
      "books:{authtorization|jwt(user_uid)}{cache_version}": {
        "ttl": "2d",
        "early_ttl": "1d"
      }
    },
    "*": {
      "proxy": {
        
      }
    }
  },
  "options": {
    "jwt_public_key_path": "/-/public/jwt", 
    "cache_header_name": "X-Cache-Keys",
    "ttl": "1h",
    "backend": "redis",
    "early_ttl": "30m",
    "etag": true,
    "refresh_query": "_refresh",
    "skip_query": "_skip",
    "manage_port": 8888,
    "metrics_port": 8080,
    "metrics_path": "/-/_metrics"
  }
}

```

roadmap
1) simple mem cache with config {"/book/{id}": {"key:{id}": {"ttl": "5m"}}
    [*] store rules in mem 
    [*] check route for regexp 
    [*] gen key
    [*] get backend options
    [*] get value from backend
    [*] proxy request and return result
2) config load from file [*]
2.1) logging + flag [*]
3) Add redis backend [*]
4) add 
   4.1) fail [*] 
   4.2) early cache [-]
   4.3) hit cache [-]
   4.2) early cache [-]
5.0) Metrics and monitoring [Progress]
5.1) k8 integration (sidecar?)
5) Manage config with http api
6) Manage config with backend 
7) Proxy headers 

9) Rules for all - global rules for example fail for all 
x) Do TODOs

```json
{
  "version": "1",
  "backends": {
    "{name}": "{dsn}"
  },
  "options": {...global_options},
  "rules": {
    "{url}": {
      "options": {
        ...url_options
      },
      "{key}": {
        ...key_options
      }
    }
  },
  "post/*": {
    "retry": {
      "type": "retry",
      "backoff": "10s",
      "count": 3
    },
    "cb": {
      "type": "circuit-breaker",
      "threshold": 0.5,
      "ttl": "5m"
    }
  }
}
```

Backends : memory, redis (client-side too), embended golang kv store
types: 
    - simple
    - early 
    - hit
    - fail 
    - rate-limit 
    - circuit-breaker
    - hot ? ( keep cache warm and regularly make requests)
    - retry
    - proxy-for metrics
    - check header consistency ( request-id )
    - idempotent key
    - concatenate



How to add new cache 
 - Add new rule.Rule
 - Add new config parser config.Config  getRuleFromOptions

How to add New Storage :
 - Add new storage in storages module and add to DSN parser StorageFormDSN 