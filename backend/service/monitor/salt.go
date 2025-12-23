package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"math/rand"
	"meme3/global"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// numberToBytes32 将整数转换为 bytes32
func numberToBytes32(num int64) []byte {
	b := make([]byte, 32)
	bigNum := big.NewInt(num)
	bytes := bigNum.Bytes()
	copy(b[32-len(bytes):], bytes)
	return b
}

// computeCreate2Address 计算 CREATE2 地址
func computeCreate2Address(factory common.Address, initCode []byte, salt int64) common.Address {
	saltBytes := numberToBytes32(salt)
	initCodeHash := crypto.Keccak256(initCode)

	packed := []byte{0xff}
	packed = append(packed, factory.Bytes()...)
	packed = append(packed, saltBytes...)
	packed = append(packed, initCodeHash...)

	hash := crypto.Keccak256(packed)
	return common.BytesToAddress(hash[12:])
}

// findValidSalt 查找满足条件的 salt
func findValidSalt(factory common.Address, initCode []byte, suffix string, maxAttempts int) (int64, common.Address) {
	rand.Seed(time.Now().UnixNano())
	startSalt := rand.Int63n(1000000)
	startTime := time.Now()

	for i := 0; i < maxAttempts; i++ {
		salt := startSalt + int64(i)
		addr := computeCreate2Address(factory, initCode, salt)

		if strings.HasSuffix(strings.ToLower(addr.Hex()), strings.ToLower(suffix)) {
			duration := time.Since(startTime).Seconds()
			fmt.Printf("✅ 找到有效 Salt!\nSalt: %d\n地址: %s\n尝试次数: %d\n耗时: %.2f 秒\n", salt, addr.Hex(), i+1, duration)
			return salt, addr
		}

		if (i+1)%10000 == 0 {
			fmt.Printf("已尝试 %d 次...\n", i+1)
		}
	}

	fmt.Println("❌ 未找到有效 salt")
	return -1, common.Address{}
}

// getDeploymentParams 调用合约获取 factory 和 initCode
func getDeploymentParams(client *ethclient.Client, contractAddress common.Address, tokenId *big.Int, name, symbol string, privateKeyHex string) (common.Address, []byte, error) {
	// 加载合约 ABI
	const contractABI = `[{"inputs":[{"internalType":"address","name":"internalPool","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"}],"name":"getDeploymentParams","outputs":[{"internalType":"address","name":"factory","type":"address"},{"internalType":"bytes32","name":"bytecodeHash","type":"bytes32"},{"internalType":"bytes","name":"initCode","type":"bytes"}],"stateMutability":"view","type":"function"}]`

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return common.Address{}, nil, err
	}

	// 构造调用数据
	data, err := parsedABI.Pack("getDeploymentParams", contractAddress, tokenId, name, symbol)
	if err != nil {
		return common.Address{}, nil, err
	}

	// 调用合约
	// // toAddr := common.HexToAddress(c_contract_token_manage_dev)
	// toAddr := common.HexToAddress(c_contract_token_manage_main)
	// if global.Config.Testnet {
	// 	if global.Config.Debug {
	// 		toAddr = common.HexToAddress(c_contract_token_manage_dev)
	// 	} else {
	// 		toAddr = common.HexToAddress(c_contract_token_manage_test)
	// 	}
	// }
	toAddr := common.HexToAddress(global.Config.RPC.TokenManageAddress)
	callMsg := ethereum.CallMsg{
		To:   &toAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return common.Address{}, nil, err
	}

	// 解析返回值
	outputs, err := parsedABI.Unpack("getDeploymentParams", result)
	if err != nil {
		return common.Address{}, nil, err
	}

	factory := outputs[0].(common.Address)
	initCode := outputs[2].([]byte)

	return factory, initCode, nil
}

func GetSalt(tokenId, tokenName, tokenSymbol string) (int64, string) {
	// 传入参数
	contractAddress := common.HexToAddress(global.Config.ContractAddress)
	tid, err := decimal.NewFromString(tokenId)
	if err != nil {
		slog.Error("get salt failed", "error", err.Error())
		return -1, ""
	}
	privateKey := "0x0000000000000000000000000000000000000000"

	// 调用合约获取 factory 和 initCode
	factory, initCode, err := getDeploymentParams(v_web_client, contractAddress, tid.BigInt(), tokenName, tokenSymbol, privateKey)
	if err != nil {
		slog.Error("get salt failed", "error", err.Error())
		return -1, ""
	}

	// 查找满足条件的 salt
	salt, addr := findValidSalt(factory, initCode, "8888", 10000000)
	if salt != -1 {
		slog.Info("generate salt successed", "salt", salt, "addr", addr)
	}
	return salt, addr.String()
}
