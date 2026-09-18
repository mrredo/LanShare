# LanShare

## Ievads JĀUZRAKSTA VAIRĀK UN VAIRĀK YAPPPP

### Aprakstītā situācija pirms produkta izveides

Lai pārsūtītu failus starp telefonu, datoru un citām ierīcēm, parasti tiek izmantots vads, mākoņpakalpojums vai ziņapmaiņas lietotne. Vada izmantošana ne vienmēr ir ērta, savukārt mākoņpakalpojumiem un ziņapmaiņas lietotnēm nepieciešams interneta savienojums un faili tiek pārsūtīti caur ārējiem serveriem.

Tas nav ērti, kad ir jāpārsūta lieli faili, kad interneta savienojums nav pieejams vai ir ļoti lēns.

### Pamatota produkta nepieciešamība

Failu pārsūtīšana starp vairākām ierīcēm ir bieža darbība. LanShare nodrošina vienkāršāku risinājumu, ļaujot pārsūtīt failus tieši lokālajā tīklā bez ārējiem failu glabāšanas pakalpojumiem.

### Pamatota produkta aktualitāte

Cilvēki ikdienā izmanto vairākas elektroniskās ierīces, tāpēc nepieciešamība pārsūtīt informāciju starp tām ir aktuāla.

LanShare ļauj failus pārsūtīt lokālajā tīklā, neizmantojot vadus vai ārējus failu pārsūtīšanas pakalpojumus.

## Uzdevuma formulējums

### Produkta nosaukums un veids

**LanShare**

Tīmekļa lietotne un serveris.

### Produkta izstrādes mērķis un pamatojums

LanShare mērķis ir nodrošināt vienkāršu failu, teksta un saišu pārsūtīšanu starp ierīcēm, izmantojot lokālo tīklu.

### Produkta mērķauditorija

LanShare mērķauditorija ir skolēni, skolotāji, darbinieki un citi lietotāji, kuri izmanto vairākas ierīces un vēlas pārsūtīt informāciju starp tām.

### Produktā realizēšanai nepieciešamie elementi

**Tehnoloģijas:**

* Go
* React
* TypeScript
* HTML
* CSS
* SQLite

**Sistēmas sastāvdaļas:**

* Go serveris ar iebūvētu mDNS servisu;
* React tīmekļa lietotne;
* SQLite datubāze;
* failu glabāšanas mape serverī.

### Produkta darbības nodrošināšana

Lai produkts darbotos, nepieciešams:

* lokālais tīkls ar vismaz 100 Mbit/s savienojuma ātrumu;
* ierīce ar Windows vai Linux operētājsistēmu;
* vismaz 2 GB RAM;
* vismaz 1 GB brīvas vietas programmas un datubāzes darbībai;
* pietiekama papildu vieta failu glabāšanai;
* mūsdienīgs interneta pārlūks, piemēram, Google Chrome, Mozilla Firefox, Microsoft Edge vai Safari.

SQLite datubāzei nav nepieciešams atsevišķs datubāzes serveris. Failu glabāšanai nepieciešamā diska vieta ir atkarīga no administratora konfigurētā glabāšanas limita un augšupielādēto failu apjoma.


### Produkta pieejamības nodrošināšanas iespējas

Kad tiek ieslēgts LanShare serveris, lokālā tīkla lietotāji var piekļūt tīmekļa lietotnei, izmantojot `lanshare.local`.

Servera atrašanu nodrošina Go aplikācijā iebūvētais mDNS serviss.

---

# Prasību specifikācija

## Sistēmas funkcionālās prasības

### Sistēmas lomas

LanShare sistēmā ir divas lomas:

* **Lietotājs** — var augšupielādēt, lejupielādēt, koplietot un dzēst savus failus.
* **Administrators** — papildus lietotāja funkcijām var pārvaldīt visus failus, servera konfigurāciju un augšupielādes iespējas.

Lietotāja faili tiek sasaistīti ar pārlūkprogrammas cookie. Administratora piekļuve tiek aizsargāta ar administratora paroli.

### 1. Pieslēgšanās sistēmai

Lietotājs savā ierīcē atver `lanshare.local`.

Pirmajā pieslēgšanās reizē serveris pārbauda, vai pārlūkprogrammā ir LanShare cookie.

**Ievades dati:**

* **Ierīces cookie** — unikāls identifikators konkrētai pārlūkprogrammai.

**Apstrāde:**

Ja cookie neeksistē, serveris ģenerē jaunu unikālu identifikatoru un saglabā to pārlūkprogrammā.

