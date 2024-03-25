package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract)
	getAllAssets(contract)
	createAsset(contract)
	readAssetByID(contract)
	transferAssetAsync(contract)
	exampleErrorHandling(contract)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

func NewChaincodeActions(gw *gateway.Gateway, channelName, chaincodeName string) *ChaincodeActions {
    network, err := gw.GetNetwork(channelName)
    if err != nil {
        log.Fatalf("Failed to get network: %v", err)
    }
    contract := network.GetContract(chaincodeName)

    return &ChaincodeActions{contract: contract}
}

// PlantSeeds submits a transaction to plant seeds with a type and weight.
func (ca *ChaincodeActions) PlantSeeds(microgreenType string, seedWeight float64) {
    _, err := ca.contract.SubmitTransaction("PlantSeeds", microgreenType, fmt.Sprintf("%f", seedWeight))
    if err != nil {
        log.Fatalf("Failed to submit PlantSeeds transaction: %v", err)
    }
    fmt.Println("Seeds planted successfully")
}

// UpdateGermination submits a transaction to update the germination status of a batch.
func (ca *ChaincodeActions) UpdateGermination(batchID string) {
    _, err := ca.contract.SubmitTransaction("UpdateGermination", batchID)
    if err != nil {
        log.Fatalf("Failed to submit UpdateGermination transaction: %v", err)
    }
    fmt.Println("Germination updated successfully")
}

// RecordWatering submits a transaction to record a watering event for a batch.
func (ca *ChaincodeActions) RecordWatering(batchID string) {
    _, err := ca.contract.SubmitTransaction("RecordWatering", batchID)
    if err != nil {
        log.Fatalf("Failed to submit RecordWatering transaction: %v", err)
    }
    fmt.Println("Watering recorded successfully")
}

// UpdateToHarvested submits a transaction to update the batch to harvested with a weight.
func (ca *ChaincodeActions) UpdateToHarvested(batchID string, harvestWeight float64) {
    _, err := ca.contract.SubmitTransaction("UpdateToHarvested", batchID, fmt.Sprintf("%f", harvestWeight))
    if err != nil {
        log.Fatalf("Failed to submit UpdateToHarvested transaction: %v", err)
    }
    fmt.Println("Updated to harvested successfully")
}

// QueryBatchesByStatus evaluates a transaction to query batches by their status.
func (ca *ChaincodeActions) QueryBatchesByStatus(status string) string {
    result, err := ca.contract.EvaluateTransaction("QueryBatchesByStatus", status)
    if err != nil {
        log.Fatalf("Failed to evaluate QueryBatchesByStatus transaction: %v", err)
    }
    return string(result)
}

// QueryBatchesByType evaluates a transaction to query batches by their type.
func (ca *ChaincodeActions) QueryBatchesByType(microgreenType string) string {
    result, err := ca.contract.EvaluateTransaction("QueryBatchesByType", microgreenType)
    if err != nil {
        log.Fatalf("Failed to evaluate QueryBatchesByType transaction: %v", err)
    }
    return string(result)
}

// QueryBatchesByDateRange evaluates a transaction to query batches by a date range.
func (ca *ChaincodeActions) QueryBatchesByDateRange(startDate, endDate string) string {
    result, err := ca.contract.EvaluateTransaction("QueryBatchesByDateRange", startDate, endDate)
    if err != nil {
        log.Fatalf("Failed to evaluate QueryBatchesByDateRange transaction: %v", err)
    }
    return string(result)
}

// QueryBatch evaluates a transaction to query a single batch by ID.
func (ca *ChaincodeActions) QueryBatch(batchID string) string {
    result, err := ca.contract.EvaluateTransaction("QueryBatch", batchID)
    if err != nil {
        log.Fatalf("Failed to evaluate QueryBatch transaction: %v", err)
    }
    return string(result)
}
	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		fmt.Println("Error Details:")

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				fmt.Printf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}