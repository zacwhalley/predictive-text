(function() { "use strict"
  // elements on page
  let inputBox;
  let okButton;
  let results;
  let resultsList;
  let errorMessage;

  function init() {
    inputBox = document.getElementById("textInput");
    okButton = document.getElementById("predictButton");
    results = document.getElementById("results");
    resultsList = document.getElementById("resultsList");
    errorMessage = document.getElementById("errorMessage");

    results.style.display = "none"
    inputBox.value = ""
    okButton.addEventListener("click", getPredictions);
    inputBox.addEventListener('keyup', getPredictions);
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
    const request = new Request("http://localhost:8080/api/prediction", {
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