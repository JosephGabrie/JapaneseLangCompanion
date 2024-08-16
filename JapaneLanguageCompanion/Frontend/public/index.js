fetch(" http://127.0.0.1:3000/ ")
    .then(res => res.json())
    .then(data => {
        const characterData = data.find(item => item.kanakanji_id === 6);
        if (characterData) {
            console.log(characterData.character);
            document.getElementById("japLetter").innerHTML = characterData.character
        } else {
            console.log("Character not found");
        }
    })


    .catch(error => console.error(error));

    