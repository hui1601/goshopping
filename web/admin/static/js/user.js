async function getUsers(){
    let response = await fetch(API_URL + "/admin/user", {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(response.status)){
        return [];
    }
    return await response.json();
}

async function deleteUser(userId){
    let response = await fetch(API_URL + "/admin/user", {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({user_id: userId})
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

async function getUserInfoById(userId){
    let response = await fetch(API_URL + "/admin/user/" + userId, {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(response.status)){
        return null;
    }
    return await response.json();
}

async function updateUser(userId, data){
    let response = await fetch(API_URL + "/admin/user/" + userId, {
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
    if(document.getElementById("error") == null){
        let error = document.createElement("div");
        error.id = "error";
        document.body.appendChild(error);
    }
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}

async function grantAdmin(userId){
    let response = await fetch(API_URL + "/admin/permission/grant/" + userId, {
        method: "GET",
        credentials: "include"
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
}

async function revokeAdmin(userId){
    let response = await fetch(API_URL + "/admin/permission/revoke/" + userId, {
        method: "GET",
        credentials: "include"
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
}

async function getAdmins(){
    let response = await fetch(API_URL + "/admin/permission", {
        method: "GET",
        credentials: "include"
    });
    if(!isOkStatus(response.status)){
        return [];
    }
    return await response.json();
}

async function renderAdminPage(){
    let users = await getUsers();
    let admins = await getAdmins();
    let table = document.createElement("table");
    let header = table.createTHead().insertRow();
    header.insertCell().innerText = "ID";
    header.insertCell().innerText = "User ID";
    header.insertCell().innerText = "Username";
    header.insertCell().innerText = "Email";
    header.insertCell().innerText = "Address";
    header.insertCell().innerText = "Phone";
    header.insertCell().innerText = "Permission";
    header.insertCell().innerText = "Action";
    for(let user of users){
        let row = table.insertRow();
        let id = row.insertCell();
        id.innerText = user.id;
        let user_id = row.insertCell();
        user_id.innerText = user.username;
        let username = row.insertCell();
        username.innerText = user.real_name;
        let email = row.insertCell();
        email.innerText = user.email;
        let address = row.insertCell();
        address.innerText = user.address;
        let phone = row.insertCell();
        phone.innerText = user.phone_number;
        let permission = row.insertCell();
        permission.innerText = admins.find(admin=>admin.id == user.id)?"Admin":"User";
        let admin = row.insertCell();
        let grantButton = document.createElement("button");
        grantButton.innerText = "Grant";
        grantButton.setAttribute("data-user-id", user.id);
        grantButton.addEventListener("click", async function(){
            if(await grantAdmin(this.getAttribute("data-user-id"))){
                window.location.reload();
            }
        });
        admin.appendChild(grantButton);
        let revokeButton = document.createElement("button");
        revokeButton.innerText = "Revoke";
        revokeButton.setAttribute("data-user-id", user.id);
        revokeButton.addEventListener("click", async function(){
            if(await revokeAdmin(this.getAttribute("data-user-id"))){
                window.location.reload();
            }
        });
        admin.appendChild(revokeButton);
        let deleteButton = document.createElement("button");
        deleteButton.innerText = "Delete";
        deleteButton.setAttribute("data-user-id", user.id);
        deleteButton.addEventListener("click", async function(){
            // prevent delete admin
            if(admins.find(admin=>admin.id == user.id)){
                if(document.getElementById("error") == null){
                    let error = document.createElement("div");
                    error.id = "error";
                    document.body.appendChild(error);
                }
                let error = document.getElementById("error");
                error.innerText = "Can't delete admin";
                return;
            }
            if(await deleteUser(this.getAttribute("data-user-id"))){
                window.location.reload();
            }
        });
        row.insertCell().appendChild(deleteButton);
        let editButton = document.createElement("button");
        editButton.innerText = "Edit";
        editButton.setAttribute("data-user-id", user.id);
        editButton.addEventListener("click", async function(){
            window.location.hash = this.getAttribute("data-user-id");
        });
        let viewPurchaseButton = document.createElement("button");
        viewPurchaseButton.innerText = "View Purchase";
        viewPurchaseButton.setAttribute("data-user-id", user.id);
        viewPurchaseButton.addEventListener("click", async function(){
            window.location.href = "/admin/purchase.html#" + this.getAttribute("data-user-id");
        });
        row.insertCell().appendChild(editButton);
    }
    document.getElementById("container").appendChild(table);
}

async function renderAdminUserInfoEdit(userId){
    let userInfo = await getUserInfoById(userId);
    let table = document.createElement("table");
    let header = table.createTHead().insertRow();
    header.insertCell().innerText = "ID";
    header.insertCell().innerText = "User ID";
    header.insertCell().innerText = "Username";
    header.insertCell().innerText = "Email";
    header.insertCell().innerText = "Address";
    header.insertCell().innerText = "Phone";
    let row = table.insertRow();
    let id = row.insertCell();
    id.innerText = userId;
    let user_id = row.insertCell();
    user_id.innerText = userInfo.username;
    let username = row.insertCell();
    let usernameInput = document.createElement("input");
    usernameInput.setAttribute("type", "text");
    usernameInput.value = userInfo.real_name;
    username.appendChild(usernameInput);
    let email = row.insertCell();
    let emailInput = document.createElement("input");
    emailInput.setAttribute("type", "text");
    emailInput.value = userInfo.email;
    email.appendChild(emailInput);
    let address = row.insertCell();
    let addressInput = document.createElement("input");
    addressInput.setAttribute("type", "text");
    addressInput.value = userInfo.address;
    address.appendChild(addressInput);
    let phone = row.insertCell();
    let phoneInput = document.createElement("input");
    phoneInput.setAttribute("type", "text");
    phoneInput.value = userInfo.phone_number;
    phone.appendChild(phoneInput);
    let updateButton = document.createElement("button");
    updateButton.innerText = "Update";
    updateButton.addEventListener("click", async function(){
        if(await updateUser(userId, {real_name: usernameInput.value, email: emailInput.value, address: addressInput.value, phone_number: phoneInput.value})){
            window.location.href = "/admin/user.html";
        }
    });
    document.getElementById("container").appendChild(table);
    document.getElementById("container").appendChild(updateButton);
}

window.addEventListener("load", async function(){
    if(window.location.hash == ""){
        await renderAdminPage();
        return;
    }
    renderAdminUserInfoEdit(window.location.hash.substring(1));
});

window.addEventListener("hashchange", async function(){
    this.window.location.reload();
});