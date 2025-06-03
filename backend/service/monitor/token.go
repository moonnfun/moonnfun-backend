package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/server/web"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	c_created_topic         = "0x1c56be39f0a7cd9b5f7d349a32369c4514e78870028ea4b721570f9a0ea3b127"
	c_created_topic_testnet = "0x1c56be39f0a7cd9b5f7d349a32369c4514e78870028ea4b721570f9a0ea3b127"
)

var TokenCreatedNotify func(tokenCreated any)

/*
创建代币事件：
CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)
*/
type TokenCreated struct {
	TokenId     decimal.Decimal `json:"tokenId"`
	Token       common.Address  `json:"token"`
	Creator     common.Address  `json:"creator"`
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	TotalSupply decimal.Decimal `json:"totalSupply"`
}

// type TokenCreated struct {
// 	TokenId     *big.Int
// 	Token       common.Address
// 	Name        string
// 	Symbol      string
// 	TotalSupply *big.Int
// }

// func WatchTokenCreated(handlers ...StreamHandler) {
// 	vaddress := C_factory_address_main
// 	tokenCreatedTopic := common.HexToHash(c_created_topic)
// 	if global.Config.Testnet {
// 		vaddress = C_factory_address_test
// 		tokenCreatedTopic = common.HexToHash(c_created_topic_testnet)
// 	}

// 	contractAddress := common.HexToAddress(vaddress)
// 	// transferSig := []byte("CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)")
// 	// transferTopic := common.BytesToHash(crypto.Keccak256(transferSig))
// 	query := ethereum.FilterQuery{
// 		Addresses: []common.Address{contractAddress},
// 		Topics:    [][]common.Hash{{tokenCreatedTopic}},
// 	}

// 	logs := make(chan types.Log)
// 	sub, err := v_wss_client.SubscribeFilterLogs(context.Background(), query, logs)
// 	if err != nil {
// 		slog.Error("SubscribeFilterLogs failed", "error", err)
// 		return
// 	}

// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			slog.Error("watch TokenCreated failed", "error", err.Error())
// 		case vLog := <-logs:
// 			slog.Info("watch TokenCreated successed", slog.Any("vLog", vLog))
// 			tokenTrading, err := getTokenCreated(vLog)
// 			if err != nil {
// 				slog.Error("get TokenCreated failed", "vLog", vLog, "error", err.Error())
// 			} else {
// 				for _, handler := range handlers {
// 					if handler != nil {
// 						go handler(tokenTrading)
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func GetFactoryAddress() string {
	address := C_factory_address_main
	if global.Config.Testnet {
		address = C_factory_address_test
	}
	return address
}

func HandleTokenCreatedTx(header *types.Header, tx *types.Transaction) {
	factoryAddress := GetFactoryAddress()
	if tx != nil && tx.To() != nil && tx.To().String() != factoryAddress {
		return
	}

	tokenCreatedTopic := common.HexToHash(c_created_topic)
	if global.Config.Testnet {
		tokenCreatedTopic = common.HexToHash(c_created_topic_testnet)
	}

	contractAddress := common.HexToAddress(factoryAddress)
	// transferSig := []byte("CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)")
	// transferTopic := common.BytesToHash(crypto.Keccak256(transferSig))
	query := ethereum.FilterQuery{
		FromBlock: header.Number,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{tokenCreatedTopic}},
	}
	logs, err := v_wss_client.FilterLogs(context.Background(), query)
	if err != nil {
		slog.Error("get TokenCreated failed", "txHash", tx.Hash().String(), "error", err.Error())
		return
	}

	for _, log := range logs {
		slog.Info("HandleTokenCreatedTx successed", slog.Any("log", log))
		tokenCreated, err := parseTokenCreated(log)
		if err != nil {
			slog.Error("get TokenCreated failed", "log", log, "error", err.Error())
		} else {
			if TokenCreatedNotify != nil {
				go TokenCreatedNotify(tokenCreated)
			}
		}
	}
}

