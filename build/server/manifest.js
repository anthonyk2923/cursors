const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {start:"_app/immutable/entry/start.Ce2KSYny.js",app:"_app/immutable/entry/app.BfFGOVZ-.js",imports:["_app/immutable/entry/start.Ce2KSYny.js","_app/immutable/chunks/BKlrPCi6.js","_app/immutable/chunks/DE0S399y.js","_app/immutable/chunks/k2L460yQ.js","_app/immutable/entry/app.BfFGOVZ-.js","_app/immutable/chunks/DE0S399y.js","_app/immutable/chunks/CuwieagD.js","_app/immutable/chunks/5xY_qTeW.js","_app/immutable/chunks/CQ37X0C3.js","_app/immutable/chunks/Co3oE8KU.js","_app/immutable/chunks/k2L460yQ.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./chunks/0-Cbzph7zM.js')),
			__memo(() => import('./chunks/1-cDUVLx6j.js')),
			__memo(() => import('./chunks/2-tIH0q-cI.js'))
		],
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();

const prerendered = new Set([]);

const base = "";

export { base, manifest, prerendered };
//# sourceMappingURL=manifest.js.map
