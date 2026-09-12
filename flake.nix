{
  nixConfig = {
    extra-substituters = ["https://nix-community.cachix.org"];
    extra-trusted-public-keys = ["nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="];
  };

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2605";
    flake-utils.url = "github:numtide/flake-utils";
    git-hooks = {
      url = "https://flakehub.com/f/cachix/git-hooks.nix/0.1";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    bun2nix = {
      url = "github:nix-community/bun2nix?ref=2.1.2";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    git-hooks,
    bun2nix,
  }:
    flake-utils.lib.eachSystem ["x86_64-linux" "aarch64-linux"] (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [bun2nix.overlays.default];
        };
        pname = "htmx-go";
        version = "0.0.1";
        vendorHash = "sha256-zT3tt5+9xrXNfolWPZ9H6LFNz4NzeQdlmY/QtKYk9NE=";
        src = pkgs.lib.cleanSourceWith {
          src = ./.;
          filter = path: _type:
            !builtins.elem (baseNameOf path) [
              ".git"
              ".jj"
              "result"
            ];
        };
        bunDeps = pkgs.bun2nix.fetchBunDeps {bunNix = ./bun.nix;};
        css = pkgs.bun2nix.mkDerivation {
          pname = "${pname}-css";
          inherit version src bunDeps;
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [pkgs.stdenv.cc.cc.lib];
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
          env.CGO_ENABLED = 0;
          subPackages = ["cmd"];
          postInstall = ''
            mv $out/bin/cmd $out/bin/${pname}
            mkdir -p $out/share/${pname}
            cp -r assets views $out/share/${pname}/
            cp ${css}/style.css $out/share/${pname}/assets/style.css
          '';
          meta.mainProgram = pname;
        };
        mkGoCheck = name: command:
          pkgs.buildGoModule {
            pname = "${pname}-${name}";
            inherit version src;
            inherit vendorHash;
            env.CGO_ENABLED = 0;
            buildPhase = command;
            doCheck = false;
            installPhase = "touch $out";
          };
        gofmt = pkgs.runCommand "${pname}-gofmt" {nativeBuildInputs = [pkgs.go];} ''
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
            nativeBuildInputs = [pkgs.actionlint];
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
            pkgs.busybox
            pkgs.cacert
            pkgs.dockerTools.fakeNss
            pkgs.sqlite
          ];
          config = {
            Cmd = ["${app}/bin/${pname}"];
            User = "65532:65532";
            Env = [
              "PORT=7883"
              "DB_URL=/var/db/prod.db"
              "ENCRYPTION_KEY=ABCDEFGHIJKLMNOPQRSTUVWXYZ"
              "VERSION=v0.0.1+${self.shortRev or "nix"}"
              "COMMIT_SHA=${self.shortRev or "unknown"}"
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
            ExposedPorts."7883/tcp" = {};
            Volumes."/data" = {};
            WorkingDir = "${app}/share/${pname}";
          };
        };
        formatter = pkgs.alejandra;
        preCommitCheck = git-hooks.lib.${system}.run {
          package = pkgs.prek;
          src = ./.;
          hooks = {
            actionlint.enable = true;
            alejandra.enable = true;
            check-added-large-files.enable = true;
            check-merge-conflicts.enable = true;
            check-yaml.enable = true;
            end-of-file-fixer.enable = true;
            gofmt.enable = true;
            shellcheck.enable = true;
            trim-trailing-whitespace.enable = true;
          };
        };
      in {
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
          pre-commit = preCommitCheck;
        };

        devShells.default = pkgs.mkShell {
          packages =
            [
              pkgs.bun
              pkgs.bun2nix
              pkgs.go
              pkgs.alejandra
            ]
            ++ preCommitCheck.enabledPackages;
          shellHook = preCommitCheck.shellHook;
        };

        inherit formatter;
      }
    );
}
