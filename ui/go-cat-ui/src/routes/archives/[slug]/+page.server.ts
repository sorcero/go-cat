import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
 



import type { PageServerLoad } from './$types';
 

export const load = (async ({ params }) => {
  console.log(params.slug)
  try {
    const imageModules = (await import(`./../../../archives/${params.slug}/infra.json`)).default //./../sorcero-releases/archives/.*/infra.json);


    console.log(imageModules)
    return imageModules
  } catch (e) {
    console.warn(e)
    throw error(404, 'Not found');  
  }
}) satisfies PageServerLoad;