package url

import "os"

const (
	baseBeamTechSandboxURL = "wise.api.beamtechnology.com:8888"
	baseBeamTechProdURL    = "wise.api.beamtechnology.com:8888"
)

func BuildURLForBeamTech(path string, urlParams Params) (string, error) {
	return BuildAbsolute(path, getBaseBeamTechURLForENV(), urlParams)
}

func BaseURLForBeanTech() string {
	return getBaseBeamTechURLForENV()
}

func getBaseBeamTechURLForENV() string {
	baseURL := protocol

	// Get URL by environment
	env := os.Getenv("API_ENV")
	switch env {
	case envDev:
		baseURL = baseURL + baseBeamTechSandboxURL
	case envStaging:
		baseURL = baseURL + baseBeamTechSandboxURL
	case envQA:
		baseURL = baseURL + baseBeamTechSandboxURL
	case envProd:
		baseURL = baseURL + baseBeamTechProdURL
	default:
		baseURL = baseURL + baseBeamTechSandboxURL
	}

	return baseURL
}
