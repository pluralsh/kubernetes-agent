package mock_redis

//go:generate mockgen.sh -source "../../redistool/expiring_hash.go" -destination "expiring_hash.go" -package "mock_redis"
