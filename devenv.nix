{ pkgs, ... }:
{
  packages = with pkgs; [ gum ];
  languages.go.enable = true;
}
