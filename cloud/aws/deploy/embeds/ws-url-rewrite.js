function handler(event) {
    var request = event.request;
  
    request.uri = "/ws";
  
    return request;
  }