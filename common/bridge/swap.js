import Decimal from "decimal.js";

export async function getSwapAmount(srcChainId, dstChainId, amount) {
    const price = await GetPrice(srcChainId, dstChainId);
    return new Decimal(amount).mul(new Decimal(price)).toFixed(0);
}

export async function GetPrice(srcChainId, dstChainId) {
    const swapPair = `${srcChainId}-${dstChainId}`;
    const bSeiBnbTestnet = swapPair.includes("97") && swapPair.includes("1328");
    const bSeiBnbMainnet = swapPair.includes("56") && swapPair.includes("1329");
    const bBnbEthTestnet = swapPair.includes("97") && swapPair.includes("84532");
    const bBnbEthMainnet = swapPair.includes("56") && swapPair.includes("8453");
    const bSeiEthTestnet = swapPair.includes("1328") && swapPair.includes("84532");
    const bSeiEthMainnet = swapPair.includes("1329") && swapPair.includes("8453");
    if (bSeiBnbTestnet || bSeiBnbMainnet) {
        const seiBnbPrice = await getSeiBnbPrice();

        // BNB => SEI, amount(BNB) * SeiBnbPrice
        if (srcChainId === "97" || srcChainId === "56") {
            return new Decimal(1).dividedBy(new Decimal(seiBnbPrice)).toString();
        }

        // SEI => BNB, amount(SEI) / SeiBnbPrice
        if (dstChainId === "97" || dstChainId === "56") {
            return new Decimal(seiBnbPrice).toString();
        }
    } else if (bBnbEthTestnet || bBnbEthMainnet) {
        const bnbEthPrice = await getBnbEthPrice();

        // BNB => ETH, amount(BNB) * SeiBnbPrice
        if (srcChainId === "97" || srcChainId === "56") {
            return new Decimal(bnbEthPrice).toString();
            
        }

        // SEI => BNB, amount(SEI) / SeiBnbPrice
        if (dstChainId === "97" || dstChainId === "56") {
            return new Decimal(1).dividedBy(new Decimal(bnbEthPrice)).toString();
        }
    } else if (bSeiEthTestnet || bSeiEthMainnet) {
        const seiBnbPrice = await getSeiBnbPrice();
        const bnbEthPrice = await getBnbEthPrice();
        const seiEthPrice = new Decimal(seiBnbPrice).mul(new Decimal(bnbEthPrice));

        // ETH => SEI, amount(BNB) * SeiBnbPrice
        if (srcChainId === "84532" || srcChainId === "8453") {
            return new Decimal(1).dividedBy(seiEthPrice).toString();
        }

        // SEI => ETH, amount(SEI) / SeiBnbPrice
        if (dstChainId === "84532" || dstChainId === "8453") {
            return new Decimal(seiEthPrice).toString();
        }
    }
    return 0;
}

export async function getSeiBnbPrice() {
    const binanceUrl = "https://api.binance.com/api/v3/ticker/price?symbol=SEIBNB";
    const binanceData = await fetch(binanceUrl).then(res => res.json());
    const seiBnbPrice = parseFloat(binanceData.price);
    return seiBnbPrice;
}

export async function getBnbEthPrice() {
    const binanceUrl = "https://api.binance.com/api/v3/ticker/price?symbol=BNBETH";
    const binanceData = await fetch(binanceUrl).then(res => res.json());
    const seiBnbPrice = parseFloat(binanceData.price);
    return seiBnbPrice;
}