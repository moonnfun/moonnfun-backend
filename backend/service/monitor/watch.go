package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"math/big"
	"meme3/global"
	"meme3/server/web"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	C_decimal      = 18
	C_supply       = 1000000000
	C_zero_address = "0x0000000000000000000000000000000000000000"

	C_Pay_Amount_test = 10
	C_Pay_Amount_main = 100
	C_Price_Amount    = "0.001"
)

var (
	QuotePrice = decimal.Zero

	v_rpc_url = ""

	v_logs_tx = true

	v_wss_client   *ethclient.Client
	v_web_client   *ethclient.Client
	v_watch_client *ethclient.Client

	v_watch_first_height int64 = 0

	v_user_contract_list = ""

	C_Web_Script_Url = "http://localhost:3157"
)

type TxData struct {
	To        string
	From      string
	Hash      string
	Time      time.Time
	BlockHash string
}

type ReceiptData struct {
	TxHash      string
	RawBuf      []byte
	BlockNumber *big.Int
}

type BlockData struct {
	Hash   string
	Height int64
	Time   time.Time
	TxList []string
}

func InitWeb3Client() error {
	v_rpc_url = global.Config.RPC.Http
	wssRpcUrl := global.Config.RPC.Websocket
	// v_payment_amount = C_Pay_Amount_main
	// v_payment_address = C_Payment_address_main
	// if global.Config.Testnet {
	// 	v_rpc_url = C_Web_RpcUrl_test
	// 	wssRpcUrl = C_Wss_RpcUrl_test
	// 	if global.Config.Debug {
	// 		v_rpc_url = C_Web_RpcUrl_dev
	// 		wssRpcUrl = C_Wss_RpcUrl_dev
	// 	}
	// }

	if !global.Config.RPC.WatchByScan {
		client, err := ethclient.Dial(wssRpcUrl)
		if err != nil {
			slog.Error("Dial websocket connect failed", "error", err)
			return err
		}
		v_watch_client = client
		v_wss_client = client
	}

	wclient, err := ethclient.Dial(v_rpc_url)
	if err != nil {
		slog.Error("Dial web connect failed", "error", err)
		return err
	}
	v_web_client = wclient

	v_user_contract_list = strings.ToUpper(global.Config.RPC.UserContractList)

	return nil
}

func WatchContractEvents(handlers ...func(*big.Int, *TxData, []byte)) {
	if err := InitWeb3Client(); err != nil {
		return
	}
	v_logs_tx = global.Config.RPC.GetLogByFetch

	var latestNumber int64
	if global.GetLatestNumber != nil {
		latestNumber = global.GetLatestNumber()
	}
	// latestNumber = 178030950
	global.Info("get latest block number successed======", "number", latestNumber)

	buyTopic := common.HexToHash(c_buyTopic)
	sellTopic := common.HexToHash(c_sellTopic)
	launchTopic := common.HexToHash(c_launchTopic)
	swapBuyTopic := common.HexToHash(c_swapBuyTopic)
	swapSellTopic := common.HexToHash(c_swapSellTopic)
	tokenCreatedTopic := common.HexToHash(c_created_topic)
	contractAddress := common.HexToAddress(global.Config.ContractAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(latestNumber),
		ToBlock:   nil,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{tokenCreatedTopic, buyTopic, sellTopic, launchTopic, swapBuyTopic, swapSellTopic}},
	}
	doWatchContractEvents(global.Config.ContractAddress, query, HandleTokenCreatedTx, HandleTokenTradingTx, HandleTokenSwapTx)
}

func doWatchContractEvents(contract string, query ethereum.FilterQuery, handlers ...func(*big.Int, *TxData, []byte)) {
	defer func() {
		go doWatchContractEvents(contract, query, handlers...)
	}()
	logCh := make(chan types.Log, 1024)
	sub, err := v_wss_client.SubscribeFilterLogs(context.Background(), query, logCh)
	if err != nil {
		slog.Error("subscribe filter logs failed", "error", err.Error())
		return
	}
	slog.Info("start watch log events", "contract", contract)

	for {
		select {
		case err := <-sub.Err():
			slog.Error("SubscribeFilterLogs failed", "error", err.Error())
			return
		case clog := <-logCh:
			clogs := []types.Log{clog}
			go HandleLogEvents(clogs, handlers...)
		}
	}
}

func FetchBlockByRetryEx(height int64, retryTimes int, handlers ...func(*big.Int, *TxData, []byte)) {
	i := 1
	for {
		select {
		case <-time.After(time.Duration(3) * time.Second):
			if i >= retryTimes {
				return
			} else {
				// block, err := v_wss_client.BlockByNumber(context.Background(), big.NewInt(height))
				// if err != nil {
				// 	slog.Error("Get Blok failed", "height", height, "times", i, "error", err)
				// 	i += 1
				// 	break
				// }
				// go handleHeader(block, handlers...)
				// return
				handleHeaderEx(height, true, handlers...)
			}
		}
	}
}

