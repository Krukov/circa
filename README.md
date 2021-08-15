# circa

Fast, smart and simple sidecar proxy 
 - Simple configuration with file or managment api
 - Prometeus metrics
 - Failover cache
 - Ratelimits
 - RequestID
 - Indepotency 
 - CircutBreaker


TODO:
 - Manage config with http api
 - Proxy headers  - Fix gzip error 
 - Performance tests
 - early 
 - rate-limit with a bucket
 - circuit-breaker []
 - hot ? ( keep cache warm and regularly make a request)



How to add new cache 
 - Add new rule.Rule
 - Add new config parser config.Config getRuleFromOptions

How to add New Storage
 - Add new storage in storages module and add to DSN parser StorageFormDSN 