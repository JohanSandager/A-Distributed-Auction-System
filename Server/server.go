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
	address                string
	port                   string
	initial_port           string
	client_array_size_int  int
	coordinator            string
	IsCoordinator          bool
	bidders_to_highest_bid map[int]int
	highest_bid            int
	highest_bidder         int
	auction_end            int
	is_first_bid           bool
}

var server_array_size = flag.String("serverArraySize", "", "Max size of the server array")

func main() {

	flag.Parse()
	client_array_size_int, _ := strconv.Atoi(*server_array_size)
	client_IP := GetOutboundIP()
	initial_port := 5000
	port := FindAvailablePort(GetOutboundIP(), client_array_size_int, initial_port)

	client := &Server{
		address:                client_IP,
		port:                   port,
		initial_port:           strconv.Itoa(initial_port),
		client_array_size_int:  client_array_size_int,
		IsCoordinator:          false,
		bidders_to_highest_bid: map[int]int{},
		highest_bid:            0,
		is_first_bid:           true,
	}

	grpcServer := grpc.NewServer()
	listener, _ := net.Listen("tcp", client.address+":"+client.port)
	pb.RegisterAuctionServer(grpcServer, client)

	// Starting listening
	log.Print("Listening at: " + client.address + ":" + client.port)
	go grpcServer.Serve(listener)

	// Starting acting
	go RunProgram(client)
	time.Sleep(10 * time.Minute)

}

func RunProgram(server *Server) {
	log.Print("Joining in, calling election")
	server.CallElection(context.Background(), &pb.CallElectionMessage{})
}

func (server *Server) Bid(ctx context.Context, send_bid_message *pb.SendBidMessage) (*pb.ResponseBidMessage, error) {
	auction_not_over := time.Now().Unix() < int64(server.auction_end) || server.is_first_bid
	log.Print(auction_not_over)
	if auction_not_over {
		if server.IsCoordinator {
			if server.is_first_bid {
				server.is_first_bid = false
				server.auction_end = int(time.Now().Unix()) + 100
				end_time := time.Unix(int64(server.auction_end), 0).String()
				log.Printf("Auction started, ending %v", end_time)
			}
			log.Print("Server is coordinator, and recieved bid")
			_, exists := server.bidders_to_highest_bid[int(send_bid_message.UniqueIdentifier)]
			if exists && int(send_bid_message.Bid) > server.highest_bid { // If the bidder is already registrered and the bid is valid
				log.Print("The bidder is already registrered and the bid is valid")
				updated_successfully := HandleValidBid(server, send_bid_message) // Handles the bid, returns true if atleast one other node was updated
				if updated_successfully {
					return &pb.ResponseBidMessage{ // Returns success to the bidder
						Status: "Success",
					}, nil
				}
				return &pb.ResponseBidMessage{ // Returns failure, because the system could not update
					Status: "Exception",
				}, nil
			} else if !exists && int(send_bid_message.Bid) > server.highest_bid { // If the bidder is not registrered but the bid is valid
				log.Print("The bidder is not registrered but the bid is valid")
				RegisterNewBidder(server, send_bid_message)
				updated_successfully := HandleValidBid(server, send_bid_message) // Handles the bid, returns true if atleast one other node was updated
				if updated_successfully {
					return &pb.ResponseBidMessage{ // Returns success to the bidder
						Status: "Success",
					}, nil
				} else {
					log.Printf("No replcias to update, system is not secure")
					return &pb.ResponseBidMessage{
						Status: "Exception",
					}, nil
				}
			} else if int(send_bid_message.Bid) <= server.highest_bid { // Bid is invalid
				log.Print("Recieved invalid bid")
				return &pb.ResponseBidMessage{
					Status: "Failure",
				}, nil
			}
		} else if send_bid_message.FromCoordinator { // If the bid call came from the coordiantor, it is an update
			log.Print("Server is updating")
			updated_successfully := HandleUpdate(server, send_bid_message) // Handles the update
			if updated_successfully {
				return &pb.ResponseBidMessage{ // Returns success to the coordinator
					Status: "Success",
				}, nil
			}
			return &pb.ResponseBidMessage{ // Returns failure, because the node could not update
				Status: "Exception",
			}, nil
		}
	} else if !auction_not_over {
		log.Print("Auction over bro")
		return &pb.ResponseBidMessage{
			Status: "Failure",
		}, nil
	}
	// If the node is not the coordinator (leader) the request is fowarded to the coordinator
	log.Print("The node is not the coordinator (leader) the request is fowarded to the coordinator")
	server_address := server.address + ":" + server.coordinator
	connection, _ := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpc_client := pb.NewAuctionClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	coordinator_repsonse, err := grpc_client.Bid(ctx, send_bid_message)
	if err != nil {
		server.CallElection(context.Background(), &pb.CallElectionMessage{})
	}
	return coordinator_repsonse, nil
}

func (server *Server) Result(ctx context.Context, request_result_message *pb.RequestResultMessage) (*pb.ResultResponseMessage, error) {
	if server.IsCoordinator {
		return &pb.ResultResponseMessage{
			Result:               int64(server.highest_bid),
			CurrentHighestBidder: int64(server.highest_bidder),
		}, nil
	}
	log.Print("The node is not the coordinator (leader) the request is fowarded to the coordinator")
	server_address := server.address + ":" + server.coordinator
	connection, _ := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpc_client := pb.NewAuctionClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	coordinator_repsonse, err := grpc_client.Result(ctx, request_result_message)
	if err != nil {
		server.CallElection(context.Background(), &pb.CallElectionMessage{})
	}
	return coordinator_repsonse, nil
}

func HandleValidBid(server *Server, send_bid_message *pb.SendBidMessage) bool {
	server.highest_bid = int(send_bid_message.Bid)
	server.highest_bidder = int(send_bid_message.UniqueIdentifier)
	one_or_more_updated_sucessfully := UpdateReplicas(server, send_bid_message)
	return one_or_more_updated_sucessfully
}

func HandleUpdate(server *Server, send_bid_message *pb.SendBidMessage) bool {
	server.highest_bid = int(send_bid_message.Bid)
	server.highest_bidder = int(send_bid_message.UniqueIdentifier)
	server.is_first_bid = send_bid_message.IsFirstBid
	server.auction_end = int(send_bid_message.EndTime)
	return true
}

func UpdateReplicas(server *Server, send_bid_message *pb.SendBidMessage) bool {
	log.Print("Updating replicas")
	one_or_more_updated_sucessfully := false
	server_port, _ := strconv.Atoi(server.port)
	initial_port, _ := strconv.Atoi(server.initial_port)
	max_port := initial_port + server.client_array_size_int
	for i := initial_port; i < max_port; i++ {
		if i != server_port {
			port := strconv.Itoa(i)
			server_address := server.address + ":" + port
			connection, _ := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			grpc_client := pb.NewAuctionClient(connection)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			response, _ := grpc_client.Bid(ctx, &pb.SendBidMessage{
				UniqueIdentifier: send_bid_message.UniqueIdentifier,
				Bid:              send_bid_message.Bid,
				FromCoordinator:  true,
				EndTime:          int64(server.auction_end),
				IsFirstBid:       false,
			})
			if response != nil && response.Status == "Success" {
				log.Print("Made an successful update")
				one_or_more_updated_sucessfully = true
			}
		}
	}
	return one_or_more_updated_sucessfully
}

func RegisterNewBidder(server *Server, send_bid_message *pb.SendBidMessage) {
	server.bidders_to_highest_bid[int(send_bid_message.UniqueIdentifier)] = int(send_bid_message.Bid)
}

/*-------------------------------------------------------------------
ELECTION CODE FROM HANDIN 4
--------------------------------------------------------------------*/

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
	connection, err := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to dial %s: %v", server_address, err)
	}
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
