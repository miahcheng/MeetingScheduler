`use strict`;

const base = "https://api.jimhua32.me";
// click event for creating the new meeting
// /meeting

document.getElementById("newMeet").addEventListener("click", (event) => {
    event.preventDefault();
    let newMeeting = {
        "MeetingName": document.getElementById("meetName").value,
        "MeetingDesc": document.getElementById("meetDes").value
    };
    fetch(base + "/meeting",
        {
            method: "POST",
            body: JSON.stringify(newMeeting),
            headers:
                {"Content-Type": "application/json",
                "Authorization": sessionStorage.getItem("auth")
                }
        }
    ).then(response => {
        console.log(response)
        if (response.status >= 400) {
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