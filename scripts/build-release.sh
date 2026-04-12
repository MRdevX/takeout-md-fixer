#!/usr/bin/env bash
# Build release artifacts locally (macOS: .app + .dmg; Windows amd64 + arm64: .exe cross-compile).
# Prerequisites: Go, Node, Wails v3 CLI (`wails3`), macOS for DMG (hdiutil).
# Usage:
#   ./scripts/build-release.sh              # uses version from build/config.yml
#   VERSION=1.2.3 ./scripts/build-release.sh # override version in config for this build only
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

# shellcheck source=lib/release-common.sh
source "${ROOT}/scripts/lib/release-common.sh"

VERSION="${VERSION:-}"
if [[ -z "${VERSION}" ]]; then
	VERSION="$(grep -E '^  version:' build/config.yml | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
fi

echo "Building release version: ${VERSION}"
apply_version_to_config "${ROOT}" "${VERSION}"

echo "Updating Wails build assets..."
(
	cd build
	wails3 update build-assets -name takeout-md-fixer -binaryname takeout-md-fixer -config config.yml -dir .
)

echo "Packaging (native platform)..."
wails3 package

if [[ "$(uname -s)" == Darwin ]]; then
	echo "Creating DMG..."
	STAGE="${ROOT}/bin/dmg-staging"
	rm -rf "${STAGE}"
	mkdir -p "${STAGE}"
	cp -R "${ROOT}/bin/takeout-md-fixer.app" "${STAGE}/"
	ln -sf /Applications "${STAGE}/Applications"
	rm -f "${ROOT}/bin/takeout-md-fixer.dmg"
	hdiutil create -volname "Takeout Metadata Fixer" -srcfolder "${STAGE}" -ov -format UDZO "${ROOT}/bin/takeout-md-fixer.dmg"
	rm -rf "${STAGE}"
	echo "DMG: bin/takeout-md-fixer.dmg"
else
	echo "Skipping DMG (macOS only)."
fi

echo "Cross-compiling Windows amd64 and arm64 exes..."
for WIN_ARCH in amd64 arm64; do
	(
		cd build
		wails3 generate syso -arch "${WIN_ARCH}" -icon windows/icon.ico -manifest windows/wails.exe.manifest -info windows/info.json -out "../wails_windows_${WIN_ARCH}.syso"
	)
	GOOS=windows GOARCH="${WIN_ARCH}" CGO_ENABLED=0 go build -tags production -trimpath -buildvcs=false -ldflags="-w -s -H windowsgui" -o "bin/takeout-md-fixer-windows-${WIN_ARCH}.exe" .
	rm -f "wails_windows_${WIN_ARCH}.syso"
done

echo "Done."
ls -la bin/takeout-md-fixer.app bin/takeout-md-fixer.dmg bin/takeout-md-fixer-windows-*.exe 2>/dev/null || ls -la bin/
