    task "backup" {
      lifecycle {
        hook    = "prestart"
        sidecar = false
      }

      driver = "docker"

      config {
        image        = "ghcr.io/sachahjkl/htmx-go@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
        command      = "/bin/sh"
        args         = ["-ec", "test -s /data/prod.db; mkdir -p /data/backups; archive=/data/backups/pre-deploy-$NOMAD_ALLOC_ID.sqlite; sqlite3 /data/prod.db \".backup '$archive'\"; test -s $archive; gzip $archive"]
        network_mode = "services"
      }

      volume_mount {
        volume      = "data"
        destination = "/data"
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
