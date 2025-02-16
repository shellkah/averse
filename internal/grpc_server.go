package internal

import (
	"context"
	"time"

	"github.com/shellkah/averse/cachepb"

	"github.com/shellkah/goutte"
)

type server struct {
	cachepb.UnimplementedCacheServiceServer
	cache *goutte.Cache[string, string]
}

func NewServer(cacheCapacity int) *server {
	return &server{cache: goutte.NewCache[string, string](cacheCapacity)}
}

func (s *server) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	value, found := s.cache.Get(req.Key)
	return &cachepb.GetResponse{
		Value: value,
		Found: found,
	}, nil
}

func (s *server) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.SetResponse, error) {
	s.cache.Set(req.Key, req.Value)
	return &cachepb.SetResponse{Success: true}, nil
}

func (s *server) SetWithTTL(ctx context.Context, req *cachepb.SetWithTTLRequest) (*cachepb.SetResponse, error) {
	ttl := time.Duration(req.TtlSeconds) * time.Second
	s.cache.SetWithTTL(req.Key, req.Value, ttl)
	return &cachepb.SetResponse{Success: true}, nil
}

func (s *server) Delete(ctx context.Context, req *cachepb.DeleteRequest) (*cachepb.DeleteResponse, error) {
	s.cache.Delete(req.Key)
	return &cachepb.DeleteResponse{Success: true}, nil
}

func (s *server) Dump(ctx context.Context, req *cachepb.DumpRequest) (*cachepb.DumpResponse, error) {
	s.cache.Dump()
	return &cachepb.DumpResponse{}, nil
}
