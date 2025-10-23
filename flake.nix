{
  description = "A website to help prepare for UPSC conducted exams.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      lib = nixpkgs.lib;
      pkgs = nixpkgs.legacyPackages.${system};

      formatters = with pkgs; [
        # keep-sorted start
        keep-sorted
        nixfmt-rfc-style
        python313Packages.mdformat
        python313Packages.mdformat-gfm
        taplo
        treefmt
        typos
        # keep-sorted end
      ];
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = formatters;
      };

      checks."${system}" =
        let
          /**
            Ensure the `$out` directory is created since the derivation will be
            marked as failed otherwise.
          */
          mkScript = script: script + "\n" + "mkdir $out";

          mkChecks =
            checks:
            lib.attrsets.mapAttrs (
              name: check:
              pkgs.runCommandLocal name {
                src = ./.;
                nativeBuildInputs = check.buildInputs;
                dontBuild = true;
              } (mkScript check.script)
            ) checks;
        in
        mkChecks {
          fmt-lint = {
            buildInputs = [ formatters ];
            script = ''
              treefmt \
                --config-file "$src/treefmt.toml" \
                --ci \
                --tree-root "$src" \
                --walk filesystem
            '';
          };
        };
    };
}
