# AudeticLink Installer

Welcome to the AudeticLink Installer! This tool helps you set up AudeticLink on your raspberry pi. Follow the steps below to download and run the installer.

---

## Quick Start Guide

### Prerequisites

- **System**: Raspberry Pi bookwork 64-bit ARM64

---

All in one go

```bash
curl -L -o audeticlink https://github.com/silvabyte/AudeticLinkInstaller/releases/download/v0.0.7/link_arm64 && \
chmod +x audeticlink
```

### Step by Step Installation

### Step 1: Download the Installer

```bash
curl -L -o audeticlink https://github.com/silvabyte/AudeticLinkInstaller/releases/download/v0.0.7/link_arm64
```

### Step 2: Make Executable

```bash
chmod +x audeticlink
```

### Step 3: Run the Installer

```bash
sudo ./audeticlink install rpi02w --github-token="<insert_gh_token>"
```
