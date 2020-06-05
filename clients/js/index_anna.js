'use strict';

let state = {
    auth: "",
    testing: "1"
};


function toggleLogin(loginD, signD) {
    document.getElementById("loginUser").style.display = loginD;
    document.getElementById("signUp").style.display = signD;
}

console.log(sessionStorage.getItem("auth"));
const base = "https://api.jimhua32.me";
const user = "/user/";
const myuser = "/user/id";
const sessions = "/sessions";
const mySession = "/sessions/mine";

function isLoggedIn() {
    // return state.auth === "";
    return sessionStorage.getItem("auth") === "";
}


document.getElementById("goToLog").addEventListener("click", (event) => {
    event.preventDefault();
    toggleLogin("block", "none");
});

document.getElementById("goToSign").addEventListener("click", (event) => {
    event.preventDefault();
    toggleLogin("none", "block");
});

document.getElementById("submitLog").addEventListener("click", function(event) {
    event.preventDefault();
    let email = document.getElementById("exampleInputEmail1").value;
    let pass = document.getElementById("exampleInputPassword1").value;
    console.log(base + sessions);
    fetch(base + sessions,
        {
            method: "POST",
            body: JSON.stringify({
                Email: email,
                Password: pass
            }),
            headers: new Headers({
                "Content-Type": "application/json",
            })
        }
    ).then((response) => {
        console.log("hello")
        console.log(response);
        if (response.status >= 400) {
            console.log("error logging in user");
            console.log(response);
            return;
        } else {
            console.log(response)
            // let token = [];
            let token = response.headers.get("Authorization");
            console.log(token);
            sessionStorage.setItem("auth", token);
            sessionStorage.setItem("loggedIn", true);
        }
        console.log("hello2");
        console.log(sessionStorage.getItem("loggedIn"));
        if (sessionStorage.getItem("loggedIn")) {
            window.location.href="home.html";
        }
    }
    );
});

console.log(sessionStorage.getItem("auth"));

// creates json for new user
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
    console.log(base + "/users");
    fetch(base + "/users",
        {
            method: "POST",
            body: JSON.stringify(newUser),
            headers: new Headers(
                {"Content-Type": "application/json",}
            )
        }
    ).then(async response => {
        if (response.status == 405 || response.status == 400) {
            console.log("error creating new user account");
            console.log(response);
            // return
        }
        console.log(response);
        let token = [];
        token = response.headers.get("Authorization").split(" ");
        console.log(token);
        sessionStorage.setItem("auth", token[1]);
        window.alert("User signed up! Please log in");
        
    })
}

// console.log(document.getElementById("submitNUser"))

document.getElementById("submitNUser").addEventListener("click", (event) => {
    event.preventDefault();
    console.log(document.getElementById("inputEmail3").value);
    // createNewUser();
    let newUser = {
        "Email": document.getElementById("inputEmail3").value,
        "Password": document.getElementById("inputPassword3").value,
        "PasswordConf": document.getElementById("inputPassword3C").value,
        "FirstName": document.getElementById("fname").value,
        "LastName": document.getElementById("lname").value
    };
    console.log(newUser.Email);
    console.log(newUser);
    console.log(base + "/users");
    fetch(base + "/users",
        {
            method: "POST",
            body: JSON.stringify(newUser),
            headers: new Headers(
                {"Content-Type": "application/json",}
            )
        }
    ).then(response => {
        if (response.status >= 400) {
            console.log("error creating new user account");
            console.log(response);
            // return
        }
        console.log(response);
        let token = [];
        // token = response.headers.get("Authorization").split(" ");
        // console.log(token);
        // sessionStorage.setItem("auth", token[1]);
        window.alert("User signed up! Please log in");
        toggleLogin("block", "none");
    });
    // exports.auth = state.auth;
});

window.logoutUser = function() {
    document.getElementById("logoutUser").addEventListener("click", (event) => {
        sessionStorage.setItem("auth", "");
        window.location.href="index.html";
    })
}
