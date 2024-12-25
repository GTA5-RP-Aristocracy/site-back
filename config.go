package siteback

type (
	Config struct {
		ReCaptchaSecret string `env:"RECAPTCHA_SECRET"`
	}
)
