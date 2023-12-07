pipeline "http_step" {
  step "http" "send_to_slack" {
    url                = "https://myapi.com/vi/api/do-something"
    method             = "post"
    insecure           = true
    ca_cert_pem        = "test"

    request_body = jsonencode({
      name = "turbie"
      app  = "flowpipe"
    })

    request_headers = {
      Accept     = "application/json"
      User-Agent = "flowpipe" // check - is this the syntax with dash in a key name???
    }
  }
}
