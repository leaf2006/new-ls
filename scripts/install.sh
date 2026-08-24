#!/usr/bin/env bash
# =============================================================================
# nls install script (Linux / macOS)
#
#   Builds new-ls from source into a binary named `nls`, then installs it to
#   /usr/local/bin (configurable via --prefix) so you can run `nls` anywhere.
#
#   Usage:
#     ./scripts/install.sh             # build + install
#     ./scripts/install.sh --prefix ~/.local/bin
#     ./scripts/install.sh --uninstall # remove the installed binary
#
#   If the Go toolchain is missing, the script offers to download an official
#   Go toolchain into a user-local cache directory (NOT your system) purely to
#   build nls. Pass --no-go-install to refuse this and get manual instructions.
# =============================================================================

set -euo pipefail

# ---------------------------------------------------------------- config ----
BINARY_NAME="nls"
DEFAULT_PREFIX="/usr/local/bin"
# Used only when the latest Go version cannot be fetched from go.dev.
GO_FALLBACK_VERSION="1.25.6"
# A temporary Go toolchain (if needed) lives here, reused across runs.
GO_CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/nls-install"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

ACTION="install"
PREFIX="$DEFAULT_PREFIX"
# prompt | always | never   (whether to auto-download a Go toolchain)
INSTALL_GO_BEHAVIOR="prompt"

# Clean up the temporary built binary on any exit path (incl. failures).
# Always returns 0 so it never overrides the script's own exit status.
BUILD_TMP=""
cleanup() { [[ -z "$BUILD_TMP" ]] || rm -f "$BUILD_TMP"; }
trap cleanup EXIT

# --------------------------------------------------------------- helpers ----
# Info/status messages go to stderr so stdout stays clean for real output
# (e.g. build_nls prints only the binary path on stdout).
say()   { printf '\033[1;34m[nls]\033[0m %s\n' "$*" >&2; }
ok()    { printf '\033[1;32m[nls]\033[0m %s\n' "$*" >&2; }
warn()  { printf '\033[1;33m[nls]\033[0m %s\n' "$*" >&2; }
die()   { printf '\033[1;31m[nls]\033[0m %s\n' "$*" >&2; exit 1; }

usage() {
  cat <<EOF
nls installer (Linux / macOS)

Builds this new-ls repo and installs the binary to ${DEFAULT_PREFIX} (or --prefix).

Usage: install.sh [options]

Options:
  -y, --yes            Skip the Go-install confirmation prompt
  --no-go-install      Never download a Go toolchain; print manual steps instead
  --prefix <dir>       Install to <dir> instead of ${DEFAULT_PREFIX}
  -u, --uninstall      Remove the installed nls binary
  -h, --help           Show this help

Note: installed for real use at \$PREFIX, so a writable dir avoids sudo;
      otherwise you will be prompted for sudo.
EOF
}

# ------------------------------------------------------------- detection ----
detect_os() {
  case "$(uname -s)" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *) die "Unsupported OS: $(uname -s). This installer supports Linux and macOS only." ;;
  esac
}

# Only needed when we must download a Go toolchain (a local Go avoids this).
detect_arch_for_go_download() {
  case "$(uname -m)" in
    x86_64|amd64)       ARCH="amd64" ;;
    aarch64|arm64)      ARCH="arm64" ;;
    *) die "Unsupported architecture: $(uname -m). Install Go manually and re-run with --no-go-install." ;;
  esac
}

# ------------------------------------------------------ go toolchain -------
ensure_go() {
  if command -v go >/dev/null 2>&1; then
    say "Using existing Go toolchain: $(go version)"
    return 0
  fi

  warn "The Go toolchain was not found on this system."

  if [[ "$INSTALL_GO_BEHAVIOR" == "never" ]]; then
    die "Aborting. Please install Go (https://go.dev/dl/) then re-run this script."
  fi

  # Interactive terminal -> ask; piped stdin (e.g. curl | bash) -> proceed.
  if [[ "$INSTALL_GO_BEHAVIOR" == "prompt" && -t 0 ]]; then
    printf '\033[1;33m[nls]\033[0m Download a temporary Go toolchain to build nls? [Y/n] '
    read -r ans
    case "${ans:-Y}" in
      [yY]*) : ;;
      *) die "Aborting. Install Go (https://go.dev/dl/) then re-run this script." ;;
    esac
  fi

  install_go_toolchain
}

