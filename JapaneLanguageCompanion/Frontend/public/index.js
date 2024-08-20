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
