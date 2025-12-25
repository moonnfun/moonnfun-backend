import { ethers, JsonRpcProvider, WebSocketProvider } from "ethers";
import * as store from "./storage.js";
import * as swap from "./swap.js";
import Decimal from 'decimal.js';

const ABI = [
    "function deposit(uint destinationChainId) external payable",
    "function withdraw(address receiver, uint amount, uint sourceChainId) external payable",
    "event Deposit(address caller, uint depositAmount, uint sourceChainId, uint destinationChainId)",
    "event Withdraw(address caller, address receiver, uint withdrawAmount, uint sourceChainId, uint destinationChainId)"
];

const RPC = {
    "97": {
        Name: "Bsc Testnet",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
    "56": {
        Name: "Bsc Mainnet",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
    "1328": {
        Name: "",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
    "1329": {
        Name: "Sei Mainnet",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
    "84532": {
        Name: "Base Sepolia",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
    "8453": {
        Name: "Base Mainnet",
        HTTP: "",
        WSS: "",
        Contract: "",
    },
}
const chainList = ["97", "1328", "84532"];
let withdrawTx = {};
let withdrawTxList = [];

export async function WatchEvents(password) {
    withdrawTxList = await store.loadJson();
    chainList.forEach(chainId => {
        watchEventsForChain(password, chainId);
    });
}

async function watchEventsForChain(password, chainId1) {
    const provider1 = new WebSocketProvider(RPC[chainId1].WSS);
    const wallet1 = new ethers.Wallet("", provider1);
    const watchContract1 = new ethers.Contract(RPC[chainId1].Contract, ABI, wallet1);

    const balance = await GetBalance(chainId1);
    console.log("start watch events,", "chain:", chainId1, "rpc:", RPC[chainId1], "balance:", balance);
    
    watchContract1.on("Deposit", async (caller, depositAmount, sourceChainId, destinationChainId) => {
        console.log(`~~~~~~~~~~~~~~~~~~~~~Received Deposit~~~~~~~~~~~~~~~~~~~~~~`);
        console.log("caller: ", caller, "depositAmount: ", depositAmount, "sourceChainId: ", sourceChainId, "destinationChainId: ", destinationChainId);
    
        const srcId = new Decimal(sourceChainId).toString();
        const dstId = new Decimal(destinationChainId).toString();
        if (chainList.includes(dstId) === false || srcId === dstId) {
            console.log(`无效的Deposit, sourceChainId: ${sourceChainId} 与 destinationChainId: ${destinationChainId} 不匹配`);
            return;
        }
        const withdrawAmount = await swap.getSwapAmount(srcId, dstId, depositAmount);

        console.log(`${srcId}->${dstId}, `, "depositAmount: ", depositAmount, "withdrawAmount: ", withdrawAmount);

        // withdraw
        let result = await withdraw(caller, sourceChainId, srcId, dstId, withdrawAmount, waleltPrivate);
        if (result !== true) {
            for (let i = 0; i < 3; i++) {
                result = await withdraw(caller, sourceChainId, srcId, dstId, withdrawAmount, waleltPrivate);
                if (result !== true) {
                    console.log(`withdraw failed, index: ${i + 1}`);
                    await new Promise(r => setTimeout(r, 1000));
                } else {
                    break;
                }
            }
        }
    });

    watchContract1.on("Withdraw", async (caller, receiver, withdrawAmount, sourceChainId, destinationChainId) => {
        console.log(`~~~~~~~~~~~~~~~~~~~~~Received Withdraw~~~~~~~~~~~~~~~~~~~~~~`);
        console.log("caller: ", caller, "receiver: ", receiver, "withdrawAmount: ", withdrawAmount, "sourceChainId: ", sourceChainId, "destinationChainId: ", destinationChainId);
    });

    // // bsc deposit
    // let options = {value: ethers.parseEther("0.01"), from: wallet1.address, gasLimit: BigInt(300000), gasPrice: (await provider1.getFeeData()).gasPrice, nonce: await wallet1.getNonce()}
    // const sentTx = await watchContract1.deposit(BigInt(chainPairs[chainId]), options);
    // await sentTx.wait(1);

    // // sei deposit
    // let options = {value: ethers.parseEther("1.013"), from: wallet1.address, gasLimit: BigInt(300000), gasPrice: (await provider1.getFeeData()).gasPrice, nonce: await wallet1.getNonce()}
    // const sentTx = await watchContract1.deposit(BigInt(chainPairs[chainId]), options);
    // await sentTx.wait(1);
}

async function withdraw(caller, sourceChainId, srcId, dstId, withdrawAmount, waleltPrivate) {
    try {
        const provider = new JsonRpcProvider(RPC[dstId].HTTP);
        const wallet = new ethers.Wallet(waleltPrivate, provider);
        const bcontract = new ethers.Contract(RPC[dstId].Contract, ABI, wallet);
        let options = {
            from: wallet.address, 
            gasLimit: BigInt(300000), 
            gasPrice: (await provider.getFeeData()).gasPrice, 
            nonce: await wallet.getNonce()
        }
        console.debug("before withdraw", "options", options)

        const sentTx = await bcontract.withdraw(caller, BigInt(withdrawAmount), sourceChainId, options);
        await sentTx.wait(1);
        console.log(`withdraw successed, withdrawTxList: `, JSON.stringify(withdrawTxList));
        return true;
    } catch(err) {
        console.log(`withdraw failed, error: ${err}, caller: ${caller}, srcId: ${srcId}, dstId: ${dstId}, withdrawAmount: ${withdrawAmount}`)
        return false;
    }
}

export async function GetBalance(chainId) {
    const provider = new JsonRpcProvider(RPC[chainId].HTTP);
    const balance = new Decimal(await provider.getBalance(RPC[chainId].Contract)).dividedBy(new Decimal(10**18));
    return balance.toNumber();
}
