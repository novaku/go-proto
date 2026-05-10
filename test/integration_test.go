package main_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/di"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/pkg/framework"
)

func TestGuestbookIntegration(t *testing.T) {
	// 1. Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&model.GuestbookEntry{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	// 2. Start Server on a random port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port

	controllerFactory := di.NewGuestbookControllerFactory(db)
	guestbookController := controllerFactory.CreateController()

	serverErrCh := make(chan error, 1)
	go func() {
		testSrv := framework.NewServer(port)
		testSrv.RegisterService(func(s *grpc.Server) {
			pb.RegisterGuestbookServiceServer(s, guestbookController)
		})
		// Close test listener so the server can bind to it
		lis.Close()

		if err := testSrv.Run(); err != nil {
			serverErrCh <- err
		}
	}()

	// 3. Wait for server to start
	time.Sleep(1 * time.Second)

	// 4. Client Connection
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGuestbookServiceClient(conn)

	// 5. Test AddEntry
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	addResp, err := c.AddEntry(ctx, &pb.AddEntryRequest{Name: "Alice", Email: "alice@example.com", Message: "Hello World"})
	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}
	if !addResp.Success {
		t.Errorf("AddEntry returned success=false, error: %s", addResp.Error)
	}

	// 6. Test ListEntries
	// 6. Test ListEntries
	listResp, err := c.ListEntries(ctx, &pb.ListEntriesRequest{})
	if err != nil {
		t.Fatalf("ListEntries failed: %v", err)
	}

	if len(listResp.Data) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(listResp.Data))
	}
	if len(listResp.Data) > 0 {
		if listResp.Data[0].Name != "Alice" {
			t.Errorf("Expected name Alice, got %s", listResp.Data[0].Name)
		}
		if listResp.Data[0].Message != "Hello World" {
			t.Errorf("Expected message Hello World, got %s", listResp.Data[0].Message)
		}
		if listResp.Data[0].Email != "alice@example.com" {
			t.Errorf("Expected email alice@example.com, got %s", listResp.Data[0].Email)
		}
	}
}
