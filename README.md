# vfetch

> **Simple, secure downloads without the complexity of package managers**

`vfetch` is a lightweight tool that downloads, verifies, and organizes files with cryptographic integrity checking. It bridges the gap between insecure `curl`/`wget` downloads and heavyweight package managers, making you conscious of security while keeping things simple.

## Why vfetch?

### The Problem with Current Approaches

**Package Managers (npm, etc.)**
- Heavy overhead and complex dependency trees
- Lock you into specific ecosystems
- Abstract away verification, making you unaware of security
- Require learning package-specific tooling

**Raw Downloads (curl, wget)**
- No integrity verification by default
- Easy to forget or skip checksum validation
- Manual hash checking is error-prone
- No organized file management

### The vfetch Philosophy

**Security by Design, Not by Accident**
- Forces you to provide checksums for every download
- Supports multiple hash algorithms (SHA256, SHA512, SHA3, BLAKE2b, BLAKE2s, BLAKE3)
- Makes verification failure explicit and loud
- **Puts you in control** - you vet the checksums, not some package registry

**Simplicity Without Compromise**
- Single binary, no dependencies
- Human-readable JSON configuration
- Predictable file organization
- No hidden magic or complex dependency resolution

**Awareness Through Responsibility**
- Every download requires a hash - no shortcuts
- You must consciously verify checksums from trusted sources
- Builds security habits through explicit verification requirements
- Makes the cost of trust visible and intentional

## Quick Start

1. **Download vfetch** 
2. **Create a config file** with your downloads and their checksums
3. **Run vfetch** and get verified, organized files

```bash
# Download and verify Go 1.21.6
vfetch -config my-tools.json
```

Example `my-tools.json`:
```json
{
  "output-dir": "/opt/tools",
  "bins-dir": "/usr/local/bin",
  "fetch": [
    {
      "name": "go",
      "url": "https://go.dev/dl/go$version.linux-amd64.tar.gz",
      "version": "1.21.6",
      "hash": "sha256:3f934f40ac360b9c01f616a9aa1796d227d8b0328bf64cb045c7b8c4ee9caea4",
      "extract": true,
      "bin-file": "go/bin/go"
    }
  ]
}
```

## Key Features

### **Mandatory Verification**
- **No downloads without checksums** - vfetch refuses to proceed without proper hashes
- **Multiple hash algorithms** supported for maximum compatibility
- **Fail-fast verification** - stops immediately on hash mismatches

### **Smart File Handling**
- **Automatic extraction** for ZIP, TAR, TAR.GZ, and GZIP archives
- **Binary symlink creation** for executable files
- **Organized output** with predictable directory structures

### **Flexible Configuration**
- **Version placeholders** in URLs (`$version` → actual version)
- **Per-item overrides** for output and binary directories
- **Documentation tracking** with optional URL fields for license, source, etc.

### **Zero Dependencies**
- Single statically-linked binary
- No runtime dependencies or package ecosystems
- Works anywhere Go runs

## Why Checksums Matter

When you download files with `curl` or `wget`, you're trusting:
- The network connection isn't compromised
- The server hasn't been hacked
- The file wasn't modified in transit
- DNS hasn't been hijacked

**vfetch makes this explicit** by requiring you to:
1. **Find official checksums** from the project's trusted sources
2. **Verify them yourself** against multiple sources when possible
3. **Take responsibility** for the integrity of what you download

This isn't paranoia - it's basic operational security that should be standard practice.

## Installation

### Download Binary
Check the [releases page](https://github.com/alvarolm/vfetch/releases) for pre-built binaries.

### Using Go Install
```bash
go install github.com/alvarolm/vfetch@latest
```

### From Source
```bash
git clone https://github.com/alvarolm/vfetch
cd vfetch
go build -o vfetch .
cp ./vfetch /usr/local/bin
```

**Remember to verify the checksum of vfetch itself!**

## Configuration Reference

See [example-config.json](example-config.json) for a comprehensive configuration example with all available options.

### Required Fields
- `name`: Human-readable identifier
- `url`: Download URL (supports `$version` placeholders)
- `version`: Version identifier
- `hash` or `hashes`: Cryptographic verification

### Optional Fields
- `extract`: Extract archives automatically
- `bin-file`: Create executable symlinks
- `output-dir`: Override global output directory
- `bin-dir`: Override global binary directory

## Examples

### Simple Binary Download
```json
{
  "name": "jq",
  "url": "https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64",
  "version": "1.6",
  "hash": "sha256:af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44",
  "bin-file": true
}
```

### Archive with Extraction
```json
{
  "name": "node",
  "url": "https://nodejs.org/dist/v$version/node-v$version-linux-x64.tar.gz",
  "version": "18.17.0",
  "hash": "sha256:...actual-hash...",
  "extract": true,
  "bin-file": "node-v18.17.0-linux-x64/bin/node"
}
```

### Multiple Hash Verification
```json
{
  "name": "critical-tool",
  "url": "https://example.com/tool.tar.gz",
  "version": "2.1.0",
  "hashes": [
    "sha256:...",
    "sha512:..."
  ],
  "extract": true
}
```

## Security Best Practices

1. **Always verify checksums** from official project sources
2. **Cross-reference hashes** from multiple trusted sources when possible
3. **Use HTTPS URLs** for downloads
4. **Keep vfetch updated** to get the latest security improvements
5. **Review configurations** before running them
6. **Store configurations in version control** for audit trails

## Comparison

| Tool | Verification | Complexity | Ecosystem Lock-in | Security Awareness |
|------|-------------|------------|-------------------|-------------------|
| vfetch | ✅ Mandatory | 🟢 Low | ❌ None | ✅ High |
| npm/pip | ⚠️ Registry-based | 🔴 High | ✅ Heavy | ❌ Hidden |
| curl/wget | ❌ Manual/Optional | 🟢 Low | ❌ None | ⚠️ User-dependent |

## Contributing

vfetch is designed to stay simple and focused. When contributing:

1. **Maintain simplicity** - avoid feature creep
2. **Security first** - never compromise on verification requirements
3. **Explicit over implicit** - make security decisions visible
4. **Test thoroughly** - especially hash verification and file handling

## License

[MIT License](LICENSE) - Use it freely, but remember: **you** are responsible for verifying what you download.

---

**Remember: Security is not a feature you can install - it's a practice you must maintain.**