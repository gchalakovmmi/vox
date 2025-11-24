const p = document.getElementById('info');
const chatDiv = document.getElementById('chat');   // <--- history container
let isFirstSpacePress = true;
let isSpacePressed = false;

document.addEventListener('keydown', async function(event) {
	if (event.code === 'Space' && !isSpacePressed) {
		event.preventDefault();
		isSpacePressed = true;

		if (isFirstSpacePress) {
			isFirstSpacePress = false;
			p.textContent = 'VOX Speaking...';

			const rsp = await fetch('/conversation/reply', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({messages: []})
			});
			if (!rsp.ok) { console.error(await rsp.text()); return; }
			const data = await rsp.json();

			// add text to history
			const para = document.createElement('p');
			para.textContent = 'VOX: ' + data.text;
			chatDiv.appendChild(para);

			// play audio
			const audioRsp = await fetch(data.audioUrl);
			if (!audioRsp.ok) { console.error(await audioRsp.text()); return; }
			const blob = await audioRsp.blob();
			const url = URL.createObjectURL(blob);
			const audio = new Audio(url);
			audio.play();
			audio.addEventListener('ended', () => URL.revokeObjectURL(url), {once: true});
		} else {
			p.textContent = 'Listening... Release to stop';
			// TODO: start recording
		}
	}
});

document.addEventListener('keyup', function(event) {
	if (event.code === 'Space') {
		event.preventDefault();
		isSpacePressed = false;
		p.textContent = 'Press and hold the spacebar to start speaking';
		// TODO: stop recording
		// TODO: send audio + history to /conversation/reply
		// TODO: add user/VOX paragraphs to chatDiv
	}
});

window.addEventListener('keydown', e => {
	if (e.code === 'Space' && e.target === document.body) e.preventDefault();
});
