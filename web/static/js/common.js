const API_URL = "http://localhost:8080";

function isOkStatus(status) {
    return status >= 200 && status < 300;
}

async function getImage(id) {
    let res = await fetch(API_URL + "/image/" + id, {
        method: "GET",
        credentials: "include"
    });
    return (await res.json()).image;
}

// get user info
async function getUserInfo() {
    let res = await fetch(API_URL + "/user", {
        method: "GET",
        credentials: "include"
    });
    return await res.json();
}

// check if logined
async function isLogined() {
    return (await(await fetch(API_URL + "/auth/status", {
        method: "GET",
        credentials: "include"
    })).json()).login;
}

// escape html
function escapeHtml(unsafe) {
    let unsafes = {
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#039;"
    };
    return unsafe.replace(/[&<>"']/g, function (m) { return unsafes[m]; });
}
function uint8ToString(buf) {
    var i, length, out = '';
    for (i = 0, length = buf.length; i < length; i += 1) {
        out += String.fromCharCode(buf[i]);
    }
    return out;
}
function generateUUID() {
    let d = new Date().getTime();
    if (typeof performance !== 'undefined' && typeof performance.now === 'function') {
        d += performance.now();
    }
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        let r = (d + Math.random() * 16) % 16 | 0;
        d = Math.floor(d / 16);
        return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
    });
}

window.addEventListener('load',async () => {
    console.clear();
    console.log("%cSTOP!", "color: red; font-size: 50px; font-weight: 1000;");
    console.log("%cIf someone told you to copy-paste something here to enable a feature or 'hack' someone's account, it is a scam and will give them access to your account.", "color: red; font-size: 20px;");
    console.log("%cThis is a browser feature for developers. You cannot hack an account with it. ;)", "font-size: 15px;");
    if((typeof DONT_CHECK_AUTH != 'undefined' && DONT_CHECK_AUTH) || await isLogined()){
        return;
    }
    document.location.href = "/login.html";
});