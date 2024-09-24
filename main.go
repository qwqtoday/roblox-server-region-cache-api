package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	security_token string
	httpClient     *http.Client
)

func init() {
	var success bool
	security_token, success = os.LookupEnv("ROBLOX_SECURITY_TOEKN")
	if !success {
		panic(errors.New("security token not found in env"))
	}
	httpClient = http.DefaultClient
}

type JoinGameInstanceRequest struct {
	PlaceId           uint64 `json:"placeId"`
	IsTeleport        bool   `json:"isTeleport"`
	GameId            string `json:"gameId"`
	GameJoinAttemptId string `json:"gameJoinAttemptId"`
}

type ServerConnection struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type UdmuxEndpoint struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type GameJoinMetadata struct {
	JoinSource              int    `json:"JoinSource"`
	RequestType             int    `json:"RequestType"`
	MatchmakingDecisionId   string `json:"MatchmakingDecisionId"`
	IsPlaceVoiceChatEnabled bool   `json:"IsPlaceVoiceChatEnabled"`
}

type Session struct {
	SessionId                string                 `json:"SessionId"`
	GameId                   string                 `json:"GameId"`
	PlaceId                  int                    `json:"PlaceId"`
	UniverseId               int                    `json:"UniverseId"`
	ClientIpAddress          string                 `json:"ClientIpAddress"`
	PlatformTypeId           int                    `json:"PlatformTypeId"`
	SessionStarted           string                 `json:"SessionStarted"`
	BrowserTrackerId         int                    `json:"BrowserTrackerId"`
	PartyId                  *string                `json:"PartyId"` // Pointer to allow null
	Age                      float64                `json:"Age"`
	Latitude                 float64                `json:"Latitude"`
	Longitude                float64                `json:"Longitude"`
	CountryId                int                    `json:"CountryId"`
	PolicyCountryId          *int                   `json:"PolicyCountryId"` // Pointer to allow null
	LanguageId               int                    `json:"LanguageId"`
	SupportedLanguageId      int                    `json:"SupportedLanguageId"`
	BlockedPlayerIds         []int                  `json:"BlockedPlayerIds"`
	JoinType                 string                 `json:"JoinType"`
	PlaySessionFlags         int                    `json:"PlaySessionFlags"`
	MatchmakingDecisionId    *string                `json:"MatchmakingDecisionId"` // Pointer to allow null
	GameJoinMetadata         GameJoinMetadata       `json:"GameJoinMetadata"`
	RandomSeed2              string                 `json:"RandomSeed2"`
	IsUserVoiceChatEnabled   bool                   `json:"IsUserVoiceChatEnabled"`
	IsUserAvatarVideoEnabled bool                   `json:"IsUserAvatarVideoEnabled"`
	SourcePlaceId            *int                   `json:"SourcePlaceId"` // Pointer to allow null
	PlayerSignals            map[string]interface{} `json:"PlayerSignals"`
	DeviceType               int                    `json:"DeviceType"`
	GameJoinRegion           string                 `json:"GameJoinRegion"`
}

