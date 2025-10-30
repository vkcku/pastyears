{
  description = "A website to help prepare for conducted exams.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages."${system}";

      buildInputs = {
        core = with pkgs; [ go ];

        dev = with pkgs; [
          # keep-sorted start
          dbmate
          sqlite
          # keep-sorted end
        ];

        # All the dependencies required for the `lint` check.
        linters = with pkgs; [
          # Formatters/linters.
          # keep-sorted start
          golangci-lint
          keep-sorted
          nixfmt-rfc-style
          python313Packages.mdformat
          python313Packages.mdformat-gfm
          sql-formatter
          taplo
          treefmt
          typos
          yamlfmt
          # keep-sorted end
        ];

        misc = with pkgs; [
          # keep-sorted start
          gopls
          yaml-language-server
          # keep-sorted end
        ];
      };
    in
    {
      devShells."${system}".default = pkgs.mkShell {
        buildInputs = nixpkgs.lib.flatten (builtins.attrValues buildInputs);

        shellHook = ''
          ROOT_DIR="$(git rev-parse --show-toplevel)"

          # For dbmate.
          export DATABASE_URL="sqlite:$ROOT_DIR/data.db"
        '';
      };

      packages."${system}" = rec {
        default = webserver;

        webserver = pkgs.buildGoModule {
          name = "pastyears-webserver";
          src = ./.;
          subPackages = [ "cmd/webserver" ];
          doCheck = false;
          # vendorHash = nixpkgs.lib.fakeHash;
          vendorHash = "sha256-Swi56SaPh4AN7LZ2a+j3p/jNf/InnbmE6AEErjqLg0g=";
        };
      };

      checks."${system}" = {
        lint = self.packages."${system}".webserver.overrideAttrs (old: {
          name = "lint";
          nativeBuildInputs = old.nativeBuildInputs ++ buildInputs.linters;
          buildPhase = ''
            XDG_CACHE_HOME="$TMPDIR" treefmt --ci
            touch "$out"
          '';
          installPhase = '''';
          fixupPhase = '''';
          env = {
            CGO_CFLAGS = "-O1";
          };
        });

        test = self.packages."${system}".webserver.overrideAttrs (old: {
          name = "test";
          nativeBuildInputs = old.nativeBuildInputs ++ buildInputs.linters;
          env = {
            CGO_CFLAGS = "-O1";
          };
          buildPhase = ''
            XDG_CACHE_HOME="$TMPDIR" go test ./...
            touch "$out"
          '';
          installPhase = '''';
          fixupPhase = '''';
        });
      };
    };
}
