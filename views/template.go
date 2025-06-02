package views

// Layout layout_info
const Layout = `{{define "layout_info"}}

	<!DOCTYPE html>
	<html lang="zh-Hant-TW">

	{{ template "head" . }}

	<body id="body">

		{{ template "header" . }}

		<main id="full-content">

			{{ template "activity_content" . }}

		</main>

		{{if eq .User.Cookie "no"}}

			<aside id="cookie-panel">
				<div class="cookie-content">
					<p>為了提供您最佳的服務，本網站會在您的電腦中放置並取用我們的Cookie，若您不願接受Cookie的寫入，您可在您使用的瀏覽器功能項中設定隱私權等級為高，即可拒絕Cookie的寫入，但可能會導致網站某些功能無法正常執行。</p>
					<button class="button button-green">接受</button>
				</div>
			</aside>

		{{end}}

	</body>

	<script>

		headerFocus()
		$("#my-profile").bind("click", function(){
			$.get("{{.Route.User}}", function(data){
				$("#full-content").html(data)
			})
		})
		
		$("#my-party").bind("click", function(){
			$.get("{{.Route.Activity}}", function(data){
				$("#full-content").html(data)
			})
		})

		$("#cookie-panel").find(".button").bind("click", function(){
			event.stopPropagation()

			var formDataBox = new FormData()
			formDataBox.append("user", "{{.User.UserID}}")
			formDataBox.append("user_id", "{{.User.UserID}}")
			formDataBox.append("cookie", "yes")
			formDataBox.append("token", "{{.Token}}")

			var settings = formSet("https://"+ selectLocalhostAPI() +"/v1/user", "PUT", formDataBox)

			$.ajax(settings).done(function(){
				/*$.get("/admin/activity?header=true", function(data){
					$("#main-content").html(data)
				})*/
				$("#cookie-panel").addClass("accept")
			}).fail(function(response){
				console.log("資料傳遞錯誤，請聯絡客服", response)
			})

		})

		function selectLocalhostAPI(){
			if(location.host === "www.hilives.net" || location.host === "hilives.net"){
				return "api.hilives.net"
			}else if(location.host === "dev.hilives.net"){
				return "apidev.hilives.net"
			}
		}

	</script>

	</html>
	
{{end}}`

// Head head
const Head = `{{define "head"}}

	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>HiLives活動管理平台</title>
		<link rel="icon" type="image/x-icon" href="/admin/assets/website/lib/favicon.svg">

		<!-- Old Data -->
		<!--script src="/admin/assets/dist/js/default.js"></script-->
		<!--link rel="stylesheet" href="/admin/assets/dist/css/default.css"-->
		<!--script src="https://cdn.jsdelivr.net/npm/jquery-twzipcode@1.7.14/jquery.twzipcode.min.js"></script-->

		<script>

			

		</script>
		

		<script src="//cdn.jsdelivr.net/npm/promise-polyfill@8/dist/polyfill.js"></script>
		
		<!-- Jquery -->
		<script src="https://code.jquery.com/jquery-3.6.0.js" integrity="sha256-H+K7U5CnXl1h5ywQfKtSj8PCmoN9aaq30gDh27Xc0jk=" crossorigin="anonymous"></script>
		
		<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.form/4.3.0/jquery.form.min.js" integrity="sha512-YUkaLm+KJ5lQXDBdqBqk7EVhJAdxRnVdT2vtCzwPHSweCzyMgYV/tgGF4/dCyqtCC2eCphz0lRQgatGVdfR0ww==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>

		<!-- Jquery UI -->
		<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js" integrity="sha256-T0Vest3yCU7pafRw9r+settMBX6JkKN06dqBnpQ8d30=" crossorigin="anonymous"></script>
		
		<!-- Font Awesome -->
		
		<!-- Anime.js -->
		<script src="https://cdnjs.cloudflare.com/ajax/libs/animejs/3.2.0/anime.min.js"></script>
		
		<!-- Sweetalert2 -->
		<script src="//cdn.jsdelivr.net/npm/sweetalert2@11"></script>
		
		<!-- QRcode -->
		<!-- script src="https://unpkg.com/qr-code-styling@1.5.0/lib/qr-code-styling.js"></script -->
		<script src="/admin/assets/dist/js/qr-code-styling.js"></script>
		<script src="/admin/assets/dist/js/create-qrcode.js"></script>

		<!-- Taiwan zip code -->
    	<script src="/admin/assets/website/js/twzipcode.js"></script>
		
		<!-- Bootstrap -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>

		<!-- Bootstrap Icon -->
    	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.7.2/font/bootstrap-icons.css">

		<!-- Main Style -->
		<link href="/admin/assets/cms/css/cms.css?ver=20221114" rel="stylesheet">

		<script src="/admin/assets/website/js/index.min.js"></script>
		<script src="/admin/assets/cms/js/cms.js?ver=20230516" charset="big5"></script>
		
	</head>

{{end}}`

