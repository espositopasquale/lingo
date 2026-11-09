import argparse
import ctypes as C
import ctypes.util
import html
from pathlib import Path
import struct
import zipfile


ROOT = Path(__file__).resolve().parent
NAVY = "#101B2D"
CYAN = "#00ADD8"
GREEN = "#009246"
WHITE = "#FFFFFF"
RED = "#CE2B37"

L = "M78 35H144Q155 35 152 47L106 267H217Q228 267 225 279L213 336Q210 348 199 348H17Q4 348 7 335L66 47Q69 35 78 35Z"
I = "M279 112H343Q355 112 352 125L307 335Q304 348 292 348H228Q216 348 219 335L264 125Q267 112 279 112Z"
N = "M385 112H449Q459 112 457 123L454 140C484 112 514 102 546 110C591 122 604 155 593 207L566 337Q564 348 554 348H490Q479 348 481 337L508 209C513 185 505 173 490 174C468 174 450 197 444 226L420 337Q418 348 408 348H343Q332 348 334 337L374 124Q376 112 385 112Z"
G = "M817 113H877Q890 113 887 126L841 337C826 408 774 436 699 432C658 430 625 418 605 402Q596 395 602 385L628 337Q633 328 643 333C667 347 685 354 711 354C745 354 761 337 767 308C744 329 718 337 691 333C620 325 599 270 615 205C632 137 687 102 743 106C775 108 798 122 807 142L811 120Q812 113 817 113Z M768 179C739 163 707 183 698 216C687 253 702 272 727 271C756 270 777 249 784 218C789 198 783 185 768 179Z"
O = "M1012 106C1081 106 1117 139 1109 207C1100 299 1049 347 973 350C901 352 864 313 876 241C887 157 938 107 1012 106Z M996 181C970 181 949 206 945 237C939 267 949 280 971 280C999 280 1019 254 1025 230C1034 189 1023 181 996 181Z"


def lin_shapes():
    return (f'<path d="{L}"/><path d="{I}"/>'
            '<ellipse cx="320" cy="44" rx="45" ry="40" transform="rotate(-16 320 44)"/>'
            f'<path d="{N}"/>')


def wordmark(prefix, light=False, mono=None):
    shapes = lin_shapes()
    go = f'<path fill-rule="evenodd" d="{G}"/><path fill-rule="evenodd" d="{O}"/>'
    if mono:
        return f'<g fill="{mono}">{shapes}{go}</g>'
    outline = (f'<g fill="none" stroke="{NAVY}" stroke-width="3" '
               f'stroke-linejoin="round">{shapes}</g>') if light else ""
    return f'''<defs><clipPath id="{prefix}-lin">{shapes}</clipPath></defs>
<g clip-path="url(#{prefix}-lin)">
<path fill="{WHITE}" d="M-100-100H1200V600H-100Z"/>
<path fill="{GREEN}" d="M-100-100H365L-90 600H-100Z"/>
<path fill="{RED}" d="M585-100H1200V600H130Z"/>
</g>{outline}<g fill="{CYAN}">{go}</g>'''


def document(name, width, height, body, description):
    key = name.replace(".", "-")
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}" role="img" aria-labelledby="{key}-title {key}-desc">
<title id="{key}-title">Lingo</title>
<desc id="{key}-desc">{html.escape(description)}</desc>
{body}
</svg>
'''


def horizontal(name, light=False, background=None, mono=None):
    bg = f'<path fill="{background}" d="M0 0H1200V500H0Z"/>' if background else ""
    body = f'{bg}<g transform="translate(44 40)">{wordmark(name, light, mono)}</g>'
    return document(name, 1200, 500, body,
                    "Lingo wordmark. Diagonal Italian tricolor on Lin; Go cyan on go."
                    if not mono else "Lingo single-color wordmark.")


def app_icon(name):
    body = (f'<rect width="1024" height="1024" rx="224" fill="{NAVY}"/>'
            f'<g transform="translate(58 332) scale(.82)">{wordmark(name)}</g>')
    return document(name, 1024, 1024, body, "Lingo full wordmark on a navy rounded-square app icon.")


def monogram(name, rounded=True):
    radius = 112 if rounded else 0
    body = f'''<rect width="512" height="512" rx="{radius}" fill="{NAVY}"/>
<defs><clipPath id="{name}-l"><path d="{L}"/></clipPath></defs>
<g transform="translate(48 57) scale(.82)" clip-path="url(#{name}-l)">
<path fill="{GREEN}" d="M-100-100H400V600H-100Z"/>
<path fill="{WHITE}" d="M220-100H285L75 600H10Z"/>
<path fill="{RED}" d="M285-100H400V600H75Z"/>
</g><path fill="{CYAN}" fill-rule="evenodd" transform="translate(-262 57) scale(.82)" d="{G}"/>'''
    return document(name, 512, 512, body, "Lingo compact Lg monogram: tricolor L and Go-cyan g on navy.")


def social_card():
    body = f'''<path fill="{NAVY}" d="M0 0H1200V630H0Z"/>
