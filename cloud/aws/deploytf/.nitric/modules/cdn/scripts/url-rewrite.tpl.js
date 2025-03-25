function handler(event) {
  var request = event.request;
  var uri = request.uri;
  var basePath = "${base_path}";

  if (uri.startsWith(basePath) || uri.startsWith(basePath + "/")) {
    request.uri = uri.replace(new RegExp("^" + basePath + "[/]*"), "/");
  }

  return request;
}
