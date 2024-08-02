package rpc

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

// TestClassAt tests the ClassAt function.
//
// The function tests the ClassAt function by creating different test sets for different environments
// (mock, testnet, and mainnet). It then iterates over each test set and performs the following steps:
//   - Calls the ClassAt function with the specified block tag and contract address.
//   - Checks the response type and performs the following actions based on the type:
//   - If the response type is DeprecatedContractClass or ContractClass:
//   - Checks if the program code exists in the response object.
//   - Checks if the expected operation exist in the provided contract address
//   - If the response type is of unknown type: log and fail the test
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress   *felt.Felt
		ExpectedOperation string
		Block             BlockID
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedOperation: "0xdeadbeef",
				Block:             WithBlockNumber(58344),
			},
		},
		"devnet": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x41A78E741E5AF2FEC34B695679BC6891742439F7AFB8484ECD7766661AD02BF"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("deployContract").String(),
				Block:             WithBlockTag("latest"),
			},
		},
		"testnet": {
			// v0 contract
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("getPublicKey").String(),
				Block:             WithBlockNumber(58344),
			},
			// v2 contract
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("name_get").String(),
				Block:             WithBlockNumber(65168),
			},
		},
		"mainnet": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x004b3d247e79c58e77c93e2c52025d0bb1727957cc9c33b33f7216f369c77be5"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("get_name").String(),
				Block:             WithBlockNumber(643360),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		fmt.Println("TestClassAt Start, test: ", testEnv)
		require := require.New(t)
		resp, err := testConfig.provider.ClassAt(context.Background(), test.Block, test.ContractAddress)
		require.NoError(err)
		fmt.Println("resp: ", resp)
		switch class := resp.(type) {
		case *DeprecatedContractClass:
			require.NotEmpty(class.Program, "code should exist")
			require.Condition(func() bool {
				fmt.Println("class.DeprecatedEntryPointsByType.External length: ", len(class.DeprecatedEntryPointsByType.External))
				for _, deprecatedCairoEntryPoint := range class.DeprecatedEntryPointsByType.External {
					fmt.Println("deprecatedCairoEntryPoint: ", deprecatedCairoEntryPoint.Selector.String())
					fmt.Println("test.ExpectedOperation: ", test.ExpectedOperation)
					if test.ExpectedOperation == deprecatedCairoEntryPoint.Selector.String() {
						return true
					}
				}
				return false
			}, "operation not found in the class")
		case *ContractClass:
			require.NotEmpty(class.SierraProgram, "code should exist")

			require.Condition(func() bool {
				for _, entryPointsByType := range class.EntryPointsByType.External {
					if test.ExpectedOperation == entryPointsByType.Selector.String() {
						return true
					}
				}
				return false
			}, "operation not found in the class")
		default:
			t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
		}
	}
}

