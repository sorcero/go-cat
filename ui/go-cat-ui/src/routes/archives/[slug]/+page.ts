import { error } from '@sveltejs/kit';
 
import type { PageLoad } from './$types';
export let prerender = true;
export const trailingSlash = 'always';

export const load = (async ({ params }) => {

  console.log(params.slug)
  if (params.slug == "[slug]") {
    prerender = false;
  }
  
  try {
    const imageModules = (await import(`./../../../archives/${params.slug}/infra.json`)).default //./../sorcero-releases/archives/.*/infra.json);

    
    console.log(imageModules)
    return imageModules
  } catch (e) {
    console.warn(e)
    throw error(404, 'Not found');  
  }
}) satisfies PageLoad;