# 📦 sack

> A lightweight, blazing-fast file archiver and Go library with a clean binary specification.

`sack` is both a CLI tool and a Go package (`lib`) for packing files and directories into `.sack` archives. It utilizes a **Footer TOC (Table of Contents)** architecture-similar to ZIP-allowing single-pass streaming without heavy RAM allocation.

---

## ⚡ Features

* **Zero-Allocation Streaming:** Pack files of any size without loading entire payloads into memory.
* **Footer TOC Architecture:** Table of contents is appended at the end of the file, storing precise `Offset`, `Size`, and POSIX `Mode` values.
* **POSIX Support:** Preserves file permissions (`chmod`) and nested directory structures upon extraction.

---

## 📐 The `.sack` Binary Format

```text
0xEE [ File 1 Data ] [ File 2 Data ] ... [ TOC Payload ] [ TOC Offset (int64) ] [ "SACK" Magic Bytes (4B) ]
```

Table of Contents (TOC) Payload:
* Count (uint32) — Total number of archived files.
* File Entries (repeated Count times):
    * NameLength (uint16) — Length of the file path string in bytes.
    * Size (int64) — Uncompressed file size in bytes.
    * Offset (int64) — Absolute byte offset of file payload from archive start.
    * Mode (uint32) — POSIX file mode and permissions flags.
    * Name ([]byte) — UTF-8 encoded relative path (NameLength bytes).