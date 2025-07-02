package transformer

import "net/http"

type RequestTransformer func(req *http.Request)
