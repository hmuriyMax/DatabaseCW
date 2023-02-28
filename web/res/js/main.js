function onloadPage() {
    // window.alert(onWindowLoadMessage)
    initClock()
    showCalendar()
}

function initClock() {
    document.getElementById("clock-text").innerText = preTimeString
    updateClock()
    let timer = setInterval(
        updateClock, 1000
    )
}

function updateClock() {
    let clock = document.getElementById("clockBlock")
    if (clock == null) {
        console.error("no #clockBlock found")
        return
    }
    let now = new Date()
    let [hour, minute, second] = [now.getHours(), now.getMinutes(), now.getSeconds()]
    let [year, month, day] = [now.getFullYear(), now.getMonth()+1, now.getDate()]
    if(day<10) day = "0"+day;
    // if(month<10) month = "0"+month;
    if(hour<10) hour = "0"+hour;
    if(minute<10) minute = "0"+minute;
    if(second<10) second = "0"+second;
    if (hour === "00" && minute === "00" && second === "00") {
        showCalendar()
    }
    clock.getElementsByClassName("time")[0].innerText = hour + ":" + minute + ":" + second
    clock.getElementsByClassName("date")[0].innerText = day + " " + getMonthName(month) + " " + year
}

onloadPage()