package main

type EmailTemplateData struct {
	Subject           string
	TempBoard         string
	MemoryTotal       string
	MemoryUsed        string
	MemoryUsedPercent string
	CPUUsagePercent   string
	RxBytes           string
	TxBytes           string
	UpTime            string
	DateTime          string
	ServerName        string
	HostName          string
}

// TODO: Darkmode and/or banner
const EmailTemplateStr = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}}</title>
</head>
<body>
    <h4>Email alert for <b>{{.ServerName}}</b> <i>[{{.HostName}}]</i></h4>
	<ul>
		<li>Board temperature: {{.TempBoard}}</li>
		<li>Memory Usage: {{.MemoryUsedPercent}} [{{.MemoryUsed}}/{{.MemoryTotal}}]</li>
		<li>CPU Usage: {{.CPUUsagePercent}}</li>
		<li>Network: Rx:{{.RxBytes}} Tx:{{.TxBytes}}</li>
		<li>UpTime: {{.UpTime}}</li>
		<li>DateTime: {{.DateTime}}</li>
		<li>ServerName: {{.ServerName}}</li>
		<li>HostName: {{.HostName}}</li>
	</ul>
	<br />
	<p>EOM</p>
</body>
</html>
`
