function handler(event) {
  var request = event.request;
  var uri = request.uri;

  // Match any '/' that occurs at the end of a URI. Replace it with a default index
  request.uri = uri.replace(/\/$/, "/index.html");

  return request;
}
