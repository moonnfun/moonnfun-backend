package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/server/web"
	"meme3/service/model"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func GetQuotePrice() decimal.Decimal {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "GetPrice")
	args = append(args, global.Config.RPC.Http) // QuotePrice统一用主网数据
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	// global.Debug("GetTokenPrice result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetQuotePrice failed", "output", string(strOutput), "error", err)
		return QuotePrice
	}

	ret := model.TokenPrice{}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	// price, err := decimal.NewFromString(gjson.Get(strResult, "priceUsd").String())
	if err := json.Unmarshal([]byte(strResult), &ret); err != nil {
		slog.Error("GetQuotePrice failed", "output", string(strOutput), "error", err)
		return QuotePrice
	}
	// global.Logf("%s: %v", tokenAddress, string(strOutput))
	QuotePrice = ret.PriceUsd
	return ret.PriceUsd
}

func GetTokenPrice(method, tokenAddress, routerAddress string) model.TokenPrice {
	if method == "" {
		method = "GetTokenPrice"
	}

	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, method)
	args = append(args, v_rpc_url)
	args = append(args, tokenAddress)
	args = append(args, routerAddress)
	args = append(args, C_Price_Amount)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	global.Debug("GetTokenPrice result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetTokenPrice failed", "output", string(strOutput), "error", err)
		return model.TokenPrice{}
	}

	ret := model.TokenPrice{}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	// price, err := decimal.NewFromString(gjson.Get(strResult, "priceUsd").String())
	if err := json.Unmarshal([]byte(strResult), &ret); err != nil {
		slog.Error("GetTokenPrice failed", "output", string(strOutput), "error", err)
		return model.TokenPrice{}
	}
	ret.CurrentTime = time.Now().UnixMilli()
	global.Debug("token price ok", tokenAddress, strResult)
	return ret
}

func ParseLog(topic, data string) (string, error) {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "ParseLog")
	args = append(args, topic)
	args = append(args, data)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	// global.Debug("GetTokenPrice result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("ParseLog failed", "output", string(strOutput), "error", err)
		return "", err
	}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	return strResult, err
}

func ParseSwapLog(topic, data string) (string, error) {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "ParseSwapLog")
	args = append(args, topic)
	args = append(args, data)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	// global.Debug("GetTokenPrice result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("ParseLog failed", "output", string(strOutput), "error", err)
		return "", err
	}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	return strResult, err
}

func GetPair(tokenAddress, routerAddress string) (string, error) {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "GetPair")
	args = append(args, v_rpc_url)
	args = append(args, tokenAddress)
	args = append(args, routerAddress)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	// global.Debug("GetTokenPrice result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetPair failed", "output", string(strOutput), "error", err)
		return "", err
	}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	return strResult, err
}

func Launch(tokenAddress, routerAddress, privateKey string) (string, error) {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "Launch")
	args = append(args, v_rpc_url)
	args = append(args, tokenAddress)
	args = append(args, routerAddress)
	args = append(args, privateKey)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	global.Debug("Launch result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("ParseLog failed", "output", string(strOutput), "error", err)
		return "", err
	}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	return strResult, err
}

func GetReverse(tokenAddress, routerAddress string, bLaunch bool) *model.TokenTrading {
	launch := "false"
	if bLaunch {
		launch = "true"
	}
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "GetReverse")
	args = append(args, v_rpc_url)
	args = append(args, tokenAddress)
	args = append(args, routerAddress)
	args = append(args, launch)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	global.Debug("GetReverse result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetReverse failed", "output", string(strOutput), "error", err)
		return nil
	}

	ret := model.TokenTrading{}
	strResult := strings.TrimSuffix(string(strOutput), "\n")
	// price, err := decimal.NewFromString(gjson.Get(strResult, "priceUsd").String())
	if err := json.Unmarshal([]byte(strResult), &ret); err != nil {
		slog.Error("GetReverse failed", "output", string(strOutput), "error", err)
		return nil
	}
	global.Debug("GetReverse successed", tokenAddress, string(strOutput))
	return &ret
}

