package integration

import (
	"context"
	"testing"
	"time"

	pb "github.com/HelenaBlack/anti-bruteforce/api/gen" //nolint:depguard
	"github.com/stretchr/testify/assert"                //nolint:depguard
	"github.com/stretchr/testify/require"               //nolint:depguard
	"google.golang.org/grpc"                            //nolint:depguard
	"google.golang.org/grpc/credentials/insecure"       //nolint:depguard
)

func TestIntegration_AuthFlow(t *testing.T) {
	// Skip if server not running
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skip("Server not available")
	}
	defer func() { _ = conn.Close() }()
	client := pb.NewAntibruteforceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	login := "testuser"
	ip := "192.168.1.1"

	// 1. Reset buckets
	_, err = client.Reset(ctx, &pb.ResetRequest{Login: login, Ip: ip})
	require.NoError(t, err)

	// Try multiple times to hit limit (assuming default limit N=10)
	for i := 0; i < 10; i++ {
		resp, err := client.Check(context.Background(), &pb.CheckRequest{Login: login, Password: "pw", Ip: ip})
		require.NoError(t, err)
		assert.True(t, resp.Ok)
	}

	// 11th time should fail
	resp, err := client.Check(context.Background(), &pb.CheckRequest{Login: login, Password: "pw", Ip: ip})
	require.NoError(t, err)
	assert.False(t, resp.Ok)

	// Add IP to whitelist
	_, err = client.AddToWhitelist(context.Background(), &pb.SubnetRequest{Subnet: "192.168.1.0/24"})
	require.NoError(t, err)

	// Should be ok now despite limit
	resp, err = client.Check(context.Background(), &pb.CheckRequest{Login: login, Password: "pw", Ip: ip})
	require.NoError(t, err)
	assert.True(t, resp.Ok)
}