// TestClassHashAt tests the ClassHashAt function.
//
// This function tests the behavior of the ClassHashAt function by providing
// different test cases for the contract hash and expected class hash. It
// verifies if the returned class hash matches the expected class hash and
// if there are any differences between the two. It also checks if the
// returned class hash is not nil. The function takes in a testing.T
// parameter and does not return anything.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
//
// itachi Pass
func TestClassHashAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      *felt.Felt
		ExpectedClassHash *felt.Felt
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:      utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		"devnet": {
			{
				ContractHash:      utils.TestHexToFelt(t, "0x41A78E741E5AF2FEC34B695679BC6891742439F7AFB8484ECD7766661AD02BF"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x7B3E05F48F0C69E4A65CE5E076A66271A527AFF2C34CE1083EC6E1526997A69"),
			},
		},
		"testnet": {
			// v0 contracts
			{
				ContractHash:      utils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x054328a1075b8820eb43caf0caa233923148c983742402dcfc38541dd843d01a"),
			},
			{
				ContractHash:      utils.TestHexToFelt(t, "0x073ad76dcf68168cbf68ea3ec0382a3605f3deaf24dc076c355e275769b3c561"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
			},
			// v2 contract
			{
				ContractHash:      utils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
			},
		},
		"mainnet": {
			{
				ContractHash:      utils.TestHexToFelt(t, "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x4c53698c9a42341e4123632e87b752d6ae470ddedeb8b0063eaa2deea387eeb"),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		fmt.Println("test: ", testEnv)
		classhash, err := testConfig.provider.ClassHashAt(context.Background(), WithBlockTag("latest"), test.ContractHash)
		require.NoError(err)
		require.NotEmpty(classhash, "should return a class")
		require.Equal(test.ExpectedClassHash, classhash)
	}
}

// TestClass is a test function that tests the behavior of the Class function.
//
// It creates a test configuration and defines a testSet containing different scenarios
// for testing the Class function. The testSet is a map where the keys represent the
// test environment and the values are an array of testSetType structs. Each struct in
// the array represents a specific test case with properties such as BlockID,
// ClassHash, ExpectedProgram, and ExpectedEntryPointConstructor.
//
// The function iterates over each test case in the testSet and performs the following steps:
// - Calls the Class function with the appropriate parameters.
// - Handles the response based on its type:
//   - If the response is of type DeprecatedContractClass:
//   - Checks if the class program starts with the expected program.
//   - If not, it reports an error.
//   - If the response is of type ContractClass:
//   - Checks if the class program ends with the expected program.
//   - Compares the constructor entry point with the expected entry point constructor.
//   - If they are not equal, it reports an error.
//
// The function is used for testing the behavior of the Class function in different scenarios.
//
// Parameters:
// - t: A *testing.T object used for reporting test failures and logging
// Returns:
//
//	none
func TestClass(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                       BlockID
		ClassHash                     *felt.Felt
		ExpectedProgram               string
		ExpectedEntryPointConstructor SierraEntryPoint
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:         WithBlockTag("pending"),
				ClassHash:       utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
		},
		"devnet": {
			//18
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     utils.TestHexToFelt(t, "0x5e30d248d37fbb0f3765c1276ca6a932af8ac45cf01c069d2c7f4539bb84c23"),
				ExpectedProgram:               "0x28c09039a0120c",
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 6, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
		"testnet": {
			// v0 class
			{
				BlockID:         WithBlockTag("latest"),
				ClassHash:       utils.TestHexToFelt(t, "0x2a747051e210a5958b3632320751a000c4ec66c8569385a37d18c0040b0b1eb"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
			// v2 classes
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     utils.TestHexToFelt(t, "0x00816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253"),
				ExpectedProgram:               utils.TestHexToFelt(t, "0x576402000a0028a9c00a010").String(),
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 34, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
				ExpectedProgram:               utils.TestHexToFelt(t, "0xe70d09071117174f17170d4fe60d09071117").String(),
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 2, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
		"mainnet": {
			// v2 class
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     utils.TestHexToFelt(t, "0x029927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
				ExpectedProgram:               utils.TestHexToFelt(t, "0x9fa00900700e00712e12500712e").String(),
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 32, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		resp, err := testConfig.provider.Class(context.Background(), test.BlockID, test.ClassHash)
		require.NoError(err)

		switch class := resp.(type) {
		case *DeprecatedContractClass:
			if !strings.HasPrefix(class.Program, test.ExpectedProgram) {
				t.Fatal("code should exist")
			}
		case *ContractClass:
			t.Log("class.SierraProgram: ", class.SierraProgram[len(class.SierraProgram)-1].String())
			t.Log("class.EntryPointsByType.Constructor[0] ", class.EntryPointsByType.Constructor[0])
			require.Equal(class.SierraProgram[len(class.SierraProgram)-1].String(), test.ExpectedProgram)
			require.Equal(class.EntryPointsByType.Constructor[0], test.ExpectedEntryPointConstructor)
		default:
			t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
		}
	}
}

// TestStorageAt tests the StorageAt function.
//
// It tests the StorageAt function by creating a test environment and executing
// a series of test cases. Each test case represents a different scenario, such
// as different contract hashes, storage keys, and expected values. The function
// checks if the actual value returned by the StorageAt function matches the
// expected value. If there is a mismatch, the test fails and an error is
// reported.
//
// Parameters:
// - t: The testing.T instance used for reporting test failures and logging
// Returns:
//
//	none
func TestStorageAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash  *felt.Felt
		StorageKey    string
		Block         BlockID
		ExpectedValue string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0xdeadbeef"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0xdeadbeef",
			},
		},
		"devnet": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				StorageKey:    "ERC20_name",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x4574686572",
			},
		},
		"testnet": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:    "_signer",
				Block:         WithBlockNumber(69399),
				ExpectedValue: "0x38bd4cad8706e3a5d167ef7af12e28268c6122df3e0e909839a103039871b9e",
			},
		},
		"mainnet": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x7f72660ca40b8ca85f9c0dd38db773f17da7a52f5fc0521cb8b8d8d44e224b8",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		value, err := testConfig.provider.StorageAt(context.Background(), test.ContractHash, test.StorageKey, test.Block)
		//t.Log("value: ", value)
		require.NoError(err)
		require.EqualValues(test.ExpectedValue, value)
	}
}

