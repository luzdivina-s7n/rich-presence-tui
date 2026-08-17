package tui
import "github.com/charmbracelet/lipgloss"
var (
	colorBg       = lipgloss.Color("#121214")
	colorPanel    = lipgloss.Color("#1b1b1f")
	colorRaised   = lipgloss.Color("#232329")
	colorSelected = lipgloss.Color("#3c3c44")
	colorText     = lipgloss.Color("#f2f2f2")
	colorMuted    = lipgloss.Color("#8f8f97")
	colorDim      = lipgloss.Color("#5f5f66")
	colorBorder   = lipgloss.Color("#3a3a41")
	appStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Foreground(colorText).
			Padding(0, 1)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			Underline(true)
	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)
	infoStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Bold(true)
	successStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)
	warningStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)
	errorStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true).
			Background(colorDim)
	selectedRow = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorSelected).
			Bold(true)
	normalRow = lipgloss.NewStyle().
			Foreground(colorText)
	dimRow = lipgloss.NewStyle().
			Foreground(colorDim)
	badgeActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)
	badgeInactive = lipgloss.NewStyle().
			Foreground(colorDim)
	fieldLabelStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
	fieldValueStyle = lipgloss.NewStyle().
			Foreground(colorText)
	fieldSelectedStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Background(colorSelected).
				Bold(true)
	boolOnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)
	boolOffStyle = lipgloss.NewStyle().
			Foreground(colorDim)
	keyCapStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)
	hintTextStyle = lipgloss.NewStyle().
			Foreground(colorDim)
	overlayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Background(colorPanel).
			Foreground(colorText).
			Padding(1, 2)
	inputPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorMuted)
	inputTextStyle = lipgloss.NewStyle().
			Foreground(colorText)
	inputPlaceholderStyle = lipgloss.NewStyle().
				Foreground(colorDim)
	inputCursorStyle = lipgloss.NewStyle().
				Foreground(colorBg).
				Background(colorMuted)
)
