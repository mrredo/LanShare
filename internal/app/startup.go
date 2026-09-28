package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"lanshare/internal/db"
	"lanshare/internal/settings"

	"golang.org/x/term"
)

func SettingsStartup(dbServices ...*db.Service) (*settings.Settings, error) {
	var dbService *db.Service
	if len(dbServices) > 0 && dbServices[0] != nil {
		dbService = dbServices[0]
	} else {
		dbService = db.NewService()
		if err := dbService.Start(); err != nil {
			return nil, fmt.Errorf("datubāzes startēšana: %w", err)
		}
	}
	return SettingsStartupWithIO(dbService, os.Stdin, os.Stdout)
}

func SettingsStartupWithIO(dbService *db.Service, in io.Reader, out io.Writer) (*settings.Settings, error) {
	reader := bufio.NewReader(in)
	repo := settings.NewRepo(dbService.DB())
	service := settings.NewService(repo)

	welcomeBanner(out)

	currentSettings, err := service.GetLatestSettings()
	if err != nil {
		return nil, fmt.Errorf("neizdevās iegūt iestatījumus: %w", err)
	}

	if currentSettings == nil || currentSettings.ID == 0 {
		currentSettings, err = initialSetup(reader, in, out, service)
		if err != nil {
			return nil, err
		}
	} else {
		displaySettingsSummary(out, "Pašreizējā servera konfigurācija", currentSettings)
	}

	for {
		fmt.Fprintln(out, "\nLanShare startēšanas izvēlne:")
		fmt.Fprintln(out, "  [1] Startēt LanShare serveri (Nospiediet Enter)")
		fmt.Fprintln(out, "  [2] Atjaunināt administratora paroli")
		fmt.Fprintln(out, "  [3] Atjaunināt servera iestatījumus")
		choice, err := readLine(reader, out, "Izvēlieties darbību [1]: ")
		if err != nil {
			return nil, err
		}
		choice = strings.TrimSpace(choice)
		if choice == "" || choice == "1" || strings.EqualFold(choice, "start") || strings.EqualFold(choice, "startēt") {
			fmt.Fprintln(out, "\n>>> Tiek startēts LanShare serveris...")
			return currentSettings, nil
		} else if choice == "2" || strings.EqualFold(choice, "parole") || strings.EqualFold(choice, "password") {
			if err := updateAdminPasswordFlow(reader, in, out, service, currentSettings); err != nil {
				fmt.Fprintf(out, "Kļūda atjaunojot paroli: %v\n", err)
			}
		} else if choice == "3" || strings.EqualFold(choice, "iestatījumi") || strings.EqualFold(choice, "settings") || strings.EqualFold(choice, "konfigurācija") {
			updated, err := configureSettingsInteractive(reader, out, currentSettings)
			if err != nil {
				fmt.Fprintf(out, "Kļūda konfigurējot iestatījumus: %v\n", err)
				continue
			}
			if err := service.UpdateSettings(updated); err != nil {
				fmt.Fprintf(out, "Neizdevās saglabāt iestatījumus: %v\n", err)
				continue
			}
			currentSettings = updated
			displaySettingsSummary(out, "Atjauninātā servera konfigurācija", currentSettings)
		} else {
			fmt.Fprintln(out, "Nederīga izvēle. Lūdzu, izvēlieties 1, 2 vai 3.")
		}
	}
}

func welcomeBanner(out io.Writer) {
	fmt.Fprintln(out, "==================================================")
	fmt.Fprintln(out, "              Laipni lūdzam LanShare!             ")
	fmt.Fprintln(out, "  Ātra un droša failu dalīšanās lokālajā tīklā ")
	fmt.Fprintln(out, "==================================================")
}

func displaySettingsSummary(out io.Writer, title string, s *settings.Settings) {
	fmt.Fprintln(out, "\n==================================================")
	fmt.Fprintf(out, "  %s\n", title)
	fmt.Fprintln(out, "==================================================")
	fmt.Fprintf(out, "  Failu mape:         %s\n", s.StoragePath)
	fmt.Fprintf(out, "  Failu mapes izmērs:       %s\n", formatBytes(s.StorageLimit))
	fmt.Fprintf(out, "  Maks. faila izmērs:      %s\n", formatBytes(s.MaxFileSize))
	fmt.Fprintf(out, "  Vai augšupielādes ir atļautas:  %s\n", formatBool(s.UploadsEnabled))
	fmt.Fprintf(out, "  Noklusējuma derīguma laiks:    %s\n", formatDuration(s.DefaultExpiry))
	fmt.Fprintln(out, "==================================================")
}

