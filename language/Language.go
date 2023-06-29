package language

type Language string

const (
	English    Language = "EN"
	Portuguese Language = "PT"
)

func GetWords(language Language) map[string]bool {

	if language == English {
		return StopWordsEnglish
	}

	if language == Portuguese {
		return StopWordsPortuguese
	}

	panic("Invalid language")

}
