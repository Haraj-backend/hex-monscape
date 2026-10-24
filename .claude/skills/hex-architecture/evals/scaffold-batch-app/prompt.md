---
tags: [hex-architecture]
max_turns: 60
timeout_seconds: 1200
allowed_tools: [Read, Glob, Grep, Skill, Write, Edit, Bash]
---

Tolong scaffold Go app baru pakai hex architecture Haraj style. App ini baca daftar URL RSS dari feeds.txt, tiap jam crawl item baru dan simpan ke data/items.jsonl, lalu tiap pagi jam 7 bikin digest (20 item terbaru) dan upload JSON-nya ke S3. Module path example.com/rssdigest. Tulis kode lengkapnya sampai bisa `go build`, plus test untuk service-nya.
