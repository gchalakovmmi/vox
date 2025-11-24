/* conversation.js -------------------------------------------------- */
const infoP	 = document.getElementById('info');
const chatDiv = document.getElementById('chat');

/* state ----------------------------------------------------------- */
const STATE = {
	IDLE:		'idle',		// waiting for user
	PLAYING:	'playing',	// greeting / VOX audio is playing
	RECORDING:	'recording',	// mic on
	UPLOADING:	'uploading'	// waiting for back-end
};
let state = STATE.IDLE;
let mediaRecorder= null;
let audioChunks	= [];
let recordStart	= 0;			// performance.now()
let audioPlayer = null;			// current HTMLAudioElement

/* helpers --------------------------------------------------------- */
function addChat(sender, text){
	const p = document.createElement('p');
	p.textContent = `${sender}: ${text}`;
	chatDiv.appendChild(p);
	chatDiv.scrollTop = chatDiv.scrollHeight;
}
function setInfo(txt){ infoP.textContent = txt; }

/* audio playback -------------------------------------------------- */
async function playAudio(url){
	state = STATE.PLAYING;
	setInfo('VOX is speaking…');
	audioPlayer = new Audio(url);
	audioPlayer.play();
	return new Promise(res => {
		audioPlayer.onended	= () => res();
		audioPlayer.onerror	= () => res(); // treat error as ended
	});
}

/* initial greeting on first space --------------------------------- */
let firstGreetingDone = false;
document.addEventListener('keydown', async e => {
	if (e.code !== 'Space' || e.repeat) return;
	e.preventDefault();

	if (state === STATE.IDLE && !firstGreetingDone){
		firstGreetingDone = true;
		const rsp = await fetch('/conversation/greeting', {
			method: 'POST',
			headers: {'Content-Type':'application/json'},
			body: JSON.stringify({messages:[]})
		});
		if (!rsp.ok){ console.error(await rsp.text()); return; }
		const data = await rsp.json();
		addChat('VOX', data.text);
		await playAudio(data.audioUrl);
		state = STATE.IDLE;
		setInfo('Hold space to answer…');
		return;
	}

	if (state !== STATE.IDLE) return;		// ignore while playing/uploading
	startRecording();
});

document.addEventListener('keyup', e => {
	if (e.code !== 'Space') return;
	e.preventDefault();
	if (state === STATE.RECORDING) stopRecording();
});

/* recording ------------------------------------------------------- */
async function startRecording(){
	audioChunks = [];
	const stream = await navigator.mediaDevices.getUserMedia({audio:true});
	mediaRecorder = new MediaRecorder(stream);
	mediaRecorder.ondataavailable = e => audioChunks.push(e.data);
	mediaRecorder.onstop = () => stream.getTracks().forEach(t => t.stop());
	mediaRecorder.start();
	recordStart = performance.now();
	state = STATE.RECORDING;
	setInfo('Recording… release to stop');
}

function stopRecording(){
	mediaRecorder.stop();
	state = STATE.UPLOADING;
	setInfo('Processing…');
	mediaRecorder.onstop = async () => {
		const blob = new Blob(audioChunks, {type:'audio/webm'});
		const duration = (performance.now() - recordStart)/1000;

		if (duration < 1){
			setInfo('Too short – hold the key while you speak');
			state = STATE.IDLE;
			return;
		}
		await uploadAudio(blob);
	};
}

/* upload ---------------------------------------------------------- */
async function uploadAudio(blob){
	const fd = new FormData();
	fd.append('audio', blob, 'recording.webm');
	const rsp = await fetch('/conversation/prompt', {
		method: 'POST',
		headers: {'Authorization':'Bearer '+getAccessToken()},
		body: fd
	});
	if (!rsp.ok){
		const msg = await rsp.text();
		setInfo(msg);
		state = STATE.IDLE;
		return;
	}
	const data = await rsp.json();
	addChat('You', data.user);
	addChat('VOX', data.vox);
	await playAudio(data.audioUrl);
	state = STATE.IDLE;
	setInfo('Hold space to answer…');
}

/* tiny helper to read the access-token cookie ------------------- */
function getAccessToken(){
	return document.cookie
		.split('; ')
		.find(row => row.startsWith('access_token='))
		?.split('=')[1] || '';
}
