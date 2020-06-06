'use strict';

let state = {
    auth: "",
    testing: "1"
};


function toggleLogin(loginD, signD) {
    document.getElementById("loginUser").style.display = loginD;
    document.getElementById("signUp").style.display = signD;
}

sessionStorage.setItem("auth", "");
const base = "https://api.jimhua32.me";
const user = "/user/";
const myuser = "/user/id";
const sessions = "/sessions";
const mySession = "/sessions/mine";

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
        if (response.status >= 400) {
            console.log("error logging in user");
            console.log(response);
            window.alert("Incorrect email or password");
            return;
        } else {
            console.log(response)
            let token = response.headers.get("Authorization");
            sessionStorage.setItem("auth", token);
            sessionStorage.setItem("loggedIn", true);
        }
        if (sessionStorage.getItem("loggedIn")) {
            window.location.href="home.html";
        }
    }
    );
});

document.getElementById("submitNUser").addEventListener("click", (event) => {
    event.preventDefault();
    let newUser = {
        "Email": document.getElementById("inputEmail3").value,
        "Password": document.getElementById("inputPassword3").value,
        "PasswordConf": document.getElementById("inputPassword3C").value,
        "FirstName": document.getElementById("fname").value,
        "LastName": document.getElementById("lname").value
    };
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
            if (response.status === 400) {
                window.alert("Error creating new user");
            }
        }
        console.log(response);
        window.alert("User signed up! Please log in");
        toggleLogin("block", "none");
    });
});

window.logoutUser = function() {
    document.getElementById("logoutUser").addEventListener("click", (event) => {
        sessionStorage.setItem("auth", "");
        window.location.href="index.html";
    })
}
