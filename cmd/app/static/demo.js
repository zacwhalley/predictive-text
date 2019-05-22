(function() { "use strict"
  let inputBox;
  let results;
  let resultsList;
  let errorMessage;
  let apiUrl;

  function init() {
    if (location.protocol !== "https:") {
      location.protocol = "https:";
    }
    
    inputBox = document.getElementById("textInput");
    results = document.getElementById("results");
    resultsList = document.getElementById("resultsList");
    errorMessage = document.getElementById("errorMessage");
    apiUrl = document.getElementById("apiUrl").innerText;

    results.style.display = "none";
    inputBox.value = "";
    
    inputBox.addEventListener("keyup", handleInputKeyup);
  }

  function handleInputKeyup(e) {
    if (!inputBox.value) {
      // user cleared input
      resetResults();
      return;
    }
    getPredictions();
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
            return;
          }
    
          // remove existing values from the list
          while (resultsList.hasChildNodes()) {
            resultsList.removeChild(resultsList.lastChild);
          }
    
          // add all new results to the list
          for (const prediction of body.predictions)  {
            if (prediction && prediction.trim()) {
              const newItem = createResult(body.input, prediction);
              resultsList.appendChild(newItem);
            }
          }  
          results.style.display = ""; // unhide results
          errorMessage.style.display = "none"; // hide error message
        });
      }
    }).catch(displayError)
  }

  function createResult(input, prediction) {
    let newLi = document.createElement("li");
    newLi.innerHTML = `
    <p>${input} <span class="predicted">${prediction}</span>
    `;
    return newLi;
  }

  function displayError(err) {
    results.style.display = "none";
    errorMessage.style.display = "";
    errorMessage.innerText = err.message;
  }

  function makePredictionRequest(input) {
    const requestUrl = `${apiUrl}/prediction?input=${input}`
    const request = new Request(requestUrl, {method: "GET"});

    return fetch(request);
  }

  function resetResults() {
    results.style.display = "none";
    errorMessage.style.display = "none";
  }

  init();
})()