Ja cookie eksistē, tas tiek izmantots lietotāja failu noteikšanai.

Cookie netiek izmantots administratora autentifikācijai.

**Rezultāts:**

Lietotājs var izmantot LanShare un redzēt savus failus.

### 2. Administratora sākotnējā konfigurācija

Pirmo reizi startējot LanShare serveri, tiek pārbaudīts, vai administratora parole jau ir iestatīta.

Ja parole nav iestatīta, serveris pieprasa to ievadīt.

**Ievades dati:**

* **Administratora parole** — parole administrācijas sadaļas aizsardzībai.

**Apstrāde:**

Serveris pārbauda ievadīto paroli un saglabā tās hešu SQLite datubāzē. Pati parole datubāzē netiek saglabāta.


**Rezultāts:**

Administrators var autentificēties administrācijas sadaļā `lanshare.local/admin`.

### 3. Administratora autentifikācija

Administrators atver `lanshare.local/admin` un ievada administratora paroli.

**Ievades dati:**

* **Administratora parole**

**Apstrāde:**

Serveris pārbauda ievadīto paroli.

Ja parole ir pareiza, tiek izveidota administratora sesija.

Ja parole nav pareiza, piekļuve administrācijas sadaļai tiek atteikta.

**Rezultāts:**

Pēc veiksmīgas autentifikācijas serveris izveido administratora sesiju un saglabā tās identifikatoru atsevišķā administratora cookie. Pēc tam administrators var izmantot administrācijas funkcijas.


### 4. Failu augšupielāde

Lietotājs izvēlas failu un pēc izvēles pievieno tekstu, URL un derīguma termiņu.

**Ievades dati:**

* **Fails**
* **Teksts**
* **URL**
* **Derīguma termiņš**
* **Ierīces cookie**

**Apstrāde:**

Serveris pārbauda, vai failu augšupielāde ir atļauta, vai nav pārsniegts maksimālais faila izmērs un vai ir pietiekami daudz brīvas vietas.

Ja failu glabāšanas mape neeksistē, serveris pārbauda konfigurācijā norādīto mapes atrašanās vietu. Ja atrašanās vieta nav norādīta, serveris izveido noklusējuma `/uploads` direktoriju.
Fails tiek saglabāts serverī, bet tā metadati tiek saglabāti SQLite datubāzē.

Failam tiek piešķirts unikāls identifikators un tā īpašnieka cookie identifikators.

**Rezultāts:**

Lietotājam tiek atgriezta faila koplietošanas saite.

```json
{
  "success": true,
  "message": "Fails tika veiksmīgi augšupielādēts",
  "download_url": "http://lanshare.local/uploads/H8e1Ja2"
}
```

Ja augšupielāde neizdodas:

```json
{
  "success": false,
  "message": "Kaut kas nogāja greizi."
}
```

### 5. Failu lejupielāde

Lietotājs atver faila koplietošanas saiti.

**Ievades dati:**

* **Faila identifikators**

**Apstrāde:**

Serveris atrod faila informāciju SQLite datubāzē un pārbauda, vai tā derīguma termiņš nav beidzies.
Ja derīguma termiņš ir beidzies, fails tiek dzēsts no servera un datubāzes.

Ja fails ir pieejams, tas tiek nosūtīts lietotājam.

**Rezultāts:**

Lietotājs saņem failu.

### 6. Failu dzēšana

Lietotājs sadaļā `Mani faili` izvēlas savu failu un apstiprina tā dzēšanu.

**Ievades dati:**

* **Faila identifikators**
* **Ierīces cookie**

**Apstrāde:**

Serveris pārbauda, vai faila īpašnieka cookie sakrīt ar pašreizējās pārlūkprogrammas cookie.

Administratoram ir tiesības dzēst jebkuru failu neatkarīgi no tā īpašnieka.
Ja tie sakrīt, fails tiek dzēsts no servera un tā ieraksts no SQLite datubāzes.

**Rezultāts:**

Fails vairs nav pieejams.

### 7. Failu saraksta apskate

Lietotājs atver `lanshare.local`.

Sadaļā `Mani faili` tiek parādīti faili, kuru īpašnieka cookie sakrīt ar pašreizējās pārlūkprogrammas cookie.

Failu sarakstā tiek parādīts:

* faila nosaukums;
* izmērs;
* augšupielādes laiks;
* derīguma termiņš;
* koplietošanas saite.

### 8. Failu koplietošana

Pēc faila augšupielādes lietotājs saņem unikālu saiti.

Saiti var nosūtīt citam lokālā tīkla lietotājam.

Saņēmējam nav nepieciešams īpašnieka cookie vai konts.

