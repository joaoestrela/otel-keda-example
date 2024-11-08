{
  description = "Otel Keda Example Application";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            delve
            git
            go_1_23
            golangci-lint
            helm-docs
            kind
            kubectl
            kubernetes-helm
            kustomize
            nixd
            nixfmt-rfc-style
            protobuf_24
            protoc-gen-go
            protoc-gen-go-grpc
            grpcurl
          ];
        };
      }
    );
}
