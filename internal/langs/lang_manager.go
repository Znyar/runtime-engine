package langs

import (
	"log/slog"
	"os"
	"os/exec"
	"runtime-engine/internal/runners"
)

type LangManager struct {
	log *slog.Logger
}

func New(log *slog.Logger) *LangManager {
	return &LangManager{
		log: log.With(slog.String("component", "langmanager")),
	}
}

func (m *LangManager) Install(lang runners.Language, version string) error {
	scriptPath := "/app/scripts/packages/" + string(lang) + "/" + version + "/install.sh"

	m.log.Info("installing language",
		slog.String("lang", string(lang)),
		slog.String("version", version),
	)

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *LangManager) Pack(lang runners.Language, version string) error {
	scriptPath := "/app/scripts/packages/" + string(lang) + "/" + version + "/pack.sh"

	m.log.Info("installing package",
		slog.String("lang", string(lang)),
		slog.String("version", version),
	)

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *LangManager) IsInstalled(lang string) bool {
	_, err := exec.LookPath(lang)
	return err == nil
}