install_go_toolchain() {
  detect_arch_for_go_download
  local go_os="$OS"

  # Reuse a previously downloaded toolchain if present and functional.
  if [[ -x "$GO_CACHE_DIR/go/bin/go" ]] && "$GO_CACHE_DIR/go/bin/go" version >/dev/null 2>&1; then
    export PATH="$GO_CACHE_DIR/go/bin:$PATH"
    say "Using cached Go toolchain: $("$GO_CACHE_DIR/go/bin/go" version)"
    return 0
  fi

  local version
  version="$(latest_go_version)"

  local url="https://go.dev/dl/${version}.${go_os}-${ARCH}.tar.gz"
  say "Downloading Go toolchain ${version} (${go_os}/${ARCH}) to build nls ..."

  local tool=""
  if command -v curl >/dev/null 2>&1; then
    tool="curl"
  elif command -v wget >/dev/null 2>&1; then
    tool="wget"
  else
    die "Need curl or wget to download the Go toolchain."
  fi

  local tarball attempt
  tarball="$(mktemp "${TMPDIR:-/tmp}/nls-go.XXXXXX.tar.gz")"
  for attempt in 1 2 3; do
    rm -f "$tarball"
    if [[ "$tool" == "curl" ]]; then
      curl -fsSL --max-time 120 "$url" -o "$tarball"
    else
      wget -q --timeout=120 --tries=3 "$url" -O "$tarball"
    fi && break
    warn "Download failed (attempt ${attempt}/3); retrying ..."
    sleep 2
  done

  if [[ ! -s "$tarball" ]]; then
    rm -f "$tarball"
    die "Failed to download the Go toolchain from $url. Check your network and re-run."
  fi

  mkdir -p "$GO_CACHE_DIR"
  rm -rf "$GO_CACHE_DIR/go"                 # clear any partial/old extraction
  tar -C "$GO_CACHE_DIR" -xzf "$tarball"
  rm -f "$tarball"

  export PATH="$GO_CACHE_DIR/go/bin:$PATH"  # session-local; system untouched
  command -v go >/dev/null 2>&1 \
    || die "Failed to set up the Go toolchain at $GO_CACHE_DIR/go."
  ok "Go toolchain ready: $(go version)"
}

latest_go_version() {
  local v="" attempt
  for attempt in 1 2 3; do
    if command -v curl >/dev/null 2>&1; then
      v="$(curl -fsSL --max-time 20 'https://go.dev/VERSION?m=text' 2>/dev/null | head -n1 || true)"
    elif command -v wget >/dev/null 2>&1; then
      v="$(wget -qO- --timeout 20 'https://go.dev/VERSION?m=text' 2>/dev/null | head -n1 || true)"
    fi
    if [[ "$v" == go[0-9.]* ]]; then
      printf '%s\n' "$v"
      return 0
    fi
    v=""
    sleep 1
  done
  warn "Could not determine the latest Go version; falling back to go${GO_FALLBACK_VERSION}."
  printf 'go%s\n' "$GO_FALLBACK_VERSION"
}

# ---------------------------------------------------------------- build ----
MIRROR_URLS=(
  "https://goproxy.cn,direct"
  "https://goproxy.io,direct"
  "https://mirrors.aliyun.com/goproxy/,direct"
  "https://mirrors.cloud.tencent.com/go/,direct"
)
MIRROR_NAMES=(
  "https://goproxy.cn"
  "https://goproxy.io"
  "https://mirrors.aliyun.com/goproxy/"
  "https://mirrors.cloud.tencent.com/go/"
)

