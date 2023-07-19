//Pour le double range slider   
let rangeMin = 0;
const range = document.querySelector(".range-selected");
const rangeInput = document.querySelectorAll(".range-input input");
const rangePrice = document.querySelectorAll(".range-price span");

rangeInput.forEach((input, index) => {
    input.addEventListener("input", (e) => {
        let minRange = parseInt(rangeInput[0].value);
        let maxRange = parseInt(rangeInput[1].value);

        if (maxRange - minRange < rangeMin) {
            if (e.target === rangeInput[0]) {
                rangeInput[1].value = minRange + rangeMin;
                maxRange = minRange + rangeMin;
            } else {
                rangeInput[0].value = maxRange - rangeMin;
                minRange = maxRange - rangeMin;
            }
        }

        rangePrice[0].textContent = minRange;
        rangePrice[1].textContent = maxRange;
        rangePrice[0].setAttribute("data-value", minRange);
        rangePrice[1].setAttribute("data-value", maxRange);
        range.style.left = ((minRange - rangeInput[0].min) / (rangeInput[0].max - rangeInput[0].min)) * 100 + "%";
        range.style.right = 100 - ((maxRange - rangeInput[1].min) / (rangeInput[1].max - rangeInput[1].min)) * 100 + "%";
        });
});


function toggleFields() {
    
    var creationCheckbox = document.querySelector('input[name="creationdate"]');
    var creationFields = document.querySelectorAll('input[name="datdebut"], input[name="datfin"]');

    var albumCheckbox = document.querySelector('input[name="firstalbum"]');
    var albumFields = document.querySelectorAll('input[name="debutalbum"], input[name="finalalbum"]');

    var membersCheckbox = document.querySelector('input[name="members"]');
    var membersField = document.querySelectorAll('input[id="member1"], input[id="member2"]');

    var locationCheckbox = document.querySelector('input[name="location"]');
    var locationField = document.querySelector('select[name="loc"]');


    // Basculer les champs pour creation date
    for (var i = 0; i < creationFields.length; i++) {
        var field = creationFields[i];
        field.disabled = !creationCheckbox.checked;
    }

    // Basculer les champs pour first album
    for (var i = 0; i < albumFields.length; i++) {
        var field = albumFields[i];
        field.disabled = !albumCheckbox.checked;
    }

    // Basculer les champs pour concert location
    locationField.disabled = !locationCheckbox.checked;

    // Basculer les champs pour members
    for (var i = 0; i < membersField.length; i++) {
        var field = membersField[i];
        field.disabled = !membersCheckbox.checked;
    }
    //    membersField.disabled = !membersCheckbox.checked;

}
var valeurInput1 = document.getElementById("member1");
var valeurAffichee1 = document.getElementById("valeurAffichee1");

valeurInput1.addEventListener("input", function(){
    valeurAffichee1.textContent = valeurInput1.value;
});

var valeurInput2 = document.getElementById("member2");
var valeurAffichee2 = document.getElementById("valeurAffichee2");
valeurInput2.addEventListener("input", function(){
    valeurAffichee2.textContent = valeurInput2.value;
});


// basculer vers filter

function toggleDiv() {
    var div = document.getElementById("filter");
    if (div.style.display === "block") {
        div.style.display = "none"; // Afficher l'élément
    } else {
        div.style.display = "block"; // Masquer l'élément
    }
}
const icon = document.querySelector('.icon');
const search = document.querySelector('.search');
icon.onclick = function(){
    search.classList.toggle('active')
}