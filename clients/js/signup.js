'use strict';
// creates json for new user

const base = "https://api.jimhua32.me";
const userP = "/users";
let hello;
console.log(sessionStorage.getItem("auth"));
function createNewUser() {
    let newUser = {
        "Email": document.getElementById("inputEmail3").value,
        "Password": document.getElementById("inputPassword3").value,
        "PasswordConf": document.getElementById("inputPassword3C").value,
        "FirstName": document.getElementById("fname").value,
        "LastName": document.getElementById("lname").value
    };
    console.log(newUser.Email);
    console.log(newUser);
    console.log(base + userP);
    fetch(base + userP,
        {
            method: "POST",
            body: JSON.stringify(newUser),
            headers: new Headers(
                {"Content-Type": "application/json"}
            )
        }
    ).then(response => {
        if (response.status == 405 || response.status == 400) {
            console.log("error creating new user account");
            console.log(response);
            // return
        }
        console.log(response);
        console.log("yo")
        hello = response;
        let token = [];
        token = response.headers.get("Authorization").split(" ");
        console.log(token);
        // state.auth = token[0];
        sessionStorage.setItem("auth", token[1]);
    })
}

console.log(hello);
console.log(document.getElementById("submitNUser"))

// document.getElementById("submitNUser").addEventListener("click", createNewUser);
// console.log(hello);
document.getElementById("submitNUser").addEventListener("click", (event) => {
    event.preventDefault();
    console.log(document.getElementById("inputEmail3").value);
    console.log(hello)
    createNewUser();
    exports.auth = state.auth;
    // window.location.href="index.html";
})