func initialSetup(reader *bufio.Reader, in io.Reader, out io.Writer, service *settings.Service) (*settings.Settings, error) {
	fmt.Fprintln(out, "\n[Sākotnējā iestatīšana]")
	fmt.Fprintln(out, "Esoša konfigurācija netika atrasta.")
	fmt.Fprintln(out, "Lūdzu, iestatiet administratora paroli, lai aizsargātu serveri.")

	password, err := promptNewPassword(reader, in, out, "Ievadiet administratora paroli: ")
	if err != nil {
		return nil, err
	}

	passHash, err := settings.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("neizdevās šifrēt paroli: %w", err)
	}

	defaultSettings := settings.DefaultSettings(passHash)

	fmt.Fprintln(out, "\nKonfigurēt servera iestatījumus:")
	fmt.Fprintln(out, "  [1] Izmantot noklusējuma iestatījumus (Ieteicams)")
	fmt.Fprintln(out, "  [2] Pielāgoti iestatījumi (konfigurēt katru opciju)")

	var targetSettings *settings.Settings
	for {
		choice, err := readLine(reader, out, "Izvēlieties darbību [1]: ")
		if err != nil {
			return nil, err
		}
		choice = strings.TrimSpace(choice)
		if choice == "" || choice == "1" || strings.EqualFold(choice, "noklusējums") || strings.EqualFold(choice, "default") {
			targetSettings = defaultSettings
			break
		} else if choice == "2" || strings.EqualFold(choice, "pielāgoti") || strings.EqualFold(choice, "custom") {
			custom, err := configureSettingsInteractive(reader, out, defaultSettings)
			if err != nil {
				return nil, err
			}
			targetSettings = custom
			break
		} else {
			fmt.Fprintln(out, "Nederīga izvēle. Lūdzu, izvēlieties 1 vai 2.")
		}
	}

	targetSettings.AdminPasswordHash = passHash
	if err := service.CreateSettings(targetSettings); err != nil {
		return nil, fmt.Errorf("neizdevās saglabāt sākotnējos iestatījumus: %w", err)
	}

	displaySettingsSummary(out, "Sākotnējie iestatījumi veiksmīgi saglabāti", targetSettings)
	return targetSettings, nil
}

func configureSettingsInteractive(reader *bufio.Reader, out io.Writer, current *settings.Settings) (*settings.Settings, error) {
	fmt.Fprintln(out, "\n--- Iestatījumu konfigurēšana (Nospiediet Enter vai 'skip', lai saglabātu vērtību) ---")

	// 1. Glabāšanas ceļš
	storagePath := current.StoragePath
	for {
		prompt := fmt.Sprintf("1. Glabāšanas ceļš [%s]: ", storagePath)
		input, err := readLine(reader, out, prompt)
		if err != nil {
			return nil, err
		}
		if isSkip(input) {
			break
		}
		storagePath = input
		break
	}

	// 2. Glabāšanas limits
	storageLimit := current.StorageLimit
	for {
		prompt := fmt.Sprintf("2. Glabāšanas limits [%s]: ", formatBytes(storageLimit))
		input, err := readLine(reader, out, prompt)
		if err != nil {
			return nil, err
		}
		parsed, err := parseBytes(input, storageLimit)
		if err != nil {
			fmt.Fprintf(out, "   Kļūda: %v. Lūdzu, mēģiniet vēlreiz.\n", err)
			continue
		}
		storageLimit = parsed
		break
	}

	// 3. Maks. faila izmērs
	maxFileSize := current.MaxFileSize
	for {
		prompt := fmt.Sprintf("3. Maks. faila izmērs [%s]: ", formatBytes(maxFileSize))
		input, err := readLine(reader, out, prompt)
		if err != nil {
			return nil, err
		}
		parsed, err := parseBytes(input, maxFileSize)
		if err != nil {
			fmt.Fprintf(out, "   Kļūda: %v. Lūdzu, mēģiniet vēlreiz.\n", err)
			continue
		}
		if parsed > storageLimit {
			fmt.Fprintf(out, "   Brīdinājums: Maks. faila izmērs (%s) pārsniedz kopējo glabāšanas limitu (%s).\n", formatBytes(parsed), formatBytes(storageLimit))
		}
		maxFileSize = parsed
		break
	}

	// 4. Augšupielādes atļautas
	uploadsEnabled := current.UploadsEnabled
	for {
		prompt := fmt.Sprintf("4. Augšupielādes atļautas [%s] (J/n): ", formatBool(uploadsEnabled))
		input, err := readLine(reader, out, prompt)
		if err != nil {
			return nil, err
		}
		parsed, err := parseBool(input, uploadsEnabled)
		if err != nil {
			fmt.Fprintf(out, "   Kļūda: %v. Lūdzu, mēģiniet vēlreiz.\n", err)
			continue
		}
		uploadsEnabled = parsed
		break
	}

	// 5. Noklusējuma faila derīgums
	defaultExpiry := current.DefaultExpiry
	for {
		prompt := fmt.Sprintf("5. Noklusējuma faila derīgums [%s]: ", formatDuration(defaultExpiry))
		input, err := readLine(reader, out, prompt)
		if err != nil {
			return nil, err
		}
		parsed, err := parseDurationSeconds(input, defaultExpiry)
		if err != nil {
			fmt.Fprintf(out, "   Kļūda: %v. Lūdzu, mēģiniet vēlreiz.\n", err)
			continue
		}
		defaultExpiry = parsed
		break
	}

	result := *current
	result.StoragePath = storagePath
	result.StorageLimit = storageLimit
	result.MaxFileSize = maxFileSize
	result.UploadsEnabled = uploadsEnabled
	result.DefaultExpiry = defaultExpiry
	return &result, nil
}

