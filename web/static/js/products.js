async function getProducts(){
    let res = await fetch(API_URL + "/products", {
        method: "GET",
        credentials: "include"
    });
    let data = await res.json();
    return data;
}

async function renderProducts(){
    let products = await getProducts();
    let productsContainer = document.getElementById("products");
    productsContainer.innerHTML = "";
    for(let product of products){
        let productElement = document.createElement("div");
        let headerElement = document.createElement("h2");
        let imgElement = document.createElement("img");
        let priceElement = document.createElement("p");
        let hrElement = document.createElement("hr");
        productElement.className = "product";
        headerElement.innerText = product.name;
        let imageData = await getImage(product.image);
        imgElement.src = imageData;
        imgElement.alt = product.name;
        imgElement.width = 300;
        priceElement.innerText = "KRW " + product.price;
        productElement.appendChild(headerElement);
        productElement.appendChild(imgElement);
        productElement.appendChild(priceElement);
        productElement.appendChild(hrElement);
        productElement.addEventListener("click", () => {
            window.location.hash = product.id;
        });
        productsContainer.appendChild(productElement);
    }
}

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

async function renderProduct(id){
    let product = await getProduct(id);
    let productElement = document.createElement("div");
    let headerElement = document.createElement("h2");
    let imgElement = document.createElement("img");
    let descriptionElement = document.createElement("p");
    let priceElement = document.createElement("p");
    let buttonElement = document.createElement("button");
    productElement.className = "product";
    headerElement.innerText = product.name;
    let imageData = await getImage(product.image);
    imgElement.src = imageData;
    imgElement.alt = product.name;
    imgElement.width = 300;
    
    // descriptionElement.innerText = product.description;
    priceElement.innerText = "KRW " + product.price;
    buttonElement.innerText = "Buy";
    productElement.appendChild(headerElement);
    productElement.appendChild(imgElement);
    productElement.appendChild(descriptionElement);
    productElement.appendChild(priceElement);
    productElement.appendChild(buttonElement);
    let productsContainer = document.getElementById("products");
    productsContainer.innerHTML = "";
    productsContainer.appendChild(productElement);
    buttonElement.addEventListener("click", async () => {
        location.href = "/payment.html#" + id;
    });
}

window.addEventListener("load", async () => {
    if(window.location.hash == ""){
        renderProducts();
    } else{
        let id = window.location.hash.substring(1);
        renderProduct(id);
    }
});
window.addEventListener("hashchange", async () => {
    window.location.reload();
});