// TestNonce is a test function for testing the Nonce functionality.
//
// It initializes a test configuration, sets up a test data set, and then performs a series of tests.
// The tests involve calling the Nonce function.
// The expected result is a successful response from the Nonce function and a matching value with the expected nonce.
// If any errors occur during the tests, the function will fail and display an error message.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress *felt.Felt
		Block           BlockID
		ExpectedNonce   *felt.Felt
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x0207acc15dc241e7d167e67e30e769719a727d3e0fa47f9e187707289885dfde"),
				Block:           WithBlockTag("latest"),
				ExpectedNonce:   utils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		"devnet": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				Block:           WithBlockTag("latest"),
				ExpectedNonce:   utils.TestHexToFelt(t, "0x0"),
			},
		},
		"testnet": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				Block:           WithBlockNumber(69399),
				ExpectedNonce:   utils.TestHexToFelt(t, "0x1"),
			},
		},
		"mainnet": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x00bE9AeF00Ec751Ba252A595A473315FBB8DA629850e13b8dB83d0fACC44E4f2"),
				Block:           WithBlockNumber(644060),
				ExpectedNonce:   utils.TestHexToFelt(t, "0x2"),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		nonce, err := testConfig.provider.Nonce(context.Background(), test.Block, test.ContractAddress)
		if err != nil {
			t.Fatal(err)
		}
		require.NoError(err)
		require.NotNil(nonce, "should return a nonce")
		require.Equal(test.ExpectedNonce, nonce)
	}
}

// TestEstimateMessageFee is a test function to test the EstimateMessageFee function.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestEstimateMessageFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		MsgFromL1
		BlockID
		ExpectedFeeEst FeeEstimate
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				MsgFromL1: MsgFromL1{FromAddress: "0x0", ToAddress: &felt.Zero, Selector: &felt.Zero, Payload: []*felt.Felt{&felt.Zero}},
				BlockID:   BlockID{Tag: "latest"},
				ExpectedFeeEst: FeeEstimate{
					GasConsumed: new(felt.Felt).SetUint64(1),
					GasPrice:    new(felt.Felt).SetUint64(2),
					OverallFee:  new(felt.Felt).SetUint64(3),
				},
			},
		},
		"devnet": {
			{
				MsgFromL1: MsgFromL1{
					FromAddress: "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645",
					ToAddress:   utils.TestHexToFelt(t, "0x05cd48fccbfd8aa2773fe22c217e808319ffcc1c5a6a463f7d8fa2da48218196"),
					Selector:    utils.TestHexToFelt(t, "0xe2bbb158"),
					Payload:     []*felt.Felt{&felt.Zero},
				},
				BlockID: WithBlockTag("latest"),
				ExpectedFeeEst: FeeEstimate{
					GasConsumed: new(felt.Felt).SetUint64(1),
					GasPrice:    new(felt.Felt).SetUint64(2),
					OverallFee:  new(felt.Felt).SetUint64(3),
				},
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		value, err := testConfig.provider.EstimateMessageFee(context.Background(), test.MsgFromL1, test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, *value, test.ExpectedFeeEst)
	}
}

