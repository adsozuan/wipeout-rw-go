package game

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

type GameSceneE int

const (
	GameSceneIntro GameSceneE = iota
	GameSceneTitle
	GameSceneMainMenu
	GameSceneHighscores
	GameSceneRace
	GameSceneNone
	NumGameScenes
)

func (g GameSceneE) String() string {
	names := [...]string{
		"GameSceneIntro",
		"GameSceneTitle",
		"GameSceneMainMenu",
		"GameSceneHighscores",
		"GameSceneRace",
		"GameSceneNone",
		"NumGameScenes",
	}

	if g < GameSceneIntro || g >= NumGameScenes {
		return "Unknown"
	}

	return names[g]
}


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

type (
	Init   func()
	Update func()
)

type GameScene interface {
	Init() error
	Update() error
}

type Game struct {
	FrameTime float64
	FrameRate float64

	RaceClass     int
	RaceType      int
	HighscoreTab  int
	Team          int
	Pilot         int
	Circuit       int
	IsAttractMode bool
	ShowCredits   bool

	IsNewLapRecord  bool
	IsNewRaceRecord bool
	BestLap         float32
	RaceTime        float32
	Lives           int
	RacePosition    int

	LapTimes          [NumPilots][NumLaps]float32
	RaceRanks         [NumPilots]PilotPoints
	ChampionshipRanks [NumPilots]PilotPoints

	CurrentScene GameSceneE
	NextScene    GameSceneE

	GameScenes map[GameSceneE]GameScene

	GlobalTextureLen int

	// TODO add camera droid ship and track

	render   *engine.Render
	platform *engine.PlatformSdl
	ui       *UI
}

func NewGame(render *engine.Render, platform *engine.PlatformSdl) (*Game, error) {
	Logger = log.New(os.Stderr, "game   |", log.Ldate|log.Ltime)
	Logger.Println("Init")
	ui := NewUI(render)

	return &Game{
		render:           render,
		platform:         platform,
		ui:               ui,
		CurrentScene:     GameSceneNone,
		NextScene:        GameSceneNone,
		GlobalTextureLen: 0,
	}, nil
}

func (g *Game) Init(startTime float64) error {
	// TODO uncomment when save is ready
	// g.platform.SetFullscreen(false)
	// g.render.SetResolution()
	// g.render.SetPostEffect()

	err := g.ui.Load()
	if err != nil {
		return err
	}

	g.GlobalTextureLen = g.render.TexturesLen()

	g.GameScenes = make(map[GameSceneE]GameScene)
	g.GameScenes[GameSceneTitle] = NewTitle(startTime, g.render)

	g.SetScene(GameSceneTitle)

	return nil
}

func (g *Game) SetScene(scene GameSceneE) {
	// TODO reset sfx
	g.NextScene = scene

	Logger.Println(g.NextScene)
}

type ResetCycleTime bool

func (g *Game) Update() ResetCycleTime {
	frameStartTime := g.platform.Now()
	resetCycleTime := false

	sh := int(g.render.Size().Y)

	var scale int
	if sh >= 720 {
		scale = sh / 360
	} else {
		scale = sh / 240
	}
	scale = max(1, scale)

	// TODO save scale

	g.ui.SetScale(scale)

	if g.NextScene != GameSceneNone {
		g.CurrentScene = g.NextScene
		g.NextScene = GameSceneNone
		g.render.TexturesReset(uint16(g.GlobalTextureLen))
		resetCycleTime = true

		if g.CurrentScene != GameSceneNone {
			g.GameScenes[g.CurrentScene].Init()
		}
	}

	if g.CurrentScene != GameSceneNone {
		g.GameScenes[g.CurrentScene].Update()
	}

	// TODO handle save

	now := g.platform.Now()
	g.FrameTime = now - frameStartTime

	if g.FrameTime > 0 {
		g.FrameRate = g.FrameRate*0.95 + 1.0/g.FrameTime*0.05
	}

	return ResetCycleTime(resetCycleTime)
}