func updateAdminPasswordFlow(reader *bufio.Reader, in io.Reader, out io.Writer, service *settings.Service, current *settings.Settings) error {
	fmt.Fprintln(out, "\n--- Administratora paroles atjaunināšana ---")
	newPass, err := promptNewPassword(reader, in, out, "Ievadiet jauno administratora paroli: ")
	if err != nil {
		return err
	}

	newHash, err := settings.HashPassword(newPass)
	if err != nil {
		return fmt.Errorf("paroles šifrēšana: %w", err)
	}

	if err := service.UpdateAdminPasswordHash(fmt.Sprint(current.ID), newHash); err != nil {
		return fmt.Errorf("administratora paroles atjaunināšana datubāzē: %w", err)
	}

	current.AdminPasswordHash = newHash
	fmt.Fprintln(out, "Administratora parole veiksmīgi atjaunināta!")
	return nil
}

func promptNewPassword(reader *bufio.Reader, in io.Reader, out io.Writer, promptLabel string) (string, error) {
	for {
		pass, err := readPassword(reader, in, out, promptLabel)
		if err != nil {
			return "", err
		}
		if pass == "" {
			fmt.Fprintln(out, "Parole nevar būt tukša. Lūdzu, mēģiniet vēlreiz.")
			continue
		}

		confirm, err := readPassword(reader, in, out, "Apstipriniet administratora paroli: ")
		if err != nil {
			return "", err
		}
		if pass != confirm {
			fmt.Fprintln(out, "Paroles nesakrīt. Lūdzu, mēģiniet vēlreiz.")
			continue
		}

		return pass, nil
	}
}

func readPassword(reader *bufio.Reader, in io.Reader, out io.Writer, prompt string) (string, error) {
	fmt.Fprint(out, prompt)
	if file, ok := in.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		bytes, err := term.ReadPassword(int(file.Fd()))
		fmt.Fprintln(out) // Pārejam jaunā rindā pēc paroles ievades
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(bytes)), nil
	}

	// Rezerves variants, ja ievade nav terminālis (piem., testi vai caurules)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func readLine(reader *bufio.Reader, out io.Writer, prompt string) (string, error) {
	fmt.Fprint(out, prompt)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func isSkip(input string) bool {
	s := strings.TrimSpace(strings.ToLower(input))
	return s == "" || s == "skip" || s == "izlaist" || s == "atstāt"
}

func parseBytes(input string, defaultVal int64) (int64, error) {
	if isSkip(input) {
		return defaultVal, nil
	}
	sUpper := strings.ToUpper(strings.TrimSpace(input))
	var multiplier int64 = 1
	trimmed := sUpper
	if strings.HasSuffix(sUpper, "TB") {
		multiplier = 1024 * 1024 * 1024 * 1024
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-2])
	} else if strings.HasSuffix(sUpper, "GB") {
		multiplier = 1024 * 1024 * 1024
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-2])
	} else if strings.HasSuffix(sUpper, "MB") {
		multiplier = 1024 * 1024
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-2])
	} else if strings.HasSuffix(sUpper, "KB") {
		multiplier = 1024
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-2])
	} else if strings.HasSuffix(sUpper, "BAITI") {
		multiplier = 1
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-5])
	} else if strings.HasSuffix(sUpper, "B") {
		multiplier = 1
		trimmed = strings.TrimSpace(sUpper[:len(sUpper)-1])
	}

	val, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || val <= 0 {
		return 0, fmt.Errorf("jābūt pozitīvam skaitlim")
	}
	return val * multiplier, nil
}

