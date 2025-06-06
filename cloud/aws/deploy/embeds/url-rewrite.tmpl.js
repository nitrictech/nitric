function handler(event) {
  var request = event.request;
  var uri = request.uri;
  var basePath = "{{.BasePath}}";

  // First apply base path removal if needed
  if (
    basePath !== "/" &&
    (uri === basePath || uri.startsWith(basePath + "/"))
  ) {
    // Inject a custom header to isolate cache keys
    request.headers["x-nitric-cache-key"] = { value: basePath };

    uri = uri.replace(new RegExp("^" + basePath + "[/]*"), "/");
  }

  // Then append index.html to the uri if it is a directory
  if (!uri.includes(".")) {
    // TODO inject root document value instead of hardcoding
    uri = uri.endsWith("/") ? uri + "index.html" : uri + "/index.html";
  }

  request.uri = uri;

  return request;
}
