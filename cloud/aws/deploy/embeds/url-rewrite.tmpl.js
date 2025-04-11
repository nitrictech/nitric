function handler(event) {
  var request = event.request;
  var uri = request.uri;
  var basePath = "{{.BasePath}}";

  if (uri === basePath || uri.startsWith(basePath + "/")) {
    request.uri = uri.replace(new RegExp("^" + basePath + "[/]*"), "/");
  }

  return request;
}
