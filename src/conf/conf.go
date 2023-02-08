package conf

// Config
const (
	AppName             = "Omega Drive"            // specify program name
	IsDev               = false                    // enable extra debugging info if set to true
	CheckParentInterval = "2s"                     // specify interval for checking parent PID. Examples are (without quotes): '2h','5m','14s'
	RcUsername          = "someUsername"           // specify username for rclone rcd server authentication
	RcPassword          = "somePassword"           // specify password for rclone rcd server authentication
	RcHost              = "http://localhost:5579/" // specify hostname for rclone rcd
	LogFilename         = "rcd.log"
)
