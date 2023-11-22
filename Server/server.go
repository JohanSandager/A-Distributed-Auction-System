package main

import (
	pb "auction/grpc"
	"context"

	"flag"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedAuctionServer
	address               string
	port                  string
	initial_port          string
	client_array_size_int int
	coordinator           string
	IsCoordinator         bool
}

var server_array_size = flag.String("serverArraySize", "", "Max size of the server array")

func main() {

	flag.Parse()
	client_array_size_int, _ := strconv.Atoi(*server_array_size)
	client_IP := GetOutboundIP()
	initial_port := 5000
	port := FindAvailablePort(GetOutboundIP(), client_array_size_int, initial_port)

	client := &Server{
		address:               client_IP,
		port:                  port,
		initial_port:          strconv.Itoa(initial_port),
		client_array_size_int: client_array_size_int,
		IsCoordinator:         false,
	}

	grpcServer := grpc.NewServer()
	listener, _ := net.Listen("tcp", client.address+":"+client.port)
	pb.RegisterAuctionServer(grpcServer, client)

	// Starting listening
	log.Print("Listening at: " + client.address + ":" + client.port)
	go grpcServer.Serve(listener)

	// Starting acting
	go RunProgram(client)
	time.Sleep(1 * time.Minute)

}

func RunProgram(server *Server) {
	log.Println("Calling election, coordinator before election: " + server.coordinator)
	server.CallElection(context.Background(), &pb.CallElectionMessage{})
	log.Println("Coordinator after election: " + server.coordinator)
}
func (client *Server) CallElection(context context.Context, call_election_message *pb.CallElectionMessage) (*pb.CallElectionResponseMessage, error) {
	log.Print("Election called!")

	address := client.address
	client_port, _ := strconv.Atoi(client.port)
	initial_port, _ := strconv.Atoi(client.initial_port)
	max_port := initial_port + client.client_array_size_int
	response_received := false

	if max_port == client_port {
		log.Print("I'm the coordinator, cause I'm the the biggest guy/girl in here!!!!")
		for i := client_port; i > initial_port; i-- {
			port := strconv.Itoa(initial_port + i)
			SendCoordinator(client, address, port)
		}
	}

	for i := client_port + 1; i <= max_port; i++ {
		port := strconv.Itoa(i)
		if SendElection(address, port) {
			response_received = true
		}
	}

	if !response_received {
		log.Print("No response means I'm the leader muhahaha")
		for i := client_port; i > initial_port; i-- {
			port := strconv.Itoa(i)
			SendCoordinator(client, address, port)
		}
	}

	return &pb.CallElectionResponseMessage{}, nil
}

func (server *Server) AssertCoordinator(context context.Context, message *pb.AssertCoordinatorMessage) (*pb.AssertCoordinatorResponseMessage, error) {
	server.coordinator = message.Port
	if message.Port == server.port {
		server.IsCoordinator = true
	} else {
		server.IsCoordinator = false
	}
	log.Print("Someone thinks they are coordinator, this guy eh: " + message.Port)
	return &pb.AssertCoordinatorResponseMessage{
		Port: server.port,
	}, nil
}

func SendCoordinator(server *Server, address string, port string) {
	server.IsCoordinator = true
	log.Printf("Telling my subjects I'm the boss around here, subject: " + port)
	server_address := address + ":" + port
	connection, _ := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpc_client := pb.NewAuctionClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	grpc_client.AssertCoordinator(ctx, &pb.AssertCoordinatorMessage{
		Port: server.port,
	})
}

func SendElection(address string, port string) bool {
	server_address := address + ":" + port
	connection, _ := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pb.NewAuctionClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err := client.CallElection(ctx, &pb.CallElectionMessage{})
	if err == nil {
		return true
	}
	return false
}

func FindAvailablePort(address string, system_size int, initial_port int) string {
	for i := 0; i < system_size; i++ {
		port := strconv.Itoa(initial_port + i)
		timeout := time.Duration(1 * time.Second)
		_, err := net.DialTimeout("tcp", address+":"+port, timeout)
		if err != nil {
			log.Printf("Hey I'm at port: %v", port)
			return port
		}
	}
	log.Fatalf("No space left")
	return "Dosn't happen"
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
