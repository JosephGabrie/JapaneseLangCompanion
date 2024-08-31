const term = document.querySelector('#term');
const definition = document.querySelector('#definition');
const checkButton = document.querySelector('#checkButton');
const nextButton = document.querySelector('#nextButton');

const words = {
    Hello: "hola",
    Bye: "Adios",
    "love you": "te amo"
};

const data = Object.entries(words);

function getRandomWord() {
    const randomTerm = data[Math.floor(Math.random() * data.length)];
    term.innerHTML = `<h3>${randomTerm[0]}</h3>`;
    definition.innerHTML = `<h3>${randomTerm[1]}</h3>`;
    definition.style.display = 'none'; // Hide definition when changing term
}

checkButton.addEventListener('click', function() {
    definition.style.display = definition.style.display === 'block' ? 'none' : 'block';
});

nextButton.addEventListener('click', function() {
    getRandomWord();
});

// Initial call to set a random word on page load
getRandomWord();
