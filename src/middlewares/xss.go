package middlewares

import "github.com/microcosm-cc/bluemonday"

func XSSMiddleware(param map[string]interface{}) map[string]interface{} {
	policy := bluemonday.UGCPolicy()

	for key, value := range param {
		if str, ok := value.(string); ok {
			param[key] = policy.Sanitize(str)
		}
	}
	return param
}