func parseTokenCreated(l types.Log) (*TokenCreated, error) {
	// contractABI, err := abi.JSON(strings.NewReader(ABI))
	// if err != nil {
	// 	return nil, err
	// }

	// tokenCreated := &TokenCreated{}
	// if err := contractABI.UnpackIntoInterface(&tokenCreated, "CreateToken", l.Data); err != nil {
	// 	return nil, err
	// }
	// supply, _ := decimal.NewFromString(tokenCreated.TotalSupply.String())
	// tokenCreated.TotalSupply = supply.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))).BigInt()
	// slog.Info("Get TokenCreated successed", "tokenCreated", tokenCreated)

	parseArgs, err := ParseLog(l.Topics[0].String(), "0x"+common.Bytes2Hex(l.Data))
	if err != nil {
		slog.Error("get TokenCreated failed", "topic", l.Topics[0].String(), "log.data", common.Bytes2Hex(l.Data), "error", err.Error())
		return nil, err
	}
	tokenCreated := &TokenCreated{}
	if err := json.Unmarshal([]byte(parseArgs), tokenCreated); err != nil {
		slog.Error("get TokenCreated failed", "parseArgs", parseArgs, "error", err.Error())
		return nil, err
	}
	slog.Info("unpack ok", slog.Any("tokenCreated", tokenCreated))

	tokenCreated.TotalSupply = tokenCreated.TotalSupply.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
	slog.Info("Get TokenCreated successed", "tokenCreated", tokenCreated)

	return tokenCreated, nil
}

func GetTokenCreated(tokenID string) ([]*TokenCreated, error) {
	v_rpc_url = C_Web_RpcUrl_main
	explorerUrl := C_Web_Explorer_main
	if global.Config.Testnet {
		v_rpc_url = C_Web_RpcUrl_test
		explorerUrl = C_Web_Explorer_test
	}

	fetchUrl := fmt.Sprintf("%s/api/v2/addresses/%s/logs", explorerUrl, GetFactoryAddress())
	fetchBuf, err := web.HttpGet(fetchUrl)
	if err != nil {
		return nil, err
	}

	ret := make([]*TokenCreated, 0)
	fetchResult := gjson.Get(string(fetchBuf), "items").Array()
	for _, result := range fetchResult {
		topics := result.Get("topics").Array()
		if len(topics) > 0 && topics[0].String() == c_created_topic {
			parseArgs, err := ParseLog(topics[0].String(), result.Get("data").String())
			if err != nil {
				slog.Error("parse log failed", "log", result.Get("data").String(), "error", err.Error())
				continue
			}

			tokenCreated := &TokenCreated{}
			if err := json.Unmarshal([]byte(parseArgs), tokenCreated); err != nil {
				slog.Error("get TokenCreated failed", "parseArgs", parseArgs, "error", err.Error())
				continue
			}
			slog.Info("unpack ok", slog.Any("tokenCreated", tokenCreated))

			tokenCreated.Token = common.HexToAddress(result.Get("address.hash").String())
			tokenCreated.TotalSupply = TrimDecimals(tokenCreated.TotalSupply)
			slog.Info("Get TokenCreated successed", "tokenCreated", tokenCreated)
			ret = append(ret, tokenCreated)
			if tokenID != "" && tokenCreated.TokenId.String() == tokenID {
				break
			}
		}
	}
	return ret, nil
}

func getTokenCreated(topic, data, tokenAddress string) *TokenCreated {
	parseArgs, err := ParseLog(topic, data)
	if err != nil {
		slog.Error("parse log failed", "topic", topic, "log", data, "error", err.Error())
		return nil
	}

	tokenCreated := &TokenCreated{}
	if err := json.Unmarshal([]byte(parseArgs), tokenCreated); err != nil {
		slog.Error("get TokenCreated failed", "parseArgs", parseArgs, "error", err.Error())
		return nil
	}
	// slog.Info("unpack ok", slog.Any("tokenCreated", tokenCreated))

	tokenCreated.TotalSupply = TrimDecimals(tokenCreated.TotalSupply)
	global.Debug("Get TokenCreated successed", "tokenCreated", tokenCreated)
	return tokenCreated
}
