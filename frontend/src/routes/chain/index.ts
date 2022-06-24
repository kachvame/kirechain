const chainEndpoint = new URL('/', process.env.BACKEND_URL).href;

export const get = async (input?: string) => {
    const url = new URL(chainEndpoint)
    if (input) {
        url.searchParams.set("in", input)
    }

    const response = await fetch(url.toString()).then((res) => res.text());

    return {
        body: response,
    };
};
