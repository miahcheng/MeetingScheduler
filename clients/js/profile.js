`use strict`;

const base = "https://api.jimhua32.me";

fetch(base + "/user/",
    {
        method: "GET",
        headers: {
            "Authorization": sessionStorage.getItem("auth"),
            "Content-Type": "application/json"
        }
    }
).then(response => {
    if (response.status >= 400) {
        console.log("error finding user");
        console.log(response);
        return;
    }
    let userInfo = document.getElementById("prof");
    let name = document.createElement("h1");
    name.classList.add("title");
    let info = document.createElement("p");
    info.classList.add("lead")
    console.log(response);
    // name.innerText = response.FirstN
})