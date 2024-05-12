async function renderMyPage(){
    let userInfo = await getUserInfo();
    if(userInfo == null || userInfo.error){
        window.location.href = "/login.html";
        return;
    }
    let userId = document.getElementById("userid");
    userId.value = userInfo.user_id;
    let username = document.getElementById("username");
    username.value = userInfo.real_name;
    let email = document.getElementById("email");
    email.value = userInfo.email;
    let address = document.getElementById("address");
    address.value = userInfo.address;
    let phone = document.getElementById("phone");
    phone.value = userInfo.phone_number;
    let password = document.getElementById("password");
    password.value = "";
    let oldPassword = document.getElementById("oldPassword");
    oldPassword.value = "";
    let updateButton = document.getElementById("edit");
    updateButton.addEventListener("click", ()=>!(updateUserInfo()||true));
    let deleteButton = document.getElementById("delete");
    deleteButton.addEventListener("click", ()=>!(deleteAccount()||true));
}

async function updateUserInfo(){
    let username = document.getElementById("username").value;
    let email = document.getElementById("email").value;
    let address = document.getElementById("address").value;
    let phone = document.getElementById("phone").value;
    let data = {
        real_name: username,
        email: email,
        address: address,
        phone_number: phone,
        new_password: document.getElementById("password").value,
        old_password: document.getElementById("oldPassword").value
    }
    let response = await fetch(API_URL + "/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify(data)
    });
    if(isOkStatus(response.status)){
        window.location.href = "/";
        return true;
    }
    let result = await response.json();
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}

async function deleteAccount(){
    let response = await fetch(API_URL + "/user", {
        method: "DELETE",
        credentials: "include"
    });
    if(isOkStatus(response.status)){
        window.location.href = "/login.html";
        return true;
    }
    let result = await response.json();
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}


window.addEventListener("load", async function(){
    await renderMyPage();
});