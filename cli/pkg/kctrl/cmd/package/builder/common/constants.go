package common

const (
	PkgBuildFileName             = "package-build.yml"
	PkgFetchContentAnnotationKey = "fetch-content-from"
)

const (
	FetchReleaseArtifactFromGithub string = "Release artifact from Github Repository"
	FetchManifestFromGithub        string = "Git Repository(Not supported)"
	FetchChartFromHelmRepo         string = "Helm Chart from Helm Repository"
	FetchChartFromGithub           string = "Helm Chart from Github repository(Not supported)"
)
