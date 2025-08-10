package langs

import (
	"log/slog"
	"os"
	"os/exec"
)

type LangManager struct {
	log *slog.Logger
}

func New(log *slog.Logger) *LangManager {
	return &LangManager{
		log: log.With(slog.String("component", "langmanager")),
	}
}

func (m *LangManager) Install(lang string, version string) error {
	scriptPath := "/app/scripts/packages/" + lang + "/" + version + "/install.sh"

	m.log.Info("installing language",
		slog.String("lang", lang),
		slog.String("version", version),
	)

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *LangManager) Pack(lang string, version string) error {
	scriptPath := "/app/scripts/packages/" + lang + "/" + version + "/pack.sh"

	m.log.Info("installing package",
		slog.String("lang", lang),
		slog.String("version", version),
	)

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *LangManager) IsInstalled(lang string) bool {
	_, err := exec.LookPath(lang)
	return err == nil
}
