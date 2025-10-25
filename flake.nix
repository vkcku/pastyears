{
  description = "A website to help prepare for UPSC conducted exams.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      lib = nixpkgs.lib;
      pkgs = nixpkgs.legacyPackages.${system};

      buildInputs = {
        # keep-sorted start block=yes newline_separated=yes
        core = with pkgs; [ go ];

        formatters = with pkgs; [
          # keep-sorted start
          golangci-lint
          keep-sorted
          nixfmt-rfc-style
          python313Packages.mdformat
          python313Packages.mdformat-gfm
          taplo
          treefmt
          typos
          # keep-sorted end
        ];

        lsps = with pkgs; [
          # keep-sorted start
          gopls
          nil
          yaml-language-server
          # keep-sorted end
        ];
        # keep-sorted end
      };
    in
    {
      packages."${system}" =
        let
          # vendorHash = lib.fakeHash;
          vendorHash = "sha256-JpkT7S+zhw46QHuMXXtPGmYAr+oFE3Jst26biVossEA=";
        in
        {
          default = pkgs.buildGoModule {
            name = "pastyears-webserver";
            src = ./.;
            subPackages = [ "cmd/webserver" ];

            # `buildGoModule` will only run the tests for the package being built. That
            # is, it effectively does `go test cmd/webserver` which means not all
            # the tests are ran. So instead, override the default `checkPhase` to run
            # all the tests across the entire project.
            #
            # REFERENCE: https://github.com/NixOS/nixpkgs/blob/39231460a6f5e193a193a44902877c1c0026f271/pkgs/build-support/go/module.nix#L313
            checkPhase = ''
              go test ./...
            '';

            inherit vendorHash;
          };

          cli = pkgs.buildGoModule {
            name = "pastyears";
            src = ./.;
            subPackages = [ "cmd/pastyears" ];
            doCheck = false;

            inherit vendorHash;
          };
        };

      devShells.${system}.default = pkgs.mkShell {
        buildInputs = lib.lists.flatten (builtins.attrValues buildInputs) ++ [
          self.packages."${system}".cli
        ];
      };

      checks."${system}" = {
        build = self.packages."${system}".default;
        build-cli = self.packages."${system}".cli;
        lint = self.packages."${system}".default.overrideAttrs (old: {
          nativeBuildInputs = old.nativeBuildInputs ++ buildInputs.formatters;
          buildPhase = ''
            # `golangci-lint` and `go` creates some cache directories using
            # `os.UserCacheDir` which takes the value from `$XDG_CACHE_HOME` or
            # sets the value as `$HOME/.cache` if `XDG_CACHE_HOME` is not found.
            # In the nix build, the $HOME directory is read-only so they both
            # fail to create files within that cache directory.
            #
            # REFERENCE: https://github.com/NixOS/nixpkgs/issues/202614
            # More specifically, this comment:
            # https://github.com/NixOS/nixpkgs/issues/202614#issuecomment-1326152971
            XDG_CACHE_HOME=$TMPDIR treefmt \
              --ci \
              --walk filesystem \
          '';
          doCheck = false;
          installPhase = ''
            touch $out
          '';
          fixupPhase = '''';
        });
      };
    };
}
