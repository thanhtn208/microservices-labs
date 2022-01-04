package main

import (
	"context"
	pb "demo/order"
	"demo/store"
	"encoding/json"
	"log"
	"net"

	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	port      = ":50051"
	aggregate = "Order"
	event     = "OrderCreated"
)

type server struct {
	pb.UnimplementedOrderServiceServer
}

// CreateOrder creates a new Order
func (s *server) CreateOrder(ctx context.Context, in *pb.Order) (*pb.OrderResponse, error) {
	store := store.OrderStore{}
	store.CreateOrder(in)
	go publishOrderCreated(in)
	return &pb.OrderResponse{IsSuccess: true}, nil
}

// publishOrderCreated publish an event via NATS server
func publishOrderCreated(order *pb.Order) {
	// Connect to NATS server
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	log.Println("Connected to " + nats.DefaultURL)
	defer natsConnection.Close()
	eventData, _ := json.Marshal(order)
	event := pb.EventStore{
		AggregateId:   order.OrderId,
		AggregateType: aggregate,
		EventId:       uuid.NewV4().String(),
		EventType:     event,
		EventData:     string(eventData),
	}
	subject := "Order.OrderCreated"
	data, _ := proto.Marshal(&event)
	// Publish message on subject
	natsConnection.Publish(subject, data)
	log.Println("Published message on subject " + subject)
}

// GetCustomers returns all customers by given filter
func (s *server) GetOrders(filter *pb.OrderFilter, stream pb.OrderService_GetOrdersServer) error {
	store := store.OrderStore{}
	orders := store.GetOrders()
	for _, order := range orders {
		if err := stream.Send(&order); err != nil { // Use stream.Recv() from Client
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	s.Serve(lis)
}
