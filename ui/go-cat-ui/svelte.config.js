import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	// pre-render for each archive page 
	// https://kit.svelte.dev/docs#ssr-and-javascript-ssr
	ssr: false,

	kit: {

		prerender: {
			crawl: true,
		},

		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		
		// strict: false // https://kit.svelte.dev/docs#configuration-strict
		adapter: adapter({
			strict: false,

			// directory
			pages: 'public',
		}),
	}
};

export default config;