func formatBytes(b int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	switch {
	case b >= TB && b%TB == 0:
		return fmt.Sprintf("%d TB", b/TB)
	case b >= GB && b%GB == 0:
		return fmt.Sprintf("%d GB", b/GB)
	case b >= MB && b%MB == 0:
		return fmt.Sprintf("%d MB", b/MB)
	case b >= KB && b%KB == 0:
		return fmt.Sprintf("%d KB", b/KB)
	default:
		return fmt.Sprintf("%d baiti", b)
	}
}

func parseDurationSeconds(input string, defaultVal *int64) (*int64, error) {
	if isSkip(input) {
		return defaultVal, nil
	}
	sLower := strings.ToLower(strings.TrimSpace(input))
	if sLower == "none" || sLower == "0" || sLower == "never" || sLower == "nekad" || sLower == "nav" || sLower == "neierobežots" || sLower == "null" || sLower == "atspējots" {
		return nil, nil
	}
	var multiplier int64 = 1
	trimmed := sLower
	if strings.HasSuffix(sLower, "dienas") || strings.HasSuffix(sLower, "diena") || strings.HasSuffix(sLower, "days") || strings.HasSuffix(sLower, "day") || strings.HasSuffix(sLower, "d") {
		multiplier = 24 * 3600
		idx := strings.IndexAny(sLower, "d")
		trimmed = strings.TrimSpace(sLower[:idx])
	} else if strings.HasSuffix(sLower, "stundas") || strings.HasSuffix(sLower, "stunda") || strings.HasSuffix(sLower, "hours") || strings.HasSuffix(sLower, "hour") || strings.HasSuffix(sLower, "h") {
		multiplier = 3600
		idx := strings.IndexAny(sLower, "h")
		if idx == -1 {
			idx = strings.Index(sLower, "st")
		}
		trimmed = strings.TrimSpace(sLower[:idx])
	} else if strings.HasSuffix(sLower, "minūtes") || strings.HasSuffix(sLower, "minūte") || strings.HasSuffix(sLower, "minutes") || strings.HasSuffix(sLower, "min") || strings.HasSuffix(sLower, "m") {
		multiplier = 60
		idx := strings.IndexAny(sLower, "m")
		trimmed = strings.TrimSpace(sLower[:idx])
	} else if strings.HasSuffix(sLower, "sekundes") || strings.HasSuffix(sLower, "sekunde") || strings.HasSuffix(sLower, "seconds") || strings.HasSuffix(sLower, "sec") || strings.HasSuffix(sLower, "s") {
		multiplier = 1
		idx := strings.IndexAny(sLower, "s")
		trimmed = strings.TrimSpace(sLower[:idx])
	}

	val, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || val <= 0 {
		return nil, fmt.Errorf("jābūt pozitīvam laikam (piem., 24h, 7d, 86400s) vai '0'/'nekad'")
	}
	res := val * multiplier
	return &res, nil
}

func formatDuration(sec *int64) string {
	if sec == nil {
		return "Nav (neierobežots)"
	}
	s := *sec
	if s%86400 == 0 {
		days := s / 86400
		if days == 1 {
			return fmt.Sprintf("1 diena (%d sekundes)", s)
		}
		return fmt.Sprintf("%d dienas (%d sekundes)", days, s)
	}
	if s%3600 == 0 {
		hours := s / 3600
		if hours == 1 {
			return fmt.Sprintf("1 stunda (%d sekundes)", s)
		}
		return fmt.Sprintf("%d stundas (%d sekundes)", hours, s)
	}
	if s%60 == 0 {
		mins := s / 60
		if mins == 1 {
			return fmt.Sprintf("1 minūte (%d sekundes)", s)
		}
		return fmt.Sprintf("%d minūtes (%d sekundes)", mins, s)
	}
	return fmt.Sprintf("%d sekundes", s)
}

func formatBool(b bool) string {
	if b {
		return "Jā"
	}
	return "Nē"
}

func parseBool(input string, defaultVal bool) (bool, error) {
	if isSkip(input) {
		return defaultVal, nil
	}
	s := strings.TrimSpace(strings.ToLower(input))
	switch s {
	case "j", "ja", "jā", "y", "yes", "true", "1", "t", "ieslēgts", "iespējots":
		return true, nil
	case "n", "ne", "nē", "no", "false", "0", "f", "izslēgts", "atspējots":
		return false, nil
	default:
		return false, fmt.Errorf("nederīga izvēle: ievadiet 'j' (jā) vai 'n' (nē)")
	}
}
