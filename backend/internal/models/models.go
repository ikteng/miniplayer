// Package models mirrors the Pydantic models in backend/app/models/entities.py
// and backend/app/models/schemas.py. Field names and JSON tags match the
// Python `model_dump(by_alias=True)` output exactly so the frontend sees an
// identical wire format.
package models

// Player mirrors app.models.entities.Player.
type Player struct {
	HasSong                     bool    `json:"has_song"`
	IsPaused                    bool    `json:"is_paused"`
	SeekbarCurrentPosition      int     `json:"seekbar_current_position"`
	SeekbarCurrentPositionHuman string  `json:"seekbar_current_position_human"`
	StatePercent                float64 `json:"state_percent"`
}

func NewPlayer() Player {
	return Player{
		HasSong:                     false,
		IsPaused:                    true,
		SeekbarCurrentPosition:      0,
		SeekbarCurrentPositionHuman: "0:00",
		StatePercent:                0.0,
	}
}

// Track mirrors app.models.entities.Track.
type Track struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	Album         string `json:"album"`
	Cover         string `json:"cover"`
	Duration      int    `json:"duration"`
	DurationHuman string `json:"duration_human"`
	URL           string `json:"url"`
}

func NewTrack() Track {
	return Track{DurationHuman: "0:00"}
}

// Query mirrors app.models.entities.Query.
type Query struct {
	Player Player `json:"player"`
	Track  Track  `json:"track"`
}

// QueryProgress mirrors app.models.entities.QueryProgress.
type QueryProgress struct {
	ProgressMs int `json:"progress_ms"`
}

// RespData mirrors app.models.entities.RespData.
type RespData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespData(code int, msg string, data interface{}) RespData {
	return RespData{Code: code, Msg: msg, Data: data}
}

// Base64Img mirrors app.models.entities.Base64Img (alias base64Img).
type Base64Img struct {
	Base64Img string `json:"base64Img"`
	Base64    string `json:"base64"`
}

// Lyric mirrors app.models.entities.Lyric.
type Lyric struct {
	Title              string `json:"title"`
	Author             string `json:"author"`
	Duration           int    `json:"duration"`
	Lyric              string `json:"lyric"`
	TLyric             string `json:"tlyric"`
	YrcLyric           string `json:"yrc_lyric"`
	YrcTLyric          string `json:"yrc_tlyric"`
	QrcLyric           string `json:"qrc_lyric"`
	QrcTLyric          string `json:"qrc_tlyric"`
	Source             string `json:"source"`
	HasLyric           bool   `json:"has_lyric"`
	HasTranslatedLyric bool   `json:"has_translated_lyric"`
	HasKaraokeLyric    bool   `json:"has_karaoke_lyric"`
}

func NewLyric() Lyric {
	return Lyric{Source: "netease"}
}

// SettingsGeneral mirrors app.models.entities.SettingsGeneral.
type SettingsGeneral struct {
	DeviceID                string `json:"device_id"`
	DeviceName              string `json:"device_name"`
	Platform                string `json:"platform"`
	AutoLaunchHomePage      bool   `json:"auto_launch_home_page"`
	RunAtStartup            bool   `json:"run_at_startup"`
	UpdateCheckFreq         int    `json:"update_check_freq"`
	SMTC                    bool   `json:"smtc"`
	FallbackPlatformEnabled bool   `json:"fallback_platform_enabled"`
	FallbackPlatform        string `json:"fallback_platform"`
	PollInterval            int    `json:"poll_interval"`
	WeSingCachePath         string `json:"we_sing_cache_path"`
	SplayerNextPort         int    `json:"splayer_next_port"`
}

func NewSettingsGeneral() SettingsGeneral {
	return SettingsGeneral{
		DeviceID:                "default",
		DeviceName:              "Default Audio Driver",
		Platform:                "netease",
		AutoLaunchHomePage:      true,
		RunAtStartup:            false,
		UpdateCheckFreq:         0,
		SMTC:                    true,
		FallbackPlatformEnabled: false,
		FallbackPlatform:        "netease",
		PollInterval:            100,
		WeSingCachePath:         "",
		SplayerNextPort:         14558,
	}
}

// SettingsLyric mirrors app.models.entities.SettingsLyric.
type SettingsLyric struct {
	UseCustomConfig        bool   `json:"use_custom_config"`
	ConfigFileName         string `json:"config_file_name"`
	LyricColor             string `json:"lyric_color"`
	LyricFontSize          int    `json:"lyric_font_size"`
	LyricFontWeight        int    `json:"lyric_font_weight"`
	StrokeColor            string `json:"stroke_color"`
	StrokeSize             int    `json:"stroke_size"`
	TranslationColor       string `json:"translation_color"`
	TranslationFontSize    int    `json:"translation_font_size"`
	TranslationFontWeight  int    `json:"translation_font_weight"`
	TranslationStrokeColor string `json:"translation_stroke_color"`
	TranslationStrokeSize  int    `json:"translation_stroke_size"`
	Gap                    int    `json:"gap"`
	TranslationGap         int    `json:"translation_gap"`
	AlbumCoverLyric        bool   `json:"album_cover_lyric"`
	Animation              bool   `json:"animation"`
	BackgroundColor        string `json:"background_color"`
	PositionType           string `json:"position_type"`
	TypeLyric              bool   `json:"type_lyric"`
	KaraokeAnimation       bool   `json:"karaoke_animation"`
	VerticalAlign          bool   `json:"vertical_align"`
	ShowInterlude          bool   `json:"show_interlude"`
	ShowTranslation        bool   `json:"show_translation"`
	FullScreenLyric        bool   `json:"full_screen_lyric"`
	EnableHover            bool   `json:"enable_hover"`
	OffsetY                int    `json:"offset_y"`
	LastUpdateTime         int    `json:"last_update_time"`
}

