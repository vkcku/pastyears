{
  description = "A website to help prepare for conducted exams.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages."${system}";

      buildInputs = {
        core = with pkgs; [ go ];

        # All the dependencies required for the `lint` check.
        linters = with pkgs; [
          # Formatters/linters.
          # keep-sorted start
          golangci-lint
          keep-sorted
          nixfmt-rfc-style
          python313Packages.mdformat
          python313Packages.mdformat-gfm
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
      };

      packages."${system}" = rec {
        default = webserver;

        webserver = pkgs.buildGoModule {
          name = "pastyears-webserver";
          src = ./.;
          subPackages = [ "cmd/webserver" ];
          doCheck = false;
          vendorHash = null;
        };
      };

      checks."${system}" = {
        lint =
          pkgs.runCommandLocal "lint"
            {
              nativeBuildInputs = buildInputs.linters ++ [
                pkgs.go # needed by `golangci-lint`
              ];
              src = ./.;
            }
            ''
              # If `XDG_CACHE_HOME` is not set, then `$HOME/.cache` is used which is
              # not writable in checks.
              XDG_CACHE_HOME="$TMPDIR" treefmt --ci --working-dir "$src"
              touch "$out"
            '';

        test =
          pkgs.runCommandLocal "test"
            {
              nativeBuildInputs = [ pkgs.go ];
              src = ./.;
            }
            ''
              cd "$src"
              XDG_CACHE_HOME="$TMPDIR" go test ./...
              touch  "$out"
            '';
      };
    };
}
