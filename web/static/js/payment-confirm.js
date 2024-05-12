function getParms(key){
    let url = new URL(window.location.href);
    return url.searchParams.get(key);
}


function getParamAndCatch(key){
    let value = getParms(key);
    if(value === null){
        if(document.getElementById("error") !== null){
            let error = document.createElement("p");
            error.id = "error";
            document.body.appendChild(error);
        }
        error.innerText = key + " is null";
        return null;
    }
    return value;
}
async function confirmPayment(){
    let paymentKey = getParamAndCatch("paymentKey"), orderId = getParamAndCatch("orderId"), amount = getParamAndCatch("amount");
    if(paymentKey === null || orderId === null || amount === null){
        return false;
    }
    if(typeof amount !== "number"){
        amount = parseInt(amount);
        if(isNaN(amount)){
            return false;
        }
    }
    let data = {
        paymentKey: paymentKey,
        orderId: orderId,
        amount: amount
    }
    let response = await fetch(API_URL + "/purchase/confirm", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify(data)
    });
    if(isOkStatus(response.status)){
        return true;
    }
    let result = await response.json();
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}

window.addEventListener("load", async function(){
    if(await confirmPayment()){
        let success = document.createElement("p");
        success.innerText = "Payment success!";
        document.body.appendChild(success);
        // setTimeout(()=>window.location.href = "/", 3000);
        return;
    }
    let fail = document.createElement("p");
    fail.innerText = "Payment failed!";
    document.body.appendChild(fail);
});