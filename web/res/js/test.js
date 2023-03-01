function openTest(){
    let layer = document.createElement("div");
    layer.id = "layer"
    let test = document.createElement("div");
    test.id = "centralBlock";
    let answerBox = document.createElement("input");
    answerBox.id = "answerBox";
    let acceptButton = document.createElement("button")
    acceptButton.id = "acceptButton";
    let closeButton = document.createElement("button")
    closeButton.id = "closeButton";


    
    layer.appendChild(test);
    test.appendChild(answerBox);
    test.appendChild(closeButton);
    test.appendChild(acceptButton);
    document.body.appendChild(layer);
}