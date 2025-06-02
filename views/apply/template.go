package apply

const Template = `{{define "template"}}

<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">

		<!-- Page Name -->
		<title>HiLive活動報名審核</title>

		<!-- Main Style -->
		<link href="/admin/assets/verification/css/audit.css" rel="stylesheet">

		<!-- Main Script -->
		<!--script src="#"></script-->

	</head>
	<body>

		{{ template "content" . }}

	</body>
</html>

{{end}}`

// Content content
const Content = `{{define "content"}}

	<div class="audit-page">
		<div class="audit-block">
			<div class="audit-content">

				{{if eq .Status "apply"}}
	
					<div class="audit-icon-block">
						<div class="icon-block">
							<img src="/admin/assets/verification/lib/img/audit-success.svg" alt="a">
						</div>
					</div>
					<div class="audit-result-block audit-success">
						<p class="audit-result">已完成活動審核</p>
					</div>
					<div class="audit-message-block">
						<p class="audit-message">活動尚未開始，敬請期待</p>
					</div>

				{{else if eq .Status "refuse"}}

					<div class="audit-icon-block">
						<div class="icon-block">
							<img src="/admin/assets/verification/lib/img/audit-error.svg" alt="a">
						</div>
					</div>
					<div class="audit-result-block audit-refuse">
						<p class="audit-result">未通過活動審核</p>
					</div>
					<div class="audit-message-block">
						<p class="audit-message">如果有疑問，請聯繫主辦單位</p>
					</div>

				{{else if eq .Status "review"}}

					<div class="audit-icon-block">
						<div class="icon-block">
							<img src="/admin/assets/verification/lib/img/audit-await.svg" alt="a">
						</div>
					</div>
					<div class="audit-result-block audit-wait">
						<p class="audit-result">請等待活動審核</p>
					</div>
					<div class="audit-message-block">
						<p class="audit-message">已提交申請，請等待審核</p>
					</div>

				{{else if eq .Status "sign"}}

					<div class="audit-icon-block">
						<div class="icon-block">
							<img src="/admin/assets/verification/lib/img/audit-success.svg" alt="a">
						</div>
					</div>
					<div class="audit-result-block audit-wait">
						<p class="audit-result">您已完成簽到</p>
					</div>
					<div class="audit-message-block">
						<p class="audit-message">活動尚未開始，敬請期待</p>
					</div>

				{{end}}
	
			</div>
		</div>
	</div>

{{end}}`
