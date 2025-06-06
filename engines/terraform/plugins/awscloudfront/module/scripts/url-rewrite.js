// Template in the base paths from the manifest
// A comma separated list of base paths
var basePaths = "${base_paths}"
var allBasePaths = basePaths.split(",").sort((a, b) => b.length - a.length);

function handler(event) {
    var request = event.request;
    var uri = request.uri;

    if (basePath !== "/") {
        for (var i = 0; i < allBasePaths.length; i++) {
            if (allBasePaths[i] === "/") {
                return request.uri;
            }

            var basePath = allBasePaths[i];
            if (uri === basePath || uri.startsWith(basePath + "/")) {
                request.headers["x-nitric-cache-key"] = { value: basePath };
                request.uri = uri.replace(new RegExp("^" + basePath + "[/]*"), "/");
                return request;
            }
        }
    }

    return request;
}
  