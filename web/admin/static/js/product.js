async function getProducts() {
    let response = await fetch(API_URL + "/products", {
        method: "GET",
        credentials: "include"
    });
    if(response.status == 200){
        return await response.json();
    }
    return [];
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

function deleteProduct(id) {
    fetch(API_URL + "/products/" + id, {
        method: "DELETE",
        credentials: "include"
    }).then(response => {
        if(response.status == 200){
            window.location.reload();
        }
    });
}

async function postImage(image) {
    // authorized.POST("/image", {"image": "(base64 encoded image)"})
    // encode image to base64
    let imageData = await image.arrayBuffer();
    let base64Image = btoa(uint8ToString(new Uint8Array(imageData)));
    let response = await fetch(API_URL + "/image", {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            image: base64Image
        })
    });
    return await response.json();
}

async function addProduct(name, price, image) {
    let postedImage = await postImage(image);
    if(postedImage == null || postedImage.id == null){
        return postedImage? postedImage.message: "Failed to post image";
    }
    if(typeof price != "number"){
        price = parseInt(price);
        if(isNaN(price)){
            return "Price is not a number";
        }
    }
    let imageId = postedImage.id;
    let response = await fetch(API_URL + "/products", {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: name,
            price: price,
            image: imageId
        })
    });
    if(response.status == 200){
        return await response.json();
    }
    return "Failed to add product";
}

async function randerProduct() {
    let products = await getProducts();
    let productTable = document.getElementById("product-list");
    for(let product of products){
        let tr = document.createElement("tr");
        let tdId = document.createElement("td");
        tdId.innerText = product.id;
        tr.appendChild(tdId);
        let tdImage = document.createElement("td");
        let img = document.createElement("img");
        let imageData = await getImage(product.image);
        img.src = imageData;
        img.width = 100;
        tdImage.appendChild(img);
        tr.appendChild(tdImage);

        let tdName = document.createElement("td");
        tdName.innerText = product.name;
        tr.appendChild(tdName);

        let tdPrice = document.createElement("td");
        tdPrice.innerText = product.price;
        tr.appendChild(tdPrice);

        let tdEdit = document.createElement("td");
        let editButton = document.createElement("button");
        editButton.setAttribute("data-id", product.id);
        editButton.innerText = "Edit";
        editButton.addEventListener("click", function(){
            window.location.hash = this.getAttribute("data-id");
        });
        tdEdit.appendChild(editButton);
        tr.appendChild(tdEdit);
        let tdDelete = document.createElement("td");
        let deleteButton = document.createElement("button");
        deleteButton.setAttribute("data-id", product.id);
        deleteButton.innerText = "Delete";
        deleteButton.addEventListener("click", async function(){
            await deleteProduct(this.getAttribute("data-id"));
            window.location.reload();
        });
        tdDelete.appendChild(deleteButton);
        tr.appendChild(tdDelete);
        productTable.appendChild(tr);
    }
    let newButton = document.createElement("button");
    newButton.innerText = "New";
    newButton.addEventListener("click", function(){
        window.location.hash = "new";
    });
    productTable.appendChild(newButton);
}

