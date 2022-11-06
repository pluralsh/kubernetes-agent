package mock_redis

//go:generate go run github.com/golang/mock/mockgen -source "../../redistool/expiring_hash.go" -destination "expiring_hash.go" -package "mock_redis"
