let productInfo = {
    name: "",
    image: "",
    description: "",
    price: 0,
};

let tossPayments;
async function getProduct(id){
    let res = await fetch(API_URL + "/products/" + id, {
        method: "GET",
        credentials: "include"
    });
    let data = await res.json();
    if(!isOkStatus(res.status)){
        return {
            name: "Product not found",
            image: "https://via.placeholder.com/300",
            description: "Product not found",
            price: 0,
        };
    }
    return data;
}

async function getPurchases(){
    let res = await fetch(API_URL + "/purchase", {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(res.status)){
        return [];
    }
    return await res.json();
}

async function renderPurchaseHistory(){
    let purchaseHistory = await getPurchases();
    let purchaseHistoryContainer = document.getElementById("container");
    purchaseHistoryContainer.innerHTML = "";
    for(let purchase of purchaseHistory){
        let productInfo = await getProduct(purchase.product_id);
        let purchaseElement = document.createElement("div");
        purchaseElement.classList.add("purchase");
        let purchaseItemImg = document.createElement("img");
        purchaseItemImg.src = await getImage(productInfo.image);
        purchaseItemImg.width = 100;
        purchaseItemImg.classList.add("purchase-img");
        purchaseElement.appendChild(purchaseItemImg);

        let purchaseItem = document.createElement("div");
        purchaseItem.classList.add("purchase-item");
        purchaseItem.innerText = 'Name: ' + productInfo.name;
        purchaseElement.appendChild(purchaseItem);

        let purchasePrice = document.createElement("div");
        purchasePrice.classList.add("purchase-price");
        purchasePrice.innerText = "Price: KRW " + purchase.amount;
        purchaseElement.appendChild(purchasePrice);

        let purchaseDate = document.createElement("div");
        purchaseDate.classList.add("purchase-date");
        purchaseDate.innerText = "Date: " + purchase.created_at;
        purchaseElement.appendChild(purchaseDate);

        let purchaseStatus = document.createElement("div");
        purchaseStatus.classList.add("purchase-status");
        purchaseStatus.innerText = "Status: " + purchase.payment_status;
        purchaseElement.appendChild(purchaseStatus);

        let purchaseRequestMessage = document.createElement("div");
        purchaseRequestMessage.classList.add("purchase-request-message");
        purchaseRequestMessage.innerText = "Request Message: " + purchase.request_message;
        purchaseElement.appendChild(purchaseRequestMessage);

        let purchaseAddress = document.createElement("div");
        purchaseAddress.classList.add("purchase-address");
        purchaseAddress.innerText = "Address: " + purchase.address;
        purchaseElement.appendChild(purchaseAddress);

        purchaseElement.appendChild(document.createElement("hr"));

        purchaseHistoryContainer.appendChild(purchaseElement);
    }
}

window.addEventListener("load", async function(){
    let userInfo = await getUserInfo();
    let productId = window.location.hash.substring(1);
    if(productId == ""){
        renderPurchaseHistory();
        return;
    }
    productInfo = await getProduct(productId);
    var clientKey = 'test_ck_ex6BJGQOVDEn5NzqZ9PqrW4w2zNb'
    tossPayments = TossPayments(clientKey);
    const addressField = document.getElementById("address");
    addressField.value = userInfo.address;


    document.getElementById("payment-button").addEventListener("click", async function() {
        let orderId = generateUUID();
        const addressField = document.getElementById("address");
        const requestMessageField = document.getElementById("request-message");
        let userInfo = await getUserInfo();
        let res = await fetch(API_URL + "/purchase/request", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            credentials: "include",
            body: JSON.stringify({
                product_id: productId,
                request_id: orderId,
                product_name: productInfo.name,
                amount: productInfo.price,
                request_message: requestMessageField.value,
                address: addressField.value
            })
        });
        if(!isOkStatus(res.status)){
            let result = await res.json();
            let error = document.getElementById("error");
            error.innerText = result.error;
            return;
        }
        tossPayments.requestPayment('카드', {
            amount: productInfo.price,
            orderId: orderId,
            orderName: productInfo.name,
            customerName: userInfo.real_name,
            successUrl: `${window.location.origin}/payment-success.html`,
            failUrl: `${window.location.origin}/payment-fail.html`
        })
    });
});