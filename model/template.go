package model

const ReleaseInfoTemplate = `* MySQL Version {{.Version}}
* Changes note:{{range .Info}}
	* {{.Version}} | {{.RelType}} {{if .IsRel}}{{else}}| Not yet released{{end}}
		* {{.URL}}{{end}}
`
