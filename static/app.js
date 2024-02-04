window.onload = fetchAllRecords();

function fetchAllRecords() {
    return function () {
        fetch('/api/tts/records')
            .then(response => response.json())
            .then(records => {
                const list = document.getElementById('ttsRecordList');
                list.innerHTML = '';
                records.forEach(record => {
                    const listItem = document.createElement('li');
                    const deleteButton = document.createElement('span');

                    listItem.textContent = record.Title;
                    listItem.dataset.recordId = record.Id;
                    listItem.onclick = () => fetchAndPlayRecord(record);

                    deleteButton.textContent = ' x';
                    deleteButton.style.cursor = 'pointer';
                    deleteButton.style.color = 'red';
                    deleteButton.style.float = 'right';
                    deleteButton.style.paddingLeft = '10px';
                    deleteButton.style.paddingRight = '10px';
                    deleteButton.onclick = (event) => {
                        event.stopPropagation(); // Prevent triggering listItem's onclick
                        deleteRecord(record.Id);
                    };

                    listItem.appendChild(deleteButton);

                    list.appendChild(listItem);
                });
            });
    };
}

function deleteRecord(recordId) {
    fetch(`/api/tts/records/${recordId}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            fetchAllRecords()();
        })
        .catch(error => console.error('Error deleting record:', error));
}

function fetchAndPlayRecord(record) {
    document.getElementById('titleInput').value = record.Title;
    document.getElementById('textInput').value = record.Text;
    fetch(`/api/tts/records/${record.Id}/audio`)
        .then(response => response.blob())
        .then(blob => {
            // Set the blob as the source of the audio player and play it
            const audioPlayer = document.getElementById('audioPlayer');
            audioPlayer.src = URL.createObjectURL(blob);
            audioPlayer.play();
        })
        .catch(error => console.error('Error fetching audio:', error));
}

function submitText() {
    const inputAreaContainer = document.getElementById('inputAreaContainer');
    const progressBarContainer = document.getElementById('progressBarContainer');

    inputAreaContainer.style.display = 'none';
    progressBarContainer.style.display = 'block';
    progressBarContainer.style.display = 'flex';
    progressBarContainer.justifyContent = 'center';
    progressBarContainer.alignItems = 'center';

    fetch('/api/tts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            text: document.getElementById('textInput').value,
            title: document.getElementById('titleInput').value
        })
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.blob();
        })
        .then(blob => {
            progressBarContainer.style.display = 'none';
            inputAreaContainer.style.display = 'block';
            const url = URL.createObjectURL(blob);
            const audioPlayer = document.getElementById('audioPlayer');
            audioPlayer.src = url;
            audioPlayer.play();
            fetchAllRecords()();
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
            progressBarContainer.style.display = 'none';
            inputAreaContainer.style.display = 'block';
        });
}

document.addEventListener('DOMContentLoaded', function() {
    const textInput = document.getElementById('textInput');
    function resizeTextarea() {
        this.style.height = 'auto';
        this.style.height = this.scrollHeight + 'px';
    }
    textInput.addEventListener('input', resizeTextarea);
    resizeTextarea.call(textInput);
});
