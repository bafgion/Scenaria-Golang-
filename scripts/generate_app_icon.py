#!/usr/bin/env python3
"""Build opaque Windows icons from branding PNGs.

Rounded icons with transparent corners look semi-transparent on the Windows
desktop. This script flattens alpha onto the app background (#252526) and
regenerates app.ico plus Wails build assets.
"""

from __future__ import annotations

from pathlib import Path

from PIL import Image

# Matches --color-sidebar in frontend/src/style.css and app-icon.svg
BG_RGB = (37, 37, 38)
ICO_SIZES = [(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)]
MASTER_PX = 512


def repo_root() -> Path:
    return Path(__file__).resolve().parent.parent


def flatten_rgba(image: Image.Image, bg: tuple[int, int, int] = BG_RGB) -> Image.Image:
    base = Image.new("RGBA", image.size, (*bg, 255))
    return Image.alpha_composite(base, image.convert("RGBA")).convert("RGB")


def save_ico(path: Path, master: Image.Image) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    master.save(path, format="ICO", sizes=ICO_SIZES)


def main() -> None:
    root = repo_root()
    branding = root / "assets" / "branding"
    square_src = branding / "app-icon-square.png"
    mark_src = branding / "app-icon-mark.png"

    if not square_src.exists():
        raise SystemExit(f"missing source icon: {square_src}")

    square = flatten_rgba(Image.open(square_src))
    square.save(branding / "app-icon-square.png", optimize=True)

    if mark_src.exists():
        mark = flatten_rgba(Image.open(mark_src))
    else:
        mark = square
    mark.save(branding / "app-icon-mark.png", optimize=True)

    master = square.resize((MASTER_PX, MASTER_PX), Image.Resampling.LANCZOS)
    save_ico(branding / "app.ico", master)

    build = root / "build"
    (build / "windows").mkdir(parents=True, exist_ok=True)
    master.save(build / "appicon.png", optimize=True)
    save_ico(build / "windows" / "icon.ico", master)

    frontend = root / "frontend" / "src" / "assets" / "branding"
    frontend.mkdir(parents=True, exist_ok=True)
    square.save(frontend / "app-icon-square.png", optimize=True)
    mark.save(frontend / "app-icon-mark.png", optimize=True)

    print(f"OK: opaque icons -> {branding / 'app.ico'}")
    print(f"OK: wails assets -> {build / 'appicon.png'}, {build / 'windows' / 'icon.ico'}")


if __name__ == "__main__":
    main()
