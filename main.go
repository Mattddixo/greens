package main

import (
	"microgreens/chaincodeactions"
	"microgreens/cli"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	// Establish gRPC connection
	clientConnection := chaincodeactions.NewGrpcConnection()
	defer clientConnection.Close()

	// Setup identity and sign function
	id := chaincodeactions.NewIdentity()
	sign := chaincodeactions.NewSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Prepare for chaincode interaction
	chaincodeName, channelName := getOverrides()
	actions := chaincodeactions.NewChaincodeActions(gw, channelName, chaincodeName)

	// Initialize and run CLI tool
	cli.RunCLI(actions)
}

// Helper function to handle environment overrides for chaincode and channel names
func getOverrides() (string, string) {
	chaincodeName := "basic" // Default chaincode name
	if ccName := os.Getenv("CHAINCODE_NAME"); ccName != "" {
		chaincodeName = ccName
	}

	channelName := "mychannel" // Default channel name
	if chName := os.Getenv("CHANNEL_NAME"); chName != "" {
		channelName = chName
	}

	return chaincodeName, channelName
}
