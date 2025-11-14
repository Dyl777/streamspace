#!/bin/bash
set -e

# StreamSpace Template Migration Script
# This script helps migrate templates to the external repository

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STREAMSPACE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TARGET_REPO="${1:-}"

usage() {
    cat << EOF
Usage: $0 <target-repository-path>

Migrates StreamSpace templates to an external repository.

Arguments:
    target-repository-path    Path to the streamspace-templates repository

Example:
    $0 /path/to/streamspace-templates

EOF
    exit 1
}

if [ -z "$TARGET_REPO" ]; then
    echo "Error: Target repository path required"
    usage
fi

if [ ! -d "$TARGET_REPO" ]; then
    echo "Error: Target repository does not exist: $TARGET_REPO"
    echo ""
    echo "Initialize it first:"
    echo "  mkdir -p $TARGET_REPO"
    echo "  cd $TARGET_REPO"
    echo "  git init"
    exit 1
fi

echo "StreamSpace Template Migration"
echo "==============================="
echo ""
echo "Source: $STREAMSPACE_ROOT"
echo "Target: $TARGET_REPO"
echo ""

# Create directory structure
echo "Creating directory structure..."
mkdir -p "$TARGET_REPO/templates"/{browsers,design,development,gaming,media,productivity,webtop}
mkdir -p "$TARGET_REPO/generated"
mkdir -p "$TARGET_REPO/icons"
mkdir -p "$TARGET_REPO/scripts"
mkdir -p "$TARGET_REPO/.github/workflows"

# Copy templates by category
echo ""
echo "Copying templates..."

# Browsers
echo "  Copying browsers..."
for template in brave chromium firefox librewolf; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/browsers/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Design
echo "  Copying design tools..."
for template in blender freecad gimp inkscape kicad krita; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/design/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Development
echo "  Copying development tools..."
for template in code-server github-desktop gitqlient; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/development/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Gaming
echo "  Copying gaming applications..."
for template in dolphin duckstation; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/gaming/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Media
echo "  Copying media applications..."
for template in audacity kdenlive; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/media/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Productivity
echo "  Copying productivity applications..."
for template in calligra libreoffice; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/productivity/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Webtop
echo "  Copying webtop environments..."
for template in webtop-alpine-i3 webtop-ubuntu-kde webtop-ubuntu-xfce; do
    if [ -f "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" ]; then
        cp "$STREAMSPACE_ROOT/manifests/templates/${template}.yaml" "$TARGET_REPO/templates/webtop/"
        echo "    ✓ ${template}.yaml"
    fi
done

# Copy generated templates
if [ -d "$STREAMSPACE_ROOT/manifests/templates-generated" ]; then
    echo "  Copying generated templates..."
    cp -r "$STREAMSPACE_ROOT/manifests/templates-generated/"* "$TARGET_REPO/generated/" 2>/dev/null || true
fi

# Copy scripts
echo ""
echo "Copying scripts..."
if [ -f "$STREAMSPACE_ROOT/scripts/generate-templates.py" ]; then
    cp "$STREAMSPACE_ROOT/scripts/generate-templates.py" "$TARGET_REPO/scripts/"
    echo "  ✓ generate-templates.py"
fi

# Create validation script
echo ""
echo "Creating validation script..."
cat > "$TARGET_REPO/scripts/validate-templates.sh" << 'EOF'
#!/bin/bash
set -e

echo "Validating StreamSpace templates..."

ERRORS=0
WARNINGS=0

