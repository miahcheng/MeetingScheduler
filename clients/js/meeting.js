let state = {
  meetings: new Map(),
  // key would be meetingID, holds struct:
  // {creatorID, members, meeting desc}
  toDisplay: null,
  toDisplayUsers: new Map(),
  add: []
};
const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
const base = "https://api.jimhua32.me";
function getuser(id, callback){
  let i = 0;
  state.meetings.get(id).Members.forEach(function(memid){
    fetch(base + "/getuser/" + parseInt(memid),
      {
          method: "GET",
          headers: {
              "Authorization": sessionStorage.getItem("auth"),
          }
      }
    ).then(response => {
      if (response.status == 401 || response.status == 405) {
        console.log("Error getting meeting information");
      }
      return response.json();
    }).then(response => {
        i = i + 1;
        console.log(state.meetings.get(id).Members.length)
        console.log(i);
        response.id = memid;
        state.toDisplayUsers.set(memid, response);
        console.log(state.toDisplayUsers);
        if (i === state.meetings.get(id).Members.length){
          callback();
        }
      })
    })
};
function setState(callback) {
  fetch(base + "/user/",
      {
          method: "GET",
          headers: {
              "Authorization": sessionStorage.getItem("auth"),
          }
      }
  ).then(response => {
    if (response.status == 400 || response.status == 405 || response.status == 401) {
      console.log("Error getting user information");
      console.log(response);
    }
    return response.json();
  }).then(response => {
    let user = response;
    let i = 0
    user.Meetings.forEach(function (id) {
    fetch(base + "/meeting/" + parseInt(id),
        {
            method: "GET",
            headers: {
                "Authorization": sessionStorage.getItem("auth"),
            }
        }
    ).then(response => {
      if (response.status == 401 || response.status == 405) {
        console.log("Error getting meeting information");
      }
      return response.json();
    }).then(response => {
        i = i+1;
        response.id = id;
        state.meetings.set(id, response);
        console.log(state.meetings);
        if (i === user.Meetings.length) {
          callback();
        }
      })
    });
  });
}
function sendState() {
  state.add.forEach(function(email) {
    console.log(email);
    console.log(state.toDisplay);
    var obj = {Email: email};
    fetch(base + "/meeting/" + parseInt(state.toDisplay),
        {
          method: "POST",
          body: JSON.stringify(obj),
          headers: {
            "Content-Type": "application/json",
            "Authorization": sessionStorage.getItem("auth"),
          }
        }
    ).then(response => {
        if (response.status === 415 || response.status === 401 || response.status === 405) {
          console.log("error adding user to meeting");
          console.log(response);
        }
        setState(function(){
          let container = document.querySelector("#content");
          container.innerHTML = "";
          renderOneMeeting(state.toDisplay);
        });
    })
  })
}

function renderMeetingList() {
  state.toDisplay = 0;
  let container = document.querySelector("#content");
  container.innerHTML = "";
  setState(function() {
    document.getElementById("submitcon").style.display = "none";
    console.log(state.meetings);
    let container = document.querySelector("#content");
    state.meetings.forEach(function (meeting) {
      console.log("hello");
      let row = document.createElement('div')
      row.classList.add("row");
      let b4 = document.createElement("button");
      b4.setAttribute("type", "button");
      b4.classList.add("btn", "btn-light", "btn-lg", "btn-block");
      b4.id = meeting.id;
      b4.innerHTML = meeting['MeetingName'];
      b4.addEventListener('click', () => {
        setState(function(){
          state.toDisplay = parseInt(b4.id);
          let container = document.querySelector("#content");
          container.innerHTML = "";
          renderOneMeeting();
        });
      });
      row.appendChild(b4)
      container.appendChild(row);
    });
  });
}

function renderOneMeeting() {
  document.getElementById("submitcon").style.display = "block";
  getuser(state.toDisplay, function(){
    renderTitleDescUsers();
    renderTimePopUp();
    renderUserPopUp();
  });
}

