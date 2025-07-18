{
  description = "CLI TODO Tool";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs =
    {
      nixpkgs,
      flake-utils,
      gomod2nix,
      self,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        lib = pkgs.lib;

        gomod2nixBuilders = gomod2nix.legacyPackages.${system};
        buildGoApplication = gomod2nixBuilders.buildGoApplication;
        mkGoEnv = gomod2nixBuilders.mkGoEnv;
        gomod2nixTool = gomod2nixBuilders.gomod2nix;

        commonArgs = {
          version = "0.1.0";
          src = ./.;
          pwd = ./.;
          modules = ./gomod2nix.toml;
        };

        app = buildGoApplication (
          commonArgs
          // {
            pname = "app";
            subPackages = [ "." ];
            meta = {
              description = "App package";
              license = lib.licenses.mit;
            };
          }
        );

        devGoEnv = mkGoEnv {
          pwd = ./.;
        };

      in
      {
        packages = {
          app = app;
          default = app;
        };

        devShells.default = pkgs.mkShell {
          packages = [
            devGoEnv
            gomod2nixTool
          ];
        };
      }
    ));
}