func HandleLatestBlocksEx(latestNumber, numberStart, numberEnd int64, handlers ...func(*big.Int, *TxData, []byte)) error {
	if numberStart > 0 && numberEnd > 0 {
		if numberStart < numberEnd {
			// for i := numberStart; i < numberEnd; i++ {
			// 	handleHeaderEx(i, false, handlers...)
			// 	slog.Info("handle the latest block successed", "number", i)
			// }
			SyncLogEvents(numberStart, numberEnd, handlers...)
		}
	} else if latestNumber > 0 {
		if latestNumber < v_watch_first_height {
			global.Config.Sync = true
			// for i := latestNumber; i < v_watch_first_height; i++ {
			// 	handleHeaderEx(i, false, handlers...)
			// 	slog.Info("handle the latest block successed", "number", i)
			// }
			SyncLogEvents(latestNumber, v_watch_first_height, handlers...)
			global.Config.Sync = false
		}
	}
	return nil
}

func handleHeaderEx(height int64, bRetry bool, handlers ...func(*big.Int, *TxData, []byte)) {
	time.Sleep(time.Duration(global.Config.RPC.WaitBlockTime) * time.Millisecond)
	blockData, err := GetBlock(height)
	if err != nil {
		slog.Error("handleHeader failed at BlockReceipts", "error", err.Error())
		if !bRetry {
			go FetchBlockByRetryEx(height, 3, handlers...)
		}
		return
	}

	for _, tx := range blockData.TxList {
		time.Sleep(time.Duration(global.Config.RPC.WaitBlockTime) * time.Millisecond)
		txData, err := GetTx(tx)
		if err != nil {
			slog.Error("handleHeader failed at GetTx", "error", err.Error())
			continue
		}
		txData.Time = blockData.Time
		txData.BlockHash = blockData.Hash
		// global.Debug("get tx ok", "receipt", receipt, "txData", txData)

		for _, handler := range handlers {
			// time.Sleep(time.Duration(100) * time.Millisecond)
			handler(big.NewInt(height), txData, nil)
		}
	}
}

func GetReceipts(height int64) ([]*ReceiptData, error) {
	ret := make([]*ReceiptData, 0)

	postJson := `{"method":"eth_getBlockReceipts","params":["height"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.ReplaceAll(postJson, "height", hexutil.Uint64(height).String())
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		slog.Error("eth_getBlockReceipts failed", "height", height)
		return ret, err
	}

	receipts := gjson.Get(string(fetchBuf), "result").Array()
	for _, receipt := range receipts {
		ret = append(ret, &ReceiptData{
			TxHash:      receipt.Get("transactionHash").String(),
			BlockNumber: big.NewInt(height),
			RawBuf:      []byte(receipt.Raw),
		})
	}
	return ret, nil
}

func GetTx(txhash string) (*TxData, error) {
	postJson := `{"method":"eth_getTransactionByHash","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "txhash", txhash, -1)
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		return nil, err
	}

	txData := &TxData{
		To:   gjson.Get(string(fetchBuf), "result.to").String(),
		From: gjson.Get(string(fetchBuf), "result.from").String(),
		Hash: txhash,
		Time: time.Now(),
	}
	return txData, nil
}