// Header header
const Header = `{{define "header"}}

	<header id="header-bar">
		<div id="brand-title">
            <a href="#首頁網址">
                <span>HiLives</span>
                <span style="color: #fafafa; font-weight: 500;">活動管理平台</span>
            </a>
        </div>
		<nav>
			<ul>
				{{if eq .IsAdmin true}}

					<li id="admin-page" class="link-button">
						<a href="{{.Route.Admin}}" role="button">
							<span>權限管理</span>
						</a>
					</li>

				{{end}}
				<li id="my-profile" class="link-button">
					<a role="button">

						{{if eq .User.Avatar ""}}

							<div class="svg-box">
								<img src="/admin/assets/website/lib/img/headpic.svg" alt="User Image">
							</div>

						{{else}}

							<div class="img-box">
								<img src="{{.User.Avatar}}" alt="User Avatar" width="100px" height="100px">
							</div>

						{{end}}
						
						<div class="username" title="{{.User.Name}}">{{.User.Name}}</div>
                	</a>
				</li>
				<li id="my-party" class="link-button focus">
					<a role="button">
                    	<span>我的活動</span>
                	</a>
				</li>
				<li id="log-out" class="link-button">
					<a href="{{.Route.Logout}}" role="button">
                    	<span>登出</span>
                	</a>
				</li>
			</ul>
		</nav>
	</header>

{{end}}`

const Sidebar = `{{define "sidebar"}}

	<nav id="sidebar">
		<ul class="menu-rank1">

			{{$ActivityID := .ActivityID }}
			{{range $key, $list := .Menu.List }}

				<li id={{$list.SidebarBtnL}} class="rank1-box">
					<div class="group-h">
						<div class="icon-block">
							<div class="icon-box">

								{{if eq $list.SidebarBtnL "info"}}

									<img src="/admin/assets/website/lib/img/sidebar-calendar.svg" alt="calendar">
								
								{{else if eq $list.SidebarBtnL "applysign"}}

									<img src="/admin/assets/website/lib/img/sidebar-list.svg" alt="list">
								
								{{else if eq $list.SidebarBtnL "interact"}}
									
									<img src="/admin/assets/website/lib/img/sidebar-gift.svg" alt="gift">

								{{else if eq $list.SidebarBtnL "staffmanage"}}

									<img src="/admin/assets/website/lib/img/sidebar-people.svg" alt="people">

								{{end}}
								
							</div>
						</div>
						<div class="right-box">
							<div class="group-h">
								<div class="link-title">{{$list.Name}}</div>
								<div class="arrow-block">
									<div class="arrow-box">
										<img src="/admin/assets/website/lib/img/sidebar-arrow.svg" alt="">
									</div>
								</div>
							</div>
						</div>
					</div>
					<ul class="menu-rank2">

						{{range $key2, $item := $list.ChildrenList}}
							{{if eq (len $item.ChildrenList) 0}}

								<li id={{$item.SidebarBtnL}} class="rank2-box">
									<a data-href="{{$item.URL}}?activity_id={{$ActivityID}}">
										<span>{{$item.Name}}</span>
									</a>
								</li>

							{{else}}

								<li id={{$item.SidebarBtnL}} class="rank2-box style-2">
									<div class="group-h">
										<span>{{$item.Name}}</span>
										<div class="arrow-block">
											<div class="arrow-box">
												<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 400"><path d="M296.26,217l-136,136a23.9,23.9,0,0,1-33.9,0l-22.6-22.6a23.9,23.9,0,0,1,0-33.9l96.4-96.4-96.4-96.4a23.9,23.9,0,0,1,0-33.9L126.26,47a23.9,23.9,0,0,1,33.9,0l136,136A23.93,23.93,0,0,1,296.26,217Z"/></svg>
											</div>
										</div>
									</div>
									<ul class="menu-rank3">

										{{range $key3, $subItem := $item.ChildrenList}}
											{{if eq (len $subItem.ChildrenList) 0}}

												{{if eq $subItem.Name "圖片播放區"}}

													<li class="rank3-box page-not-open">
														<a href="#not-open"><span>圖片播放區</span></a>
													</li>

												{{else if eq $subItem.Name "倒數簽到"}}

													<li class="rank3-box page-not-open">
														<a href="#not-open"><span>倒數簽到</span></a>
													</li>
														
												{{else}}

													<li id={{$subItem.SidebarBtnL}} class="rank3-box">
														<a data-href="{{$subItem.URL}}?activity_id={{$ActivityID}}" data-label-name="{{$subItem.SidebarBtnL}}">
															<span>{{$subItem.Name}}</span>
														</a>
													</li>
													
												{{end}}

											{{else}}

												<li class="rank3-box" style="justify-content: flex-start; cursor: none; pointer-events: none;">
													<span style="color: #767676 !important">--{{$subItem.Name}}--</span>
												</li>

												{{range $key4, $fourItem := $subItem.ChildrenList}}

													<li id={{$fourItem.SidebarBtnL}} class="rank3-box">
														<a data-href="{{$fourItem.URL}}?activity_id={{$ActivityID}}" data-label-name="{{$fourItem.SidebarBtnL}}">
															<span>{{$fourItem.Name}}</span>
														</a>
													</li>

												{{end}}

											{{end}}
										
										{{end}}
									
									</ul>
								</li>

							{{end}}
						{{end}}

					</ul>
				</li>

			{{end}}

		</ul>
	</nav>
	<script>
		dfgdgdgdgd()
	</script>

{{end}}`