function renderTitleDescUsers() {
  let firstRow = document.createElement('div')
  firstRow.classList.add("row");
  let firstCol = document.createElement('div')
  firstCol.classList.add("col", "text-left");
  let secondCol = document.createElement('div')
  secondCol.classList.add("col", "float-right");
  let container = document.querySelector("#content");
  let header = document.createElement('h1');
  header.innerHTML = "Event: " + state.meetings.get(state.toDisplay)["MeetingName"];
  let inp = document.createElement("input");
  inp.classList.add("btn", "btn-secondary", "btn-lg", "btn-block", "pull-right");
  inp.setAttribute("type", "submit");
  inp.setAttribute("value", "Add users to meeting");
  inp.setAttribute("data-toggle", "modal");
  inp.setAttribute("data-target", "#basicExampleModal1");
  inp.addEventListener('click', () => {
    let container = document.querySelector("#content");
    container.innerHTML = "";
    setState(renderOneMeeting());
  })
  firstCol.appendChild(header);
  secondCol.appendChild(inp);
  firstRow.appendChild(firstCol);
  firstRow.appendChild(secondCol)
  container.appendChild(firstRow);
  let desc = document.createElement('h4');
  desc.innerHTML = "Description: " + state.meetings.get(state.toDisplay)["MeetingDesc"];
  firstCol.appendChild(desc);
  let usrt = document.createElement('h2');
  usrt.innerHTML = "Users:"
  firstCol.appendChild(usrt)
  inp = document.createElement("input");
  inp.classList.add("btn", "btn-secondary", "btn-sm", "btn-block", "pull-right");
  inp.setAttribute("type", "submit");
  inp.setAttribute("value", "Delete this Meeting");
  inp.addEventListener('click', () => {
    fetch(base + "/meeting/" + parseInt(state.toDisplay),
        {
            method: "DELETE",
            headers: {
                "Authorization": sessionStorage.getItem("auth"),
            }
        }
    ).then(response => {
      if (response.status == 400 || response.status == 405 || response.status == 401) {
        console.log("Error getting user information");
        console.log(response);
      }
      console.log(response);
      console.log("ok");
      state.meetings.delete(state.toDisplay);
      renderMeetingList();
    })
  });
  secondCol.appendChild(inp);
  let usr = document.createElement('h4');
  usr.innerHTML = "";
  state.toDisplayUsers.forEach(function (user) {
    usr.innerHTML = usr.innerHTML + user.FirstName + " " + user.LastName + ", ";
  });
  usr.innerHTML = usr.innerHTML.substring(0, usr.innerHTML.length - 2);
  firstCol.appendChild(usr)
  firstRow.appendChild(firstCol);
  firstRow.appendChild(secondCol)
  container.appendChild(firstRow);
}