func GetBlock(height int64) (*BlockData, error) {
	postJson := `{"method":"eth_getBlockByNumber","params":["height",false],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "height", hexutil.EncodeBig(big.NewInt(height)), -1)
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		return nil, err
	}
	// slog.Debug("get block ok", "rpc", v_rpc_url, "block", string(fetchBuf), "postJson", postJson)

	timestamp, err := hexutil.DecodeBig(gjson.Get(string(fetchBuf), "result.timestamp").String())
	if err != nil {
		return nil, err
	}
	// slog.Debug("get timestamp ok", "timestamp", timestamp)

	blockData := &BlockData{
		Height: height,
		Hash:   gjson.Get(string(fetchBuf), "result.hash").String(),
		Time:   time.Unix(timestamp.Int64(), 0),
	}

	txList := gjson.Get(string(fetchBuf), "result.transactions").Array()
	for _, tx := range txList {
		blockData.TxList = append(blockData.TxList, tx.String())
	}
	return blockData, nil
}

func SyncLogEvents(start, end int64, handlers ...func(*big.Int, *TxData, []byte)) error {
	slog.Info("start sync logs", "start", start, "end", end)
	workingStart := start
	workingEnd := start
	filterLogsMax := 2000
	if global.Config.RPC.FilterLogsMax > 0 {
		filterLogsMax = global.Config.RPC.FilterLogsMax
	}

	contractAddress := common.HexToAddress(global.Config.ContractAddress)
	tokenCreatedTopic := common.HexToHash(c_created_topic)
	buyTopic := common.HexToHash(c_buyTopic)
	sellTopic := common.HexToHash(c_sellTopic)
	launchTopic := common.HexToHash(c_launchTopic)
	swapBuyTopic := common.HexToHash(c_swapBuyTopic)
	swapSellTopic := common.HexToHash(c_swapSellTopic)

	for workingEnd < end {
		time.Sleep(time.Duration(3) * time.Second)
		workingStart = workingEnd
		workingEnd += int64(filterLogsMax)
		if workingEnd >= end {
			workingEnd = end
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(workingStart),
			ToBlock:   big.NewInt(workingEnd),
			Addresses: []common.Address{contractAddress},
			Topics:    [][]common.Hash{{tokenCreatedTopic, buyTopic, sellTopic, launchTopic, swapBuyTopic, swapSellTopic}},
		}
		filterLogs, err := v_web_client.FilterLogs(context.Background(), query)
		if err != nil {
			slog.Error("get TokenTrading failed", "start", start, "end", end, "error", err.Error())
			continue
		}
		slog.Info("filter logs successed", "start", start, "end", end, "logs", len(filterLogs))

		HandleLogEvents(filterLogs[:], handlers...)
	}
	return nil
}

func HandleLogEvents(filterLogs []types.Log, handlers ...func(*big.Int, *TxData, []byte)) error {
	for _, flog := range filterLogs {
		time.Sleep(time.Duration(global.Config.RPC.WaitBlockTime) * time.Millisecond)
		height := int64(flog.BlockNumber)
		blockData, err := GetBlock(height)
		if err != nil {
			slog.Error("handleHeader failed at BlockReceipts", "error", err.Error())
			continue
		}
		slog.Info("fetch block successed", "height", height, "topics", flog.Topics)

		time.Sleep(time.Duration(global.Config.RPC.WaitBlockTime) * time.Millisecond)
		txData, err := GetTx(flog.TxHash.String())
		if err != nil {
			slog.Error("handleHeader failed at GetTx", "error", err.Error())
			continue
		}
		txData.Time = blockData.Time
		// global.Debug("get tx ok", "blockData", blockData, "txData", txData)

		handleLogs := []types.Log{flog}
		var inputLogs = map[string]map[string][]types.Log{
			"result": {
				"logs": handleLogs,
			},
		}
		logsBuf, _ := json.Marshal(inputLogs)

		for _, handler := range handlers {
			handler(big.NewInt(height), txData, logsBuf)
		}
	}
	return nil
}

func GetLogs(client *ethclient.Client, tx *TxData, query ethereum.FilterQuery, receiptBuf []byte) ([]types.Log, error) {
	time.Sleep(time.Duration(global.Config.RPC.WaitBlockTime) * time.Millisecond)
	if v_logs_tx {
		return GetLogsByTx(client, tx, query, receiptBuf)
	} else {
		return GetLogsByFilter(client, query)
	}
}

func GetLogsByFilter(client *ethclient.Client, query ethereum.FilterQuery) ([]types.Log, error) {
	return client.FilterLogs(context.Background(), query)
}

func GetLogsByHash(txhash string) []byte {
	postJson := `{"method":"eth_getTransactionReceipt","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "txhash", txhash, -1)
	response, _ := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	return response
}

func GetLogsByTx(client *ethclient.Client, tx *TxData, query ethereum.FilterQuery, receiptBuf []byte) ([]types.Log, error) {
	if receiptBuf == nil {
		postJson := `{"method":"eth_getTransactionReceipt","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
		postJson = strings.Replace(postJson, "txhash", tx.Hash, -1)
		response, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
		if err != nil {
			return nil, err
		}
		receiptBuf = response[:]
		// slog.Info("get tx logs successed", "tx", tx.Hash().String(), "result", string(fetchBuf))
	}

	// if tx.From == "0x36c6dfba13ca516bf03e22055f8578f022f4542f" {
	// 	slog.Debug("fetch log ok", "logs", string(receiptBuf))
	// }

	ret := make([]types.Log, 0)
	rlogs := gjson.Get(string(receiptBuf), "result.logs").Array()
	matchCount := len(query.Topics)
	for _, rlog := range rlogs {
		topics := rlog.Get("topics").Array()
		for _, hashList := range query.Topics {
			for _, topicHash := range hashList {
				for _, topic := range topics {
					if topic.String() == topicHash.String() {
						matchCount -= 1
					}
				}
			}
		}
		if matchCount == 0 {
			tlog := types.Log{}
			if err := json.Unmarshal([]byte(rlog.Raw), &tlog); err != nil {
				return ret, nil
			}
			ret = append(ret, tlog)
		}
	}
	return ret, nil
}

func GetTxData(tx *types.Transaction) *TxData {
	txData := &TxData{
		Hash: tx.Hash().String(),
		Time: tx.Time(),
	}
	if tx.To() != nil {
		txData.To = tx.To().String()
	} else {
		// slog.Error("GetTxData with invalid tx", "tx", tx)
	}
	return txData
}

func IsContractAddress(address string) bool {
	postJson := `{"method":"eth_getCode","params":["txhash", "latest"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "txhash", common.HexToAddress(address).String(), -1)
	response, _ := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	return gjson.Get(string(response), "result").String() != "0x"
}
