function onloadPage() {
    // window.alert(onWindowLoadMessage)
    initClock()
    showCalendar()
    animateBanner().then()
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

async function animateBanner() {
    let banner = document.getElementById("banner")
    let width = banner.offsetWidth
    let height = banner.offsetHeight
    let Vx = 1
    let Vy = 1
    let timer = setInterval(function (){
        let screenWidth = document.documentElement.clientWidth
        let screenHeight = document.documentElement.clientHeight
        let top = banner.offsetTop
        let left = banner.offsetLeft
        let right = left + width
        let bottom = top + height

        if (right >= screenWidth || left <= 0) {
            Vx *= -1
        }
        if (bottom >= screenHeight || top <= 0) {
            Vy *= -1
        }

        banner.style.left = (left + Vx).toString() + 'px'
        banner.style.top = (top + Vy).toString() + 'px'
    }, 10)
}

onloadPage()