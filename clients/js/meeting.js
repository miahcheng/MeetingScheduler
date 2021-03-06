let state = {
  meetings: new Map(),
  // key would be meetingID, holds struct:
  // {creatorID, members, meeting desc}
  toDisplay: null
};
const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
let testing = {
  meetingID : 1,
  creatorID : 3,
  members : [1,2,3],
  meetingtitle: "stuff to do",
  meetingdesc : "doing stuff",
  freeTime: days.map(function(someValue) {
    var tmp = {};
    tmp[someValue] = ["100", "200", "400"];
    return tmp;
  })
}
let testing1 = {
  meetingID : 2,
  creatorID : 3,
  members : [1,2,3],
  meetingtitle: "stuff to do2",
  meetingdesc : "doing stuff",
  freeTime: days.map(function(someValue) {
    var tmp = {};
    tmp[someValue] = ["100", "200", "400"];
    return tmp;
  })
}
function setState(){
  //Get user
  //Get meetingID
  state.meetings.set(testing.meetingID, {creatorID:testing.creatorID, members:testing.members, meetingtitle:testing.meetingtitle, meetingdesc:testing.meetingdesc, freeTime: testing.freeTime});
  state.meetings.set(testing1.meetingID, {creatorID:testing1.creatorID, members:testing1.members, meetingtitle:testing1.meetingtitle, meetingdesc:testing1.meetingdesc, freeTime: testing1.freeTime});
}
function sendState(toSend){
  //Post meetingID
  testing.members = [1,2,3,4];
}
function renderMeetingList(){
  document.getElementById("submitcon").style.display = "none";
  setState();
  let container = document.querySelector("#content");
  state.meetings.forEach(function(meeting){
    let row = document.createElement('div')
    row.classList.add("row");
    let b4 = document.createElement("button");
    b4.setAttribute("type", "button");
    b4.classList.add("btn", "btn-light", "btn-lg", "btn-block");
    b4.id = meeting.meetingID;
    b4.innerHTML = meeting.meetingtitle;
    b4.addEventListener('click', () => {
      setState();
      state.toDisplay = parseInt(b4.id);
      let container = document.querySelector("#content");
      container.innerHTML = "";
      renderOneMeeting();
    });
    row.appendChild(b4)
    container.appendChild(row);
  });
}
function renderOneMeeting(){
  document.getElementById("submitcon").style.display = "block";
  renderTitleDescUsers();
  renderTimePopUp();
  renderUserPopUp();
}
function renderTitleDescUsers(){
  let firstRow = document.createElement('div')
  firstRow.classList.add("row");
  let firstCol = document.createElement('div')
  firstCol.classList.add("col", "text-left");
  let secondCol = document.createElement('div')
  secondCol.classList.add("col", "float-right");
  let container = document.querySelector("#content");
  let header = document.createElement('h1');
  header.innerHTML = "Event: " + state.meetings.get(1).meetingtitle;
  let inp = document.createElement("input");
  inp.classList.add("btn", "btn-secondary", "btn-lg", "btn-block", "pull-right");
  inp.setAttribute("type", "submit");
  inp.setAttribute("value", "Add users to meeting");
  inp.setAttribute("data-toggle","modal")
  inp.setAttribute("data-target","#basicExampleModal1")
  firstCol.appendChild(header);
  secondCol.appendChild(inp);
  firstRow.appendChild(firstCol);
  firstRow.appendChild(secondCol)
  container.appendChild(firstRow);
  let desc = document.createElement('div');
  desc.classList.add("row");
  desc.innerHTML = "Description: " + state.meetings.get(1).meetingdesc;
  container.appendChild(desc);
  let usrt = document.createElement('h4');
  usrt.classList.add("row");
  usrt.innerHTML = "Users:"
  container.appendChild(usrt);
  let usr = document.createElement('div');
  usr.classList.add("row");
  usr.innerHTML = "";
  state.meetings.get(1).members.forEach(function(user){
    usr.innerHTML = usr.innerHTML + parseInt(user) + ", ";
  });
  usr.innerHTML = usr.innerHTML.substring(0, usr.innerHTML.length - 2 );
  container.appendChild(usr);

}
function renderUserPopUp(){
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
  b1.setAttribute("aria-label","Close");
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
  f2.setAttribute("aria-describedby","emailHelp");
  f2.setAttribute("placeholder","hi@gmail.com, bye@gmail.com");
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
    sendState(toParse);
    setState(state.toDisplay);
    let container = document.querySelector("#content");
    container.innerHTML = "";
    renderOneMeeting(state.toDisplay);
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
  m1.appendChild(m2)
  let cont = document.querySelector("#submitcon");
  cont.appendChild(m1);
}
function renderTimePopUp(){
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
  b1.setAttribute("aria-label","Close");
  let b2 = document.createElement("span");
  b2.setAttribute("aria-hidden", "true");
  b2.innerHTML = "&times";
  let m6 = document.createElement("div");
  m6.classList.add("modal-body");
  let parent = document.createElement("div")
  parent.classList.add("col")
  days.forEach(function(day, i){
    let toAdd = document.createElement("div")
    toAdd.classList.add("row")
    toAdd.innerHTML = JSON.stringify(state.meetings.get(1).freeTime[i]);
    toAdd.innerHTML = toAdd.innerHTML.replace(/[{}""[\]]/g, "")
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
  m1.appendChild(m2)
  let cont = document.querySelector("#submitcon");
  cont.appendChild(m1);
}
//renderOneMeeting();
renderMeetingList();
