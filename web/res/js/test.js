function openTest(){
    let layer = document.createElement("div");
    layer.id = "layer"
    let test = document.createElement("div");
    test.id = "centralBlock";
    // let answerBox = document.createElement("input");
    // answerBox.id = "answerBox";
    // let acceptButton = document.createElement("button")
    // acceptButton.id = "acceptButton";
    // let closeButton = document.createElement("button")
    // closeButton.id = "closeButton";


    
    layer.appendChild(test);
    // test.appendChild(answerBox);
    // test.appendChild(closeButton);
    // test.appendChild(acceptButton);
    document.body.appendChild(layer);

    test.innerHTML = startHTML
}

let sessionID = 0
let lastQuestionID = 0

async function startTest() {
    let testBlock = document.getElementById('centralBlock')
    let info = testBlock.getElementsByClassName('info')[0]
    let response = await fetch(urlStart, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
    });
    if (!response.ok) {
        info.innerText = await response.text()
    }
    sessionID = await response.text();
    document.cookie = "session-id=" + sessionID
    response = await fetch(urlNext, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `sessionID=${sessionID}`
    });

    info = testBlock.getElementsByClassName('info')[0]
    if (!response.ok) {
        info.innerText = await response.text()
        return
    }
    let result = await response.json();
    testBlock.innerHTML = getTestHTML(result.id, result.question, result.image)
    lastQuestionID = result.id
}

async function next() {
    let answer = document.getElementById('answer')
    let testBlock = document.getElementById('centralBlock')
    let response = await fetch(urlAssert, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `sessionID=${sessionID}&problemID=${lastQuestionID}&answer=${answer.value}`
    });
    let info = testBlock.getElementsByClassName('info')[0]
    if (!response.ok) {
        info.innerText = await response.text()
        return
    }


    response = await fetch(urlNext, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `sessionID=${sessionID}`
    });

    info = testBlock.getElementsByClassName('info')[0]
    if (!response.ok) {
        info.innerText = await response.text()
        return
    }
    if (response.status === 204) {
        stopTest()
    }
    let result = await response.json();
    testBlock.innerHTML = getTestHTML(result.id, result.question, result.image)
    lastQuestionID = result.id
}

async function stopTest() {
    let testBlock = document.getElementById('centralBlock')
    response = await fetch(urlStop, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `sessionID=${sessionID}`
    });

    if (!response.ok) {
        testBlock.innerText = await response.text()
        return
    }
    let result = await response.json();
    testBlock.innerHTML = getScoreHTML(result.correct, result.all)
}

async function closeLayer() {
    let layer = document.getElementById('layer')
    layer.remove()
}

function getTestHTML(num, ques, src) {
    return `<div class="textBlock">
    <h1>Вопрос ${num}</h1>
    <div class="questionText">${ques}</div>
    <div class="info"></div>
    <img src="${src}" alt="Без изображения">
</div>
<input id="answer" placeholder="Ваш ответ" oninput="updateActive()">
<div id="buttonBlock">
    <button class="testButtonInactive" onClick="next()">Ответить</button>
    <button class="closeButton" onclick="closeLayer()">Закрыть</button>
</div>`
}

function getScoreHTML(cor, all) {
    let score = cor/all*150
    return `<div class="textBlock">
    <h1>Ваш результат: ${score}</h1>
    <div>Неплохо!</div>
</div>
<div id="buttonBlock">
    <button class="closeButton" onclick="closeLayer()">Закрыть</button>
</div>`
}

function updateActive(){
    let ans = document.getElementById("answer")
    let button = document.getElementById("buttonBlock").children[0]
    if (ans.value.length > 0){
        button.className = "testButton";
    } else {
        button.className = "testButtonInactive";
    }
}