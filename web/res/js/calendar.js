function showCalendar() {
    let table = document.getElementById("jsCalendar")
    let [data, today] = generateCalendar()
    table.innerHTML = `
    <tr><td>ПН</td><td>ВТ</td><td>СР</td><td>ЧТ</td><td>ПТ</td><td>СБ</td><td>ВС</td></tr>
    `
    for (let week in data) {
        let tr = document.createElement("tr")
        for (let day in data[week]) {
            let td = document.createElement("td")
            if (data[week][day] === today) {
                td.id = "todayElement"
            }
            td.innerText = data[week][day]
            tr.appendChild(td)
        }
        table.appendChild(tr)
    }
}

function generateCalendar() {
    let now = new Date()
    let [month, day] = [now.getMonth(), now.getDate()]
    let weeks = []
    let iterDate = now
    iterDate.setDate(1)
    while (true) {
        let week = []
        for (let i = 0; i < 7; i++) {
            let weekDay = iterDate.getDay() > 0 ? iterDate.getDay()-1 : 6
            if (iterDate.getMonth() !== month || weekDay !== i) {
                week.push("")
            } else {
                week.push(iterDate.getDate())
                iterDate.setDate(iterDate.getDate() + 1)
            }
        }
        weeks.push(week)
        if (iterDate.getMonth() !== month) {
            break
        }
    }
    return [weeks, day]
}