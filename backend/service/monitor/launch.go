package monitor

import (
	"encoding/json"
	"log/slog"
	"meme3/service/model"
	"meme3/service/store"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	c_launchTopicx      = "0xe4b8baaa1c927a55ac28b7693040a55b0b48044cf0572d1fc3beb0f418054a20"
	c_launch_privateKey = "0x9f0e5956e79ca0dc773ebd70d4e2840532b21bc3b854b6a6fdf7358287b736d6"
)

type LaunchLiquid struct {
	Token common.Address `json:"token"`
	Pair  common.Address `json:"lpToken"`
}

var TokenLaunchNotify func(tokenLaunchLiquid any) bool

// func HandleLaunchLiquidTxx(height *big.Int, tx *TxData, receiptBuf []byte) {
// 	contractHash := common.HexToAddress(global.Config.ContractAddress)
// 	launchTopic := common.HexToHash(c_launchTopic)
// 	blockHash := common.HexToHash(tx.BlockHash)
// 	query := ethereum.FilterQuery{
// 		BlockHash: &blockHash,
// 		Addresses: []common.Address{contractHash},
// 		Topics:    [][]common.Hash{{launchTopic}},
// 	}

// 	logs := make([]types.Log, 0)
// 	if common.HexToAddress(tx.To).Cmp(common.HexToAddress(global.Config.ContractAddress)) != 0 {
// 		if !IsContractAddress(tx.To) {
// 			return
// 		}
// 		if logsBuf := GetLogsByHash(tx.Hash); len(logsBuf) > 0 {
// 			tlogs, err := GetLogs(v_wss_client, tx, query, logsBuf)
// 			if err != nil || len(tlogs) == 0 {
// 				// slog.Debug("get LaunchLiquid failed", "txHash", tx.Hash, "error", err)
// 				return
// 			} else {
// 				logs = tlogs[:]
// 			}
// 		} else {
// 			return
// 		}
// 	}

// 	if logs != nil && len(logs) == 0 {
// 		tlogs, err := GetLogs(v_wss_client, tx, query, receiptBuf)
// 		if err != nil {
// 			slog.Error("get LaunchLiquid failed", "txHash", tx.Hash, "error", err.Error())
// 			return
// 		}
// 		logs = tlogs[:]
// 	}
// 	// global.DebugForce("get logs successed", "logs", logs, "routerAddress", GetRouter())

// 	for _, log := range logs {
// 		slog.Info("HandleLaunchLiquidTx successed", slog.Any("log", log))
// 		if err := getLaunchLiquid(log, tx, tx.Time); err != nil {
// 			slog.Error("get LaunchLiquid failed", "log", log, "error", err.Error())
// 		}
// 	}
// }

func getLaunchLiquid(log types.Log, tx *TxData, txTime time.Time) error {
	topics := ""
	liquidEvent := &LaunchLiquid{}
	for _, topic := range log.Topics {
		topics += topic.String() + ","
	}
	topics = strings.TrimSuffix(topics, ",")
	if err := doGetLaunchLiquid(topics, "0x"+common.Bytes2Hex(log.Data), liquidEvent); err != nil {
		return err
	}

	token, err := store.DBGet[model.Token](model.C_Token, bson.M{"address": liquidEvent.Token.String()})
	if err != nil {
		return err
	}
	token.Launch = true
	token.Progress = 100
	token.Pair = liquidEvent.Pair.String()
	token.LaunchTime = time.Now().UnixMilli()

	token.DBID = primitive.ObjectID{}
	if err := store.DBSet(model.C_Token, token, bson.M{"address": token.Address}); err != nil {
		return err
	}
	RegisterPair(token.Pair)

	if TokenLaunchNotify != nil {
		TokenLaunchNotify(liquidEvent)
	}

	// // launch dex
	// output, err := Launch(token.Address, global.Config.ContractAddress, c_launch_privateKey)
	// if err != nil {
	// 	return err
	// }
	// if !strings.Contains(output, "successed") {
	// 	return fmt.Errorf("launch token failed, token: %+v, output: %s", token, output)
	// }
	return nil
}

func doGetLaunchLiquid(topic, data string, output any) error {
	parseArgs, err := ParseLog(topic, data)
	if err != nil {
		slog.Error("parse liquid log failed", "topic", topic, "log", data, "error", err.Error())
		return err
	}

	if err := json.Unmarshal([]byte(parseArgs), output); err != nil {
		slog.Error("get LaunchLiquid failed", "parseArgs", parseArgs, "error", err.Error())
		return err
	}
	slog.Info("Get LaunchLiquid successed", "output", output)
	return nil
}
