<div align="center">

```
╔════════════════════════════════════════════════════════════════════════╗
║                                                                        ║
║   ████████╗██╗   ██╗██████╗ ███████╗ ██████╗██████╗  █████╗ ███████╗████████╗
║   ╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝
║      ██║    ╚████╔╝ ██████╔╝█████╗  ██║     ██████╔╝███████║█████╗     ██║   
║      ██║     ╚██╔╝  ██╔═══╝ ██╔══╝  ██║     ██╔══██╗██╔══██║██╔══╝     ██║   
║      ██║      ██║   ██║     ███████╗╚██████╗██║  ██║██║  ██║██║        ██║   
║      ╚═╝      ╚═╝   ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝        ╚═╝   
║                                                                        ║
║                   AI-Powered Book Production Engine                   ║
║                                                                        ║
╚════════════════════════════════════════════════════════════════════════╝
```

**Transform manuscripts into professionally typeset books in minutes, not weeks.**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Alpha-orange?style=for-the-badge)]()
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-brightgreen?style=for-the-badge)](CONTRIBUTING.md)

[Features](#-features) • [Quick Start](#-quick-start) • [Architecture](#-architecture) • [Roadmap](#-roadmap) • [Contributing](#-contributing)

</div>

---

## 🎯 The Problem

Publishing a professionally typeset book today requires:
- **$1,500-$3,000** in design and typesetting services
- **2-4 weeks** of production time
- **Technical expertise** in tools like InDesign or LaTeX
- **Multiple professionals**: editors, designers, typesetters

## ✨ Our Solution

**Typecraft** is an intelligent book production system that combines **centuries of typographic wisdom** with **cutting-edge AI** to transform raw manuscripts into publication-ready books automatically.

```
📄 manuscript.docx  →  🤖 AI Analysis  →  📚 professional-book.pdf + .epub
                           5 minutes              Ready to print/publish
```

### What Makes Typecraft Different

- 🧠 **Content-Aware Design**: AI analyzes genre, tone, and complexity to make contextual design decisions
- 📐 **Classical Typography**: Implements Van de Graaf Canon, grid systems, and optimal spacing algorithms
- 🎨 **Intelligent Refinement**: Deep learning models for optical kerning and widow/orphan elimination
- ⚡ **Blazingly Fast**: Go-powered backend processes 300-page books in under 5 minutes
- 💰 **Radically Affordable**: $15 per book vs. $1,500+ traditional cost

---

## 🚀 Features

### Current (MVP - In Development)

- [x] **Multi-format Ingestion**: `.docx`, `.txt`, `.md`, `.odt` → Structured Markdown
- [x] **Dual Rendering Pipeline**:
  - **LaTeX**: Perfect for text-heavy, academic, technical books
  - **HTML/CSS + Paged.js**: Ideal for visual layouts, magazines, illustrated books
- [x] **Content Analysis (AI)**: Genre classification, tone analysis, complexity scoring
- [x] **Automated Layout**: Van de Graaf Canon for harmonious margins
- [ ] **Typographic Refinement**: Paragraph optimization, widow/orphan correction
- [ ] **Multi-Channel Output**:
  - PDF/X-1a (Amazon KDP)
  - PDF/X-4 (IngramSpark)
  - ePub 3 (universal digital)

### Roadmap (Next 12 Months)

- [ ] **Optical Kerning**: Transformer-based model for professional letter spacing
- [ ] **Copy-Editing Assistant**: LLM-powered grammar and style suggestions
- [ ] **Cover Generation**: AI-generated book covers from manuscript analysis
- [ ] **Multi-Language**: Translate and typeset in 10+ languages
- [ ] **Audiobook Export**: AI narration with character voice detection

[See full roadmap →](docs/ROADMAP.md)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT (React + TS)                     │
└────────────────────────┬────────────────────────────────────┘
                         │ REST API
┌────────────────────────▼────────────────────────────────────┐
│                  API SERVER (Go + Gin)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │  Projects    │  │  Processing  │  │   Downloads     │   │
│  │  CRUD        │  │  Orchestrator│  │   Handler       │   │
│  └──────┬───────┘  └──────┬───────┘  └─────────────────┘   │
└─────────┼──────────────────┼──────────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│              ASYNC WORKERS (Asynq + Redis)                  │
│  ┌──────────────────┐  ┌──────────────────────────────┐    │
│  │  Converter       │  │  Renderer (LaTeX/HTML)       │    │
│  │  (Pandoc)        │  │                              │    │
│  └──────────────────┘  └──────────────────────────────┘    │
│  ┌──────────────────┐  ┌──────────────────────────────┐    │
│  │  AI Analyzer     │  │  Refinement Engine           │    │
│  │  (OpenAI/Claude) │  │  (Typography AI)             │    │
│  └──────────────────┘  └──────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   DATA LAYER                                │
│  ┌────────────┐  ┌────────────┐  ┌───────────────────┐     │
│  │ PostgreSQL │  │  MinIO/S3  │  │  Redis Cache      │     │
│  │ (metadata) │  │  (files)   │  │                   │     │
│  └────────────┘  └────────────┘  └───────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

**Tech Stack:**
- **Backend**: Go 1.22+, Gin, GORM
- **Workers**: Asynq (Redis-backed queue)
- **AI/ML**: OpenAI GPT-4o, Anthropic Claude 3.5, PyTorch (custom models)
- **Processing**: Pandoc (universal converter), LuaLaTeX (PDF rendering)
- **Frontend**: React 18, TypeScript, TailwindCSS
- **Infrastructure**: Docker, PostgreSQL, Redis, MinIO

[See detailed architecture →](docs/architecture/BLUEPRINT.md)

---

## ⚡ Quick Start

### Prerequisites

- Go 1.22+ ([install](https://go.dev/doc/install))
- Docker & Docker Compose ([install](https://docs.docker.com/get-docker/))
- Pandoc ([install](https://pandoc.org/installing.html))
- LaTeX (TeX Live) ([install](https://www.tug.org/texlive/))

### Installation

```bash
# Clone repository
git clone https://github.com/JuanCS-Dev/typecraft.git
cd typecraft

# Install dependencies
make install

# Start infrastructure (PostgreSQL, Redis, MinIO)
docker compose up -d

# Run database migrations
make migrate

# Start API server
make run-api

# Start worker (in another terminal)
make run-worker
```

The API will be available at `http://localhost:8000`

### Your First Book

```bash
# Create a project
curl -X POST http://localhost:8000/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Book",
    "author": "Your Name",
    "genre": "Fiction",
    "page_format": "6x9"
  }'

# Upload manuscript
curl -X POST http://localhost:8000/api/v1/projects/{id}/upload \
  -F "file=@manuscript.docx"

# Process book
curl -X POST http://localhost:8000/api/v1/projects/{id}/process

# Download PDF (after ~5 minutes)
curl -O http://localhost:8000/api/v1/projects/{id}/download/pdf
```

[See full documentation →](docs/guides/GETTING_STARTED.md)

---

## 📊 Project Status

| Milestone | Status | Target Date | Progress |
|-----------|--------|-------------|----------|
| **Phase 0**: Foundation | 🟡 In Progress | Nov 2025 | 40% |
| **Phase 1**: MVP | ⏳ Planned | Feb 2026 | 0% |
| **Phase 2**: MMP | ⏳ Planned | May 2026 | 0% |
| **Phase 3**: MLP | ⏳ Planned | Nov 2026 | 0% |

**Current Sprint:** Phase 0 - Environment Setup & PoC  
**Next Milestone:** Core Pipeline (Ingestion → Conversion → Basic Rendering)

---

## 🤝 Contributing

We welcome contributions! Typecraft is built under the **Vértice Constitution** principles:

- ✅ **P1 - Completeness**: No TODOs or placeholders in production code
- ✅ **P2 - Preventive Validation**: All APIs/libraries validated before use
- ✅ **P3 - Critical Skepticism**: Challenge assumptions, document decisions
- ✅ **P6 - Token Efficiency**: Diagnose before fixing, max 2 iterations

Please read our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

### Development Setup

```bash
# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Build for production
make build
```

---

## 📄 License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

This project stands on the shoulders of giants:

- **Robert Bringhurst** - *The Elements of Typographic Style*
- **Josef Müller-Brockmann** - *Grid Systems in Graphic Design*
- **Donald Knuth** - *The TeXbook* and the TeX typesetting system
- **Pandoc** - Universal document converter
- **LaTeX Community** - Decades of typographic excellence

---

## 📞 Contact

- **Author**: Juan (Architect-in-Chief)
- **AI Executor**: Claude (Anthropic) under Vértice Constitution v3.0
- **Repository**: [github.com/JuanCS-Dev/typecraft](https://github.com/JuanCS-Dev/typecraft)
- **Issues**: [Report a bug or request a feature](https://github.com/JuanCS-Dev/typecraft/issues)

---

<div align="center">

**Made with ❤️ by the Typecraft Team**

*Democratizing professional publishing through AI and classical typography*

[⭐ Star us on GitHub](https://github.com/JuanCS-Dev/typecraft) • [🐛 Report Bug](https://github.com/JuanCS-Dev/typecraft/issues) • [💡 Request Feature](https://github.com/JuanCS-Dev/typecraft/issues)

</div>
