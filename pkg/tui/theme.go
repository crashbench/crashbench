package tui

import "fmt"

// Bohemian & Handwritten Minimalist TrueColor ANSI Palette
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// 24-bit TrueColor tokens matching web theme
	ColorTerracotta = "\033[38;2;200;90;50m"   // #C85A32
	ColorSage       = "\033[38;2;86;110;80m"   // #566E50
	ColorOchre      = "\033[38;2;199;139;54m"  // #C78B36
	ColorCream      = "\033[38;2;250;247;242m" // #FAF7F2
	ColorMuted      = "\033[38;2;132;119;109m" // #84776D
	ColorStone      = "\033[38;2;168;156;146m" // #A89C92
	ColorEspresso   = "\033[38;2;43;36;32m"    // #2B2420

	// Combined styles
	ColorBoldTerracotta = Bold + ColorTerracotta
	ColorBoldSage       = Bold + ColorSage
	ColorBoldOchre      = Bold + ColorOchre
	ColorBoldCream      = Bold + ColorCream
	ColorItalicSage     = Italic + ColorSage

	// Backgrounds
	BgTerracotta = "\033[48;2;200;90;50m"
	BgSage       = "\033[48;2;86;110;80m"
	BgOchre      = "\033[48;2;199;139;54m"
	BgPaperTint  = "\033[48;2;245;239;230m"
	BgEspresso   = "\033[48;2;27;23;20m"
)

// Helper styling functions
func Terracotta(s string) string { return ColorTerracotta + s + Reset }
func Sage(s string) string       { return ColorSage + s + Reset }
func Ochre(s string) string      { return ColorOchre + s + Reset }
func Cream(s string) string      { return ColorCream + s + Reset }
func Muted(s string) string      { return ColorMuted + s + Reset }
func Stone(s string) string      { return ColorStone + s + Reset }
func Espresso(s string) string   { return ColorEspresso + s + Reset }

func BoldTerracotta(s string) string { return Bold + ColorTerracotta + s + Reset }
func BoldSage(s string) string       { return Bold + ColorSage + s + Reset }
func BoldOchre(s string) string      { return Bold + ColorOchre + s + Reset }
func BoldCream(s string) string      { return Bold + ColorCream + s + Reset }
func HandNote(s string) string       { return Italic + ColorTerracotta + s + Reset }

// Version holds the current release version of CrashBench, resolvable automatically
var Version = "v1.1.0"

// SetVersion allows setting the version dynamically from build info or CLI flags
func SetVersion(v string) {
	if v != "" {
		Version = v
	}
}

// Banner returns the artisan Bohemian linocut stamp banner with dynamic version
func Banner() string {
	return fmt.Sprintf(`
  %s✦ ───────────────────────────────────────────────────────────────────────────── ✦%s
   %s[CB]%s  %sC R A S H B E N C H%s  %s%s%s  %s~ arena 2.0 ~%s
        %sThe AI Agent Chaos & Operational Reliability Benchmark%s
        %sBradley-Terry MLE (λ=0.15) · 300 Bayesian Bootstraps · Zero-Leak DFA%s
  %s✦ ───────────────────────────────────────────────────────────────────────────── ✦%s
`,
		ColorOchre, Reset,
		ColorBoldTerracotta, Reset,
		ColorBoldCream, Reset,
		ColorOchre, Version, Reset,
		ColorItalicSage, Reset,
		ColorStone, Reset,
		ColorMuted, Reset,
		ColorOchre, Reset,
	)
}
