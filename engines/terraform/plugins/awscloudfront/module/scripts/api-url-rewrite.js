function handler(event) {
  var request = event.request;
  var uri = request.uri;
  // Strip off the first part of the path
  request.uri = uri.replace(/^\/[^\/]+/, "");
  return request;
}