function randerNewProduct() {
    let productTable = document.getElementById("product-list");
    let form = document.createElement("form");
    form.setAttribute("id", "new-product-form");
    let nameLabel = document.createElement("label");
    nameLabel.setAttribute("for", "name");
    nameLabel.innerText = "Name: ";
    form.appendChild(nameLabel);
    let nameInput = document.createElement("input");
    nameInput.setAttribute("type", "text");
    nameInput.setAttribute("id", "name");
    form.appendChild(nameInput);
    form.appendChild(document.createElement("br"));
    let priceLabel = document.createElement("label");
    priceLabel.setAttribute("for", "price");
    priceLabel.innerText = "Price: ";
    form.appendChild(priceLabel);
    let priceInput = document.createElement("input");
    priceInput.setAttribute("type", "number");
    priceInput.setAttribute("id", "price");
    form.appendChild(priceInput);
    form.appendChild(document.createElement("br"));
    let imageLabel = document.createElement("label");
    imageLabel.setAttribute("for", "image");
    imageLabel.innerText = "Image: ";
    form.appendChild(imageLabel);
    let imageInput = document.createElement("input");
    imageInput.setAttribute("type", "file");
    imageInput.setAttribute("id", "image");
    form.appendChild(imageInput);
    form.appendChild(document.createElement("br"));
    let submitButton = document.createElement("button");
    submitButton.setAttribute("type", "submit");
    submitButton.innerText = "Submit";
    form.appendChild(submitButton);
    form.addEventListener("submit", async function(e){
        e.preventDefault();
        let name = document.getElementById("name").value;
        let price = document.getElementById("price").value;
        let image = document.getElementById("image").files[0];
        let product = await addProduct(name, price, image);
        if(product != null && typeof product == "string"){
            if(document.getElementById("error")){
                document.getElementById("error").remove();
            }
            let p = document.createElement("p");
            p.innerText = "Failed to add product: " + product;
            p.id = "error";
            p.style.color = "red";
            let productTable = document.getElementById("product-list");
            productTable.appendChild(p);
            return;
        }
        if(product){
            window.location.hash = "";
        }
    });
    productTable.appendChild(form);
}

async function randerEditProduct(id) {
    let product = await getProduct(id);
    let form = document.createElement("form");
    form.setAttribute("id", "edit-product-form");
    let nameLabel = document.createElement("label");
    nameLabel.setAttribute("for", "name");
    nameLabel.innerText = "Name: ";
    form.appendChild(nameLabel);
    let nameInput = document.createElement("input");
    nameInput.setAttribute("type", "text");
    nameInput.setAttribute("id", "name");
    nameInput.setAttribute("value", product.name);
    form.appendChild(nameInput);
    form.appendChild(document.createElement("br"));
    let priceLabel = document.createElement("label");
    priceLabel.setAttribute("for", "price");
    priceLabel.innerText = "Price: ";
    form.appendChild(priceLabel);
    let priceInput = document.createElement("input");
    priceInput.setAttribute("type", "number");
    priceInput.setAttribute("id", "price");
    priceInput.setAttribute("value", product.price);
    form.appendChild(priceInput);
    form.appendChild(document.createElement("br"));
    let imageLabel = document.createElement("label");
    imageLabel.setAttribute("for", "image");
    imageLabel.innerText = "Image: ";
    form.appendChild(imageLabel);
    let imageInput = document.createElement("input");
    imageInput.setAttribute("type", "file");
    imageInput.setAttribute("id", "image");
    form.appendChild(imageInput);
    form.appendChild(document.createElement("br"));
    let submitButton = document.createElement("button");
    submitButton.setAttribute("type", "submit");
    submitButton.innerText = "Submit";
    form.appendChild(submitButton);
    form.appendChild(document.createElement("br"));
    form.addEventListener("submit", async function(e){
        e.preventDefault();
        let name = document.getElementById("name").value;
        let price = document.getElementById("price").value;
        let imageId = product.image;
        if(document.getElementById("image").files.length > 0){
            let image = document.getElementById("image").files[0];
            imageId = (await postImage(image)).id;
        }
        if(typeof price != "number"){
            price = parseInt(price);
            if(isNaN(price)){
                return "Price is not a number";
            }
        }
        let response = await fetch(API_URL + "/products/" + id, {
            method: "PUT",
            credentials: "include",
            body: JSON.stringify({
                name: name,
                price: price,
                image: imageId
            })
        });
        if(response.status == 200){
            window.location.hash = "";
        }
    });
    document.body.appendChild(form);
}

window.addEventListener("load", async ()=>{
    switch(window.location.hash){
        case "":
            randerProduct();
            break;
        case "#new":
            randerNewProduct();
            break;
        default:
            let id = window.location.hash.slice(1);
            randerEditProduct(id);
            break;
    }
});

window.addEventListener("hashchange", ()=>{
    window.location.reload();
});