select_mirror() {
  local idx
  while true; do
    echo "" >&2
    warn "国内可用的 Go 模块镜像源："
    local i
    for (( i=0; i<${#MIRROR_NAMES[@]}; i++ )); do
      printf '  %d) %s\n' "$((i+1))" "${MIRROR_NAMES[$i]}" >&2
    done
    echo "" >&2
    printf '\033[1;33m[nls]\033[0m 请输入编号 (1-%d) 选择镜像源: ' "${#MIRROR_NAMES[@]}" >&2
    read -r idx
    if [[ "$idx" =~ ^[1-4]$ ]]; then
      echo "${MIRROR_URLS[$((idx-1))]}"
      return 0
    fi
    warn "无效输入，请输入 1 到 ${#MIRROR_NAMES[@]} 之间的数字。"
  done
}

try_build() {
  local proxy="$1"
  local out="$2"
  if (cd "$PROJECT_DIR" && GOPROXY="$proxy" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$out" ./cmd/new-ls 2>&1); then
    return 0
  else
    return 1
  fi
}

build_nls() {
  say "Building ${BINARY_NAME} ..."
  [[ -f "$PROJECT_DIR/go.mod" ]] \
    || die "install.sh must live inside the new-ls repo (in scripts/)."
  local out
  out="$(mktemp "${TMPDIR:-/tmp}/nls-build.XXXXXX")"
  BUILD_TMP="$out"

  local build_log

  if [[ -n "${GOPROXY:-}" ]]; then
    if try_build "$GOPROXY" "$out"; then
      ok "Build finished."
      printf '%s\n' "$out"
      return 0
    fi
    rm -f "$out"
    BUILD_TMP=""
    die "Build failed. Check your GOPROXY setting ($GOPROXY) and network."
  fi

  build_log="$(try_build "https://proxy.golang.org,direct" "$out" 2>&1)" && {
    ok "Build finished."
    printf '%s\n' "$out"
    return 0
  }

  if echo "$build_log" | grep -qiE 'proxy\.golang\.org|connection refused|dial tcp|no such host|timeout|i/o timeout'; then
    warn "无法连接 proxy.golang.org，可能是网络环境限制。"
    printf '\033[1;33m[nls]\033[0m 是否切换为国内镜像源？[y/N] ' >&2
    read -r ans
    case "${ans:-N}" in
      [yY]*) ;;
      *)
        rm -f "$out"
        BUILD_TMP=""
        die "Build failed. 请设置 GOPROXY 环境变量或配置网络后重试。"
        ;;
    esac

    local selected_proxy
    selected_proxy="$(select_mirror)"
    say "使用镜像源: ${selected_proxy%%,*}"

    if try_build "$selected_proxy" "$out"; then
      ok "Build finished."
      printf '%s\n' "$out"
      return 0
    fi
  fi

  rm -f "$out"
  BUILD_TMP=""
  die "Build failed. 请检查 Go 工具链和网络连接。"
}

# --------------------------------------------------------------- install ---
install_binary() {
  local src="$1"
  local dst="$PREFIX/$BINARY_NAME"
  mkdir -p "$PREFIX" 2>/dev/null || true
  if [[ -w "$PREFIX" ]]; then
    install -m 0755 "$src" "$dst"
  elif command -v sudo >/dev/null 2>&1; then
    sudo install -m 0755 "$src" "$dst"
  else
    die "Cannot write to $PREFIX and 'sudo' is not available. Use --prefix <writable-dir> or run as root."
  fi
  ok "Installed $dst"
}

uninstall_nls() {
  local dst="$PREFIX/$BINARY_NAME"
  if [[ ! -e "$dst" ]]; then
    warn "nls is not installed at $dst."
    return 0
  fi
  if [[ -w "$PREFIX" ]]; then
    rm -f "$dst"
  elif command -v sudo >/dev/null 2>&1; then
    sudo rm -f "$dst"
  else
    die "Cannot write to $PREFIX. Remove $dst manually or run with sudo."
  fi
  ok "Removed $dst"
}

# ---------------------------------------------------------------- main -----
detect_os

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)        usage; exit 0 ;;
    -u|--uninstall)   ACTION="uninstall" ;;
    -y|--yes)         INSTALL_GO_BEHAVIOR="always" ;;
    --no-go-install)  INSTALL_GO_BEHAVIOR="never" ;;
    --prefix)
      [[ $# -ge 2 ]] || die "--prefix requires a directory argument."
      PREFIX="$2"; shift ;;
    *) die "Unknown option: $1 (run install.sh --help)" ;;
  esac
  shift
done

if [[ "$ACTION" == "uninstall" ]]; then
  uninstall_nls
  exit 0
fi

ensure_go
bin="$(build_nls)"
BUILD_TMP="$bin"          # set here so the EXIT trap can clean up on any failure
install_binary "$bin"
rm -f "$bin"

printf '\033[1;32m[nls]\033[0m Done. Run \033[1mnls\033[0m to use it ' >&2
printf '(re-open your shell if %s is not on PATH).\n' "$PREFIX" >&2
