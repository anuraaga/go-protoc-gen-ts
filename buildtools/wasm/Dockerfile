FROM rust:1-slim-bookworm

RUN apt-get update && apt-get install -y binaryen curl libssl-dev pkg-config && rustup target add wasm32-wasi
RUN cargo install cargo-patch

RUN mkdir -p /workspace && curl -L https://github.com/thesayyn/protoc-gen-ts/archive/578f3a0a508878d7ffc277d4bb7c47832c2cc5ae.tar.gz | tar -xz --strip-components 1 -C /workspace
WORKDIR /workspace

RUN cargo-patch
RUN cargo build --release --target wasm32-wasi

RUN wasm-opt -o target/wasm32-wasi/release/protoc-gen-ts.wasm --flatten --rereloop --converge -O3 target/wasm32-wasi/release/protoc-gen-ts.wasm

CMD ["cp", "target/wasm32-wasi/release/protoc-gen-ts.wasm", "/out/"]
