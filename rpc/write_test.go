package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

func TestDeclareTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		DeclareTx     BroadcastAddDeployTxnType
		ExpectedResp  AddDeclareTransactionResponse
		ExpectedError error
	}
	// Class Hash
	//content, err := os.ReadFile("./tests/testContracts/testing_origin.contract_class.json")
	content, err := os.ReadFile("./tests/testContracts/token_bridge.sierra.json")

	require.NoError(t, err)

	var class ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)

	// Compiled Class Hash
	//content2, err := os.ReadFile("./tests/hello_world_compiled.casm.json")
	require.NoError(t, err)

	testSet := map[string][]testSetType{
		"devnet": {
			{
				DeclareTx: BroadcastDeclareTxnV2{
					Type:    TransactionType_Declare,
					Version: TransactionV2,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7c1119bf8ca3df81c877b17798e6edc697b5303ad820b4fabe6a43aafc1f17b"),
						utils.TestHexToFelt(t, "0xdae9b53c4818adac054351064269fa49833a02926a88becc23474773177ec5"),
					},
					Nonce:             utils.TestHexToFelt(t, "0x1"),
					MaxFee:            utils.TestHexToFelt(t, "0x2386f26fc10000"),
					SenderAddress:     utils.TestHexToFelt(t, "0x75d0acc17175e0187d0685bdeb6e778d72a2253c75de4df75e233fd12b036c"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x6766f1aafa6a11df6c3a6f7ff5e0ad3cdbdb85a6602b5697cc44c850043fe1d"),
					ContractClass:     class,
				},
				ExpectedResp: AddDeclareTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x47c217194148981d378156d8c493926687223c79619ad8321b850f697859f43")},
				ExpectedError: nil,
			},
		},
		"mainnet": {},
		"mock": {
			{
				DeclareTx: DeclareTxnV2{},
				ExpectedResp: AddDeclareTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3")},
				ExpectedError: nil,
			},
			{
				DeclareTx: DeclareTxnV3{
					Type:          TransactionType_Declare,
					Version:       TransactionV3,
					Signature:     []*felt.Felt{},
					Nonce:         utils.TestHexToFelt(t, "0x0"),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
						L2Gas: ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:                   "",
					PayMasterData:         []*felt.Felt{},
					SenderAddress:         utils.TestHexToFelt(t, "0x0"),
					CompiledClassHash:     utils.TestHexToFelt(t, "0x0"),
					ClassHash:             utils.TestHexToFelt(t, "0x0"),
					AccountDeploymentData: []*felt.Felt{},
				},
				ExpectedResp: AddDeclareTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x48776db363442bcfec44b979dbdab1f2033cb25c7b3950a0cd7c238bb5e4785")},
				ExpectedError: nil,
			},
		},
		"testnet": {{
			DeclareTx: DeclareTxnV1{},
			ExpectedResp: AddDeclareTransactionResponse{
				TransactionHash: utils.TestHexToFelt(t, "0x55b094dc5c84c2042e067824f82da90988674314d37e45cb0032aca33d6e0b9")},
			ExpectedError: errors.New("Invalid Params"),
		},
		},
	}[testEnv]

	for _, test := range testSet {
		if test.DeclareTx == nil && testEnv == "testnet" {
			declareTxJSON, err := os.ReadFile("./tests/write/declareTx.json")
			if err != nil {
				t.Fatal("should be able to read file", err)
			}
			var declareTx AddDeclareTxnInput
			require.Nil(t, json.Unmarshal(declareTxJSON, &declareTx), "Error unmarshalling decalreTx")
			test.DeclareTx = declareTx
		}

		resp, err := testConfig.provider.AddDeclareTransaction(context.Background(), test.DeclareTx)
		t.Log("resp", resp)
		t.Log("err", err)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
		} else {
			require.Equal(t, (*resp.TransactionHash).String(), (*test.ExpectedResp.TransactionHash).String())
		}

	}
}

func TestAddInvokeTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTx      BroadcastInvokeTxnType
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			InvokeTx: InvokeTxnV1{
				Type:    TransactionType_Invoke,
				Version: TransactionV1,
				Signature: []*felt.Felt{
					utils.TestHexToFelt(t, "0x4d7204d421567cf1f247f573f64ff71de1f6b1105fbfffa0761cec0af65badb"),
					utils.TestHexToFelt(t, "0x5d1f03ecfbf55ba215c0f80b3cbae6ef9fef42664927e3d604b12370c9602a1"),
				},
				Nonce:         utils.TestHexToFelt(t, "0x0"),
				SenderAddress: utils.TestHexToFelt(t, "0x4"),
				Calldata: []*felt.Felt{
					utils.TestHexToFelt(t, "0x1"),
					utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
					utils.TestHexToFelt(t, "0x3"),
					utils.TestHexToFelt(t, "0x75d0acc17175e0187d0685bdeb6e778d72a2253c75de4df75e233fd12b036c"),
					utils.TestHexToFelt(t, "0x3635c9adc5dea00000"),
					utils.TestHexToFelt(t, "0x0"),
				},
				MaxFee: utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
			},
			ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0x651e7e07eeba00c6be3593077266bcec8362db0905f14385b8efb9b30e9cce8")},
			ExpectedError: nil,
		},
		},
		"mainnet": {},
		"mock": {
			{
				InvokeTx:     InvokeTxnV1{SenderAddress: new(felt.Felt).SetUint64(123)},
				ExpectedResp: AddInvokeTransactionResponse{&felt.Zero},
				ExpectedError: &RPCError{
					Code:    ErrUnexpectedError.Code,
					Message: ErrUnexpectedError.Message,
					Data:    "Something crazy happened"},
			},
			{
				InvokeTx:      InvokeTxnV1{},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0xdeadbeef")},
				ExpectedError: nil,
			},
			{
				InvokeTx: InvokeTxnV3{
					Type:    TransactionType_Invoke,
					Version: TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x71a9b2cd8a8a6a4ca284dcddcdefc6c4fd20b92c1b201bd9836e4ce376fad16"),
						utils.TestHexToFelt(t, "0x6bef4745194c9447fdc8dd3aec4fc738ab0a560b0d2c7bf62fbf58aef3abfc5"),
					},
					Nonce:         utils.TestHexToFelt(t, "0xe97"),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x5af3107a4000",
						},
						L2Gas: ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:           "",
					PayMasterData: []*felt.Felt{},
					SenderAddress: utils.TestHexToFelt(t, "0x3f6f3bc663aedc5285d6013cc3ffcbc4341d86ab488b8b68d297f8258793c41"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x2"),
						utils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						utils.TestHexToFelt(t, "0x27c3334165536f239cfd400ed956eabff55fc60de4fb56728b6a4f6b87db01c"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x4c312760dfd17a954cdd09e76aa9f149f806d88ec3e402ffaf5c4926f568a42"),
						utils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x5"),
						utils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						utils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x7fe4fd616c7fece1244b3616bb516562e230be8c9f29668b46ce0369d5ca829"),
						utils.TestHexToFelt(t, "0x287acddb27a2f9ba7f2612d72788dc96a5b30e401fc1e8072250940e024a587"),
					},
					AccountDeploymentData: []*felt.Felt{},
				},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd")},
				ExpectedError: nil,
			},
		},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.AddInvokeTransaction(context.Background(), test.InvokeTx)
		if test.ExpectedError != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			require.Equal(t, *resp, test.ExpectedResp)
		}

	}
}

func TestAddDeployAccountTansaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		DeployTx      BroadcastAddDeployTxnType
		ExpectedResp  AddDeployAccountTransactionResponse
		ExpectedError error
	}
	testSet := map[string][]testSetType{
		"devnet": {{
			DeployTx: DeployAccountTxn{
				MaxFee:    utils.TestHexToFelt(t, "0x16345785d8a0000"),
				Type:      TransactionType_DeployAccount,
				Version:   TransactionV1,
				ClassHash: utils.TestHexToFelt(t, "0x1a736d6ed154502257f02b1ccdf4d9d1089f80811cd6acad48e6b6a9d1f2003"),
				Signature: []*felt.Felt{
					utils.TestHexToFelt(t, "0x17efc659a121c93cefabaea3fe40f03b43ef0ec82eaf0d37824e25c361220fd"),
					utils.TestHexToFelt(t, "0x63726c7b5385dd702979b4459dae04c060da3d98c1fc7e6e19f0028afa9aded"),
				},
				Nonce:               utils.TestHexToFelt(t, "0x0"),
				ContractAddressSalt: utils.TestHexToFelt(t, "0x3db4f4fad96e5444bdc1d0286f42948763af1b7bdb7873c08d28cb5129d4aac"),
				ConstructorCalldata: []*felt.Felt{
					utils.TestHexToFelt(t, "0x36e3b4894dcf3f7ff967c4cd748fe3c6aee367b4ce13b8676bd10e8aaafb1b4"),
					utils.TestHexToFelt(t, "0x0"),
				},
			},
			ExpectedResp: AddDeployAccountTransactionResponse{
				TransactionHash: utils.TestHexToFelt(t, "0x7b49482ee823e61ab561348cae844a4e550038be2f4678eb9e6a1b595bc48a6"),
				ContractAddress: utils.TestHexToFelt(t, "0x75d0acc17175e0187d0685bdeb6e778d72a2253c75de4df75e233fd12b036c")},
			ExpectedError: nil,
		}},
		"mainnet": {},
		"mock": {
			{
				DeployTx: DeployAccountTxn{},
				ExpectedResp: AddDeployAccountTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e"),
					ContractAddress: utils.TestHexToFelt(t, "0x0"),
				},
				ExpectedError: nil,
			},
			{
				DeployTx: DeployAccountTxnV3{
					Type:      TransactionType_DeployAccount,
					Version:   TransactionV3,
					ClassHash: utils.TestHexToFelt(t, "0x2338634f11772ea342365abd5be9d9dc8a6f44f159ad782fdebd3db5d969738"),
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x6d756e754793d828c6c1a89c13f7ec70dbd8837dfeea5028a673b80e0d6b4ec"),
						utils.TestHexToFelt(t, "0x4daebba599f860daee8f6e100601d98873052e1c61530c630cc4375c6bd48e3"),
					},
					Nonce:         new(felt.Felt),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x5af3107a4000",
						},
						L2Gas: ResourceBounds{
							MaxAmount:       "",
							MaxPricePerUnit: "",
						},
					},
					Tip:                 "",
					PayMasterData:       []*felt.Felt{},
					ContractAddressSalt: new(felt.Felt),
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x5cd65f3d7daea6c63939d659b8473ea0c5cd81576035a4d34e52fb06840196c"),
					},
				},
				ExpectedResp: AddDeployAccountTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e"),
					ContractAddress: utils.TestHexToFelt(t, "0x0")},
				ExpectedError: nil,
			},
		},
	}[testEnv]

	for _, test := range testSet {

		resp, err := testConfig.provider.AddDeployAccountTransaction(context.Background(), test.DeployTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
		} else {
			require.Equal(t, (*resp.TransactionHash).String(), (*test.ExpectedResp.TransactionHash).String())
		}

	}
}
