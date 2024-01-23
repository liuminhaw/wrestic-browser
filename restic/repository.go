package restic

const (
	resticCmd   string = "restic"
	passwordEnv string = "RESTIC_PASSWORD"
)

type Repository interface {
	Connect() error
}
