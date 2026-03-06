#!/bin/sh

set -e

REPO="utkarsh867/tfcdash"
BINARY_NAME="tfcdash"

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       echo "unknown" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64)   echo "x86_64" ;;
        amd64)    echo "x86_64" ;;
        aarch64)  echo "arm64" ;;
        arm64)    echo "arm64" ;;
        armv7*)   echo "arm" ;;
        *)        echo "unknown" ;;
    esac
}

get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": "v?([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": "v?([^"]+)".*/\1/')
    else
        echo "Error: curl or wget is required to install ${BINARY_NAME}" >&2
        exit 1
    fi
    echo "$version"
}

get_download_url() {
    os="$1"
    arch="$2"
    version="${3:-latest}"

    if [ "$version" = "latest" ]; then
        version=$(get_latest_version)
    fi

    echo "https://github.com/${REPO}/releases/download/v${version}/${BINARY_NAME}_${os}_${arch}.tar.gz"
}

install_binary() {
    os="$1"
    arch="$2"
    version="$3"

    url=$(get_download_url "$os" "$arch" "$version")

    echo "Installing ${BINARY_NAME} ${version} for ${os}/${arch}..."

    tmp_dir=$(mktemp -d)
    cd "$tmp_dir"

    if command -v curl >/dev/null 2>&1; then
        curl -sL "$url" -o "${BINARY_NAME}.tar.gz"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$url" -O "${BINARY_NAME}.tar.gz"
    else
        echo "Error: curl or wget is required to install ${BINARY_NAME}" >&2
        exit 1
    fi

    tar xzf "${BINARY_NAME}.tar.gz"
    rm -f "${BINARY_NAME}.tar.gz"

    install_dir="${HOME}/.local/bin"
    mkdir -p "$install_dir"
    cp "$BINARY_NAME" "$install_dir/"
    rm -rf "$tmp_dir"

    echo "Installed ${BINARY_NAME} to ${install_dir}"

    if [ -w "/usr/local/bin" ]; then
        cp "$install_dir/$BINARY_NAME" /usr/local/bin/
        echo "Also copied to /usr/local/bin"
    fi
}

main() {
    os=$(detect_os)
    arch=$(detect_arch)

    if [ "$os" = "unknown" ]; then
        echo "Error: Unsupported operating system" >&2
        exit 1
    fi

    if [ "$arch" = "unknown" ]; then
        echo "Error: Unsupported architecture" >&2
        exit 1
    fi

    version="${1:-latest}"

    install_binary "$os" "$arch" "$version"

    echo ""
    echo "Installation complete!"
    echo "Make sure ${HOME}/.local/bin is in your PATH"
}

main "$@"
