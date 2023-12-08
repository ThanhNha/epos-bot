package phab_bot

import (
	"log"

	"github.com/spf13/viper"
)

func GetNotifyTypes() []string {
	notifyTypes := viper.GetStringSlice("telegram.notify_types")
	if len(notifyTypes) == 0 {
		notifyTypes = []string{"TASK", "DREV", "WIKI"}
		log.Printf("No notify types specified, defaulting to %v", notifyTypes)
	}
	return notifyTypes
}

func CreateNotifyTypesMap(notifyTypes []string) map[string]bool {
	notifyTypesMap := make(map[string]bool)
	for _, v := range notifyTypes {
		notifyTypesMap[v] = true
	}
	return notifyTypesMap
}