type JoinScript struct {
	ClientPort                   int                `json:"ClientPort"`
	MachineAddress               string             `json:"MachineAddress"`
	ServerPort                   int                `json:"ServerPort"`
	ServerConnections            []ServerConnection `json:"ServerConnections"`
	UdmuxEndpoints               []UdmuxEndpoint    `json:"UdmuxEndpoints"`
	DirectServerReturn           bool               `json:"DirectServerReturn"`
	PepperId                     int                `json:"PepperId"`
	TokenValue                   string             `json:"TokenValue"`
	PingUrl                      string             `json:"PingUrl"`
	PingInterval                 int                `json:"PingInterval"`
	UserName                     string             `json:"UserName"`
	DisplayName                  string             `json:"DisplayName"`
	HasVerifiedBadge             bool               `json:"HasVerifiedBadge"`
	SeleniumTestMode             bool               `json:"SeleniumTestMode"`
	UserId                       int64              `json:"UserId"`
	RobloxLocale                 string             `json:"RobloxLocale"`
	GameLocale                   string             `json:"GameLocale"`
	SuperSafeChat                bool               `json:"SuperSafeChat"`
	FlexibleChatEnabled          bool               `json:"FlexibleChatEnabled"`
	CharacterAppearance          string             `json:"CharacterAppearance"`
	ClientTicket                 string             `json:"ClientTicket"`
	GameId                       string             `json:"GameId"`
	PlaceId                      int                `json:"PlaceId"`
	BaseUrl                      string             `json:"BaseUrl"`
	ChatStyle                    string             `json:"ChatStyle"`
	CreatorId                    int                `json:"CreatorId"`
	CreatorTypeEnum              string             `json:"CreatorTypeEnum"`
	MembershipType               string             `json:"MembershipType"`
	AccountAge                   int                `json:"AccountAge"`
	CookieStoreFirstTimePlayKey  string             `json:"CookieStoreFirstTimePlayKey"`
	CookieStoreFiveMinutePlayKey string             `json:"CookieStoreFiveMinutePlayKey"`
	CookieStoreEnabled           bool               `json:"CookieStoreEnabled"`
	IsUnknownOrUnder13           bool               `json:"IsUnknownOrUnder13"`
	GameChatType                 string             `json:"GameChatType"`
	SessionId                    string             `json:"SessionId"`
	AnalyticsSessionId           string             `json:"AnalyticsSessionId"`
	DataCenterId                 int                `json:"DataCenterId"`
	UniverseId                   int                `json:"UniverseId"`
	FollowUserId                 int                `json:"FollowUserId"`
	CharacterAppearanceId        int64              `json:"characterAppearanceId"`
	CountryCode                  string             `json:"CountryCode"`
	RandomSeed1                  string             `json:"RandomSeed1"`
	ClientPublicKeyData          string             `json:"ClientPublicKeyData"`
	RccVersion                   string             `json:"RccVersion"`
	ChannelName                  string             `json:"ChannelName"`
	VerifiedAMP                  int                `json:"VerifiedAMP"`
	PrivateServerOwnerID         int                `json:"PrivateServerOwnerID"`
	PrivateServerID              string             `json:"PrivateServerID"`
	EventID                      string             `json:"EventID"`
	EphemeralEarlyPubKey         string             `json:"EphemeralEarlyPubKey"`
	PartyId                      string             `json:"PartyId"`
}

type JoinGameInstanceResponse struct {
	JobId                string       `json:"jobId"`
	Status               int          `json:"status"`
	StatusData           *interface{} `json:"statusData"` // Pointer to allow null
	JoinScriptUrl        string       `json:"joinScriptUrl"`
	AuthenticationUrl    string       `json:"authenticationUrl"`
	AuthenticationTicket string       `json:"authenticationTicket"`
	Message              *interface{} `json:"message"` // Pointer to allow null
	JoinScript           JoinScript   `json:"joinScript"`
	QueuePosition        int          `json:"queuePosition"`
}

func main() {
	serverIpsCache := cache.New(5*time.Minute, 10*time.Minute)

	app := fiber.New()
	// Define a route for the GET method on the root path '/'
	app.Get("/place/:place_id/ip/:server_id", func(c fiber.Ctx) error {
		var err error
		placeId, err := strconv.ParseUint(c.Params("place_id", ""), 10, 64)
		if err != nil {
			return err
		}
		serverId := c.Params("server_id", "")

		cacheId := fmt.Sprintf("%d:%s", placeId, serverId)
		ip, found := serverIpsCache.Get(cacheId)
		if found {
			return c.SendString(ip.(string))
		}
		data := JoinGameInstanceRequest{
			PlaceId:           placeId,
			IsTeleport:        false,
			GameId:            serverId,
			GameJoinAttemptId: serverId,
		}

		b, err := json.Marshal(data)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "https://gamejoin.roblox.com/v1/join-game-instance", bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Referer", fmt.Sprintf("https://www.roblox.com/games/%d", placeId))
		req.Header.Set("User-Agent", "Roblox/WinInet")

		req.AddCookie(&http.Cookie{
			Name:   ".ROBLOSECURITY",
			Value:  security_token,
			Domain: ".roblox.com",
			Path:   "/",
		})

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		var result JoinGameInstanceResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		if len(result.JoinScript.UdmuxEndpoints) == 0 {
			return errors.New("no udmux endpoints found")
		}

		ip = result.JoinScript.UdmuxEndpoints[0].Address
		serverIpsCache.SetDefault(
			cacheId,
			ip,
		)

		return c.SendString(ip.(string))
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
