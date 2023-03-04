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


const urlStart = '/test/start';
const urlNext = '/test/next';
const urlAssert = '/test/assert';
const urlStop = '/test/stop';

const startHTML = `<div class="textBlock">
  <h1>Пройдите наш лучший тест на IQ!</h1>
  <div>Необходимо будет решить 2 задачи на логику. Максимальный балл: 150</div>
</div>
<div class="info"></div>
<div id="buttonBlock">
    <button class="" onclick="startTest()">Начать тест</button>
</div>`
