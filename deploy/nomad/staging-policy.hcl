namespace "staging" {
  capabilities = ["list-jobs", "parse-job", "read-job", "submit-job"]
}

host_volume "htmx-go-staging-data" {
  policy = "write"
}
