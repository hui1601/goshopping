async function register(){
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    let real_name = document.getElementById("real_name").value;
    let email = document.getElementById("email").value;
    let address = document.getElementById("address").value;
    let phone_number = document.getElementById("phone_number").value;
    let data = {
        username: username,
        password: password,
        real_name: real_name,
        email: email,
        address: address,
        phone_number: phone_number
    }
    let response = await fetch(API_URL + "/auth/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify(data)
    });
    if(isOkStatus(response.status)){
        window.location.href = "/login.html";
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

window.addEventListener("load", function(){
    let registerButton = document.getElementById("register");
    registerButton.addEventListener("click", ()=>!(register()||true));
});