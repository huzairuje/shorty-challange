package utils

import "regexp"

const (
	NotMatchingAnyRoute         = "Not Matching of Any Routes"
	SomethingWentWrong          = "Oops, Something Went Wrong"
	BadRequest                  = "Bad Request"
	NotFound                    = "Not Found"
	UrlIsNotSet                 = "url not present"
	ShortCodeExist              = "The desired shortcode is already in use. Shortcodes are case-sensitive"
	ShortCodeFailedRegexPattern = "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{6}$."
	ShortCodeIsNotExist         = "the shortcode cannot be found in the system"
)

func IsValidShortCode(shortCode string) bool {
	regexQueryParam := regexp.MustCompile(`^[0-9a-zA-Z_]{6}$`)
	return regexQueryParam.MatchString(shortCode)
}
