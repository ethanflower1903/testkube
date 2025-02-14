package analytics

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"

	"github.com/denisbrodbeck/machineid"
	v1 "github.com/mjpitz/go-ga/client/v1"
	"github.com/mjpitz/go-ga/client/v1/gatypes"

	"github.com/kubeshop/testkube/cmd/tools/commands"
)

var testkubeTrackingID = "UA-204665550-8" //this is default but it can be set using ldflag -X github.com/kubeshop/testkube/pkg/analytics.testkubeTrackingID=UA-204665550-8

// SendAnonymousInfo will send event to GA
func SendAnonymousInfo() {

	var isEnabled bool
	if val, ok := os.LookupEnv("TESTKUBE_ANALYTICS_ENABLED"); ok {
		isEnabled, _ = strconv.ParseBool(val)
	}
	if isEnabled {
		client := v1.NewClient(testkubeTrackingID, "golang")
		payload := &gatypes.Payload{
			HitType:                           "event",
			DataSource:                        "api-server",
			NonInteractionHit:                 true,
			DisableAdvertisingPersonalization: true,
			Users: gatypes.Users{
				ClientID: MachineID(),
			},
			Event: gatypes.Event{
				EventCategory: "beacon",
				EventAction:   "testkube-heartbeat",
			},
			Apps: gatypes.Apps{
				ApplicationName:    "testkube",
				ApplicationVersion: commands.Version,
			},
		}
		client.SendPost(payload)
	}
}

// SendAnonymouscmdInfo will send CLI event to GA
func SendAnonymouscmdInfo() {
	client := v1.NewClient(testkubeTrackingID, "golang")
	event := "command"
	command := []string{}
	if len(os.Args) > 1 {
		command = os.Args[1:]
		event = command[0]
	}
	payload := &gatypes.Payload{
		HitType:                           "event",
		DataSource:                        "CLI",
		DisableAdvertisingPersonalization: true,
		Users: gatypes.Users{
			ClientID: MachineID(),
		},
		Event: gatypes.Event{
			EventCategory: event,
			EventAction:   "execution",
			EventLabel:    event,
		},
		Apps: gatypes.Apps{
			ApplicationName:    "testkube",
			ApplicationVersion: commands.Version,
		},
		CustomDimensions: gatypes.StringList(command),
	}
	client.SendPost(payload)
}

// MachineID returns unique user machine ID
func MachineID() string {
	id, _ := generate()
	return id
}

// Generate returns protected id for the current machine
func generate() (string, error) {
	id, err := machineid.ProtectedID("testkube")
	if err != nil {
		return fromHostname()
	}
	return id, err
}

// fromHostname generates a machine id hash from hostname
func fromHostname() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:]), nil
}
