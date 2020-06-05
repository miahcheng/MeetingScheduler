`use strict`;

const base = "https://api.jimhua32.me";
// click event for creating the new meeting
// /meeting
console.log(sessionStorage.getItem("auth"));
function createNewMeeting() {
    let newMeeting = {
        "MeetingName": document.getElementById("meetName").value,
        "MeetingDesc": document.getElementById("meetDes").value
    };
    console.log(newMeeting);
    console.log(sessionStorage.getItem("auth"));
    fetch(base + "/meeting",
        {
            method: "POST",
            body: JSON.stringify(newMeeting),
            headers: new Headers(
                {"Content-Type": "application/json",
                "Authorization": sessionStorage.getItem("auth"),
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
}

document.getElementById("newMeet").addEventListener("click", (event) => {
    event.preventDefault();
    createNewMeeting();
});