import type { RequestHandler } from '@sveltejs/kit';

const chainEndpoint = new URL('/', process.env.BACKEND_URL).href;

export const get: RequestHandler = async ({ url: { searchParams } }) => {
  const url = new URL(chainEndpoint);
  const input = searchParams.get('input');
  if (input) {
    url.searchParams.set('in', input);
  }
  const response = await fetch(url.toString()).then((res) => res.text());

  return {
    body: response,
  };
};
