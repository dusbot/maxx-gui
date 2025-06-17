package report

import (
	"bytes"
	"fmt"
	"html/template"
	"maxxgui/backend/model"
	"sort"
	"strings"
	"time"

	"github.com/dusbot/maxx/libs/slog"
	utils_ "github.com/dusbot/maxx/libs/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var ()

func DoGenCrackReport(zh bool, task *model.CrackTask, results []*model.CrackResult) (ok bool, content string) {
	langs := map[string]map[string]string{
		"zh": {
			"ReportTitle":  "安全评估报告",
			"ToggleTheme":  "切换主题",
			"ScoreLabel":   "整体安全评分",
			"StatusGood":   "安全状况良好",
			"StatusRisk":   "存在一定风险",
			"StatusDanger": "高危，请立即整改",
			"TaskInfo":     "任务信息",
			"TaskID":       "任务ID",
			"StartTime":    "开始时间",
			"EndTime":      "结束时间",
			"Targets":      "目标",
			"Usernames":    "用户字典",
			"Passwords":    "密码字典",
			"Proxies":      "代理",
			"Thread":       "线程数",
			"Interval":     "扫描间隔",
			"MaxRuntime":   "最大运行时间设置",
			"TotalCrack":   "破解总数",
			"WeakCount":    "弱口令数",
			"Analysis":     "统计分析",
			"WeakDetail":   "弱口令结果明细",
			"NoWeak":       "无弱口令结果",
			"TargetCrack":  "各Target破解次数",
			"ServicePie":   "各Service破解占比",
			"Advice":       "安全建议",
			"Advice1":      "定期更换弱口令，使用复杂密码。",
			"Advice2":      "加强对高风险服务的访问控制。",
			"Advice3":      "对暴露在公网的服务进行多因素认证。",
			"Target":       "目标",
			"Service":      "服务",
			"Username":     "用户名",
			"Password":     "密码",
			"Fold":         "（点击折叠/展开）",
		},
		"en": {
			"ReportTitle":  "Security Assessment Report",
			"ToggleTheme":  "Toggle Theme",
			"ScoreLabel":   "Overall Security Score",
			"StatusGood":   "Good Security",
			"StatusRisk":   "Some Risks Exist",
			"StatusDanger": "High Risk, Immediate Action Required",
			"TaskInfo":     "Task Info",
			"TaskID":       "Task ID",
			"StartTime":    "Start Time",
			"EndTime":      "End Time",
			"Targets":      "Targets",
			"Usernames":    "Usernames",
			"Passwords":    "Passwords",
			"Proxies":      "Proxies",
			"Thread":       "Thread",
			"Interval":     "Interval",
			"MaxRuntime":   "Max Runtime",
			"TotalCrack":   "Total Attempts",
			"WeakCount":    "Weak Passwords",
			"Analysis":     "Statistics",
			"WeakDetail":   "Weak Password Details",
			"NoWeak":       "No weak password found",
			"TargetCrack":  "Attempts per Target",
			"ServicePie":   "Service Distribution",
			"Advice":       "Security Advice",
			"Advice1":      "Regularly change weak passwords and use complex ones.",
			"Advice2":      "Strengthen access control for high-risk services.",
			"Advice3":      "Enable multi-factor authentication for public services.",
			"Target":       "Target",
			"Service":      "Service",
			"Username":     "Username",
			"Password":     "Password",
			"Fold":         " (Click to fold/unfold)",
		},
	}
	lang := "en"
	if zh {
		lang = "zh"
	}

	weakResults := make([]map[string]string, 0)
	for _, r := range results {
		if r.Password != "" {
			weakResults = append(weakResults, map[string]string{
				"Target":   r.Target,
				"Service":  r.Service,
				"Username": r.Username,
				"Password": r.Password,
			})
		}
	}

	targetTotalCount := map[string]int{}
	targetWeakCount := map[string]int{}
	serviceCount := map[string]int{}
	weakPasswordCount := 0
	totalTargets := utils_.ParseNetworkInput(task.Targets)
	totalCrackCount := (len(strings.Split(task.Usernames, ",")) + 1) * (len(strings.Split(task.Passwords, ",")) + 1)
	for _, target := range totalTargets {
		targetTotalCount[target] = totalCrackCount
	}

	for _, r := range results {
		serviceCount[r.Service]++
		if r.Password != "" {
			weakPasswordCount++
			targetWeakCount[r.Target]++
		}
	}

	targets := make([]string, 0, len(targetTotalCount))
	for k := range targetTotalCount {
		targets = append(targets, k)
	}
	sort.Strings(targets)

	totalVals := make([]opts.BarData, 0, len(targets))
	weakVals := make([]opts.BarData, 0, len(targets))
	for _, k := range targets {
		totalVals = append(totalVals, opts.BarData{Value: targetTotalCount[k]})
		weakVals = append(weakVals, opts.BarData{Value: targetWeakCount[k]})
	}

	barHeight := 60 + 40*len(targets)
	bar := charts.NewBar()
	showLegendOpt := true
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: fmt.Sprintf("%dpx", barHeight),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "category",
			Data: targets,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{Show: &showLegendOpt}),
	)
	barShow := true
	bar.AddSeries(langs[lang]["TotalCrack"], totalVals).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: &barShow, Position: "right"}))
	bar.AddSeries(langs[lang]["WeakCount"], weakVals).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: &barShow, Position: "right"}))

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "320px",
		}),
	)
	pieData := make([]opts.PieData, 0, len(serviceCount))
	for k, v := range serviceCount {
		pieData = append(pieData, opts.PieData{Name: k, Value: v})
	}
	pie.AddSeries("service", pieData)

	total := 0
	for _, v := range targetTotalCount {
		total += v
	}
	score := 100
	if total > 0 {
		score = 100 - (weakPasswordCount * 100 / total)
	}

	status := langs[lang]["StatusGood"]
	if score < 80 {
		status = langs[lang]["StatusRisk"]
	}
	if score < 60 {
		status = langs[lang]["StatusDanger"]
	}

	advice := []string{
		langs[lang]["Advice1"],
		langs[lang]["Advice2"],
		langs[lang]["Advice3"],
	}

	var barBuf, pieBuf bytes.Buffer
	bar.Render(&barBuf)
	barHTML := barBuf.String()
	pie.Render(&pieBuf)
	pieHTML := pieBuf.String()

	reportTpl := `
<!DOCTYPE html>
<html lang="` + lang + `">
<head>
    <meta charset="UTF-8">
    <title>{{index .Lang "ReportTitle"}}</title>
    <link href="https://fonts.googleapis.com/css?family=Roboto:900,400" rel="stylesheet">
    <style>
        :root {
            --bg-main: linear-gradient(135deg, #232526 0%, #414345 100%);
            --card-bg: rgba(30,40,60,0.95);
            --section-bg: rgba(20,30,40,0.85);
            --text-main: #fff;
            --accent: #00eaff;
            --shadow: #00eaff33;
        }
        [data-theme="light"] {
            --bg-main: linear-gradient(135deg, #f8fafc 0%, #e3e8ee 100%);
            --card-bg: #fff;
            --section-bg: #f5f7fa;
            --text-main: #222;
            --accent: #0077ff;
            --shadow: #0077ff22;
        }
        body {
            background: var(--bg-main);
            color: var(--text-main);
            font-family: 'Roboto', Arial, sans-serif;
            margin: 0; padding: 0;
            transition: background 0.5s, color 0.5s;
        }
        .theme-toggle {
            position: fixed;
            top: 24px; right: 36px;
            background: var(--card-bg);
            color: var(--accent);
            border: none;
            border-radius: 20px;
            padding: 8px 22px;
            font-size: 1em;
            cursor: pointer;
            box-shadow: 0 2px 8px var(--shadow);
            transition: background 0.3s, color 0.3s;
            z-index: 99;
        }
        .header {
            background: rgba(0,0,0,0.7);
            padding: 40px 0 20px 0;
            text-align: center;
            box-shadow: 0 4px 24px #000a;
        }
        .header h1 {
            font-size: 2.8em;
            letter-spacing: 2px;
            margin: 0;
            color: var(--accent);
            text-shadow: 0 2px 12px var(--shadow);
        }
        .score-ring {
            width: 140px; height: 140px;
            margin: 30px auto 10px auto;
            position: relative;
        }
        .score-ring svg {
            transform: rotate(-90deg);
        }
        .score-ring .score-text {
            position: absolute;
            top: 50%; left: 50%;
            transform: translate(-50%, -50%);
            font-size: 2.2em;
            font-weight: bold;
            color: var(--text-main);
            text-shadow: 0 2px 8px #000a;
        }
        .task-info-card {
            max-width: 900px;
            margin: 32px auto 0 auto;
            background: var(--card-bg);
            border-radius: 18px;
            box-shadow: 0 2px 16px var(--shadow);
            padding: 28px 40px 18px 40px;
            display: flex;
            flex-direction: column;
            gap: 12px;
        }
        .task-info-title {
            font-size: 1.3em;
            color: var(--accent);
            margin-bottom: 10px;
            font-weight: bold;
            letter-spacing: 1px;
        }
        .task-info-table {
            width: 100%;
            border-collapse: collapse;
            font-size: 1.08em;
            word-break: break-all;
        }
        .task-info-table th, .task-info-table td {
            padding: 7px 12px;
            text-align: left;
        }
        .task-info-table th {
            color: var(--accent);
            font-weight: 500;
            width: 120px;
            letter-spacing: 1px;
        }
        .task-info-table tr {
            border-bottom: 1px solid #ffffff11;
        }
        .stats-cards {
            display: flex;
            justify-content: center;
            gap: 32px;
            margin: 30px 0;
            flex-wrap: wrap;
        }
        .card {
            background: var(--card-bg);
            border-radius: 18px;
            box-shadow: 0 2px 16px var(--shadow);
            padding: 24px 36px;
            min-width: 220px;
            text-align: center;
            margin-bottom: 16px;
        }
        .card-title {
            font-size: 1.1em;
            color: var(--accent);
            margin-bottom: 8px;
            letter-spacing: 1px;
        }
        .card-value {
            font-size: 2em;
            font-weight: bold;
            margin-bottom: 4px;
        }
        .section {
            margin: 40px auto;
            max-width: 1100px;
            background: var(--section-bg);
            border-radius: 18px;
            box-shadow: 0 2px 24px var(--shadow);
            padding: 32px 40px;
        }
        .section h2 {
            color: var(--accent);
            font-size: 1.5em;
            margin-bottom: 18px;
            letter-spacing: 1px;
        }
        .chart-card {
            background: var(--card-bg);
            border-radius: 14px;
            box-shadow: 0 2px 12px var(--shadow);
            padding: 18px 10px 12px 10px;
            margin-bottom: 32px;
            max-width: 100%;
            width: 99%;
            margin-left: auto;
            margin-right: auto;
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .chart-title {
            color: var(--accent);
            font-size: 1.08em;
            margin-bottom: 8px;
            text-align: center;
            letter-spacing: 1px;
        }
        .chart-inner {
            width: 100%;
            max-width: 1380px;
            min-width: 900px;
            overflow-x: auto;
            overflow-y: hidden;
            display: flex;
            justify-content: flex-start;
            align-items: center;
            padding: 0;
        }
        .chart-inner > div, .chart-inner > #chart-container {
            width: 100% !important;
            min-width: 900px !important;
            max-width: 1380px !important;
        }
        .advice-list li {
            margin: 12px 0;
            font-size: 1.1em;
            background: linear-gradient(90deg, var(--accent)33 0%, #fff0 100%);
            border-left: 4px solid var(--accent);
            padding: 8px 16px;
            border-radius: 6px;
            box-shadow: 0 1px 6px var(--shadow);
        }
        @media (max-width: 900px) {
            .chart-card { max-width: 100%; }
            .task-info-card { padding: 18px 8vw 10px 8vw; }
            .chart-inner { max-width: 100%; }
        }
    </style>
    <script>
        function toggleTheme() {
            const html = document.documentElement;
            const theme = html.getAttribute('data-theme');
            if (theme === 'light') {
                html.setAttribute('data-theme', 'dark');
                localStorage.setItem('theme', 'dark');
            } else {
                html.setAttribute('data-theme', 'light');
                localStorage.setItem('theme', 'light');
            }
        }
        window.onload = function() {
            const saved = localStorage.getItem('theme');
            if(saved) document.documentElement.setAttribute('data-theme', saved);
        }
    </script>
</head>
<body>
    <button class="theme-toggle" onclick="toggleTheme()">{{index .Lang "ToggleTheme"}}</button>
    <div class="header">
        <h1>{{index .Lang "ReportTitle"}}</h1>
        <div class="score-ring">
            <svg width="140" height="140">
                <circle cx="70" cy="70" r="60" stroke="#222" stroke-width="16" fill="none"/>
                <circle cx="70" cy="70" r="60" stroke="{{.ScoreColor}}" stroke-width="16" fill="none"
                    stroke-dasharray="377" stroke-dashoffset="{{.ScoreOffset}}" />
            </svg>
            <div class="score-text">{{.Score}}</div>
        </div>
        <div style="font-size:1.2em;color:var(--text-main);margin-bottom:10px;">{{index .Lang "ScoreLabel"}}</div>
        <div style="font-size:1.1em;color:{{.ScoreColor}};margin-bottom:10px;">{{.Status}}</div>
    </div>
    <div class="task-info-card">
        <div class="task-info-title">{{index .Lang "TaskInfo"}}</div>
        <table class="task-info-table">
            <tr>
                <th>{{index .Lang "Targets"}}</th>
                <td style="word-break: break-all;">{{.Task.Targets | formatMultiline}}</td>
            </tr>
            <tr>
                <th>{{index .Lang "Usernames"}}</th>
                <td style="word-break: break-all;">{{.Task.Usernames | formatMultiline}}</td>
            </tr>
            <tr>
                <th>{{index .Lang "Passwords"}}</th>
                <td style="word-break: break-all;">{{.Task.Passwords | formatMultiline}}</td>
            </tr>
            <tr>
                <th>{{index .Lang "Proxies"}}</th>
                <td style="word-break: break-all;">{{.Task.Proxies | formatMultiline}}</td>
            </tr>
        </table>
    </div>
    <div class="stats-cards">
        <div class="card">
            <div class="card-title">{{index .Lang "TotalCrack"}}</div>
            <div class="card-value">{{.Total}}</div>
        </div>
        <div class="card">
            <div class="card-title">{{index .Lang "WeakCount"}}</div>
            <div class="card-value">{{.WeakCount}}</div>
        </div>
    </div>
    <div class="section">
        <h2>{{index .Lang "Analysis"}}</h2>
        <details class="chart-card" style="max-width:1400px;" open>
            <summary class="chart-title" style="cursor:pointer;">{{index .Lang "WeakDetail"}}{{index .Lang "Fold"}}</summary>
            <div style="overflow-x:auto;">
                <table style="width:100%;border-collapse:collapse;">
                    <thead>
                        <tr style="background:var(--section-bg);">
                            <th style="padding:8px 12px;">{{index .Lang "Target"}}</th>
                            <th style="padding:8px 12px;">{{index .Lang "Service"}}</th>
                            <th style="padding:8px 12px;">{{index .Lang "Username"}}</th>
                            <th style="padding:8px 12px;">{{index .Lang "Password"}}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .WeakResults}}
                        <tr>
                            <td style="padding:6px 12px;">{{.Target}}</td>
                            <td style="padding:6px 12px;">{{.Service}}</td>
                            <td style="padding:6px 12px;">{{.Username}}</td>
                            <td style="padding:6px 12px;">{{.Password}}</td>
                        </tr>
                        {{end}}
                        {{if not .WeakResults}}
                        <tr>
                            <td colspan="4" style="text-align:center;padding:12px;">{{index .Lang "NoWeak"}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </details>
        <details class="chart-card" open>
            <summary class="chart-title" style="cursor:pointer;">{{index .Lang "TargetCrack"}}{{index .Lang "Fold"}}</summary>
            <div class="chart-inner">
                {{.BarChart}}
            </div>
        </details>
        <details class="chart-card" open>
            <summary class="chart-title" style="cursor:pointer;">{{index .Lang "ServicePie"}}{{index .Lang "Fold"}}</summary>
            <div class="chart-inner">
                {{.PieChart}}
            </div>
        </details>
    </div>
    <div class="section">
        <h2>{{index .Lang "Advice"}}</h2>
        <ul class="advice-list">
            {{range .Advice}}<li>{{.}}</li>{{end}}
        </ul>
    </div>
</body>
</html>
`

	funcMap := template.FuncMap{
		"formatTime": func(ts int64) string {
			if ts == 0 {
				return "-"
			}
			return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
		},
		"formatMultiline": func(s string) template.HTML {
			s = strings.ReplaceAll(s, ",", "<br>")
			s = strings.ReplaceAll(s, "\n", "<br>")
			return template.HTML(s)
		},
	}
	scoreColor := "#00eaf"
	if score < 80 {
		scoreColor = "#ffb300"
	}
	if score < 60 {
		scoreColor = "#ff3b3b"
	}
	scoreOffset := 377 - int(float64(score)/100*377)
	data := map[string]interface{}{
		"Task":        task,
		"Score":       score,
		"ScoreColor":  scoreColor,
		"ScoreOffset": scoreOffset,
		"Status":      status,
		"Total":       total,
		"WeakCount":   weakPasswordCount,
		"BarChart":    template.HTML(barHTML),
		"PieChart":    template.HTML(pieHTML),
		"Advice":      advice,
		"WeakResults": weakResults,
		"Lang":        langs[lang],
	}
	tpl, err := template.New("report").Funcs(funcMap).Parse(reportTpl)
	if err != nil {
		return
	}
	var out bytes.Buffer
	tpl.Execute(&out, data)
	content = out.String()
	slog.Printf(slog.INFO, "安全评估报告已生成")
	ok = true
	return
}
