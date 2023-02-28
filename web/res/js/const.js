const onWindowLoadMessage = "Добро пожаловать в проект курсовой работы по базам данных."
const preTimeString = "В Москве сейчас"

function getMonthName(month) {
    switch (month) {
        case 1: return "января"
        case 2: return "февраля"
        case 3: return "марта"
        case 4: return "апреля"
        case 5: return "мая"
        case 6: return "июня"
        case 7: return "июля"
        case 8: return "августа"
        case 9: return "сентября"
        case 10: return "октября"
        case 11: return "ноября"
        case 12: return "декабря"
        default: return month.toString()
    }
}