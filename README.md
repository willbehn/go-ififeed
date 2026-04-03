[![GitHub Release](https://img.shields.io/github/v/release/willbehn/go-ififeed)](https://github.com/willbehn/go-ififeed/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/willbehn/go-ififeed)](go.mod)
![License](https://img.shields.io/github/license/willbehn/go-ififeed)
# go-ififeed

**ififeed** er et verktøy for deg som er IFI-student og er lei av å
måtte klikke deg inn på Mine studier eller semestersiden for hvert enkelt
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

### 2. Velg emnene du tar

Kjør `ififeed` en gang for å opprette en config-fil, som havner her:

```
~/.config/ififeed/courses.yaml
```

Deretter kan du redigere `courses.yaml` filen ved å legge til kursene du tar på formen:

```yaml
Courses:
  - code: "<emnekode>"
    title: "<tittel>"
    semester: "<v25>"
```

### 3. ififeed er klar til bruk!
Du kan nå kjøre ififeed med kommandoen:

```bash
ififeed
```

## Videre utvikling
- Gjøre det mulig å legge til kurs via terminalen, og ikke måtte endre config filen manuelt
- Legg til caching av beskjeder
- Legg til søk og sortering av beskjeder basert på emne