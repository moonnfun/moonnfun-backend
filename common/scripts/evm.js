import Decimal from 'decimal.js';
import * as ethers from "ethers";
import * as ABI from "./ABI.js";
import * as ABIex from "./ABIex.js";
import { FactoryABI, RouterABI } from './factory.js';

export const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

let provider = {};

Decimal.set({ toExpNeg: -1000, toExpPos: 1000 });

export async function initChain(rpcURL) {
    provider = new ethers.JsonRpcProvider(rpcURL)
}

export async function GetTokenPrice(tokenAddress, routerAddress, buyAmount) {
    try {
        const routerV2 = new ethers.Contract(routerAddress, ABI.ABI, provider);

        const quoteAMount = parseFloat(buyAmount);
        const amountIn = quoteAMount * 10 ** 18;
        // const tradingTokens = [Wrapper_Address, tokenAddress];
        const amountsOut = await routerV2.getAmountOut(ZERO_ADDRESS, BigInt(amountIn), tokenAddress);

        const priceNative = await GetPrice();
        const tokenAmount = new Decimal(amountsOut).dividedBy(10 ** 18);
        const price = new Decimal(amountIn).dividedBy(amountsOut);
        const priceUsd = new Decimal(priceNative.priceUsd).mul(price);
        console.debug(`pancakeV2 getAmountsOut successed, token: ${tokenAddress}, tokenAmount: ${tokenAmount}, quoteAmount: ${quoteAMount}`);
        return {
            price: price,
            priceUsd: priceUsd
        }
    } catch (err) {
        return { 'Error': `${err}` };
    }
}

export async function GetReverse(tokenAddress, routerAddress, launch) {
    try {
        let decimals = 0;
        let launchAmount = 0;
        let tokenLaunchAmount = 0;

        const routerV2 = new ethers.Contract(routerAddress, ABI.ABI, provider);
        const amountsOut = await routerV2.getReserve(tokenAddress);
        console.debug(`get reverse successed, token: ${tokenAddress}, amountsOut: ${amountsOut}`);
        return {
            tokenAmount: new Decimal(amountsOut[0]).sub(tokenLaunchAmount).dividedBy(10 ** 18),
            totalAmount: new Decimal(amountsOut[1]).sub(launchAmount).dividedBy(10 ** 18)
        }
    } catch (err) {
        return { 'Error': `${err}` };
    }
}

export async function GetPrice() {
    try {
        // 1. Binance: SEIUSDT
        const binanceUrl = "https://api.binance.com/api/v3/ticker/price?symbol=BNBUSDT";
        const binanceData = await fetch(binanceUrl).then(res => res.json());
        const binancePrice = parseFloat(binanceData.price);

        // 2. Kraken: SEI/USD
        const krakenUrl = "https://api.kraken.com/0/public/Ticker?pair=SEIUSD";
        const krakenData = await fetch(krakenUrl).then(res => res.json());
        const krakenKey = Object.keys(krakenData.result)[0]; // 取第一个键
        const krakenPrice = parseFloat(krakenData.result[krakenKey].c[0]); // c[0] = last trade price

        // 3. Coinbase: SEI-USD
        const coinbaseUrl = "https://api.exchange.coinbase.com/products/SEI-USD/ticker";
        const coinbaseData = await fetch(coinbaseUrl).then(res => res.json());
        const coinbasePrice = parseFloat(coinbaseData.price);

        const prices = [binancePrice, krakenPrice, coinbasePrice];
        return {
            priceUsd: new Decimal(prices.reduce((sum, p) => sum + p, 0) / prices.length)
        }
    } catch (err) {
        console.error("Error fetching SEI prices:", err.message);
    }
}


export async function ParseLog(topic, data) {
    try {
        const iface = new ethers.Interface(ABI.ABI);
        const log = {
            topics: [topic],
            data: data,
        };
        const parsed = iface.parseLog(log);
        let kv = parsed.args.toObject();
        for (let k in kv) {
            kv[k] = typeof kv[k] === 'bigint' ? kv[k].toString() : kv[k];
        }
        return JSON.stringify(kv);
    } catch (error) {
        const iface = new ethers.Interface(ABIex.ABI);
        const log = {
            topics: [topic],
            data: data,
        };
        const parsed = iface.parseLog(log);
        let kv = parsed.args.toObject();
        for (let k in kv) {
            kv[k] = typeof kv[k] === 'bigint' ? kv[k].toString() : kv[k];
        }
        return JSON.stringify(kv);
    }
}

export async function ParseSwapLog(topic, data) {
    try {
        const topics = topic.split(",")
        const iface = new ethers.Interface(ABI.SwapABI);
        const log = {
            topics: topics,
            data: data,
        };
        const parsed = iface.parseLog(log);
        console.debug(parsed)
        let kv = parsed.args.toObject();
        for (let k in kv) {
            kv[k] = typeof kv[k] === 'bigint' ? kv[k].toString() : kv[k];
        }
        return JSON.stringify(kv);
    } catch (error) {
        const topics = topic.split(",")
        const iface = new ethers.Interface(ABIex.SwapABI);
        const log = {
            topics: topics,
            data: data,
        };
        const parsed = iface.parseLog(log);
        console.debug(parsed)
        let kv = parsed.args.toObject();
        for (let k in kv) {
            kv[k] = typeof kv[k] === 'bigint' ? kv[k].toString() : kv[k];
        }
        return JSON.stringify(kv);
    }
}

export async function GetTokenBalance(tokenAddress, walletAddress) {
    const ERC20_ABI = [
        "function balanceOf(address wallet) view returns (uint256)"
    ];
    const token = new ethers.Contract(tokenAddress, ERC20_ABI, provider);
    try {
        const balance = await token.balanceOf(walletAddress) / BigInt(10**18);
        console.debug(`get token balance, tokenAddress: ${tokenAddress}, walletAddress: ${walletAddress}, balance: ${balance}`)
        return balance.toString();
    } catch (error) {
        console.error('GetTokenBalance failed:', error);
    }
}

export async function Token0(pairAddress) {
    try {
        const PAIR_ABI = [
            "function token0() external view returns (address)",
            "function token1() external view returns (address)",
            "function fee() external view returns (uint24)"
        ];
        const lpPair = new ethers.Contract(pairAddress, PAIR_ABI, provider);

        const result = await lpPair.token0();
        console.debug(`check wrapperAddress is token0 or not, pairAddress: ${pairAddress}, result: ${result}`)
        return result
    } catch (error) {
        console.error('Token0 failed:', error);
    }
}

export async function IsListToken(id, tokenAddress, walletAddress, contractAddress) {
    const PAY_ABI = [
        "function isListToken(address account,address token, uint id) view returns (bool)"
    ];
    const payContract = new ethers.Contract(contractAddress, PAY_ABI, provider);
    try {
        const output = await payContract.isListToken(walletAddress, tokenAddress, BigInt(id));
        console.debug(`get IsListToken, id: ${id}, tokenAddress: ${tokenAddress}, walletAddress: ${walletAddress}, contractAddress: ${contractAddress}`)
        return output;
    } catch (error) {
        console.error('IsListToken failed:', error);
    }
}