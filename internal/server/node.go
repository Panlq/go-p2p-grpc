package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gpb "github/panlq-github/go-p2p-grpc/api/gen/pb/hello"
)

const KeyPrefix = "go-p2p-grpc"

type Node struct {
	Name string
	Addr string

	SDAddress string
	SDKV      api.KV

	Peers map[string]gpb.HelloServiceClient
}

func NewNode(conf Config) *Node {
	return &Node{Name: conf.NodeName, Addr: conf.NodeAddr, SDAddress: conf.ServiceDiscoveryAddress}
}

func (node *Node) SayHello(ctx context.Context, stream *gpb.HelloRequest) (*gpb.HelloReply, error) {
	return &gpb.HelloReply{Message: "Hello from " + node.Name}, nil
}

func (node *Node) StartListening() {
	grpcServer := grpc.NewServer() // n is for serving purpose
	gpb.RegisterHelloServiceServer(grpcServer, node)

	lis, err := net.Listen("tcp", node.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (node *Node) registerService() error {
	config := api.DefaultConfig()

	config.Address = node.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Println("failed to contact discovery service: ", err)
		return err
	}

	kv := consul.KV()
	pair := &api.KVPair{Key: KeyPrefix + node.Name, Value: []byte(node.Addr)}
	if _, err := kv.Put(pair, nil); err != nil {
		log.Println("failed to register service: ", err)
		return err
	}

	// store the kv for future use
	node.SDKV = *kv

	log.Println("Successfully registered service to consul.")
	return nil
}

func (node *Node) Start() error {
	node.Peers = make(map[string]gpb.HelloServiceClient)

	go node.StartListening()

	if err := node.registerService(); err != nil {
		return err
	}

	for {
		node.BroadcastMessage("Hello from " + node.Name)
		// sleep for 5 seconds
		time.Sleep(10 * time.Second)
	}
}

func (node *Node) BroadcastMessage(message string) {
	// get all nodes -- inefficient, but this is just an example
	kvpairs, _, err := node.SDKV.List(KeyPrefix, nil)
	if err != nil {
		log.Println("failed to get keypairs from service discovery", err)
		return
	}

	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, KeyPrefix+node.Name) == 0 {
			// skip self
			continue
		}

		if node.Peers[kventry.Key] == nil {
			fmt.Println("new member online: ", kventry.Key)
			// connection not established previously
			node.SetupClient(kventry.Key, string(kventry.Value))
		}
	}
}

func (node *Node) SetupClient(name string, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("unable to connect to %s:%v", addr, err)
		return
	}

	defer conn.Close()

	// store the connection
	node.Peers[name] = gpb.NewHelloServiceClient(conn)
	response, err := node.Peers[name].SayHello(context.Background(), &gpb.HelloRequest{Name: node.Name})
	if err != nil {
		log.Printf("error making request to %s:%v", name, err)
	}

	log.Printf("greeting from other node: %s", response.GetMessage())
}
