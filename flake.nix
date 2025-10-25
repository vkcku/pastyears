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
        core = with pkgs; [
          go
          postgresql_18
        ];

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
          vendorHash = "sha256-Bwv/TTpr/oTj3ufNugb/pNeYH4rXEyAqAPEQNtBuA20=";
        in
        {
          default = pkgs.buildGoModule {
            name = "pastyears-webserver";
            src = ./.;
            subPackages = [ "cmd/webserver" ];
            doCheck = false;

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

        shellHook = ''
          export PGHOST="$TMPDIR/pastyears/pg"
          export PGDATABASE="pastyears"
          export PGPORT=5432
        '';
      };

      checks."${system}" = {
        test = self.packages."${system}".default.overrideAttrs (old: {
          name = "test";
          buildPhase = "touch $out";
          installPhase = "";
          checkPhase = ''
            go test ./...
          '';
          fixupPhase = "";
          doCheck = true;
        });

        lint = self.packages."${system}".default.overrideAttrs (old: {
          nativeBuildInputs =
            old.nativeBuildInputs ++ buildInputs.formatters ++ [ self.packages.${system}.cli ];
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
            XDG_CACHE_HOME=$TMPDIR "${self.packages.${system}.cli}/bin/pastyears" lint
          '';
          doCheck = false;
          installPhase = ''
            touch $out
          '';
          fixupPhase = '''';
          env = old.env // {
            CI = 1;
          };
        });
      };
    };
}
