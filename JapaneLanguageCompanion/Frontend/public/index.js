async function fetchKanaKanjiData() {
    try {
        const response = await fetch('http://127.0.0.1:3000');
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        const data = await response.json();
        console.log(data);  // This logs the JSON data to the console
        displayKanaKanjiData(data);  // Function to display the data on the front end
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function displayKanaKanjiData(data) {
    const kanaKanjiContainer = document.getElementById('kanaKanjiContainer');
    kanaKanjiContainer.innerHTML = '';  // Clear any existing content

    data.forEach(item => {
        const kanaKanjiElement = document.createElement('div');
        kanaKanjiElement.textContent = `${item.Character} - ${item.Romanization}`;
        kanaKanjiContainer.appendChild(kanaKanjiElement);
    });
}

// Call the function when the page loads or when needed
fetchKanaKanjiData();
