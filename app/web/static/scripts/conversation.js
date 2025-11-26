/* global STATE, conversationID, getAccessToken */
const infoP  = document.getElementById('info');
const chatDiv= document.getElementById('chat');

const STATE = {
    IDLE:       'idle',
    PLAYING:    'playing',
    RECORDING:  'recording',
    UPLOADING:  'uploading'
};
let state = STATE.IDLE;
let mediaRecorder = null;
let audioChunks   = [];
let recordStart   = 0;
let audioPlayer   = null;
let conversationID = -1;

/* ---------- helpers ---------- */
function addChat(sender, text){
    const wrap = document.createElement('div');
    wrap.className = sender === 'You'
        ? 'd-flex justify-content-end mb-2'
        : 'd-flex justify-content-start mb-2';

    const bubble = document.createElement('div');
    bubble.className = sender === 'You'
        ? 'bg-secondary bg-opacity-25 text-light px-3 py-2 rounded-start-4 rounded-end-0'
        : 'bg-secondary bg-opacity-25 text-light px-3 py-2 rounded-end-4 rounded-start-0';
    bubble.innerHTML = `<strong>${sender}:</strong> ${text}`;
    wrap.appendChild(bubble);
    chatDiv.appendChild(wrap);
    chatDiv.scrollTop = chatDiv.scrollHeight;
}

function setInfo(txt){ infoP.textContent = txt; }

/* ---------- prefetch greeting + mic permission ---------- */
let prefetched = false;
async function prefetch(){
    if (prefetched) return;
    prefetched = true;

    /* 1. request mic (and instantly release) */
    try {
        const temp = await navigator.mediaDevices.getUserMedia({audio:true});
        temp.getTracks().forEach(t => t.stop());
    } catch (e) {
        console.warn('Mic permission denied:', e);
    }

    /* 2. fetch greeting */
    fetch('/conversation/greeting', {
        method:'POST',
        headers:{'Content-Type':'application/json'},
        body: JSON.stringify({messages:[]})
    })
    .then(r => r.json())
    .then(data => {
        addChat('Voxy', data.text);
        playAudio(data.audioUrl);
    });
}

/* ---------- audio playback ---------- */
async function playAudio(url){
    state = STATE.PLAYING;
    setInfo('Voxy is speaking…');
    audioPlayer = new Audio(url);
    audioPlayer.play();
    await new Promise(res => {
        audioPlayer.onended = () => res();
        audioPlayer.onerror = () => res();
    });
    /* ─── NEW: make sure UI is unlocked ─── */
    state = STATE.IDLE;
    setInfo('Hold space to answer…');
}

/* ---------- space handler ---------- */
let firstGreetingDone = false;
document.addEventListener('keydown', async e => {
    if (e.code !== 'Space' || e.repeat) return;
    e.preventDefault();

    /* very first space: prefetch greeting + mic permission */
    if (!firstGreetingDone){
        firstGreetingDone = true;
        prefetch();               // starts both mic + greeting fetch
        setInfo('Hold space to answer…');
        return;
    }

    if (state !== STATE.IDLE) return;
    startRecording();
});

document.addEventListener('keyup', e => {
    if (e.code !== 'Space') return;
    e.preventDefault();
    if (state === STATE.RECORDING) stopRecording();
});

/* ---------- recording ---------- */
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

/* ---------- upload ---------- */
async function uploadAudio(blob){
    const fd = new FormData();
    fd.append('audio', blob, 'recording.webm');
    const rsp = await fetch('/conversation/prompt?conversationID='+conversationID, {
        method:'POST',
        headers:{'Authorization':'Bearer '+getAccessToken()},
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
    /* skip correction when it is exactly "Correct" */
    if (data.correction.trim() !== 'Correct') addChat('Correction', data.correction);
    addChat('Voxy', data.vox);

    if (conversationID === -1) {
        conversationID = data.conversationID;
        const endBtn = document.getElementById('end');
        endBtn.classList.remove('d-none');
        endBtn.href = '/conversation/analysis?id='+conversationID;
    }
    await playAudio(data.audioUrl);
    state = STATE.IDLE;
    setInfo('Hold space to answer…');
}

/* ---------- token helper ---------- */
function getAccessToken(){
    return document.cookie
        .split('; ')
        .find(row => row.startsWith('access_token='))
        ?.split('=')[1] || '';
}
