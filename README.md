# go-ififeed

![Made with VHS](https://vhs.charm.sh/vhs-5cSeP5vpvW9eqbssNVZlcY.gif)

**ififeed** er et verktøy for deg som er IFI-student og er lei av å
måtte klikke deg inn på Mine studier eller semestersiden for hvert enkelt
emne bare for å sjekke nye beskjeder.

Med **ififeed** kan du hente alle beskjeder fra alle emner
du tar rett i terminalen! Ingen flere avbrekk fra den produktive progge-økta!

## Kom i gang

### 1. Last ned binary-filen

Gå til [Releases](https://github.com/willbehn/go-ifi-feed/releases) og last ned riktig binary-fil:

| Plattform | Fil |
|---|---|
| macOS (Apple Silicon) | `ififeed-darwin-arm64` |
| Linux | `ififeed-linux-amd64` |

Gjør filen kjørbar og flytt den til PATH:

```bash
chmod +x ififeed-darwin-arm64
mv ififeed-darwin-arm64 /usr/local/bin/ififeed
```

### 2. Kjør oppsett

Kjør `ififeed` en gang for å opprette config-filen, som havner her:

```
~/.config/ififeed/courses.yaml
```

Deretter kan du redigere `courses.yaml`filen ved å legge til kursene du tar på formen:

```yaml
Courses:
  - code: "<emnekode>"
    title: "<tittel>"
```

### 3. ififeed er klar til bruk!

```bash
ififeed
```
