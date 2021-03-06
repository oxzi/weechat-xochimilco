{ config ? {}, overlays ? [] }:

let
  mkNixpkgs = { rev, sha256 }: import (builtins.fetchTarball {
    inherit sha256;
    url = "https://github.com/NixOS/nixpkgs/archive/${rev}.tar.gz";
  });
  nixpkgs = mkNixpkgs (import ./nixpkgs.nix);

  overlaysAll = [(import ./overlays.nix)] ++ overlays;
in nixpkgs { inherit config; overlays=overlaysAll; }