### 9. Failu automātiska dzēšana

Serveris periodiski pārbauda failu derīguma termiņus.

Ja termiņš ir beidzies, fails tiek dzēsts no servera un tā ieraksts no SQLite datubāzes.

### 10. Servera atrašana lokālajā tīklā

Go serveris izmanto mDNS, lai LanShare būtu pieejams ar adresi `lanshare.local`.

Lietotājam nav nepieciešams zināt servera IP adresi.

### 11. Administrācijas panelis

Administrators var piekļūt `lanshare.local/admin`.

Administrācijas panelī administrators var:

* apskatīt izmantoto un brīvo vietu;
* apskatīt failu skaitu;
* apskatīt failu sarakstu;
* dzēst jebkuru failu;
* dzēst visus failus;
* mainīt servera konfigurāciju;
* iestatīt maksimālo pieejamo vietu failiem;
* iestatīt maksimālo faila izmēru;
* apturēt vai atļaut jaunu failu augšupielādi;
* pārvaldīt failu derīguma termiņa iestatījumus;
* apskatīt servera darbības informāciju.

Administrācijas panelis ir pieejams tikai autentificētam administratoram.

### 12. Servera konfigurācija

Servera konfigurāciju var pārvaldīt administrācijas panelī.

Konfigurācijā var tikt saglabāti:

* failu glabāšanas mapes atrašanās vieta;
* failiem pieejamās vietas limits;
* maksimālais faila izmērs;
* vai ir atļauta failu augšupielāde;
* noklusētais faila derīguma termiņš;
* citi servera darbības iestatījumi.

Konfigurācijas dati tiek glabāti SQLite datubāzē.

Mainot konfigurāciju, jaunie iestatījumi tiek izmantoti bez nepieciešamības restartēt serveri.

---

## Sistēmas galvenā funkcionalitāte

LanShare galvenā funkcionalitāte ir failu pārsūtīšana lokālajā tīklā.

Lietotājs atver `lanshare.local`, augšupielādē failu un saņem unikālu koplietošanas saiti. Savus failus lietotājs var apskatīt sadaļā `Mani faili`, lejupielādēt un dzēst.

Cookie identificē lietotāju un tā piederošos failus.

## Sistēmas papildfunkcionalitāte

Papildfunkcionalitāte:

* mDNS servera atrašanai;
* failu derīguma termiņi;
* failu koplietošana ar unikālām saitēm;
* teksta un URL pievienošana failiem;
* administratora panelis;
* servera konfigurācija;
* izmantotās vietas kontrole;
* failu augšupielādes apturēšana;
* administratora iespēja dzēst jebkuru failu;
* automātiska failu dzēšana pēc derīguma termiņa.

---

# Sistēmas nefunkcionālās prasības

## Darbības vides prasības

**Serverim nepieciešams:**

* dators ar Windows vai Linux;
* lokālais tīkls;
* LanShare servera izpildāmais fails;
* pietiekama vieta failu glabāšanai.

SQLite darbībai nav nepieciešams atsevišķs datubāzes serveris.

**Klientam nepieciešams:**

* dators, telefons vai cita ierīce;
* moderna tīmekļa pārlūkprogramma;
* savienojums ar to pašu lokālo tīklu.

## Drošība, datu aizsardzība un uzticamība

Faili tiek pārsūtīti lokālajā tīklā un netiek nosūtīti uz ārējiem failu glabāšanas pakalpojumiem.

Failiem var piekļūt ierīces, kuras var sasniegt LanShare serveri lokālajā tīklā.

Faila īpašnieka noteikšanai tiek izmantots cookie. IP un MAC adrese netiek izmantota kā lietotāja identifikators.

Administrācijas panelis ir aizsargāts ar administratora paroli.

Administratoram ir tiesības pārvaldīt visus failus un servera konfigurāciju.

Sistēmai jānovērš situācija, kurā viena lietotāja cookie ļauj dzēst cita lietotāja failus.

## Saskarne un dizains

LanShare ir tīmekļa lietotne latviešu valodā.

Saskarnei jābūt vienkāršai un responsīvai, lai to varētu izmantot datorā, telefonā un planšetē.

Galvenās sadaļas:

* `Mani faili`;
* faila augšupielāde;
* koplietošana;
* administrācijas panelis.

Administrācijas sadaļa ir pieejama tikai administratoram.

## Veiktspēja

Failu pārsūtīšanas ātrumu nosaka lokālā tīkla un servera aparatūras veiktspēja.

Sistēma nedrīkst būtiski ierobežot pieejamo tīkla ātrumu.

