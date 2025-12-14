{ pkgs, lib, ... }:
let
  # Define your tool as a Nix package
  git-quill = pkgs.buildGoModule {
    pname = "git-quill";
    version = "0.1.0";

    # Use the current directory as the source
    src = lib.cleanSource ./.;

    # 1. SET THIS TO null FIRST.
    # Nix will complain and give you the correct hash to paste here.
    vendorHash = "sha256-gAi4a3ZrmImd3m9TX0Te0PhVlwRiRqVV8+vdsyzdflg=";
    # vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

    # Optimize binary size (strip debug symbols)
    ldflags = [
      "-s"
      "-w"
    ];

    # Skip tests during dev builds for speed
    doCheck = false;
  };
in
{
  packages =
    (with pkgs; [
      goreleaser
    ])
    ++ [
      # git-quill
    ];
  languages.go.enable = true;
}
