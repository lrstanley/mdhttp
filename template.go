// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mdhttp

const htmlTemplate = `<!DOCTYPE html>
<html>
	<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta charset="UTF-8">
	<title>{{.Title}} &middot; {{ .FileInfo.Name }}</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha256-LA89z+k9fjgMKQ/kq4OO2Mrf8VltYml/VES+Rg0fh20=" crossorigin="anonymous" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.9.0/github-markdown.min.css" integrity="sha256-AcEoN20lbehGx2iA6fVk4geoblRBOAUmz+gEM9QE59Y=" crossorigin="anonymous" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css" integrity="sha256-eZrrJcwDc/3uDhsdt61sL2oOBY362qM3lon1gyExkL0=" crossorigin="anonymous" />
	<style>
		html, body {
			font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol";
		}
		body { position: relative; }

		.main {
			margin-top: 40px;
			margin-bottom: 40px;
		}

		.markdown-title {
			display: block;
			padding: 9px 10px 10px;
			margin: 0;
			font-size: 14px;
			font-weight: 600;
			line-height: 17px;
			background-color: #f6f8fa;
			border: 1px solid rgba(27,31,35,0.15);
			border-bottom: 0;
			border-radius: 3px 3px 0 0;
		}

		.markdown-title > span { margin-right: 10px; }

		.markdown-body {
			box-sizing: border-box;
			padding: 20px 45px;
			word-wrap: break-word;
			background-color: #FFFFFF;
			border: 1px solid #DDDDDD;
			border-bottom-right-radius: 3px;
			border-bottom-left-radius: 3px;
		}

		@media (max-width: 767px) {
			.markdown-body { padding: 15px; }
		}

		h1 > .anchor, h2 > .anchor, h3 > .anchor, h4 > .anchor {
			font-size: 20px;
			position: relative;
			text-decoration: none;
			color: #639fe2;
		}

		h1 > .anchor:hover, h2 > .anchor:hover, h3 > .anchor:hover, h4 > .anchor:hover { color: #0366d6; }
		h1 > .anchor { top: 10px; }
		h2 > .anchor { top: 5px; }
		h3 > .anchor { top: 2.5px; }
		h4 > .anchor { font-size: 18px; }

		.sidebar { top: 40px; }
		.sidebar .sidebar-content {
			box-sizing: border-box;
			padding: 15px;
			word-wrap: break-word;
			background-color: #FFFFFF;
			border: 1px solid #DDDDDD;
			border-bottom-right-radius: 3px;
			border-bottom-left-radius: 3px;
		}

		.html-gen-toc > ul { position: relative; }
		.html-gen-toc li { padding: 0 0 0.3rem 0.4rem; }
		.html-gen-toc li:last-child { padding: 0 0 0rem 0.4rem; }
		.html-gen-toc > ul:last-child { padding: 0 0 0.6rem 0; }

		.html-gen-toc li > a {
			display: block;
			text-decoration: none;
			font-size: 0.9rem;
			color: #99979c;
			padding: 0;
		}

		.html-gen-toc li > a:before {
			content: '· ';
			font-weight: 700;
			color: #000000;
			margin-left: -0.3rem;
		}

		.html-gen-toc li > a.active {
			color: #000000;
			border-left: 5px solid rgb(0, 93, 214);
			padding: 0 0 0 20px;
			margin-left: -25px;
		}

		.html-gen-toc li > a:hover {
			font-size: 0.9rem;
			color: #007bff;
			text-decoration: none;
		}

		.footer {
			font-size: 13px;
			color: grey;
			margin: 10px 0;
		}
	</style>
</head>

<body data-spy="scroll" data-target="#scrollspy-parent" data-offset="40">
	<div class="container main">
		{{ $toc := .TOC }}
		<div class="row">
			<div class="col-sm-12{{ if (ne $toc "") }} col-lg-9{{ end }}">
				<h3 class="markdown-title">
					<span><i class="fa fa-book"></i> {{.Title}}</span>
					<span><i class="fa fa-file-code-o"></i> {{ .FileInfo.Name }}</span>
					<span><i class="fa fa-clock-o"></i> {{ .FileInfo.ModTime.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</span>
				</h3>
				<div class="markdown-body">{{ .Body }}</div>
			</div>

			{{ if (ne $toc "") }}
				<div class="col-sm-12 col-lg-3 ml-auto">
					<div class="sticky-top sidebar">
						<h3 class="markdown-title"><i class="fa fa-list-ul"></i> Table of Contents</h3>

						<div class="sidebar-content html-gen-toc flex-column pb-0" id="scrollspy-parent">{{ $toc }}</div>
					</div>
				</div>
			{{ end }}
		</div>

		<footer class="footer">
			Generated with <a target="_blank" href="https://github.com/russross/blackfriday">blackfriday</a> &middot;
			Served with <a target="_blank" href="https://github.com/lrstanley/mdhttp">mdhttp</a>
	</div>

	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.2.4/jquery.min.js" integrity="sha256-BbhdlvQf/xTY9gja0Dq3HiwQF8LaCRTXxZKRutelT44=" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha256-5+02zu5UULQkO7w1GIr6vftCgMfFdZcAHeDtFnKZsBs=" crossorigin="anonymous"></script>

	<script type="text/javascript">
		$(function () {
			$.each($('.markdown-body h1, .markdown-body h2, .markdown-body h3, .markdown-body h4, .markdown-body h5, .markdown-body h6'), function(index, value) {
				var anchor = $(this).attr("id");
				$(this).html("<a class='anchor' href='#" + anchor + "'>#</a> " + $(this).html());
			});

			setTimeout(function() {
				if (window.location.hash != "") {
					el = document.getElementById(window.location.hash.replace("#", ""));
					if (el !== null) { el.scrollIntoView(); }
				}
			}, 500);
		});
	</script>
</body>
</html>`
