package web

const serviceWorkerJS = `const CDN = "https://static1.e621.net/";
const LOCAL = "/s1/";

self.addEventListener("install", () => self.skipWaiting());
self.addEventListener("activate", e => e.waitUntil(self.clients.claim()));

self.addEventListener("fetch", event => {
	const url = event.request.url;
	if (!url.startsWith(CDN))
		return;

	const headers = new Headers();
	const range = event.request.headers.get("range");
	if (range) {
		headers.set("Range", range);
	}

	event.respondWith(fetch(LOCAL + url.slice(CDN.length), { headers }));
});`