package main

import (
	"log"
	"net/url"
)

func UpdateSmartTVApp(m3uURL string, macAddress string) error {

	log.Println("update m3uURL of device with mac address: " + macAddress)

	resp, err := call("POST", "https://siptv.app/scripts/up_url_only.php", url.Values{
		"mac":           {macAddress},
		"sel_countries": {"OSN"},
		"sel_logos":     {"0"},
		"detect_epg":    {"on"},
		"lang":          {"en"},
		"url1":          {m3uURL},
		"epg1":          {""},
		"pin":           {""},
		"url_count":     {"1"},
		"file_selected": {"0"},
		"plist_order":   {"0"},
	}, "")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
