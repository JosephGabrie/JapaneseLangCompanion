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

    
    // Do some reasearch on thisDocument Object Model 
const term= document.querySelector('.term');
const checkButton = document.querySelector('.check');
const definition = document.querySelector('.definition')
const nextButton = document.querySelector('.next')

words = {
    Hello: "Hola",
    Goodbye: "Adios",
    "I drink water": "Yo tomo agua"
 }

data = Object.entries(words)

function getRandomWord(){
    randomTerm=data[Math.floor(Math.random() * data.length)]
    term.innerHTML = `<h3>${randomTerm[0]}</h3>`;
    definition.innerHTML = ` <h3>${randomTerm[1]}</h3`
}

data = Object.entries(words)
console.log(data[0][0])

checkButton.addEventListener('click',function(){
    definition.style.display = 'block';
});
nextButton.addEventListener('click',function(){
    getRandomWord();
});


    