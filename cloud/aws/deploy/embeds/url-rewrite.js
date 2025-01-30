function handler(event) {
  var request = event.request;
  var uri = request.uri;

  // If the URI has no extension and doesn't end with a trailing slash, append '/index.html'
  if (!uri.includes(".") && !uri.endsWith("/")) {
    request.uri = uri + "/index.html";
  }
  // If the URI ends with a trailing slash, replace it with '/index.html'
  else if (uri.endsWith("/")) {
    request.uri = uri.replace(/\/$/, "/index.html");
  }

  return request;
}
