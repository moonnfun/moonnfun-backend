import express from "express";
import readline from 'readline';
import * as swap from './swap.js';
import * as store from './storage.js';
import * as monitor from './monitor.js';

let originalLog = console.log;
console.log = (...args) => {
    const timestamp = new Date().toLocaleString(); // 当前时间 ISO 格式
    originalLog(`\n[${timestamp}]`, ...args); // 打印时间和原始参数
};
// console.debug = function(){};

process.on('uncaughtException', function (err) {
    console.error('uncaughtException', err);
})

process.on('unhandledRejection', function (err, promise) {
    console.error('unhandledRejection', err);
})

const askHiddenInput = (query) => {
    return new Promise((resolve) => {
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });

        rl.question(query, (input) => {
            rl.history = rl.history.slice(1); // 防止记录历史
            resolve(input);
            rl.close();
        });

        // 在用户输入时隐藏字符
        rl._writeToOutput = () => {
            rl.output.write('*');
        };
    });
};

async function initLog() {
    let originalLog = console.log;
    console.log = (...args) => {
        const timestamp = new Date().toLocaleString(); // 当前时间 ISO 格式
        originalLog(`[${timestamp}]`, ...args); // 打印时间和原始参数
    };
    // console.debug = (...args) => {}
}

async function runMonitor(password) {
    await monitor.WatchEvents(password);
}

function runServer() {
    const app = express();
    const port = process.argv[process.argv.length - 1];
  
    app.get('/', (req, res) => {
        res.send('Welcome to bridge server!');
    });

    app.get('/bapi/price', async (req, res) => {
        const srcChainId = req.query.srcChainId;
        const dstChainId = req.query.dstChainId;
        const response = {
            data: await swap.GetPrice(srcChainId, dstChainId),
            error: ""
        }
        res.send(JSON.stringify(response));
    });

    app.get('/bapi/balance', async (req, res) => {
        const chainId = req.query.chainId;
        const response = {
            data: await monitor.GetBalance(chainId),
            error: ""
        }
        res.send(JSON.stringify(response));
    });

    app.get('/bapi/withdraw/tx', async (req, res) => {
        res.send(JSON.stringify(await monitor.GetWithdrawTx(req.query.address, req.query.txhash)));
    });

    app.listen(port, () => {
        console.log(`Server is running on port ${port}`);
    });
}

async function main() {
    console.log(process.argv)
    const password = await askHiddenInput("Please enter the admin password: ");

    await initLog();
    await store.initJson("./history.json", '[]');
    await runMonitor(password);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
    .then(() => runServer())
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });