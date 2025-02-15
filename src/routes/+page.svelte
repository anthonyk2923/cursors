<script>
	import { onMount } from "svelte";

	const socket = new WebSocket("ws://localhost:8080");
	let curs = [];
	let mouseX = 1;
	let mouseY = 1;
	let divContainer;
	let rect = { left: 0, top: 0, width: 0, height: 0 };
	let sizeInd = 1;
	let scaleO = 1.4;
	let scaleW = 1.8;
	socket.addEventListener("open", () =>
		console.log("WebSocket connection opened"),
	);
	socket.addEventListener("error", (event) =>
		console.error("WebSocket error:", event),
	);
	socket.addEventListener("close", () =>
		console.log("WebSocket connection closed"),
	);

	socket.addEventListener("message", (event) => {
		console.log("Message from server:", event.data);
		const data = JSON.parse(event.data);
		if (Array.isArray(data)) {
			curs = data.map((cur) => ({
				...cur,
				adjustedX: (cur.x - rect.left) / scale,
				adjustedY: (cur.y - rect.top) / scale,
			}));
		} else {
			console.error("Received data is not an array:", data);
		}
	});

	const handleMouseMove = (event) => {
		if (!divContainer) return;
		rect = divContainer.getBoundingClientRect();

		mouseX = (event.clientX - rect.left) / scale;
		mouseY = (event.clientY - rect.top) / scale;

		if (mouseY % 10 <= 2 || mouseX % 10 <= 2) {
			if (
				mouseX >= 0 &&
				mouseX <= rect.width / scale &&
				mouseY >= 0 &&
				mouseY <= rect.height / scale
			) {
				sendData(event.clientX, event.clientY);
				console.log("done");
			}
		}

		console.log("Mouse moved to:", mouseX, mouseY);
	};

	function sendData(x, y) {
		const data = JSON.stringify({ x, y });
		console.log("Sending:", data);
		socket.send(data);
	}

	onMount(() => {
		if (divContainer) {
			rect = divContainer.getBoundingClientRect();
		}
	});

	$: scale = (2 / sizeInd) * scaleO;
</script>

<svelte:body on:mousemove={handleMouseMove} />

<div class="flex items-center justify-center h-screen cursor-none">
	<div
		bind:this={divContainer}
		class="bg-stone-950 relative"
		style="
		width: {96 * scale * scaleW}px;
		height: {96 * scale}px;
		transform: scale({scale});
	"
	>
		{#each curs as cur}
			<div
				class="absolute"
				style="left: {cur.adjustedX}px; top: {cur.adjustedY}px;"
			>
				<svg
					class="w-3 h-3"
					fill="none"
					stroke="rgb({cur.color.r}, {cur.color
						.g}, {cur.color.b})"
					stroke-width="5"
					version="1.1"
					xmlns="http://www.w3.org/2000/svg"
					xmlns:xlink="http://www.w3.org/1999/xlink"
					viewBox="0 0 203.079 203.079"
					xml:space="preserve"
				>
					<path
						d="M192.231,104.082V102c0-12.407-10.094-22.5-22.5-22.5c-2.802,0-5.484,0.519-7.961,1.459  C159.665,70.722,150.583,63,139.731,63c-2.947,0-5.76,0.575-8.341,1.61C128.667,55.162,119.624,48,109.231,48  c-2.798,0-5.496,0.541-8,1.516V22.5c0-12.407-10.094-22.5-22.5-22.5s-22.5,10.093-22.5,22.5v66.259  c-3.938-5.029-8.673-9.412-14.169-11.671c-6.133-2.52-12.587-2.219-18.667,0.872c-11.182,5.686-15.792,19.389-10.277,30.548  l27.95,56.563c0.79,1.552,19.731,38.008,54.023,38.008h40c31.54,0,57.199-25.794,57.199-57.506l-0.031-41.491H192.231z   M135.092,188.079h-40c-24.702,0-40.091-28.738-40.646-29.796l-27.88-56.42c-1.924-3.893-0.33-8.519,3.629-10.532  c2.182-1.11,4.081-1.223,6.158-0.372c8.281,3.395,16.41,19.756,19.586,29.265l2.41,7.259l12.883-4.559V22.5  c0-4.136,3.364-7.5,7.5-7.5s7.5,3.364,7.5,7.5V109h0.136h14.864h0.136V71c0-4.187,3.748-8,7.864-8c4.262,0,8,3.505,8,7.5v15v26h15  v-26c0-4.136,3.364-7.5,7.5-7.5s7.5,3.364,7.5,7.5V102v16.5h15V102c0-4.136,3.364-7.5,7.5-7.5s7.5,3.364,7.5,7.5v10.727h0.035  l0.025,32.852C177.291,169.014,158.36,188.079,135.092,188.079z"
					/>
				</svg>
			</div>
		{/each}
	</div>
</div>
