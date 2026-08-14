{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
    bun2nix = {
      url = "github:nix-community/bun2nix?ref=2.1.2";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      bun2nix,
    }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" ] (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ bun2nix.overlays.default ];
        };
        pname = "htmx-go";
        version = "0.0.1";
        vendorHash = "sha256-T5v3g8iQsvI840x0JEpUoYJHp4ZTOe+1Kxj5PQiAS8k=";
        src = pkgs.lib.cleanSourceWith {
          src = ./.;
          filter =
            path: _type:
            !builtins.elem (baseNameOf path) [
              ".git"
              ".jj"
              "result"
            ];
        };
        bunDeps = pkgs.bun2nix.fetchBunDeps { bunNix = ./bun.nix; };
        css = pkgs.bun2nix.mkDerivation {
          pname = "${pname}-css";
          inherit version src bunDeps;
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [ pkgs.stdenv.cc.cc.lib ];
          buildPhase = ''
            runHook preBuild
            bun ./node_modules/@tailwindcss/cli/dist/index.mjs -i ./style/style.css -o style.css
            runHook postBuild
          '';
          installPhase = ''
            runHook preInstall
            install -Dm644 style.css $out/style.css
            runHook postInstall
          '';
        };
        app = pkgs.buildGoModule {
          inherit pname version src;
          inherit vendorHash;
          subPackages = [ "cmd" ];
          postInstall = ''
            mv $out/bin/cmd $out/bin/${pname}
            mkdir -p $out/share/${pname}
            cp -r assets views $out/share/${pname}/
            cp ${css}/style.css $out/share/${pname}/assets/style.css
          '';
          meta.mainProgram = pname;
        };
        mkGoCheck =
          name: command:
          pkgs.buildGoModule {
            pname = "${pname}-${name}";
            inherit version src;
            inherit vendorHash;
            buildPhase = command;
            doCheck = false;
            installPhase = "touch $out";
          };
        gofmt = pkgs.runCommand "${pname}-gofmt" { nativeBuildInputs = [ pkgs.go ]; } ''
          unformatted=$(find ${src} -name '*.go' -type f -exec gofmt -l {} +)
          if [ -n "$unformatted" ]; then
            echo "$unformatted"
            exit 1
          fi
          touch $out
        '';
        actionlint =
          pkgs.runCommand "${pname}-actionlint"
            {
              nativeBuildInputs = [ pkgs.actionlint ];
            }
            ''
              actionlint -config-file ${src}/.github/actionlint.yaml ${src}/.github/workflows/*.yml
              touch $out
            '';
        dockerImage = pkgs.dockerTools.buildLayeredImage {
          name = pname;
          tag = version;
          contents = [
            app
            pkgs.cacert
            pkgs.dockerTools.fakeNss
          ];
          config = {
            Cmd = [ "${app}/bin/${pname}" ];
            Env = [
              "PORT=7883"
              "DB_URL=/var/db/prod.db"
              "ENCRYPTION_KEY=ABCDEFGHIJKLMNOPQRSTUVWXYZ"
              "VERSION=v0.0.1+${self.shortRev or "nix"}"
              "COMMIT_SHA=${self.shortRev or "unknown"}"
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
            ExposedPorts."7883/tcp" = { };
            Volumes."/var/db" = { };
            WorkingDir = "${app}/share/${pname}";
          };
        };
        formatter = pkgs.writeShellApplication {
          name = "nix-fmt";
          runtimeInputs = [ pkgs.nixfmt ];
          text = ''
            if [ "$#" -eq 0 ]; then
              set -- flake.nix bun.nix
            fi
            exec nixfmt "$@"
          '';
        };
      in
      {
        packages = {
          default = app;
          inherit css dockerImage;
        };

        apps.default = {
          type = "app";
          program = "${app}/bin/${pname}";
        };

        checks = {
          inherit
            actionlint
            css
            dockerImage
            gofmt
            ;
          build = app;
          tests = mkGoCheck "tests" "go test ./...";
          vet = mkGoCheck "vet" "go vet ./...";
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.bun
            pkgs.bun2nix
            pkgs.go
            pkgs.nixfmt
          ];
        };

        inherit formatter;
      }
    );
}
