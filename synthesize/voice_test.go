package synthesize

import (
	"net/http"
	"testing"

	"github.com/mrwormhole/errdiff"
)

var testVoices = []Voice{
	AfrikaansVoice,
	AlbanianVoice,
	AmharicVoice,
	ArabicVoice,
	BengaliVoice,
	BosnianVoice,
	BulgarianVoice,
	CantoneseVoice,
	CatalanVoice,
	ChineseSimplifiedVoice,
	ChineseTraditionalVoice,
	CroatianVoice,
	CzechVoice,
	DanishVoice,
	DutchVoice,
	EnglishVoice,
	EstonianVoice,
	FilipinoVoice,
	FinnishVoice,
	FrenchVoice,
	FrenchCanadianVoice,
	GalicianVoice,
	GermanVoice,
	GreekVoice,
	GujaratiVoice,
	HausaVoice,
	HebrewVoice,
	HindiVoice,
	HungarianVoice,
	IcelandicVoice,
	IndonesianVoice,
	ItalianVoice,
	JapaneseVoice,
	JavaneseVoice,
	KhmerVoice,
	KoreanVoice,
	LatinVoice,
	LatvianVoice,
	LithuanianVoice,
	MalayVoice,
	MalayalamVoice,
	MarathiVoice,
	MyanmarVoice,
	NepaliVoice,
	NorwegianVoice,
	PolishVoice,
	PortugueseBrazilianVoice,
	PortugueseVoice,
	PunjabiVoice,
	RomanianVoice,
	RussianVoice,
	SerbianVoice,
	SinhalaVoice,
	SlovakVoice,
	SpanishVoice,
	SundaneseVoice,
	SwahiliVoice,
	SwedishVoice,
	TamilVoice,
	TeluguVoice,
	ThaiVoice,
	UkrainianVoice,
	UrduVoice,
	VietnameseVoice,
	WelshVoice,
}

func TestRun_AllVoices(t *testing.T) {
	t.Parallel()

	type Test struct {
		name    string
		client  *http.Client
		opt     Opt
		wantErr error
	}
	var tests []Test

	for _, voice := range testVoices {
		opt := Opt{
			Speed: NormalSpeed,
			Voice: voice,
			Text:  "hello",
		}
		tests = append(tests, Test{
			name:    "language voice:" + string(opt.Voice),
			client:  &http.Client{},
			opt:     opt,
			wantErr: nil,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			audio, err := Run(t.Context(), tt.client, tt.opt)
			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Errorf("Run(%v): err diff=\n%s", tt.opt, diff)
			}

			if len(audio) < 1 {
				t.Errorf("Run(%v): audio len(%v) must be greater than 1", tt.opt, len(audio))
			}
		})
	}
}
