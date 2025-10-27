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
    in
    {
      devShells."${system}".default = pkgs.mkShell {
        buildInputs = with pkgs; [
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
      };
    };
}
