package game

import (
	e "github.com/adsozuan/wipeout-rw-go/engine"
)

const (
	SaveDataMagic = 0x64736f77
)

type Save struct {
	Magic   uint32
	IsDirty bool

	SfxVolume   float32
	MusicVolume float32
	UiScale     byte
	ShowFps     bool
	Fullscreen  bool
	ScreenRes   int
	PostEffect  int

	HasRapierClass   uint32
	HasBonusCircuits uint32

	Buttons [NumGameActions][2]e.Button

	HighscoresName [4]byte
	Highscores     [NumRaceClasses][NumCircuits][NumHighscoreTabs]HighScores
}

func NewSave() Save {
	s := Save{
		Magic:       SaveDataMagic,
		IsDirty:     true,
		SfxVolume:   0.6,
		MusicVolume: 0.5,
		UiScale:     0,
		ShowFps:     false,
		Fullscreen:  false,
		ScreenRes:   0,
		PostEffect:  0,

		HasRapierClass:   0,
		HasBonusCircuits: 0,

		Buttons: [NumGameActions][2]e.Button{
			AUp:         {e.InputKeyUp, e.InputGamepadDpadUp},
			ADown:       {e.InputKeyDown, e.InputGamepadDpadDown},
			ALeft:       {e.InputKeyLeft, e.InputGamepadDpadLeft},
			ARight:      {e.InputKeyRight, e.InputGamepadDpadRight},
			ABrakeLeft:  {e.InputKeyC, e.InputGamepadLShoulder},
			ABrakeRight: {e.InputKeyV, e.InputGamepadRShoulder},
			AThrust:     {e.InputKeyX, e.InputGamepadA},
			AFire:       {e.InputKeyZ, e.InputGamepadX},
			AChangeView: {e.InputKeyA, e.InputGamepadY},
		},

		Highscores: [NumRaceClasses][NumCircuits][NumHighscoreTabs]HighScores{
			RaceClassVenom: {
				{
					HighscoreTabRace: HighScores{
						LapRecord: 85.83,
						Entries: [NumHighscores]HighScoreEntry{
							{"WIP", 254.50}, {"EOU", 271.17}, {"TPC", 289.50}, {"NOT", 294.50}, {"PSX", 314.50},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 85.83,
						Entries: [NumHighscores]HighScoreEntry{
							{"MVE", 254.50}, {"ALM", 271.17}, {"POL", 289.50}, {"NIK", 294.50}, {"DAR", 314.50},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 55.33,
						Entries: [NumHighscores]HighScoreEntry{
							{"AJY", 159.33}, {"AJS", 172.67}, {"DLS", 191.00}, {"MAK", 207.67}, {"JED", 219.33},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 55.33,
						Entries: [NumHighscores]HighScoreEntry{
							{"DAR", 159.33}, {"STU", 172.67}, {"MOC", 191.00}, {"DOM", 207.67}, {"NIK", 219.33},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 57.5,
						Entries: [NumHighscores]HighScoreEntry{
							{"JD", 171.00}, {"AJC", 189.33}, {"MSA", 202.67}, {"SD", 219.33}, {"TIM", 232.67},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 57.5,
						Entries: [NumHighscores]HighScoreEntry{
							{"PHO", 171.00}, {"ENI", 189.33}, {"XR", 202.67}, {"ISI", 219.33}, {"NG", 232.67},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 85.17,
						Entries: [NumHighscores]HighScoreEntry{
							{"POL", 251.33}, {"DAR", 263.00}, {"JAS", 283.00}, {"ROB", 294.67}, {"DJR", 314.82},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 85.17,
						Entries: [NumHighscores]HighScoreEntry{
							{"DOM", 251.33}, {"DJR", 263.00}, {"MPI", 283.00}, {"GOC", 294.67}, {"SUE", 314.82},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 80.17,
						Entries: [NumHighscores]HighScoreEntry{
							{"NIK", 236.17}, {"SAL", 253.17}, {"DOM", 262.33}, {"LG", 282.67}, {"LNK", 298.17},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 80.17,
						Entries: [NumHighscores]HighScoreEntry{
							{"NIK", 236.17}, {"ROB", 253.17}, {"AM", 262.33}, {"JAS", 282.67}, {"DAR", 298.17},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 61.67,
						Entries: [NumHighscores]HighScoreEntry{
							{"HAN", 182.33}, {"PER", 196.33}, {"FEC", 214.83}, {"TPI", 228.83}, {"ZZA", 244.33},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 61.67,
						Entries: [NumHighscores]HighScoreEntry{
							{"FC", 182.33}, {"SUE", 196.33}, {"ROB", 214.83}, {"JEN", 228.83}, {"NT", 244.33},
						},
					},
				},
				{
					HighscoreTabRace: HighScores{
						LapRecord: 63.83,
						Entries: [NumHighscores]HighScoreEntry{
							{"CAN", 195.40}, {"WEH", 209.23}, {"AVE", 227.90}, {"ABO", 239.90}, {"NUS", 240.73},
						},
					},
					HighscoreTabTimeTrial: HighScores{
						LapRecord: 63.83,
						Entries: [NumHighscores]HighScoreEntry{
							{"DJR", 195.40}, {"NIK", 209.23}, {"JAS", 227.90}, {"NCW", 239.90}, {"LOU", 240.73},
						},
					},
				},
			},
			RaceClassRapier: {
				{
                HighscoreTabRace: HighScores{
                    LapRecord: 69.50,
                    Entries: [NumHighscores]HighScoreEntry{
                        {"AJY", 200.67}, {"DLS", 213.50}, {"AJS", 228.67}, {"MAK", 247.67}, {"JED", 263.00},
                    },
                },
                HighscoreTabTimeTrial: HighScores{
                    LapRecord: 69.50,
                    Entries: [NumHighscores]HighScoreEntry{
                        {"NCW", 200.67}, {"LEE", 213.50}, {"STU", 228.67}, {"JAS", 247.67}, {"ROB", 263.00},
                    },
                },
            },
            {
                HighscoreTabRace: HighScores{
                    LapRecord: 47.33,
                    Entries: [NumHighscores]HighScoreEntry{
                        {"BOR", 134.58}, {"ING", 147.00}, {"HIS", 162.25}, {"COR", 183.08}, {"ES", 198.25},
                    },
                },
                HighscoreTabTimeTrial: HighScores{
                    LapRecord: 47.33,
                    Entries: [NumHighscores]HighScoreEntry{
                        {"NIK", 134.58}, {"POL", 147.00}, {"DAR", 162.25}, {"STU", 183.08}, {"ROB", 198.25},
                    },
                },
            },
			},
		},
	}

	return s
}

type HighScores struct {
	Entries   [NumHighscores]HighScoreEntry
	LapRecord float32
}

type HighScoreEntry struct {
	Name string
	Time float32
}
