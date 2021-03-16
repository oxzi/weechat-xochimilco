{ config ? {}, overlays ? [], pkgs ? import ./contrib/nix { inherit config overlays; } }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    clang-tools
    go
    golangci-lint
    reuse
    weechat
  ];
}
