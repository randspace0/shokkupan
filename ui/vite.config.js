import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

const config = {
	plugins: [tailwindcss(), sveltekit()],
	clearScreen: false,
	logLevel: "info",
	server: {
		proxy: {
			"/api/exec/socket": {
				target: "ws://localhost:5330",
				changeOrigin: true,
				ws: true,
			},
			"/api": {
				target: "http://localhost:5330",
				changeOrigin: true,
			}
		}
	},

}

export default config;