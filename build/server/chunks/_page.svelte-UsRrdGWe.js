import { p as push, e as ensure_array_like, d as stringify, c as pop } from './index-ChYQ0XmM.js';
import { e as escape_html } from './escaping-CqgfEcN3.js';

const replacements = {
  translate: /* @__PURE__ */ new Map([
    [true, "yes"],
    [false, "no"]
  ])
};
function attr(name, value, is_boolean = false) {
  if (value == null || !value && is_boolean || value === "" && name === "class") return "";
  const normalized = name in replacements && replacements[name].get(value) || value;
  const assignment = is_boolean ? "" : `="${escape_html(normalized, true)}"`;
  return ` ${name}${assignment}`;
}
function _page($$payload, $$props) {
  push();
  const socket = new WebSocket("wss://cursorsserver.onrender.com/ws");
  let curs = [];
  let bullets = [];
  let socketId;
  socket.addEventListener("open", () => console.log("WS connected"));
  socket.addEventListener("error", console.error);
  socket.addEventListener("close", () => console.log("WS closed"));
  socket.addEventListener("message", (event) => {
    try {
      const data = JSON.parse(event.data);
      if (typeof data === "object" && data.type && data.payload) {
        switch (data.type) {
          case "p":
            if (Array.isArray(data.payload)) {
              curs = data.payload;
            }
            break;
          case "b":
            handleBullet(data.payload);
            break;
          case "i":
            if (!socketId) {
              const idData = JSON.parse(event.data);
              socketId = idData.payload.id;
              console.log("Assigned socket ID:", socketId);
            }
            break;
          default:
            console.warn("Unknown message type:", data.type);
        }
      }
    } catch (error) {
      console.log("Raw event data:", event.data);
    }
  });
  function handleBullet(bulletData) {
    bullets = [...bullets, bulletData];
    console.log("Before timeout:", bullets);
    setTimeout(
      () => {
        bullets = bullets.filter((b) => b !== bulletData);
        console.log("After timeout:", bullets);
      },
      500
    );
  }
  const each_array = ensure_array_like(bullets);
  const each_array_1 = ensure_array_like(curs);
  $$payload.out += `<div class="h-screen w-screen cursor-none select-none"><div class="bg-stone-950 relative w-full h-full pointer-events-none">`;
  {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]--> <!--[-->`;
  for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
    let bul = each_array[$$index];
    $$payload.out += `<svg class="absolute inset-0 w-full h-full pointer-events-none select-none"><line${attr("x1", bul.fromPoint.x)}${attr("y1", bul.fromPoint.y)}${attr("x2", bul.isHit.x)}${attr("y2", bul.isHit.y)}${attr("stroke", `rgb(${bul.fromPoint.color.r}, ${bul.fromPoint.color.g}, ${bul.fromPoint.color.b}`)} stroke-width="2" stroke-linecap="square" stroke-dasharray="2,20" stroke-opacity="0.5"></line></svg>`;
  }
  $$payload.out += `<!--]--> <!--[-->`;
  for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
    let cur = each_array_1[$$index_1];
    $$payload.out += `<div class="absolute pointer-events-none select-none"${attr("style", `left: ${stringify(cur.x)}px; top: ${stringify(cur.y)}px;`)}><svg class="w-4 h-4 pointer-events-none select-none" fill="none"${attr("stroke", `rgb(${cur.color.r}, ${cur.color.g}, ${cur.color.b}`)} stroke-width="2" viewBox="0 0 24 24"><path d="M12 2C8.134 2 5 5.134 5 9c0 5 7 13 7 13s7-8 7-13c0-3.866-3.134-7-7-7z"></path><circle cx="12" cy="9" r="2.5"></circle></svg></div>`;
  }
  $$payload.out += `<!--]--></div></div>`;
  pop();
}

export { _page as default };
//# sourceMappingURL=_page.svelte-UsRrdGWe.js.map