Parastiem API pieprasījumiem atbildes laiks nedrīkst pārsniegt **2 sekundes**. Lielu failu augšupielādei un lejupielādei šis ierobežojums neattiecas.

Sistēmai jāapstrādā vairāki vienlaicīgi pieprasījumi.

## Galvenās nefunkcionālās prasības

* Darbība lokālajā tīklā.
* Windows un Linux atbalsts serverim.
* Tīmekļa pārlūkprogrammas izmantošana klientā.
* Responsīva saskarne.
* Administratora autentifikācija.
* Failu glabāšanas vietas kontrole.
* Failu izmēra un augšupielādes ierobežojumu konfigurēšana.


## Papildu nefunkcionālās prasības

* Lietotājam nav nepieciešams instalēt atsevišķu programmu.
* Serverim jābūt vienkārši konfigurējamam.
* Administrācijas panelim jābūt pieejamam - `lanshare.local/admin`.
* Serverim jāspēj darboties bez interneta savienojuma.
* Sistēmai jānodrošina kļūdu apstrāde un datu konsekvence.

## Uzdevuma risināšanas līdzekļu apraksts un izvēles pamatojums

### Iespējamo risinājuma līdzekļu un valodu apraksts

[Norādi trīs piemērotākās programmēšanas valodu alternatīvas un īsi raksturo katru no tām.]
Programmēšanas valodu alternatīvas ir:
1. PHP:
   * Plašs bibliotēku un ietvaru klāsts.
2. Java:
   * Labi piemērota vairāku vienlaicīgu pieprasījumu apstrādei.
3. Javascript:
   * Var izmantot gan klienta, gan servera puses izstrādei.
   * Samazināta mentālā slodze, jo nav jāpārslēdzas starp redaktoriem, bet visu raksti vienā vietā un vienā valodā.

Tehnoloģiju alternatīvas ir:
1. Laravel ietvars:
    * Nodrošina nepieciešamās funkcijas darbam ar datubāzi, autentifikāciju un tīmekļa pieprasījumiem. 
    * Ļauj ātrāk izstrādāt tīmekļa lietotnes.
2. Flutter:
   * Var izmantot mobilo un WEB lietotņu izstrādei.
   * Nākotnē LanShare var paplašināt ar mobilo lietotni.
3. MySQL:
   * Piemērots lielākām sistēmām ar daudziem lietotājiem.
   * Nodrošina labu veiktspēju un vienlaicīgu pieprasījumu apstrādi.

### Izvēlēto risinājuma līdzekļu un valodu apraksts

#### Izvēlētās programmēšanas valodas

1. **Go programmēšanas valoda**

   * Pluss: Ātra, kompilēta valoda ar labu veiktspēju.
   * Pluss: Ērti piemērota serveru un tīkla lietotņu izstrādei.
   * Pluss: Atbalsta vienlaicīgu uzdevumu izpildi ar gorutīnām. Gorutīnas sākotnēji izmanto tikai dažus KB atmiņas, tāpēc iespējams efektīvi apstrādāt lielu skaitu vienlaicīgu uzdevumu.
   * Mīnuss: Salīdzinot ar dažām citām valodām, ir mazāk iebūvētas funkcionalitātes.
   * Izvēlēta, jo LanShare serverim nepieciešama laba veiktspēja, tīkla pieprasījumu apstrāde un failu pārsūtīšana.
   * Atšķirībā no Java, Go programmas var kompilēt vienā izpildāmā failā, tāpēc servera uzstādīšana ir vienkāršāka.
   * Atšķirībā no JavaScript, Go ir kompilēta valoda un tāpēc ir daudz ātrāka.

2. **JavaScript, HTML un CSS**

   * Pluss: Standarta tīmekļa tehnoloģijas, kuras atbalsta praktiski visas mūsdienu pārlūkprogrammas.
   * Pluss: Ļauj izveidot interaktīvu un responsīvu tīmekļa saskarni.
   * Pluss: Nav nepieciešams instalēt atsevišķu programmu klienta ierīcē.
   * Izvēlētas, jo LanShare ir tīmekļa lietotne, kurai jādarbojas datoros un mobilajās ierīcēs.
   * Atšķirībā no Flutter un WebAssembly, lietotnes izstrāde ar šīm tehnoloģijām ir daudz ātrāka un vienkāršāka.
   * HTML nodrošina lapas struktūru, CSS – vizuālo noformējumu, bet JavaScript – interaktivitāti un saziņu ar serveri.

#### Izvēlēto tehnoloģiju pamatojums

