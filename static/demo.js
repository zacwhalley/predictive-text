(function() { "use strict"
  let inputBox;
  let predictBtn;
  let results;
  let resultsList;
  let errorMessage;
  let apiUrl;

  function init() {
    inputBox = document.getElementById("textInput");
    results = document.getElementById("results");
    resultsList = document.getElementById("resultsList");
    predictBtn = document.getElementById("predictBtn");
    errorMessage = document.getElementById("errorMessage");
    apiUrl = document.getElementById("apiUrl").innerText;

    results.style.display = "none"
    inputBox.value = ""
    
    //inputBox.addEventListener("keyup", handleInputKeyup);
    predictBtn.addEventListener("click", getPredictions);
  }

  function handleInputKeyup(e) {
    const spaceCode = 32;
    const enterCode = 18;

    if (!inputBox.value) {
      // user cleared input
      resetResults();
      return
    }

    if (e.keyCode === spaceCode || e.keyCode === enterCode) {
      getPredictions();
    }
  }
  function getPredictions() {
    // get text from input box
    const input = inputBox.value.trim();
    if (!input) {
      resetResults();
      return;
    }

    // mock call to get predictions
    makePredictionRequest(input).then(response => {
      if (response.status === 200) {
        response.json().then(body => {
          // display predictions in result
          if (!body || !body.predictions || body.predictions.length < 1) {
            return
          }
    
          // remove existing values from the list
          while (resultsList.hasChildNodes()) {
            resultsList.removeChild(resultsList.lastChild);
          }
    
          // add all new results to the list
          for (const prediction of body.predictions)  {
            let newLi = document.createElement("li")
            newLi.innerText = body.input + " " + prediction
            resultsList.appendChild(newLi)
          }  
          results.style.display = ""; // unhide results
          errorMessage.style.display = "none"; // hide error message
        });
      }
    }).catch(displayError)
  }

  function displayError(err) {
    results.style.display = "none";
    errorMessage.style.display = "";
    errorMessage.innerText = err.message
  }

  function makePredictionRequest(input) {
    const request = new Request(apiUrl + "/prediction", {
      method: "POST",
      body: JSON.stringify({input: input}),
    });

    return fetch(request)
  }

  function resetResults() {
    results.style.display = "none";
    errorMessage.style.display = "none";
  }

  init();
})()