package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

var (
	name                string = "testnet"
	someMainnetContract string = "0x024dE48Fb640DB135B3dc85ef0FE2789e032FbCA2fca54E58aB8dB93ca22F767"
	contractMethod      string = "getName"
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then makes a contract call and prints the response.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	fmt.Println("Starting simpeCall example")
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	url := os.Getenv("INTEGRATION_BASE")
	fmt.Print("Connecting to the client at: ", url)
	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	contractAddress, err := utils.HexToFelt(someMainnetContract)
	if err != nil {
		panic(err)
	}
	fmt.Println("Contract Address: ", contractAddress)
	chainID1, err := clientv02.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("ChainID: ", chainID1)

	TxHash, err := clientv02.BlockWithTxHashes(context.Background(), rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err)
	}
	fmt.Println("BlockWithTxHashes: ", TxHash)

	BlockHashAndNumberOutput, _ := clientv02.BlockHashAndNumber(context.Background())
	blocknum, _ := clientv02.BlockNumber(context.Background())

	fmt.Println("BlockHashAndNumberOutput: ", BlockHashAndNumberOutput)
	fmt.Println("BlockNumber: ", blocknum)

	txnNums, _ := clientv02.BlockTransactionCount(context.Background(), rpc.BlockID{Tag: "latest"})
	fmt.Println("BlockTransactionCount: ", txnNums)

	// Make read contract call
	// tx := rpc.FunctionCall{
	// 	ContractAddress:    contractAddress,
	// 	EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
	// }

	// fmt.Println("Making Call() request")
	// callResp, err := clientv02.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, callResp[0]))
}
