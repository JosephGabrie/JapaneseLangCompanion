<<<<<<< HEAD
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
=======
fetch(" http://127.0.0.1:3000/ ")
    .then(res => res.json())
    .then(data => {
        const arrayLength = data.length;
        
        for(let i=0; i< arrayLength; i++) {
            const characterData= data.find(item => item.kanakanji_id === i);
            if (characterData) {
                console.log(characterData.character);
                document.getElementById("japLetter").innerHTML = characterData.character
            }
            console.log(data[i]);

    
        }
        
    })
    .catch(error => console.error(error))
>>>>>>> origin/production
