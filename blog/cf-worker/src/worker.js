function unauthorized() {
  return new Response("Authentication required", {
    status: 401,
    headers: {
      "WWW-Authenticate": 'Basic realm="Audos SDK build log"',
      "Cache-Control": "private, no-store",
    },
  });
}

export default {
  async fetch(request, env) {
    const authHeader = request.headers.get("authorization");
    if (authHeader) {
      const [scheme, encoded] = authHeader.split(" ");
      if (scheme === "Basic" && encoded) {
        const decoded = atob(encoded);
        const idx = decoded.indexOf(":");
        const user = decoded.slice(0, idx);
        const pass = decoded.slice(idx + 1);
        if (user === env.BASIC_AUTH_USER && pass === env.BASIC_AUTH_PASS) {
          const assetResponse = await env.ASSETS.fetch(
            new Request(request.url, { method: request.method, headers: {} })
          );
          const response = new Response(assetResponse.body, assetResponse);
          response.headers.set("Cache-Control", "private, no-store");
          return response;
        }
      }
    }
    return unauthorized();
  },
};