for file in templates/**/*.yaml generated/*.yaml; do
    if [ ! -f "$file" ]; then
        continue
    fi

    echo "Validating $file..."

    # Check for required fields
    if ! grep -q "apiVersion: stream.space/v1alpha1" "$file"; then
        echo "  ERROR: Missing or incorrect apiVersion in $file"
        ERRORS=$((ERRORS + 1))
    fi

    if ! grep -q "kind: Template" "$file"; then
        echo "  ERROR: Missing kind: Template in $file"
        ERRORS=$((ERRORS + 1))
    fi

    if ! grep -q "displayName:" "$file"; then
        echo "  ERROR: Missing displayName in $file"
        ERRORS=$((ERRORS + 1))
    fi

    if ! grep -q "baseImage:" "$file"; then
        echo "  ERROR: Missing baseImage in $file"
        ERRORS=$((ERRORS + 1))
    fi

    # Check for recommended fields
    if ! grep -q "description:" "$file"; then
        echo "  WARNING: Missing description in $file"
        WARNINGS=$((WARNINGS + 1))
    fi

    if ! grep -q "category:" "$file"; then
        echo "  WARNING: Missing category in $file"
        WARNINGS=$((WARNINGS + 1))
    fi

    if ! grep -q "icon:" "$file"; then
        echo "  WARNING: Missing icon in $file"
        WARNINGS=$((WARNINGS + 1))
    fi

    echo "  ✓ $file validated"
done

echo ""
echo "Validation Summary:"
echo "  Errors: $ERRORS"
echo "  Warnings: $WARNINGS"

if [ $ERRORS -gt 0 ]; then
    echo ""
    echo "❌ Validation failed with $ERRORS errors"
    exit 1
else
    echo ""
    echo "✅ All templates validated successfully"
    if [ $WARNINGS -gt 0 ]; then
        echo "⚠️  $WARNINGS warnings (recommended fields missing)"
    fi
fi
EOF

chmod +x "$TARGET_REPO/scripts/validate-templates.sh"
echo "  ✓ validate-templates.sh"

# Create README
echo ""
echo "Creating README.md..."
cat > "$TARGET_REPO/README.md" << 'EOF'
# StreamSpace Templates

Official template repository for StreamSpace - Cloud-native desktop streaming platform.

## Overview

This repository contains application templates for StreamSpace sessions. Each template defines a containerized desktop application that can be streamed via web browser.

## Template Categories

- **Browsers**: Web browsers (Firefox, Chromium, Brave, etc.)
- **Design**: 3D modeling, graphic design, CAD applications
- **Development**: IDEs, code editors, git clients
- **Gaming**: Emulators and gaming applications
- **Media**: Audio/video editing software
- **Productivity**: Office suites and productivity tools
- **Webtop**: Full desktop environments

## Template Structure

Templates are Kubernetes Custom Resources (CRDs) with the following format:

```yaml
apiVersion: stream.space/v1alpha1
kind: Template
metadata:
  name: template-name
  namespace: workspaces
spec:
  displayName: "Display Name"
  description: "Detailed description"
  category: "Category Name"
  icon: "https://..."
  baseImage: "docker.io/image:tag"
  defaultResources:
    memory: 2Gi
    cpu: 1000m
  ports:
    - name: vnc
      containerPort: 3000
      protocol: TCP
  env: []
  volumeMounts: []
  kasmvnc:
    enabled: true
    port: 3000
  capabilities: []
  tags: []
```

## Usage

### Adding to StreamSpace

1. Navigate to **Repositories** in StreamSpace UI
2. Click **Add Repository**
3. Enter repository URL: `https://github.com/JoshuaAFerguson/streamspace-templates`
4. Select branch: `main`
5. Click **Add and Sync**

### Creating Templates

See [TEMPLATE_MIGRATION_GUIDE.md](../streamspace/TEMPLATE_MIGRATION_GUIDE.md) for template creation guidelines.

## Available Templates

### Browsers
- Brave Browser
- Chromium
- Firefox
- LibreWolf

### Design & Graphics
- Blender 3D
- FreeCAD
- GIMP
- Inkscape
- KiCAD
- Krita

### Development Tools
- Code Server (VS Code)
- GitHub Desktop
- GitQlient

### Gaming & Emulation
- Dolphin Emulator
- DuckStation

### Media & Audio
- Audacity
- Kdenlive

### Productivity & Office
- Calligra Suite
- LibreOffice

### Desktop Environments
- Webtop Alpine i3
- Webtop Ubuntu KDE
- Webtop Ubuntu XFCE

## Validation

Run the validation script to check all templates:

```bash
./scripts/validate-templates.sh
```

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add or modify templates
4. Run validation script
5. Submit a pull request

## License

MIT License - See LICENSE file.

Individual applications have their own licenses.
EOF

echo "  ✓ README.md"

# Create LICENSE
echo ""
echo "Creating LICENSE..."
cat > "$TARGET_REPO/LICENSE" << 'EOF'
MIT License

Copyright (c) 2024 StreamSpace

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

echo "  ✓ LICENSE"

# Create .gitignore
echo ""
echo "Creating .gitignore..."
cat > "$TARGET_REPO/.gitignore" << 'EOF'
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
env/
venv/
.venv/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Temporary files
*.tmp
*.bak
*.log
EOF

echo "  ✓ .gitignore"

# Run validation
echo ""
echo "Running validation..."
cd "$TARGET_REPO"
./scripts/validate-templates.sh

# Count templates
TEMPLATE_COUNT=$(find templates -name "*.yaml" | wc -l)
GENERATED_COUNT=$(find generated -name "*.yaml" 2>/dev/null | wc -l)

echo ""
echo "==============================="
echo "Migration Complete!"
echo "==============================="
echo ""
echo "Summary:"
echo "  Templates migrated: $TEMPLATE_COUNT"
echo "  Generated templates: $GENERATED_COUNT"
echo "  Target repository: $TARGET_REPO"
echo ""
echo "Next steps:"
echo "  1. Review migrated templates in $TARGET_REPO"
echo "  2. Initialize git repository:"
echo "     cd $TARGET_REPO"
echo "     git init"
echo "     git add ."
echo "     git commit -m 'Initial commit: StreamSpace templates'"
echo "  3. Add remote and push:"
echo "     git remote add origin https://github.com/JoshuaAFerguson/streamspace-templates.git"
echo "     git branch -M main"
echo "     git push -u origin main"
echo "  4. Add repository in StreamSpace UI"
echo ""
