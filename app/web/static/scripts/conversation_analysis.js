document.addEventListener("DOMContentLoaded", () => {
	const params = new URLSearchParams(window.location.search);
	const id = params.get("id");
	if (!id) return;

	const feedbackDiv = document.getElementById("feedback");
	feedbackDiv.textContent = "Generating feedback…";

	fetch(`/conversation/analysis/feedback?id=${id}`, {
		headers: { Authorization: "Bearer " + getAccessToken() }
	})
	.then(r => {
		if (!r.ok) throw new Error(r.statusText);
		return r.text();
	})
	.then(text => feedbackDiv.textContent = text)
	.catch(err => feedbackDiv.textContent = "Error: " + err.message);
});

function getAccessToken() {
	return document.cookie
		.split("; ")
		.find(row => row.startsWith("access_token="))
		?.split("=")[1] ?? "";
}
