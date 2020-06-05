'use strict';

let state = {
    auth: "",
    testing: "1"
};

const base = "https://api.blah.com";
const user = "/user/";
const myuser = "/user/id";
const sessions = "/sessions";
const mySession = "/sessions/mine";

function isLoggedIn() {
    // return state.auth === "";
    return sessionStorage.getItem("auth") === "";
}

function loginUser() {
    // let form = document.getElementById("loginAll");
    let email = document.getElementById("exampleInputEmail1").value;
    let pass = document.getElementById("exampleInputPassword1").value;
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
    ).then(response => {
        if (response.status == 415 || response.status == 405) {
            console.log("error logging in user");
            console.log(response);
            return;
        }
        let token = [];
        token = response.headers.get("Authorization").split(" ");
        console.log(token);
        // state.auth = token[0];
        sessionStorage.setItem("auth", token);
    }
    )
}

// logs in the user
document.getElementById("submitLog").addEventListener("click", (event) => {
    event.preventDefault();
    console.log(document.getElementById("exampleInputEmail1").value);
    console.log(document.getElementById("exampleInputPassword1").value);
    loginUser();
    // sessionStorage.setItem(auth, '1');
    window.location.href="index.html";
});

// creates json for new user
function createNewUser() {
    let newUser = {
        "Email": document.getElementById("inputEmail3").value,
        "Password": document.getElementById("inputPassword3").value,
        "PasswordConf": document.getElementById("inputPassword3C").value,
        "UserName": document.getElementById("username").value,
        "FirstName": document.getElementById("fname").value,
        "LastName": document.getElementById("lname").value
    };
    console.log(newUser);
    fetch(base + "/user",
        {
            method: "POST",
            body: JSON.stringify(newUser),
            headers: new Headers(
                {"Content-Type": "application/json",}
            )
        }
    ).then(response => {
        if (response.status == 405 || response.status == 400) {
            console.log("error creating new user account");
            console.log(response);
            return
        }
        // let token = [];
        // token = response.headers.get("Authorization").split(" ");
        // state.auth = token;
        
    })
}

// console.log(document.getElementById("submitNUser"))

document.getElementById("submitNUser").addEventListener("click", (event) => {
    event.preventDefault();
    createNewUser();
    exports.auth = state.auth;
    window.location.href="index.html";
})

// click event for creating the new meeting
// /meeting
document.getElementById("newMeet").addEventListener("click", function(event) {
    let newMeeting = {
        "MeetingName": document.getElementById("meetName").value,
        "MeetingDesc": document.getElementById("meetDes").value
    };
    console.log(newMeeting);
    fetch(base + "/meeting",
        {
            method: "POST",
            body: JSON.stringify(newMeeting),
            headers: newHeaders(
                {"Content-Type": "application/json",
                "Authorization": state.auth,
                }
            )
        }
    ).then(response => {
        if (response.status == 405 || response.status == 400) {
            console.log("Error creating new meeting");
            console.log(response);
            return
        }
        let token = [];
        token = response.headers.get("Content-Type");
        console.log(token);
        window.alert("New Meeting Created!");
    })
});
