function fetchAllRecords() {
    return function () {
        fetch('/api/tts/records')
            .then(response => response.json())
            .then(records => {
                const list = document.getElementById('ttsRecordList');
                records.forEach(record => {
                    const listItem = document.createElement('li');
                    listItem.textContent = record.Title;
                    listItem.dataset.recordId = record.Id;
                    listItem.onclick = () => fetchAndPlayRecord(record);
                    list.appendChild(listItem);
                });
            });
    };
}

window.onload = fetchAllRecords();

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
            const url = URL.createObjectURL(blob);
            const audioPlayer = document.getElementById('audioPlayer');
            audioPlayer.src = url;
            audioPlayer.play();
            fetchAllRecords()();
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}