function renderUserPopUp() {
  let m1 = document.createElement("div");
  m1.classList.add("modal", "fade");
  m1.id = "basicExampleModal1";
  m1.setAttribute("tabindex", "-1");
  m1.setAttribute("role", "dialog");
  m1.setAttribute("aria-labelledby", "exampleModalLabel");
  m1.setAttribute("aria-hidden", "true");
  let m2 = document.createElement("div");
  m2.classList.add("modal-dialog");
  m2.setAttribute("role", "document");
  let m3 = document.createElement("div");
  m3.classList.add("modal-content");
  let m4 = document.createElement("div");
  m4.classList.add("modal-header");
  let m5 = document.createElement("h5");
  m5.classList.add("modal-title");
  m5.id = "exampleModalLabel";
  m5.innerHTML = "Add User Emails below to add";
  let b1 = document.createElement("button");
  b1.setAttribute("type", "button");
  b1.classList.add("close")
  b1.setAttribute("data-dismiss", "modal");
  b1.setAttribute("aria-label", "Close");
  let b2 = document.createElement("span");
  b2.setAttribute("aria-hidden", "true");
  b2.innerHTML = "&times";
  let m6 = document.createElement("div");
  m6.classList.add("modal-body");
  let parent = document.createElement("div")
  parent.classList.add("col")
  let fillOut = document.createElement("form");
  let f1 = document.createElement("label");
  f1.setAttribute("for", "exampleemail");
  let f2 = document.createElement("input");
  f2.setAttribute("type", "email");
  f2.classList.add('form-control');
  f2.id = "exampleemail";
  f2.setAttribute("aria-describedby", "emailHelp");
  f2.setAttribute("placeholder", "hi@gmail.com, bye@gmail.com");
  fillOut.appendChild(f1);
  fillOut.appendChild(f2);
  parent.appendChild(fillOut)
  m6.appendChild(parent);
  let m7 = document.createElement("div");
  m7.classList.add("modal-footer");
  let b3 = document.createElement("button");
  b3.setAttribute("type", "button");
  b3.classList.add("btn", "btn-secondary");
  b3.setAttribute("data-dismiss", "modal");
  b3.innerHTML = "Close";
  let b4 = document.createElement("button");
  b4.setAttribute("type", "button");
  b4.classList.add("btn", "btn-primary", "mr-auto");
  b4.setAttribute("data-dismiss", "modal");
  b4.innerHTML = "Add user";
  b4.addEventListener('click', () => {
    let toParse = document.querySelector(".form-control").value
    toParse.replace(" ", "");
    toParse = toParse.split(",")
    console.log(toParse);
    state.add = toParse;
    sendState();
  });
  m7.appendChild(b4);
  m7.appendChild(b3);
  b1.appendChild(b2);
  m4.appendChild(m5);
  m4.appendChild(b1);
  m3.appendChild(m4);
  m3.appendChild(m6);
  m3.appendChild(m7);
  m2.appendChild(m3);
  m1.appendChild(m2);
  let cont = document.querySelector("#submitcon");
  cont.appendChild(m1);
}

function renderTimePopUp() {
  let m1 = document.createElement("div");
  m1.classList.add("modal", "fade");
  m1.id = "basicExampleModal";
  m1.setAttribute("tabindex", "-1");
  m1.setAttribute("role", "dialog");
  m1.setAttribute("aria-labelledby", "exampleModalLabel");
  m1.setAttribute("aria-hidden", "true");
  let m2 = document.createElement("div");
  m2.classList.add("modal-dialog");
  m2.setAttribute("role", "document");
  let m3 = document.createElement("div");
  m3.classList.add("modal-content");
  let m4 = document.createElement("div");
  m4.classList.add("modal-header");
  let m5 = document.createElement("h5");
  m5.classList.add("modal-title");
  m5.id = "exampleModalLabel";
  m5.innerHTML = "Free Times";
  let b1 = document.createElement("button");
  b1.setAttribute("type", "button");
  b1.classList.add("close")
  b1.setAttribute("data-dismiss", "modal");
  b1.setAttribute("aria-label", "Close");
  let b2 = document.createElement("span");
  b2.setAttribute("aria-hidden", "true");
  b2.innerHTML = "&times";
  let m6 = document.createElement("div");
  m6.classList.add("modal-body");
  let parent = document.createElement("div")
  parent.classList.add("col")
  days.forEach(function (day) {
    let toAdd = document.createElement("div")
    toAdd.classList.add("row")
    toAdd.innerHTML = JSON.stringify(state.meetings.get(state.toDisplay)[day]);
    toAdd.innerHTML = toAdd.innerHTML.replace(/[{}""[\]]/g, "");
    toAdd.innerHTML = day + ":" + toAdd.innerHTML;
    parent.appendChild(toAdd);
  });
  m6.appendChild(parent);
  let m7 = document.createElement("div");
  m7.classList.add("modal-footer");
  let b3 = document.createElement("button");
  b3.setAttribute("type", "button");
  b3.classList.add("btn", "btn-secondary");
  b3.setAttribute("data-dismiss", "modal");
  b3.innerHTML = "Close";
  m7.appendChild(b3);
  b1.appendChild(b2);
  m4.appendChild(m5);
  m4.appendChild(b1);
  m3.appendChild(m4);
  m3.appendChild(m6);
  m3.appendChild(m7);
  m2.appendChild(m3);
  m1.appendChild(m2);
  let cont = document.querySelector("#submitcon");
  cont.appendChild(m1);
}
renderMeetingList();