func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	testBlockNumber := uint64(15643)
	type testSetType struct {
		txs           []BroadcastTxn
		simFlags      []SimulationFlag
		blockID       BlockID
		expectedResp  []FeeEstimate
		expectedError error
	}
	testSet := map[string][]testSetType{
		"mainnet": {
			{
				txs: []BroadcastTxn{
					InvokeTxnV0{
						Type:    TransactionType_Invoke,
						Version: TransactionV0,
						MaxFee:  utils.TestHexToFelt(t, "0x95e566845d000"),
						FunctionCall: FunctionCall{
							ContractAddress:    utils.TestHexToFelt(t, "0x45e92c365ba0908382bc346159f896e528214470c60ae2cd4038a0fff747b1e"),
							EntryPointSelector: utils.TestHexToFelt(t, "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad"),
							Calldata: utils.TestHexArrToFelt(t, []string{
								"0x1",
								"0x4a3621276a83251b557a8140e915599ae8e7b6207b067ea701635c0d509801e",
								"0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
								"0x0",
								"0x3",
								"0x3",
								"0x697066733a2f2f516d57554c7a475135556a52616953514776717765347931",
								"0x4731796f4757324e6a5a76564e77776a66514577756a",
								"0x0",
								"0x2"}),
						},
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x63e4618ca2e323a45b9f860f12a4f5c4984648f1d110aa393e79d596d82abcc"),
							utils.TestHexToFelt(t, "0x2844257b088ad4f49e2fe3df1ea6a8530aa2d21d8990112b7e88c4bd0ce9d50"),
						},
					},
				},
				simFlags:      []SimulationFlag{},
				blockID:       BlockID{Number: &testBlockNumber},
				expectedError: nil,
				expectedResp: []FeeEstimate{
					{
						GasConsumed:     utils.TestHexToFelt(t, "0x3074"),
						GasPrice:        utils.TestHexToFelt(t, "0x350da9915"),
						DataGasConsumed: &felt.Zero,
						DataGasPrice:    &felt.Zero,
						OverallFee:      utils.TestHexToFelt(t, "0xa0a99fc14d84"),
						FeeUnit:         UnitWei,
					},
				},
			},
			{

				txs: []BroadcastTxn{
					DeployAccountTxn{

						Type:    TransactionType_DeployAccount,
						Version: TransactionV1,
						MaxFee:  utils.TestHexToFelt(t, "0xdec823b1380c"),
						Nonce:   utils.TestHexToFelt(t, "0x0"),
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x41dbc4b41f6506502a09eb7aea85759de02e91f49d0565776125946e54a2ec6"),
							utils.TestHexToFelt(t, "0x85dcf2bc8e3543071a6657947cc9c157a9f6ad7844a686a975b588199634a9"),
						},
						ContractAddressSalt: utils.TestHexToFelt(t, "0x74ddc51af144d1bd805eb4184d07453d7c4388660270a7851fec387e654a50e"),
						ClassHash:           utils.TestHexToFelt(t, "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918"),
						ConstructorCalldata: utils.TestHexArrToFelt(t, []string{
							"0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2",
							"0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463",
							"0x2",
							"0x74ddc51af144d1bd805eb4184d07453d7c4388660270a7851fec387e654a50e",
							"0x0",
						}),
					},
				},
				simFlags:      []SimulationFlag{},
				blockID:       BlockID{Hash: utils.TestHexToFelt(t, "0x1b0df1bafcb826b1fc053495aef5cdc24d0345cbfa1259b15939d01b89dc6d9")},
				expectedError: nil,
				expectedResp: []FeeEstimate{
					{
						GasConsumed:     utils.TestHexToFelt(t, "0x1154"),
						GasPrice:        utils.TestHexToFelt(t, "0x378f962c4"),
						DataGasConsumed: &felt.Zero,
						DataGasPrice:    &felt.Zero,
						OverallFee:      utils.TestHexToFelt(t, "0x3c2c41636c50"),
						FeeUnit:         UnitWei,
					},
				},
			},
		},
		"devnet": {{

			txs: []BroadcastTxn{
				DeployAccountTxn{

					Type:    TransactionType_DeployAccount,
					Version: TransactionV1,
					MaxFee:  utils.TestHexToFelt(t, "0x16345785d8a0000"),
					Nonce:   utils.TestHexToFelt(t, "0x0"),
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x23db7ec0ab875d20817d1c58e4cea132ef00eed3e529eb86e2dd61bda722575"),
						utils.TestHexToFelt(t, "0x5074825698cb654b0c7d6dba2254d3292eb5536e7dba90b39d6539d2ccc0b86"),
					},
					ContractAddressSalt: utils.TestHexToFelt(t, "0x3db4f4fad96e5444bdc1d0286f42948763af1b7bdb7873c08d28cb5129d4aac"),
					ClassHash:           utils.TestHexToFelt(t, "0x1a736d6ed154502257f02b1ccdf4d9d1089f80811cd6acad48e6b6a9d1f2003"),
					ConstructorCalldata: utils.TestHexArrToFelt(t, []string{
						"0x6e769c37229fe0ea06364d4f91e8b60b17c35f48546cd00a3ede0356eda26bc",
						"0x0",
					}),
				},
			},
			simFlags:      []SimulationFlag{},
			blockID:       WithBlockNumber(16),
			expectedError: nil,
			expectedResp: []FeeEstimate{
				{
					GasConsumed:     utils.TestHexToFelt(t, "0xd06"),
					GasPrice:        utils.TestHexToFelt(t, "0x1"),
					DataGasConsumed: &felt.Zero,
					DataGasPrice:    &felt.Zero,
					OverallFee:      utils.TestHexToFelt(t, "0xd06"),
					FeeUnit:         UnitWei,
				},
			},
		}},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.EstimateFee(context.Background(), test.txs, test.simFlags, test.blockID)
		t.Log("resp: ", resp)
		require.Equal(t, test.expectedError, err)
		require.Equal(t, test.expectedResp, resp)
	}
}