<path fill="{GREEN}" d="M72 70H108V76H72Z"/>
<path fill="{WHITE}" d="M108 70H144V76H108Z"/>
<path fill="{RED}" d="M144 70H180V76H144Z"/>
<g transform="translate(100 129) scale(.90)">{wordmark("social")}</g>
<text x="72" y="565" font-family="DejaVu Sans, sans-serif" font-size="22" letter-spacing="1.2" fill="#AEBBCD">github.com/espositopasquale/lingo</text>'''
    return document("social-card", 1200, 630, body, "Lingo repository social sharing card.")


def preview(assets):
    body = '''<rect width="1600" height="1570" fill="#E9EEF4"/>
<text x="64" y="75" font-family="DejaVu Sans, sans-serif" font-size="34" font-weight="bold" fill="#101B2D">Lingo / Brand assets</text>
<text x="64" y="110" font-family="DejaVu Sans, sans-serif" font-size="18" fill="#526275">One identity. Formats for documentation, apps, print and sharing.</text>'''
    cards = [
        ("wordmark-on-light.svg", "01 / Wordmark · light backgrounds", WHITE),
        ("wordmark-on-dark.svg", "02 / Wordmark · dark backgrounds", NAVY),
        ("logo-light.svg", "03 / Light lockup", WHITE),
        ("logo-dark.svg", "04 / Dark lockup", NAVY),
        ("icon.svg", "05 / App icon · full name", "#DAE2EC"),
        ("monogram.svg", "06 / Compact mark · avatar and favicon", "#DAE2EC"),
        ("mono-dark.svg", "07 / Single color · dark ink", WHITE),
        ("mono-light.svg", "08 / Single color · white ink", NAVY),
        ("social-card.svg", "09 / Social sharing · 1200 × 630", NAVY),
        ("favicon.svg", "10 / Favicon · small-size previews", "#DAE2EC"),
    ]
    for index, (filename, label, bg) in enumerate(cards):
        x = 64 + (index % 2) * 760
        y = 150 + (index // 2) * 278
        body += f'<rect x="{x}" y="{y}" width="712" height="220" rx="18" fill="{bg}"/>'
        svg = assets[filename]
        if filename in ("icon.svg", "monogram.svg"):
            sx, sy, sw, sh = x + 266, y + 20, 180, 180
        elif filename == "favicon.svg":
            for offset, size in ((140, 16), (200, 32), (280, 48), (390, 96)):
                start = svg.index(">")
                nested = svg[start + 1:svg.rindex("</svg>")]
                nested = nested.replace("favicon-", f"preview-{size}-favicon-")
                body += f'<svg x="{x + offset}" y="{y + (220-size)//2}" width="{size}" height="{size}" viewBox="0 0 512 512">{nested}</svg>'
            body += f'<text x="{x}" y="{y+252}" font-family="DejaVu Sans, sans-serif" font-size="16" fill="#33465D">{html.escape(label)}</text>'
            continue
        elif filename == "social-card.svg":
            sx, sy, sw, sh = x + 160, y + 7, 392, 206
        else:
            sx, sy, sw, sh = x + 42, y + 6, 628, 208
        start = svg.index(">")
        viewbox = svg.split('viewBox="', 1)[1].split('"', 1)[0]
        nested = svg[start + 1:svg.rindex("</svg>")]
        body += f'<svg x="{sx}" y="{sy}" width="{sw}" height="{sh}" viewBox="{viewbox}">{nested}</svg>'
        body += f'<text x="{x}" y="{y+252}" font-family="DejaVu Sans, sans-serif" font-size="16" fill="#33465D">{html.escape(label)}</text>'
    return document("preview", 1600, 1570, body, "Contact sheet of all Lingo logo variants and favicon sizes.")


class Rectangle(C.Structure):
    _fields_ = [(name, C.c_double) for name in ("x", "y", "width", "height")]


class Renderer:
    """Thin bindings to the existing system SVG renderer; no image manipulation."""

    def __init__(self):
        self.rsvg = C.CDLL(self.library("rsvg-2"))
        self.cairo = C.CDLL(self.library("cairo"))
        self.gobject = C.CDLL(self.library("gobject-2.0"))
        self.glib = C.CDLL(self.library("glib-2.0"))
        self.bind(self.rsvg, "rsvg_handle_new_from_data", [C.c_char_p, C.c_size_t, C.POINTER(C.c_void_p)], C.c_void_p)
        self.bind(self.rsvg, "rsvg_handle_render_document", [C.c_void_p, C.c_void_p, C.POINTER(Rectangle), C.POINTER(C.c_void_p)], C.c_int)
        self.bind(self.cairo, "cairo_image_surface_create", [C.c_int, C.c_int, C.c_int], C.c_void_p)
        self.bind(self.cairo, "cairo_surface_status", [C.c_void_p], C.c_int)
        self.bind(self.cairo, "cairo_create", [C.c_void_p], C.c_void_p)
        self.bind(self.cairo, "cairo_status", [C.c_void_p], C.c_int)
        self.bind(self.cairo, "cairo_surface_write_to_png", [C.c_void_p, C.c_char_p], C.c_int)
        self.bind(self.cairo, "cairo_destroy", [C.c_void_p], None)
        self.bind(self.cairo, "cairo_surface_destroy", [C.c_void_p], None)
        self.bind(self.gobject, "g_object_unref", [C.c_void_p], None)
        self.bind(self.glib, "g_error_free", [C.c_void_p], None)

    @staticmethod
    def library(name):
        path = ctypes.util.find_library(name)
        if not path:
            raise RuntimeError(f"Missing system library {name}; install librsvg/Cairo or use --svg-only")
        return path

    @staticmethod
    def bind(library, name, arguments, result):
        function = getattr(library, name)
        function.argtypes, function.restype = arguments, result

    def render(self, source, target, width, height):
        data = source.encode("utf-8")
        error = C.c_void_p()
        handle = self.rsvg.rsvg_handle_new_from_data(data, len(data), C.byref(error))
        if not handle:
            if error:
                self.glib.g_error_free(error)
            raise RuntimeError(f"Cannot parse SVG for {target.name}")
        surface = self.cairo.cairo_image_surface_create(0, width, height)
        context = self.cairo.cairo_create(surface)
        try:
            if self.cairo.cairo_surface_status(surface) or self.cairo.cairo_status(context):
                raise RuntimeError(f"Cannot create render surface for {target.name}")
            viewport = Rectangle(0, 0, width, height)
            if not self.rsvg.rsvg_handle_render_document(handle, context, C.byref(viewport), C.byref(error)):
                raise RuntimeError(f"Cannot render {target.name}")
            if self.cairo.cairo_surface_write_to_png(surface, str(target).encode()):
                raise RuntimeError(f"Cannot write {target}")
        finally:
            if error:
                self.glib.g_error_free(error)
            self.cairo.cairo_destroy(context)
            self.cairo.cairo_surface_destroy(surface)
            self.gobject.g_object_unref(handle)


def write_ico(target, entries):
    frames = [path.read_bytes() for _, path in entries]
    header = struct.pack("<HHH", 0, 1, len(entries))
    offset = 6 + 16 * len(entries)
    directory = bytearray()
    for (size, _), data in zip(entries, frames):
        dimension = 0 if size == 256 else size
        directory.extend(struct.pack("<BBBBHHII", dimension, dimension, 0, 0, 1, 32, len(data), offset))
        offset += len(data)
    target.write_bytes(header + directory + b"".join(frames))


def write_archive():
    target = ROOT.parent / "lingo-brand-kit.zip"
    files = {ROOT / "README.md", ROOT / "generate.py"}
    for extension in ("*.svg", "*.png", "*.ico"):
        files.update(ROOT.glob(extension))
    files.update(ROOT.parent / name for name in ("lingo-icon.png", "lingo-icon.prompt.txt"))
    with zipfile.ZipFile(target, "w", compression=zipfile.ZIP_DEFLATED) as archive:
        for path in sorted(files):
            if not path.is_file():
                continue
            info = zipfile.ZipInfo(path.relative_to(ROOT.parent).as_posix())
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o100644 << 16
            archive.writestr(info, path.read_bytes())
    print(f"Packaged {target.name}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--svg-only", action="store_true")
    args = parser.parse_args()
    assets = {
        "wordmark-on-dark.svg": horizontal("wordmark-on-dark"),
        "wordmark-on-light.svg": horizontal("wordmark-on-light", light=True),
        "logo-dark.svg": horizontal("logo-dark", background=NAVY),
        "logo-light.svg": horizontal("logo-light", light=True, background=WHITE),
        "icon.svg": app_icon("icon"),
        "monogram.svg": monogram("monogram"),
        "mono-dark.svg": horizontal("mono-dark", mono=NAVY),
        "mono-light.svg": horizontal("mono-light", mono=WHITE),
        "social-card.svg": social_card(),
        "favicon.svg": monogram("favicon"),
    }
    assets["preview.svg"] = preview(assets)
    for name, source in assets.items():
        (ROOT / name).write_text(source, encoding="utf-8")
    print(f"Generated {len(assets)} SVG files")
    if args.svg_only:
        return
    renderer = Renderer()
    for name, source in assets.items():
        if name == "favicon.svg":
            continue
        if name == "preview.svg":
            width, height = 1600, 1570
        elif name == "social-card.svg":
            width, height = 1200, 630
        elif name in ("icon.svg", "monogram.svg"):
            width = height = 512
        else:
            width, height = 1200, 500
        renderer.render(source, ROOT / name.replace(".svg", ".png"), width, height)
    renderer.render(assets["icon.svg"], ROOT / "icon-1024.png", 1024, 1024)
    renderer.render(monogram("apple-touch", rounded=False), ROOT / "apple-touch-icon.png", 180, 180)
    entries = []
    for size in (16, 32, 48, 64, 128, 256):
        path = ROOT / f"favicon-{size}.png"
        renderer.render(assets["favicon.svg"], path, size, size)
        if size in (16, 32, 48, 256):
            entries.append((size, path))
    write_ico(ROOT / "favicon.ico", entries)
    print("Generated PNG exports, apple-touch icon, and multi-size ICO")
    write_archive()


if __name__ == "__main__":
    main()