const Form = `{{define "form"}}

	{{range $key, $data := .FormInfo.FieldList}}
		{{if $data.IsHide}}
			<input type="hidden" name="{{$data.Field}}" value="{{$data.Value}}">
		{{else}}

			{{if eq $data.FormType.String "text"}}

				<fieldset class="input-group-3">
					<label class="input-title" for="{{$data.Field}}">{{$data.Header}}</label>
					<input class="input form-control" id="{{$data.Field}}" type="{{$data.TypeName}}" name="{{$data.Field}}" value="{{$data.Value}}" maxlength="20" placeholder="{{(index $data.Placeholder 0)}}" autocomplete="off" required="{{$data.Must}}">
					{{if ne $data.HelpMsg ""}}
						<p class="help-message">{{$data.HelpMsg}}</p>
					{{end}}
				</fieldset>

			{{else if eq $data.FormType.String "select_single"}}

				<fieldset class="input-group-3">
					<label class="input-title" for="{{$data.Field}}">{{$data.Header}}</label>
					<select class="input input-select form-select" id="{{$data.Field}}" name="{{$data.Field}}" placeholder="{{$data.Placeholder}}" required="{{$data.Must}}">
						<option value="">請選擇活動類型</option>
						{{range $key, $v := .FieldOptions }}
							<option value='{{$v.Value}}' {{$v.SelectedLabel}}>{{$v.Text}}</option>
						{{end}}
					</select>
					{{if ne $data.HelpMsg ""}}
						<p class="help-message">{{$data.HelpMsg}}</p>
					{{end}}
				</fieldset>

			{{else if eq $data.FormType.String "number"}}

				<fieldset class="input-group-3">
					<label class="input-title" for="{{$data.Field}}">{{$data.Header}}</label>
					<input class="input form-control" id="{{$data.Field}}" type="{{$data.TypeName}}" name="{{$data.Field}}" value="{{$data.Value}}" min="0" max="10000" step="1" placeholder="{{(index $data.Placeholder 0)}}" required="{{$data.Must}}">
					{{if ne $data.HelpMsg ""}}
						<p class="help-message">{{$data.HelpMsg}}</p>
					{{end}}
				</fieldset>

			{{else if eq $data.FormType.String "select_city"}}

				<fieldset class="input-group-3">
					<label class="input-title" for="{{$data.Field}}">{{$data.Header}}</label>
					<div id="citytown-box">
						<div data-role="county" data-css="input-select form-select" data-name="city" data-label="{{(index $data.Placeholder 0)}}" {{if eq $data.Value ""}}{{else}} data-value="{{$data.Value}}" {{end}} style="margin-bottom: 0.313rem;"></div>
						<div data-role="district" data-css="input-select form-select" data-name="town" data-label="{{(index $data.Placeholder 1)}}" {{if eq $data.Value ""}}{{else}} data-value="{{$data.Value2}}" {{end}}></div>
					</div>
					<script>
						var twzipcode = new TWzipcode("#citytown-box")
					</script>
					{{if ne $data.HelpMsg ""}}
						<p class="help-message">{{$data.HelpMsg}}</p>
					{{end}}
				</fieldset>

			{{else if eq $data.FormType.String "datetime_range"}}

				<fieldset class="input-group-3">
					<label class="input-title" for="{{$data.Field}}">{{$data.Header}}</label>
					<input type="datetime-local" name="{{$data.Field}}" value="{{$data.Value}}" class="input input-localtime form-control" placeholder="{{(index $data.Placeholder 0)}}" required="{{$data.Must}}">
					<input type="datetime-local" name="end_time" value="{{$data.Value2}}" class="input input-localtime form-control" placeholder="{{(index $data.Placeholder 1)}}" required="{{$data.Must}}">
					{{if ne $data.HelpMsg ""}}
						<p class="help-message">{{$data.HelpMsg}}</p>
					{{end}}
				</fieldset>

			{{end}}

		{{end}}
	{{end}}
	<input type="hidden" name="token" value="{{.Token}}">
{{end}}`

//sidebar可以去掉吧
//form等個人資訊和活動表單有更新也可以去掉吧
//其他寫固定的