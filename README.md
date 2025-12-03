# ififeed

![Made with VHS](https://vhs.charm.sh/vhs-56bf3ZDqfNxs40A9SYrn4H.gif)

**ififeed** er et lite verktøy for deg som er IFI-student og er lei av å
måtte klikke deg inn på *Mine studier* eller semestersiden for hvert enkelt
emne bare for å sjekke nye beskjeder.

Med **ififeed** kan du hente alle beskjeder fra alle emner
du tar, rett i terminalen. Ingen flere avbrekk fra den produktive progge-økta!

## 🚀 Kom i gang

### 1. Klon repoet

``` bash
git clone https://github.com/willbehn/go-ififeed.git
cd go-ififeed
```

### 2. Legg inn ønskede kurs

Konfigurer hvilke emner du tar ved å legge dem inn i konfigurasjonsfilen
(`courses.yaml`) slik:

``` yaml
Courses:
  - code: "IN5060"
    title: "Kvantitativ ytelsesanalyse"
  - code: "IN1010"
    title: "Objektorientert programmering"
```

### 3. Kjør programmet

```bash
go build -o ififeed
./ififeed
```


