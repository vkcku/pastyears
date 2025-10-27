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

        # All the dependencies required for the `lint` check.
        linters = with pkgs; [
          # Formatters/linters.
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

        misc = with pkgs; [
          yaml-language-server
        ];
      };
    in
    {
      devShells."${system}".default = pkgs.mkShell {
        buildInputs = nixpkgs.lib.flatten (builtins.attrValues buildInputs);
      };

      checks."${system}" = {
        lint =
          pkgs.runCommandLocal "lint"
            {
              nativeBuildInputs = buildInputs.linters;
              src = ./.;
            }
            ''
              # If `XDG_CACHE_HOME` is not set, then `$HOME/.cache` is used which is
              # not writable in checks.
              XDG_CACHE_HOME="$TMPDIR" treefmt --ci --working-dir "$src"
              touch "$out"
            '';
      };
    };
}
