package utils

import (
	"fmt"
	"strconv"
)

func GenPages(firstUrl string) []string {
	var convertedUrl []string
	ur := `https://freelance.habr.com/tasks?categories=development_all_inclusive%2Cdevelopment_backend%2Cdevelopment_frontend%2Cdevelopment_prototyping%2Cdevelopment_ios%2Cdevelopment_android%2Cdevelopment_desktop%2Cdevelopment_bots%2Cdevelopment_games%2Cdevelopment_1c_dev%2Cdevelopment_scripts%2Cdevelopment_voice_interfaces%2Cdevelopment_other&page=`
	for page := 0; page <= 3; page++ {
		url := fmt.Sprint(ur + strconv.Itoa(page))
		convertedUrl = append(convertedUrl, url)
	}
	convertedUrl[0] = firstUrl
	return convertedUrl
}
