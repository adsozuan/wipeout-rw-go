package wipeout

import (
	"log"
	"os"

	"github.com/adsozuan/wipeout-rw-go/engine"
)

const (
	NumAIOpponents      = 7
	NumPilotsPerTeam    = 2
	NumNonBonusCircuits = 6
	NumMusicTracks      = 11
	NumHighscores       = 5

	NumLaps        = 3
	NumLives       = 3
	QualifyingRank = 3
	SaveDataMagic  = 0x64736f77
)

type Action int

const (
	AUp Action = iota
	ADown
	ALeft
	ARight
	ABrakeLeft
	ABrakeRight
	AThrust
	AFire
	AChangeView
	NumGameActions

	AMenuUp
	AMenuDown
	AMenuLeft
	AMenuRight
	AMenuBack
	AMenuSelect
	AMenuStart
	AMenuQuit
)

type GameScene int

const (
	GameSceneIntro GameScene = iota
	GameSceneTitle
	GameSceneMainMenu
	GameSceneHighscores
	GameSceneRace
	GameSceneNone
	NumGameScenes
)

type RaceClassE int

const (
	RaceClassVenom RaceClassE = iota
	RaceClassRapier
	NumRaceClasses
)

type RaceTypeE int

const (
	RaceTypeChampionship RaceTypeE = iota
	RaceTypeSingle
	RaceTypeTimeTrial
	NumRaceTypes
)

type HighscoreTab int

const (
	HighscoreTabTimeTrial HighscoreTab = iota
	HighscoreTabRace
	NumHighscoreTabs
)

type PilotE int

const (
	PilotJohnDekka PilotE = iota
	PilotDanielChang
	PilotArialTetsuo
	PilotAnastasiaCherovoski
	PilotKelSolaar
	PilotArianTetsuo
	PilotSofiaDeLaRente
	PilotPaulJackson
	NumPilots
)

type TeamE int

const (
	TeamAGSystems TeamE = iota
	TeamAuricom
	TeamQirex
	TeamFeisar
	NumTeams
)

type CircuitE int

const (
	CircuitAltimaVII CircuitE = iota
	CircuitKarbonisV
	CircuitTerramax
	CircuitKorodera
	CircuitArridosIV
	CircuitSilverstream
	CircuitFirestar
	NumCircuits
)

// Logger is a package-level logger
var Logger *log.Logger

func init() {
	// Create a logger that writes to standard error with a timestamp
	Logger = log.New(os.Stderr, "wipeout|", log.Ldate|log.Ltime)
}

type GameDefinition struct {
	RaceClasses      [NumRaceClasses]RaceClass
	RaceTypes        [NumRaceTypes]RaceType
	Pilots           [NumPilots]Pilot
	Teams            [NumTeams]Team
	AiSettings       [NumRaceClasses][NumAIOpponents]AiSetting
	Circuits         [NumCircuits]Circuit
	ShipModelToPilot [NumPilots]int
	RaceModeForRank  [NumPilots]int
	MusicTracks      [NumMusicTracks]MusicTrack
	Credits          string
	Congratulations  Congratulations
}

type Congratulations struct {
	Venom             [15]string
	VenomAllCircuits  [19]string
	Rapier            [26]string
	RapierAllCircuits [24]string
}

type RaceClass struct {
	Name string
}

type RaceType struct {
	Name string
}

type Pilot struct {
	Name      string
	Portrait  string
	LogoModel int
	Team      int
}

type AiSetting struct {
	ThrustMax       float32
	ThrustMagnitude float32
	FightBack       bool
}

type TeamAttributes struct {
	Mass        float32
	ThrustMax   float32
	Resistance  float32
	TurnRate    float32
	TurnRateMax float32
	Skid        float32
}

type Team struct {
	Name           string
	LogoModel      int
	Pilots         [NumPilotsPerTeam]Pilot
	TeamAttributes [NumRaceClasses]TeamAttributes
}

type CircuitSettings struct {
	Path         string
	StartLinePos float32
	BehindSpeed  float32
	SpreadBase   float32
	SpreadFactor float32
	SkyYOffset   float32
}

type Circuit struct {
	Name           string
	IsBonusCircuit bool
	Settings       [NumRaceClasses]CircuitSettings
}

type MusicTrack struct {
	Path string
	Name string
}

type PilotPoints struct {
    Pilot  uint16
    Points uint16
}

type Game struct {
    FrameTime          float32
    FrameRate          float32

    RaceClass          int
    RaceType           int
    HighscoreTab       int
    Team               int
    Pilot              int
    Circuit            int
    IsAttractMode      bool
    ShowCredits        bool

    IsNewLapRecord     bool
    IsNewRaceRecord    bool
    BestLap            float32
    RaceTime           float32
    Lives              int
    RacePosition       int

    LapTimes           [NumPilots][NumLaps]float32
    RaceRanks          [NumPilots]PilotPoints
    ChampionshipRanks  [NumPilots]PilotPoints

	CurrentScene GameScene
	NextScene GameScene

	GlobalTextureLen int

	// TODO add camera droid ship and track

	render engine.Render
	platform engine.PlatformSdl
	ui UI

}

func NewGame(render engine.Render, platform engine.PlatformSdl, ui UI) *Game  {
	return &Game{
		render: render,
		platform: platform,
		ui: ui,
		CurrentScene: GameSceneNone,
		NextScene: GameSceneNone,
		GlobalTextureLen: 0,
	}
}

func (g *Game) Init()  {
	// TODO uncomment when save is ready
	// g.platform.SetFullscreen(false)
	// g.render.SetResolution()
	// g.render.SetPostEffect()

	g.ui.Load()
	
}