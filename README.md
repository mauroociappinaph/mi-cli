# mi-cli

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/tuusuario/mi-cli?include_prereleases)](https://github.com/tuusuario/mi-cli/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/tuusuario/mi-cli/ci.yml?branch=main)](https://github.com/tuusuario/mi-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/tuusuario/mi-cli)](https://goreportcard.com/report/github.com/tuusuario/mi-cli)
[![License](https://img.shields.io/github/license/tuusuario/mi-cli)](LICENSE)
[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/tuusuario/mi-cli)

> Un CLI moderno y rápido construido con **Go + Cobra + Viper** — binario único, multiplataforma, listo para producción.

---

## ✨ Características

- ⚡ **Startup instantáneo** (~2-5ms) — sin runtime, binario estático
- 🌍 **Multiplataforma** — Linux, macOS, Windows (amd64/arm64)
- ⚙️ **Configuración flexible** — flags, env vars (`MI_CLI_*`), archivo YAML
- 📦 **Distribución simple** — un solo binario, Homebrew, Scoop, Winget, APT, RPM, Docker
- 🔄 **Releases automáticos** — GoReleaser + GitHub Actions (tags → artifacts + changelog)
- 🧪 **CI completo** — test (race detector), lint (golangci-lint), build multi-OS
- 📝 **Open source ready** — MIT license, contributing guide, code of conduct

---

## 🚀 Instalación

### Binario directo (recomendado)

```bash
# Linux/macOS
curl -fsSL https://github.com/tuusuario/mi-cli/releases/latest/download/mi-cli_Linux_x86_64.tar.gz | tar xz
sudo mv mi-cli /usr/local/bin/

# macOS (Apple Silicon)
curl -fsSL https://github.com/tuusuario/mi-cli/releases/latest/download/mi-cli_Darwin_arm64.tar.gz | tar xz
sudo mv mi-cli /usr/local/bin/

# Windows (PowerShell)
irm https://github.com/tuusuario/mi-cli/releases/latest/download/mi-cli_Windows_x86_64.zip | tar xz
# mover mi-cli.exe a tu PATH
```

### Gestores de paquetes

```bash
# Homebrew (macOS/Linux)
brew tap tuusuario/tap && brew install mi-cli

# Scoop (Windows)
scoop bucket add tuusuario https://github.com/tuusuario/scoop-bucket
scoop install mi-cli

# Winget (Windows)
winget install tuusuario.mi-cli

# Docker
docker run --rm ghcr.io/tuusuario/mi-cli:latest version
```

### Desde fuente

```bash
git clone https://github.com/tuusuario/mi-cli
cd mi-cli
make install
```

---

## 📖 Uso

```bash
mi-cli [command] [flags]

Commands:
  greet [name]    Saluda a alguien
  version         Muestra la versión
  help            Ayuda sobre cualquier comando

Flags:
  -c, --config string   Archivo de configuración (default: $HOME/.mi-cli.yaml)
  -h, --help            Ayuda
  -v, --verbose         Salida verbosa
  -o, --output string   Formato de salida (text|json) (default "text")
```

### Ejemplos

```bash
# Saludo simple
mi-cli greet
# ¡Hola, mundo! 👋

# Saludo personalizado
mi-cli greet Mauro
# ¡Hola, Mauro! 👋

# Output JSON para scripting
mi-cli greet --output json
# {"greeting": "¡Hola, mundo! 👋"}

# Configuración via environment variable
MI_CLI_VERBOSE=true mi-cli greet
# Usando config: /home/user/.mi-cli.yaml
# ¡Hola, mundo! 👋
```

### Archivo de configuración (`~/.mi-cli.yaml`)

```yaml
verbose: true
output: "json"
lang: "es"
```

---

## 🛠 Desarrollo

### Requisitos

- Go 1.23+
- Make
- golangci-lint (para linting): `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- goreleaser (para releases): `go install github.com/goreleaser/goreleaser@latest`

### Comandos útiles

```bash
make help        # Lista todos los targets
make build       # Compila para plataforma actual
make test        # Tests con race detector + coverage
make lint        # Ejecuta golangci-lint
make fmt         # Formatea con gofmt
make check       # Pipeline completo (fmt + vet + lint + test)
make snapshot    # Build snapshot con goreleaser (sin publicar)
make release     # Release completo (requiere tag git)
make dev ARGS="greet mundo"  # Build y ejecuta
```

### Estructura del proyecto

```
.
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Comando raíz, flags globales, viper
│   ├── version.go          # Subcomando version
│   └── greet.go            # Ejemplo de subcomando
├── internal/               # Lógica privada (no exportable)
├── pkg/                    # Código reusable por terceros
├── .goreleaser.yml         # Configuración de releases
├── Makefile                # Targets de desarrollo
├── go.mod / go.sum
└── .github/
    ├── workflows/ci.yml    # CI/CD pipeline
    └── dependabot.yml      # Updates automáticos de deps
```

---

## 🧪 Testing

```bash
# Tests unitarios con coverage
make test

# Ver reporte HTML
make cover

# Tests rápidos (sin race detector)
make test-short
```

---

## 📦 Release

Los releases son **automáticos** al pushear un tag:

```bash
# Crear y pushear tag
git tag v1.0.0
git push origin v1.0.0
```

Esto dispara:
1. ✅ CI pipeline (test, lint, build multi-OS)
2. 🏗 GoReleaser builda para todas las plataformas
3. 📝 Genera changelog desde commits convencionales
4. 📤 Publica artifacts en GitHub Releases
5. 🐳 Pushea imagen Docker a GHCR
6. 🍺 Actualiza Homebrew tap, Scoop bucket, Winget

### Versionado

Usamos [SemVer](https://semver.org/) + [Conventional Commits](https://www.conventionalcommits.org/):

| Tipo de commit | Version bump |
|----------------|--------------|
| `fix:`         | PATCH        |
| `feat:`        | MINOR        |
| `BREAKING CHANGE:` | MAJOR   |

---

## 🤝 Contribuir

1. Fork del repo
2. Crea tu feature branch (`git checkout -b feat/amazing-feature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push al branch (`git push origin feat/amazing-feature`)
5. Abre un Pull Request

### Guidelines

- Código: `make check` debe pasar
- Tests: añadir tests para nueva funcionalidad
- Commits: [Conventional Commits](https://www.conventionalcommits.org/)
- Docs: actualizar README si cambia la CLI

---

## 📄 Licencia

MIT License — ver [LICENSE](LICENSE) para detalles.

---

## 🙏 Créditos

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Viper](https://github.com/spf13/viper) — Configuración
- [GoReleaser](https://goreleaser.com/) — Releases automáticos
- [golangci-lint](https://golangci-lint.run/) — Linting

---

<div align="center">
  <sub>Hecho con ❤️ en Go • <a href="https://github.com/tuusuario/mi-cli/issues">Reportar bug</a> • <a href="https://github.com/tuusuario/mi-cli/issues">Pedir feature</a></sub>
</div>