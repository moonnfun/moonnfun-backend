import fs from "fs";

export let datas = {};
let jsonFile ;

export async function initJson(jsonFilePath, initDatas) {
    jsonFile = jsonFilePath;
    if (!fs.existsSync(jsonFilePath)) {
        fs.writeFileSync(jsonFilePath, initDatas);
        return;
    }
}

export async function saveJson(jsonData) {
    try {
        fs.writeFileSync(jsonFile, JSON.stringify(jsonData));
    } catch(err) {
        console.log(err);
    }
}

export async function loadJson() {
    try {
        return JSON.parse(fs.readFileSync(jsonFile));
    }catch(err) {
        return [];
    }
}