1. **React**

   * Izvēlēts tīmekļa saskarnes izstrādei.
   * Ļauj sadalīt saskarni atkārtoti izmantojamās komponentēs.
   * Atvieglo dinamisku datu, piemēram, failu saraksta un augšupielādes statusa, attēlošanu.
   * Salīdzinot ar vienkāršu JavaScript, React nodrošina ērtāku lielākas un interaktīvākas saskarnes izstrādi.

2. **SQLite**

   * Izvēlēta datu glabāšanai, piemēram, failu metadatiem un servera konfigurācijai.
   * Nav nepieciešams atsevišķs datubāzes serveris.
   * Vienkārši uzstādāma un piemērota nelielai lokālā tīkla sistēmai.
   * Atšķirībā no MySQL nav nepieciešams atsevišķi darbināt datubāzes serveri, tāpēc LanShare uzstādīšana ir vienkāršāka.

3. **mDNS**

   * Izvēlēts servera automātiskai atrašanai lokālajā tīklā.
   * Ļauj lietotājiem piekļūt serverim, izmantojot `lanshare.local`, nevis meklējot servera IP adresi.
   * Atšķirībā no manuālas IP adreses ievadīšanas lietotājam nav jāzina servera IP adrese.


---

# Sistēmas modelēšana un projektēšana

## Sistēmas struktūras modelis

[Iekļauj produkta sistēmas struktūras aprakstu.]

[Pievieno sistēmas shēmu, ER-diagrammu, klašu diagrammu vai līdzīgu diagrammu.]

### `files`

| Lauks | Datu tips | Ierobežojumi | Apraksts |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY, NOT NULL | Unikāls faila identifikators. |
| `filename` | TEXT | NOT NULL | Faila nosaukums. |
| `storage_path` | TEXT | NOT NULL, UNIQUE | Faila atrašanās vieta serverī. |
| `size` | INTEGER | NOT NULL, ≥ 0 | Faila izmērs baitos. |
| `owner_cookie` | TEXT | NOT NULL | Pārlūka sīkdatnes identifikators, kas nosaka faila īpašnieku. |
| `text` | TEXT | NULL | Lietotāja pievienotais teksts. |
| `url` | TEXT | NULL | Lietotāja pievienotā saite. |
| `uploaded_at` | DATETIME | NOT NULL | Faila augšupielādes datums un laiks. |
| `expires_at` | DATETIME | NULL | Datums un laiks, kad failam beidzas derīguma termiņš. |

### `settings`

| Lauks | Datu tips | Ierobežojumi | Apraksts |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY | Iestatījumu ieraksta identifikators. |
| `storage_path` | TEXT | NOT NULL, UNIQUE | Ceļš uz mapi, kurā tiek glabāti faili. |
| `storage_limit` | INTEGER | NOT NULL, > 0 | Maksimālais glabātuves izmērs baitos. |
| `max_file_size` | INTEGER | NOT NULL, > 0 | Maksimālais viena faila izmērs baitos. |
| `uploads_enabled` | BOOLEAN | NOT NULL | Norāda, vai lietotājiem ir atļauts augšupielādēt failus. |
| `default_expiry` | INTEGER | NULL, > 0 | Noklusējuma faila derīguma ilgums sekundēs. |
| `admin_password_hash` | TEXT | NOT NULL | Administratora paroles jaucējkods. |

### `admin_sessions`

| Lauks | Datu tips | Ierobežojumi | Apraksts |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY, NOT NULL | Administratora sesijas unikāls identifikators. |
| `created_at` | DATETIME | NOT NULL | Sesijas izveides datums un laiks. |
| `expires_at` | DATETIME | NOT NULL | Datums un laiks, kad administratora sesija beidzas. |

**Piezīme:** LanShare nav atsevišķas `users` tabulas, jo parastie lietotāji tiek identificēti ar pārlūka `owner_cookie`.

[Pievieno diagrammu kopumu, kas apraksta produkta galvenās struktūras.]

[Pievieno diagrammas, kas apraksta produkta papildfunkciju struktūras.]

## Funkcionālais un dinamiskais sistēmas modelis

[Iekļauj lietojuma gadījumu diagrammas, datu plūsmu diagrammas, secību diagrammas, komunikāciju diagrammas vai scenārijus.]

[Iekļauj algoritmu shēmas, stāvokļu diagrammas, aktivitāšu diagrammas vai izvēlētās risināšanas metodes aprakstu.]

[Pārliecinies, ka diagrammas atbilst sistēmas funkcionālajām prasībām, un īsi apraksti katru diagrammu.]
