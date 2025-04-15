function handler(event) {
  var request = event.request;
  var uri = request.uri;
  // Strip off the "/{api-type}/{name}" part of the path
  request.uri = uri.replace(/^\/[^\/]+\/[^\/]+/, "");
  return request;
}
