// data to render
const times = [ "12:00 - 12:30", "12:30 - 1:00", "1:00 - 1:30", "1:30 - 2:00", "2:00 - 2:30", "2:30 - 3:00", "3:00 - 3:30", "3:30 - 4:00", "4:00 - 4:30", "4:30 - 5:00", "5:00 - 5:30",
"5:30 - 6:00", "6:00 - 6:30 ,","6:30 - 7:00","7:00 - 7:30" ,"7:30 - 8:00" ,"8:00 - 8:30" ,"8:30 - 9:00", "9:00 - 9:30","9:30 - 10:00" ,"10:00 - 10:30" ,"10:30 - 11:00" ,"11:00 - 11:30" ,"11:30 - 12:00","12:00 - 12:30" ,"12:30 - 1:00"];
const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
const base = "https://api.blah.com";
let state = {
  selected: new Map(),
  auth: ""
};
function renderOneDay(day) {
  let dayArea = document.createElement("div");
  dayArea.classList.add("day");
  let nameofday = document.createElement("h5");
  nameofday.innerHTML = day;
  dayArea.appendChild(nameofday);
  times.forEach(function (time) {
    dayArea.appendChild(renderOneTime(time, "AM", day));
  });
  times.forEach(function (time) {
    dayArea.appendChild(renderOneTime(time, "PM", day));
  });
  return dayArea;
}
function renderOneTime(time, amorpm, day) {
  let oneTime = document.createElement("div");
  oneTime.classList.add("form-group", "form-check");
  let check = document.createElement("input");
  check.type = "checkbox";
  check.classList.add("form-check-input");
  let parsedTime = time.split(" ");;
  parsedTime = parsedTime[0].replace(":", "");
  if (amorpm === "PM") {
    parsedTime = parseInt(parsedTime) + 1200;
    parsedTime = parsedTime.toString();
  }
  check.id = day + ":" + parsedTime;
  if (state.selected.get(day).includes(parsedTime)){
    check.checked = true;
  }
  oneTime.appendChild(check);
  text = document.createElement("label");
  text.classList.add("form-check-label");
  text.for = "exampleCheck1";
  text.innerHTML = time + " " + amorpm;
  oneTime.appendChild(text);
  return oneTime;
}
function renderOneWeek() {
  let weekArea = document.querySelector("#schedule");
  days.forEach(function(day) {
    weekArea.appendChild(renderOneDay(day));
  });
  return weekArea;
}
function newMap(){
  days.forEach(function(day){
    state.selected.set(day, []);
  });
}
function setState(){
  state.auth = sessionStorage.getItem('auth');
  console.log(state.auth);
  if (state.selected.size === 0){
    newMap();
  }
  //GET USER, set state.selected to GET USER JSON
  fetch(base + "/user/",
      {
          method: "GET",
          body: "",
          headers: new Headers({
              "Authentication":state.auth,
          })
      }
  ).then(response => {
      let user = JSON.parse(response);
      days.forEach(function(day){
        user.set(day, user.${day})
      });
  })
}
function sendState(){
  if (state.selected.size === 0){
    newMap();
  }
  //GET USER, set state.selected to GET USER JSON
  fetch(base + "/user/",
      {
          method: "PATCH",
          body: JSON.stringify(state.selected),
          headers: new Headers({
              "Content-Type": "application/json",
              "Authentication":state.auth,
          })
      }
  ).then(response => {
      if (response.status == 405 || response.status == 405) {
          console.log("error logging in user");
          console.log(response);
          return;
      }
      let token = [];
      token = response.headers.get("Authorization").split(" ");
      console.log(token);
      state.auth = token[0];
  })
  console.log(state.selected);
}
// call the function to render chats
setState();
renderOneWeek();
let button = document.querySelector('#submit');
button.addEventListener('click', () => {
  newMap()
  document.querySelectorAll(".form-check-input").forEach(function(checkmark){
    if (checkmark.checked === true) {
      let parsed = checkmark.id.split(":");
      var arr = [];
      arr = arr.concat(state.selected.get(parsed[0]));
      arr.push(parsed[1]);
      state.selected.set(parsed[0], arr);
    }
  });
    sendState();

});
