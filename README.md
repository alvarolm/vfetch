![vfetch-logo](https://github.com/user-attachments/assets/ea4eab93-655c-4340-9d68-d0f7007f8ce8)

<svg
   width="400"
   height="100"
   viewBox="0 0 105.83333 26.458333"
   version="1.1"
   id="svg5"
   xml:space="preserve"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:svg="http://www.w3.org/2000/svg"><defs
     id="defs2" /><g
     id="layer1"><rect
       style="opacity:0.9;fill:#000000;fill-opacity:1;stroke-width:0.2;stroke-linecap:round"
       id="rect1134"
       width="105.83333"
       height="26.458332"
       x="0"
       y="0"
       ry="0.24815215" /><g
       id="g2216"
       transform="translate(-3.0828817)"
       style="font-size:15.1903px;-inkscape-font-specification:sans-serif;fill:#f9f9f9;stroke-width:0.956872;stroke-linecap:round"><path
         d="m 35.573636,7.8973705 4.238093,10.6332105 h 1.853217 L 45.963801,7.8973705 H 44.095394 L 40.783909,16.540651 37.502804,7.8973705 Z"
         style="opacity:1"
         id="path2161" /><path
         d="M 47.593975,7.8973705 V 18.530581 h 1.792455 v -4.359617 h 4.830516 V 12.575983 H 49.38643 V 9.5075423 h 5.331796 l 0.01519,-1.6101718 z"
         style="font-family:Montserrat;-inkscape-font-specification:Montserrat;opacity:1"
         id="path2158" /><path
         d="M 56.515496,7.8973705 V 18.530581 h 7.777434 v -1.610172 h -5.984978 v -2.931728 h 5.195082 V 12.393699 H 58.307952 V 9.5075423 h 5.802694 V 7.8973705 Z"
         style="font-family:Montserrat;-inkscape-font-specification:Montserrat;opacity:1"
         id="path2155" /><path
         d="m 65.406629,7.8973705 v 1.6101718 h 3.357056 v 9.0230387 h 1.792456 V 9.5075423 h 3.372246 V 7.8973705 Z"
         style="font-family:Montserrat;-inkscape-font-specification:Montserrat;opacity:1"
         id="path2152" /><path
         d="m 82.713194,10.798718 1.048131,-1.306366 c -1.032941,-1.0177501 -2.53678,-1.6405524 -4.04062,-1.6405524 -3.144392,0 -5.59003,2.3241164 -5.59003,5.3469854 0,3.053251 2.415257,5.407747 5.529269,5.407747 1.503839,0 3.03806,-0.683563 4.116571,-1.746884 l -1.063321,-1.184844 c -0.805086,0.774705 -1.898787,1.260795 -2.962108,1.260795 -2.126642,0 -3.797575,-1.655743 -3.797575,-3.752004 0,-2.096261 1.670933,-3.7368139 3.797575,-3.7368139 1.078511,0 2.187403,0.5164702 2.962108,1.3519369 z"
         style="font-family:Montserrat;-inkscape-font-specification:Montserrat;opacity:1"
         id="path2149" /><path
         d="M 85.786446,7.8973705 V 18.530581 h 1.792456 v -4.405187 h 5.635601 v 4.405187 h 1.792456 V 7.8973705 H 93.214503 V 12.515222 H 87.578902 V 7.8973705 Z"
         style="font-family:Montserrat;-inkscape-font-specification:Montserrat;opacity:1"
         id="path2146" /></g><path
       d="M 20.338807,17.646577 17.021504,13.936744 H 19.59585 V 8.9836972 h 1.485915 v 4.9530468 l 2.623381,0.02264 z M 14.39515,18.889792 H 26.282464 V 7.5685403 H 14.39515 Z m -1.485913,1.415156 H 27.768378 V 6.1533852 H 12.909237 Z"
       id="download-[#1458]"
       style="opacity:0.59156916;fill:#f9f9f9;fill-rule:evenodd;stroke:none;stroke-width:0.707578" /></g></svg>



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
# Download, verify and install esbuild
vfetch -config my-tools.json
```

Example `my-tools.json`:
```json
{
  "output-dir": "/home/user/tools",
  "bins-dir": "/home/user/.bin",
  "fetch": [
    {
      "name": "esbuild",
      "url": "https://registry.npmjs.org/@esbuild/linux-x64/-/linux-x64-$VERSION.tgz",
      "version": "0.25.10",
      "hash": "sha256:25a7b968b8e5172baaa8f44f91b71c1d2d7e760042c691f22ab59527d870d145",
      "bin-file": "/package/bin/esbuild",
      "extract": true
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

[LICENSE](LICENSE) - Use it freely, but remember: **you** are responsible for verifying what you download.

---

**Remember: Security is not a feature you can install - it's a practice you must maintain.**
