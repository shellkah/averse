// server_test.go
package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/shellkah/averse/cachepb"
	"github.com/shellkah/averse/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startTestServer(t *testing.T) (cachepb.CacheServiceClient, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	svc := internal.NewServer(100)
	cachepb.RegisterCacheServiceServer(grpcServer, svc)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("gRPC server stopped serving: %v", err)
		}
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	client := cachepb.NewCacheServiceClient(conn)

	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
		lis.Close()
	}

	return client, cleanup
}

func TestSetAndGet(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setReq := &cachepb.SetRequest{
		Key:   "foo",
		Value: "bar",
	}
	setRes, err := client.Set(ctx, setReq)
	if err != nil {
		t.Fatalf("Set RPC failed: %v", err)
	}
	if !setRes.Success {
		t.Fatalf("Expected set success, got false")
	}

	getReq := &cachepb.GetRequest{
		Key: "foo",
	}
	getRes, err := client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get RPC failed: %v", err)
	}
	if !getRes.Found {
		t.Fatalf("Expected key 'foo' to be found")
	}
	if getRes.Value != "bar" {
		t.Fatalf("Expected value 'bar', got '%s'", getRes.Value)
	}
}

func TestSetWithTTL(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ttlSeconds := int64(1)
	setTTLReq := &cachepb.SetWithTTLRequest{
		Key:        "temp",
		Value:      "data",
		TtlSeconds: ttlSeconds,
	}
	setTTLRes, err := client.SetWithTTL(ctx, setTTLReq)
	if err != nil {
		t.Fatalf("SetWithTTL RPC failed: %v", err)
	}
	if !setTTLRes.Success {
		t.Fatalf("Expected SetWithTTL success, got false")
	}

	getReq := &cachepb.GetRequest{Key: "temp"}
	getRes, err := client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get RPC failed: %v", err)
	}
	if !getRes.Found || getRes.Value != "data" {
		t.Fatalf("Expected to find key 'temp' with value 'data'")
	}

	time.Sleep(2 * time.Second)

	getRes, err = client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get RPC failed after TTL expiration: %v", err)
	}
	if getRes.Found {
		t.Fatalf("Expected key 'temp' to have expired")
	}
}

func TestDelete(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setReq := &cachepb.SetRequest{
		Key:   "deleteKey",
		Value: "to be deleted",
	}
	_, err := client.Set(ctx, setReq)
	if err != nil {
		t.Fatalf("Set RPC failed: %v", err)
	}

	delReq := &cachepb.DeleteRequest{
		Key: "deleteKey",
	}
	delRes, err := client.Delete(ctx, delReq)
	if err != nil {
		t.Fatalf("Delete RPC failed: %v", err)
	}
	if !delRes.Success {
		t.Fatalf("Expected delete success, got false")
	}

	getReq := &cachepb.GetRequest{Key: "deleteKey"}
	getRes, err := client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get RPC failed: %v", err)
	}
	if getRes.Found {
		t.Fatalf("Expected key 'deleteKey' to be deleted")
	}
}
