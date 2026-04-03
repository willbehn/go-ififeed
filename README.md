[![GitHub Release](https://img.shields.io/github/v/release/willbehn/go-ififeed)](https://github.com/willbehn/go-ififeed/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/willbehn/go-ififeed)](go.mod)
![License](https://img.shields.io/github/license/willbehn/go-ififeed)
# go-ififeed

**ififeed** er et verktøy for deg som er IFI-student og er lei av å
måtte klikke deg inn på Mine Studier eller semestersiden for hvert enkelt
emne bare for å sjekke nye beskjeder.

Med **ififeed** kan du hente alle beskjeder fra alle emner
du tar rett i terminalen! Ingen flere avbrekk fra den produktive progge-økta!

![Made with VHS](assets/demo.gif)

## Kom i gang

### 1. Last ned binary-filen

Gå til [Releases](https://github.com/willbehn/go-ififeed/releases) og last ned riktig binary-fil basert på plattform:

| Plattform | Fil |
|---|---|
| macOS (Apple Silicon) | `ififeed-darwin-arm64` |
| Linux | `ififeed-linux-amd64` |

Gjør filen kjørbar og flytt den til PATH:

```bash
chmod +x ififeed-darwin-arm64
mv ififeed-darwin-arm64 /usr/local/bin/ififeed
```

### 2. Legg til emnene du tar

Du kan legge til emner direkte fra terminalen, for eksempel:
```bash
ififeed add IN1000 v25 "Introduksjon til objektorientert programmering"
```

Eller rediger config-filen manuelt (opprettes automatisk første gang du kjører `ififeed`):

På macOS blir config filen lagt på:
```
~/Library/Application Support/ififeed/courses.yaml
```

For andre OS: Varierer, se hva os.UserConfigDir() bruker på ditt OS.

Eksempel på config fil:
```yaml
Courses:
  - code: "IN1000"
    title: "Introduksjon til objektorientert programmering"
    semester: h25

  - code: "IN1010"
    title: "Objektorientert programmering"
    semester: h25
```

### 3. ififeed er klar til å brukes!

Root kommando:
```bash
ififeed
```

## Kommandoer

| Kommando | Argumenter | Beskrivelse |
|---|---|---|
| `ififeed` | — | Åpner TUI med beskjeder fra alle emner du har lagt til|
| `ififeed add` | `<emnekode> <semester> <tittel>` | Legger til et emne i config |
| `ififeed remove` | `<emnekode> <semester>` | Fjerner et emne fra config |
| `ififeed list` | — | Viser alle emner i config |

### Semesterformat

Semester skrives på formen `v25` (vår 2025) eller `h25` (høst 2025) etc.

### Eksempler på bruk

```bash
ififeed add IN2140 v25 "Introduksjon til operativsystemer og datakommunikasjon"
>ififeed: la til IN2140 v25
ififeed remove IN1000 h24
>ififeed: fjernet IN2140 v25
```

## Videre utvikling
- Legg til caching av beskjeder
- Legg til søk og sortering av beskjeder basert på emne