func NewSettingsLyric() SettingsLyric {
	return SettingsLyric{
		LyricColor:             "#f5f5f5",
		LyricFontSize:          28,
		LyricFontWeight:        700,
		StrokeColor:            "#000000",
		TranslationColor:       "#f5f5f5",
		TranslationFontSize:    28,
		TranslationFontWeight:  700,
		TranslationStrokeColor: "#000000",
		BackgroundColor:        "#4d000000",
		PositionType:           "bottom",
		KaraokeAnimation:       true,
		VerticalAlign:          true,
		ShowInterlude:          true,
		ShowTranslation:        true,
	}
}

// SettingsLyricCommon mirrors app.models.entities.SettingsLyricCommon.
type SettingsLyricCommon struct {
	LyricSource         string `json:"lyric_source"`
	AutoSelectBestLyric bool   `json:"auto_select_best_lyric"`
	Romanization        bool   `json:"romanization"`
	Translate           bool   `json:"translate"`
	ChineseConversion   string `json:"chinese_conversion"`
}

func NewSettingsLyricCommon() SettingsLyricCommon {
	return SettingsLyricCommon{
		LyricSource:         "netease",
		AutoSelectBestLyric: true,
		Translate:           true,
		ChineseConversion:   "off",
	}
}

// SettingsOutput mirrors app.models.entities.SettingsOutput.
type SettingsOutput struct {
	OutputTitle    bool   `json:"output_title"`
	OutputAuthor   bool   `json:"output_author"`
	OutputCover    bool   `json:"output_cover"`
	OutputCustom   bool   `json:"output_custom"`
	Template       string `json:"template"`
	EnableWhenNone bool   `json:"enable_when_none"`
	CustomNone     string `json:"custom_none"`
}

func NewSettingsOutput() SettingsOutput {
	return SettingsOutput{
		OutputTitle:  true,
		OutputAuthor: true,
		OutputCover:  true,
		OutputCustom: true,
		Template:     "{title} - {author}",
	}
}

// SettingsWidget mirrors app.models.entities.SettingsWidget.
type SettingsWidget struct {
	Cover           int    `json:"cover"`
	CoverBlur       bool   `json:"cover_blur"`
	CoverGlow       bool   `json:"cover_glow"`
	FirstLogin      bool   `json:"first_login"`
	HideDelay       int    `json:"hide_delay"`
	HideEqualizer   bool   `json:"hide_equalizer"`
	HideOnPause     bool   `json:"hide_on_pause"`
	IsDemo          bool   `json:"is_demo"`
	Player          int    `json:"player"`
	PlayerColors    bool   `json:"player_colors"`
	SongChangeOnly  bool   `json:"song_change_only"`
	Theme           int    `json:"theme"`
	TintColor       string `json:"tint_color"`
	VisibleDuration int    `json:"visible_duration"`
}

func NewSettingsWidget() SettingsWidget {
	return SettingsWidget{
		HideDelay:       10,
		PlayerColors:    true,
		TintColor:       "#40a4d8",
		VisibleDuration: 5,
	}
}

// SettingsWrapper mirrors app.models.entities.SettingsWrapper.
type SettingsWrapper struct {
	General     SettingsGeneral        `json:"general"`
	Lyric       SettingsLyric          `json:"lyric"`
	LyricCommon SettingsLyricCommon    `json:"lyric_common"`
	Output      SettingsOutput         `json:"output"`
	Widget      SettingsWidget         `json:"widget"`
	Plugins     map[string]interface{} `json:"plugins"`
}

func NewSettingsWrapper() SettingsWrapper {
	return SettingsWrapper{
		General:     NewSettingsGeneral(),
		Lyric:       NewSettingsLyric(),
		LyricCommon: NewSettingsLyricCommon(),
		Output:      NewSettingsOutput(),
		Widget:      NewSettingsWidget(),
		Plugins:     map[string]interface{}{},
	}
}

// CoverConvertDTO mirrors app.models.entities.CoverConvertDTO.
type CoverConvertDTO struct {
	CoverURL string `json:"cover_url"`
}

// CoverVideoDTO mirrors app.models.entities.CoverVideoDTO.
type CoverVideoDTO struct {
	SongTitle  string `json:"songTitle"`
	SongAuthor string `json:"songAuthor"`
}
