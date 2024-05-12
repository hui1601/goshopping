async function logout(){
    let response = await fetch(API_URL + "/auth/logout", {
        method: "POST",
        credentials: "include"
    });
    if(isOkStatus(response.status)){
        // window.location.href = "/";
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
    if(!(await isLogined())){
        window.location.href = "/login.html";
        return;
    }
    let logoutButton = document.getElementById("logout");
    logoutButton.addEventListener("click", ()=>!((logout()||true) && (this.document.body.appendChild(this.document.createElement("div").appendChild(this.document.createTextNode("logout successful. ")).appendChild(this.document.createElement("a").appendChild(this.document.createTextNode("Login")).setAttribute("href", "/login.html"))))));
});