func GetTokenBalance(tokenAddress, walletAddress string) decimal.Decimal {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "GetTokenBalance")
	args = append(args, v_rpc_url)
	args = append(args, tokenAddress)
	args = append(args, walletAddress)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	global.Debug("GetTokenBalance result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetTokenBalance failed", "output", string(strOutput), "error", err)
		return decimal.NewFromInt(-1)
	}

	ret, err := decimal.NewFromString(strings.TrimSuffix(string(strOutput), "\n"))
	if err != nil {
		slog.Error("GetTokenBalance failed", "output", string(strOutput), "error", err)
		return decimal.NewFromInt(-1)
	}
	global.Debug("GetTokenBalance successed", tokenAddress, string(strOutput))
	return ret
}

func GetToken0(pairAddress string) string {
	dir, _ := os.Getwd()
	args := make([]string, 0)
	cpath := filepath.Join(dir, "scripts", "web3.js")
	args = append(args, cpath)
	args = append(args, "Token0")
	args = append(args, v_rpc_url)
	args = append(args, pairAddress)
	cmd := exec.Command(global.Config.NodePath, args...)

	largs := make([]string, 0)
	largs = append(largs, global.Config.NodePath)
	largs = append(largs, cpath)
	largs = append(largs, args[1:]...)
	// global.Log("before call GetTokenPrice, cmd: ", strings.Join(largs, " "))

	strOutput, err := cmd.Output()
	global.Debug("GetToken0 result:", string(strOutput), "error:", err)
	if err != nil {
		slog.Error("GetToken0 failed", "output", string(strOutput), "error", err)
		return ""
	}

	ret := strings.TrimSuffix(string(strOutput), "\n")
	global.Debug("GetToken0 successed", pairAddress, string(strOutput))
	return ret
}

func IsListToken(id, tokenAddress, walletAddress string) string {
	tokenAddress = common.HexToAddress(tokenAddress).String()
	walletAddress = common.HexToAddress(walletAddress).String()
	fetchUrl := fmt.Sprintf("%s/%s?testnet=%v&id=%s&tokenAddress=%s&walletAddress=%s&contractAddress=%s", C_Web_Script_Url, "IsListToken", global.Config.Testnet, id, tokenAddress, walletAddress, "")
	strOutput, err := web.HttpGet(fetchUrl)
	if err != nil {
		slog.Error("IsListToken failed", "output", string(strOutput), "error", err, "fetchUrl", fetchUrl)
		return err.Error()
	}
	global.Debug("fetch successed", "fetchUrl", fetchUrl, "output", string(strOutput))

	return strings.TrimSuffix(string(strOutput), "\n")
}

func GetTxSender(txhash string) string {
	postJson := `{"method":"eth_getTransactionByHash","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "txhash", txhash, -1)
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		return ""
	}
	return gjson.Get(string(fetchBuf), "result.from").String()
}

func GetTxReceiver(rpcUrl, txhash string) string {
	if v_rpc_url == "" {
		v_rpc_url = rpcUrl
	}
	postJson := `{"method":"eth_getTransactionByHash","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.Replace(postJson, "txhash", txhash, -1)
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		return ""
	}
	return gjson.Get(string(fetchBuf), "result.to").String()
}

func IsTxOK(txhash string) (bool, int64) {
	postJson := `{"method":"eth_getTransactionReceipt","params":["txhash"],"id":1,"jsonrpc":"2.0"}`
	postJson = strings.ReplaceAll(postJson, "txhash", txhash)
	fetchBuf, err := web.HttpPostJson(v_rpc_url, bytes.NewBuffer([]byte(postJson)))
	if err != nil {
		return false, 0
	}
	to := gjson.Get(string(fetchBuf), "result.to").String()
	global.Debug("eth_getTransactionReceipt ok", "txhash", txhash, "result", string(fetchBuf))
	number := hexutil.MustDecodeUint64(gjson.Get(string(fetchBuf), "result.blockNumber").String())
	bTxOK := strings.ToUpper(to) == strings.ToUpper(global.Config.ContractAddress) && gjson.Get(string(fetchBuf), "result.status").String() == "0x1"
	return bTxOK, int64(number)
}

func TrimDecimals(val decimal.Decimal) decimal.Decimal {
	return val.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
}
