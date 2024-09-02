const term = document.querySelector('#term');
const definition = document.querySelector('#definition');
const checkButton = document.querySelector('#checkButton');
const nextButton = document.querySelector('#nextButton');

let hiraganaData = []; // This will hold the fetched data

fetch("japanese.json")
    .then(response => response.json())
    .then(data => {
        hiraganaData = data.hiragana; // Save the data into the hiraganaData array
        getRandomWord(); // Initial call to set a random word on page load
    })
    .catch(error => console.error('Error:', error));

function getRandomWord() {
    if (hiraganaData.length === 0) {
        return; // Exit if no data is available
    }

    const randomItem = hiraganaData[Math.floor(Math.random() * hiraganaData.length)];
    term.innerHTML = `<h3>${randomItem.character}</h3>`;
    definition.innerHTML = `<h3>${randomItem.romanji}</h3>`;
    definition.style.display = 'none'; // Hide definition when changing term
}

checkButton.addEventListener('click', function() {
    definition.style.display = definition.style.display === 'block' ? 'none' : 'block';
});

nextButton.addEventListener('click', function() {
    getRandomWord();
});
