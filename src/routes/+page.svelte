<script>
	import { onMount } from "svelte";

	const socket = new WebSocket("ws://localhost:8080");
	//const socket = new WebSocket("ws://192.168.1.3:8080");
	let curs = [];
	let bullets = [];
	let mouseX = 0;
	let mouseY = 0;
	let divContainer;
	let drawLine = false;
	let linePoint = { x: 0, y: 0 };
	let socketId;
	socket.addEventListener("open", () => console.log("WS connected"));
	socket.addEventListener("error", console.error);
	socket.addEventListener("close", () => console.log("WS closed"));
	socket.addEventListener("message", (event) => {
		try {
			const data = JSON.parse(event.data);
			if (
				typeof data === "object" &&
				data.type &&
				data.payload
			) {
				switch (data.type) {
					case "p":
						if (
							Array.isArray(
								data.payload,
							)
						) {
							curs = data.payload;
						}
						break;

					case "b":
						handleBullet(data.payload);
						break;

					case "i":
						if (!socketId) {
							const idData =
								JSON.parse(
									event.data,
								);
							socketId =
								idData.payload
									.id;
							console.log(
								"Assigned socket ID:",
								socketId,
							);
						}
						break;

					default:
						console.warn(
							"Unknown message type:",
							data.type,
						);
				}
			}
		} catch (error) {
			console.log("Raw event data:", event.data);
		}
	});
	function handleBullet(bulletData) {
		bullets = [...bullets, bulletData];
		console.log("Before timeout:", bullets);

		setTimeout(() => {
			bullets = bullets.filter((b) => b !== bulletData);
			console.log("After timeout:", bullets);
		}, 500);
	}
	const handleMouseMove = (event) => {
		if (!divContainer) return;
		const rect = divContainer.getBoundingClientRect();

		mouseX = event.clientX - rect.left;
		mouseY = event.clientY - rect.top;

		if (event.buttons === 1) {
			drawLine = true;
			for (let i = 0; i < curs.length; i++) {
				if (curs[i].user_id == socketId) {
					linePoint = {
						x: curs[i].x,
						y: curs[i].y,
					};
				}
			}
		} else {
			if (drawLine) {
				sendBullet(mouseX, mouseY, linePoint);
			}
			drawLine = false;
			linePoint = {};
			if (
				mouseX >= 0 &&
				mouseX <= rect.width &&
				mouseY >= 0 &&
				mouseY <= rect.height
			) {
				sendPosition();
			}
		}
	};

	const sendPosition = () => {
		socket.send(
			JSON.stringify({
				type: "p",
				payload: { x: mouseX, y: mouseY },
			}),
		);
	};
	const sendBullet = (mouseX, mouseY, linePoint) => {
		// Calc direction vector (dx, dy)
		const dx = mouseX - linePoint.x;
		const dy = mouseY - linePoint.y;

		// Normalize direction vector
		const magnitude = Math.sqrt(dx * dx + dy * dy);
		const normalizedDx = dx / magnitude;
		const normalizedDy = dy / magnitude;

		socket.send(
			JSON.stringify({
				type: "b",
				payload: { dx: normalizedDx, dy: normalizedDy },
			}),
		);
	};
</script>

<svelte:body on:mousemove={handleMouseMove} />

<div class="h-screen w-screen cursor-none select-none">
	<div
		bind:this={divContainer}
		class="bg-stone-950 relative w-full h-full pointer-events-none"
	>
		{#if drawLine}
			<svg
				class="absolute inset-0 w-full h-full pointer-events-none select-none"
			>
				<line
					x1={linePoint.x}
					y1={linePoint.y}
					x2={mouseX}
					y2={mouseY}
					stroke="white"
					stroke-width="2"
					stroke-dasharray="8,5"
					stroke-opacity="0.5"
				/>
			</svg>
		{/if}
		{#each bullets as bul}
			<svg
				class="absolute inset-0 w-full h-full pointer-events-none select-none"
			>
				<line
					x1={bul.fromPoint.x}
					y1={bul.fromPoint.y}
					x2={bul.isHit.x}
					y2={bul.isHit.y}
					stroke={`rgb(${bul.fromPoint.color.r}, ${bul.fromPoint.color.g}, ${bul.fromPoint.color.b}`}
					stroke-width="2"
					stroke-linecap="square"
					stroke-dasharray="2,20"
					stroke-opacity="0.5"
				/>
			</svg>
		{/each}
		{#each curs as cur}
			<div
				class="absolute pointer-events-none select-none"
				style="left: {cur.x}px; top: {cur.y}px;"
			>
				<svg
					class="w-4 h-4 pointer-events-none select-none"
					fill="none"
					stroke={`rgb(${cur.color.r}, ${cur.color.g}, ${cur.color.b}`}
					stroke-width="2"
					viewBox="0 0 24 24"
				>
					<path
						d="M12 2C8.134 2 5 5.134 5 9c0 5 7 13 7 13s7-8 7-13c0-3.866-3.134-7-7-7z"
					/>
					<circle cx="12" cy="9" r="2.5" />
				</svg>
			</div>
		{/each}
	</div>
</div>
