// [{"product_id": productID, "request_id": requestID, "amount": amount, "created_at": createdAt, "request_message": requestMessage, "address": address, "payment_status": paymentStatus}]
async function getPurchases(){
    let response = await fetch(API_URL + "/admin/purchase", {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(response.status)){
        return [];
    }
    return await response.json();
}

async function getPurchase(purchaseId){
    let response = await getPurchases();
    return response.find(purchase => purchase.request_id == purchaseId);
}

async function getUserPurchase(userId){
    let response = await fetch(API_URL + "/admin/user/purchase/" + userId, {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(response.status)){
        return [];
    }
    return await response.json();
}

async function editPurchase(purchaseId, purchase){
    let response = await fetch(API_URL + "/admin/purchase/" + purchaseId, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify(purchase)
    });
    if(isOkStatus(response.status)){
        return true;
    }
    let result = await response.json();
    if(document.getElementById("error") == null){
        let error = document.createElement("div");
        error.id = "error";
        document.body.appendChild(error);
    }
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}

// create table and insert into id=container
// when click edit button, change hash to #edit-{purchase_id}
async function renderPurchase(){
    let purchases = await getPurchases();
    let container = document.getElementById("container");
    let table = document.createElement("table");
    container.appendChild(table);
    let thead = document.createElement("thead");
    table.appendChild(thead);
    let tr = document.createElement("tr");
    thead.appendChild(tr);
    let ths = ["Product ID", "Request ID", "Amount", "Created At", "Request Message", "Address", "Payment Status", "Edit"];
    for(let th of ths){
        let thElement = document.createElement("th");
        thElement.innerText = th;
        tr.appendChild(thElement);
    }
    let tbody = document.createElement("tbody");
    table.appendChild(tbody);
    for(let purchase of purchases){
        let tr = document.createElement("tr");
        tbody.appendChild(tr);
        let tds = [purchase.product_id, purchase.request_id, purchase.amount, purchase.created_at, purchase.request_message, purchase.address, purchase.payment_status];
        for(let td of tds){
            let tdElement = document.createElement("td");
            tdElement.innerText = td;
            tr.appendChild(tdElement);
        }
        let td = document.createElement("td");
        tr.appendChild(td);
        let button = document.createElement("button");
        button.innerText = "Edit";
        button.setAttribute("data-id", purchase.request_id);
        button.addEventListener("click", function(){
            window.location.hash = "edit-" + this.getAttribute("data-id");
        });
        td.appendChild(button);
    }
}

// get user id from hash
// get user purchase by user id
// create table and insert into id=container
// when click edit button, change hash to #edit-{purchase_id}
async function renderUserPurchase(){
    let userId = window.location.hash.substring(1);
    let purchases = await getUserPurchase(userId);
    let container = document.getElementById("container");
    let table = document.createElement("table");
    container.appendChild(table);
    let thead = document.createElement("thead");
    table.appendChild(thead);
    let tr = document.createElement("tr");
    thead.appendChild(tr);
    let ths = ["Product ID", "Request ID", "Amount", "Created At", "Request Message", "Address", "Payment Status", "Edit"];
    for(let th of ths){
        let thElement = document.createElement("th");
        thElement.innerText = th;
        tr.appendChild(thElement);
    }
    let tbody = document.createElement("tbody");
    table.appendChild(tbody);
    for(let purchase of purchases){
        let tr = document.createElement("tr");
        tbody.appendChild(tr);
        let tds = [purchase.product_id, purchase.request_id, purchase.amount, purchase.created_at, purchase.request_message, purchase.address, purchase.payment_status];
        for(let td of tds){
            let tdElement = document.createElement("td");
            tdElement.innerText = td;
            tr.appendChild(tdElement);
        }
        let td = document.createElement("td");
        tr.appendChild(td);
        let button = document.createElement("button");
        button.innerText = "Edit";
        button.setAttribute("data-id", purchase.request_id);
        button.addEventListener("click", function(){
            window.location.hash = "edit-" + this.getAttribute("data-id");
        });
        td.appendChild(button);
    }
}

// get purchase id from hash
// get purchase by purchase id
// create form and insert into id=container
// when click update button, call editPurchase
async function renderEditPurchase(){
    let purchaseId = window.location.hash.substring(6);
    let purchase = await getPurchase(purchaseId);
    let container = document.getElementById("container");
    let form = document.createElement("form");
    container.appendChild(form);
    let labels = ["Product ID", "Request ID", "Amount", "Created At", "Request Message", "Address", "Payment Status"];
    let readonly = ["Product ID", "Request ID", "Created At"];
    for(let label of labels){
        let labelElement = document.createElement("label");
        labelElement.setAttribute("for", label.toLowerCase().replace(" ", "_"));
        labelElement.innerText = label + ": ";
        form.appendChild(labelElement);
        let inputElement = document.createElement("input");
        inputElement.setAttribute("type", "text");
        inputElement.setAttribute("id", label.toLowerCase().replace(" ", "_"));
        if(readonly.includes(label)){
            inputElement.setAttribute("readonly", "readonly");
        }
        inputElement.value = purchase[label.toLowerCase().replace(" ", "_")];
        form.appendChild(inputElement);
        form.appendChild(document.createElement("br"));
    }
    let updateButton = document.createElement("button");
    updateButton.innerText = "Update";
    updateButton.addEventListener("click", async function(){
        let purchase = {};
        // let purchaseId = window.location.hash.substring(5);
        for(let label of labels){
            purchase[label.toLowerCase().replace(" ", "_")] = document.getElementById(label.toLowerCase().replace(" ", "_")).value;
        }
        purchase.request_id = undefined;
        purchase.product_id = undefined;
        purchase.created_at = undefined;
        purchase.amount = parseInt(purchase.amount);
        if(isNaN(purchase.amount)){
            purchase.amount = 0;
        }
        if(await editPurchase(purchaseId, purchase)){
            window.location.hash = "";
        }
    });
    form.appendChild(updateButton);
}

window.addEventListener("load", async function(){
    if(window.location.hash == ""){
        renderPurchase();
        return;
    }
    if(window.location.hash.substring(1, 5) == "edit"){
        renderEditPurchase();
        return;
    }
    renderUserPurchase();
});

window.addEventListener("hashchange", async function(){
    this.window.location.reload();
});