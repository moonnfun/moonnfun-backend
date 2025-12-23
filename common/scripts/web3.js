let logFunc = console.log;
console.log = function(){};
// console.debug = function(){};

import * as evm from "./evm.js";

async function main() {
    if (process.argv[2] === "GetPrice") {
        let rpcUrl = process.argv[3];
        await evm.initChain(rpcUrl);
        logFunc(JSON.stringify(await evm.GetPrice(rpcUrl)));
    } else if (process.argv[2] === "GetTokenPrice") {
        let rpcUrl = process.argv[3];
        let tokenAddress = process.argv[4];
        let routerAddress = process.argv[5];
        let buyAmount = process.argv[6];
        await evm.initChain(rpcUrl);
        logFunc(JSON.stringify(await evm.GetTokenPrice(tokenAddress, routerAddress, buyAmount)));
    } else if (process.argv[2] === "GetTokenPriceEx") {
        let rpcUrl = process.argv[3];
        let tokenAddress = process.argv[4];
        let routerAddress = process.argv[5];
        await evm.initChain(rpcUrl);
        logFunc(JSON.stringify(await evm.GetTokenPriceEx(tokenAddress, routerAddress)));
    } else if (process.argv[2] === "Token0") {
        let rpcUrl = process.argv[3];
        let pairAddress = process.argv[4];
        await evm.initChain(rpcUrl);
        logFunc(await evm.Token0(pairAddress));
    } else if (process.argv[2] === "GetReverse") {
        let rpcUrl = process.argv[3];
        let tokenAddress = process.argv[4];
        let routerAddress = process.argv[5];
        let strLaunch = process.argv[6];
        await evm.initChain(rpcUrl);
        logFunc(JSON.stringify(await evm.GetReverse(tokenAddress, routerAddress, strLaunch)));
    } else if (process.argv[2] === "ParseLog") {
        let logTopic = process.argv[3];
        let logData = process.argv[4];
        logFunc(await evm.ParseLog(logTopic, logData));
    } else if (process.argv[2] === "ParseSwapLog") {
        let logTopic = process.argv[3];
        let logData = process.argv[4];
        logFunc(await evm.ParseSwapLog(logTopic, logData));
    } else if (process.argv[2] === "GetPair") {
        let rpcUrl = process.argv[3];
        let tokenAddress = process.argv[4];
        let routerAddress = process.argv[5];
        await evm.initChain(rpcUrl);
        logFunc(await evm.GetPair(tokenAddress, routerAddress));
    } else if (process.argv[2] === "GetTokenBalance") {
        let rpcUrl = process.argv[3];
        let tokenAddress = process.argv[4];
        let walletAddress = process.argv[5];
        await evm.initChain(rpcUrl);
        logFunc(await evm.GetTokenBalance(tokenAddress, walletAddress));
    } else if (process.argv[2] === "IsListToken") {
        let rpcUrl = process.argv[3];
        let id = process.argv[4];
        let tokenAddress = process.argv[5];
        let walletAddress = process.argv[6];
        let contractAddress = process.argv[7];
        await evm.initChain(rpcUrl);
        logFunc(await evm.IsListToken(id, tokenAddress, walletAddress, contractAddress));
    } else {
        logFunc("invalid operation")
    }
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });