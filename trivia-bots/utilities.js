
module.exports = {
    sleep: sleep
}

async function sleep(s){
    return new Promise(resolve => setTimeout(resolve,s * 1000));
}