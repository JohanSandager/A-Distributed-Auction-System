package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	pb "auction/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpc_client       pb.AuctionClient
	unique_identifier int64
}

var sAddress = flag.String("sAddress", "", "Server address")

func main() {
	flag.Parse()

	connection := ConnectToServer(*sAddress)
	client := Client{
		grpc_client:       GetClientFromConnection(connection),
		unique_identifier: GetUniqueIdentifier(),
	}
	go RunProgram(&client)
	time.Sleep(1 * time.Minute)
}

func RunProgram(client *Client) {
	scanner := bufio.NewScanner(os.Stdin)
	log.Print("Ready to take bid:")

	for scanner.Scan() {
		input_text := scanner.Text()
		input, err := strconv.Atoi(input_text)
		if input_text == "result" {
			GetResult(client)
			continue
		} else if err != nil {
			log.Print("Invalid input")
			continue
		}
		SendBid(input, client)
	}
}

// From handin 3
func ConnectToServer(server_address string) *grpc.ClientConn {
	port := FindSomeNode(5, 5006, server_address)
	address := server_address + ":" + port
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to %v", server_address)
	}

	log.Printf("Connected to the server at %v\n", address)

	return connection
}

// From handin3
func GetClientFromConnection(connection *grpc.ClientConn) pb.AuctionClient {
	client := pb.NewAuctionClient(connection)
	return client
}

func SendBid(value int, client *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	response, err := client.grpc_client.Bid(ctx, &pb.SendBidMessage{
		UniqueIdentifier: client.unique_identifier,
		Bid:              int64(value),
	})

	if err != nil {
		log.Printf("Some error this is %v", err)
		client.grpc_client = GetClientFromConnection(ConnectToServer(*sAddress)) // The "interface" node crashed, finding new node
		SendBid(value, client)
	} else if response.Status == "Success" {
		log.Printf("The bid of %v was successfull", value)
	} else if response.Status == "Failure" {
		log.Printf("The bid of %v was not successfull, due to invalid bid", value)
	} else if response.Status == "Exception" {
		log.Printf("The bid of %v was not successfull, due to system failure", value)
	} else {
		log.Fatalf("ERROR")
	}
}

func GetResult(client *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	response, err := client.grpc_client.Result(ctx, &pb.RequestResultMessage{})

	if err != nil {
		client.grpc_client = GetClientFromConnection(ConnectToServer(*sAddress)) // The "interface" node crashed, finding new node
		GetResult(client)
	} else {
		am_i_highest := client.unique_identifier == response.CurrentHighestBidder
		log.Printf("Highest bid: %v \n Highest bidder is: %v \n Is it me? %v", response.Result, response.CurrentHighestBidder, am_i_highest)
	}
}

func GetUniqueIdentifier() int64 {
	return int64(rand.Intn(999999999999999999))
}

func FindSomeNode(system_size int, initial_port int, address string) string {
	for i := system_size; i > 0; i-- {
		port := strconv.Itoa(initial_port - i)
		log.Printf("Trying at port: %v", port)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		ln, err := grpc.DialContext(ctx, address+":"+port, grpc.WithInsecure(), grpc.WithBlock())
		if err == nil {
			log.Printf("Someone at port: %v", port)
			ln.Close()
			return port
		}
	}
	log.Fatalf("No space left")
	return "Dosn't happen"
}
