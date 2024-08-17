fetch("http://127.0.0.1:3000/")
  .then(res => res.json())
  .then(data => {
    const container = document.getElementById("characterInputs"); // Replace with your container ID

    data.forEach((character, index) => {
      const label = document.createElement("label");
      label.textContent = character.character;

      const input = document.createElement("input");
      input.type = "text";
      input.id = `input_${index}`; // Unique ID for each input

      const containerDiv = document.createElement("div");
      containerDiv.appendChild(label);
      containerDiv.appendChild(input);

      container.appendChild(containerDiv);
    });
  })
  .catch(error => console.error(error));