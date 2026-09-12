[[ if eq (var "environment" .) "staging" ]]
job "htmx-go" {
  namespace   = [[ var "environment" . | quote ]]
  datacenters = ["homelab"]
  type        = "service"

  meta {
    image = [[ var "image" . | quote ]]
  }

  group "web" {
    count = 1

    update {
      max_parallel      = 1
      health_check      = "checks"
      min_healthy_time  = "10s"
      healthy_deadline  = "2m"
      progress_deadline = "5m"
      auto_revert       = true
    }

    restart {
      attempts = 3
      interval = "10m"
      delay    = "15s"
      mode     = "fail"
    }

    reschedule {
      attempts       = 3
      interval       = "1h"
      delay          = "30s"
      delay_function = "exponential"
      max_delay      = "5m"
      unlimited      = false
    }

    network {
      mode = "host"

      port "http" {
        to = 7883
      }
    }

    volume "data" {
      type            = "host"
      source          = "htmx-go-staging-data"
      attachment_mode = "file-system"
      access_mode     = "single-node-writer"
      sticky          = true
    }

    task "web" {
      driver = "docker"

      config {
        image        = [[ var "image" . | quote ]]
        network_mode = "services"
        ports        = ["http"]
      }

      env {
        DB_URL         = "/data/prod.db"
        ENCRYPTION_KEY = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        PORT           = "7883"
      }

      volume_mount {
        volume      = "data"
        destination = "/data"
      }

      service {
        name     = "htmx-go-staging"
        provider = "nomad"
        port     = "http"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.htmx-go-staging.entrypoints=websecure",
          "traefik.http.routers.htmx-go-staging.middlewares=htmx-go-staging-noindex",
          "traefik.http.routers.htmx-go-staging.rule=Host(`[[ var "domain" . ]]`)",
          "traefik.http.routers.htmx-go-staging.tls.domains[0].main=[[ var "domain" . ]]",
          "traefik.http.middlewares.htmx-go-staging-noindex.headers.customresponseheaders.X-Robots-Tag=noindex, nofollow",
        ]

        check {
          name     = "HTTP health"
          type     = "http"
          path     = "/api/health"
          interval = "10s"
          timeout  = "2s"

          check_restart {
            limit           = 3
            grace           = "30s"
            ignore_warnings = false
          }
        }
      }

      resources {
        cpu    = 500
        memory = 512
      }

      logs {
        max_files     = 5
        max_file_size = 10
      }

      kill_timeout = "30s"
    }
  }
}
[[ else ]]
job "htmx-go" {
  namespace   = [[ var "environment" . | quote ]]
  datacenters = ["homelab"]
  type        = "service"

  meta {
    image = [[ var "image" . | quote ]]
  }

  group "web" {
    count = 1

    update {
      max_parallel      = 1
      health_check      = "checks"
      min_healthy_time  = "10s"
      healthy_deadline  = "2m"
      progress_deadline = "5m"
      auto_revert       = true
    }

    restart {
      attempts = 3
      interval = "10m"
      delay    = "15s"
      mode     = "fail"
    }

    reschedule {
      attempts       = 3
      interval       = "1h"
      delay          = "30s"
      delay_function = "exponential"
      max_delay      = "5m"
      unlimited      = false
    }

    network {
      mode = "host"

      port "http" {
        to = 7883
      }
    }

    volume "data" {
      type            = "host"
      source          = "htmx-go-production-data"
      attachment_mode = "file-system"
      access_mode     = "single-node-writer"
      sticky          = true
    }

    task "backup" {
      lifecycle {
        hook    = "prestart"
        sidecar = false
      }

      driver = "docker"

      config {
        image        = [[ var "image" . | quote ]]
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

    task "web" {
      driver = "docker"

      config {
        image        = [[ var "image" . | quote ]]
        network_mode = "services"
        ports        = ["http"]
      }

      env {
        DB_URL         = "/data/prod.db"
        ENCRYPTION_KEY = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        PORT           = "7883"
      }

      volume_mount {
        volume      = "data"
        destination = "/data"
      }

      service {
        name     = "htmx-go-production"
        provider = "nomad"
        port     = "http"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.htmx-go-production.entrypoints=websecure",
          "traefik.http.routers.htmx-go-production.rule=Host(`[[ var "domain" . ]]`)",
        ]

        check {
          name     = "HTTP health"
          type     = "http"
          path     = "/api/health"
          interval = "10s"
          timeout  = "2s"

          check_restart {
            limit           = 3
            grace           = "30s"
            ignore_warnings = false
          }
        }
      }

      resources {
        cpu    = 500
        memory = 512
      }

      logs {
        max_files     = 5
        max_file_size = 10
      }

      kill_timeout = "30s"
    }
  }
}
[[ end ]]
