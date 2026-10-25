import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	kit: {
		adapter: adapter({
			pages: "dist",
			fallback: "app.html"
		}),
		alias: {
			$src: "./src",
			$components: "./src/components",
			$common: "./src/components/common",
		},
	},
	preprocess: vitePreprocess(),
};

export default config;