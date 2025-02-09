package Parameter

var ( //配置键值对，默认值
	MaxCoroutine                 = 512
	MinCoroutine                 = 256
	ConfigPath                   = "./conf/ids.conf"
	LogPath                      = "./log"
	AlertPath                    = "./log/alert.log"
	HardwareDataRefreshFrequency = 500
	IPBanAllowedDiff             = 5
	EmailSendInterval            = 5 * 1000 * 1000 * 1000 * 60

	SqlUser                = "root"
	SqlPassword            = "root"
	SqlHost                = "127.0.0.1:3307"
	SqlSchemaName          = "IDSdata"
	SqlRecordTableName     = "ActionDescribeTriggerLog"
	SqlRecordTableFullName string

	WhiteListPath     string
	BlackListPath     string
	BlackPortListPath string

	SaveIptablesShPath     string
	IptablesConfigSavePath string

	IdsMode = 1
	IdsName = "node1"
	IdsID   = 1

	Account  = "a"
	Password = "12345"

	SmtpHost     string
	SmtpPort     int
	ToEmail      string
	AuthEmail    string
	AuthPassword string

	InitBanIpShPath     string
	AddIpToBanShPath    string
	DeleteIpFromShPath  string
	SaveBanIpShPath     string
	BanIpConfigSavePath string

	InitBanIpNetWorkShPath     string
	AddIpNetworkToBanShPath    string
	DeleteIpNetworkFromShPath  string
	SaveBanIpNetworkShPath     string
	BanIpNetworkConfigSavePath string

	InitWhiteListShPath    string
	AddIpToWhiteListShPath string
	DeleteWhiteIpSh        string
	SaveWhiteListShPath    string
	WhiteListSavePath      string

	InitWhiteListNetworkShPath    string
	AddIpToWhiteListNetworkShPath string
	DeleteWhiteIpNetworkSh        string
	SaveWhiteListNetworkShPath    string
	WhiteListNetworkSavePath      string

	InitBanPortShPath      string
	AddPortToBanShPath     string
	DeleteBanPortFromBanSh string
	SaveBanPortShPath      string
	BanPortSavePath        string
)
