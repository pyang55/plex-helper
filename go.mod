module plex-helper

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.2 // indirect
	github.com/aws/aws-sdk-go v1.40.42 // indirect
	github.com/c-bata/go-prompt v0.2.6
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/manifoldco/promptui v0.8.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/stromland/cobra-prompt v0.0.0-20181123224253-940a0a2bd0d3
)

replace github.com/c-bata/go-prompt => github.com/pyang55/go-prompt v0.2.7

replace github.com/stromland/cobra-prompt => github.com/pyang55/cobra-prompt v0.2.0
