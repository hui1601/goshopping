async function login(){
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    let data = {
        username: username,
        password: password
    }
    let response = await fetch(API_URL + "/auth/login", {
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
    if(document.getElementById("error") == null){
        let error = document.createElement("div");
        error.id = "error";
        document.body.appendChild(error);
    }
    let error = document.getElementById("error");
    error.innerText = result.error;
    return false;
}

window.addEventListener("load", async function(){
    let loginButton = document.getElementById("login");
    if(await isLogined()){
        window.location.href = "/";
        return;
    }
    loginButton.addEventListener("click", ()=>!(login()||true));
});