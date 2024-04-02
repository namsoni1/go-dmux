package sideline_impls

import (
	"fmt"
	co "github.com/flipkart-incubator/go-dmux/config"
	"github.com/flipkart-incubator/go-dmux/logging"
	"github.com/flipkart-incubator/go-dmux/metrics"
	"log"
)

//

type DmuxCustom struct {
}

func (d *DmuxCustom) DmuxStart(path string, sidelineImp interface{}) {
	//log.Println(checkMessageSideline.SidelineMessage())

	dconf := co.DMuxConfigSetting{
		FilePath: path,
	}
	conf := dconf.GetDmuxConf()

	dmuxLogging := new(logging.DMuxLogging)
	//_ = new(logging.DMuxLogging)

	fmt.Printf("config of dmuxstart: %v \n", conf)

	//start showing metrics at the endpoint
	metrics.Start(conf.MetricPort)

	for _, item := range conf.DMuxItems {
		log.Println(item.ConnType)
		if item.SidelineEnable {
			fmt.Printf("sideline enabled: %v \n", conf)
			go func(connType co.ConnectionType, connConf interface{}, logDebug bool) {
				connType.Start(connConf, logDebug, sidelineImp)
			}(item.ConnType, item.Connection, dmuxLogging.EnableDebug)
		} else {
			fmt.Printf("sideline disabled: %v \n", conf)
			go func(connType co.ConnectionType, connConf interface{}, logDebug bool) {
				connType.Start(connConf, logDebug, nil)
			}(item.ConnType, item.Connection, dmuxLogging.EnableDebug)
		}
	}

	//main thread halts. TODO make changes to listen to kill and reboot
	select {}
}
